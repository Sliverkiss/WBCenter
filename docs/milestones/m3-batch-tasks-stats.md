# 里程碑 M3：一键任务 + Token/积分统计

> 项目：/root/WBCenter。先读 `docs/DEV-HANDBOOK.md`（第 1 节边界、第 4 节纪律），再读本任务书。
> M0-M2 已完成：四绿工具链 + 契约/类型/MSW 地基 + 账号池监控（probe + 卡片网格 + 12153 徽标，82 测试）。

## 背景

PR #60 的「一键任务」语义：批量执行签到/保活/旅行/活跃/积分探测，逐账号结果展示（含禁用跳过、12153 计数、活跃 N 条+streak 回读）。PR #60 还有 Token 与积分统计聚合（请求行 hook 聚合、趋势展示）。

WBCenter 后端已有：`internal/control/service.go` 的 `Run(ctx, action, uid)` 方法（支持 checkin/travel/refresh，已在 `/api/accounts/{uid}/actions/{action}` 使用）；`internal/control/state.go` 的自动化任务体系（`Automations()` + `RunRecord`）；`internal/upstream/client.go` 的 `UserResource`（积分）、`DailyCheckin`、`TravelOnce` 等。M2 新增了 `POST /api/probe`（全量探测）。

契约参考：`docs/api-contract.md` §2.2 `GET /api/stats/summary`（M1 已定义契约+mock，后端未实现）。

## 任务（TDD，test 先于 feat，前后端都要）

### A. 后端：`GET /api/stats/summary` 统计汇总（契约见 docs/api-contract.md §2.2）
1. `internal/control/` 新增统计服务方法：聚合全账号的积分总览（总余额/今日发放/今日已用/今日剩余/套餐数）+ 账号状态分布（健康/过期/会话死/禁用）+ 自动化运行统计（今日成功/失败次数）。
2. 数据来源：复用 M2 的 `lastProbe` 缓存（最近一次探测结果）或按需实时查上游（选缓存策略，在报告里说明）；自动化统计来自 `state.go` 的 RunRecord 历史。
3. handler `GET /api/stats/summary`：登录态校验。
4. Go 测试：httptest mock 上游 + state fixtures，覆盖：正常聚合、空账号池、部分账号上游失败降级、缓存命中。

### B. 后端：一键批量任务端点
1. 契约扩展：`POST /api/batch-actions`——请求 `{ action: "checkin"|"travel"|"refresh", uids?: string[] }`（uids 缺省=全量），响应 `{ results: [{uid, nickname, ok, message, skipped?, error?}] }`（与 PR #60 语义对齐：禁用跳过、12153 计数）。
2. 实现：复用 `Run(ctx, action, uid)`，逐账号执行（受控并发，与 probe 共用或独立信号量），结果聚合。禁用账号 skip（不调上游），12153 SessionDead 标记但不视为可重试失败。
3. Go 测试：批量执行、禁用跳过标记、单账号失败不中断、12153 标记、指定 uid 子集。

### C. 前端：一键任务面板（新增页面或 Automation 页升级）
1. 在 Automation 页或新增独立「一键任务」区：四个操作按钮（签到/保活/旅行/刷新凭据），点击后调 `POST /api/batch-actions`，逐账号结果表（uid/昵称/状态徽标/消息/跳过标记/失败错误）。
2. **自动化所有权防双跑**（手册第 3 节 M3 要求）：面板默认只开新活动探测，接管网关任务需部署层确认——在 UI 上体现为自动化开关旁的提示文案，不自动接管。
3. MSW handler 扩展：batch-actions 返回 4 类账号结果样本（成功/跳过/失败/12153）。
4. 组件测试：按钮 pending 禁用、结果逐行渲染、跳过标记、失败文案、空结果。

### D. 前端：统计页（新增或 Credits 页升级）
1. 引入 echarts（按需 `echarts/core` + `useEcharts` 封装 hook，**唯一新运行时依赖**，不引 echarts-for-react 全量包）。
2. 统计汇总卡（总积分/今日发放/已用/剩余 + 账号状态分布饼图/条形）+ 趋势折线图（每日积分趋势——数据来自 `GET /api/stats/summary` 的历史快照或 state 里的 probe 历史记录；如无历史数据则展示当前快照 + 「历史趋势待 M5 日志系统后补」降级标注）。
3. 轮询清理：页面卸载时 dispose echarts 实例 + clear interval（测试覆盖）。
4. 组件测试：数据转换单测（原始响应 → 图表 option）、空数据降级、轮询清理验证。

### E. 契约同步
1. `docs/api-contract.md` 补 `POST /api/batch-actions` 契约（如果 M1 没预定义，按本任务书定义补写）。
2. `src/types/` + `src/services/` + `src/mocks/` 同步新端点类型与方法。
3. 契约变更顺序：改文档 → 改类型 → 改 mock → 改组件。

## 验收标准

1. `bash scripts/check.sh` 四绿 + `go test -race ./internal/control/ ./internal/upstream/` 全绿。
2. `POST /api/batch-actions` 契约实现与文档逐字段一致（Hermes 抽查）；禁用跳过、12153 标记有测试覆盖。
3. `GET /api/stats/summary` 实现与契约 §2.2 一致；缓存策略在报告里说明。
4. echarts 引入按需（`echarts/core`），不引全量包；`grep -r 'echarts' internal/webui/dist/assets/*.js` 产物大小增长合理（报告附构建前后对比）。
5. 轮询清理测试通过（页面卸载后 echarts dispose + interval clear）。
6. 前端新测试 ≥15 个；Any 在新页面清零。
7. 原子提交：后端 batch-actions（test→feat）、后端 stats（test→feat）、echarts 引入+统计页（test→feat）、一键任务面板（test→feat）、契约同步各自独立。

## 纪律
同手册第 4 节。新依赖只允许 echarts/core 相关包（echarts/core + echarts/charts + echarts/components + echarts/renderers 按需子包）。不引状态管理库（zustand 不需要，fetch+hooks 够用）。只动 /root/WBCenter；禁 push/gh/重型 MCP。

## 交付报告
`.claude/reports/m3-report.md` + stdout 全文：commits / 文件清单 / batch-actions 契约对齐表 / stats 缓存策略说明 / echarts 构建前后大小对比 / RED→GREEN 证据 / 四绿+race 证据 / 覆盖率 / 偏差说明。
