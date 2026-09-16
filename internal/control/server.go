package control

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	cookieName         = "wbcc_session"
	loginAttemptLimit  = 10
	loginAttemptWindow = 15 * time.Minute

	// loginPruneInterval 整表清理登录失败记录的最小间隔，见 pruneAttemptsLocked。
	// 用它限流是为了避免「表很大但记录都还没过期」时把每次登录都放大成一次 O(n) 扫描。
	loginPruneInterval = time.Minute
)

type session struct{ expires time.Time }
type loginAttempt struct {
	count int
	first time.Time
}
type Server struct {
	svc      *Service
	static   http.Handler
	mu       sync.Mutex
	sessions map[string]session
	attempts map[string]loginAttempt

	// lastPrune 记录上次整表清理 attempts 的时刻（见 pruneAttemptsLocked），由 mu 保护。
	lastPrune time.Time
}

func NewServer(svc *Service, static http.Handler) *Server {
	return &Server{svc: svc, static: static, sessions: map[string]session{}, attempts: map[string]loginAttempt{}}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/session", s.session)
	mux.HandleFunc("POST /api/login", s.login)
	mux.HandleFunc("POST /api/logout", s.logout)
	mux.HandleFunc("GET /api/overview", s.overview)
	mux.HandleFunc("GET /api/accounts", s.accounts)
	mux.HandleFunc("GET /api/credits", s.credits)
	mux.HandleFunc("GET /api/credits/{uid}", s.creditDetail)
	mux.HandleFunc("GET /api/models", s.models)
	mux.HandleFunc("GET /api/activities", s.activities)
	mux.HandleFunc("GET /api/scheduler-tasks", s.schedulerTasks)
	mux.HandleFunc("GET /api/scheduler-tasks/{uid}/{taskID}", s.schedulerTaskDetail)
	mux.HandleFunc("GET /api/mock-scheduler-tasks", s.mockSchedulerTasks)
	mux.HandleFunc("POST /api/mock-scheduler-tasks", s.mockSchedulerCreate)
	mux.HandleFunc("PUT /api/mock-scheduler-tasks/{id}", s.mockSchedulerUpdate)
	mux.HandleFunc("DELETE /api/mock-scheduler-tasks/{id}", s.mockSchedulerDelete)
	mux.HandleFunc("GET /api/automations", s.automations)
	mux.HandleFunc("PUT /api/automations/{id}", s.automationPut)
	mux.HandleFunc("POST /api/automations/{id}/run", s.automationRun)
	mux.HandleFunc("POST /api/accounts/{uid}/actions/{action}", s.accountAction)
	mux.HandleFunc("POST /api/oauth/start", s.oauthStart)
	mux.HandleFunc("POST /api/oauth/{id}/poll", s.oauthPoll)
	mux.Handle("/", s.static)
	return s.headers(s.auth(mux))
}
func (s *Server) headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; base-uri 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; connect-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:")
		next.ServeHTTP(w, r)
	})
}
func (s *Server) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api/login" || r.URL.Path == "/api/session" {
			next.ServeHTTP(w, r)
			return
		}
		if !s.valid(r) {
			jsonErr(w, http.StatusUnauthorized, "未登录或会话已过期")
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			origin := r.Header.Get("Origin")
			if origin == "" || !sameOrigin(origin, r.Host) {
				jsonErr(w, http.StatusForbidden, "跨站请求被拒绝")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
func sameOrigin(origin, host string) bool {
	origin = strings.TrimSuffix(origin, "/")
	return strings.EqualFold(origin, "http://"+host) || strings.EqualFold(origin, "https://"+host)
}
func (s *Server) valid(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.sessions[c.Value]
	if !ok || time.Now().After(v.expires) {
		delete(s.sessions, c.Value)
		return false
	}
	return true
}
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
func source(r *http.Request) string {
	h, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return h
	}
	return r.RemoteAddr
}
func (s *Server) session(w http.ResponseWriter, r *http.Request) {
	cfg := s.svc.Config()
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": s.valid(r), "read_only": cfg.ReadOnly, "using_default_password": cfg.Password == "workbuddy", "timezone": cfg.Timezone})
}

// pruneAttemptsLocked 回收已过窗口的登录失败记录（调用方须持有 s.mu）。
//
// 为什么需要：attempts 的条目只在「同一来源在窗口结束后再次尝试」或「该来源登录
// 成功」时才被删除（见 login）。于是一个**只失败一次、再也不回来**的来源会留下一条
// 永久记录，而来源数量并不受控——IPv6 下一个 /64 就能造出天文数字的不同地址，每个
// 地址只发一个请求即可，等于给未认证的远端留了一条无界增长服务端内存的路径。
// 这些过期记录也再不会被用到：只有该来源自己回来时才会被读到，届时同样按过期处理。
//
// 用时间间隔限流而不是「每次登录都扫」：表很大且记录都还没过期时，逐次整表扫描会被
// 放大成 O(n)/请求；限流后清理开销摊薄到每分钟至多一次，单次成本仍是一次 map 遍历。
func (s *Server) pruneAttemptsLocked(now time.Time) {
	if now.Sub(s.lastPrune) < loginPruneInterval {
		return
	}
	s.lastPrune = now
	for ip, rec := range s.attempts {
		if now.Sub(rec.first) >= loginAttemptWindow {
			delete(s.attempts, ip)
		}
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if !decode(w, r, &in) {
		return
	}
	key := source(r)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneAttemptsLocked(time.Now())
	attempt := s.attempts[key]
	if !attempt.first.IsZero() && time.Since(attempt.first) >= loginAttemptWindow {
		delete(s.attempts, key)
		attempt = loginAttempt{}
	}
	if attempt.count >= loginAttemptLimit {
		jsonErr(w, http.StatusTooManyRequests, "登录失败次数过多，请在 15 分钟后再试")
		return
	}
	cfg := s.svc.Config()
	if subtle.ConstantTimeCompare([]byte(in.Username), []byte(cfg.Username)) != 1 || subtle.ConstantTimeCompare([]byte(in.Password), []byte(cfg.Password)) != 1 {
		if attempt.first.IsZero() {
			attempt.first = time.Now()
		}
		attempt.count++
		s.attempts[key] = attempt
		jsonErr(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	token, err := randomToken()
	if err != nil {
		jsonErr(w, 500, "创建会话失败")
		return
	}
	s.sessions[token] = session{expires: time.Now().Add(12 * time.Hour)}
	delete(s.attempts, key)
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, MaxAge: 43200})
	writeJSON(w, 200, map[string]any{"ok": true})
}
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(cookieName); err == nil {
		s.mu.Lock()
		delete(s.sessions, c.Value)
		s.mu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, MaxAge: -1})
	writeJSON(w, 200, map[string]any{"ok": true})
}
func (s *Server) overview(w http.ResponseWriter, r *http.Request) {
	accounts, warns := s.svc.Accounts()
	active := 0
	for _, a := range accounts {
		if !a.Expired {
			active++
		}
	}
	autos := s.svc.State().Automations()
	enabled := 0
	for _, a := range autos {
		if a.Enabled {
			enabled++
		}
	}
	writeJSON(w, 200, map[string]any{"accounts": len(accounts), "healthy": active, "expired": len(accounts) - active, "automations": enabled, "warnings": warns, "updated_at": time.Now()})
}
func (s *Server) accounts(w http.ResponseWriter, r *http.Request) {
	items, warns := s.svc.Accounts()
	writeJSON(w, 200, map[string]any{"accounts": items, "warnings": warns})
}
func (s *Server) credits(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.svc.Credits(r.Context(), "")})
}
func (s *Server) creditDetail(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.svc.Credits(r.Context(), r.PathValue("uid"))})
}
func (s *Server) models(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.svc.Models(r.Context())})
}
func (s *Server) activities(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.svc.Activities(r.Context())})
}
func (s *Server) schedulerTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.svc.SchedulerTasks(r.Context()), "write_enabled": false, "note": "这是各 WorkBuddy 账号的云端定时任务。创建接口请求体尚未经过真实认证验证，当前保持只读。"})
}
func (s *Server) schedulerTaskDetail(w http.ResponseWriter, r *http.Request) {
	item, err := s.svc.SchedulerTaskDetail(r.Context(), r.PathValue("uid"), r.PathValue("taskID"))
	if err != nil {
		jsonErr(w, http.StatusBadGateway, "账号云端定时任务详情读取不可用")
		return
	}
	writeJSON(w, 200, map[string]any{"item": redactTaskDetail(item)})
}
func (s *Server) mockSchedulerTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.svc.MockSchedulerTasks(), "mode": "mock", "note": "本地 Mock 骨架：仅写入控制台 data/state.json，不会提交到 WorkBuddy 上游。"})
}

type mockSchedulerInput struct {
	AccountUID string `json:"account_uid"`
	Name       string `json:"name"`
	Cron       string `json:"cron"`
	Prompt     string `json:"prompt"`
	Enabled    bool   `json:"enabled"`
}

func (s *Server) mockSchedulerCreate(w http.ResponseWriter, r *http.Request) {
	if s.svc.Config().ReadOnly {
		jsonErr(w, 403, "服务端已开启只读模式")
		return
	}
	var in mockSchedulerInput
	if !decode(w, r, &in) {
		return
	}
	task, err := s.svc.CreateMockSchedulerTask(in.AccountUID, in.Name, in.Cron, in.Prompt, in.Enabled)
	if err != nil {
		jsonErr(w, 400, err.Error())
		return
	}
	writeJSON(w, 201, task)
}
func (s *Server) mockSchedulerUpdate(w http.ResponseWriter, r *http.Request) {
	if s.svc.Config().ReadOnly {
		jsonErr(w, 403, "服务端已开启只读模式")
		return
	}
	var in mockSchedulerInput
	if !decode(w, r, &in) {
		return
	}
	task, err := s.svc.UpdateMockSchedulerTask(r.PathValue("id"), in.Name, in.Cron, in.Prompt, in.Enabled)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			jsonErr(w, 404, "Mock 任务不存在")
		} else {
			jsonErr(w, 400, err.Error())
		}
		return
	}
	writeJSON(w, 200, task)
}
func (s *Server) mockSchedulerDelete(w http.ResponseWriter, r *http.Request) {
	if s.svc.Config().ReadOnly {
		jsonErr(w, 403, "服务端已开启只读模式")
		return
	}
	if err := s.svc.State().DeleteMockSchedulerTask(r.PathValue("id")); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			jsonErr(w, 404, "Mock 任务不存在")
		} else {
			jsonErr(w, 500, "删除 Mock 任务失败")
		}
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true})
}

// redactTaskDetail ensures a future upstream response cannot accidentally turn
// this management view into a credential viewer.
func redactTaskDetail(item map[string]any) map[string]any {
	out := make(map[string]any, len(item))
	for k, v := range item {
		key := strings.ToLower(k)
		if strings.Contains(key, "token") || strings.Contains(key, "authorization") || strings.Contains(key, "cookie") || strings.Contains(key, "secret") || strings.Contains(key, "password") {
			continue
		}
		switch child := v.(type) {
		case map[string]any:
			out[k] = redactTaskDetail(child)
		case []any:
			if len(child) > 100 {
				out[k] = append([]any(nil), child[:100]...)
			} else {
				out[k] = child
			}
		default:
			out[k] = v
		}
	}
	return out
}
func (s *Server) automations(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"items": s.svc.State().Automations(), "runs": s.svc.State().Runs()})
}
func (s *Server) automationPut(w http.ResponseWriter, r *http.Request) {
	if s.svc.Config().ReadOnly {
		jsonErr(w, 403, "服务端已开启只读模式")
		return
	}
	var in struct {
		Enabled     bool `json:"enabled"`
		EveryMinute int  `json:"every_minutes"`
	}
	if !decode(w, r, &in) {
		return
	}
	a, err := s.svc.State().SetAutomation(r.PathValue("id"), in.Enabled, in.EveryMinute)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			jsonErr(w, 400, err.Error())
			return
		}
		jsonErr(w, 404, "自动化任务不存在")
		return
	}
	writeJSON(w, 200, a)
}
func (s *Server) automationRun(w http.ResponseWriter, r *http.Request) {
	if s.svc.Config().ReadOnly {
		jsonErr(w, http.StatusForbidden, "服务端已开启只读模式")
		return
	}
	id := r.PathValue("id")
	var action string
	for _, a := range s.svc.State().Automations() {
		if a.ID == id {
			action = a.Action
		}
	}
	if action == "" {
		jsonErr(w, 404, "自动化任务不存在")
		return
	}
	msg, err := s.svc.Run(r.Context(), action, "")
	if err != nil {
		_ = s.svc.State().Complete(id, err.Error(), false)
		jsonErr(w, 502, err.Error())
		return
	}
	_ = s.svc.State().Complete(id, msg, true)
	writeJSON(w, 200, map[string]any{"ok": true, "message": msg})
}
func (s *Server) accountAction(w http.ResponseWriter, r *http.Request) {
	if s.svc.Config().ReadOnly {
		jsonErr(w, http.StatusForbidden, "服务端已开启只读模式")
		return
	}
	action := r.PathValue("action")
	if action != "checkin" && action != "travel" && action != "refresh" {
		jsonErr(w, 400, "不支持的动作")
		return
	}
	msg, err := s.svc.Run(r.Context(), action, r.PathValue("uid"))
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "message": msg})
}
func (s *Server) oauthStart(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Region string `json:"region"`
	}
	if !decode(w, r, &in) {
		return
	}
	id, url, err := s.svc.StartOAuth(in.Region)
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"id": id, "url": url})
}
func (s *Server) oauthPoll(w http.ResponseWriter, r *http.Request) {
	res, err := s.svc.PollOAuth(r.PathValue("id"))
	if err != nil {
		jsonErr(w, 502, err.Error())
		return
	}
	writeJSON(w, 200, res)
}
func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Body == nil {
		jsonErr(w, 400, "缺少请求体")
		return false
	}
	defer r.Body.Close()
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(dst); err != nil {
		jsonErr(w, 400, "请求格式错误")
		return false
	}
	return true
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func jsonErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}
