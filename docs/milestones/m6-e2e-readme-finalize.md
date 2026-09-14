# 里程碑 M6：Playwright E2E + README 重写 + 企业化收尾 v1.0.0

> 项目：/root/WBCenter。先读 `docs/DEV-HANDBOOK.md`（第 1 节边界、第 4 节纪律），再读本任务书。
> M0-M5 已完成：四绿工具链 + 契约/类型/MSW + 账号池监控 + 一键任务/统计 + 模型价格/活动 + OAuth/日志/安全（129 测试）。

## 背景

这是最后一个里程碑——把 demo 坯子打磨成企业级合格项目。M0-M5 已经把功能全部实现，M6 做端到端验证、文档收尾、构建验证和版本定标。

## 任务

### A. Playwright E2E
1. 引入 Playwright（`@playwright/test` devDependency）。`web/` 下新建 `e2e/` 目录。
2. E2E 测试通过真实 dev server（`npm run dev` + MSW 启用）跑，不依赖后端。
3. 关键用户流（至少 4 个）：
   - **登录→概览**：打开登录页 → 输入 admin/密码 → 进入概览 → 验证统计卡渲染
   - **账号池**：导航到账号管理 → 验证卡片网格渲染 → 点全量探测 → 验证结果
   - **一键任务**：导航到自动化 → 点一键任务按钮 → 验证逐账号结果表
   - **模型价格**：导航到模型管理 → 验证价格列/限免标识 → 切换 pricing 视图
4. Playwright 配置：`web/playwright.config.ts`（baseURL、CI 模式、截图/trace on failure）。
5. `scripts/check.sh` 加 `cd web && npx playwright test`（或 `npm run e2e` 脚本），但 E2E 只在 CI/手动跑（check.sh 保留四绿，E2E 单独 `npm run e2e`）。

### B. README 重写
1. 完全重写 `README.md`（现有的是 demo 坯子的简版），覆盖：
   - 项目定位（独立 Web 管理面板，遵循 workbuddy2api issue #61 边界：只共享 auths/、数据直取上游）
   - 功能矩阵（M0-M5 实现的全部能力，分模块列出）
   - 安装与配置（config.example.json 字段说明、环境变量 WBCC_* 列表、.env 用法）
   - 部署（docker-compose 一键部署、auths 目录映射、data 目录权限、公网安全建议）
   - 开发（四绿脚本、前端目录结构、测试命令、E2E 命令）
   - 安全说明（登录限速/同源校验/HttpOnly+SameSite/CSP/只读模式）
   - 自动化所有权（默认只开新活动探测，接管网关任务需部署层确认）
2. 功能矩阵必须与实际实现一致——对照 `docs/api-contract.md` 和 `internal/control/server.go` 路由表核实，不夸大不遗漏。

### C. Dockerfile/compose 复核
1. 检查 `Dockerfile` 三阶段构建（web → go build → alpine runtime）是否完整、无冗余。
2. 检查 `docker-compose.yml` 端口映射、volume 映射（auths + data）、环境变量。
3. 确保 `.dockerignore` 排除 `.claude/`、`web/node_modules/`、`web/coverage/` 等。
4. 如果有问题，修正。

### D. 构建验证
1. `docker build -t workbuddy-control-center .` 从零构建成功（不需要 docker compose，纯 Dockerfile）。
2. 运行容器 + healthcheck 通过。
3. 前端 dist 产物正确嵌入 go:embed（`internal/webui/dist/` 有 index.html + assets）。

### E. 版本定标
1. `git tag v1.0.0`（在 M6 最后一个 commit 上，不 push tag）。
2. README 顶部标注 `Version: 1.0.0`。

## 验收标准

1. `bash scripts/check.sh` 四绿（E2E 不进四绿，单独 `npm run e2e` 全绿）。
2. `npm run e2e` ≥4 个 E2E 测试全绿。
3. README 功能矩阵与 `server.go` 路由表 + `api-contract.md` 逐条核对一致。
4. `docker build` 从零构建成功 + healthcheck 通过。
5. `.dockerignore` 排除正确。
6. `git tag v1.0.0` 存在。
7. 原子提交：E2E（test→feat）、README、Dockerfile 复核、版本定标各自独立。

## 纪律
同手册第 4 节。Playwright 是唯一新 devDependency。只动 /root/WBCenter；禁 push/gh/重型 MCP。

## 交付报告
`.claude/reports/m6-report.md` + stdout 全文：commits / 文件清单 / E2E 测试清单与结果 / README 功能矩阵核对表 / Dockerfile 构建输出 / 版本 tag 证据 / 四绿+e2e 证据 / 偏差说明。
