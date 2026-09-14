# M5 OAuth 完整流 + 日志环形缓冲 + 安全回归 —— 验收记录（Hermes）

日期：2026-09-15　基线：`e49dd81`(M4 验收) → `10b771e`（7 个原子 commit）

## 验收结论：✅ 通过

## 逐项核对（对照任务书 7 条验收标准）

| # | 验收标准 | 结果 | 证据 |
|---|---|---|---|
| 1 | 四绿 + race | ✅ | Hermes 复跑：go vet/test（4 包 ok）+ web lint(0)/typecheck(0)/test(**129/129**, 14 文件)/build ✅；`go test -race` control 2.97s + upstream 1.02s 全绿 |
| 2 | OAuth 状态机 waiting→success/timeout/区域路由 + MSW 轮询序列 | ✅ | 9 个 Go 测试（11217 归一/5min 超时/cn-global 路由/凭据写入/终态幂等/只读守卫）；前端 MSW waiting→waiting→success 序列 calls===3 断言 |
| 3 | GET /api/logs 与契约一致 + 环形覆盖测试 | ✅ | logbuf.go 容量 500 + redactLog 脱敏 + mutex 并发安全；8 个 Go 测试（环形覆盖/limit/并发-race/脱敏/时区/401） |
| 4 | 日志页列表+过滤+轮询清理 | ✅ | Logs.tsx 级别过滤 + 10s 轮询 + 卸载清理 interval；7 个前端测试 |
| 5 | 安全回归 4 类 | ✅ | 8 个安全测试（超任务书）：限速429+按源隔离 / 跨源403 / 未登录401 / logout失效 / Cookie HttpOnly+SameSite=Strict+Path / 用户枚举防护 / 安全响应头(CSP/nosniff/frame-deny/no-referrer) |
| 6 | 前端新测试 ≥15 + Any 清零 | ✅ | 前端新增 17 个（OAuth 9 + Logs 7 + 1 路由），总计 129；`grep Any OAuth.tsx Logs.tsx` = 0 |
| 7 | 原子提交 test→feat | ✅ | 4 组：OAuth(41b5f8f/48220ec) → 日志(a0eeeb0/6f59258) → 安全(a2e58df) → 前端(3e7e840/10b771e) |

## Hermes 抽查

- OAuth 状态机终态幂等：success/error/timeout 后复读不再调上游（pollCalls 不增）——面板侧加固
- 日志脱敏：`redactLog` 在落缓冲前替换 token/cookie/password/secret/authorization 键值为 ***（`TestLogsRedactSecrets` 覆盖）
- 11217 归一：`client.PollLogin` 对 4xx 业务码 11217 返回 (nil,nil) → 状态机标 waiting 继续轮询（与 harness `buddy-oauth.ts:482` 一致）
- 事件埋点：登录成功/失败、OAuth 全状态、探测聚合、批量任务、单账号操作均写入 logbuf
- 安全超出任务书：用户枚举防护（用户名错/密码错同文案）+ 安全响应头（CSP/frame-deny/nosniff/no-referrer）——都是真实加固不是凑数

## 后端联调验证项（Hermes 记录，部署阶段执行）

1. OAuth 真实流 cn/global 三步（state/poll/account）
2. GET /api/logs 真实运行事件
3. 安全头在真实 HTTP 层验证

## 后续
- M6 任务书：Playwright E2E + README 重写 + 企业化收尾 v1.0.0
