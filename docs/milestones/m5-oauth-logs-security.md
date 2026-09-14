# 里程碑 M5：OAuth 完整流 + 运行日志环形缓冲 + 安全硬化回归

> 项目：/root/WBCenter。先读 `docs/DEV-HANDBOOK.md`（第 1 节边界、第 4 节纪律），再读本任务书。
> M0-M4 已完成：四绿工具链 + 契约/类型/MSW + 账号池监控 + 一键任务/统计 + 模型价格/活动（112 测试）。

## 背景

### PR #60 OAuth 语义
OAuth 设备登录三步：state（发起授权）→ poll（轮询结果）→ account（账号入池）。区域路由 cn/global。11217（token not ready）归一为 waiting 继续轮询。凭据原子写入 `auths/`。

### 现有实现
- 后端 `internal/upstream/client.go` 已有 `StartLogin`/`PollLogin`/`RefreshToken`。
- 后端 `internal/control/server.go` 已有 `POST /api/oauth/start` + `POST /api/oauth/{id}/poll`。
- 前端 `web/src/pages/OAuth.tsx` 已有基础 UI（区域选择+发起+轮询），但轮询逻辑简陋（3s 定时器，无超时/重试/状态展示）。
- 安全：登录限速/同源校验/HttpOnly+SameSite=Strict 会话已有（README 记录），但缺回归测试。
- 日志：无运行日志查看能力。PR #60 有请求级日志/环形缓冲。

### 契约参考
- `docs/api-contract.md` 已有 OAuth 端点契约 + `GET /api/logs` 契约（M1 预定义，后端未实现）。

## 任务（TDD，test 先于 feat，前后端都要）

### A. 后端：OAuth 状态机完善
1. OAuth 流程状态管理：`StartLogin` 返回 state+url 后，面板维护一个内存中的登录会话（id → {state, region, created_at, status}），poll 查这个会话。11217 → status="waiting" 继续轮询；成功 → status="success" + 写入 auths；失败 → status="error" + 错误文案。
2. 轮询超时：5 分钟上限（与 harness `buddy-oauth.ts:482` 一致），超时 → status="timeout"。
3. 区域路由：cn → copilot.tencent.com，global → www.workbuddy.ai（复用 `client.go` 的 `regionBases`）。
4. Go 测试：mock 上游 OAuth 端点，覆盖 waiting→success / waiting→error / timeout / cn vs global 路由 / 凭据写入 auths。

### B. 后端：运行日志环形缓冲（契约 `GET /api/logs`）
1. `internal/control/` 新增日志缓冲：内存环形 buffer（容量 500 行），记录面板自身运行事件（自动化任务执行/探测结果/OAuth 登录/错误等）。
2. `GET /api/logs` handler：返回最近 N 条日志（支持 `?limit=` 查询参数，默认 100）。
3. 日志格式：`{time, level, action, message}`（time 用 cfg.Timezone）。
4. Go 测试：写入日志 → 读取验证、环形覆盖（写超过容量后最旧的被覆盖）、limit 参数。

### C. 前端：OAuth 页升级
1. `web/src/pages/OAuth.tsx` 升级：完整三步流程展示——授权链接（新窗口打开）、轮询状态（waiting/success/error/timeout + 倒计时）、成功后账号信息、失败重试按钮。
2. 5 分钟超时倒计时显示。
3. 区域选择带说明文案（中国大陆/国际版）。
4. MSW handler 扩展：oauth start → 返回 fake state+url；oauth poll → 依次返回 waiting→waiting→success（模拟真实轮询序列）。
5. 组件测试：发起→轮询中→成功、发起→轮询中→超时、区域路由、重试。

### D. 前端：日志查看页
1. 新增日志页或嵌入概览页/设置页：调 `GET /api/logs?limit=100`，展示最近日志列表（时间/级别/动作/消息），支持级别过滤（all/info/error）。
2. 自动刷新（10s 轮询）+ 卸载清理 interval。
3. MSW handler：返回 5 条样本日志（含 info/error 级别）。
4. 组件测试：日志列表渲染、级别过滤、轮询清理。

### E. 安全硬化回归测试
1. 登录限速：同一 IP 短时间多次失败登录 → 429（后端已有逻辑，补测试）。
2. 同源校验：写操作（POST/PUT/DELETE）从不同 Origin 来 → 403。
3. 会话：未登录访问写操作 → 401；logout 后旧 session 不可用。
4. HttpOnly + SameSite=Strict：response Set-Cookie 头验证。
5. Go 测试：覆盖以上 4 类安全场景。

## 验收标准

1. `bash scripts/check.sh` 四绿 + `go test -race ./internal/control/ ./internal/upstream/` 全绿。
2. OAuth 状态机：waiting→success / timeout / 区域路由 有测试覆盖；MSW 模拟轮询序列。
3. `GET /api/logs` 与契约一致；环形覆盖测试通过。
4. 日志页：列表+过滤+轮询清理测试通过。
5. 安全回归：4 类安全场景测试通过（限速/同源/会话/Cookie 属性）。
6. 前端新测试 ≥15 个；Any 在改动页面清零。
7. 原子提交：OAuth 后端（test→feat）、日志后端（test→feat）、安全回归（test）、OAuth 前端（test→feat）、日志前端（test→feat）各自独立。

## 纪律
同手册第 4 节。只动 /root/WBCenter；禁 push/gh/重型 MCP；不引新依赖。

## 交付报告
`.claude/reports/m5-report.md` + stdout 全文：commits / 文件清单 / OAuth 状态机对齐表 / 日志环形缓冲验证 / 安全回归清单 / RED→GREEN 证据 / 四绿+race 证据 / 覆盖率 / 偏差说明。
