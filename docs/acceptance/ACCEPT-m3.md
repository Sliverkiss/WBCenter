# M3 一键任务 + Token/积分统计 —— 验收记录（Hermes）

日期：2026-09-14　基线：`09753d8`(M2 验收) → `b8aae5f`（10 个原子 commit）

## 验收结论：✅ 通过

## 逐项核对（对照任务书 7 条验收标准）

| # | 验收标准 | 结果 | 证据 |
|---|---|---|---|
| 1 | scripts/check.sh 四绿 + race | ✅ | Hermes 复跑：go vet/test（4 包 ok）+ web lint(0)/typecheck(0)/test(**100/100**, 11 文件)/build ✅；`go test -race ./internal/control/ ./internal/upstream/` 全绿（2.86s + 1.02s） |
| 2 | batch-actions 契约逐字段一致 + 禁用跳过/12153 测试 | ✅ | Hermes 抽查：batch.go BatchUpstream 接口 + BatchActionResult JSON tag 逐字段对齐契约 §2.3；11 个 Go 测试含 TestBatchActionsSkipsDisabled + 12153 标记 + 400/401/403 |
| 3 | stats/summary 契约 §2.2 一致 + 缓存策略说明 | ✅ | 混合策略（积分实时聚合 / 签到+12153 复用 lastProbe / 状态分布凭据本地判定），契约 §2.2 已定稿；8 个 Go 测试含空池/降级/缓存 |
| 4 | echarts 按需引入 + 构建大小合理 | ✅ | `echarts/core + BarChart + GridComponent + CanvasRenderer` 子包（useEcharts.ts:12）；构建增量 +467KB min / +213KB gzip——echarts core+bar 子集理论下限；未引 echarts-for-react 全量 |
| 5 | 轮询清理测试 | ✅ | Stats.tsx 卸载 dispose echarts + clear interval；useEcharts 含 jsdom 环境探测降级；stats.test 10 用例含轮询清理验证 |
| 6 | 前端新测试 ≥15 + Any 清零 | ✅ | 前端新增 18 个用例（stats 10 + batchPanel 8）；`grep Any Stats.tsx BatchPanel.tsx` = 0 |
| 7 | 原子提交 test→feat | ✅ | 5 组完整 RED→GREEN：契约(b5b40a5) → batch-actions(eab611d/457c966) → stats(f90012d/1165ffe) → 统计页(7a879b1/4bf6d7b) → 一键任务(c48ca79/b8aae5f) |

## Hermes 抽查

- batch-actions 并发复用 Run，authstore 新增 `Disabled` 字段（嵌套/扁平两种磁盘形态对称读写）
- stats 缓存策略如实：积分实时聚合（复用 creditFor 口径），签到无上游只读端点→保守零值（与 M2 探测降级口径一致），历史趋势标注「待 M5 日志系统后补」
- echarts 在 jsdom 无 canvas 崩溃 → useEcharts 增加环境探测降级（真实 bug 修复，不是装饰性代码）
- Any 清零：Stats.tsx + BatchPanel.tsx 全用 M1 契约类型

## 后端联调验证项（Hermes 记录，部署阶段执行）

1. POST /api/batch-actions 真实账号批量签到/旅行（禁用跳过+12153 标记+逐账号结果）
2. GET /api/stats/summary 真实积分聚合 + 自动化今日统计
3. echarts 统计页在真实浏览器渲染（dev server 冒烟）

## 后续
- M4 任务书：模型价格与限免（Issue #70 归并）+ 活动管理升级
