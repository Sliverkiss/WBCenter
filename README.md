# WBCenter · WorkBuddy 本地控制台

独立运行的 WorkBuddy 账号管理与自动化控制台，面向本机或私有网络部署。项目遵循 `workbuddy2api` Issue #61 的边界：控制台只读取共享 `auths/` 凭据，直接访问 WorkBuddy 上游，不读取、修改或调用网关的 `config.json`、HTTP 服务或 Docker Socket。

## 功能

- 本地管理员登录：HttpOnly、SameSite=Strict 会话、同源写校验、登录失败限速。
- OAuth 添加国内版 / 国际版账号，凭据原子写入 `auths/`；支持多账号列表、刷新、签到、旅行和保活。
- 积分管理：当前积分、今日套餐额度、已用、剩余；按账号展开上游套餐明细，不伪造本地积分流水。
- 模型管理：读取上游模型，支持按账号和名称筛选。
- 活动管理：探测新成长任务；按账号汇总任务数量和全部奖励；查看任务明细；执行已验证的任务完成、回读确认和领奖流程；执行进度、耗时与中文明细日志可查询。
- 自动化管理：新活动探测、活动任务自动完成、签到、旅行巡检、账号保活，均可启停、调整间隔、立即运行并查看中文运行记录。
- 账号云端定时任务：只读查询 WorkBuddy 账号自己的云端任务和详情。
- Mock 任务：云端 scheduler 写入契约尚未稳定时，提供仅保存到本地 `data/state.json` 的创建、编辑、启停和删除骨架。

活动自动完成使用上游已验证的任务映射，任务完成后必须重新读取进度才会领奖。未知任务、真实捐款等不可安全自动化的任务会明确标记为人工完成；不会把 HTTP 200 直接当作业务成功。

## 与网关的关系

```text
浏览器 → WBCenter 控制台 → WorkBuddy 上游
                         └→ 只读挂载 workbuddy2api/auths
```

控制台与网关独立部署、独立升级。Compose 为活动任务复用网关仓库中已经验证的 `scripts/task_runner.py`，脚本目录以只读方式挂载；不会修改网关源码。部署目录默认结构如下：

```text
Workbuddy/
├── workbuddy2api/
└── workbuddy-control-center/
```

## 本机部署

准备 Docker、Docker Compose 和相邻的 `workbuddy2api` 仓库。首次部署生成本地配置并设置唯一强密码：

```sh
cp config.example.json control-center.json
# 编辑 control-center.json：至少替换 password，并确认 auth_dir 指向 ../workbuddy2api/auths
chmod 600 control-center.json
docker compose up -d --build
```

访问 <http://127.0.0.1:8787>。配置文件、账号凭据、运行状态和日志均已排除在 Git 及 Docker 构建上下文之外。默认仅监听本机；通过反向代理或 Tailnet 暴露时，应使用 HTTPS 和访问控制。

也可以用环境变量覆盖配置：`WBCC_LISTEN`、`WBCC_AUTH_DIR`、`WBCC_DATA_DIR`、`WBCC_USERNAME`、`WBCC_PASSWORD`、`WBCC_READ_ONLY`、`WBCC_TIMEZONE`、`WBCC_TIMEOUT_SECONDS`、`WBCC_TASK_RUNNER`。生产环境建议通过环境变量或受限权限的本地配置注入密码，不要将密码写入 README、日志或提交记录。

## 自动化默认值

首次启动默认启用“新活动探测”和“活动任务自动完成”（每日一次），其他自动化按当前配置决定。签到、旅行、保活和活动完成都会产生真实账号行为；启用前请确认没有其他服务执行相同任务，避免重复上报或触发频控。开启 `WBCC_READ_ONLY=true` 后，所有写操作和自动化立即被拒绝。

## 安全边界

- 仅使用本人授权的 WorkBuddy 账号，并妥善保护 `auths/` 中的明文凭据。
- 不提交 `control-center.json`、`auths/`、`data/`、令牌、Cookie、管理员密码或运行日志。
- 所有写操作要求管理员会话和同源请求；活动任务代码使用白名单，不接受任意上游路径或 payload。
- 云端 scheduler 当前保持只读；Mock 任务不会同步到 WorkBuddy。

## 开发与验证

```sh
# 前端
cd web && npm ci && npm run build

# Go（需要 Go ≥ 1.23）
go test ./...
```

项目采用 MIT 许可证。WorkBuddy 上游接口为非公开、可能变化的客户端接口，使用时请遵守目标平台服务条款和所在地法律。
