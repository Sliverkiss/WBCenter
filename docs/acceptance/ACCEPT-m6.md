# M6 Playwright E2E + README + 企业化收尾 —— 验收记录（Hermes）

日期：2026-09-15　基线：`1caebc0`(M5 验收) → `b710b6c`（4 个原子 commit + tag v1.0.0）

## 验收结论：✅ 通过 — M0-M6 全部完成

## 逐项核对（对照任务书 7 条验收标准）

| # | 验收标准 | 结果 | 证据 |
|---|---|---|---|
| 1 | 四绿（E2E 单独） | ✅ | Hermes 复跑：go vet/test 4 包 ok + web lint(0)/typecheck(0)/test(**129/129**, 14 文件)/build ✅ |
| 2 | E2E ≥4 全绿 | ✅ | `npx playwright test` → **5 passed (9.1s)**：登录→概览 / 错误密码 / 账号池探测 / 一键任务 / 模型价格 |
| 3 | README 功能矩阵与路由表逐条一致 | ✅ | 27 条路由全覆盖对照（server.go 27 条 = README 功能矩阵 12 模块）；配置字段 8 项 + WBCC_* 环境变量 5 个逐字段一致；安全说明 7 项与 M5 测试一致 |
| 4 | docker build 从零成功 + healthcheck | ✅ | 执行体报告：三阶段构建（npm ci→vite build→go build→alpine runtime），`docker ps` Up (healthy)，`GET /api/session` 200，`GET /` 200 `<title>WorkBuddy Control Center</title>`（go:embed SPA 正确嵌入） |
| 5 | .dockerignore 排除正确 | ✅ | 排除 .git/.claude/docs/data/.env/node_modules/coverage 等；保留 internal/webui/dist/index.html 占位（Dockerfile 阶段依赖） |
| 6 | git tag v1.0.0 | ✅ | `git tag --points-at HEAD` = `v1.0.0`，打在 b710b6c 上 |
| 7 | 原子提交 | ✅ | E2E RED→GREEN(5f174e3/80494e2) → README(a813720) → Dockerfile(b710b6c) |

## Hermes 抽查

- E2E 通过 dev server + MSW 跑（不依赖后端），mock 会话态 localStorage 双写设计合理（reload 后持久）
- README 159 行，功能矩阵 27 条路由 12 模块全覆盖，Version: 1.0.0 标注在顶部
- .dockerignore 关键修正：保留 `internal/webui/dist/index.html` 占位（Dockerfile COPY 依赖），初版过度排除已修正
- Vitest `exclude: ['e2e/**']` 防止 Playwright spec 被误收为 vitest 用例

## 全里程碑总结

| 里程碑 | 内容 | 测试数 | commit 数 |
|---|---|---|---|
| M0 | 工程基线（四绿/token双主题/组件化/路由） | 39 | 8 |
| M1 | 契约终稿+TS类型+MSW双环境 | 65(+26) | 7 |
| M2 | 账号池监控（probe/卡片网格/12153） | 82(+17) | 8 |
| M3 | 一键任务+统计（batch-actions/echarts） | 100(+18) | 10 |
| M4 | 模型价格与限免+活动升级 | 112(+12) | 5 |
| M5 | OAuth完整流+日志缓冲+安全回归 | 129(+17) | 7 |
| M6 | E2E+README+企业化收尾 | 129+5E2E | 4 |
| **总计** | | **129 单测 + 5 E2E** | **49 commit** |
