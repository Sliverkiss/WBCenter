package control

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	UID            string                  `json:"uid"`
	Nickname       string                  `json:"nickname"`
	Current        int64                   `json:"current"`
	TodayAllocated float64                 `json:"today_allocated"`
	TodayConsumed  float64                 `json:"today_consumed"`
	TodayRemaining float64                 `json:"today_remaining"`
	Packages       int                     `json:"packages"`
	PackageDetails []upstream.DailyPackage `json:"package_details,omitempty"`
	FetchedAt      time.Time               `json:"fetched_at"`
	Error          string                  `json:"error,omitempty"`
}
type Activity struct {
	UID      string  `json:"uid"`
	Nickname string  `json:"nickname"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
	Target   float64 `json:"target"`
	Reward   float64 `json:"reward"`
	Action   string  `json:"action"`
	New      bool    `json:"new"`
	Error    string  `json:"error,omitempty"`
}

type ActivityRun struct {
	ID              string    `json:"id"`
	UID             string    `json:"uid"`
	Nickname        string    `json:"nickname"`
	Status          string    `json:"status"`
	Total           int       `json:"total"`
	Completed       int       `json:"completed"`
	Succeeded       int       `json:"succeeded"`
	Already         int       `json:"already"`
	Skipped         int       `json:"skipped"`
	Failed          int       `json:"failed"`
	Credit          int       `json:"credit"`
	Energy          int       `json:"energy"`
	StartedAt       time.Time `json:"started_at"`
	FinishedAt      time.Time `json:"finished_at,omitempty"`
	DurationSeconds int64     `json:"duration_seconds"`
	Summary         string    `json:"summary"`
	Logs            []string  `json:"logs"`
	doneCodes       map[string]bool
	names           map[string]string
}

var completableGrowthTasks = map[string]bool{
	"create_canvas": true, "template_5": true, "expert_5": true,
	"Expert_team_use_3": true, "skill_1": true, "automation_1": true,
	"playbook_prompt": true, "Expert_lighthouse": true, "Buddy_App": true,
	"Buddy_App_QQ": true, "Hp_Appearance": true, "chat_5": true,
	"Model_chat_GLM5.2": true, "black_cat": true, "RichMeow_Chat": true,
	"Library_read": true,
}

type Model struct {
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Error    string `json:"error,omitempty"`
}

// SchedulerTask belongs to a WorkBuddy account and represents its own cloud
// scheduler task.
type SchedulerTask struct {
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Updated  string `json:"updated,omitempty"`
	Error    string `json:"error,omitempty"`
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
	activityRuns map[string]*ActivityRun
}

func NewService(cfg Config, store *authstore.Store, up *upstream.Client, state *State) *Service {
	return &Service{cfg: cfg, store: store, up: up, state: state, accountLocks: map[string]*sync.Mutex{}, logins: map[string]loginFlow{}, activityRuns: map[string]*ActivityRun{}}
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
	r.PackageDetails = packages
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
			out = append(out, Activity{UID: a.UID, Nickname: a.Nickname, Error: activityReadError(err)})
			continue
		}
		keys := make([]string, 0, len(tasks))
		for _, t := range tasks {
			keys = append(keys, a.UID+":"+t.Code)
		}
		newTasks := s.state.MarkSeenBatch(a.UID, keys)
		for _, t := range tasks {
			key := a.UID + ":" + t.Code
			action := "manual"
			if completableGrowthTasks[t.Code] {
				action = "complete"
			}
			if strings.EqualFold(t.AcceptStatus, "claimed") {
				action = "done"
			}
			out = append(out, Activity{UID: a.UID, Nickname: a.Nickname, Code: t.Code, Name: t.Name, Status: t.Status, Progress: t.Progress, Target: t.Target, Reward: t.Reward, Action: action, New: newTasks[key] || s.state.IsRecentTask(key, 24*time.Hour)})
		}
	}
	return out
}

// StartActivityRun starts an observable background execution. Account and task
// codes are taken from the authenticated upstream task list and passed as argv
// without a shell.
func (s *Service) StartActivityRun(uid, code string) (ActivityRun, error) {
	if s.cfg.ReadOnly {
		return ActivityRun{}, fmt.Errorf("服务端已开启只读模式")
	}
	a, err := s.auth(uid)
	if err != nil {
		return ActivityRun{}, err
	}
	if code != "" && !completableGrowthTasks[code] {
		return ActivityRun{}, fmt.Errorf("该任务需要在 WorkBuddy 内人工完成")
	}
	runner := os.Getenv("WBCC_TASK_RUNNER")
	if runner == "" {
		runner = "/gateway-scripts/task_runner.py"
	}
	if info, err := os.Stat(runner); err != nil || info.IsDir() {
		return ActivityRun{}, fmt.Errorf("成长任务执行器不可用")
	}
	tasks, err := s.up.GrowthTasks(a)
	if err != nil {
		return ActivityRun{}, fmt.Errorf("读取成长任务失败")
	}
	codes, names := []string{}, map[string]string{}
	for _, task := range tasks {
		if !completableGrowthTasks[task.Code] || strings.EqualFold(task.AcceptStatus, "claimed") || (code != "" && task.Code != code) {
			continue
		}
		codes = append(codes, task.Code)
		names[task.Code] = task.Name
	}
	if len(codes) == 0 {
		return ActivityRun{}, fmt.Errorf("当前没有需要执行的可自动任务")
	}
	run := &ActivityRun{ID: fmt.Sprintf("activity-%d", time.Now().UnixNano()), UID: uid, Nickname: a.Nickname, Status: "running", Total: len(codes), StartedAt: time.Now(), Summary: "任务已开始", names: names, doneCodes: map[string]bool{}}
	s.mu.Lock()
	s.activityRuns[run.ID] = run
	for len(s.activityRuns) > 30 {
		var oldest string
		var oldestAt time.Time
		for id, item := range s.activityRuns {
			if oldest == "" || item.StartedAt.Before(oldestAt) {
				oldest, oldestAt = id, item.StartedAt
			}
		}
		delete(s.activityRuns, oldest)
	}
	s.mu.Unlock()
	go s.executeActivityRun(run.ID, runner, codes)
	return s.ActivityRun(run.ID)
}

func (s *Service) ActivityRun(id string) (ActivityRun, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	run, ok := s.activityRuns[id]
	if !ok {
		return ActivityRun{}, os.ErrNotExist
	}
	copy := *run
	copy.Logs = append([]string(nil), run.Logs...)
	copy.doneCodes, copy.names = nil, nil
	if copy.Status == "running" {
		copy.DurationSeconds = int64(time.Since(copy.StartedAt).Seconds())
	}
	return copy, nil
}

func (s *Service) ActivityRuns() []ActivityRun {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]ActivityRun, 0, len(s.activityRuns))
	for _, item := range s.activityRuns {
		copy := *item
		copy.Logs = append([]string(nil), item.Logs...)
		copy.doneCodes, copy.names = nil, nil
		if copy.Status == "running" {
			copy.DurationSeconds = int64(time.Since(copy.StartedAt).Seconds())
		}
		out = append(out, copy)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartedAt.After(out[j].StartedAt) })
	return out
}

var runnerSummaryRE = regexp.MustCompile(`accounts=(\d+) total=(\d+) ok=(\d+) already=(\d+) skipped=(\d+) pending=(\d+) fail=(\d+) credit=\+?(\d+) energy=\+?(\d+)`)
var runnerTaskRE = regexp.MustCompile(`^\[task_runner\]\s+\S+\s+(\S+):\s+(.*)$`)

func (s *Service) executeActivityRun(id, runner string, codes []string) {
	s.mu.Lock()
	run := s.activityRuns[id]
	s.mu.Unlock()
	args := []string{runner, run.UID, "--yes", "--gap", "1.1"}
	for _, code := range codes {
		args = append(args, "--only", code)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "python3", args...)
	cmd.Env = append(os.Environ(), "WB2A_AUTHS="+s.cfg.AuthDir, "PYTHONUNBUFFERED=1")
	unlock := s.lock(run.UID)
	defer unlock()
	pipe, err := cmd.StdoutPipe()
	if err == nil {
		cmd.Stderr = cmd.Stdout
		err = cmd.Start()
	}
	if err == nil {
		scanner := bufio.NewScanner(pipe)
		buf := make([]byte, 64*1024)
		scanner.Buffer(buf, 256*1024)
		for scanner.Scan() {
			s.updateActivityRun(id, scanner.Text())
		}
		err = cmd.Wait()
	}
	s.mu.Lock()
	run = s.activityRuns[id]
	run.FinishedAt = time.Now()
	run.DurationSeconds = int64(run.FinishedAt.Sub(run.StartedAt).Seconds())
	if err != nil {
		run.Status = "failed"
		run.Failed++
		run.Summary = fmt.Sprintf("执行中断：已处理 %d/%d 项，失败 %d 项，耗时 %s", run.Completed, run.Total, run.Failed, formatDuration(run.DurationSeconds))
	} else {
		run.Status = "completed"
		run.Completed = run.Total
		run.Summary = fmt.Sprintf("执行完成：共 %d 项，成功 %d 项，已领取 %d 项，跳过 %d 项，失败 %d 项；获得积分 %d、能量 %d；耗时 %s", run.Total, run.Succeeded, run.Already, run.Skipped, run.Failed, run.Credit, run.Energy, formatDuration(run.DurationSeconds))
	}
	s.mu.Unlock()
}

func (s *Service) updateActivityRun(id, raw string) {
	line := strings.TrimSpace(raw)
	if line == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	run := s.activityRuns[id]
	if m := runnerSummaryRE.FindStringSubmatch(line); len(m) == 10 {
		run.Succeeded, _ = strconv.Atoi(m[3])
		run.Already, _ = strconv.Atoi(m[4])
		run.Skipped, _ = strconv.Atoi(m[5])
		run.Failed, _ = strconv.Atoi(m[7])
		run.Credit, _ = strconv.Atoi(m[8])
		run.Energy, _ = strconv.Atoi(m[9])
		return
	}
	readable, code, terminal := readableRunnerLine(line, run.names)
	if readable != "" {
		if len(run.Logs) >= 500 {
			run.Logs = run.Logs[1:]
		}
		run.Logs = append(run.Logs, readable)
		run.Summary = readable
	}
	if terminal && !run.doneCodes[code] {
		run.doneCodes[code] = true
		run.Completed++
	}
}

func readableRunnerLine(line string, names map[string]string) (string, string, bool) {
	m := runnerTaskRE.FindStringSubmatch(line)
	if len(m) != 3 {
		return "", "", false
	}
	code, detail := m[1], m[2]
	name := names[code]
	if name == "" {
		name = code
	}
	prefix := "任务「" + name + "」："
	switch {
	case strings.HasPrefix(detail, "accept "):
		return prefix + "已接受任务", code, false
	case strings.HasPrefix(detail, "report ") && strings.Contains(detail, "code=0"):
		parts := strings.Fields(detail)
		progress := ""
		if len(parts) > 1 {
			progress = parts[1]
		}
		return prefix + "进度上报 " + progress + " 成功", code, false
	case strings.HasPrefix(detail, "claim ") && strings.Contains(detail, "already_claimed"):
		return prefix + "奖励此前已经领取", code, true
	case strings.HasPrefix(detail, "claim ") && strings.Contains(detail, " ok"):
		return prefix + "奖励领取成功", code, true
	case strings.Contains(detail, "已领，跳过"):
		return prefix + "奖励已领取，本次跳过", code, true
	case strings.Contains(detail, "不可伪造") || strings.Contains(detail, "非映射任务"):
		return prefix + "需要人工完成，本次跳过", code, true
	case strings.Contains(detail, "任务不存在"):
		return prefix + "当前账号没有该任务", code, true
	case strings.Contains(detail, "WARN 待下次"):
		return prefix + "本次未达到目标进度，请稍后重试", code, true
	case strings.Contains(detail, " -> ERR"):
		return prefix + "执行失败，请查看上游状态", code, true
	case strings.Contains(detail, "query re-read"):
		return prefix + "已回读并核对任务进度", code, false
	}
	return "", code, false
}

func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%d 秒", seconds)
	}
	return fmt.Sprintf("%d 分 %d 秒", seconds/60, seconds%60)
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
			out = append(out, Model{UID: a.UID, Nickname: a.Nickname, ID: first(item, "id", "model_id", "modelId", "code"), Name: first(item, "name", "display_name", "displayName")})
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
			out = append(out, SchedulerTask{UID: a.UID, Nickname: a.Nickname, ID: first(item, "task_id", "taskId", "id"), Name: first(item, "name", "title"), Status: first(item, "status", "state"), Updated: first(item, "updated_at", "updatedAt")})
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

func activityReadError(err error) string {
	var upstreamErr *upstream.Error
	if errors.As(err, &upstreamErr) {
		switch upstreamErr.Kind {
		case upstream.ErrSessionDead:
			return "活动读取失败，请刷新凭据或重新登录"
		case upstream.ErrSoftRate:
			return "活动读取受限，请稍后再试"
		case upstream.ErrNotFound:
			return "此账号暂不支持成长活动接口"
		}
	}
	return "上游活动任务查询失败"
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
		newCount := 0
		for _, item := range items {
			if item.New {
				newCount++
			}
		}
		return fmt.Sprintf("已探测 %d 条活动任务，近 24 小时新发现 %d 条", len(items), newCount), nil
	}
	if action == "activity_complete" {
		return s.runActivityAutomation(ctx)
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

func (s *Service) runActivityAutomation(ctx context.Context) (string, error) {
	accounts, _ := s.store.List()
	succeeded, skipped, failed, tasks, credit, energy := 0, 0, 0, 0, 0, 0
	started := time.Now()
	for _, account := range accounts {
		run, err := s.StartActivityRun(account.UID, "")
		if err != nil {
			if strings.Contains(err.Error(), "没有需要执行") {
				skipped++
				continue
			}
			failed++
			continue
		}
		for run.Status == "running" {
			select {
			case <-ctx.Done():
				return "", fmt.Errorf("活动任务自动完成已停止")
			case <-time.After(time.Second):
			}
			run, err = s.ActivityRun(run.ID)
			if err != nil {
				break
			}
		}
		if err != nil || run.Status == "failed" {
			failed++
			continue
		}
		succeeded++
		tasks += run.Total
		credit += run.Credit
		energy += run.Energy
	}
	return fmt.Sprintf("已处理 %d 个账号、%d 项活动任务；账号成功 %d，已无待办 %d，失败 %d；获得积分 %d、能量 %d；耗时 %s", len(accounts), tasks, succeeded, skipped, failed, credit, energy, formatDuration(int64(time.Since(started).Seconds()))), nil
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
