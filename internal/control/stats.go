package control

import (
	"context"
	"time"

	"workbuddy-control-center/internal/authstore"
)

// StatsSummary 是 GET /api/stats/summary 的响应结构（契约 §2.2）。
// JSON tag 与契约逐字一致。
type StatsSummary struct {
	AccountsTotal              int       `json:"accounts_total"`
	CreditsCurrentTotal        int64     `json:"credits_current_total"`
	CreditsTodayAllocatedTotal float64   `json:"credits_today_allocated_total"`
	CreditsTodayConsumedTotal  float64   `json:"credits_today_consumed_total"`
	CreditsTodayRemainingTotal float64   `json:"credits_today_remaining_total"`
	CheckinDoneToday           int       `json:"checkin_done_today"`
	CheckinPendingToday        int       `json:"checkin_pending_today"`
	StatusHealthy              int       `json:"status_healthy"`
	StatusExpired              int       `json:"status_expired"`
	StatusSessionDead          int       `json:"status_session_dead"`
	StatusDisabled             int       `json:"status_disabled"`
	AutomationRunsTodayOK      int       `json:"automation_runs_today_ok"`
	AutomationRunsTodayFailed  int       `json:"automation_runs_today_failed"`
	GeneratedAt                time.Time `json:"generated_at"`
}

// creditsFor 抽离单账号积分聚合，供 /api/credits 与 stats 复用；
// 测试通过 creditsFetcher 字段替换，避免 httptest 起真实上游。
type creditsFetcher interface {
	creditsFor(ctx context.Context, a *authstore.Account) CreditSummary
}

// serviceCreditsFetcher 是生产实现：复用 Service.creditFor（UserResource + DailyFreePackages）。
type serviceCreditsFetcher struct{ s *Service }

func (f serviceCreditsFetcher) creditsFor(ctx context.Context, a *authstore.Account) CreditSummary {
	return f.s.creditFor(ctx, a)
}

func (s *Service) creditsFetcherOrDefault() creditsFetcher {
	if s.creditsFetcher != nil {
		return s.creditsFetcher
	}
	return serviceCreditsFetcher{s}
}

// StatsSummary 聚合统计（契约 §2.2）。缓存策略：
//   - 积分四项：按需实时聚合（与 /api/credits 同口径），单账号失败按 0 计入不拖垮整体；
//   - 签到计数与 12153 分布：复用 lastProbe 缓存（未探测时保守零值，不主动触发上游）；
//   - 自动化今日统计：本地 RunRecord 历史按服务端时区的自然日过滤。
func (s *Service) StatsSummary(ctx context.Context) StatsSummary {
	accounts, _ := s.Accounts()
	fetcher := s.creditsFetcherOrDefault()
	rows, _ := s.store.List()

	sum := StatsSummary{AccountsTotal: len(accounts), GeneratedAt: time.Now()}
	for _, a := range accounts {
		if a.Expired {
			sum.StatusExpired++
		} else {
			sum.StatusHealthy++
		}
	}
	for _, a := range rows {
		if a.Disabled {
			sum.StatusDisabled++
		}
		c := fetcher.creditsFor(ctx, a)
		if c.Error != "" {
			continue // 失败账号积分按 0 计入
		}
		sum.CreditsCurrentTotal += c.Current
		sum.CreditsTodayAllocatedTotal += c.TodayAllocated
		sum.CreditsTodayConsumedTotal += c.TodayConsumed
		sum.CreditsTodayRemainingTotal += c.TodayRemaining
	}

	// 签到计数与 12153：复用 lastProbe 缓存。
	s.mu.Lock()
	probe := s.lastProbe
	s.mu.Unlock()
	done := 0
	for _, r := range probe {
		if r.SessionDead {
			sum.StatusSessionDead++
		}
		if r.OK && r.Checkin != nil && r.Checkin.CheckedIn {
			done++
		}
	}
	if len(probe) > 0 {
		sum.CheckinDoneToday = done
		sum.CheckinPendingToday = len(accounts) - done
		if sum.CheckinPendingToday < 0 {
			sum.CheckinPendingToday = 0
		}
	} else {
		sum.CheckinDoneToday = 0
		sum.CheckinPendingToday = len(accounts)
	}

	// 自动化今日统计：本地 RunRecord，按服务端配置时区的自然日过滤。
	loc := time.Local
	if tz := s.cfg.Timezone; tz != "" {
		if l, err := time.LoadLocation(tz); err == nil {
			loc = l
		}
	}
	now := time.Now().In(loc)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	for _, r := range s.state.Runs() {
		if r.At.Before(dayStart) {
			continue
		}
		if r.OK {
			sum.AutomationRunsTodayOK++
		} else {
			sum.AutomationRunsTodayFailed++
		}
	}
	return sum
}
