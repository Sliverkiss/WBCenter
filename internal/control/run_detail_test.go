package control

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// runDetailTransport 让 travel 动作在「不联网」的前提下走到各个分支（命名与其他测试文件区分）。
type runDetailTransport struct {
	buddyStatus int // 猫档案接口的 HTTP 状态码（0 视为 200）
	travelState string
}

func (t runDetailTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body := `{"code":0,"data":{}}`
	switch {
	case strings.HasSuffix(r.URL.Path, upstream.BuddyInfoPath):
		code := t.buddyStatus
		if code == 0 {
			code = http.StatusOK
		}
		payload := `{"code":0,"data":{"buddy":{"id":1,"name":"咪咪"}}}`
		if code != http.StatusOK {
			payload = `{"code":500,"msg":"buddy service unavailable"}`
		}
		return &http.Response{
			StatusCode: code,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(payload)),
			Request:    r,
		}, nil
	case strings.HasSuffix(r.URL.Path, upstream.TravelStatusPath):
		state := t.travelState
		if state == "" {
			state = "traveling"
		}
		body = fmt.Sprintf(`{"code":0,"data":{"state":%q}}`, state)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    r,
	}, nil
}

// newRunDetailService 建一个装着 accounts 个账号的 Service，上游由 transport 应答。
func newRunDetailService(t *testing.T, rt http.RoundTripper, accounts int) (*Service, string) {
	t.Helper()
	dir := t.TempDir()
	store, err := authstore.New(dir + "/auths")
	if err != nil {
		t.Fatal(err)
	}
	first := ""
	for i := 0; i < accounts; i++ {
		uid := fmt.Sprintf("account-%06d", i)
		if first == "" {
			first = uid
		}
		if err := store.Save(&authstore.Account{
			UID: uid, Nickname: fmt.Sprintf("账号%d", i),
			AccessToken: "test-access-token", RefreshToken: "test-refresh-token",
			ExpiresAt: time.Now().Add(time.Hour).Unix(), Domain: "https://example.invalid",
		}); err != nil {
			t.Fatal(err)
		}
	}
	state, err := NewState(dir + "/data")
	if err != nil {
		t.Fatal(err)
	}
	cfg := Config{
		Username: "admin", Password: "test-password",
		AuthDir: dir + "/auths", DataDir: dir + "/data",
		TimeoutSeconds: 5, Timezone: "Asia/Shanghai",
	}
	up := upstream.New(cfg.Timeout())
	up.HTTP = &http.Client{Transport: rt, Timeout: cfg.Timeout()}
	return NewService(cfg, store, up, state), first
}

// TestRunKeepsSkipDetail Run 必须把「跳过」也带回来。
//
// TravelOnce 在 traveling / 已派出 / 门槛未达标 等分支返回 (Action=skip, err=nil, Message=…)。
// 旧实现把整个 *TravelResult 丢掉（`_, err = s.up.TravelOnce(a)`），只统计 ok/failed，
// 于是运维在自动化历史里永远只看到「成功 1，失败 0」——一次 skip 与一次真的派出、
// 一次真的领奖，在面板上完全无法区分。
func TestRunKeepsSkipDetail(t *testing.T) {
	svc, uid := newRunDetailService(t, runDetailTransport{travelState: "traveling"}, 1)

	msg, err := svc.Run(context.Background(), "travel", uid)
	if err != nil {
		t.Fatalf("Run 失败: %v", err)
	}
	if !strings.HasPrefix(msg, "成功 1，失败 0") {
		t.Errorf("原有计数语义被破坏：message = %q", msg)
	}
	if !strings.Contains(msg, "猫正在旅行中") {
		t.Errorf("Run 丢掉了逐账号明细：message = %q，期望包含 %q", msg, "猫正在旅行中")
	}
}

// TestRunKeepsFailureReason Run 必须把失败原因带回来，否则「失败 N」无法排查。
func TestRunKeepsFailureReason(t *testing.T) {
	svc, uid := newRunDetailService(t, runDetailTransport{buddyStatus: http.StatusInternalServerError}, 1)

	msg, err := svc.Run(context.Background(), "travel", uid)
	if err != nil {
		t.Fatalf("Run 失败: %v", err)
	}
	if !strings.HasPrefix(msg, "成功 0，失败 1") {
		t.Errorf("原有计数语义被破坏：message = %q", msg)
	}
	if !strings.Contains(msg, "查询猫档案失败") {
		t.Errorf("Run 只报数量不报原因：message = %q，期望包含 %q", msg, "查询猫档案失败")
	}
}

// TestRunDetailMessageStaysBounded 明细不能把消息撑成无界字符串（它会写进状态文件）。
func TestRunDetailMessageStaysBounded(t *testing.T) {
	const n = 8
	svc, _ := newRunDetailService(t, runDetailTransport{buddyStatus: http.StatusInternalServerError}, n)

	msg, err := svc.Run(context.Background(), "travel", "") // uid 为空 = 全部账号
	if err != nil {
		t.Fatalf("Run 失败: %v", err)
	}
	if !strings.HasPrefix(msg, fmt.Sprintf("成功 0，失败 %d", n)) {
		t.Errorf("原有计数语义被破坏：message = %q", msg)
	}
	if len(msg) > 600 {
		t.Errorf("明细未做上限，消息长度 %d 字节：%q", len(msg), msg)
	}
	const limit = 5 // 与 Run 里的 detailLimit 对齐
	if !strings.Contains(msg, fmt.Sprintf("另有 %d 条明细未展开", n-limit)) {
		t.Errorf("被截断的明细条数没有交代：message = %q", msg)
	}
}
