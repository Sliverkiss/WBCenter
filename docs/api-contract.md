# WBCenter API 契约（M1 终稿）

> 本文档是面板前端与 Go 后端之间的唯一契约权威。
> 存量端点字段与 `internal/control/server.go`、`internal/control/service.go`、`internal/control/state.go` 的 JSON tag 逐字一致（snake_case，不转 camelCase）。
> 时间字段为 RFC 3339 字符串（Go `time.Time` 的 JSON 序列化形式），如 `"2026-09-14T08:30:00+08:00"`。
> 鉴权：除 `POST /api/login` 与 `GET /api/session` 外，所有 `/api/*` 端点需携带会话 Cookie `wbcc_session`（HttpOnly、SameSite=Strict，12 小时过期）。

## 0. 通用约定

### 错误信封

所有非 2xx 响应统一为：

```json
{ "error": "中文错误信息" }
```

### 通用状态码

| 状态码 | 触发场景 | error 示例 |
|---|---|---|
| 401 | 未登录或会话过期（所有需鉴权端点） | `未登录或会话已过期` |
| 403 | 非 GET 请求 Origin 校验失败 | `跨站请求被拒绝` |
| 403 | 写操作且服务端开启只读模式 | `服务端已开启只读模式` |
| 400 | 请求体缺失或格式错误 | `缺少请求体` / `请求格式错误` |
| 502 | 上游 WorkBuddy 调用失败 | 各端点具体文案 |

写操作（POST/PUT/DELETE）另受 CSRF 同源校验：必须带 `Origin` 头且与 Host 同源。

---

## 1. 存量端点（22 条，对照 server.go 路由注册顺序）

### 1.1 `GET /api/session`

会话状态。无鉴权要求。

**响应 200：**

```json
{
  "authenticated": true,
  "read_only": false,
  "using_default_password": false,
  "timezone": "Asia/Shanghai"
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| authenticated | bool | 当前请求是否持有有效会话 |
| read_only | bool | 服务端只读模式开关（只读时所有写端点返回 403） |
| using_default_password | bool | 面板口令仍为出厂默认 `workbuddy` 时为 true（前端应提示改密） |
| timezone | string | 服务端配置的时区名 |

### 1.2 `POST /api/login`

面板登录。无鉴权要求。限速：同一来源 15 分钟内失败 10 次后返回 429。

**请求体：**

```json
{ "username": "admin", "password": "••••••" }
```

**响应 200：** `{"ok": true}`，同时 `Set-Cookie: wbcc_session=<token>; HttpOnly; SameSite=Strict; Path=/; Max-Age=43200`

**错误：**
- 401 `用户名或密码错误`
- 429 `登录失败次数过多，请在 15 分钟后再试`

### 1.3 `POST /api/logout`

登出。清除服务端会话并下发过期 Cookie。

**响应 200：** `{"ok": true}`

### 1.4 `GET /api/overview`

概览统计卡。

**响应 200：**

```json
{
  "accounts": 5,
  "healthy": 4,
  "expired": 1,
  "automations": 2,
  "session_dead": 1,
  "warnings": ["账号文件 alice.json 读取失败，已跳过"],
  "updated_at": "2026-09-14T08:30:00+08:00"
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| accounts | number | 账号总数 |
| healthy | number | Token 未过期账号数 |
| expired | number | Token 已过期账号数 |
| automations | number | 已启用的本地自动化任务数 |
| session_dead | number | 最近一次全量探测中被判 12153 会话失效的账号数（未运行过探测时为 0） |
| warnings | string[] | auths 目录读取告警（可空数组） |
| updated_at | string(时间) | 响应生成时间 |

### 1.5 `GET /api/accounts`

账号列表。

**响应 200：**

```json
{
  "accounts": [
    {
      "uid": "u-1001",
      "nickname": "阿明",
      "domain": "copilot.tencent.com",
      "expires_at": 1760000000,
      "expired": false,
      "needs_refresh": false
    }
  ],
  "warnings": []
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| uid | string | 账号唯一标识 |
| nickname | string | 昵称（凭据文件内） |
| domain | string | 上游域（国内 `copilot.tencent.com` / 国际版域） |
| expires_at | number | Token 过期时间（Unix 秒） |
| expired | bool | 是否已过期 |
| needs_refresh | bool | 距过期不足 24 小时，建议刷新 |
| warnings | string[] | auths 目录读取告警 |

### 1.6 `GET /api/credits` 与 1.7 `GET /api/credits/{uid}`

积分汇总。带 `{uid}` 时只返回该账号（仍包一层 `items` 数组）。

**响应 200：**

```json
{
  "items": [
    {
      "uid": "u-1001",
      "nickname": "阿明",
      "current": 1200,
      "today_allocated": 100,
      "today_consumed": 40,
      "today_remaining": 60,
      "packages": 2,
      "fetched_at": "2026-09-14T08:30:00+08:00",
      "error": "已取得当前积分；今日套餐明细暂不可用"
    }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| uid / nickname | string | 账号标识与昵称 |
| current | number | 当前剩余积分（上游 UserResource.Remain） |
| today_allocated | number | 今日免费套餐发放额度合计 |
| today_consumed | number | 今日已用额度合计 |
| today_remaining | number | 今日剩余额度合计 |
| packages | number | 套餐数量（上游 UserResource.Packages） |
| fetched_at | string(时间) | 本次上游拉取时间 |
| error | string（可省略） | 部分失败说明。`上游积分概要查询失败`（全失败）或 `已取得当前积分；今日套餐明细暂不可用`（降级）；成功时字段缺省 |

### 1.8 `GET /api/models`

按账号查询上游可用模型。

**响应 200：**

```json
{
  "items": [
    { "uid": "u-1001", "nickname": "阿明", "id": "wb-pro", "name": "WorkBuddy Pro", "raw": { "id": "wb-pro" } },
    { "uid": "u-1002", "nickname": "阿红", "id": "", "name": "", "error": "上游模型查询失败" }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| id / name | string | 模型标识与显示名（从上游多候选键提取，提取不到为空串） |
| raw | object（可省略） | 上游原始条目（成功时存在） |
| error | string（可省略） | 该账号上游查询失败原因 |

### 1.9 `GET /api/activities`

成长任务（活动）列表，新发现任务带 `new: true`。

**响应 200：**

```json
{
  "items": [
    { "uid": "u-1001", "nickname": "阿明", "code": "T-DAILY", "name": "每日签到", "status": "可领取", "reward": 100, "new": true }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| code / name / status | string | 上游成长任务编号、名称、状态 |
| reward | number | 奖励额度 |
| new | bool | 完成基线后首次探测到的新任务 |

### 1.10 `GET /api/scheduler-tasks`

各账号云端定时任务（只读）。

**响应 200：**

```json
{
  "items": [
    { "uid": "u-1001", "nickname": "阿明", "id": "task-1", "name": "每日摘要", "status": "active", "updated": "2026-09-01 10:00", "raw": {} },
    { "uid": "u-1002", "nickname": "阿红", "id": "", "name": "", "status": "", "error": "上游拒绝授权，请刷新凭据" }
  ],
  "write_enabled": false,
  "note": "这是各 WorkBuddy 账号的云端定时任务。创建接口请求体尚未经过真实认证验证，当前保持只读。"
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| items[].id / name / status / updated | string | 从上游多候选键提取；失败时为空串 |
| items[].error | string（可省略） | 按上游状态码区分：`上游拒绝授权，请刷新凭据`(401) / `上游拒绝访问，此账号没有云端任务权限`(403) / `上游暂未提供云端任务接口`(404) / `账号云端定时任务读取不可用`(其他) |
| write_enabled | bool | 恒为 false（只读承诺） |
| note | string | 只读说明文案 |

### 1.11 `GET /api/scheduler-tasks/{uid}/{taskID}`

单个云端任务详情（脱敏后透传：剔除 token/authorization/cookie/secret/password 键，数组截断 100 条）。

**响应 200：** `{"item": { ...上游原始字段（已脱敏）... }}`

**错误：** 502 `账号云端定时任务详情读取不可用`

### 1.12 `GET /api/mock-scheduler-tasks`

本地 Mock 任务列表（仅落 `data/state.json`，不提交上游）。

**响应 200：**

```json
{
  "items": [
    {
      "id": "mock-1726000000000000000",
      "account_uid": "u-1001",
      "name": "每日摘要",
      "cron": "0 9 * * *",
      "prompt": "汇总今日日程",
      "enabled": true,
      "created_at": "2026-09-14T08:00:00+08:00",
      "updated_at": "2026-09-14T08:00:00+08:00",
      "nickname": "阿明"
    }
  ],
  "mode": "mock",
  "note": "本地 Mock 骨架：仅写入控制台 data/state.json，不会提交到 WorkBuddy 上游。"
}
```

### 1.13 `POST /api/mock-scheduler-tasks`

**请求体：**

```json
{ "account_uid": "u-1001", "name": "每日摘要", "cron": "0 9 * * *", "prompt": "汇总今日日程", "enabled": true }
```

**响应 201：** 单个 MockSchedulerTask 对象（同 1.12 的 items[] 结构）。

**错误：** 400（校验失败：`任务名称需为 1 至 100 个字符` / `调度表达式需为 1 至 100 个字符` / `任务内容不能超过 4000 个字符`，或 account_uid 不存在）；403 只读模式。

### 1.14 `PUT /api/mock-scheduler-tasks/{id}`

**请求体：** 同 1.13（account_uid 字段被忽略，按路径 id 定位）。

**响应 200：** 更新后的 MockSchedulerTask 对象。

**错误：** 404 `Mock 任务不存在`；400 校验失败；403 只读模式。

### 1.15 `DELETE /api/mock-scheduler-tasks/{id}`

**响应 200：** `{"ok": true}`

**错误：** 404 `Mock 任务不存在`；403 只读模式。

### 1.16 `GET /api/automations`

本地自动化任务与最近运行记录。

**响应 200：**

```json
{
  "items": [
    {
      "id": "checkin",
      "name": "每日签到",
      "action": "checkin",
      "enabled": false,
      "every_minutes": 1440,
      "last_run_at": "2026-09-14T08:00:00+08:00",
      "last_result": "成功 3，失败 0",
      "next_run_at": "2026-09-15T08:00:00+08:00"
    }
  ],
  "runs": [
    { "at": "2026-09-14T08:00:00+08:00", "action": "checkin", "ok": true, "message": "成功 3，失败 0" }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| items[].id | string | 内置四任务之一：`activity-probe` / `checkin` / `travel` / `keepalive` |
| items[].action | string | 执行动作：`activity_probe` / `checkin` / `travel` / `refresh` |
| items[].every_minutes | number | 执行间隔（5 至 10080 分钟） |
| items[].last_run_at / last_result / next_run_at | （可省略） | 未运行/未启用时缺省 |
| runs | array | 最近运行记录（最多 100 条，新→旧） |

### 1.17 `PUT /api/automations/{id}`

**请求体：** `{"enabled": true, "every_minutes": 60}`

**响应 200：** 更新后的 Automation 对象（同 1.16 items[]）。

**错误：** 404 `自动化任务不存在`；400 `执行间隔需在 5 至 10080 分钟之间`；403 只读模式。

### 1.18 `POST /api/automations/{id}/run`

立即执行一次该自动化动作（作用于全部账号）。

**响应 200：** `{"ok": true, "message": "成功 3，失败 0"}`（activity_probe 时为 `已探测 N 条活动任务`）

**错误：** 404 `自动化任务不存在`；502 上游执行失败；403 只读模式。

### 1.19 `POST /api/accounts/{uid}/actions/{action}`

单账号动作。`action ∈ {checkin, travel, refresh}`，请求体为空。

**响应 200：** `{"ok": true, "message": "成功 1，失败 0"}`

**错误：** 400 `不支持的动作`；502 上游失败；403 只读模式。

### 1.20 `POST /api/oauth/start`

发起 OAuth 设备授权。

**请求体：** `{"region": "cn"}`（`cn` 或 `global`，其他值报 502）

**响应 200：** `{"id": "1726000000000000000", "url": "https://copilot.tencent.com/oauth/authorize?..."}`

**错误：** 502（区域非法或上游失败）

### 1.21 `POST /api/oauth/{id}/poll`

轮询授权结果（授权会话 10 分钟过期）。

**响应 200（等待中）：** `{"status": "pending"}`

**响应 200（成功）：** `{"status": "success", "uid": "u-1001", "nickname": "阿明"}`

**错误：** 502 `授权会话不存在或已过期` / `授权会话已过期` / `服务端已开启只读模式` / 上游失败

> 说明：1.20 与 1.21 计两条路由，存量合计 22 条（1.1–1.21 编号中 1.6/1.7 为两条路由）。

---

## 2. 新增端点契约（M2–M5 使用；M1 只定契约+mock，后端实现随后续里程碑补齐）

### 2.1 `POST /api/probe`（M2 全量探测）

对账号池逐账号执行「签到状态 + 积分 + 猫猫旅行」探测，聚合返回。

**请求体：** `{}`

**响应 200：**

```json
{
  "results": [
    {
      "uid": "u-1001",
      "nickname": "阿明",
      "ok": true,
      "checkin": { "checked_in": true, "streak_days": 7 },
      "credits": { "current": 1200, "today_remaining": 60 },
      "travel": { "status": "traveling", "location": "杭州", "departed_at": "2026-09-14T06:00:00+08:00", "arrive_at": "2026-09-14T18:00:00+08:00" }
    },
    {
      "uid": "u-1002",
      "nickname": "阿红",
      "ok": false,
      "error": "上游会话失效（12153）"
    }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| results[].ok | bool | 该账号探测是否整体成功；失败时 error 必填、其余子对象缺省 |
| results[].session_dead | bool | 该账号本次探测被上游判 12153 会话失效（前端据此渲染徽标） |
| results[].checkin.checked_in | bool | 今日是否已签到。**数据来源标注**：上游无独立「查签到」只读端点（实测结论链 Phase 3-B），此值来自面板本地自动化任务运行记录（最近一次 checkin 动作是否成功）；面板未运行过 checkin 时 checked_in=false、streak_days=0，为占位而非上游事实 |
| results[].checkin.streak_days | number | 连续签到天数。同上数据来源标注，当前实现恒 0 占位 |
| results[].credits.current / today_remaining | number | 探测时刻的积分快照 |
| results[].travel.status | string | 猫猫状态机：`idle` / `traveling` / `arrived`（来自上游 `state` 字段原值） |
| results[].travel.location / departed_at / arrive_at | （可省略） | ⚠️ **上游未提供**：实测 `GET /activity/growth/buddy/travel/status` 响应只含 `state/daily_limit_reached/record_id/reward_credit`（见 `internal/upstream/client.go` TravelState 与 Apifox 导出 growth 域），无地点与起止时间。本契约保留字段以备上游日后透出，当前后端不返回这些字段，前端按缺失降级渲染 |
| results[].error | string（可省略） | 该账号失败原因（逐账号降级，不拖垮整批） |

**错误：** 401 未登录；403 只读模式。整批请求本身恒 200，失败下沉到 results[].error。

### 2.2 `GET /api/stats/summary`（M3 统计汇总卡）

**响应 200：**

```json
{
  "accounts_total": 5,
  "credits_current_total": 8600,
  "credits_today_allocated_total": 500,
  "credits_today_consumed_total": 210,
  "checkin_done_today": 3,
  "checkin_pending_today": 2,
  "generated_at": "2026-09-14T08:30:00+08:00"
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| accounts_total | number | 账号总数 |
| credits_current_total | number | 全池当前积分合计 |
| credits_today_allocated_total / consumed_total | number | 今日额度发放/消耗合计 |
| checkin_done_today / pending_today | number | 今日已签/未签账号数 |
| generated_at | string(时间) | 汇总生成时间 |

**错误：** 401 未登录。

### 2.3 `GET /api/logs`（M5 面板运行日志）

面板后端环形缓冲日志快照（新→旧，最多 200 行）。敏感字段（token/cookie/密码）在后端写入缓冲前已脱敏。

**响应 200：**

```json
{
  "lines": [
    { "at": "2026-09-14T08:30:00+08:00", "level": "info", "message": "账号 u-1001 签到成功" },
    { "at": "2026-09-14T08:29:10+08:00", "level": "warn", "message": "账号 u-1002 上游会话失效（12153）" }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| lines[].at | string(时间) | 日志时间 |
| lines[].level | string | `info` / `warn` / `error` |
| lines[].message | string | 已脱敏的日志内容 |

**错误：** 401 未登录。

### 2.4 `GET /api/models/pricing`（M4 模型价格与限免，Issue #70）

> ⚠️ 上游字段未核实，M4 实测定稿。Apifox 导出中计费域存在套餐/价格相关端点，但「按模型的价格与限免标识」字段名尚未在真实响应中确认；本契约为占位结构，M4 实测后允许调整字段名（届时同步修订本文档与 `src/types/`）。

**响应 200：**

```json
{
  "items": [
    {
      "id": "wb-pro",
      "name": "WorkBuddy Pro",
      "free": false,
      "price": 9.9,
      "price_unit": "CNY/天",
      "free_quota": 0,
      "note": "上游字段未核实，M4 实测定稿"
    },
    {
      "id": "wb-lite",
      "name": "WorkBuddy Lite",
      "free": true,
      "price": 0,
      "price_unit": "",
      "free_quota": 100,
      "note": "限免：每日 100 次"
    }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| id / name | string | 模型标识与显示名（与 1.8 对齐） |
| free | bool | 是否限免 |
| price | number | 单价（限免为 0） |
| price_unit | string | 计价单位文案，如 `CNY/天`、`CNY/千次` |
| free_quota | number | 限免额度（0 表示无） |
| note | string | 补充说明（含占位标注） |

**错误：** 401 未登录；上游价格端点不可用时返回 200 且 items 为空数组 + 顶层 `"warning": "上游未提供模型价格数据"`（降级而非 502）。
