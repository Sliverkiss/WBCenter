package control

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// newCookieTestServer 构造一个可直接驱动 handler 的 Server。
// 刻意不复用其他测试文件的辅助函数：本文件要能被单独合并。
func newCookieTestServer(t *testing.T) *Server {
	t.Helper()
	dir := t.TempDir()
	store, err := authstore.New(dir + "/auths")
	if err != nil {
		t.Fatal(err)
	}
	state, err := NewState(dir + "/data")
	if err != nil {
		t.Fatal(err)
	}
	cfg := Config{
		Username: "admin", Password: "test-password",
		AuthDir: dir + "/auths", DataDir: dir + "/data",
		TimeoutSeconds: 5, Timezone: "Asia/Shanghai",
	}
	return NewServer(NewService(cfg, store, upstream.New(cfg.Timeout()), state), http.NotFoundHandler())
}

// loginCookie 提交一次成功登录，返回会话 Cookie。
func loginCookie(t *testing.T, srv *Server, hdrs map[string]string, tlsConn bool) *http.Cookie {
	t.Helper()
	raw, err := json.Marshal(map[string]string{"username": "admin", "password": "test-password"})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(raw))
	for k, v := range hdrs {
		req.Header.Set(k, v)
	}
	if tlsConn {
		req.TLS = &tls.ConnectionState{}
	}
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("登录应成功，实际 %d：%s", rec.Code, rec.Body.String())
	}
	for _, c := range rec.Result().Cookies() {
		if c.Name == cookieName {
			return c
		}
	}
	t.Fatalf("登录响应没有下发 %s", cookieName)
	return nil
}

// TestSessionCookieIsSecureOnHTTPS HTTPS 链路上会话 Cookie 必须带 Secure。
//
// 面板自身只监听明文 HTTP（默认 127.0.0.1:8787），而 README 明确写着
// 「若经反向代理开放到局域网或公网，请使用 HTTPS…」，即生产部署是「反代终结 TLS」。
// 此时浏览器<->反代是 https，但 Set-Cookie 没有 Secure：只要用户（或一次 http 跳转、
// 一个 http 子域/同域链接）以明文 http:// 访问同一主机，浏览器就会照样回传这个会话
// Cookie，会话标识直接在明文链路里裸奔。HttpOnly / SameSite=Strict 都挡不住这一点。
func TestSessionCookieIsSecureOnHTTPS(t *testing.T) {
	srv := newCookieTestServer(t)
	cases := []struct {
		name    string
		hdrs    map[string]string
		tlsConn bool
	}{
		{"反向代理 X-Forwarded-Proto", map[string]string{"X-Forwarded-Proto": "https"}, false},
		{"多级代理 X-Forwarded-Proto 列表", map[string]string{"X-Forwarded-Proto": "https, http"}, false},
		{"直连 TLS", nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := loginCookie(t, srv, tc.hdrs, tc.tlsConn)
			if !c.Secure {
				t.Error("HTTPS 链路下会话 Cookie 缺少 Secure：明文 http:// 请求同样会带上它，会话标识可被中间人直接拿走")
			}
		})
	}
}

// TestSessionCookieNotSecureOnPlainHTTP 默认的明文直连部署不能被 Safe 化改坏：
// 加上 Secure 会让本地 http://127.0.0.1:8787 的登录直接失效。
func TestSessionCookieNotSecureOnPlainHTTP(t *testing.T) {
	srv := newCookieTestServer(t)
	c := loginCookie(t, srv, nil, false)
	if c.Secure {
		t.Error("明文直连（默认部署）下不能加 Secure，否则本地登录立刻失效")
	}
	if !c.HttpOnly {
		t.Error("HttpOnly 丢失：会话 Cookie 会对脚本可见")
	}
	if c.SameSite != http.SameSiteStrictMode {
		t.Errorf("SameSite=Strict 丢失，实际 %v", c.SameSite)
	}
	if c.Path != "/" {
		t.Errorf("Path 丢失，实际 %q", c.Path)
	}
}

// TestLogoutCookieKeepsSecureOnHTTPS 登出时清除 Cookie 的属性要与登录一致，
// 否则同一份会话 Cookie 会出现两套属性。
func TestLogoutCookieKeepsSecureOnHTTPS(t *testing.T) {
	srv := newCookieTestServer(t)
	c := loginCookie(t, srv, map[string]string{"X-Forwarded-Proto": "https"}, false)

	req := httptest.NewRequest(http.MethodPost, "/api/logout", bytes.NewReader(nil))
	req.Header.Set("X-Forwarded-Proto", "https")
	// /api/logout 不是 GET/HEAD，auth 中间件要求 Origin 与 Host 同源（CSRF 守卫），
	// 不带会被 403「跨站请求被拒绝」拦下。
	req.Header.Set("Origin", "http://"+req.Host)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: c.Value})
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("登出应成功，实际 %d：%s", rec.Code, rec.Body.String())
	}
	var cleared *http.Cookie
	for _, got := range rec.Result().Cookies() {
		if got.Name == cookieName {
			cleared = got
		}
	}
	if cleared == nil {
		t.Fatalf("登出响应没有回写 %s", cookieName)
	}
	if cleared.Value != "" || cleared.MaxAge >= 0 {
		t.Errorf("登出应清空会话 Cookie，实际 value=%q maxAge=%d", cleared.Value, cleared.MaxAge)
	}
	if !cleared.Secure {
		t.Error("HTTPS 链路下登出回写的 Cookie 也应带 Secure，否则属性不一致")
	}
}
