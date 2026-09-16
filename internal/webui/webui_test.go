package webui

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

// assetRefRe 匹配 index.html 里对构建产物（/assets/…）的引用。
var assetRefRe = regexp.MustCompile(`(?:src|href)="(/assets/[^"]+)"`)

// TestCommittedDistIsSelfConsistent 守住「仓库里提交的 dist 必须自洽」这条不变量。
//
// 背景（真实缺陷）：仓库里提交的是**一份 Vite 构建产物** index.html，它引用
// /assets/index-<hash>.js 与 .css，但 .gitignore 忽略了 internal/webui/dist/assets/，
// 真正的产物并不在仓库里。于是全新克隆 `go build` 出来的二进制：
//
//   - Handler() 能在 embed 里找到 dist/index.html，因此**不会**走 notBuilt 提示页
//     （而 webui.go 的注释明确写着「前端未构建（只有占位文件）：返回可读的提示页
//     而不是 404 迷宫」——设计意图被这份构建产物绕过了）；
//   - 浏览器取 /assets/index-<hash>.js 时被 SPA 回退返回 text/html 的 index.html，
//     叠加 X-Content-Type-Options: nosniff，浏览器拒绝把它当模块执行 → 纯空白页，
//     控制台只剩 MIME 报错，用户完全看不出「前端没构建」。
//
// 因此断言：若 dist 下存在 index.html，它引用的每个 /assets/ 产物都必须真的在
// embed 里（即仓库确实自带完整产物）；否则就不该把 index.html 提交进仓库，
// 让 Handler() 按设计回落到 notBuilt 提示页。
func TestCommittedDistIsSelfConsistent(t *testing.T) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		t.Fatalf("fs.Sub(dist): %v", err)
	}
	raw, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		// 未提交 index.html：Handler() 会返回可读的构建提示页，符合占位设计。
		return
	}
	refs := assetRefRe.FindAllStringSubmatch(string(raw), -1)
	if len(refs) == 0 {
		return // 自包含的占位页，不引用任何产物
	}
	for _, m := range refs {
		p := strings.TrimPrefix(m[1], "/")
		if _, err := fs.Stat(sub, p); err != nil {
			t.Errorf("dist/index.html 引用了仓库中不存在的产物 %s：go build 出的服务会返回空白页而非构建提示页（%v）", m[1], err)
		}
	}
}
