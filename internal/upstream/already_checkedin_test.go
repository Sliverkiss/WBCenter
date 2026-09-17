package upstream

import (
	"errors"
	"net/http"
	"testing"
)

// TestIsAlreadyCheckedInRejectsUnrelatedFailures 只有「确实已经签过」才算幂等成功。
//
// IsAlreadyCheckedIn 的判定串来自 err.Error()，而 Error.Error() 会把上游原始 msg/body
// 拼进去。所以那些「只是提到了签到操作本身」或「只是碰巧含 already」的真实失败，
// 会被误判成「今日已签到」的幂等成功——DailyCheckin 直接当成成功返回，Run 也计成功：
// 用户看到「签到成功」，而账号 token 其实已经失效、或上游 5xx 根本没签到。
func TestIsAlreadyCheckedInRejectsUnrelatedFailures(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"token 过期（401，实为会话失效）", &Error{Kind: ErrSessionDead, Status: http.StatusUnauthorized, Msg: `{"msg":"token already expired"}`}},
		{"签到接口业务失败", &Error{Kind: ErrClient, Status: http.StatusBadRequest, Msg: "code=400 msg=checkin failed: invalid request"}},
		{"签到链路 5xx", &Error{Kind: ErrServer, Status: http.StatusInternalServerError, Msg: "code=500 msg=daily-checkin service unavailable"}},
		{"接口不存在", &Error{Kind: ErrNotFound, Status: http.StatusNotFound, Msg: "上游接口不存在（域名或路径可能已变更）"}},
		{"限流", &Error{Kind: ErrSoftRate, Status: http.StatusTooManyRequests, Msg: "too many requests"}},
		{"余额不足", &Error{Kind: ErrHardCredit, Status: http.StatusPaymentRequired, Msg: "积分不足，请充值"}},
		{"上游返回 HTML 错误页", &Error{Kind: ErrClient, Status: http.StatusBadRequest, Msg: "上游返回了非 JSON 的错误页"}},
		{"网络错误", errors.New("网络请求失败: context deadline exceeded")},
	}
	for _, c := range cases {
		if IsAlreadyCheckedIn(c.err) {
			t.Errorf("%s: 真实失败被误判成「已签到」幂等成功 → 用户会看到「签到成功」：%v", c.name, c.err)
		}
	}
}

// TestIsAlreadyCheckedInKeepsIdempotentSuccesses 明确表达「已经签过」的措辞必须仍然识别。
func TestIsAlreadyCheckedInKeepsIdempotentSuccesses(t *testing.T) {
	for _, s := range []string{
		"code=10001 msg=今天已签到",
		"用户已签到",
		"今天已经签到过啦",
		"already checked in",
		"already checkin today",
		"daily checkin done",
	} {
		if !IsAlreadyCheckedIn(errors.New(s)) {
			t.Errorf("应识别为已签到的幂等成功: %q", s)
		}
	}
}
