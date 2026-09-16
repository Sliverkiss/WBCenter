package webui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

// newTestHandler 用内存 FS 冒充 dist 内容，让回落规则的测试不依赖仓库里
// 恰好放了什么产物（占位 .gitkeep 还是完整构建产物都能跑）。
func newTestHandler() http.Handler {
	return handlerFor(fstest.MapFS{
		"index.html":          {Data: []byte("<!doctype html><title>app</title>")},
		"assets/app-a1b2.js":  {Data: []byte("console.log(1)")},
		"assets/app-a1b2.css": {Data: []byte("body{}")},
	})
}

func get(h http.Handler, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

// TestMissingAssetIs404NotIndexHTML /assets/ 下是 Vite 产物（文件名带内容哈希），
// 不存在就是真的不存在，必须 404。
//
// 回落 index.html 的后果是浏览器把 text/html 的 HTML 当成 JS/CSS 模块，叠加
// X-Content-Type-Options: nosniff 后只剩一片空白页和一条 MIME 报错：既看不出
// 「产物缺失」，也拿不到 404 去定位（这正是「提交了引用不存在 assets 的
// index.html」那个 bug 表现得如此难查的原因）。前端路由也不会以 /assets/ 开头，
// 所以这个前缀不需要 SPA 回落。
func TestMissingAssetIs404NotIndexHTML(t *testing.T) {
	h := newTestHandler()
	for _, p := range []string{
		"/assets/app-deadbeef.js",
		"/assets/index-BaSNl-6I.js",
		"/assets/app-a1b2.js.map",
		"/assets/missing-style.css",
	} {
		rec := get(h, p)
		if rec.Code != http.StatusNotFound {
			t.Errorf("GET %s → %d（Content-Type: %s），期望 404：回落 HTML 会让浏览器把 index.html 当模块执行，只剩空白页",
				p, rec.Code, rec.Header().Get("Content-Type"))
			continue
		}
		if ct := rec.Header().Get("Content-Type"); strings.Contains(ct, "text/html") {
			t.Errorf("GET %s → Content-Type %s，不能拿 HTML 冒充静态资源", p, ct)
		}
	}
}

// TestExistingAssetServed 存在的产物正常返回，并带上可长缓存的响应头。
func TestExistingAssetServed(t *testing.T) {
	h := newTestHandler()
	rec := get(h, "/assets/app-a1b2.js")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /assets/app-a1b2.js → %d，期望 200", rec.Code)
	}
	if got := rec.Body.String(); got != "console.log(1)" {
		t.Errorf("产物内容错误: %q", got)
	}
	if cc := rec.Header().Get("Cache-Control"); !strings.Contains(cc, "immutable") {
		t.Errorf("带内容哈希的产物应可长缓存，Cache-Control=%q", cc)
	}
}

// TestSPAFallbackStillWorksForRoutes 修掉 /assets/ 之后，前端路由的回落必须照旧。
func TestSPAFallbackStillWorksForRoutes(t *testing.T) {
	h := newTestHandler()
	for _, p := range []string{"/", "/dashboard", "/accounts/abc/settings", "/deep/link/no/where"} {
		rec := get(h, p)
		if rec.Code != http.StatusOK {
			t.Errorf("GET %s → %d，期望 200（回落 index.html）", p, rec.Code)
			continue
		}
		if !strings.Contains(rec.Body.String(), "app") {
			t.Errorf("GET %s 未返回 index.html 内容: %q", p, rec.Body.String())
		}
	}
}

// TestNotBuiltWhenIndexMissing 没有 index.html（未构建前端）时给可读提示页。
func TestNotBuiltWhenIndexMissing(t *testing.T) {
	h := handlerFor(fstest.MapFS{"placeholder.txt": {Data: []byte("x")}})
	rec := get(h, "/")
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("GET / → %d，期望 503 提示页", rec.Code)
	}
}
