package control

import (
	"regexp"
	"strings"
	"sync"
	"time"
)

// logLine 是 GET /api/logs 的单行（契约 §2.4 + M5 action 维度）。
type logLine struct {
	At      string `json:"at"`
	Level   string `json:"level"`
	Action  string `json:"action"`
	Message string `json:"message"`
}

// logBuffer 面板运行日志环形缓冲（M5）。
//
// 容量固定，写满后最旧记录被覆盖；snapshot 按新→旧返回。
// 所有写入在落缓冲前经 redactLog 脱敏（token/cookie/password 类键值）。
type logBuffer struct {
	mu     sync.Mutex
	lines  []logLine
	head   int // 下一次写入位置
	count  int // 已写入条数（<= cap）
	loc    *time.Location
}

// newLogBuffer 创建容量为 cap 的环形缓冲；timezone 无效时回退 Local。
func newLogBuffer(capacity int, timezone string) *logBuffer {
	if capacity <= 0 {
		capacity = 500
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.Local
	}
	return &logBuffer{lines: make([]logLine, capacity), loc: loc}
}

// add 追加一条日志；level 非法值归一为 info。
func (b *logBuffer) add(level, action, message string) {
	switch level {
	case "info", "warn", "error":
	default:
		level = "info"
	}
	line := logLine{
		At:      time.Now().In(b.loc).Format(time.RFC3339),
		Level:   level,
		Action:  action,
		Message: redactLog(message),
	}
	b.mu.Lock()
	b.lines[b.head] = line
	b.head = (b.head + 1) % len(b.lines)
	if b.count < len(b.lines) {
		b.count++
	}
	b.mu.Unlock()
}

// snapshot 按新→旧返回最近 limit 条；limit<=0 或超过存量时返回全部。
func (b *logBuffer) snapshot(limit int) []logLine {
	b.mu.Lock()
	defer b.mu.Unlock()
	n := b.count
	if limit > 0 && limit < n {
		n = limit
	}
	out := make([]logLine, n)
	for i := 0; i < n; i++ {
		// head-1 是最新一条，向前回绕。
		idx := (b.head - 1 - i + len(b.lines)) % len(b.lines)
		out[i] = b.lines[idx]
	}
	return out
}

// logSecretPatterns 日志脱敏规则：键值形态的敏感字段值替换为 ***。
var logSecretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(token|cookie|password|secret|authorization)("?\s*[:=]\s*"?)[^"\s&]+`),
}

// redactLog 在日志写入缓冲前脱敏，避免凭据类内容进入面板日志面。
func redactLog(msg string) string {
	for _, re := range logSecretPatterns {
		msg = re.ReplaceAllString(msg, `${1}${2}***`)
	}
	return strings.TrimSpace(msg)
}
