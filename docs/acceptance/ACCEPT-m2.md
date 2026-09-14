# M2 账号池监控 —— 验收记录（Hermes）

日期：2026-09-14　基线：`4efac75`(M1 验收) → `efeb200`（8 个原子 commit）

## 验收结论：✅ 通过

## 逐项核对（对照任务书 5 条验收标准）

| # | 验收标准 | 结果 | 证据 |
|---|---|---|---|
| 1 | scripts/check.sh 四绿 + race | ✅ | Hermes 复跑：go vet/test（control/upstream/authstore/fsutil 全 ok）+ web lint(0)/typecheck(0)/test(**82/82**, 9 文件)/build ✅；`go test -race ./internal/control/ ./internal/upstream/` 全绿（1.31s + 1.02s） |
| 2 | /api/probe 与契约 §2.1 逐字段一致 | ✅ | Hermes 抽查：probe.go `ProbeConcurrency=6` 信号量 + 30s ctx 超时；server.go:196 overview handler 含 `session_dead` 联动字段；checkin 字段标注本地来源（非编造上游）；travel.location/departed_at/arrive_at 按缺失降级（契约标注上游未提供）——真实上游联调留到部署阶段 |
| 3 | 卡片状态组合测试 ≥10 + Accounts Any 清零 | ✅ | 前端 17 个新用例（accounts.test 15 + overview 2）；`grep Any Accounts.tsx` = 0 命中 |
| 4 | 探测并发受控 | ✅ | `const ProbeConcurrency = 6` + `sem := make(chan struct{}, ProbeConcurrency)` 信号量（probe.go:17,77）；TestProbeConcurrencyIsBounded 验证 |
| 5 | 原子提交 test→feat | ✅ | eec26b2 契约 → 2f0afd6/de41240 后端 RED→GREEN → 9cbad23/c260f19 Accounts RED→GREEN → 136b10d/40e834d Overview RED→GREEN → efeb200 测试修正，4 组完整配对 |

## Hermes 抽查（mock 盲区排查）

- probe.go 并发模型：每请求 `map[int]*toolCallState` + 信号量 6，失败隔离不中断（TestProbeSingleAccountFailureDoesNotStopOthers 覆盖）
- 12153 SessionDead 标记：上游 IsSessionDead 判定 → `session_dead: true` → 前端 `<Badge ok={false}>会话失效</Badge>`（Accounts.tsx:72）
- Accounts 卡片从表格升级为网格（AccountCard 组件），保持 WBCenter 观感（global.css 新增 account-grid/account-card/account-credits 样式，薄荷绿色系不变）
- 执行体如实标注了两个降级：checkin 无上游只读端点（本地自动化记录来源）、travel 地点/时间上游未透出（契约保留字段+前端降级）——非编造，符合任务书防编造条款

## 后端联调验证项（Hermes 记录，部署阶段执行）

1. 真实 auths 起后端，POST /api/probe 全量探测真实账号
2. 12153 会话死账号在真实上游场景下的标记准确性
3. overview.session_dead 与 probe 结果联动

## 后续

- M3 任务书：一键任务 + Token/积分统计（echarts 趋势）
