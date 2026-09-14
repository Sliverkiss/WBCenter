package control

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

type Account struct {
	UID          string `json:"uid"`
	Nickname     string `json:"nickname"`
	Domain       string `json:"domain"`
	ExpiresAt    int64  `json:"expires_at"`
	Expired      bool   `json:"expired"`
	NeedsRefresh bool   `json:"needs_refresh"`
}
type CreditSummary struct {
	UID            string    `json:"uid"`
	Nickname       string    `json:"nickname"`
	Current        int64     `json:"current"`
	TodayAllocated float64   `json:"today_allocated"`
	TodayConsumed  float64   `json:"today_consumed"`
	TodayRemaining float64   `json:"today_remaining"`
	Packages       int       `json:"packages"`
	FetchedAt      time.Time `json:"fetched_at"`
	Error          string    `json:"error,omitempty"`
}
type Activity struct {
	UID      string  `json:"uid"`
	Nickname string  `json:"nickname"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	Reward   float64 `json:"reward"`
	TaskType string  `json:"task_type"`
	New      bool    `json:"new"`
}
type Model struct {
	UID            string         `json:"uid"`
	Nickname       string         `json:"nickname"`
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Credits        string         `json:"credits"`
	SupportsImages bool           `json:"supports_images"`
	DescriptionZh  string         `json:"description_zh"`
	DescriptionEn  string         `json:"description_en"`
	Badges         []string       `json:"badges"`
	Raw            map[string]any `json:"raw,omitempty"`
	Error          string         `json:"error,omitempty"`
}

// ModelPricing 是 GET /api/models/pricing 的单条响应（契约 §2.5）。
// credits 描述串原样透传；Free 为面板推导值（badges 非空或 credits 为 x0 形态）。
type ModelPricing struct {
	UID            string   `json:"uid"`
	Nickname       string   `json:"nickname"`
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Credits        string   `json:"credits"`
	Free           bool     `json:"free"`
	Badges         []string `json:"badges"`
	SupportsImages bool     `json:"supports_images"`
	Note           string   `json:"note"`
}

// ModelPricingResponse 是 /api/models/pricing 的顶层响应。
type ModelPricingResponse struct {
	Items            []ModelPricing `json:"items"`
	PricingAvailable bool           `json:"pricing_available"`
	Warning          string         `json:"warning,omitempty"`
}

// LotteryDraw 是抽奖历史条目（契约 §2.6）。
type LotteryDraw struct {
	Prize string `json:"prize"`
	At    string `json:"at"`
}

// LotteryItem 是 GET /api/activities/lottery 的单账号聚合（契约 §2.6）。
// 上游字段未核实：多候选键提取失败给零值，Note 标注占位。
type LotteryItem struct {
	UID          string        `json:"uid"`
	Nickname     string        `json:"nickname"`
	Chances      int           `json:"chances"`
	DrawsTotal   int           `json:"draws_total"`
	Recent       []LotteryDraw `json:"recent"`
	RewardsTotal int           `json:"rewards_total"`
	Note         string        `json:"note"`
	Error        string        `json:"error,omitempty"`
}

// SchedulerTask belongs to a WorkBuddy account and represents its own cloud
// scheduler task.
type SchedulerTask struct {
	UID      string         `json:"uid"`
	Nickname string         `json:"nickname"`
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Status   string         `json:"status"`
	Updated  string         `json:"updated,omitempty"`
	Raw      map[string]any `json:"raw,omitempty"`
	Error    string         `json:"error,omitempty"`
}
type MockSchedulerTaskView struct {
	MockSchedulerTask
	Nickname string `json:"nickname"`
}
type loginFlow struct {
	Region     upstream.Region
	State, URL string
	Created    time.Time
}

type Service struct {
	cfg          Config
	store        *authstore.Store
	up           *upstream.Client
	state        *State
	mu           sync.Mutex
	accountLocks map[string]*sync.Mutex
	logins       map[string]loginFlow
	// prober 用于测试替换 Probe 的上游调用；生产为 nil，Probe 走 s.up。
	prober ProbeUpstream
	// batcher 用于测试替换 BatchActions 的上游调用；生产为 nil，走 s.up。
	batcher BatchUpstream
	// creditsFetcher 用于测试替换 StatsSummary 的积分聚合；生产为 nil，走 creditFor。
	creditsFetcher creditsFetcher
	// modelsLister 用于测试替换 Models/ModelPricing 的上游调用；生产为 nil，走 s.up。
	modelsLister ModelsUpstream
	// growthLister 用于测试替换 Activities 的上游调用；生产为 nil，走 s.up。
	growthLister GrowthUpstream
	// lotteryReader 用于测试替换 LotteryOverview 的上游调用；生产为 nil，走 s.up。
	lotteryReader LotteryUpstream
	// lastProbe 是最近一次 POST /api/probe 的聚合结果；供 /api/overview 的
	// session_dead 卡片读取，未运行过探测时为空。
	lastProbe []ProbeResult
}

func NewService(cfg Config, store *authstore.Store, up *upstream.Client, state *State) *Service {
	return &Service{cfg: cfg, store: store, up: up, state: state, accountLocks: map[string]*sync.Mutex{}, logins: map[string]loginFlow{}}
}
// ModelsUpstream 是 Models/ModelPricing 所依赖的上游最小接口。
// 生产实现是 *upstream.Client，测试用桩替换以验证逐账号降级语义。
type ModelsUpstream interface {
	AvailableModels(a *authstore.Account) ([]map[string]any, error)
}

// GrowthUpstream 是 Activities 所依赖的上游最小接口。
type GrowthUpstream interface {
	GrowthTasks(a *authstore.Account) ([]upstream.GrowthTask, error)
}

// LotteryUpstream 是 LotteryOverview 所依赖的上游最小接口（只读四端点）。
type LotteryUpstream interface {
	LotterySummary(a *authstore.Account) (map[string]any, error)
	LotteryChances(a *authstore.Account) (map[string]any, error)
	LotteryDraws(a *authstore.Account, page, pageSize int) (map[string]any, error)
	LotteryRewards(a *authstore.Account, page, pageSize int) (map[string]any, error)
}

func (s *Service) modelsListerOrDefault() ModelsUpstream {
	if s.modelsLister != nil {
		return s.modelsLister
	}
	return s.up
}
func (s *Service) growthListerOrDefault() GrowthUpstream {
	if s.growthLister != nil {
		return s.growthLister
	}
	return s.up
}
func (s *Service) lotteryReaderOrDefault() LotteryUpstream {
	if s.lotteryReader != nil {
		return s.lotteryReader
	}
	return s.up
}

func (s *Service) Config() Config { return s.cfg }
func (s *Service) Accounts() ([]Account, []string) {
	rows, warns := s.store.List()
	out := make([]Account, 0, len(rows))
	for _, a := range rows {
		out = append(out, Account{UID: a.UID, Nickname: a.Nickname, Domain: a.Domain, ExpiresAt: a.ExpiresAt, Expired: a.Expired(), NeedsRefresh: a.NeedsRefresh(24 * time.Hour)})
	}
	return out, warns
}
func (s *Service) auth(uid string) (*authstore.Account, error) { return s.store.Get(uid) }
func (s *Service) lock(uid string) func() {
	s.mu.Lock()
	l := s.accountLocks[uid]
	if l == nil {
		l = &sync.Mutex{}
		s.accountLocks[uid] = l
	}
	s.mu.Unlock()
	l.Lock()
	return l.Unlock
}

func (s *Service) Credits(ctx context.Context, uid string) []CreditSummary {
	rows, _ := s.store.List()
	out := make([]CreditSummary, 0, len(rows))
	for _, a := range rows {
		if uid != "" && uid != a.UID {
			continue
		}
		out = append(out, s.creditFor(ctx, a))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Nickname < out[j].Nickname })
	return out
}
func (s *Service) creditFor(ctx context.Context, a *authstore.Account) CreditSummary {
	r := CreditSummary{UID: a.UID, Nickname: a.Nickname, FetchedAt: time.Now()}
	credits, err := s.up.UserResource(a)
	if err != nil {
		r.Error = "上游积分概要查询失败"
		return r
	}
	r.Current = credits.Remain
	r.Packages = credits.Packages
	packages, err := s.up.DailyFreePackages(a, credits.PackageCodes)
	if err != nil {
		r.Error = "已取得当前积分；今日套餐明细暂不可用"
		return r
	}
	for _, p := range packages {
		r.TodayAllocated += p.Total
		r.TodayConsumed += p.Used
		r.TodayRemaining += p.Remaining
	}
	return r
}

func (s *Service) Activities(ctx context.Context) []Activity {
	rows, _ := s.store.List()
	up := s.growthListerOrDefault()
	out := []Activity{}
	for _, a := range rows {
		tasks, err := up.GrowthTasks(a)
		if err != nil {
			continue
		}
		keys := make([]string, 0, len(tasks))
		for _, t := range tasks {
			keys = append(keys, a.UID+":"+t.Code)
		}
		newTasks := s.state.MarkSeenBatch(a.UID, keys)
		for _, t := range tasks {
			key := a.UID + ":" + t.Code
			out = append(out, Activity{UID: a.UID, Nickname: a.Nickname, Code: t.Code, Name: t.Name, Status: t.Status, Reward: t.Reward, TaskType: t.TaskType, New: newTasks[key]})
		}
	}
	return out
}

// SchedulerTaskDetail reads a single WorkBuddy account scheduler task. It
// remains read-only; creation and changes are not inferred from an unverified
// write contract.
func (s *Service) SchedulerTaskDetail(ctx context.Context, uid, taskID string) (map[string]any, error) {
	a, err := s.auth(uid)
	if err != nil {
		return nil, err
	}
	return s.up.ConsoleTaskDetail(a, taskID)
}
func (s *Service) Models(ctx context.Context) []Model {
	rows, _ := s.store.List()
	up := s.modelsListerOrDefault()
	out := []Model{}
	for _, a := range rows {
		items, err := up.AvailableModels(a)
		if err != nil {
			out = append(out, Model{UID: a.UID, Nickname: a.Nickname, Error: "上游模型查询失败"})
			continue
		}
		for _, item := range items {
			out = append(out, Model{
				UID:            a.UID,
				Nickname:       a.Nickname,
				ID:             first(item, "id", "model_id", "modelId", "code"),
				Name:           first(item, "name", "display_name", "displayName"),
				Credits:        first(item, "credits"),
				SupportsImages: firstBool(item, "supportsImages", "supports_images"),
				DescriptionZh:  first(item, "descriptionZh", "description_zh", "desc"),
				DescriptionEn:  first(item, "descriptionEn", "description_en"),
				Badges:         extractFreeBadges(item),
				Raw:            item,
			})
		}
	}
	return out
}

// ModelPricing 聚合账号池内全部模型的 credits 描述串（契约 §2.5）。
// 降级语义：单账号失败跳过该账号；全部失败或全部模型无 credits 时
// PricingAvailable=false 且 Warning="上游未提供模型价格数据"。不编造数值单价。
func (s *Service) ModelPricing(ctx context.Context) ModelPricingResponse {
	rows, _ := s.store.List()
	up := s.modelsListerOrDefault()
	resp := ModelPricingResponse{Items: []ModelPricing{}}
	for _, a := range rows {
		items, err := up.AvailableModels(a)
		if err != nil {
			continue // 单账号失败按契约降级（该账号条目缺省），不 502
		}
		for _, item := range items {
			credits := first(item, "credits")
			badges := extractFreeBadges(item)
			p := ModelPricing{
				UID:            a.UID,
				Nickname:       a.Nickname,
				ID:             first(item, "id", "model_id", "modelId", "code"),
				Name:           first(item, "name", "display_name", "displayName"),
				Credits:        credits,
				Badges:         badges,
				SupportsImages: firstBool(item, "supportsImages", "supports_images"),
			}
			p.Free = len(badges) > 0 || isZeroCredits(credits)
			if credits == "" && len(badges) == 0 {
				p.Note = "该模型上游未提供价格与限免数据"
			}
			if credits != "" {
				resp.PricingAvailable = true
			}
			resp.Items = append(resp.Items, p)
		}
	}
	if !resp.PricingAvailable {
		resp.Warning = "上游未提供模型价格数据"
	}
	return resp
}

// LotteryOverview 聚合各账号抽奖概要（契约 §2.6，只读）。
// 上游 schema 在 Apifox 导出中为空占位，多候选键提取失败给零值并在 Note 标注。
func (s *Service) LotteryOverview(ctx context.Context) []LotteryItem {
	rows, _ := s.store.List()
	up := s.lotteryReaderOrDefault()
	out := make([]LotteryItem, 0, len(rows))
	for _, a := range rows {
		item := LotteryItem{UID: a.UID, Nickname: a.Nickname, Recent: []LotteryDraw{}}
		chances, errC := up.LotteryChances(a)
		draws, errD := up.LotteryDraws(a, 1, 5)
		rewards, errR := up.LotteryRewards(a, 1, 5)
		if errC != nil && errD != nil && errR != nil {
			item.Error = "上游抽奖概要查询失败"
			out = append(out, item)
			continue
		}
		if errC == nil {
			item.Chances = int(firstFloat(chances, "chances", "remaining", "count", "total"))
		}
		if errD == nil {
			item.DrawsTotal = int(firstFloat(draws, "total", "count"))
			item.Recent = extractLotteryDraws(draws)
		}
		if errR == nil {
			item.RewardsTotal = int(firstFloat(rewards, "total", "count"))
		}
		// 全部字段为推导零值 → 标注占位（上游 schema 未核实）。
		if item.Chances == 0 && item.DrawsTotal == 0 && item.RewardsTotal == 0 && len(item.Recent) == 0 {
			item.Note = "上游字段未核实的占位"
		}
		out = append(out, item)
	}
	return out
}

// extractLotteryDraws 从抽奖历史响应中提取最近记录（最多 5 条）。
func extractLotteryDraws(draws map[string]any) []LotteryDraw {
	rows := findObjectRowsLocal(draws, "draws", "logs", "items", "list", "records")
	out := make([]LotteryDraw, 0, len(rows))
	for i, row := range rows {
		if i >= 5 {
			break
		}
		prize := first(row, "prize_name", "prizeName", "prize", "name", "title", "reward_name", "rewardName")
		at := first(row, "created_at", "createdAt", "at", "time", "draw_time", "drawTime")
		if prize == "" && at == "" {
			continue
		}
		out = append(out, LotteryDraw{Prize: prize, At: at})
	}
	return out
}

// findObjectRowsLocal 在抽奖响应 map 中按候选键找对象数组。
func findObjectRowsLocal(m map[string]any, keys ...string) []map[string]any {
	for _, k := range keys {
		if raw, ok := m[k].([]any); ok {
			out := make([]map[string]any, 0, len(raw))
			for _, item := range raw {
				if row, ok := item.(map[string]any); ok {
					out = append(out, row)
				}
			}
			return out
		}
	}
	return nil
}

// freeBadgeKeywords 限免/活动徽标关键词（大小写不敏感子串匹配）。
var freeBadgeKeywords = []string{"免费", "限免", "free", "trial"}

// extractFreeBadges 从上游 tags/badges 数组中提取含限免关键词的条目，原样透出。
// 上游字段未核实显式限免标记，该提取是面板侧推导，非上游原话。
func extractFreeBadges(item map[string]any) []string {
	var raw any
	for _, k := range []string{"badges", "tags"} {
		if v, ok := item[k]; ok {
			raw = v
			break
		}
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil
	}
	out := []string{}
	for _, v := range arr {
		s, ok := v.(string)
		if !ok || s == "" {
			continue
		}
		lower := strings.ToLower(s)
		for _, kw := range freeBadgeKeywords {
			if strings.Contains(lower, kw) {
				out = append(out, s)
				break
			}
		}
	}
	return out
}

// isZeroCredits 判定 credits 描述串是否为「免费」形态（x0 / 0 credits 等）。
// 仅作 free 推导信号，不改变透传原文。
func isZeroCredits(credits string) bool {
	s := strings.ToLower(strings.TrimSpace(credits))
	if s == "" {
		return false
	}
	s = strings.TrimPrefix(s, "x")
	s = strings.TrimSpace(s)
	// 形如 "0 credits" / "0.0 credits" / "0" 视为免费。
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return false
	}
	if n, err := strconv.ParseFloat(fields[0], 64); err == nil {
		return n == 0
	}
	return false
}

// firstFloat 从上游响应 map 提取数值字段（多候选键，string/float64 宽容转换）。
// 与 upstream 包内同名私有函数同口径；抽奖 schema 未核实，数值可能以字符串形态出现。
func firstFloat(row map[string]any, keys ...string) float64 {
	for _, k := range keys {
		if v, ok := row[k]; ok {
			switch n := v.(type) {
			case float64:
				return n
			case int:
				return float64(n)
			case int64:
				return float64(n)
			case string:
				if f, err := strconv.ParseFloat(strings.TrimSpace(n), 64); err == nil {
					return f
				}
			}
		}
	}
	return 0
}

// firstBool 从上游模型条目提取 bool 字段（多候选键）。
func firstBool(row map[string]any, keys ...string) bool {
	for _, k := range keys {
		if v, ok := row[k]; ok {
			switch b := v.(type) {
			case bool:
				return b
			case string:
				return strings.EqualFold(b, "true") || b == "1"
			case float64:
				return b != 0
			}
		}
	}
	return false
}
func (s *Service) SchedulerTasks(ctx context.Context) []SchedulerTask {
	rows, _ := s.store.List()
	out := []SchedulerTask{}
	for _, a := range rows {
		items, err := s.up.ConsoleTasks(a)
		if err != nil {
			out = append(out, SchedulerTask{UID: a.UID, Nickname: a.Nickname, Error: schedulerReadError(err)})
			continue
		}
		for _, item := range items {
			out = append(out, SchedulerTask{UID: a.UID, Nickname: a.Nickname, ID: first(item, "task_id", "taskId", "id"), Name: first(item, "name", "title"), Status: first(item, "status", "state"), Updated: first(item, "updated_at", "updatedAt"), Raw: item})
		}
	}
	return out
}

func (s *Service) MockSchedulerTasks() []MockSchedulerTaskView {
	accounts, _ := s.Accounts()
	names := make(map[string]string, len(accounts))
	for _, account := range accounts {
		names[account.UID] = account.Nickname
	}
	tasks := s.state.MockSchedulerTasks()
	out := make([]MockSchedulerTaskView, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, MockSchedulerTaskView{MockSchedulerTask: task, Nickname: names[task.AccountUID]})
	}
	return out
}

func (s *Service) CreateMockSchedulerTask(accountUID, name, cron, prompt string, enabled bool) (MockSchedulerTaskView, error) {
	if _, err := s.auth(accountUID); err != nil {
		return MockSchedulerTaskView{}, err
	}
	task, err := s.state.CreateMockSchedulerTask(accountUID, name, cron, prompt, enabled)
	if err != nil {
		return MockSchedulerTaskView{}, err
	}
	account, _ := s.auth(accountUID)
	return MockSchedulerTaskView{MockSchedulerTask: task, Nickname: account.Nickname}, nil
}

func (s *Service) UpdateMockSchedulerTask(id, name, cron, prompt string, enabled bool) (MockSchedulerTaskView, error) {
	task, err := s.state.UpdateMockSchedulerTask(id, name, cron, prompt, enabled)
	if err != nil {
		return MockSchedulerTaskView{}, err
	}
	account, err := s.auth(task.AccountUID)
	if err != nil {
		return MockSchedulerTaskView{}, err
	}
	return MockSchedulerTaskView{MockSchedulerTask: task, Nickname: account.Nickname}, nil
}

func schedulerReadError(err error) string {
	var upstreamErr *upstream.Error
	if errors.As(err, &upstreamErr) {
		switch upstreamErr.Status {
		case 401:
			return "上游拒绝授权，请刷新凭据"
		case 403:
			return "上游拒绝访问，此账号没有云端任务权限"
		case 404:
			return "上游暂未提供云端任务接口"
		}
	}
	return "账号云端定时任务读取不可用"
}
func first(row map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := row[k]; ok {
			if x := fmt.Sprint(v); x != "" && x != "<nil>" {
				return x
			}
		}
	}
	return ""
}

func (s *Service) StartOAuth(region string) (string, string, error) {
	r, err := upstream.NormalizeRegion(region)
	if err != nil {
		return "", "", err
	}
	state, url, err := s.up.StartLogin(r)
	if err != nil {
		return "", "", err
	}
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	s.mu.Lock()
	s.logins[id] = loginFlow{Region: r, State: state, URL: url, Created: time.Now()}
	s.mu.Unlock()
	return id, url, nil
}
func (s *Service) PollOAuth(id string) (map[string]any, error) {
	s.mu.Lock()
	flow, ok := s.logins[id]
	s.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("授权会话不存在或已过期")
	}
	if time.Since(flow.Created) > 10*time.Minute {
		return nil, fmt.Errorf("授权会话已过期")
	}
	a, err := s.up.PollLogin(flow.Region, flow.State)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return map[string]any{"status": "pending"}, nil
	}
	if s.cfg.ReadOnly {
		return nil, fmt.Errorf("服务端已开启只读模式")
	}
	if err := s.store.Save(a); err != nil {
		return nil, err
	}
	s.mu.Lock()
	delete(s.logins, id)
	s.mu.Unlock()
	return map[string]any{"status": "success", "uid": a.UID, "nickname": a.Nickname}, nil
}

func (s *Service) Run(ctx context.Context, action, uid string) (string, error) {
	if s.cfg.ReadOnly {
		return "", fmt.Errorf("服务端已开启只读模式")
	}
	if action == "activity_probe" {
		items := s.Activities(ctx)
		return fmt.Sprintf("已探测 %d 条活动任务", len(items)), nil
	}
	rows, _ := s.store.List()
	targets := rows
	if uid != "" {
		a, err := s.auth(uid)
		if err != nil {
			return "", err
		}
		targets = []*authstore.Account{a}
	}
	ok, failed := 0, 0
	for _, a := range targets {
		unlock := s.lock(a.UID)
		var err error
		switch action {
		case "checkin":
			_, err = s.up.DailyCheckin(a)
		case "travel":
			_, err = s.up.TravelOnce(a)
		case "refresh":
			err = s.up.RefreshToken(a)
			if err == nil {
				err = s.store.Save(a)
			}
		default:
			unlock()
			return "", fmt.Errorf("未知动作")
		}
		unlock()
		if err != nil {
			failed++
		} else {
			ok++
		}
		time.Sleep(150 * time.Millisecond)
	}
	return fmt.Sprintf("成功 %d，失败 %d", ok, failed), nil
}
func (s *Service) Tick(ctx context.Context) {
	for _, a := range s.state.ClaimDue(time.Now()) {
		message, err := s.Run(ctx, a.Action, "")
		if err != nil {
			_ = s.state.Complete(a.ID, err.Error(), false)
		} else {
			_ = s.state.Complete(a.ID, message, true)
		}
	}
}
func (s *Service) State() *State { return s.state }
