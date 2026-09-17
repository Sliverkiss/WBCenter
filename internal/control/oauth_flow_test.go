package control

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// oauthTestTransport 让 StartOAuth 不必真的打网络（命名与其他测试文件区分，避免冲突）。
type oauthTestTransport struct{}

func (oauthTestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body := `{"code":0,"data":{"state":"st-from-upstream","authUrl":"https://example.invalid/login?state=x"}}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    r,
	}, nil
}

func newOAuthTestService(t *testing.T) *Service {
	t.Helper()
	dir := t.TempDir()
	store, err := authstore.New(dir + "/auths")
	if err != nil {
		t.Fatal(err)
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
	up.HTTP = &http.Client{Transport: oauthTestTransport{}, Timeout: cfg.Timeout()}
	return NewService(cfg, store, up, state)
}

// TestPollOAuthReclaimsExpiredFlow 被轮询到「已过期」的授权流程必须同时被回收。
//
// logins 的条目只在**登录成功**时删除。于是「点开授权页又关掉」或「一直没完成」的流程
// 会永久留在 map 里：它承载着一个 state 与一个 authUrl，而每次发起授权都会新增一条，
// 面板进程生命周期内只增不减——纯粹的记账内存泄漏。
func TestPollOAuthReclaimsExpiredFlow(t *testing.T) {
	svc := newOAuthTestService(t)
	svc.logins["expired-flow"] = loginFlow{
		Region:  upstream.RegionCN,
		State:   "stale-state",
		URL:     "https://example.invalid/login?state=stale",
		Created: time.Now().Add(-2 * 10 * time.Minute),
	}

	if _, err := svc.PollOAuth("expired-flow"); err == nil {
		t.Fatal("过期流程应当报错")
	}
	if n := len(svc.logins); n != 0 {
		t.Errorf("过期授权流程未被回收：logins 仍有 %d 条（放弃的授权会永久留档 → 无界增长）", n)
	}
}

// TestStartOAuthReclaimsAbandonedFlows 被放弃、再也不会被轮询的流程，
// 应当在下次发起授权时被顺手清理（否则它们永远没有回收时机）。
func TestStartOAuthReclaimsAbandonedFlows(t *testing.T) {
	svc := newOAuthTestService(t)
	for i := 0; i < 50; i++ {
		svc.logins[fmt.Sprintf("abandoned-%d", i)] = loginFlow{
			Region:  upstream.RegionCN,
			State:   fmt.Sprintf("state-%d", i),
			URL:     "https://example.invalid/login?state=x",
			Created: time.Now().Add(-time.Hour),
		}
	}

	if _, _, err := svc.StartOAuth("cn"); err != nil {
		t.Fatalf("发起授权失败: %v", err)
	}
	if n := len(svc.logins); n != 1 {
		t.Errorf("被放弃的授权流程未被回收：logins = %d 条，期望只保留刚新建的 1 条", n)
	}
}

// TestStartOAuthKeepsFreshFlows 未过期的流程不能被误清（否则正在进行的授权会失败）。
func TestStartOAuthKeepsFreshFlows(t *testing.T) {
	svc := newOAuthTestService(t)
	svc.logins["fresh"] = loginFlow{Region: upstream.RegionCN, Created: time.Now()}

	if _, _, err := svc.StartOAuth("cn"); err != nil {
		t.Fatalf("发起授权失败: %v", err)
	}
	if _, ok := svc.logins["fresh"]; !ok {
		t.Error("未过期的授权流程被误删，正在进行的授权会中断")
	}
	if n := len(svc.logins); n != 2 {
		t.Errorf("logins = %d 条，期望 2 条（1 个在途 + 1 个新建）", n)
	}
}
