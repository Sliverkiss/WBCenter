package upstream

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// TestNewConfiguresProxyFromEnvironment 传输层必须保留 HTTP_PROXY 支持与 HTTP/2。
//
// New() 手搓 &http.Transport{...}，这不会继承 http.DefaultTransport 的默认值。最要紧的
// 那个遗漏是 Proxy（默认值是 http.ProxyFromEnvironment）：少了它，面板在需要走代理的
// 环境（企业出口代理、容器侧车代理）里根本连不上上游，而且不报配置错——只表现为连不通。
// 同一批丢失的默认值还包括 ForceAttemptHTTP2（HTTP/2 被关掉）与 TLSHandshakeTimeout
// （握手不再有超时）。
func TestNewConfiguresProxyFromEnvironment(t *testing.T) {
	c := New(5 * time.Second)
	tr, ok := c.HTTP.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("HTTP.Transport 不是 *http.Transport：%T", c.HTTP.Transport)
	}
	if tr.Proxy == nil {
		t.Error("New() 的传输层没有配置 Proxy：HTTP_PROXY/HTTPS_PROXY/NO_PROXY 全部被忽略")
	}
	if !tr.ForceAttemptHTTP2 {
		t.Error("New() 的传输层没有开启 ForceAttemptHTTP2：HTTP/2 会被关掉")
	}
	if tr.TLSHandshakeTimeout == 0 {
		t.Error("New() 的传输层没有 TLSHandshakeTimeout：TLS 握手没有超时")
	}
	// 克隆默认传输层不能把面板特意调过的连接池参数丢掉（账号多时要复用连接）。
	if tr.MaxIdleConnsPerHost != 10 || tr.MaxIdleConns != 50 {
		t.Errorf("连接池参数丢失：MaxIdleConns=%d MaxIdleConnsPerHost=%d，期望 50/10",
			tr.MaxIdleConns, tr.MaxIdleConnsPerHost)
	}
}

// TestNewHonorsHTTPProxyEnv 用子进程做行为验证：设置 HTTP_PROXY 后请求必须真的走代理。
//
// 必须在子进程里跑：http.ProxyFromEnvironment 对当前进程的环境变量只解析一次并缓存，
// 同进程内改环境变量会受执行顺序影响，测试不稳定；子进程的环境是全新的，结果确定。
func TestNewHonorsHTTPProxyEnv(t *testing.T) {
	if os.Getenv(proxyChildMarker) == "1" {
		t.Skip("本用例由父测试驱动")
	}
	var hits int64
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "proxied")
	}))
	defer proxy.Close()

	cmd := exec.Command(os.Args[0], "-test.run=TestNewHonorsHTTPProxyEnvChild", "-test.v")
	cmd.Env = proxyChildEnv(proxy.URL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("子进程失败: %v\n%s", err, out)
	}
	// 防「子进程跳过了、所以当然没命中代理」的假通过：必须确认它真的发起了请求。
	if !strings.Contains(string(out), proxyChildAttempt) {
		t.Fatalf("子进程没有执行到请求路径，本次结果不可信:\n%s", out)
	}
	if atomic.LoadInt64(&hits) == 0 {
		t.Error("设置 HTTP_PROXY 后请求没有经过代理：New() 的传输层忽略了环境变量代理设置")
	}
}

const (
	proxyChildMarker = "WBCC_PROXY_ENV_CHILD"
	// proxyChildAttempt 由子进程打在自己的输出里，供父测试确认它真的走到了请求那一步。
	proxyChildAttempt = "PROXY_CHILD_ATTEMPTED"
)

// TestNewHonorsHTTPProxyEnvChild 由父测试在子进程里执行。
func TestNewHonorsHTTPProxyEnvChild(t *testing.T) {
	if os.Getenv(proxyChildMarker) != "1" {
		t.Skip("仅在子进程里执行")
	}
	c := New(5 * time.Second)
	// .invalid 是 RFC 2606 保留域名，永远不解析：只有真的走了代理才可能拿到响应。
	resp, err := c.HTTP.Get("http://upstream-proxy-probe.invalid/health")
	if err == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}
	t.Logf("%s err=%v", proxyChildAttempt, err)
}

// proxyChildEnv 构造子进程环境：清掉继承来的所有代理变量，只留测试用的 HTTP_PROXY。
func proxyChildEnv(proxyURL string) []string {
	env := make([]string, 0, len(os.Environ())+4)
	for _, kv := range os.Environ() {
		lower := strings.ToLower(kv)
		if strings.HasPrefix(lower, "http_proxy=") || strings.HasPrefix(lower, "https_proxy=") ||
			strings.HasPrefix(lower, "no_proxy=") || strings.HasPrefix(lower, "all_proxy=") {
			continue
		}
		env = append(env, kv)
	}
	return append(env,
		proxyChildMarker+"=1",
		"HTTP_PROXY="+proxyURL,
		"http_proxy="+proxyURL,
		"NO_PROXY=",
		"no_proxy=",
	)
}
