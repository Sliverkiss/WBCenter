# WorkBuddy Control Center

**Version: 1.0.0**

本地、独立的 WorkBuddy 账号控制台（Web 管理面板）。它遵循 workbuddy2api Issue #61 的边界：

- 与网关的唯一共享面是 `auths/` 账号凭据目录（只读凭据 + 面板写入新账号）。
- 账号、积分、模型、活动和云端定时任务数据**直连 WorkBuddy 上游接口**（`copilot.tencent.com` / 国际版域），不经过网关。
- 不读取、不写入 `workbuddy2api/config.json`，不调用网关 HTTP 接口，不使用 Docker Socket。

UI 全中文；浅/暗双主题（薄荷绿 `#54c7ad` 设计 token 体系）。

## 功能矩阵

| 模块 | 能力 | 主要接口 |
|---|---|---|
| 登录与会话 | 本地口令登录、HttpOnly + SameSite=Strict 会话 Cookie、登录限速（窗口内超限 429）、退出登录 | `POST /api/login`、`POST /api/logout`、`GET /api/session` |
| 概览 | 健康/总数/过期/会话失效/已启用自动化统计卡、快速入口 | `GET /api/overview` |
| 账号池监控 | 账号卡片网格（Token 三态徽标、CN/Global 区域徽标、12153 会话失效徽标）、全量探测（签到状态/积分/猫猫旅行状态机逐账号聚合）、单账号操作（刷新凭据/签到/旅行巡检） | `GET /api/accounts`、`POST /api/probe`、`POST /api/accounts/{uid}/actions/{action}` |
| 积分管理 | 积分总览与按账号查询；今日区域展示上游「套餐发放/已用/剩余」，不把套餐切片伪造为交易流水 | `GET /api/credits`、`GET /api/credits/{uid}` |
| 模型管理 | 按账号直连上游查询可用模型；模型价格与限免视图（Issue #70；价格为上游 credits 描述串原样透传，上游未提供数值单价时明确标注「上游未提供」降级展示） | `GET /api/models`、`GET /api/models/pricing` |
| 统计分析 | Token 与积分汇总卡（今日发放/消耗/剩余、签到分布、状态分布、自动化运行成败） | `GET /api/stats/summary` |
| 自动化管理 | 面板本地持久化定时任务（启停/间隔/立即运行/最近运行记录）；一键任务面板（批量签到/旅行巡检/刷新凭据，禁用账号跳过、12153 单独标记、逐账号结果表） | `GET/PUT /api/automations/{id}`、`POST /api/automations/{id}/run`、`POST /api/batch-actions` |
| 账号云端定时任务 | 各账号 WorkBuddy 云端 scheduler 任务**只读**列表（`/console/as/tasks`）；上游 `access_denied` 时逐账号降级 | `GET /api/scheduler-tasks`、`GET /api/scheduler-tasks/{uid}/{taskID}` |
| 本地 Mock 任务骨架 | 上游 scheduler 创建接口契约未经验证前的占位 CRUD，仅写入面板 `data/state.json`，不提交上游 | `GET/POST /api/mock-scheduler-tasks`、`PUT/DELETE /api/mock-scheduler-tasks/{id}` |
| 活动管理 | 成长任务/新任务探测、抽奖概览（只读） | `GET /api/activities`、`GET /api/activities/lottery` |
| 添加账号（OAuth） | 国内版/国际版设备授权完整流（start → poll → account），11217 waiting 归一、倒计时、超时与错误终态处理 | `POST /api/oauth/start`、`POST /api/oauth/{id}/poll` |
| 运行日志 | 面板运行日志查看（环形缓冲、级别过滤、10s 轮询） | `GET /api/logs` |

## 安装与配置

### 配置文件

```sh
cp config.example.json control-center.json
# 编辑 control-center.json：至少设置 password 和 auth_dir
```

| 字段 | 默认 | 说明 |
|---|---|---|
| `listen` | `127.0.0.1:8787` | 监听地址。默认只绑本机 |
| `auth_dir` | `../workbuddy2api/auths` | 与网关共享的凭据目录（唯一共享面） |
| `data_dir` | `./data` | 面板自身状态目录（自动化配置、运行记录、日志） |
| `username` | `admin` | 面板登录用户名 |
| `password` | 无（必填） | 面板登录密码，不设置则拒绝启动 |
| `read_only` | `false` | 只读模式：禁用一切写操作（账号动作/OAuth/Mock CRUD/自动化变更） |
| `timezone` | `Asia/Shanghai` | 统计与自动化调度时区 |
| `timeout_seconds` | `30` | 上游请求超时（小于 5 秒按 30 处理） |

### 环境变量（优先级高于配置文件）

| 变量 | 对应字段 |
|---|---|
| `WBCC_LISTEN` | `listen` |
| `WBCC_AUTH_DIR` | `auth_dir` |
| `WBCC_DATA_DIR` | `data_dir` |
| `WBCC_USERNAME` | `username` |
| `WBCC_PASSWORD` | `password` |

命令行参数：`-config <路径>` 指定配置文件（默认 `control-center.json`）。

## 部署

### docker-compose（推荐）

```sh
WBCC_PASSWORD='独立强密码' docker compose up -d --build
```

compose 配置要点（见 `docker-compose.yml`）：

- 端口映射 `8787:8787`，`WBCC_LISTEN=0.0.0.0:8787`（容器内监听全网卡，宿主侧可用 `127.0.0.1:8787:8787` 收紧）。
- Volume：`../workbuddy2api/auths:/auths`（与网关共享凭据，唯一共享面）、`./data:/data`（面板状态持久化）。
- `WBCC_PASSWORD` 必须从本机环境注入，未设置时 Compose 拒绝启动。
- 容器以非 root（UID 10001）运行；`/data` 由镜像内 `chown` 保证可写，宿主 `./data` 目录若属 root 需先 `chown -R 10001:10001 data`。

### 纯 Docker

```sh
docker build -t workbuddy-control-center .
docker run --rm -p 127.0.0.1:8787:8787 \
  -e WBCC_PASSWORD='替换为强密码' \
  -e WBCC_AUTH_DIR=/auths \
  -v /绝对路径/workbuddy2api/auths:/auths \
  -v "$PWD/data:/data" \
  workbuddy-control-center
```

镜像为三阶段构建：前端（node:20-alpine）→ Go 编译（golang:1.23-alpine，前端 dist 经 `go:embed` 打入二进制）→ 运行（alpine:3.20，非 root，自带 `/api/session` healthcheck）。

访问 `http://127.0.0.1:8787`。

### 公网安全建议

默认只绑定本机。若经反向代理开放到局域网或公网：

- 必须上 HTTPS（会话 Cookie 依赖传输层保护）。
- 使用独立强密码（不要与网关或系统账号复用）。
- 用反代 ACL/IP 白名单收紧访问面；登录限速是最后一道防线，不是访问控制。
- 审计场景开 `read_only: true`，把面板降为纯观测台。

## 安全说明

- **登录限速**：同一来源在固定窗口内失败次数超限返回 429，自动解封。
- **同源写校验**：非 GET 请求校验 `Origin` 与 `Host` 同源，跨源返回 403（配合 SameSite=Strict Cookie 构成 CSRF 双保险）。
- **会话 Cookie**：`HttpOnly; SameSite=Strict; Path=/`，12 小时过期；logout 立即失效。
- **安全响应头**：`Content-Security-Policy`（default-src 'self'）、`X-Frame-Options: DENY` 等。
- **用户枚举防护**：登录失败统一「用户名或密码错误」，不区分用户不存在与密码错误。
- **只读模式**：`read_only: true` 时全部写接口返回 403。
- 凭据只落 `auths/`（与网关共享）与 `data/`（面板状态），均不依赖网关配置。

## 自动化所有权

为避免同一账号双跑，面板默认只启用「新活动探测」。如需面板接管签到、旅行、保活，必须先在部署层关闭网关中同名定时任务——这不改网关源码，但由操作者确认完成。面板内一键任务（批量签到/旅行/刷新）为手动触发，不受此限制。

## 开发

### 四绿检查

```sh
bash scripts/check.sh
# = go vet + go test + web: lint/typecheck/test/build
```

### 前端结构（`web/`）

```
src/
  components/   通用组件（Head/StatCard/Badge/BatchPanel/…）
  pages/        10 个页面（Overview/Accounts/Credits/Models/Stats/Automation/Scheduler/Activities/OAuth/Logs/Login）
  services/     fetch 封装（401 全局登出、错误信封解析）
  hooks/        useHashRoute/useTheme/useNotice/useEcharts
  types/        与 docs/api-contract.md 逐字段对齐（字段名照抄后端 JSON tag）
  mocks/        MSW handlers + fixtures（dev 与 test 双环境共用）
  styles/       设计 token（浅/暗双主题 CSS 变量）+ 全局样式
e2e/            Playwright 用例（不进 vitest）
```

### 测试命令

```sh
cd web
npm test            # vitest（jsdom + MSW，129 用例）
npm run coverage    # 覆盖率
npm run e2e         # Playwright E2E（自动拉起 dev server + MSW，5 用例，需先 npx playwright install chromium）
```

E2E 覆盖主链路：登录→概览统计卡、错误密码拦截、账号池卡片网格+全量探测、一键任务逐账号结果表、模型价格与限免视图。E2E 不进四绿，单独跑。

### 后端

```sh
go test ./...       # 含 race 安全的集成/单元测试
go build ./cmd/server
```

## 边界与免责

第三方本地工具，与 WorkBuddy 官方无关。所有数据来自上游接口直取或共享凭据目录；面板不修改网关任何配置。云端 scheduler 创建接口因请求体契约未经认证验证（且实测 `access_denied`）保持只读，本地 Mock 骨架不会同步上游，待权限与契约确认后以独立适配器接入。
