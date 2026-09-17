package control

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"workbuddy-control-center/internal/authstore"
	"workbuddy-control-center/internal/upstream"
)

// captureTransport 记录上游请求（路径与 body），让测试不需要真的打网络。
type captureTransport struct {
	mu    sync.Mutex
	paths []string
	bodys [][]byte
}

func (c *captureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	raw, _ := io.ReadAll(r.Body)
	c.mu.Lock()
	c.paths = append(c.paths, r.URL.Path)
	c.bodys = append(c.bodys, raw)
	c.mu.Unlock()
	// 带一个 PackageCode，让 creditFor 继续走到「今日切片」那一跳。
	body := `{"code":0,"data":{"Response":{"Data":{"TotalDosage":100,"Accounts":[` +
		`{"PackageCode":"p_tcaca","CycleCapacitySize":100,"CycleCapacityRemain":50}]}}}}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    r,
	}, nil
}

func (c *captureTransport) bodyFor(t *testing.T, path string) map[string]any {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, p := range c.paths {
		if p != path {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal(c.bodys[i], &m); err != nil {
			t.Fatalf("解析 %s 请求体失败: %v", path, err)
		}
		return m
	}
	t.Fatalf("未捕获到 %s 请求，捕获到：%v", path, c.paths)
	return nil
}

// wallClockClose 判断 got 作为 loc 时区的墙钟，是否落在 want 附近。
func wallClockClose(got string, want time.Time, loc *time.Location, tol time.Duration) bool {
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", got, loc)
	if err != nil {
		return false
	}
	d := parsed.Sub(want)
	if d < 0 {
		d = -d
	}
	return d <= tol
}

// TestDailyWindowFollowsConfiguredTimezone 面板的「今天」必须按配置时区算。
//
// 部署形态决定进程时区通常不是用户所在时区：Dockerfile 装了 tzdata 却在运行阶段
// 没有 ENV TZ，compose 也没有 TZ，于是 Alpine 里 time.Local 就是 UTC；而 config
// 的 timezone 只被 LoadConfig 校验、被 /api/session 展示，从未被任何计算使用。
//
// 后果不是「显示时差」这么轻：今日额度切片窗口（SlicePeriodStartTime/EndTime）与
// 积分查询区间（PackageEndTimeRangeBegin）都是**墙钟字符串**，由上游按该时区解释。
// 用 UTC 格式化，国内用户每天 00:00–08:00 请求到的是**前一天**的切片窗口，面板上
// 「今日发放 / 今日已用 / 今日剩余」整天有三个多小时是错的。
func TestDailyWindowFollowsConfiguredTimezone(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("加载时区失败: %v", err)
	}
	// 模拟默认容器：进程时区是 UTC（挂钟与配置时区相差 8 小时，永不相同）。
	saved := time.Local
	time.Local = time.UTC
	defer func() { time.Local = saved }()

	dir := t.TempDir()
	store, err := authstore.New(dir + "/auths")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(&authstore.Account{
		UID: "account-123456", Nickname: "测试账号",
		AccessToken: "test-access-token", RefreshToken: "test-refresh-token",
		ExpiresAt: time.Now().Add(time.Hour).Unix(), Domain: "https://example.invalid",
	}); err != nil {
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
	rt := &captureTransport{}
	up := upstream.New(cfg.Timeout())
	up.HTTP = &http.Client{Transport: rt, Timeout: cfg.Timeout()}

	// 关键：装配时就必须把配置时区交给上游客户端。
	svc := NewService(cfg, store, up, state)
	svc.Credits(context.Background(), "")

	want := time.Now().In(loc)

	// 1) 积分查询区间的墙钟必须按配置时区。
	resource := rt.bodyFor(t, "/v2/billing/meter/get-user-resource")
	gotBegin, _ := resource["PackageEndTimeRangeBegin"].(string)
	if !wallClockClose(gotBegin, want, loc, 30*time.Second) {
		t.Errorf("积分查询区间起点用了进程时区(%s)的墙钟，而不是配置时区(%s)：\n  got  %q\n  want ≈ %q",
			time.Local, loc, gotBegin, want.Format("2006-01-02 15:04:05"))
	}

	// 2) 今日切片窗口必须是配置时区的自然日边界。
	free := rt.bodyFor(t, "/billing/meter/get-user-resource-free-packages")
	gotStart, _ := free["SlicePeriodStartTime"].(string)
	gotEnd, _ := free["SlicePeriodEndTime"].(string)
	wantStart := time.Date(want.Year(), want.Month(), want.Day(), 0, 0, 0, 0, loc)
	wantEnd := wantStart.Add(24*time.Hour - time.Millisecond)
	if gotStart != wantStart.Format("2006-01-02 15:04:05") {
		t.Errorf("今日切片窗口起点不是配置时区(%s)的自然日边界：got %q，期望 %q", loc, gotStart, wantStart.Format("2006-01-02 15:04:05"))
	}
	if gotEnd != wantEnd.Format("2006-01-02 15:04:05") {
		t.Errorf("今日切片窗口终点不正确：got %q，期望 %q", gotEnd, wantEnd.Format("2006-01-02 15:04:05"))
	}
}
