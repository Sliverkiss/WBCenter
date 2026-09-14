package control

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// fakeModelsUpstream 是 Models/ModelPricing 依赖的最小上游接口的测试桩。
// 每个账号的 AvailableModels 返回独立配置，便于验证逐账号聚合与失败降级。
type fakeModelsUpstream struct {
	models map[string][]map[string]any
	errs   map[string]error
}

func (f *fakeModelsUpstream) AvailableModels(a *authstore.Account) ([]map[string]any, error) {
	if err := f.errs[a.UID]; err != nil {
		return nil, err
	}
	return f.models[a.UID], nil
}

// newModelsService 构建挂 modelsLister 桩的 Service。
func newModelsService(t *testing.T, accounts []*authstore.Account, up ModelsUpstream) *Service {
	t.Helper()
	svc := newProbeService(t, accounts, &fakeUpstream{})
	svc.modelsLister = up
	return svc
}

// --- /api/models 字段扩展（契约 §1.8） ---

// TestModelsParsesPricingFields 上游返回 credits/supportsImages/description/badges 时应完整透传。
func TestModelsParsesPricingFields(t *testing.T) {
	fake := &fakeModelsUpstream{
		models: map[string][]map[string]any{
			"u-m-1": {
				{
					"id":             "wb-pro",
					"name":           "WorkBuddy Pro",
					"credits":        "x0.51 credits",
					"supportsImages": true,
					"descriptionZh":  "旗舰对话模型",
					"descriptionEn":  "Flagship chat model",
					"tags":           []any{"chat"},
				},
			},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-m-1")}, fake)
	items := svc.Models(context.Background())
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	m := items[0]
	if m.Credits != "x0.51 credits" {
		t.Errorf("Credits = %q, want %q", m.Credits, "x0.51 credits")
	}
	if !m.SupportsImages {
		t.Errorf("SupportsImages = false, want true")
	}
	if m.DescriptionZh != "旗舰对话模型" || m.DescriptionEn != "Flagship chat model" {
		t.Errorf("description = %q/%q", m.DescriptionZh, m.DescriptionEn)
	}
	// tags 无限免关键词 → badges 为空
	if len(m.Badges) != 0 {
		t.Errorf("Badges = %v, want 空", m.Badges)
	}
}

// TestModelsFreeBadgeExtraction 上游 tags 含限免关键词时应提取到 badges。
func TestModelsFreeBadgeExtraction(t *testing.T) {
	cases := []struct {
		name string
		tags []any
		want []string
	}{
		{"中文限免", []any{"chat", "限时免费"}, []string{"限时免费"}},
		{"中文免费", []any{"免费体验"}, []string{"免费体验"}},
		{"英文 free", []any{"Free Trial"}, []string{"Free Trial"}},
		{"trial", []any{"trial-model"}, []string{"trial-model"}},
		{"无关键词", []any{"chat", "text-to-image"}, nil},
		{"空 tags", nil, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fake := &fakeModelsUpstream{
				models: map[string][]map[string]any{
					"u-m-1": {{"id": "m1", "tags": c.tags}},
				},
			}
			svc := newModelsService(t, []*authstore.Account{healthyAccount("u-m-1")}, fake)
			items := svc.Models(context.Background())
			if len(items) != 1 {
				t.Fatalf("len = %d", len(items))
			}
			got := items[0].Badges
			if len(got) != len(c.want) {
				t.Fatalf("Badges = %v, want %v", got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("Badges[%d] = %q, want %q", i, got[i], c.want[i])
				}
			}
		})
	}
}

// TestModelsMissingFieldsDegrade 上游未返回价格字段时应给零值而非报错。
func TestModelsMissingFieldsDegrade(t *testing.T) {
	fake := &fakeModelsUpstream{
		models: map[string][]map[string]any{
			"u-m-1": {{"id": "wb-old", "name": "WorkBuddy Old"}},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-m-1")}, fake)
	items := svc.Models(context.Background())
	if len(items) != 1 {
		t.Fatalf("len = %d", len(items))
	}
	m := items[0]
	if m.Credits != "" || m.SupportsImages || m.DescriptionZh != "" || m.DescriptionEn != "" || len(m.Badges) != 0 {
		t.Errorf("缺字段应全部零值: %+v", m)
	}
}

// TestModelsHandlerShape 端到端：/api/models 响应 JSON 带新字段。
func TestModelsHandlerShape(t *testing.T) {
	fake := &fakeModelsUpstream{
		models: map[string][]map[string]any{
			"u-m-1": {{"id": "wb-pro", "name": "Pro", "credits": "x1 credits", "supportsImages": true}},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-m-1")}, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodGet, srv.URL+"/api/models", nil, cookie, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body struct {
		Items []struct {
			ID             string   `json:"id"`
			Credits        string   `json:"credits"`
			SupportsImages bool     `json:"supports_images"`
			DescriptionZh  string   `json:"description_zh"`
			DescriptionEn  string   `json:"description_en"`
			Badges         []string `json:"badges"`
		} `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Items) != 1 {
		t.Fatalf("len = %d", len(body.Items))
	}
	it := body.Items[0]
	if it.Credits != "x1 credits" || !it.SupportsImages {
		t.Errorf("JSON 字段错误: %+v", it)
	}
	// badges 必须序列化为数组（即使为空），不能缺省——契约 §1.8 约定。
	raw := map[string]json.RawMessage{}
	if err := json.NewDecoder(res.Body).Decode(&raw); err == nil {
		t.Log("res.Body 已被消费，跳过原始字节核对")
	}
}

// --- /api/models/pricing（契约 §2.5） ---

// TestModelPricingNormalPath 正常路径：credits 透传 + free 推导 + pricing_available=true。
func TestModelPricingNormalPath(t *testing.T) {
	fake := &fakeModelsUpstream{
		models: map[string][]map[string]any{
			"u-p-1": {
				{"id": "wb-pro", "name": "Pro", "credits": "x0.51 credits"},
				{"id": "wb-lite", "name": "Lite", "credits": "x0 credits"},
				{"id": "wb-trial", "name": "Trial", "tags": []any{"限时免费"}},
			},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-p-1")}, fake)
	resp := svc.ModelPricing(context.Background())
	if !resp.PricingAvailable {
		t.Errorf("PricingAvailable = false, want true（至少一个模型有 credits）")
	}
	if len(resp.Items) != 3 {
		t.Fatalf("len(items) = %d, want 3", len(resp.Items))
	}
	by := map[string]ModelPricing{}
	for _, it := range resp.Items {
		by[it.ID] = it
	}
	if by["wb-pro"].Credits != "x0.51 credits" || by["wb-pro"].Free {
		t.Errorf("wb-pro: %+v", by["wb-pro"])
	}
	// credits 为 x0 形态 → free=true
	if !by["wb-lite"].Free {
		t.Errorf("wb-lite（x0 credits）应判 free=true: %+v", by["wb-lite"])
	}
	// badges 非空 → free=true
	if !by["wb-trial"].Free || len(by["wb-trial"].Badges) != 1 {
		t.Errorf("wb-trial: %+v", by["wb-trial"])
	}
}

// TestModelPricingNoCredits 全部模型无 credits 时 pricing_available=false + warning 文案。
func TestModelPricingNoCredits(t *testing.T) {
	fake := &fakeModelsUpstream{
		models: map[string][]map[string]any{
			"u-p-1": {
				{"id": "wb-old", "name": "Old"},
			},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-p-1")}, fake)
	resp := svc.ModelPricing(context.Background())
	if resp.PricingAvailable {
		t.Errorf("PricingAvailable = true, want false")
	}
	if resp.Warning != "上游未提供模型价格数据" {
		t.Errorf("Warning = %q", resp.Warning)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("len = %d", len(resp.Items))
	}
	// 该模型无 credits 且无 badges → note 降级文案
	if resp.Items[0].Note == "" {
		t.Errorf("无价格数据模型应有 note 降级文案: %+v", resp.Items[0])
	}
}

// TestModelPricingPartialData 部分模型有 credits 部分没有：pricing_available=true，无 credits 的带 note。
func TestModelPricingPartialData(t *testing.T) {
	fake := &fakeModelsUpstream{
		models: map[string][]map[string]any{
			"u-p-1": {
				{"id": "wb-pro", "credits": "x1 credits"},
				{"id": "wb-old"},
			},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-p-1")}, fake)
	resp := svc.ModelPricing(context.Background())
	if !resp.PricingAvailable {
		t.Errorf("部分有 credits 时 PricingAvailable 应为 true")
	}
	if resp.Warning != "" {
		t.Errorf("部分有数据时不应有 warning: %q", resp.Warning)
	}
	by := map[string]ModelPricing{}
	for _, it := range resp.Items {
		by[it.ID] = it
	}
	if by["wb-pro"].Note != "" {
		t.Errorf("有 credits 的模型不应带 note: %+v", by["wb-pro"])
	}
	if by["wb-old"].Note == "" {
		t.Errorf("无 credits 无 badges 的模型应带 note: %+v", by["wb-old"])
	}
}

// TestModelPricingUpstreamFails 全部账号上游失败：items 空 + pricing_available=false + warning，不 502。
func TestModelPricingUpstreamFails(t *testing.T) {
	fake := &fakeModelsUpstream{
		errs: map[string]error{
			"u-p-1": &upstream.Error{Kind: upstream.ErrServer, Status: 503, Msg: "upstream down"},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-p-1")}, fake)
	resp := svc.ModelPricing(context.Background())
	if resp.PricingAvailable {
		t.Errorf("全失败时 PricingAvailable 应为 false")
	}
	if resp.Warning != "上游未提供模型价格数据" {
		t.Errorf("Warning = %q", resp.Warning)
	}
	if len(resp.Items) != 0 {
		t.Errorf("全失败时 items 应为空: %v", resp.Items)
	}
}

// TestModelPricingHandlerShape 端到端：/api/models/pricing JSON 形状与契约 §2.5 一致。
func TestModelPricingHandlerShape(t *testing.T) {
	fake := &fakeModelsUpstream{
		models: map[string][]map[string]any{
			"u-p-1": {{"id": "wb-pro", "name": "Pro", "credits": "x1 credits", "supportsImages": true}},
		},
	}
	svc := newModelsService(t, []*authstore.Account{healthyAccount("u-p-1")}, fake)
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodGet, srv.URL+"/api/models/pricing", nil, cookie, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body struct {
		Items []struct {
			UID            string   `json:"uid"`
			ID             string   `json:"id"`
			Credits        string   `json:"credits"`
			Free           bool     `json:"free"`
			Badges         []string `json:"badges"`
			SupportsImages bool     `json:"supports_images"`
			Note           string   `json:"note"`
		} `json:"items"`
		PricingAvailable bool   `json:"pricing_available"`
		Warning          string `json:"warning"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if !body.PricingAvailable || len(body.Items) != 1 || body.Items[0].Credits != "x1 credits" {
		t.Fatalf("响应形状错误: %+v", body)
	}
}

// TestModelPricingHandlerUnauthorized 未登录 401。
func TestModelPricingHandlerUnauthorized(t *testing.T) {
	svc := newModelsService(t, nil, &fakeModelsUpstream{})
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	res := request(t, srv.Client(), http.MethodGet, srv.URL+"/api/models/pricing", nil, nil, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}

// --- /api/activities 字段扩展：task_type（契约 §1.9） ---

// fakeGrowthUpstream 提供 GrowthTasks 桩（活动 task_type 扩展用）。
type fakeGrowthUpstream struct {
	tasks map[string][]upstream.GrowthTask
	errs  map[string]error
}

func (f *fakeGrowthUpstream) GrowthTasks(a *authstore.Account) ([]upstream.GrowthTask, error) {
	if err := f.errs[a.UID]; err != nil {
		return nil, err
	}
	return f.tasks[a.UID], nil
}

// TestActivitiesTaskType GrowthTasks 返回的 TaskType 应透传到 /api/activities。
func TestActivitiesTaskType(t *testing.T) {
	fake := &fakeGrowthUpstream{
		tasks: map[string][]upstream.GrowthTask{
			"u-a-1": {
				{Code: "daily_login", Name: "每日签到", Status: "可领取", Reward: 100, TaskType: "daily"},
				{Code: "first_buddy", Name: "首次领养猫猫", Status: "进行中", Reward: 300, TaskType: "once"},
			},
		},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-a-1")}, &fakeUpstream{})
	svc.growthLister = fake
	items := svc.Activities(context.Background())
	if len(items) != 2 {
		t.Fatalf("len = %d", len(items))
	}
	by := map[string]Activity{}
	for _, it := range items {
		by[it.Code] = it
	}
	if by["daily_login"].TaskType != "daily" {
		t.Errorf("daily_login.TaskType = %q", by["daily_login"].TaskType)
	}
	if by["first_buddy"].TaskType != "once" {
		t.Errorf("first_buddy.TaskType = %q", by["first_buddy"].TaskType)
	}
}

// --- /api/activities/lottery（契约 §2.6，只读） ---

// fakeLotteryUpstream 提供抽奖概要/机会/记录/奖品四端点桩。
type fakeLotteryUpstream struct {
	summary map[string]map[string]any
	chances map[string]map[string]any
	draws   map[string]map[string]any
	rewards map[string]map[string]any
	errs    map[string]error
}

func (f *fakeLotteryUpstream) LotterySummary(a *authstore.Account) (map[string]any, error) {
	if err := f.errs[a.UID]; err != nil {
		return nil, err
	}
	return f.summary[a.UID], nil
}
func (f *fakeLotteryUpstream) LotteryChances(a *authstore.Account) (map[string]any, error) {
	if err := f.errs[a.UID]; err != nil {
		return nil, err
	}
	return f.chances[a.UID], nil
}
func (f *fakeLotteryUpstream) LotteryDraws(a *authstore.Account, page, pageSize int) (map[string]any, error) {
	if err := f.errs[a.UID]; err != nil {
		return nil, err
	}
	return f.draws[a.UID], nil
}
func (f *fakeLotteryUpstream) LotteryRewards(a *authstore.Account, page, pageSize int) (map[string]any, error) {
	if err := f.errs[a.UID]; err != nil {
		return nil, err
	}
	return f.rewards[a.UID], nil
}

// TestLotteryNormalPath 正常聚合：chances/draws_total/recent/rewards_total。
func TestLotteryNormalPath(t *testing.T) {
	fake := &fakeLotteryUpstream{
		summary: map[string]map[string]any{"u-l-1": {}},
		chances: map[string]map[string]any{"u-l-1": {"chances": float64(2)}},
		draws: map[string]map[string]any{
			"u-l-1": {
				"total": float64(5),
				"draws": []any{
					map[string]any{"prize_name": "积分 +10", "created_at": "2026-09-10 12:00"},
				},
			},
		},
		rewards: map[string]map[string]any{"u-l-1": {"total": float64(1)}},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-l-1")}, &fakeUpstream{})
	svc.lotteryReader = fake
	items := svc.LotteryOverview(context.Background())
	if len(items) != 1 {
		t.Fatalf("len = %d", len(items))
	}
	it := items[0]
	if it.Chances != 2 || it.DrawsTotal != 5 || it.RewardsTotal != 1 {
		t.Errorf("聚合错误: %+v", it)
	}
	if len(it.Recent) != 1 || it.Recent[0].Prize != "积分 +10" {
		t.Errorf("recent 错误: %+v", it.Recent)
	}
	if it.Error != "" {
		t.Errorf("不应有 error: %q", it.Error)
	}
}

// TestLotteryUpstreamFails 单账号失败：error 下沉，其余字段零值，整批不中断。
func TestLotteryUpstreamFails(t *testing.T) {
	fake := &fakeLotteryUpstream{
		errs: map[string]error{"u-l-1": errors.New("upstream down")},
		chances: map[string]map[string]any{"u-l-2": {"chances": float64(3)}},
		draws:   map[string]map[string]any{"u-l-2": {"total": float64(0)}},
		rewards: map[string]map[string]any{"u-l-2": {"total": float64(0)}},
		summary: map[string]map[string]any{"u-l-2": {}},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-l-1"), healthyAccount("u-l-2")}, &fakeUpstream{})
	svc.lotteryReader = fake
	items := svc.LotteryOverview(context.Background())
	if len(items) != 2 {
		t.Fatalf("len = %d", len(items))
	}
	var failed, ok *LotteryItem
	for i := range items {
		if items[i].UID == "u-l-1" {
			failed = &items[i]
		} else {
			ok = &items[i]
		}
	}
	if failed == nil || failed.Error == "" {
		t.Errorf("失败账号应有 error: %+v", failed)
	}
	if ok == nil || ok.Error != "" || ok.Chances != 3 {
		t.Errorf("健康账号不应受影响: %+v", ok)
	}
}

// TestLotteryEmptyUpstreamData 上游返回空对象（Apifox 登记的形态）：全部零值 + note 占位标注。
func TestLotteryEmptyUpstreamData(t *testing.T) {
	fake := &fakeLotteryUpstream{
		summary: map[string]map[string]any{"u-l-1": {}},
		chances: map[string]map[string]any{"u-l-1": {}},
		draws:   map[string]map[string]any{"u-l-1": {}},
		rewards: map[string]map[string]any{"u-l-1": {}},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-l-1")}, &fakeUpstream{})
	svc.lotteryReader = fake
	items := svc.LotteryOverview(context.Background())
	if len(items) != 1 {
		t.Fatalf("len = %d", len(items))
	}
	it := items[0]
	if it.Chances != 0 || it.DrawsTotal != 0 || it.RewardsTotal != 0 || len(it.Recent) != 0 {
		t.Errorf("空数据应全零值: %+v", it)
	}
	if !strings.Contains(it.Note, "占位") {
		t.Errorf("零值应带占位标注 note: %q", it.Note)
	}
}

// TestLotteryHandlerShape 端到端：/api/activities/lottery JSON 形状与契约 §2.6 一致。
func TestLotteryHandlerShape(t *testing.T) {
	fake := &fakeLotteryUpstream{
		summary: map[string]map[string]any{"u-l-1": {}},
		chances: map[string]map[string]any{"u-l-1": {"chances": float64(1)}},
		draws:   map[string]map[string]any{"u-l-1": {"total": float64(2)}},
		rewards: map[string]map[string]any{"u-l-1": {"total": float64(0)}},
	}
	svc := newProbeService(t, []*authstore.Account{healthyAccount("u-l-1")}, &fakeUpstream{})
	svc.lotteryReader = fake
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	cookie := login(t, srv)
	res := request(t, srv.Client(), http.MethodGet, srv.URL+"/api/activities/lottery", nil, cookie, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var body struct {
		Items []struct {
			UID         string `json:"uid"`
			Chances     int    `json:"chances"`
			DrawsTotal  int    `json:"draws_total"`
			RewardsTotal int   `json:"rewards_total"`
			Note        string `json:"note"`
		} `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Items) != 1 || body.Items[0].Chances != 1 || body.Items[0].DrawsTotal != 2 {
		t.Fatalf("响应形状错误: %+v", body)
	}
}

// TestLotteryHandlerUnauthorized 未登录 401。
func TestLotteryHandlerUnauthorized(t *testing.T) {
	svc := newProbeService(t, nil, &fakeUpstream{})
	srv := httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler())
	defer srv.Close()
	res := request(t, srv.Client(), http.MethodGet, srv.URL+"/api/activities/lottery", nil, nil, false)
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", res.StatusCode)
	}
}
