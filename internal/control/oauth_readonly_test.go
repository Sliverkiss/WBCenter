package control

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// oauthGateTransport 记录上游是否真的被调用过（不打网络）。
type oauthGateTransport struct {
	mu    sync.Mutex
	calls int
}

func (t *oauthGateTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.mu.Lock()
	t.calls++
	t.mu.Unlock()
	body := `{"code":0,"data":{"state":"st-from-upstream","authUrl":"https://example.invalid/login?state=x"}}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    r,
	}, nil
}

func (t *oauthGateTransport) count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.calls
}

// newOAuthGateServer 起一个只会被记录、不会真连上游的面板，用于验证只读门禁。
func newOAuthGateServer(t *testing.T, readOnly bool) (*httptest.Server, *oauthGateTransport) {
	t.Helper()
	dir := t.TempDir()
	if _, err := authstore.New(dir + "/auths"); err != nil {
		t.Fatal(err)
	}
	state, err := NewState(dir + "/data")
	if err != nil {
		t.Fatal(err)
	}
	store, err := authstore.New(dir + "/auths")
	if err != nil {
		t.Fatal(err)
	}
	cfg := Config{
		Username: "admin", Password: "test-password",
		AuthDir: dir + "/auths", DataDir: dir + "/data",
		TimeoutSeconds: 5, Timezone: "Asia/Shanghai", ReadOnly: readOnly,
	}
	rt := &oauthGateTransport{}
	up := upstream.New(cfg.Timeout())
	up.HTTP = &http.Client{Transport: rt, Timeout: cfg.Timeout()}
	svc := NewService(cfg, store, up, state)
	return httptest.NewServer(NewServer(svc, http.NotFoundHandler()).Handler()), rt
}

// TestOAuthStartBlockedInReadOnlyMode 只读模式必须在「添加账号」的**入口**就挡住。
//
// 面板把只读承诺暴露给了前端（/api/session 的 read_only），其余 6 个写入端点也都在
// 处理函数开头返回 403，只有 /api/oauth/start 例外。后果不是「多一次点击」，而是：
// 用户完整走完 WorkBuddy 网页授权（打开上游授权页、同意、登录），凭据真的换回来了，
// 走到最后一步 PollOAuth 才被告知「服务端已开启只读模式」——用户白跑一趟，
// 而且上游侧已经留下了这次授权的痕迹。
func TestOAuthStartBlockedInReadOnlyMode(t *testing.T) {
	ts, rt := newOAuthGateServer(t, true)
	defer ts.Close()
	cookie := login(t, ts)

	res := request(t, ts.Client(), http.MethodPost, ts.URL+"/api/oauth/start", map[string]string{"region": "cn"}, cookie, true)
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("只读模式下 /api/oauth/start 状态码 = %d，期望 403；响应体：%s", res.StatusCode, strings.TrimSpace(string(raw)))
	}
	if n := rt.count(); n != 0 {
		t.Errorf("只读模式下仍然调用了上游 %d 次：只读承诺没有覆盖「添加账号」入口", n)
	}
}

// TestOAuthStartWorksWhenWritable 守卫：可写模式下入口必须照常工作（防止门禁写过头）。
func TestOAuthStartWorksWhenWritable(t *testing.T) {
	ts, rt := newOAuthGateServer(t, false)
	defer ts.Close()
	cookie := login(t, ts)

	res := request(t, ts.Client(), http.MethodPost, ts.URL+"/api/oauth/start", map[string]string{"region": "cn"}, cookie, true)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(res.Body)
		t.Fatalf("可写模式下 /api/oauth/start 状态码 = %d，期望 200；响应体：%s", res.StatusCode, strings.TrimSpace(string(raw)))
	}
	var out struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out.ID == "" || out.URL == "" {
		t.Errorf("可写模式下没有返回授权信息：%#v", out)
	}
	if n := rt.count(); n != 1 {
		t.Errorf("可写模式下上游调用次数 = %d，期望 1", n)
	}
}
