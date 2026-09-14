# 里程碑 M2：账号池监控（PR #60 核心重达）

> 项目：/root/WBCenter。先读 `docs/DEV-HANDBOOK.md`（第 1 节边界、第 4 节纪律），再读本任务书。
> M0/M1 已完成：四绿工具链 + 契约/类型/MSW 地基（65 测试）。本里程碑是第一个**前后端都动**的里程碑：把 PR #60 的账号池监控能力用独立面板形态重达。

## 背景

PR #60（内嵌面板）的账号池监控语义（重达目标，不抄代码）：
- **账号卡片积分富版**：套餐明细、7 天到期预警
- **签到状态**、**猫猫旅行状态机**（位置/出发到达时间）
- **12153 徽标**（SessionDead 会话死，上游 code 12153 判定，见 `internal/upstream/client.go` 的 `IsSessionDead`）
- **全量探测**：批量、逐账号结果
- **单账号操作**：刷新凭据/签到/旅行巡检

WBCenter 后端已有能力（对照着用，别重复造）：`internal/upstream/client.go` 34 个方法（`UserResource` 积分、`DailyFreePackages` 套餐切片、`DailyCheckin`、`TravelStatus`/`BuddyInfo`、`IsSessionDead`/`IsAlreadyCheckedIn` 等错误分类）；`internal/control/service.go` 的 `Accounts()`/`Credits()`/`Run(action,uid)`。现有 `/api/accounts/{uid}/actions/{action}` 只支持 checkin/travel/refresh。

上游接口契约参考：`docs/reference/apifox-workbuddy-endpoints.md`（growth 域有出行状态/猫猫信息/成长任务的完整契约；计费域有签到与积分）。实测结论见 `/root/loop-workbuddy-reverse.md`（只读，不在仓库内，需要时用绝对路径读）。

## 任务（TDD，test 先于 feat，Go 与前端都要）

### A. 后端：`POST /api/probe` 全量探测（契约见 docs/api-contract.md §2.1）
1. `internal/control/` 新增 `Probe` 服务方法：对每个账号并发（受控并发，建议 4-8，带 context 超时）拉取：`UserResource`（积分）、签到状态、`TravelStatus`+`BuddyInfo`（猫猫状态机）。逐账号结果聚合，失败不中断（错误进 `error` 字段）。
2. 签到状态：上游没有独立"查签到"只读端点（实测结论：checkin-status 是空壳，权威判定是 checkin-activity-status——见 loop 文档 Phase 3-B）。**做法**：探测只带"上次签到"面板侧状态（来自自动化任务运行记录），不瞎编上游字段；契约里 checkin 字段标注数据来源。如果你在 Apifox 导出里找到可用的只读端点，可以改用，但要在报告里给出端点证据。
3. `TravelState`/`Buddy` 序列化对齐契约（位置/出发/到达/server_now 已有解析）。
4. SessionDead 检测：探测时上游返回 12153 → `session_dead: true` 徽标数据。
5. handler `POST /api/probe`：登录态校验、read_only 拒写、context 超时 30s（复用现有 middleware 模式）。
6. Go 测试：httptest mock 上游（`internal/upstream/client_test.go` 已有模式），覆盖：正常聚合、单账号失败不中断、12153 标记、超时降级。

### B. 前端：账号池页面重做（`pages/Accounts.tsx`）
1. 从表格升级为**卡片网格**（保持 WBCenter 观感：`.card` 语汇 + tokens）：每账号一张卡——昵称/uid/区域徽标、积分（大数字）+ 套餐明细折叠、7 天内到期 Token 红色预警、签到状态、猫猫旅行状态（idle/traveling/arrived + 位置与到达时间）、12153 徽标。
2. **全量探测按钮**：调 `POST /api/probe`，进度反馈（结果逐账号展示：成功/失败/error 文案），完成后卡片数据刷新。MSW handler 已有，按契约扩展 fixtures（健康/过期/12153/上游失败四类账号样本）。
3. 单账号操作保留现有三个（刷新凭据/签到/旅行巡检），操作后局部刷新该卡。
4. 组件测试：卡片各状态组合（含空数据/加载/错误）、探测按钮流转（pending→结果）、操作反馈。类型全用 M1 的 `Account`/`ProbeResult`（`Any` 逐步清退，本页面清零）。

### C. 概览页联动
`/api/overview` 卡片区加"会话死账号"计数（来自 probe 结果或 accounts 过滤）+ 快捷入口跳探测。后端 overview handler 相应扩展（字段加 `session_dead`，契约同步更新——契约变更顺序：改文档→改类型→改 mock→改组件）。

## 验收标准

1. `bash scripts/check.sh` 四绿（Go 新测试全绿 + race：`go test -race ./internal/control/ ./internal/upstream/` 也要绿）。
2. `POST /api/probe` 契约实现与 `docs/api-contract.md` §2.1 逐字段一致（Hermes 抽查）；真实上游联调标记为「后端联调验证项」（Hermes 执行）。
3. 账号卡片覆盖全部状态组合的测试（≥10 个新前端测试）；Accounts 页面 `Any` 清零。
4. 探测并发受控（代码里有并发上限常量/信号量，测试覆盖并发行为或至少单测验证并发度参数）。
5. 原子提交：后端 probe（test→feat）、前端卡片（test→feat）、概览联动、契约更新各自独立。

## 纪律
同手册第 4 节。特别强调：上游接口字段不确定时**先查 Apifox 导出文档**，仍不确定就在契约里标注占位并如实写进报告偏差——禁止编造上游字段。只动 /root/WBCenter；禁 push/gh/重型 MCP；新依赖禁止（echarts 是 M3 的事，本里程碑不引）。

## 交付报告
`.claude/reports/m2-report.md` + stdout 全文：commits / 文件清单 / probe 契约对齐表 / 新测试清单 / 四绿+race 证据 / 契约变更记录 / 偏差说明。
