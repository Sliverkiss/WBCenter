// Package webui 把前端构建产物（web/dist）以 embed 方式打包进二进制。
//
// 为什么用 embed 而不是依赖外部目录：GUI 要作为单文件/单容器交付，
// 产物里没有 Node 运行时也不影响运行。
package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// distFS 前端构建产物。构建前 dist 下只有占位文件，保证 go build 永远可用。
//
//go:embed all:dist
var distFS embed.FS

// Handler 返回 SPA 静态资源 handler：
//   - 命中真实文件 → 直接返回（带长缓存，文件名含内容哈希时）
//   - 其余路径     → 返回 index.html（前端路由接管）
func Handler() (http.Handler, error) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, err
	}
	return handlerFor(sub), nil
}

// handlerFor 基于给定的 dist 文件系统构造 SPA handler。
//
// 独立成函数是为了能用 testing/fstest.MapFS 精确覆盖回落规则
// （哪些路径该回落 index.html、哪些该直接 404），不必依赖仓库里恰好放了什么产物。
func handlerFor(sub fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(sub))
	index, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		// 前端未构建（只有占位文件）：返回可读的提示页而不是 404 迷宫。
		return http.HandlerFunc(notBuilt)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
		if clean == "/" {
			serveIndex(w, index)
			return
		}
		// 静态资源存在则交给 FileServer。
		if f, err := sub.Open(strings.TrimPrefix(clean, "/")); err == nil {
			_ = f.Close()
			if strings.HasPrefix(clean, "/assets/") {
				// Vite 产物文件名带内容哈希，可安全长缓存。
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			fileServer.ServeHTTP(w, r)
			return
		}
		// /assets/ 下是 Vite 产物（文件名带内容哈希），不存在就是真的不存在：
		// 回落到 index.html 会让浏览器把 text/html 的 HTML 当成 JS/CSS 模块，
		// 叠加 X-Content-Type-Options: nosniff 后只剩一片空白页和一条 MIME 报错，
		// 既看不出「产物缺失」，也拿不到 404 去定位。前端路由不会以 /assets/ 开头，
		// 这个前缀不需要 SPA 回落。
		if strings.HasPrefix(clean, "/assets/") {
			http.NotFound(w, r)
			return
		}
		// 其余一律回落 index.html（前端路由 / 深链接刷新）。
		serveIndex(w, index)
	})
}

func serveIndex(w http.ResponseWriter, index []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// index.html 绝不能缓存：否则前端发版后用户拿到旧 HTML 引用不存在的 assets。
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	_, _ = w.Write(index)
}

// notBuilt 前端未构建时的提示页。
func notBuilt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusServiceUnavailable)
	_, _ = w.Write([]byte(`<!doctype html><html lang="zh-CN"><head><meta charset="utf-8">
<title>WorkBuddy GUI · 前端未构建</title>
<style>
body{font-family:system-ui,-apple-system,"Segoe UI",sans-serif;background:#0f1117;color:#e6e8ee;
display:flex;align-items:center;justify-content:center;min-height:100vh;margin:0;padding:24px}
.card{max-width:640px;background:#171a23;border:1px solid #262b38;border-radius:12px;padding:28px}
h1{font-size:19px;margin:0 0 12px}code{background:#0b0d13;padding:2px 6px;border-radius:5px;color:#7ee787}
pre{background:#0b0d13;padding:12px;border-radius:8px;overflow:auto;color:#7ee787;font-size:13px}
p{line-height:1.7;color:#a9b0c0;font-size:14px}
</style></head><body><div class="card">
<h1>前端产物尚未构建</h1>
<p>后端已正常运行，但二进制内没有打包前端资源。请在项目根目录执行：</p>
<pre>cd web &amp;&amp; npm install &amp;&amp; npm run build
cd .. &amp;&amp; go build -o wbgui ./cmd/server</pre>
<p>然后重新启动本服务。API 接口此时已可用，例如 <code>/api/session</code>。</p>
</div></body></html>`))
}
