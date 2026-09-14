# 里程碑 M1：API 契约终稿 + TS 类型 + MSW 双环境

> 项目：/root/WBCenter。先读 `docs/DEV-HANDBOOK.md`（第 1 节边界、第 4 节纪律），再读本任务书。
> M0 已完成：前端已组件化 + 四绿工具链（39 测试）。本里程碑做「契约先行」的契约层，为 M2+ 的功能重达打地基。

## 背景

当前前端页面用的是 `Any = Record<string, unknown>` 透传，没有类型约束。后端 Go 侧 22 条路由（`internal/control/server.go`）是契约的**现状权威**，本里程碑把它文档化为正式契约，并定义 M2-M5 需要的**新增端点**契约（后端实现随后续里程碑补，本里程碑只定契约+mock）。

## 任务（TDD，test 先于 feat）

### A. 契约终稿 `docs/api-contract.md`
逐端点写：方法、路径、请求体、响应 JSON 示例、字段语义表、错误信封（现有约定：非 2xx 返回 `{"error": "中文错误"}`）。分两部分：

1. **存量 22 端点**：逐条对照 `internal/control/server.go` 的 handler 实现 + `internal/control/service.go` 写准确（字段名照抄 Go JSON tag，snake_case 不转 camelCase）。特别注意：
   - `/api/session`（authenticated/read_only/timezone/using_default_password）
   - `/api/accounts`（uid/nickname/domain/expires_at/expired/needs_refresh）
   - `/api/credits` 与 `/api/credits/{uid}`（current/today_allocated/today_consumed/today_remaining/packages/fetched_at/error）
   - `/api/automations`、`/api/scheduler-tasks`、`/api/mock-scheduler-tasks` CRUD、`/api/oauth/start|poll`、`/api/accounts/{uid}/actions/{action}`
2. **新增端点契约**（M2-M5 要用的，只定义不实现）：
   - `POST /api/probe`——全量探测：请求 `{}`，响应 `{ results: [{uid, nickname, ok, checkin?: {...}, credits?: {...}, travel?: {...}, error?}] }`（对照 PR #60 语义：批量、逐账号结果）
   - `GET /api/stats/summary`——统计汇总卡数据
   - `GET /api/logs`——面板运行日志（环形缓冲快照，`{ lines: [...] }`）
   - `GET /api/models/pricing`——模型价格与限免（Issue #70；若上游契约未定，契约里字段留占位并标注「上游字段未核实，M4 实测定稿」）
   - 每个新端点给出错误分支（401/上游失败降级）

### B. TS 类型层 `src/types/`
1. 把契约逐字段转成 interface（命名对齐后端语义，如 `Account`、`CreditSummary`、`ProbeResult`、`AutomationItem`、`OAuthStartResponse`）。禁止改字段名。
2. 写类型守卫或 zod-free 的窄化函数（不引 zod，手写 `isAccount(x)` 之类，够用即可）——至少给 Account/CreditSummary/ProbeResult 写。

### C. services 层强化
1. `src/services/api.ts` 保持 fetch 封装，新增各域方法：`getAccounts()`、`getCredits()`、`startProbe()`、`getStatsSummary()`、`getLogs()`、`getModelPricing()` 等，返回值直接用新类型（页面本里程碑**不大改**，只把 Login/Overview 两个页面的 Any 换成新类型作为示范，其余页面 M2+ 逐页迁移）。
2. 401 统一处理：`request` 里 401 时触发全局登出回调（简单事件或回调注入，不引状态库），跳转登录 hash。测试覆盖：401 → 回调触发 + hash 变 `#/`。

### D. MSW 双环境
1. 引入 `msw@2`（允许，唯一新运行时 devDependency；handlers 不进 bundle）。
2. `src/mocks/handlers.ts`：存量 22 端点 + 新端点全部按契约返回确定性 fixtures（`src/mocks/fixtures.ts`）。错误分支也要有（401/上游失败）。
3. dev 环境：`main.tsx` 检测 `import.meta.env.DEV && import.meta.env.VITE_USE_MSW !== '0'` 时启用 worker（`npx msw init public/ --save` 生成 worker 文件）。
4. test 环境：Vitest setup 里 `server.listen({ onUnhandledRequest: 'error' })` + 每用例 resetHandlers。
5. **把 M0 冒烟测试里的 fetch mock 全部迁移到 MSW**（这是红线：同一套 handler 服务 dev+test，测试不再自己造 fetch 假数据）。

### E. 收尾
1. `.gitignore` 补 `web/coverage/`（M0 遗留）。
2. `scripts/check.sh` 不变（npm test 已含新测试）。

## 验收标准

1. 四绿复跑通过（`bash scripts/check.sh`）。
2. `docs/api-contract.md` 覆盖 22 存量端点 + ≥4 新端点，字段名与 Go 实现逐字一致（Hermes 会抽查 3 个端点对照 Go 源码）。
3. MSW dev+test 双环境：`npm run dev` 起来后页面全部走 mock（截图不必，报告说明验证方式即可）；测试里 `onUnhandledRequest: 'error'` 下全绿（证明没有裸 fetch 漏网）。
4. 401 → 登出流转测试通过。
5. 类型守卫单测 ≥ 3 个（含非法输入返回 false 的负例）。
6. 原子提交：契约文档单独 commit，类型+守卫 test→feat，MSW 迁移单独 commit。

## 纪律
同手册第 4 节：只动 /root/WBCenter；禁 push/gh/重型 MCP；新依赖只允许 msw@2 及其 peer；后端 Go 代码本里程碑**零改动**（契约新增端点不实现）；commit 前四绿。

## 交付报告
`.claude/reports/m1-report.md` + stdout 全文：commits / 文件清单 / 契约-类型对齐抽查表 / MSW 双环境证据 / RED→GREEN 证据 / 四绿证据 / 覆盖率 / 偏差说明。
