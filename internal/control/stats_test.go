package control

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// fakeStatsUpstream 组合 probe/credits 两类上游调用，供 StatsSummary 测试。
// 积分四项走 UserResource + DailyFreePackages（与 /api/credits 同口径）。
type fakeStatsUpstream struct {
	fakeUpstream
	credits         map[string]*upstream.Credits
	packages        map[string][]upstream.DailyPackage
	userResourceErr map[string]error
	packagesErr     map[string]error
}

func (f *fakeStatsUpstream) UserResource(a *authstore.Account) (*upstream.Credits, error) {
	if err := f.userResourceErr[a.UID]; err != nil {
		return nil, err
	}
	if c, ok := f.credits[a.UID]; ok {
		return c, nil
	}
	return &upstream.Credits{}, nil
}

// newStatsService 构建挂 fakeStatsUpstream 的 Service。
// 通过包内字段替换 creditsFetcher，不动生产构造路径。
func newStatsService(t *testing.T, accounts []*authstore.Account, up *fakeStatsUpstream) *Service {
	t.Helper()
	svc := newProbeService(t, accounts, &up.fakeUpstream)
	svc.creditsFetcher = up
	return svc
}

func (f *fakeStatsUpstream) creditsFor(ctx context.Context, a *authstore.Account) CreditSummary {
	r := CreditSummary{UID: a.UID, Nickname: a.Nickname, FetchedAt: time.Now()}
	c, err := f.UserResource(a)
	if err != nil {
		r.Error = "上游积分概要查询失败"
		return r
	}
	r.Current = c.Remain
	r.Packages = c.Packages
	if err := f.packagesErr[a.UID]; err != nil {
		r.Error = "已取得当前积分；今日套餐明细暂不可用"
		return r
	}
	for _, p := range f.packages[a.UID] {
		r.TodayAllocated += p.Total
		r.TodayConsumed += p.Used
		r.TodayRemaining += p.Remaining
	}
	return r
}

// TestStatsSummaryAggregatesCredits 正常聚合：积分四项 = 各账号求和。
func TestStatsSummaryAggregatesCredits(t *testing.T) {
	fake := &fakeStatsUpstream{
		credits: map[string]*upstream.Credits{
			"u-stat-1": {Remain: 1200, Packages: 2},
			"u-stat-2": {Remain: 350, Packages: 1},
		},
		packages: map[string][]upstream.DailyPackage{
			"u-stat-1": {{Total: 100, Used: 40, Remaining: 60}},
			"u-stat-2": {{Total: 50, Used: 10, Remaining: 40}},
		},
	}
	svc := newStatsService(t, []*authstore.Account{healthyAccount("u-stat-1"), healthyAccount("u-stat-2")}, fake)
	sum := svc.StatsSummary(context.Background())
	if sum.AccountsTotal != 2 {
		t.Errorf("accounts_total = %d, want 2", sum.AccountsTotal)
	}
	if sum.CreditsCurrentTotal != 1550 {
		t.Errorf("credits_current_total = %d, want 1550", sum.CreditsCurrentTotal)
	}
	if sum.CreditsTodayAllocatedTotal != 150 || sum.CreditsTodayConsumedTotal != 50 || sum.CreditsTodayRemainingTotal != 100 {
		t.Errorf("今日额度合计错误: alloc=%v consumed=%v remain=%v, want 150/50/100",
			sum.CreditsTodayAllocatedTotal, sum.CreditsTodayConsumedTotal, sum.CreditsTodayRemainingTotal)
	}
	if sum.GeneratedAt.IsZero() {
		t.Error("generated_at 不应为零值")
	}
}

// TestStatsSummaryEmptyPool 空账号池：全零值且不报错。
func TestStatsSummaryEmptyPool(t *testing.T) {
	fake := &fakeStatsUpstream{}
	svc := newStatsService(t, nil, fake)
	sum := svc.StatsSummary(context.Background())
	if sum.AccountsTotal != 0 || sum.CreditsCurrentTotal != 0 {
		t.Errorf("空池应全零: %+v", sum)
	}
}

// TestStatsSummaryUpstreamFailureDegrades 单账号上游失败：该账号积分按 0 计入，不拖垮整体。
func TestStatsSummaryUpstreamFailureDegrades(t *testing.T) {
	fake := &fakeStatsUpstream{
		credits: map[string]*upstream.Credits{
			"u-stat-ok": {Remain: 500, Packages: 1},
		},
		packages: map[string][]upstream.DailyPackage{
			"u-stat-ok": {{Total: 100, Used: 20, Remaining: 80}},
		},
		userResourceErr: map[string]error{
			"u-stat-bad": &upstream.Error{Kind: upstream.ErrServer, Status: 503, Msg: "down"},
		},
	}
	svc := newStatsService(t, []*authstore.Account{healthyAccount("u-stat-ok"), healthyAccount("u-stat-bad")}, fake)
	sum := svc.StatsSummary(context.Background())
	if sum.AccountsTotal != 2 {
		t.Errorf("accounts_total = %d, want 2（失败账号仍计入总数）", sum.AccountsTotal)
	}
	if sum.CreditsCurrentTotal != 500 {
		t.Errorf("credits_current_total = %d, want 500（失败账号按 0 计入）", sum.CreditsCurrentTotal)
	}
	if sum.CreditsTodayAllocatedTotal != 100 {
		t.Errorf("allocated = %v, want 100", sum.CreditsTodayAllocatedTotal)
	}
}

// TestStatsSummaryCheckinFromLastProbe 签到计数来自 lastProbe 缓存；未探测时 pending=total。
func TestStatsSummaryCheckinFromLastProbe(t *testing.T) {
	fake := &fakeStatsUpstream{}
	svc := newStatsService(t, []*authstore.Account{healthyAccount("u-ck-1"), healthyAccount("u-ck-2")}, fake)

	// 未运行过探测：done=0、pending=total（保守零值，不主动触发上游）。
	sum := svc.StatsSummary(context.Background())
	if sum.CheckinDoneToday != 0 || sum.CheckinPendingToday != 2 {
		t.Errorf("未探测时签到计数应保守: done=%d pending=%d, want 0/2", sum.CheckinDoneToday, sum.CheckinPendingToday)
	}

	// 注入探测缓存：一个已签、一个未签。
	svc.mu.Lock()
	svc.lastProbe = []ProbeResult{
		{UID: "u-ck-1", Nickname: "n1", OK: true, Checkin: &ProbeCheckin{CheckedIn: true}},
		{UID: "u-ck-2", Nickname: "n2", OK: true, Checkin: &ProbeCheckin{CheckedIn: false}},
	}
	svc.mu.Unlock()
	sum = svc.StatsSummary(context.Background())
	if sum.CheckinDoneToday != 1 || sum.CheckinPendingToday != 1 {
		t.Errorf("探测后签到计数: done=%d pending=%d, want 1/1", sum.CheckinDoneToday, sum.CheckinPendingToday)
	}
}

// TestStatsSummaryStatusDistribution 状态分布：健康/过期/禁用按凭据判定，12153 来自探测缓存。
func TestStatsSummaryStatusDistribution(t *testing.T) {
	expired := healthyAccount("u-exp-01")
	expired.ExpiresAt = time.Now().Add(-time.Hour).Unix()
	fake := &fakeStatsUpstream{}
	svc := newStatsService(t, []*authstore.Account{
		healthyAccount("u-ok-01"), expired, disabledAccount("u-dis-01"),
	}, fake)
	// 注入 12153 探测缓存（一个会话死）。
	svc.mu.Lock()
	svc.lastProbe = []ProbeResult{
		{UID: "u-ok-01", Nickname: "n", OK: true},
		{UID: "u-exp-01", Nickname: "n", OK: false, SessionDead: true, Error: "上游会话失效（12153）"},
	}
	svc.mu.Unlock()
	sum := svc.StatsSummary(context.Background())
	if sum.StatusHealthy != 2 { // u-ok-01 + u-dis-01（禁用与过期是两个维度）
		t.Errorf("status_healthy = %d, want 2", sum.StatusHealthy)
	}
	if sum.StatusExpired != 1 {
		t.Errorf("status_expired = %d, want 1", sum.StatusExpired)
	}
	if sum.StatusSessionDead != 1 {
		t.Errorf("status_session_dead = %d, want 1", sum.StatusSessionDead)
	}
	if sum.StatusDisabled != 1 {
		t.Errorf("status_disabled = %d, want 1", sum.StatusDisabled)
	}
}

// TestStatsSummaryAutomationRunsToday 自动化今日统计来自 RunRecord 历史（服务端时区）。
func TestStatsSummaryAutomationRunsToday(t *testing.T) {
	fake := &fakeStatsUpstream{}
	svc := newStatsService(t, nil, fake)
	st := svc.State()
	now := time.Now()
	yesterday := now.Add(-26 * time.Hour)
	st.InjectRuns([]RunRecord{
		{At: now.Add(-time.Hour), Action: "checkin", OK: true, Message: "成功 3，失败 0"},
		{At: now.Add(-2 * time.Hour), Action: "travel", OK: true, Message: "成功 2，失败 0"},
		{At: now.Add(-3 * time.Hour), Action: "checkin", OK: false, Message: "上游失败"},
		{At: yesterday, Action: "checkin", OK: true, Message: "昨天的不计入"},
	})
	sum := svc.StatsSummary(context.Background())
	if sum.AutomationRunsTodayOK != 2 || sum.AutomationRunsTodayFailed != 1 {
		t.Errorf("今日自动化统计 = ok:%d failed:%d, want 2/1",
			sum.AutomationRunsTodayOK, sum.AutomationRunsTodayFailed)
	}
}

// TestStatsSummaryHandlerShape 端到端：GET /api/stats/summary 契约形状。
func TestStatsSummaryHandlerShape(t *testing.T) {
	fake := &fakeStatsUpstream{
		credits:  map[string]*upstream.Credits{"u-hs-01": {Remain: 100}},
		packages: map[string][]upstream.DailyPackage{"u-hs-01": {{Total: 10, Used: 3, Remaining: 7}}},
	}
	svc := newStatsService(t, []*authstore.Account{healthyAccount("u-hs-01")}, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodGet, srv.URL+"/api/stats/summary", nil, cookie, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body StatsSummary
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.AccountsTotal != 1 || body.CreditsCurrentTotal != 100 || body.CreditsTodayRemainingTotal != 7 {
		t.Errorf("契约字段错误: %+v", body)
	}
	if body.GeneratedAt == "" {
		t.Error("generated_at 缺省")
	}
}

// TestStatsSummaryHandlerUnauthorized 未登录 401。
func TestStatsSummaryHandlerUnauthorized(t *testing.T) {
	fake := &fakeStatsUpstream{}
	svc := newStatsService(t, nil, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	res := request(t, srv.Client(), http.MethodGet, srv.URL+"/api/stats/summary", nil, nil, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}
