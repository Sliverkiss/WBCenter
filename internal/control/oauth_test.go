package control

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// fakeOAuthUpstream 是 OAuth 状态机依赖的最小上游接口的测试桩。
type fakeOAuthUpstream struct {
	start func(region upstream.Region) (string, string, error)
	poll  func(region upstream.Region, state string) (*authstore.Account, error)

	startRegions []upstream.Region
	pollRegions  []upstream.Region
	pollCalls    atomic.Int32
}

func (f *fakeOAuthUpstream) StartLogin(region upstream.Region) (string, string, error) {
	f.startRegions = append(f.startRegions, region)
	if f.start != nil {
		return f.start(region)
	}
	return "state-123", "https://example.invalid/auth?state=state-123", nil
}

func (f *fakeOAuthUpstream) PollLogin(region upstream.Region, state string) (*authstore.Account, error) {
	f.pollRegions = append(f.pollRegions, region)
	f.pollCalls.Add(1)
	if f.poll != nil {
		return f.poll(region, state)
	}
	return nil, nil
}

// newOAuthService 构建挂 OAuth 桩的 Service（auths 落盘到临时目录）。
func newOAuthService(t *testing.T, up OAuthUpstream) (*Service, *authstore.Store) {
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
	cfg := Config{Username: "admin", Password: "test-password", AuthDir: dir + "/auths", DataDir: dir + "/data", TimeoutSeconds: 5, Timezone: "Asia/Shanghai"}
	svc := NewService(cfg, store, upstream.New(time.Second), state)
	svc.oauth = up
	return svc, store
}

// TestOAuthWaitingThenSuccess 轮询序列 waiting→waiting→success：
// 前两次返回 waiting，第三次写入凭据并返回 success+uid+nickname。
func TestOAuthWaitingThenSuccess(t *testing.T) {
	fake := &fakeOAuthUpstream{}
	acct := &authstore.Account{UID: "u-oauth-1", Nickname: "新账号", AccessToken: "at", RefreshToken: "rt", Domain: "copilot.tencent.com"}
	calls := 0
	fake.poll = func(region upstream.Region, state string) (*authstore.Account, error) {
		calls++
		if calls < 3 {
			return nil, nil
		}
		return acct, nil
	}
	svc, store := newOAuthService(t, fake)

	id, url, err := svc.StartOAuth("cn")
	if err != nil {
		t.Fatalf("StartOAuth: %v", err)
	}
	if id == "" || !strings.Contains(url, "state-123") {
		t.Fatalf("StartOAuth 返回异常: id=%q url=%q", id, url)
	}

	for i := 0; i < 2; i++ {
		res, err := svc.PollOAuth(context.Background(), id)
		if err != nil {
			t.Fatalf("第 %d 次 PollOAuth 报错: %v", i+1, err)
		}
		if res["status"] != "waiting" {
			t.Fatalf("第 %d 次轮询 status = %v, want waiting", i+1, res["status"])
		}
	}
	res, err := svc.PollOAuth(context.Background(), id)
	if err != nil {
		t.Fatalf("成功轮询报错: %v", err)
	}
	if res["status"] != "success" || res["uid"] != "u-oauth-1" || res["nickname"] != "新账号" {
		t.Fatalf("成功响应字段错误: %+v", res)
	}
	// 凭据已写入 auths。
	saved, err := store.Get("u-oauth-1")
	if err != nil {
		t.Fatalf("凭据未写入 auths: %v", err)
	}
	if saved.AccessToken != "at" || saved.Nickname != "新账号" {
		t.Fatalf("凭据内容错误: %+v", saved)
	}
	// 终态幂等：再轮询仍返回 success，不再次调上游。
	before := fake.pollCalls.Load()
	res, err = svc.PollOAuth(context.Background(), id)
	if err != nil || res["status"] != "success" {
		t.Fatalf("终态复读应仍 success: res=%+v err=%v", res, err)
	}
	if fake.pollCalls.Load() != before {
		t.Fatal("终态后不应再次调用上游")
	}
}

// TestOAuthWaitingThenError 上游 5xx 错误应归一为 error 状态并带文案，且会话进入终态。
func TestOAuthWaitingThenError(t *testing.T) {
	fake := &fakeOAuthUpstream{}
	calls := 0
	fake.poll = func(region upstream.Region, state string) (*authstore.Account, error) {
		calls++
		if calls == 1 {
			return nil, nil
		}
		return nil, &upstream.Error{Kind: upstream.ErrServer, Status: 503, Msg: "上游不可用"}
	}
	svc, _ := newOAuthService(t, fake)
	id, _, err := svc.StartOAuth("cn")
	if err != nil {
		t.Fatalf("StartOAuth: %v", err)
	}
	res, err := svc.PollOAuth(context.Background(), id)
	if err != nil || res["status"] != "waiting" {
		t.Fatalf("首次轮询应 waiting: res=%+v err=%v", res, err)
	}
	res, err = svc.PollOAuth(context.Background(), id)
	if err != nil {
		t.Fatalf("error 终态不应以 Go error 返回: %v", err)
	}
	if res["status"] != "error" {
		t.Fatalf("status = %v, want error", res["status"])
	}
	if msg, _ := res["message"].(string); msg == "" {
		t.Fatalf("error 状态缺 message: %+v", res)
	}
	// 终态后再轮询同样返回 error（幂等），不再次调上游。
	before := fake.pollCalls.Load()
	res, err = svc.PollOAuth(context.Background(), id)
	if err != nil || res["status"] != "error" {
		t.Fatalf("终态复读应仍 error: res=%+v err=%v", res, err)
	}
	if fake.pollCalls.Load() != before {
		t.Fatal("终态后不应再次调用上游")
	}
}

// TestOAuthTimeout 超过 5 分钟轮询上限应归一为 timeout 状态（终态）。
func TestOAuthTimeout(t *testing.T) {
	fake := &fakeOAuthUpstream{}
	svc, _ := newOAuthService(t, fake)
	id, _, err := svc.StartOAuth("cn")
	if err != nil {
		t.Fatalf("StartOAuth: %v", err)
	}
	// 直接把会话创建时间改到 6 分钟前，模拟超时。
	svc.mu.Lock()
	flow := svc.logins[id]
	flow.Created = time.Now().Add(-6 * time.Minute)
	svc.logins[id] = flow
	svc.mu.Unlock()

	res, err := svc.PollOAuth(context.Background(), id)
	if err != nil {
		t.Fatalf("timeout 不应以 Go error 返回: %v", err)
	}
	if res["status"] != "timeout" {
		t.Fatalf("status = %v, want timeout", res["status"])
	}
	if fake.pollCalls.Load() != 0 {
		t.Fatal("超时后不应再调上游")
	}
}

// TestOAuthRegionRouting cn/global 应路由到对应上游 region。
func TestOAuthRegionRouting(t *testing.T) {
	fake := &fakeOAuthUpstream{}
	svc, _ := newOAuthService(t, fake)

	idCN, _, err := svc.StartOAuth("cn")
	if err != nil {
		t.Fatalf("StartOAuth cn: %v", err)
	}
	idGlobal, _, err := svc.StartOAuth("global")
	if err != nil {
		t.Fatalf("StartOAuth global: %v", err)
	}
	if _, err := svc.PollOAuth(context.Background(), idCN); err != nil {
		t.Fatalf("PollOAuth cn: %v", err)
	}
	if _, err := svc.PollOAuth(context.Background(), idGlobal); err != nil {
		t.Fatalf("PollOAuth global: %v", err)
	}
	if len(fake.startRegions) != 2 || fake.startRegions[0] != upstream.RegionCN || fake.startRegions[1] != upstream.RegionGlobal {
		t.Fatalf("StartLogin 区域路由错误: %v", fake.startRegions)
	}
	if len(fake.pollRegions) != 2 || fake.pollRegions[0] != upstream.RegionCN || fake.pollRegions[1] != upstream.RegionGlobal {
		t.Fatalf("PollLogin 区域路由错误: %v", fake.pollRegions)
	}
}

// TestOAuthInvalidRegion 非法区域应报错且不创建会话。
func TestOAuthInvalidRegion(t *testing.T) {
	svc, _ := newOAuthService(t, &fakeOAuthUpstream{})
	if _, _, err := svc.StartOAuth("moon"); err == nil {
		t.Fatal("非法区域应报错")
	}
}

// TestOAuthPollUnknownID 未知会话 id 应报「不存在或已过期」。
func TestOAuthPollUnknownID(t *testing.T) {
	svc, _ := newOAuthService(t, &fakeOAuthUpstream{})
	if _, err := svc.PollOAuth(context.Background(), "no-such-id"); err == nil {
		t.Fatal("未知会话应报错")
	}
}

// TestOAuthPoll11217Waiting 上游 11217（token not ready）应归一为 waiting 继续轮询，而非 error。
func TestOAuthPoll11217Waiting(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasPrefix(r.URL.Path, "/v2/plugin/auth/state") {
			_, _ = w.Write([]byte(`{"code":0,"data":{"state":"st-11217","authUrl":"https://example.invalid/auth"}}`))
			return
		}
		// 轮询端点返回 11217。
		_, _ = w.Write([]byte(`{"code":11217,"msg":"token not ready"}`))
	}))
	defer srv.Close()

	client := upstream.New(time.Second)
	client.BaseOverride = srv.URL
	a, err := client.PollLogin(upstream.RegionCN, "st-11217")
	if err != nil {
		t.Fatalf("11217 应归一为等待（nil, nil），实际 err = %v", err)
	}
	if a != nil {
		t.Fatalf("11217 不应返回账号: %+v", a)
	}
}

// TestOAuthPollSuccessWritesAccount 真实 client 全链路：state→poll(11217)→poll(成功)→login/account。
func TestOAuthPollSuccessWritesAccount(t *testing.T) {
	var pollCount atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.Contains(r.URL.Path, "/v2/plugin/auth/state"):
			_, _ = w.Write([]byte(`{"code":0,"data":{"state":"st-ok","authUrl":"https://example.invalid/auth?state=st-ok"}}`))
		case strings.Contains(r.URL.Path, "/v2/plugin/auth/token"):
			if pollCount.Add(1) == 1 {
				_, _ = w.Write([]byte(`{"code":11217,"msg":"token not ready"}`))
				return
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"accessToken":"real-at","refreshToken":"real-rt","expiresIn":7200,"domain":"copilot.tencent.com"}}`))
		case strings.Contains(r.URL.Path, "/v2/plugin/login/account"):
			if r.Header.Get("Authorization") != "Bearer real-at" {
				t.Errorf("login/account 缺少 Bearer 头: %q", r.Header.Get("Authorization"))
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"uid":"u-real","nickname":"真实账号","enterpriseId":"e-1"}}`))
		default:
			t.Errorf("未知上游路径: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	client := upstream.New(time.Second)
	client.BaseOverride = srv.URL

	// Service 层需要走真实 client：构造一个把 StartLogin/PollLogin 委托给真实 client 的适配桩。
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
	svc := NewService(cfg, store, client, state)

	id, url, err := svc.StartOAuth("cn")
	if err != nil {
		t.Fatalf("StartOAuth: %v", err)
	}
	if !strings.Contains(url, "st-ok") {
		t.Fatalf("授权 URL 应带 state: %q", url)
	}
	res, err := svc.PollOAuth(context.Background(), id)
	if err != nil || res["status"] != "waiting" {
		t.Fatalf("11217 应归一为 waiting: res=%+v err=%v", res, err)
	}
	res, err = svc.PollOAuth(context.Background(), id)
	if err != nil {
		t.Fatalf("成功轮询: %v", err)
	}
	if res["status"] != "success" || res["uid"] != "u-real" || res["nickname"] != "真实账号" {
		t.Fatalf("成功响应错误: %+v", res)
	}
	saved, err := store.Get("u-real")
	if err != nil {
		t.Fatalf("凭据未写入: %v", err)
	}
	if saved.AccessToken != "real-at" || saved.EnterpriseID != "e-1" {
		t.Fatalf("凭据字段错误: %+v", saved)
	}
}

// TestOAuthReadOnlyGuard 只读模式下成功轮询应报 502 语义错误（不写凭据）。
func TestOAuthReadOnlyGuard(t *testing.T) {
	fake := &fakeOAuthUpstream{}
	acct := &authstore.Account{UID: "u-ro", Nickname: "只读", AccessToken: "at", Domain: "copilot.tencent.com"}
	fake.poll = func(region upstream.Region, state string) (*authstore.Account, error) { return acct, nil }
	svc, store := newOAuthService(t, fake)
	svc.cfg.ReadOnly = true

	id, _, err := svc.StartOAuth("cn")
	if err != nil {
		t.Fatalf("StartOAuth: %v", err)
	}
	if _, err := svc.PollOAuth(context.Background(), id); err == nil {
		t.Fatal("只读模式成功轮询应报错")
	}
	if _, err := store.Get("u-ro"); err == nil {
		t.Fatal("只读模式不应写入凭据")
	}
}
