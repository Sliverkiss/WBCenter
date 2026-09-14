package control

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// ProbeConcurrency 控制全量探测的同时在上游调用上限。
// 6 是「对 N 账号回源 ≠ 慢到不可接受」与「不把上游当成压测对象」之间的折中；
// 测试断言同时在飞的 UserResource 调用峰值 ≤ 该值。
const ProbeConcurrency = 6

// ProbeUpstream 是 Probe 所依赖的上游最小接口。
// 生产实现是 *upstream.Client，测试用桩替换以验证并发与降级语义。
type ProbeUpstream interface {
	UserResource(a *authstore.Account) (*upstream.Credits, error)
	TravelStatus(a *authstore.Account) (*upstream.TravelState, error)
	BuddyInfo(a *authstore.Account) (*upstream.Buddy, error)
}

// ProbeCheckin 签到状态子对象。
// 数据来源标注：上游没有「查签到」只读端点，此值来自面板本地自动化任务运行记录；
// 未运行过 checkin 时返回零值占位（见契约 §2.1 字段标注）。
type ProbeCheckin struct {
	CheckedIn  bool `json:"checked_in"`
	StreakDays int  `json:"streak_days"`
}

// ProbeCredits 探测时刻的积分快照。
type ProbeCredits struct {
	Current        int64 `json:"current"`
	TodayRemaining int64 `json:"today_remaining"`
}

// ProbeTravel 猫猫状态机。location/departed_at/arrive_at 上游未提供（见契约 §2.1 标注），
// 当前实现只回传 status 原值。
type ProbeTravel struct {
	Status     string `json:"status"`
	Location   string `json:"location,omitempty"`
	DepartedAt string `json:"departed_at,omitempty"`
	ArriveAt   string `json:"arrive_at,omitempty"`
}

// ProbeResult 单账号探测结果。契约：失败时 Error 必填、其余子对象缺省。
type ProbeResult struct {
	UID         string        `json:"uid"`
	Nickname    string        `json:"nickname"`
	OK          bool          `json:"ok"`
	SessionDead bool          `json:"session_dead,omitempty"`
	Checkin     *ProbeCheckin `json:"checkin,omitempty"`
	Credits     *ProbeCredits `json:"credits,omitempty"`
	Travel      *ProbeTravel  `json:"travel,omitempty"`
	Error       string        `json:"error,omitempty"`
}

// Probe 对每个账号并发拉取积分/签到状态/猫猫状态机并聚合。失败下沉到单账号 Error，
// 不中断整批；整体 ctx 取消时提前返回。
//
// 签到状态：契约已标注「上游无独立只读端点」，本实现仅返回占位（{false,0}），
// 待面板自动化任务运行记录接入后再真实填充——禁止编造上游字段（手册纪律）。
func (s *Service) Probe(ctx context.Context) ([]ProbeResult, error) {
	if s.cfg.ReadOnly {
		return nil, fmt.Errorf("服务端已开启只读模式")
	}
	rows, _ := s.store.List()
	results := make([]ProbeResult, len(rows))
	if len(rows) == 0 {
		return results, nil
	}
	// 信号量控制并发；ctx 取消时尽快退出。
	sem := make(chan struct{}, ProbeConcurrency)
	var wg sync.WaitGroup
	for i, a := range rows {
		i, a := i, a
		results[i] = ProbeResult{UID: a.UID, Nickname: a.Nickname}
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				results[i].Error = "探测已取消"
				return
			}
			s.probeOne(ctx, a, &results[i])
		}()
	}
	wg.Wait()
	// 输出顺序按 UID 稳定排序（聚合顺序与上游响应时序解耦，便于前端按 uid 索引）。
	sort.Slice(results, func(i, j int) bool { return results[i].UID < results[j].UID })
	s.mu.Lock()
	s.lastProbe = results
	s.mu.Unlock()
	failed, dead := 0, 0
	for _, r := range results {
		if r.Error != "" {
			failed++
		}
		if r.SessionDead {
			dead++
		}
	}
	level := "info"
	if failed > 0 {
		level = "warn"
	}
	s.Logf(level, "probe", "全量探测完成：%d 账号，失败 %d，会话失效 %d", len(results), failed, dead)
	return results, nil
}

// SessionDeadCount 返回最近一次探测被判 12153 的账号数（未运行过探测时为 0）。
// 供 /api/overview 卡片渲染；契约字段 session_dead。
func (s *Service) SessionDeadCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, r := range s.lastProbe {
		if r.SessionDead {
			n++
		}
	}
	return n
}

// prober 返回当前上游客户端；测试通过 svc.prober 替换为桩。
func (s *Service) proberOrDefault() ProbeUpstream {
	if s.prober != nil {
		return s.prober
	}
	return s.up
}

func (s *Service) probeOne(ctx context.Context, a *authstore.Account, out *ProbeResult) {
	up := s.proberOrDefault()
	// 任何一步失败都视为该账号整体失败（契约：ok=false 时其余子对象缺省）。
	credits, err := up.UserResource(a)
	if err != nil {
		markProbeError(out, err)
		return
	}
	out.Credits = &ProbeCredits{Current: credits.Remain, TodayRemaining: credits.Remain}

	travel, err := up.TravelStatus(a)
	if err != nil {
		markProbeError(out, err)
		return
	}
	// BuddyInfo 失败降级为不影响 travel 字段——无猫账号同样要看 travel 状态。
	_, buddyErr := up.BuddyInfo(a)
	_ = buddyErr

	out.Travel = &ProbeTravel{Status: travel.State}
	// 占位签到状态：见 Probe 注释。
	out.Checkin = &ProbeCheckin{CheckedIn: false, StreakDays: 0}
	out.OK = true
}

// markProbeError 把上游错误写到结果：12153 单独打 SessionDead 徽标，错误文案中文友好。
func markProbeError(out *ProbeResult, err error) {
	out.OK = false
	if upstream.IsSessionDead(err) {
		out.SessionDead = true
		out.Error = "上游会话失效（12153）"
		return
	}
	var ue *upstream.Error
	if errors.As(err, &ue) {
		out.Error = ue.Error()
		return
	}
	out.Error = err.Error()
}
