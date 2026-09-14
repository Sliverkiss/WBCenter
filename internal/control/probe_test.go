package control

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// fakeUpstream 是 Probe 依赖的最小上游接口的测试桩。
// 通过记录调用次数验证并发度与聚合语义。
type fakeUpstream struct {
	userResource func(a *authstore.Account) (*upstream.Credits, error)
	travelStatus func(a *authstore.Account) (*upstream.TravelState, error)
	buddyInfo    func(a *authstore.Account) (*upstream.Buddy, error)

	userResourceCalls atomic.Int32
	travelCalls       atomic.Int32
	buddyCalls        atomic.Int32
	// 用于观测并发上限：记录同时在飞的 UserResource 峰值。
	inflight    atomic.Int32
	maxInflight atomic.Int32
}

func (f *fakeUpstream) UserResource(a *authstore.Account) (*upstream.Credits, error) {
	f.userResourceCalls.Add(1)
	cur := f.inflight.Add(1)
	for {
		old := f.maxInflight.Load()
		if cur <= old || f.maxInflight.CompareAndSwap(old, cur) {
			break
		}
	}
	defer f.inflight.Add(-1)
	if f.userResource == nil {
		return &upstream.Credits{}, nil
	}
	return f.userResource(a)
}
func (f *fakeUpstream) TravelStatus(a *authstore.Account) (*upstream.TravelState, error) {
	f.travelCalls.Add(1)
	if f.travelStatus == nil {
		return &upstream.TravelState{State: "idle"}, nil
	}
	return f.travelStatus(a)
}
func (f *fakeUpstream) BuddyInfo(a *authstore.Account) (*upstream.Buddy, error) {
	f.buddyCalls.Add(1)
	if f.buddyInfo == nil {
		return &upstream.Buddy{ID: 1, Name: "阿喵"}, nil
	}
	return f.buddyInfo(a)
}

// newProbeService 构建一个挂 fakeUpstream 的 Service。
// 通过包内字段替换，不改动生产构造路径。
func newProbeService(t *testing.T, accounts []*authstore.Account, up ProbeUpstream) *Service {
	t.Helper()
	dir := t.TempDir()
	store, err := authstore.New(dir + "/auths")
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range accounts {
		if err := store.Save(a); err != nil {
			t.Fatal(err)
		}
	}
	state, err := NewState(dir + "/data")
	if err != nil {
		t.Fatal(err)
	}
	cfg := Config{Username: "admin", Password: "x", AuthDir: dir + "/auths", DataDir: dir + "/data", TimeoutSeconds: 5, Timezone: "Asia/Shanghai"}
	svc := NewService(cfg, store, upstream.New(time.Second), state)
	svc.prober = up
	return svc
}

func healthyAccount(uid string) *authstore.Account {
	return &authstore.Account{
		UID:          uid,
		Nickname:     "账号-" + uid,
		AccessToken:  "at",
		RefreshToken: "rt",
		ExpiresAt:    time.Now().Add(time.Hour).Unix(),
		Domain:       "copilot.tencent.com",
	}
}

// TestProbeAggregatesHealthyAccount 健康账号应同时返回积分、签到与旅行三个子对象。
func TestProbeAggregatesHealthyAccount(t *testing.T) {
	fake := &fakeUpstream{
		userResource: func(a *authstore.Account) (*upstream.Credits, error) {
			return &upstream.Credits{Remain: 1200, Used: 40, Size: 1240, Packages: 2}, nil
		},
		travelStatus: func(a *authstore.Account) (*upstream.TravelState, error) {
			return &upstream.TravelState{State: "traveling", RecordID: 9}, nil
		},
		buddyInfo: func(a *authstore.Account) (*upstream.Buddy, error) {
			return &upstream.Buddy{ID: 1, Name: "阿喵"}, nil
		},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-1001")}, fake)
	results, err := svc.Probe(context.Background())
	if err != nil {
		t.Fatalf("Probe: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	r := results[0]
	if r.UID != "u-1001" || r.Nickname != "账号-u-1001" {
		t.Errorf("结果缺 uid/nickname: %+v", r)
	}
	if !r.OK {
		t.Errorf("ok=false, error=%q", r.Error)
	}
	if r.Credits == nil || r.Credits.Current != 1200 {
		t.Errorf("credits 错误: %+v", r.Credits)
	}
	if r.Travel == nil || r.Travel.Status != "traveling" {
		t.Errorf("travel 错误: %+v", r.Travel)
	}
	if r.Checkin == nil {
		t.Errorf("checkin 子对象缺失（应始终返回，数据来源在契约中标注）: %+v", r)
	}
	if fake.userResourceCalls.Load() != 1 || fake.travelCalls.Load() != 1 || fake.buddyCalls.Load() != 1 {
		t.Errorf("调用次数 = userResource:%d travel:%d buddy:%d, want 全 1",
			fake.userResourceCalls.Load(), fake.travelCalls.Load(), fake.buddyCalls.Load())
	}
}

// TestProbeSingleAccountFailureDoesNotStopOthers 单账号失败不拖垮整批。
func TestProbeSingleAccountFailureDoesNotStopOthers(t *testing.T) {
	fake := &fakeUpstream{
		userResource: func(a *authstore.Account) (*upstream.Credits, error) {
			if a.UID == "u-bad" {
				return nil, &upstream.Error{Kind: upstream.ErrServer, Status: 503, Msg: "upstream down"}
			}
			return &upstream.Credits{Remain: 100}, nil
		},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-good"), healthyAccount("u-bad")}, fake)
	results, err := svc.Probe(context.Background())
	if err != nil {
		t.Fatalf("Probe: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	var good, bad *ProbeResult
	for i := range results {
		switch results[i].UID {
		case "u-good":
			good = &results[i]
		case "u-bad":
			bad = &results[i]
		}
	}
	if good == nil || !good.OK {
		t.Errorf("健康账号应成功: %+v", good)
	}
	if bad == nil || bad.OK || bad.Error == "" {
		t.Errorf("故障账号应失败且带 error: %+v", bad)
	}
}

// TestProbeSessionDeadMarked 上游 12153 应生成 session_dead=true 徽标。
func TestProbeSessionDeadMarked(t *testing.T) {
	fake := &fakeUpstream{
		userResource: func(a *authstore.Account) (*upstream.Credits, error) {
			return nil, &upstream.Error{Kind: upstream.ErrSessionDead, Status: 401, Msg: "code=12153"}
		},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-dead")}, fake)
	results, err := svc.Probe(context.Background())
	if err != nil {
		t.Fatalf("Probe: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d", len(results))
	}
	r := results[0]
	if r.OK {
		t.Errorf("ok 应为 false: %+v", r)
	}
	if !r.SessionDead {
		t.Errorf("session_dead 应为 true: %+v", r)
	}
	if !strings.Contains(r.Error, "12153") && !strings.Contains(r.Error, "会话失效") && !strings.Contains(r.Error, "session") {
		t.Errorf("error 文案应体现 12153/会话失效: %q", r.Error)
	}
}

// TestProbeConcurrencyIsBounded 受控并发：8 个账号同时在飞的 UserResource 调用数不应超过上限。
func TestProbeConcurrencyIsBounded(t *testing.T) {
	const n = 8
	fake := &fakeUpstream{
		userResource: func(a *authstore.Account) (*upstream.Credits, error) {
			time.Sleep(20 * time.Millisecond) // 撑开窗口让并发可见
			return &upstream.Credits{Remain: 1}, nil
		},
	}
	accs := make([]*authstore.Account, 0, n)
	for i := 0; i < n; i++ {
		accs = append(accs, healthyAccount("u-"+string(rune('a'+i))))
	}
	svc := newProbeService(t, accs, fake)
	if _, err := svc.Probe(context.Background()); err != nil {
		t.Fatalf("Probe: %v", err)
	}
	peak := fake.maxInflight.Load()
	if peak > ProbeConcurrency {
		t.Errorf("并发峰值 %d 超过上限 %d", peak, ProbeConcurrency)
	}
	if fake.userResourceCalls.Load() != n {
		t.Errorf("userResource 调用数 = %d, want %d", fake.userResourceCalls.Load(), n)
	}
}

// TestProbeContextCancelled 整体 ctx 取消时探测应尽早返回，不挂死。
func TestProbeContextCancelled(t *testing.T) {
	fake := &fakeUpstream{
		userResource: func(a *authstore.Account) (*upstream.Credits, error) {
			time.Sleep(200 * time.Millisecond)
			return &upstream.Credits{}, nil
		},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-1"), healthyAccount("u-2")}, fake)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := svc.Probe(ctx)
	// ctx 取消可视为整体错误，或以降级结果返回；两种都接受，但要在 1s 内回来。
	if time.Since(start) > time.Second {
		t.Errorf("Probe 在 ctx 取消后未尽早返回（耗时 %v）", time.Since(start))
	}
	_ = err
}

// TestProbeHandlerShape 端到端：POST /api/probe 走完整路由链，验证信封与 JSON 形状。
func TestProbeHandlerShape(t *testing.T) {
	fake := &fakeUpstream{}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-1001")}, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()

	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/probe", map[string]any{}, cookie, true)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body struct {
		Results []ProbeResult `json:"results"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Results) != 1 || body.Results[0].UID != "u-1001" {
		t.Fatalf("结果形状错误: %+v", body)
	}
}

// TestProbeHandlerUnauthorized 未登录访问应 401。
func TestProbeHandlerUnauthorized(t *testing.T) {
	fake := &fakeUpstream{}
	svc := newProbeService(t, nil, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/probe", map[string]any{}, nil, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}

// TestProbeHandlerReadOnly 只读模式下应拒绝（探测是触发上游动作的写语义）。
func TestProbeHandlerReadOnly(t *testing.T) {
	fake := &fakeUpstream{}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-1001")}, fake)
	svc.cfg.ReadOnly = true
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/probe", map[string]any{}, cookie, true)
	defer res.Body.Close()
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403", res.StatusCode)
	}
}
