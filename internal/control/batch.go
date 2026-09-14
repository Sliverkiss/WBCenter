package control

import (
	"context"
	"fmt"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// BatchUpstream 是 BatchActions 所依赖的上游最小接口。
// 生产实现是 *upstream.Client，测试用桩替换以验证跳过/隔离语义。
type BatchUpstream interface {
	DailyCheckin(a *authstore.Account) (*upstream.CheckinResult, error)
	TravelOnce(a *authstore.Account) (*upstream.TravelResult, error)
	RefreshToken(a *authstore.Account) error
}

// BatchActionResult 单账号批量动作结果。契约 §2.3：
// 失败时 Error 必填；禁用时 Skipped=true 且 OK=true；12153 打 SessionDead 标记。
type BatchActionResult struct {
	UID         string `json:"uid"`
	Nickname    string `json:"nickname"`
	OK          bool   `json:"ok"`
	Message     string `json:"message,omitempty"`
	Skipped     bool   `json:"skipped,omitempty"`
	SessionDead bool   `json:"session_dead,omitempty"`
	Error       string `json:"error,omitempty"`
}

// batchActionDelay 是逐账号执行之间的节流间隔，与 Run 保持一致（上游友好）。
const batchActionDelay = 150 * time.Millisecond

// batcher 返回当前上游客户端；测试通过 svc.batcher 替换为桩。
func (s *Service) batcherOrDefault() BatchUpstream {
	if s.batcher != nil {
		return s.batcher
	}
	return s.up
}

// BatchActions 对账号池批量执行指定动作（checkin/travel/refresh），逐账号聚合结果。
// 语义对照契约 §2.3 与 PR #60：
//   - 禁用账号跳过（不调上游），skipped=true 且 ok=true；
//   - 单账号失败不中断整批，error 下沉到 results[].error；
//   - 12153 会话失效打 session_dead 标记；
//   - uids 为空 = 全量；不存在的 uid 静默忽略；
//   - refresh 成功后回写凭据文件。
//
// 逐账号串行执行并复用 per-uid 锁（与单账号动作/自动化互斥），账号间保留节流间隔。
func (s *Service) BatchActions(ctx context.Context, action string, uids []string) ([]BatchActionResult, error) {
	if s.cfg.ReadOnly {
		return nil, fmt.Errorf("服务端已开启只读模式")
	}
	if action != "checkin" && action != "travel" && action != "refresh" {
		return nil, fmt.Errorf("不支持的动作")
	}
	rows, _ := s.store.List()
	wanted := make(map[string]bool, len(uids))
	for _, uid := range uids {
		wanted[uid] = true
	}
	up := s.batcherOrDefault()
	out := make([]BatchActionResult, 0, len(rows))
	for _, a := range rows {
		if len(uids) > 0 && !wanted[a.UID] {
			continue
		}
		out = append(out, s.batchOne(ctx, up, action, a))
		if ctx.Err() != nil {
			break
		}
	}
	return out, nil
}

func (s *Service) batchOne(ctx context.Context, up BatchUpstream, action string, a *authstore.Account) BatchActionResult {
	r := BatchActionResult{UID: a.UID, Nickname: a.Nickname}
	if a.Disabled {
		r.OK = true
		r.Skipped = true
		r.Message = "账号已禁用，已跳过"
		return r
	}
	unlock := s.lock(a.UID)
	defer unlock()
	var err error
	switch action {
	case "checkin":
		_, err = up.DailyCheckin(a)
		r.Message = "签到成功"
	case "travel":
		_, err = up.TravelOnce(a)
		r.Message = "旅行巡检完成"
	case "refresh":
		err = up.RefreshToken(a)
		if err == nil {
			err = s.store.Save(a)
		}
		r.Message = "凭据已刷新"
	}
	if err != nil {
		r.OK = false
		r.Message = ""
		markBatchError(&r, err)
		return r
	}
	r.OK = true
	// 节流：与 Run 的 150ms 一致，避免把上游当压测对象。
	select {
	case <-time.After(batchActionDelay):
	case <-ctx.Done():
	}
	return r
}

// markBatchError 与 probe 的 markProbeError 同口径：12153 单独打 SessionDead，文案中文友好。
func markBatchError(r *BatchActionResult, err error) {
	if upstream.IsSessionDead(err) {
		r.SessionDead = true
		r.Error = "上游会话失效（12153）"
		return
	}
	r.Error = err.Error()
}
