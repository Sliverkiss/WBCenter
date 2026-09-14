# 里程碑 M4：模型价格与限免（Issue #70 归并）+ 活动管理升级

> 项目：/root/WBCenter。先读 `docs/DEV-HANDBOOK.md`（第 1 节边界、第 4 节纪律），再读本任务书。
> M0-M3 已完成：四绿工具链 + 契约/类型/MSW + 账号池监控 + 一键任务/统计（100 测试）。

## 背景

### Issue #70 归并诉求
在 Web 面板中展示各模型请求价格与限时免费标识，方便不安装客户端即可随时查看。

### 上游接口现状
- `internal/upstream/client.go` 的 `FetchModels`（`/console/enterprises/personal/models` CN / `/v2/enterprises/personal/models` global）返回模型列表，已有 `credits`、`reasoning.supportedEfforts` 等字段，但**没有解析价格/限免相关字段**。
- Apifox 导出 `docs/reference/apifox-workbuddy-endpoints.md` 的 growth 域有模型/计费相关端点，可作为契约参考。
- 逆向分析报告 `/root/workbuddy2api/.claude/reports/deepseek-harness-analysis.md` §6.3 指出：上游模型定义含 `credits` 字段（如 `"credits":"x0.51 credits"`），`supportsImages` 字段也未解析。
- **如果上游确实没有价格字段**：在面板上做占位并明确标注「上游未提供价格数据」，不编造。

### 现有模型页
`web/src/pages/Models.tsx`：当前只展示按账号可用模型列表（模型 ID/名称/状态），无价格/限免信息。后端 `/api/models` 已有 `AvailableModels` 方法。

## 任务（TDD，test 先于 feat，前后端都要）

### A. 后端：模型列表字段扩展
1. `internal/upstream/client.go` 的 `FetchModels` / `ModelInfo` 结构：新增解析 `credits`（价格描述串）、`supportsImages`、`descriptionZh`/`descriptionEn` 字段（如果上游返回了的话）。**先查 Apifox 导出和 harness 分析报告确认上游真实返回的字段名**，不要猜。
2. `internal/control/` 的 `/api/models` handler：响应里带上新字段。
3. 如果上游确实没有独立的价格数值字段（只有 `credits` 描述串），在响应里原样透传该串，前端解析展示。如果连 `credits` 都没有，在响应里加 `pricing_available: false` 标记。
4. Go 测试：mock 上游模型列表 JSON（含 credits 字段），验证解析和透传。

### B. 后端：模型价格探测端点（契约已预定义）
1. 契约 `docs/api-contract.md` §2.4 `GET /api/models/pricing`——M1 已定义契约+mock，后端未实现。
2. 实现逻辑：遍历模型列表，提取每个模型的 `credits` 字段；如果上游模型定义里有限免标识（如 `tags` 含 `badge:限时免费` 之类），也提取。
3. 如果上游数据不可用或不完整：返回 `pricing_available: false` + 可用部分，不编造。
4. Go 测试：正常解析 / 上游无价格字段降级 / 部分模型有部分无。

### C. 前端：模型管理页升级
1. `web/src/pages/Models.tsx` 升级：模型卡片/表格增加价格列（credits 描述串）、限免徽标（如有）、图片支持标识（supportsImages）。
2. 无价格数据时显示「上游未提供」降级文案，不空白。
3. 新增「模型价格」独立视图或在现有模型页内加 tab/section，调 `GET /api/models/pricing` 展示价格汇总。
4. MSW handler 扩展：fixtures 含有价格/无限免/无价格三类模型样本。
5. 组件测试：价格列渲染、限免徽标、降级文案、空数据。

### D. 前端：活动管理升级
1. `web/src/pages/Activities.tsx` 升级：成长任务列表增加奖励展示、任务类型分组（single/重复性）、lottery 概览（只读：抽奖机会/记录/奖品概要）。
2. 后端如果需要新增端点支持 lottery 概要（契约参考 Apifox growth 域的 lottery 端点），按需添加但标注「上游字段未核实的占位」。
3. MSW handler 扩展：活动任务 + lottery 样本。
4. 组件测试：任务分组渲染、奖励展示、lottery 概要只读渲染。

### E. 契约同步
1. `docs/api-contract.md` 更新 `/api/models` 和 `/api/models/pricing` 的字段（根据实测上游真实返回）。
2. `src/types/` + `src/services/` + `src/mocks/` 同步。
3. 契约变更顺序：改文档 → 改类型 → 改 mock → 改组件。

## 验收标准

1. `bash scripts/check.sh` 四绿 + `go test -race ./internal/control/ ./internal/upstream/` 全绿。
2. `/api/models/pricing` 实现与契约 §2.4 一致；上游无价格数据时返回 `pricing_available: false`（不编造）。
3. 模型页展示价格列/限免标识/图片支持；无数据时降级文案。
4. 活动页任务分组 + 奖励 + lottery 概要只读渲染。
5. 前端新测试 ≥15 个；Any 在改动的页面清零。
6. 原子提交：后端模型扩展（test→feat）、后端 pricing（test→feat）、前端模型页（test→feat）、前端活动页（test→feat）、契约同步各自独立。

## 纪律
同手册第 4 节。特别强调：**上游字段不确定时先查 Apifox 导出文档和 harness 分析报告，查不到就占位并如实标注「上游未提供」——禁止编造**。只动 /root/WBCenter；禁 push/gh/重型 MCP；不引新依赖。

## 交付报告
`.claude/reports/m4-report.md` + stdout 全文：commits / 文件清单 / pricing 契约对齐表 / 上游字段实测核对表（哪些字段上游真有 vs 没有）/ RED→GREEN 证据 / 四绿+race 证据 / 覆盖率 / 偏差说明。
