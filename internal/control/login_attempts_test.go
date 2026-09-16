package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// newTestControlServer 构造一个可直接驱动 handler 的 Server，
// 便于检查内部记账状态（attempts）。与 newTestHTTPServer 的区别是它不做
// 真实监听，返回 *Server 本身。
func newTestControlServer(t *testing.T) *Server {
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

// attemptLogin 以给定来源地址提交一次登录（走完整 Handler，含中间件）。
func attemptLogin(t *testing.T, srv *Server, remoteAddr, user, pass string) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(map[string]string{"username": user, "password": pass})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(raw))
	req.RemoteAddr = remoteAddr + ":54321"
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	return rec
}

// TestStaleLoginAttemptsAreCollected 早已过期窗口的失败记录必须被回收。
//
// attempts 的条目只在「同一来源在窗口结束后再次尝试」或「该来源登录成功」时才
// 被删除，于是一个**只失败一次、再也不回来**的来源会留下一条永久记录。来源数量
// 不受控（IPv6 下一个 /64 就能造出天文数字的不同地址，每个地址发一个请求即可），
// 这条路径等于让未认证的远端可以无界增长服务端内存，而那些过期记录再也不会有任何
// 用处——只有该来源自己回来时才会被读到，届时同样按过期处理。
func TestStaleLoginAttemptsAreCollected(t *testing.T) {
	srv := newTestControlServer(t)
	stale := time.Now().Add(-2 * loginAttemptWindow)
	for i := 0; i < 500; i++ {
		ip := fmt.Sprintf("10.%d.%d.%d", i>>16&0xff, i>>8&0xff, i&0xff)
		srv.attempts[ip] = loginAttempt{count: 1, first: stale}
	}
	if got := len(srv.attempts); got != 500 {
		t.Fatalf("前置条件不成立：attempts=%d", got)
	}

	// 一次来自新来源的失败登录：应当顺带回收早已过期的记录。
	attemptLogin(t, srv, "203.0.113.7", "admin", "wrong-password")

	if got := len(srv.attempts); got > 1 {
		t.Errorf("过期失败记录未被回收：attempts 仍有 %d 条（只失败一次、再也不回来的来源会永久留档 → 无界增长）", got)
	}
	if _, ok := srv.attempts["203.0.113.7"]; !ok {
		t.Error("当前来源的失败记录不应被回收，否则 per-IP 限流会失效")
	}
}

// TestLiveLoginAttemptsAreKept 窗口内的记录必须保留，且限流语义不变。
func TestLiveLoginAttemptsAreKept(t *testing.T) {
	srv := newTestControlServer(t)
	const ip = "203.0.113.9"
	for i := 0; i < loginAttemptLimit; i++ {
		if rec := attemptLogin(t, srv, ip, "admin", "wrong-password"); rec.Code != http.StatusUnauthorized {
			t.Fatalf("第 %d 次失败登录 → %d，期望 401", i+1, rec.Code)
		}
	}
	rec, ok := srv.attempts[ip]
	if !ok {
		t.Fatal("窗口内的失败记录被误删，per-IP 限流会失效")
	}
	if rec.count != loginAttemptLimit {
		t.Errorf("失败计数=%d，期望 %d", rec.count, loginAttemptLimit)
	}
	if got := attemptLogin(t, srv, ip, "admin", "wrong-password"); got.Code != http.StatusTooManyRequests {
		t.Errorf("超过上限后 → %d，期望 429", got.Code)
	}
}
