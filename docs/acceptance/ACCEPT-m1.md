# M1 契约终稿 + TS 类型 + MSW —— 验收记录（Hermes）

日期：2026-09-14　基线：`4e1aaa1` → `2ebc52a`（7 个原子 commit）

## 验收结论：✅ 通过

## 逐项核对（对照任务书 6 条验收标准）

| # | 验收标准 | 结果 | 证据 |
|---|---|---|---|
| 1 | 四绿复跑 | ✅ | Hermes 复跑 check.sh：go vet/test ✅ + lint 0 / typecheck 0 / test **65/65**（8 文件）/ build ✅ |
| 2 | 契约覆盖 22 存量 + ≥4 新端点，字段与 Go 逐字一致 | ✅ | 抽查 1：`Account`（service.go:15-22）六字段 uid/nickname/domain/expires_at/expired/needs_refresh 逐字对齐 TS interface（types/index.ts:51-59）；抽查 2：`CreditSummary`（service.go:23-33）九字段含 `error?` omitempty 对齐（types/index.ts:70-80）；抽查 3：新端点 4 个全部有契约（api-contract.md §2.1-2.4 probe/stats/logs/pricing）且 services 方法齐备（api.ts:108-111）。契约文档 510 行 29 节 |
| 3 | MSW 双环境 + onUnhandledRequest:error 全绿 | ✅ | setup.ts:8 确认 `server.listen({ onUnhandledRequest: 'error' })`；main.tsx 确认 DEV+VITE_USE_MSW 门控动态 import；dist 产物 grep msw = 0（handlers 不进 bundle）；65/65 全绿证明无裸 fetch 漏网 |
| 4 | 401 → 登出流转 | ✅ | api.test.ts 覆盖：401 触发回调 + hash 重置 `#/`；含登录端点 401 不误触发全局登出的负例 |
| 5 | 类型守卫单测 ≥3 含负例 | ✅ | 17 个（6 组守卫，null/undefined/数组/类型错误负例） |
| 6 | 原子提交纪律 | ✅ | 契约独立 commit（1ff61a8）→ 守卫 test→feat（6764087/0f1aeea）→ msw 依赖独立（e30f427）→ services test→feat（73bcd14/786a0b0）→ MSW 迁移独立（2ebc52a） |

## Hermes 抽查（mock 盲区排查）

- authstore 的磁盘 Account（accessToken/refreshToken）与 control 的 API Account 是两个不同结构——执行体正确只对齐了 API 层的，未混淆。
- fixtures 2 账号数据集（健康+过期降级）比单账号 stub 更强，迁移中真实暴露并修复 3 个测试失败（多元素/断言歧义/401 handler 覆盖）——这是真测试不是走过场。
- 覆盖率 85.79% 语句（types 100%/services 94.64%），与自述一致。
- 偏差 5 条核实均合理（lint 修复归并 commit、登录测试补口令输入、fixtures 双账号、新端点仅契约、联调项移交）。

## 后端联调验证项（Hermes 记录，部署阶段执行）

1. 真实 auths 逐端点比对（重点 credits 降级文案、scheduler-tasks 错误映射）
2. 会话过期 >12h 写操作 401 流转
3. OAuth 真实流（cn/global 三步）

## 后续

- M2 任务书：`docs/milestones/m2-account-pool.md`（账号池监控——PR #60 核心重达，含后端 probe 端点实现）
