package control

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// M5 安全硬化回归：限速/同源/会话/Cookie 属性四类场景。
// 既有覆盖（server_integration_test.go）：401 未登录读、无 Origin 写 403、
// HttpOnly+SameSite=Strict、只读模式写 403。本文件补齐缺口并防回归。

// TestSecurityLoginRateLimit 同一来源连续失败登录超过上限后返回 429；
// 成功登录后计数清零（429 期间即使密码正确也拒绝）。
func TestSecurityLoginRateLimit(t *testing.T) {
	ts := newTestHTTPServer(t, false, false)
	defer ts.Close()

	bad := map[string]string{"username": "admin", "password": "wrong-password"}
	for i := 0; i < loginAttemptLimit; i++ {
		res := request(t, ts.Client(), http.MethodPost, ts.URL+"/api/login", bad, nil, false)
		res.Body.Close()
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("第 %d 次失败登录状态码 = %d, want 401", i+1, res.StatusCode)
		}
	}
	// 第 11 次（达到上限后）即使密码正确也应 429。
	res := request(t, ts.Client(), http.MethodPost, ts.URL+"/api/login", map[string]string{"username": "admin", "password": "test-password"}, nil, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusTooManyRequests {
		body, _ := io.ReadAll(res.Body)
		t.Fatalf("限速后登录状态码 = %d, want 429（body: %s）", res.StatusCode, body)
	}
}

// TestSecurityCrossOriginWriteRejected 写操作携带不同源 Origin 应 403。
func TestSecurityCrossOriginWriteRejected(t *testing.T) {
	ts := newTestHTTPServer(t, false, false)
	defer ts.Close()
	cookie := login(t, ts)

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/batch-actions", strings.NewReader(`{"action":"checkin"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://evil.example.com")
	req.AddCookie(cookie)
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("跨源写操作状态码 = %d, want 403", res.StatusCode)
	}

	// 对照组：同源 Origin 放行（batch-actions 无账号也返回 200 空结果）。
	req2, err := http.NewRequest(http.MethodPost, ts.URL+"/api/batch-actions", strings.NewReader(`{"action":"checkin"}`))
	if err != nil {
		t.Fatal(err)
	}
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Origin", ts.URL)
	req2.AddCookie(cookie)
	res2, err := ts.Client().Do(req2)
	if err != nil {
		t.Fatal(err)
	}
	defer res2.Body.Close()
	if res2.StatusCode != http.StatusOK {
		t.Fatalf("同源写操作状态码 = %d, want 200", res2.StatusCode)
	}
}

// TestSecurityUnauthenticatedWriteRejected 未登录访问写操作应 401（先于同源校验）。
func TestSecurityUnauthenticatedWriteRejected(t *testing.T) {
	ts := newTestHTTPServer(t, false, false)
	defer ts.Close()
	res := request(t, ts.Client(), http.MethodPost, ts.URL+"/api/batch-actions", map[string]any{"action": "checkin"}, nil, true)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("未登录写操作状态码 = %d, want 401", res.StatusCode)
	}
}

// TestSecurityLogoutInvalidatesSession logout 后旧 session cookie 立即不可用。
func TestSecurityLogoutInvalidatesSession(t *testing.T) {
	ts := newTestHTTPServer(t, false, false)
	defer ts.Close()
	cookie := login(t, ts)

	// 登录后可读。
	res := request(t, ts.Client(), http.MethodGet, ts.URL+"/api/overview", nil, cookie, false)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("登录后读状态码 = %d", res.StatusCode)
	}

	// logout（写操作，需要同源 Origin）。
	res = request(t, ts.Client(), http.MethodPost, ts.URL+"/api/logout", map[string]any{}, cookie, true)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("logout 状态码 = %d", res.StatusCode)
	}

	// 旧 cookie 立即失效。
	res = request(t, ts.Client(), http.MethodGet, ts.URL+"/api/overview", nil, cookie, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("logout 后旧会话状态码 = %d, want 401", res.StatusCode)
	}
}

// TestSecuritySessionCookieAttributes Set-Cookie 头验证：HttpOnly + SameSite=Strict + Path=/。
func TestSecuritySessionCookieAttributes(t *testing.T) {
	ts := newTestHTTPServer(t, false, false)
	defer ts.Close()
	res := request(t, ts.Client(), http.MethodPost, ts.URL+"/api/login", map[string]string{"username": "admin", "password": "test-password"}, nil, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login 状态码 = %d", res.StatusCode)
	}
	raw := res.Header.Get("Set-Cookie")
	for _, attr := range []string{"HttpOnly", "SameSite=Strict", "Path=/"} {
		if !strings.Contains(raw, attr) {
			t.Fatalf("Set-Cookie 缺少 %q: %q", attr, raw)
		}
	}
	if strings.Contains(raw, "Secure") {
		t.Log("提示：本地 HTTP 部署不强制 Secure 属性（生产经反代 HTTPS 时由部署层保证）")
	}
}

// TestSecurityLoginRateLimitWindow 不同来源 IP 的限速计数互相隔离。
// 直接构造 Server 并伪造 RemoteAddr（HTTP 层无法覆盖连接级 RemoteAddr）。
func TestSecurityLoginRateLimitWindow(t *testing.T) {
	dir := t.TempDir()
	store, err := authstore.New(dir + "/auths")
	if err != nil {
		t.Fatal(err)
	}
	state, err := NewState(dir + "/data")
	if err != nil {
		t.Fatal(err)
	}
	cfg := Config{Username: "admin", Password: "test-password", AuthDir: dir + "/auths", DataDir: dir + "/data", TimeoutSeconds: 5, Timezone: "Asia/Shanghai"}
	handler := NewServer(NewService(cfg, store, upstream.New(time.Second), state), http.NotFoundHandler()).Handler()

	loginFrom := func(remoteAddr, password string) int {
		req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"username":"admin","password":"`+password+`"}`))
		req.RemoteAddr = remoteAddr
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec.Code
	}

	const attacker = "203.0.113.7:50001"
	for i := 0; i < loginAttemptLimit; i++ {
		if code := loginFrom(attacker, "wrong"); code != http.StatusUnauthorized {
			t.Fatalf("第 %d 次失败登录 = %d, want 401", i+1, code)
		}
	}
	if code := loginFrom(attacker, "test-password"); code != http.StatusTooManyRequests {
		t.Fatalf("攻击者来源限速后应 429（即使密码正确），实际 %d", code)
	}
	// 另一来源独立计数，不受攻击者影响。
	if code := loginFrom("192.0.2.1:9999", "test-password"); code != http.StatusOK {
		t.Fatalf("新来源应可正常登录（200），实际 %d", code)
	}
}

// TestSecurityFailedLoginDoesNotRevealUserExistence 用户名错与密码错返回相同文案与状态码。
func TestSecurityFailedLoginDoesNotRevealUserExistence(t *testing.T) {
	ts := newTestHTTPServer(t, false, false)
	defer ts.Close()
	read := func(payload map[string]string) (int, string) {
		res := request(t, ts.Client(), http.MethodPost, ts.URL+"/api/login", payload, nil, false)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(body)
	}
	s1, b1 := read(map[string]string{"username": "no-such-user", "password": "x"})
	s2, b2 := read(map[string]string{"username": "admin", "password": "x"})
	if s1 != s2 || b1 != b2 {
		t.Fatalf("登录失败响应可区分用户是否存在: (%d,%s) vs (%d,%s)", s1, b1, s2, b2)
	}
}

// TestSecurityHeadersPresent 安全响应头回归（CSP/nosniff/frame-deny）。
func TestSecurityHeadersPresent(t *testing.T) {
	ts := newTestHTTPServer(t, false, false)
	defer ts.Close()
	res := request(t, ts.Client(), http.MethodGet, ts.URL+"/api/session", nil, nil, false)
	defer res.Body.Close()
	checks := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "no-referrer",
	}
	for h, want := range checks {
		if got := res.Header.Get(h); got != want {
			t.Fatalf("安全头 %s = %q, want %q", h, got, want)
		}
	}
	csp := res.Header.Get("Content-Security-Policy")
	if !strings.Contains(csp, "default-src 'self'") || !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Fatalf("CSP 不达标: %q", csp)
	}
}
