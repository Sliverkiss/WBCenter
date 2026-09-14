package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// TestLogBufferWriteAndRead 写入日志后按新→旧顺序读回，字段完整。
func TestLogBufferWriteAndRead(t *testing.T) {
	lb := newLogBuffer(10, "Asia/Shanghai")
	lb.add("info", "oauth", "账号 u-1001 登录成功")
	lb.add("error", "checkin", "账号 u-1002 签到失败")

	lines := lb.snapshot(100)
	if len(lines) != 2 {
		t.Fatalf("len(lines) = %d, want 2", len(lines))
	}
	// 新→旧：后写的 error 在前。
	if lines[0].Level != "error" || lines[1].Level != "info" {
		t.Fatalf("顺序应为新→旧: %+v", lines)
	}
	if lines[0].Action != "checkin" || lines[0].Message != "账号 u-1002 签到失败" {
		t.Fatalf("字段缺失: %+v", lines[0])
	}
	if lines[0].At == "" {
		t.Fatal("at 时间戳为空")
	}
	// at 应按配置时区渲染（Asia/Shanghai = +08:00）。
	if !strings.Contains(lines[0].At, "+08:00") {
		t.Fatalf("at 未按配置时区渲染: %q", lines[0].At)
	}
}

// TestLogBufferRingOverwrite 写入超过容量后最旧的被覆盖，长度恒定。
func TestLogBufferRingOverwrite(t *testing.T) {
	lb := newLogBuffer(3, "UTC")
	for i := 0; i < 10; i++ {
		lb.add("info", "test", fmt.Sprintf("msg-%02d", i))
	}
	lines := lb.snapshot(100)
	if len(lines) != 3 {
		t.Fatalf("环形缓冲长度 = %d, want 3", len(lines))
	}
	// 最新的三条是 msg-09/08/07（新→旧），最旧的 msg-00..06 已被覆盖。
	if lines[0].Message != "msg-09" || lines[1].Message != "msg-08" || lines[2].Message != "msg-07" {
		t.Fatalf("环形覆盖错误: %+v", lines)
	}
}

// TestLogBufferLimit limit 参数裁剪返回条数；0/负数/超容量按容量兜底。
func TestLogBufferLimit(t *testing.T) {
	lb := newLogBuffer(5, "UTC")
	for i := 0; i < 4; i++ {
		lb.add("info", "test", "msg")
	}
	if got := lb.snapshot(2); len(got) != 2 {
		t.Fatalf("limit=2 返回 %d 条", len(got))
	}
	if got := lb.snapshot(0); len(got) != 4 {
		t.Fatalf("limit=0 应返回全部 %d 条，实际 %d", 4, len(got))
	}
	if got := lb.snapshot(99); len(got) != 4 {
		t.Fatalf("limit 超存量应返回全部，实际 %d", len(got))
	}
}

// TestLogBufferConcurrent 并发写读不竞争（-race 验证）。
func TestLogBufferConcurrent(t *testing.T) {
	lb := newLogBuffer(50, "UTC")
	done := make(chan struct{})
	for w := 0; w < 4; w++ {
		go func(id int) {
			for i := 0; i < 200; i++ {
				lb.add("info", "test", "并发写入")
			}
			done <- struct{}{}
		}(w)
	}
	for r := 0; r < 2; r++ {
		go func() {
			for i := 0; i < 200; i++ {
				_ = lb.snapshot(10)
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 6; i++ {
		<-done
	}
	if got := lb.snapshot(100); len(got) != 50 {
		t.Fatalf("容量 50 缓冲应有 50 条，实际 %d", len(got))
	}
}

// TestLogsHandlerShape GET /api/logs 契约形状：lines[] + limit 查询参数。
func TestLogsHandlerShape(t *testing.T) {
	srv := newTestHTTPServer(t, false, false)
	defer srv.Close()
	cookie := loginForTest(t, srv.URL)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/api/logs?limit=100", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("GET /api/logs 状态码 = %d", resp.StatusCode)
	}
	var body struct {
		Lines []struct {
			At      string `json:"at"`
			Level   string `json:"level"`
			Action  string `json:"action"`
			Message string `json:"message"`
		} `json:"lines"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Lines == nil {
		t.Fatal("lines 字段缺失（应为数组，允许为空）")
	}
	// 登录事件应已被记录。
	found := false
	for _, l := range body.Lines {
		if l.Action == "login" {
			found = true
			if l.Level != "info" || l.Message == "" || l.At == "" {
				t.Fatalf("登录日志字段不完整: %+v", l)
			}
		}
	}
	if !found {
		t.Fatalf("缺少登录事件日志: %+v", body.Lines)
	}
}

// TestLogsHandlerLimitParam ?limit= 裁剪条数；非法值按默认 100。
func TestLogsHandlerLimitParam(t *testing.T) {
	srv := newTestHTTPServer(t, false, false)
	defer srv.Close()
	cookie := loginForTest(t, srv.URL)

	get := func(url string) []map[string]any {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		req.AddCookie(cookie)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var body struct {
			Lines []map[string]any `json:"lines"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		return body.Lines
	}

	// 制造 5 条日志（重复登录失败/成功混合）。
	for i := 0; i < 4; i++ {
		loginForTest(t, srv.URL)
	}
	all := get(srv.URL + "/api/logs")
	if len(all) < 5 {
		t.Fatalf("应至少有 5 条登录日志，实际 %d", len(all))
	}
	limited := get(srv.URL + "/api/logs?limit=2")
	if len(limited) != 2 {
		t.Fatalf("limit=2 应返回 2 条，实际 %d", len(limited))
	}
	// 新→旧：limit 后的第一条与全量第一条相同。
	if len(all) > 0 && len(limited) > 0 && all[0]["message"] != limited[0]["message"] {
		t.Fatalf("limit 裁剪应保持新→旧顺序: %+v vs %+v", all[0], limited[0])
	}
	// 非法 limit 不报 400，按默认处理。
	if bad := get(srv.URL + "/api/logs?limit=abc"); bad == nil {
		t.Fatal("非法 limit 应回退默认而非报错")
	}
}

// TestLogsHandlerUnauthorized 未登录访问 /api/logs 应 401。
func TestLogsHandlerUnauthorized(t *testing.T) {
	srv := newTestHTTPServer(t, false, false)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/api/logs")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("未登录 GET /api/logs 状态码 = %d, want 401", resp.StatusCode)
	}
}

// TestLogsRedactSecrets 写入缓冲前脱敏：token/cookie/密码类内容不落盘。
func TestLogsRedactSecrets(t *testing.T) {
	lb := newLogBuffer(10, "UTC")
	lb.add("info", "oauth", "token=abc123 cookie=wbcc_session=xyz password=hunter2 正常内容")
	lines := lb.snapshot(10)
	if len(lines) != 1 {
		t.Fatalf("lines = %d", len(lines))
	}
	msg := lines[0].Message
	for _, secret := range []string{"abc123", "xyz", "hunter2"} {
		if strings.Contains(msg, secret) {
			t.Fatalf("日志未脱敏，泄露 %q: %q", secret, msg)
		}
	}
	if !strings.Contains(msg, "正常内容") {
		t.Fatalf("脱敏不应破坏正常内容: %q", msg)
	}
}

// loginForTest 完成一次成功登录并返回会话 cookie（复用 integration helper 模式）。
// 注意：成功登录会向日志缓冲写入一条 action=login 记录，handler 测试依赖此行为。
func loginForTest(t *testing.T, base string) *http.Cookie {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, base+"/api/login", strings.NewReader(`{"username":"admin","password":"test-password"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("测试登录失败: %d", resp.StatusCode)
	}
	for _, c := range resp.Cookies() {
		if c.Name == cookieName {
			return c
		}
	}
	t.Fatal("登录响应缺少会话 cookie")
	return nil
}
