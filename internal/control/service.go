package control

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

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
	New      bool    `json:"new"`
}
type Model struct {
	UID      string         `json:"uid"`
	Nickname string         `json:"nickname"`
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Raw      map[string]any `json:"raw,omitempty"`
	Error    string         `json:"error,omitempty"`
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
}

func NewService(cfg Config, store *authstore.Store, up *upstream.Client, state *State) *Service {
	return &Service{cfg: cfg, store: store, up: up, state: state, accountLocks: map[string]*sync.Mutex{}, logins: map[string]loginFlow{}}
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
	out := []Activity{}
	for _, a := range rows {
		tasks, err := s.up.GrowthTasks(a)
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
			out = append(out, Activity{UID: a.UID, Nickname: a.Nickname, Code: t.Code, Name: t.Name, Status: t.Status, Reward: t.Reward, New: newTasks[key]})
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
	out := []Model{}
	for _, a := range rows {
		items, err := s.up.AvailableModels(a)
		if err != nil {
			out = append(out, Model{UID: a.UID, Nickname: a.Nickname, Error: "上游模型查询失败"})
			continue
		}
		for _, item := range items {
			out = append(out, Model{UID: a.UID, Nickname: a.Nickname, ID: first(item, "id", "model_id", "modelId", "code"), Name: first(item, "name", "display_name", "displayName"), Raw: item})
		}
	}
	return out
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
	// detailLimit 限制明细条数：这条消息会被写进状态文件（automation 的 LastResult 与
	// Runs 历史），账号多的时候不能把它撑成无界字符串。
	const detailLimit = 5
	details := make([]string, 0, detailLimit)
	dropped := 0
	addDetail := func(a *authstore.Account, text string) {
		if text == "" {
			return
		}
		if len(details) >= detailLimit {
			dropped++
			return
		}
		name := strings.TrimSpace(a.Nickname)
		if name == "" {
			name = a.UID
		}
		details = append(details, name+"："+clipMessage(text))
	}
	for _, a := range targets {
		unlock := s.lock(a.UID)
		var err error
		note := ""
		switch action {
		case "checkin":
			var res *upstream.CheckinResult
			res, err = s.up.DailyCheckin(a)
			if res != nil {
				if res.Already {
					note = "今日已签到（幂等跳过）"
				} else {
					note = res.Message
				}
			}
		case "travel":
			var res *upstream.TravelResult
			res, err = s.up.TravelOnce(a)
			if res != nil {
				note = res.Message
			}
		case "refresh":
			err = s.up.RefreshToken(a)
			if err == nil {
				err = s.store.Save(a)
			}
			if err == nil {
				note = "凭据已刷新"
			}
		default:
			unlock()
			return "", fmt.Errorf("未知动作")
		}
		unlock()
		if err != nil {
			failed++
			// 明细优先用上游给的人类可读描述（如「查询猫档案失败: …」），没有时退回错误
			// 文本——原本两条都丢，运维只能看到一个数字，无从排查。
			if note == "" {
				note = err.Error()
			}
			addDetail(a, "失败："+note)
		} else {
			ok++
			addDetail(a, note)
		}
		time.Sleep(150 * time.Millisecond)
	}
	msg := fmt.Sprintf("成功 %d，失败 %d", ok, failed)
	if len(details) > 0 {
		msg += "；" + strings.Join(details, "；")
	}
	if dropped > 0 {
		msg += fmt.Sprintf("…（另有 %d 条明细未展开）", dropped)
	}
	return msg, nil
}

// clipMessage 限制单条明细的长度，按字符边界截断，避免把中文劈成半个字
// （半个字会产出非法 UTF-8，写进响应 JSON 时变成乱码）。
func clipMessage(s string) string {
	const max = 80
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) <= max {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); {
		_, w := utf8.DecodeRuneInString(s[i:])
		if b.Len()+w > max-3 { // 留 3 字节给省略号
			break
		}
		b.WriteString(s[i : i+w])
		i += w
	}
	return b.String() + "…"
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
