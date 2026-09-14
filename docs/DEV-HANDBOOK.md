# WBCenter 开发手册

> 本文档是 WorkBuddy Control Center（WBCenter）的企业化改造总纲。执行体（Claude Code, `--model dashscope/kimi-k3`）看不到对话历史，**每个里程碑的任务书必须自包含**；本手册是任务书共用的背景层。
> 维护者：Hermes（规划+验收+部署）。执行体：Claude Code（TDD 实现者）。
> 纪律基线：本地原子提交、禁 push/fetch、四绿验收、Hermes 每里程碑验收后才进入下一里程碑。

## 1. 项目定位（不可妥协的边界）

WBCenter 是 workbuddy2api 网关的**独立 Web 管理面板**，遵循 workbuddy2api Issue #61：

- 与网关的唯一共享面：`auths/` 凭据目录（只读凭据+面板写入新账号）。
- 展示与操作数据**直连上游 WorkBuddy 官方接口**（`copilot.tencent.com` / 国际版域），不经过网关、不读网关 config.json、不调 Docker Socket、不调网关 HTTP 接口。
- 目标：把 PR #60（内嵌面板）的全部能力以独立形态重达，并纳入 Issue #70 归并的「模型价格与限免」诉求。
- 代码语言全中文注释/文档；UI 中文。

## 2. 现状盘点（2026-09-14）

### 2.1 已有（demo 坯子，commit 42cc365）
- Go 后端 `internal/control`（约 1600 行）+ `internal/upstream`（client.go 1122 行，34 个上游方法：OAuth、RefreshToken、DailyCheckin、DailyFreePackages、GrowthTasks、ConsoleTasks、AvailableModels、UserResource、Travel 状态机、BuddyInfo/First/Agreement 等），`internal/authstore`（auths 目录原子读写）。
- 22 个 `/api/*` 路由（登录/会话/概览/账号/积分/模型/活动/自动化/云端 scheduler 只读+本地 Mock CRUD/OAuth 轮询/账号动作）。
- 前端 `web/`：React 18 + Vite，**单文件 App.tsx（95 行）+ 单行 styles.css**——8 个页面全挤在一个文件里，无路由、无组件拆分、无测试、无设计 token 体系、无暗色主题。这就是"风格要保持"的 WBCenter CSS 基线（薄荷绿 `#54c7ad` + 纸白 + 大标题衬线感、圆角 18px、侧栏 276px 布局骨架、`.black` 胶囊按钮、`.badge`/`.stat`/`.card` 语汇）。
- 后端测试：`server_integration_test.go` 176 行 + `state_test.go` 140 行 + `client_test.go` 238 行（基础薄但存在）。
- 部署：docker-compose（`../workbuddy2api/auths:/auths`、`./data:/data`，UID 10001），已验证与网关 auths 共享可用。

### 2.2 上游接口知识库（面板数据源权威）
- **Apifox 253 端点导出**：`docs/reference/apifox-workbuddy-endpoints.md`（约 360KB，10 个域：growth 成长/猫猫、计费/签到、消息中心、云盘 tdrive、云助理 cloudagent、技能与插件市场、开放授权、会话 conversations、智能邮箱、专家与运营市场）。含请求/响应 schema 与示例。
- **实测结论链**：`/root/loop-workbuddy-reverse.md`（212 行，Phase 0-4 执行记录：领猫链路 agreement→buddy/first→depart→claim、签到 100/天、streak/redeem 7/14/28d、lottery EV19/抽、travel 5~10 分/天、12153 会话死、上游 scheduler access_denied 现状、checkin 双端点回退等）。
- **国内版客户端逆向产物**：`/root/wb_analysis/`（asar 解包 v5.5.4：`cli/`、`main/`、`renderer/`、`native/`、`node_modules/`、`tencent-docs/`；系统提示词 4 份；`WorkBuddy-setup.exe`）。执行体需要新端点契约细节时，先查 Apifox 导出，不够再定位到 asar 源码文件。
- **PR #60 功能语义**（重达目标，来源 PR 描述+commits）：账号池卡片（积分富版=套餐明细/7天到期预警、签到状态、猫猫旅行状态机、12153 徽标、SessionDeadFails）、一键任务+全量探测（批量、逐账号结果）、Token 与积分统计聚合、请求级日志/环形缓冲、OAuth 设备登录（state/poll/account 三步）、上游原语直调（禁用跳过/12153 计数/活跃 N 条+streak 回读）。参考实现不可抄代码（未合入 master，且其架构是内嵌 go:embed，我们是独立面板），只对齐**功能语义与字段口径**。

### 2.3 风格基线（保持 WBCenter 现有观感）
- 设计 token：薄荷绿 `#54c7ad`（primary）、墨色 `#101010`、纸白 `#fff`、muted `#707070`、line `#e9e9e6`、soft `#effaf7`、radius 18px；大标题 letter-spacing 负值；胶囊按钮 `.black`；侧栏 276px wordmark 布局。
- 升级方向：整理成 CSS 变量设计 token 文件（浅/暗双主题 `[data-theme='dark']`），但**色相与气质不变**（暗色用深灰蓝底+薄荷绿点缀，不引入新色相）。

## 3. 里程碑路线图（TDD + 原子提交 + Hermes 验收）

每个里程碑一份任务书 `docs/milestones/mN-*.md`，由 Hermes 在派发前写好（自包含：背景/约束/验收标准/提交纪律）。执行体完成的定义 = 任务书验收标准全绿 + 本地 commit 完成。

| # | 里程碑 | 内容 | 关键验收点 |
|---|---|---|---|
| M0 | 工程基线 | 前端工程化改造：tsconfig 严格、ESLint、Vitest+RTL 测试环境、目录结构（components/pages/services/hooks/types/styles）、设计 token 文件（浅/暗双主题）、布局骨架组件化（保留现有观感）、路由（react-router 或轻量 hash 路由，执行体定）、CI 脚本（`npm run lint/test/typecheck/build` 四绿） | 四绿一次跑通；观感与基线一致（截图/样式抽查）；token 色值核对 |
| M1 | 契约终稿 + API 层 | `docs/api-contract.md`：把现有 22 个 `/api/*` 端点 + 新增端点（探测/一键任务/统计/OAuth 日志等）写成请求/响应契约；`src/types/` 与契约逐字段对齐（字段名照抄后端 JSON tag，不转 camelCase）；`src/services/` fetch 封装（401 统一跳登录、错误信封解析）；MSW handlers + fixtures（dev+test 双环境） | 契约-类型一致性抽查；MSW 双环境跑通；401 流转测试 |
| M2 | 账号池监控 | PR #60 核心重达：账号卡片（积分富版：套餐明细/7 天到期预警；签到状态；猫猫旅行状态机位置/出发到达时间；12153 徽标）、全量探测（批量、逐账号结果表）、单账号操作（刷新凭据/签到/旅行巡检）。后端按需补探测聚合端点（对照 Apifox growth/计费域契约），上游调用语义对照 `internal/upstream/client.go` 既有方法 | 组件状态组合测试；探测结果逐账号呈现；RED→GREEN 证据链 |
| M3 | 一键任务 + 统计 | 一键任务面板（签到/保活/旅行/活跃/积分，逐账号结果，语义对照 PR #60「禁用跳过/12153 计数」）；Token 与积分统计页（趋势图 echarts 按需引入、汇总卡）。**注意自动化所有权**：默认只开新活动探测，接管网关任务需部署层确认（README 既有约定） | 数据转换单测；任务结果矩阵测试；轮询清理测试 |
| M4 | 模型与活动 | 模型管理升级：按账号可用模型 + **模型价格与限免标识**（Issue #70 归并诉求；价格数据源在上游模型/计费域接口，契约以 Apifox 导出为准，找不到实测字段就先做占位并明确标注「上游未提供」）；活动管理升级（成长任务/新任务探测/lottery 概览只读） | 契约字段核对；无价格数据时的降级展示测试 |
| M5 | OAuth + 日志 + 安全 | OAuth 设备登录完整流（state/poll/account、区域路由 cn/global、11217 waiting 归一）；面板运行日志查看（环形缓冲，后端新增）；安全硬化复查（登录限速/同源校验/CSRF/会话过期/HttpOnly 既有项回归） | OAuth 全流程测试（mock 上游）；安全项回归清单 |
| M6 | 打磨 + E2E + 企业化收尾 | Playwright E2E（登录→概览→账号→探测→一键任务主链路）；README 重写（安装/配置/部署/安全说明/功能矩阵）；Dockerfile/compose 复核；版本 tag v1.0.0；全量四绿 | E2E 全绿；dist 构建可嵌入 go:embed；README 与实现一致 |

## 4. 执行纪律（写入每份任务书）

1. **模型指派**：`claude -p "<任务书路径>" --model dashscope/kimi-k3 --output-format json`，`--model` 用 CLI flag（settings.json env 块会静默覆盖进程环境变量）；不设 `--max-turns`（用户明确）；workdir `/root/WBCenter`。
2. **禁止**：push/fetch/gh、拉起重型 MCP 或 node 长驻子进程（纯静态分析）、任务书外依赖引入（新依赖须任务书明示）、改 workbuddy2api 仓库任何文件。
3. **提交纪律**：每逻辑单元单 commit，`type(scope): 中文描述`；test 提交先于对应 feat（RED→GREEN 证据链）；commit 前四绿（`go test ./...` + `cd web && npm run lint && npm run typecheck && npm test && npm run build`，按 M0 定稿的脚本名）。
4. **RED 门槛**：改业务代码前必须先看到 RED（编译失败或测试失败且失败原因正确）。
5. **证据边界**：不虚假联调——所有"真实上游联调"标记为「后端联调验证项」，由 Hermes 在验收/部署阶段用真实 auths 执行；前端验收只看 mock/MSW 层证据。
6. **验收报告**：交付时 stdout 打印固定格式（commits 清单/文件清单/RED 证据/GREEN 证据/四绿证据/覆盖率/偏差说明），报告文件写到 `.claude/reports/`（不入库）。
7. **OOM 防护**：任务书一次 ≤2 个强相关子任务；长任务断点续传（commit 即存档）。

## 5. Hermes 验收协议（每里程碑）

1. 收到执行体完成通知 → `git log --oneline` 核对原子提交序列与任务书验收标准的对应关系。
2. 复跑四绿（不接受自述）。
3. 按 `executor-acceptance` 技能抓 mock 盲区：抽 1-2 个关键路径核对实现与契约语义一致。
4. 前端观感抽查（dev server 起容器冒烟或构建后 embed 跑）。
5. 写验收记录 `docs/acceptance/ACCEPT-mN.md`（入库），通过后才派发 M+1；不通过则最小修复任务书重派。
6. 全部里程碑完成后：构建镜像 → 部署（保留旧镜像）→ 真实上游冒烟 → push（脱敏检查后，用户确认）。

## 6. 关键资料索引

| 资料 | 路径 |
|---|---|
| Apifox 全量端点导出（253/10 域） | `docs/reference/apifox-workbuddy-endpoints.md` |
| 逆向实测结论链 | `/root/loop-workbuddy-reverse.md` |
| 国内版客户端 asar 解包 | `/root/wb_analysis/`（app/ asar/ tools/ unpack/） |
| PR #60（功能语义参考） | https://github.com/Sliverkiss/workbuddy2api/pull/60 |
| Issue #61（边界定义） | https://github.com/Sliverkiss/workbuddy2api/issues/61 |
| Issue #70 归并（模型价格） | 见 Issue #61 维护者归并记录 |
| workbuddy2api 仓库（只读参考，禁改） | `/root/workbuddy2api` |
| 网关 auths（真实凭据，勿入库勿打印） | `/root/workbuddy2api/auths` |
| WBCenter 部署配置 | `docker-compose.yml` + `.env`（已忽略入库） |
