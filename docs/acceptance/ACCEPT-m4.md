# M4 模型价格与限免 + 活动管理升级 —— 验收记录（Hermes）

日期：2026-09-15　基线：`1ada1bd`(M3 验收) → `f62b7d1`（5 个原子 commit）

## 验收结论：✅ 通过

## 逐项核对（对照任务书 6 条验收标准）

| # | 验收标准 | 结果 | 证据 |
|---|---|---|---|
| 1 | 四绿 + race | ✅ | Hermes 复跑：go vet/test（4 包 ok）+ web lint(0)/typecheck(0)/test(**112/112**, 12 文件)/build ✅；`go test -race` control 2.9s + upstream 1.0s 全绿 |
| 2 | /api/models/pricing 与契约一致 + 无价格不编造 | ✅ | Hermes 抽查：service.go `pricing_available` + `note` 字段；`TestModelPricingNoCredits`/`TestModelPricingUpstreamFails` 覆盖降级；M1 编造的 price/price_unit/free_quota 已删除，改为 credits 透传 |
| 3 | 模型页价格/限免/图片支持 + 降级文案 | ✅ | Models.tsx 3 处「上游未提供」降级（空 credits / pricing_available=false / 无数据 note）；supports_images 透出 |
| 4 | 活动页分组 + 奖励 + lottery 概要只读 | ✅ | Activities.tsx 重复性/单次双分组；lottery 只读无抽奖按钮（测试 `queryByRole('button', {name:/抽奖/})` 为 null） |
| 5 | 前端新测试 ≥15 + Any 清零 | ✅ | 新增 12 个（models.test.tsx），总计 112；`grep Any Models.tsx Activities.tsx` = 0 |
| 6 | 原子提交 test→feat | ✅ | 契约(866e770) → 后端 RED→GREEN(214955f/7ba5a8b) → 前端 RED→GREEN(99acc15/f62b7d1) |

## Hermes 抽查（防编造核实）

- **上游字段实测核对表**（报告 §3）逐条核实：`credits`/`supportsImages`/`tags`/`task_type` 确认上游存在（harness 报告+Apifox 导出为证）；数值单价/显式限免字段确认不存在 → `free` 为面板推导值并明确标注「非上游原话」；lottery 四端点存在但 schema 空白 → 字段多候选键提取 + `note: "上游字段未核实的占位"`
- **M1 编造字段已删除**：price/price_unit/free_quota 从契约和类型中清除，改为 credits 描述串透传——遵守了任务书「禁止编造」红线
- 前端覆盖率 90.51% 行 / 81.93% 分支（Models.tsx 100% / Activities.tsx 98.87%）
- 偏差 5 条核实合理（后端合并提交因接口共享 / 编造字段删除 / free 推导 / lottery 占位 / 联调项移交）

## 后端联调验证项（Hermes 记录，部署阶段执行）

1. GET /api/models/pricing 真实账号 credits 字段形态
2. GET /api/activities/lottery 真实抽奖概要字段名验证
3. 模型页/活动页在真实浏览器渲染

## 后续
- M5 任务书：OAuth 完整流 + 运行日志环形缓冲 + 安全硬化回归
