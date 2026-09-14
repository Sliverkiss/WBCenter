package control

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// fakeBatchUpstream 是 BatchActions 依赖的最小上游接口的测试桩。
// 记录每账号调用次数，验证禁用跳过与失败隔离语义。
type fakeBatchUpstream struct {
	checkin func(a *authstore.Account) error
	travel  func(a *authstore.Account) error
	refresh func(a *authstore.Account) error

	checkinCalls map[string]int
	travelCalls  map[string]int
	refreshCalls map[string]int
}

func newFakeBatchUpstream() *fakeBatchUpstream {
	return &fakeBatchUpstream{
		checkinCalls: map[string]int{},
		travelCalls:  map[string]int{},
		refreshCalls: map[string]int{},
	}
}
func (f *fakeBatchUpstream) DailyCheckin(a *authstore.Account) (*upstream.CheckinResult, error) {
	f.checkinCalls[a.UID]++
	if f.checkin != nil {
		return &upstream.CheckinResult{}, f.checkin(a)
	}
	return &upstream.CheckinResult{}, nil
}
func (f *fakeBatchUpstream) TravelOnce(a *authstore.Account) (*upstream.TravelResult, error) {
	f.travelCalls[a.UID]++
	if f.travel != nil {
		return &upstream.TravelResult{}, f.travel(a)
	}
	return &upstream.TravelResult{}, nil
}
func (f *fakeBatchUpstream) RefreshToken(a *authstore.Account) error {
	f.refreshCalls[a.UID]++
	if f.refresh != nil {
		return f.refresh(a)
	}
	return nil
}

// disabledAccount 构造凭据层标记禁用的账号。
func disabledAccount(uid string) *authstore.Account {
	a := healthyAccount(uid)
	a.Disabled = true
	return a
}

// newBatchService 与 newProbeService 同构，挂 batcher 桩。
func newBatchService(t *testing.T, accounts []*authstore.Account, up BatchUpstream) *Service {
	t.Helper()
	svc := newProbeService(t, accounts, &fakeUpstream{})
	svc.batcher = up
	return svc
}

// TestBatchActionsAllAccounts 全量批量签到：三个账号全部执行并聚合结果。
func TestBatchActionsAllAccounts(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, []*authstore.Account{
		healthyAccount("u-b1"), healthyAccount("u-b2"), healthyAccount("u-b3"),
	}, fake)
	res, err := svc.BatchActions(context.Background(), "checkin", nil)
	if err != nil {
		t.Fatalf("BatchActions: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("len(results) = %d, want 3", len(res))
	}
	for _, r := range res {
		if !r.OK {
			t.Errorf("账号 %s 应成功: %+v", r.UID, r)
		}
		if r.Message == "" {
			t.Errorf("账号 %s 缺 message: %+v", r.UID, r)
		}
	}
	if fake.checkinCalls["u-b1"] != 1 || fake.checkinCalls["u-b2"] != 1 || fake.checkinCalls["u-b3"] != 1 {
		t.Errorf("checkin 调用数 = %v, want 每账号 1 次", fake.checkinCalls)
	}
}

// TestBatchActionsDisabledSkipped 禁用账号跳过且不调上游；skipped=true、ok=true。
func TestBatchActionsDisabledSkipped(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, []*authstore.Account{
		healthyAccount("u-active"), disabledAccount("u-off"),
	}, fake)
	res, err := svc.BatchActions(context.Background(), "checkin", nil)
	if err != nil {
		t.Fatalf("BatchActions: %v", err)
	}
	var active, off *BatchActionResult
	for i := range res {
		switch res[i].UID {
		case "u-active":
			active = &res[i]
		case "u-off":
			off = &res[i]
		}
	}
	if active == nil || !active.OK || active.Skipped {
		t.Errorf("启用账号应正常执行: %+v", active)
	}
	if off == nil || !off.Skipped || !off.OK {
		t.Errorf("禁用账号应 skipped=true 且 ok=true: %+v", off)
	}
	if off != nil && off.Message == "" {
		t.Errorf("禁用账号应有跳过说明文案: %+v", off)
	}
	if fake.checkinCalls["u-off"] != 0 {
		t.Errorf("禁用账号不应调用上游, got %d 次", fake.checkinCalls["u-off"])
	}
	if fake.checkinCalls["u-active"] != 1 {
		t.Errorf("启用账号应调用 1 次, got %d", fake.checkinCalls["u-active"])
	}
}

// TestBatchActionsSingleFailureDoesNotStopOthers 单账号失败不中断整批。
func TestBatchActionsSingleFailureDoesNotStopOthers(t *testing.T) {
	fake := newFakeBatchUpstream()
	fake.checkin = func(a *authstore.Account) error {
		if a.UID == "u-bad" {
			return &upstream.Error{Kind: upstream.ErrServer, Status: 503, Msg: "upstream down"}
		}
		return nil
	}
	svc := newBatchService(t, []*authstore.Account{
		healthyAccount("u-good"), healthyAccount("u-bad"), healthyAccount("u-good2"),
	}, fake)
	res, err := svc.BatchActions(context.Background(), "checkin", nil)
	if err != nil {
		t.Fatalf("BatchActions: %v", err)
	}
	if len(res) != 3 {
		t.Fatalf("len(results) = %d, want 3（失败不中断整批）", len(res))
	}
	var bad *BatchActionResult
	okCount := 0
	for i := range res {
		if res[i].UID == "u-bad" {
			bad = &res[i]
		}
		if res[i].OK {
			okCount++
		}
	}
	if okCount != 2 {
		t.Errorf("成功数 = %d, want 2", okCount)
	}
	if bad == nil || bad.OK || bad.Error == "" {
		t.Errorf("失败账号应 ok=false 且带 error: %+v", bad)
	}
	// 失败账号之后的账号也必须被执行（不中断）。
	if fake.checkinCalls["u-good2"] != 1 {
		t.Errorf("失败后续账号未执行: %v", fake.checkinCalls)
	}
}

// TestBatchActionsSessionDeadMarked 12153 应打 session_dead 标记且不视为可重试失败。
func TestBatchActionsSessionDeadMarked(t *testing.T) {
	fake := newFakeBatchUpstream()
	fake.checkin = func(a *authstore.Account) error {
		return &upstream.Error{Kind: upstream.ErrSessionDead, Status: 401, Msg: "code=12153"}
	}
	svc := newBatchService(t, []*authstore.Account{healthyAccount("u-dead")}, fake)
	res, err := svc.BatchActions(context.Background(), "checkin", nil)
	if err != nil {
		t.Fatalf("BatchActions: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("len(results) = %d", len(res))
	}
	r := res[0]
	if r.OK {
		t.Errorf("ok 应为 false: %+v", r)
	}
	if !r.SessionDead {
		t.Errorf("session_dead 应为 true: %+v", r)
	}
	if !strings.Contains(r.Error, "12153") && !strings.Contains(r.Error, "会话失效") {
		t.Errorf("error 文案应体现 12153/会话失效: %q", r.Error)
	}
}

// TestBatchActionsUIDSubset 指定 uids 子集时只执行子集；不存在的 uid 静默忽略。
func TestBatchActionsUIDSubset(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, []*authstore.Account{
		healthyAccount("u-s1"), healthyAccount("u-s2"), healthyAccount("u-s3"),
	}, fake)
	res, err := svc.BatchActions(context.Background(), "travel", []string{"u-s2", "u-ghost"})
	if err != nil {
		t.Fatalf("BatchActions: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("len(results) = %d, want 1（只有 u-s2 在池内）", len(res))
	}
	if res[0].UID != "u-s2" {
		t.Errorf("结果账号 = %s, want u-s2", res[0].UID)
	}
	if fake.travelCalls["u-s1"] != 0 || fake.travelCalls["u-s3"] != 0 {
		t.Errorf("子集外账号不应执行: %v", fake.travelCalls)
	}
	if fake.travelCalls["u-s2"] != 1 {
		t.Errorf("u-s2 应执行 1 次, got %d", fake.travelCalls["u-s2"])
	}
}

// TestBatchActionsRefreshSavesAccount refresh 动作成功后应回写凭据文件。
func TestBatchActionsRefreshSavesAccount(t *testing.T) {
	fake := newFakeBatchUpstream()
	fake.refresh = func(a *authstore.Account) error {
		a.AccessToken = "rotated"
		return nil
	}
	svc := newBatchService(t, []*authstore.Account{healthyAccount("u-r1")}, fake)
	if _, err := svc.BatchActions(context.Background(), "refresh", nil); err != nil {
		t.Fatalf("BatchActions: %v", err)
	}
	got, err := svc.store.Get("u-r1")
	if err != nil {
		t.Fatalf("store.Get: %v", err)
	}
	if got.AccessToken != "rotated" {
		t.Errorf("refresh 后凭据未回写, accessToken = %q", got.AccessToken)
	}
}

// TestBatchActionsUnknownAction 未知动作返回错误。
func TestBatchActionsUnknownAction(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, []*authstore.Account{healthyAccount("u-x")}, fake)
	_, err := svc.BatchActions(context.Background(), "explode", nil)
	if err == nil {
		t.Fatal("未知动作应返回错误")
	}
}

// TestBatchActionsReadOnly 只读模式拒绝。
func TestBatchActionsReadOnly(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, []*authstore.Account{healthyAccount("u-x")}, fake)
	svc.cfg.ReadOnly = true
	_, err := svc.BatchActions(context.Background(), "checkin", nil)
	if err == nil || !strings.Contains(err.Error(), "只读") {
		t.Fatalf("只读模式应拒绝: %v", err)
	}
}

// ---- handler 层 ----

// TestBatchActionsHandlerShape 端到端：信封形状与契约一致。
func TestBatchActionsHandlerShape(t *testing.T) {
	fake := newFakeBatchUpstream()
	fake.checkin = func(a *authstore.Account) error {
		if a.UID == "u-h-dead" {
			return &upstream.Error{Kind: upstream.ErrSessionDead, Status: 401, Msg: "code=12153"}
		}
		return nil
	}
	svc := newBatchService(t, []*authstore.Account{
		healthyAccount("u-h-ok"), disabledAccount("u-h-off"), healthyAccount("u-h-dead"),
	}, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/batch-actions",
		map[string]any{"action": "checkin"}, cookie, true)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body struct {
		Results []BatchActionResult `json:"results"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Results) != 3 {
		t.Fatalf("len(results) = %d, want 3", len(body.Results))
	}
	byUID := map[string]BatchActionResult{}
	for _, r := range body.Results {
		byUID[r.UID] = r
		if r.Nickname == "" {
			t.Errorf("缺 nickname: %+v", r)
		}
	}
	if !byUID["u-h-ok"].OK || byUID["u-h-ok"].Skipped {
		t.Errorf("成功账号: %+v", byUID["u-h-ok"])
	}
	if !byUID["u-h-off"].Skipped {
		t.Errorf("禁用账号应 skipped: %+v", byUID["u-h-off"])
	}
	if !byUID["u-h-dead"].SessionDead {
		t.Errorf("12153 账号应 session_dead: %+v", byUID["u-h-dead"])
	}
}

// TestBatchActionsHandlerBadAction 非法 action 返回 400。
func TestBatchActionsHandlerBadAction(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, nil, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	for _, body := range []map[string]any{
		{"action": "explode"},
		{"uids": []string{"u-1"}}, // 缺 action
	} {
		res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/batch-actions", body, cookie, true)
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("body=%v status = %d, want 400", body, res.StatusCode)
		}
		res.Body.Close()
	}
}

// TestBatchActionsHandlerUnauthorized 未登录 401。
func TestBatchActionsHandlerUnauthorized(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, nil, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/batch-actions",
		map[string]any{"action": "checkin"}, nil, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}

// TestBatchActionsHandlerReadOnly 只读模式 403。
func TestBatchActionsHandlerReadOnly(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, []*authstore.Account{healthyAccount("u-ro")}, fake)
	svc.cfg.ReadOnly = true
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/batch-actions",
		map[string]any{"action": "checkin"}, cookie, true)
	defer res.Body.Close()
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403", res.StatusCode)
	}
}

// TestBatchActionsHandlerUIDsSubsetViaHTTP handler 层 uids 子集过滤。
func TestBatchActionsHandlerUIDsSubsetViaHTTP(t *testing.T) {
	fake := newFakeBatchUpstream()
	svc := newBatchService(t, []*authstore.Account{
		healthyAccount("u-p1"), healthyAccount("u-p2"),
	}, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodPost, srv.URL+"/api/batch-actions",
		map[string]any{"action": "travel", "uids": []string{"u-p2"}}, cookie, true)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body struct {
		Results []BatchActionResult `json:"results"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Results) != 1 || body.Results[0].UID != "u-p2" {
		t.Fatalf("子集结果错误: %+v", body.Results)
	}
}

// 防止未使用告警（errors 留给后续断言扩展）。
var _ = errors.Is
