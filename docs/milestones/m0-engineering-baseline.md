# 里程碑 M0：工程基线（前端工程化改造）

> 项目：/root/WBCenter（WorkBuddy Control Center——workbuddy2api 的独立 Web 管理面板）
> 你看不到规划对话。总纲与边界先读 `docs/DEV-HANDBOOK.md`（必读，特别是第 1 节边界与第 4 节纪律），再读本任务书。

## 背景

WBCenter 前端 `web/` 目前是 demo 坯子：React 18 + Vite，但只有单文件 `src/App.tsx`（95 行，8 个页面函数全挤在里面）+ 单行 `src/styles.css`。没有路由、没有组件拆分、没有测试、没有 lint/typecheck 脚本。后端（Go，`internal/`）本里程碑**不改**。

现有观感是产品的风格基线，必须保持：薄荷绿 `#54c7ad`、墨 `#101010`、纸白 `#fff`、muted `#707070`、line `#e9e9e6`、soft `#effaf7`、radius 18px、大标题 letter-spacing 负值、胶囊黑按钮 `.black`、276px 侧栏 wordmark 布局。改造前先通读 `web/src/App.tsx` 和 `web/src/styles.css`，把现有页面（overview/accounts/credits/models/automation/scheduler/activities/oauth + 登录页）全部行为原样保留。

## 任务（TDD：每步先写失败测试再实现，test commit 先于 feat commit）

### A. 四绿脚本与工具链
1. `web/package.json` 增加脚本：`lint`（ESLint 9 flat config + typescript-eslint + react 插件）、`typecheck`（tsc --noEmit 或 tsc -b）、`test`（Vitest + @testing-library/react + jsdom，建议 `test` 单跑 + `test:run` CI 模式）、`build`（保留现有）。
2. `npm run lint && npm run typecheck && npm test -- --run && npm run build` 四绿一次跑通（允许 `npm install` 安装 devDependencies，但**只允许**本任务书列出的类别：ESLint 体系 / Vitest / Testing Library / jsdom / @vitest/coverage-v8。不引 UI 框架、不引 tailwind、不引状态管理库、不引 axios——fetch 已够用）。
3. 给 Go 侧留一个根级 `make` 或脚本 `scripts/check.sh`：`go vet ./... && go test ./... && cd web && npm run lint && npm run typecheck && npm test -- --run && npm run build`（Go 侧现有测试必须保持全绿，不许改 Go 代码，如果 Go 测试本来就红，停下来在报告里说明，不要修）。

### B. 设计 token 化 + 暗色主题
1. 新建 `web/src/styles/tokens.css`：把 styles.css 里散落的色值/圆角/间距收敛为 CSS 变量（`:root` 浅色 + `[data-theme='dark']` 暗色）。暗色用深灰蓝底（如 `#0d1117` 系）+ 现有薄荷绿点缀，**不引入新色相**；文字/边框/表面色按暗色惯例推导。
2. `App.tsx` 加主题切换（侧栏底部小按钮，localStorage 持久化 `data-theme`，默认浅色），document.documentElement 上挂 `data-theme`。
3. 测试：渲染时默认浅色、切换后 html 的 data-theme 变化、localStorage 持久化（jsdom 下 mock localStorage）。

### C. 组件拆分与目录结构
按此结构重组（保持观感与行为不变）：
```
src/
  components/    # Button/Badge/Card/StatCard/EmptyState/Loading/ErrorView/Drawer/Notice/Head/PageShell(侧栏+主区) 等自绘原语
  pages/        # Overview/Accounts/Credits/Models/Automation/Scheduler/Activities/OAuth/Login 八页
  services/     # api.ts（现有 request/get/post 封装迁入，暂不做 MSW——那是 M1）
  types/        # 先放 Any 之类过渡类型 + Page 类型
  styles/       # tokens.css + 全局样式
  hooks/        # useNotice 等可复用 hook
  App.tsx       # 只剩布局 + 路由装配
```
1. 拆分过程保持 fetch 调用 URL 与行为逐字不变（对照 `internal/control/server.go` 的 22 条路由）。
2. 每个页面组件至少一个冒烟测试（渲染 + 关键数据出现在 DOM；fetch 用 vi.stubGlobal 或 mock，不引 MSW）。
3. 组件原语（Button/Badge/Card 等）各至少 2 个用例（含边界：disabled/bad 状态）。

### D. 路由
用 hash 路由（自己写 30 行内的小 hook 即可，不引 react-router；如果发现引 react-router 更省总代码量，允许引入但要在报告里说明理由）。8 个页面各一个路由，登录态守卫（未登录一律渲染 Login）。测试：hash 切换渲染对应页面、未登录守卫。

## 验收标准（全部满足才算完成）

1. `bash scripts/check.sh` 四绿（Go vet/test + 前端 lint/typecheck/test/build）。
2. `git log --oneline` 显示原子提交序列：test 提交先于对应 feat（RED→GREEN 证据链），提交信息 `type(scope): 中文描述`。
3. 现有 8 个页面行为与观感零回归（报告里附每页拆分前后的职责对照表；构建产物 `web` dist 可生成）。
4. 前端测试数量 ≥ 25 个，全部通过。
5. 暗色主题可用且无新色相；token 文件里浅色值与现 styles.css 逐字一致（报告附对照）。

## 纪律（违反即任务失败）

- 只动 `/root/WBCenter`；禁改 `/root/workbuddy2api`、禁 push/fetch/gh、禁拉 MCP 或 node 长驻子进程（测试跑完进程要退出）。
- 新依赖只限任务书列出类别；ESLint 报的现有代码问题允许顺手修（这属于本次工程化范围）。
- 每完成一个子任务立即 commit（断点续传存档）；commit 前 `cd web && npm run lint && npm run typecheck && npm test -- --run` 必须绿。
- 最终在 stdout 打印交付报告：commits 清单 / 文件清单 / RED→GREEN 证据 / 四绿证据 / 测试清单 / 覆盖率（coverage-v8）/ 偏差说明（任何没做到的点如实说明，禁止粉饰）。

## 交付报告
写到 `.claude/reports/m0-report.md`（该目录不入库），stdout 打印全文。
