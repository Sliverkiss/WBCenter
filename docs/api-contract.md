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
    {
      "uid": "u-1001",
      "nickname": "阿明",
      "id": "wb-pro",
      "name": "WorkBuddy Pro",
      "credits": "x0.51 credits",
      "supports_images": true,
      "description_zh": "旗舰对话模型",
      "description_en": "Flagship chat model",
      "badges": ["限时免费"],
      "raw": { "id": "wb-pro" }
    },
    { "uid": "u-1002", "nickname": "阿红", "id": "", "name": "", "error": "上游模型查询失败" }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| id / name | string | 模型标识与显示名（从上游多候选键提取，提取不到为空串） |
| credits | string | 上游价格描述串原样透传（如 `"x0.51 credits"`）。**上游实测来源**：`/v2/enterprises/personal/models` 响应 `data.models[].credits`（harness 分析报告 §6.1/§6.3）。上游未返回该字段时为空串 |
| supports_images | bool | 是否支持图片输入。上游字段名 `supportsImages`（camelCase，harness buddy.ts:563-583 解析口径），本契约按面板惯例统一 snake_case 落盘。**上游未返回时为 false** |
| description_zh / description_en | string | 中英文描述。**上游字段未核实**：harness 分析报告未提及该字段，作为可选占位透传（多候选键 `descriptionZh/description_zh/desc` 与 `descriptionEn/description_en`），上游未提供时为空串 |
| badges | string[] | 限免/活动徽标列表。从上游 `tags` / `badges` 数组提取含「免费/限免/free/trial」关键词的条目原样透出；上游无此类条目标记时为空数组 |
| raw | object（可省略） | 上游原始条目（成功时存在） |
| error | string（可省略） | 该账号上游查询失败原因 |

### 1.9 `GET /api/activities`

成长任务（活动）列表，新发现任务带 `new: true`。

**响应 200：**

```json
{
  "items": [
    { "uid": "u-1001", "nickname": "阿明", "code": "T-DAILY", "name": "每日签到", "status": "可领取", "reward": 100, "task_type": "daily", "new": true }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| code / name / status | string | 上游成长任务编号、名称、状态 |
| reward | number | 奖励额度 |
| task_type | string | 任务类型。**上游实测来源**：`/v2/activity/growth/tasks` 响应 `data.tasks[].task_type`（Apifox 导出 growth 域 §成长任务清单）。用于前端「单次/重复性」分组展示；上游未返回时为空串 |
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
  "credits_today_remaining_total": 290,
  "checkin_done_today": 3,
  "checkin_pending_today": 2,
  "status_healthy": 3,
  "status_expired": 1,
  "status_session_dead": 1,
  "status_disabled": 1,
  "automation_runs_today_ok": 4,
  "automation_runs_today_failed": 1,
  "generated_at": "2026-09-14T08:30:00+08:00"
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| accounts_total | number | 账号总数 |
| credits_current_total | number | 全池当前积分合计（上游查询失败账号按 0 计入） |
| credits_today_allocated_total / consumed_total | number | 今日额度发放/消耗合计（同上口径） |
| credits_today_remaining_total | number | 今日剩余额度合计（同上口径） |
| checkin_done_today / pending_today | number | 今日已签/未签账号数。数据来源：最近一次 `POST /api/probe` 缓存快照的 `checkin.checked_in`（该字段本身的数据来源见 §2.1 标注）；未运行过探测时 done=0、pending=accounts_total |
| status_healthy / status_expired | number | Token 未过期/已过期账号数（凭据本地判定，不调上游） |
| status_session_dead | number | 最近一次探测被判 12153 的账号数（缓存；未探测为 0，与 /api/overview 同源） |
| status_disabled | number | 凭据标记 `disabled: true` 的账号数（见 §2.3 禁用约定）；与其他状态可重叠 |
| automation_runs_today_ok / failed | number | 面板本地自动化任务今日（服务端时区）成功/失败次数，来自 RunRecord 历史 |
| generated_at | string(时间) | 汇总生成时间 |

**数据来源与缓存策略（M3 定稿）：** 积分四项为**按需实时聚合**（与 `GET /api/credits` 同口径逐账号查询，单账号失败降级按 0 计入，不拖垮整体）；签到计数与 12153 分布**复用 lastProbe 缓存**（未探测时为保守零值，不主动触发上游调用）；自动化统计来自本地 RunRecord（无上游调用）。无历史趋势数据——趋势图待 M5 日志系统后补，本端点只返回当前快照。

**错误：** 401 未登录。

### 2.3 `POST /api/batch-actions`（M3 一键批量任务）

对账号池批量执行指定动作（签到/旅行/刷新凭据），逐账号返回结果。与 PR #60「一键任务」语义对齐：禁用账号跳过不调上游、12153 会话失效标记但不视为可重试失败。

**请求体：**

```json
{ "action": "checkin", "uids": ["u-1001", "u-1002"] }
```

| 字段 | 类型 | 语义 |
|---|---|---|
| action | string | 必填，`checkin` / `travel` / `refresh` 之一 |
| uids | string[]（可省略） | 目标账号子集；缺省或空数组 = 全量账号。不存在的 uid 静默忽略 |

**响应 200：**

```json
{
  "results": [
    { "uid": "u-1001", "nickname": "阿明", "ok": true, "message": "签到成功" },
    { "uid": "u-1002", "nickname": "阿红", "ok": true, "skipped": true, "message": "账号已禁用，已跳过" },
    { "uid": "u-1003", "nickname": "阿蓝", "ok": false, "session_dead": true, "error": "上游会话失效（12153）" },
    { "uid": "u-1004", "nickname": "阿绿", "ok": false, "error": "上游签到接口调用失败" }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| results[].uid / nickname | string | 账号标识与昵称 |
| results[].ok | bool | 该账号本次动作是否成功（skipped=true 时 ok=true，表示「按预期跳过」而非上游成功） |
| results[].message | string（可省略） | 成功或跳过时的说明文案 |
| results[].skipped | bool（可省略） | 该账号被跳过未调上游（当前唯一触发：账号已禁用） |
| results[].session_dead | bool（可省略） | 上游判 12153 会话失效（前端渲染徽标；不视为可重试失败） |
| results[].error | string（可省略） | 失败原因（ok=false 时必填） |

**账号禁用约定（M3 新增，凭据文件层）：** auths 目录下凭据 JSON 新增可选字段 `"disabled": true`（面板不主动写入，由运维手工标记或后续里程碑提供开关）；批量动作对禁用账号跳过且不计失败。禁用账号仍可被探测（探测是只读诊断），仅写动作跳过。

**错误：** 400 `不支持的动作`；401 未登录；403 只读模式/跨站。整批请求本身恒 200（单账号失败下沉到 results[].error，不中断整批）。

### 2.4 `GET /api/logs`（M5 面板运行日志）

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

### 2.5 `GET /api/models/pricing`（M4 模型价格与限免，Issue #70）

> **M4 实测定稿**。上游字段核对结论：
> - **存在**（harness 分析报告 §6.1/§6.3 实测）：`data.models[].credits`（价格描述串，如 `"x0.51 credits"`）、`supportsImages`、`tags`（含 `text-to-image` 等分类标签）。
> - **不存在**（Apifox 导出 + harness 分析报告均未发现）：按模型的数值单价（CNY/千次等）、`free_quota` 限免额度字段、`descriptionZh/En` 描述字段、显式「限时免费」徽标字段。
> - **候选限免信号**：`tags`/`badges` 数组若含「免费/限免/free/trial」关键词条目，原样透传到 `badges` 并置 `free=true`；上游未提供时 `free=false`。
>
> 因此本端点定稿为 **credits 描述串透传 + pricing_available 标记**：不编造数值单价与额度。

**响应 200：**

```json
{
  "items": [
    {
      "uid": "u-1001",
      "nickname": "阿明",
      "id": "wb-pro",
      "name": "WorkBuddy Pro",
      "credits": "x0.51 credits",
      "free": false,
      "badges": [],
      "supports_images": true,
      "note": ""
    },
    {
      "uid": "u-1001",
      "nickname": "阿明",
      "id": "wb-lite",
      "name": "WorkBuddy Lite",
      "credits": "x0 credits",
      "free": true,
      "badges": ["限时免费"],
      "supports_images": false,
      "note": ""
    }
  ],
  "pricing_available": true
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| items[].uid / nickname | string | 账号标识与昵称（与 1.8 对齐；同一模型可被多账号看到，按账号-模型成行） |
| items[].id / name | string | 模型标识与显示名 |
| items[].credits | string | 上游 `credits` 价格描述串原样透传；上游未提供时为空串 |
| items[].free | bool | 是否限免。`true` 的触发条件：`badges` 非空，或 credits 串归一化后形如 `x0`/`0 credits`。上游无显式限免字段，该字段为面板推导值，非上游原话 |
| items[].badges | string[] | 从上游 `tags`/`badges` 提取的限免/活动类徽标（关键词：免费/限免/free/trial，大小写不敏感） |
| items[].supports_images | bool | 是否支持图片输入（同 1.8） |
| items[].note | string | 补充说明；`credits` 与 `badges` 均为空时为 `该模型上游未提供价格与限免数据`，否则为空串 |
| pricing_available | bool | 账号池内是否至少一个模型返回了 `credits` 描述串。false 时前端展示「上游未提供模型价格数据」降级文案 |
| warning | string（可省略） | 全部账号上游查询失败或所有模型均无 credits 时为 `上游未提供模型价格数据`；此时 pricing_available=false |

**错误：** 401 未登录。上游失败按账号降级（该账号条目缺省），不返回 502。

### 2.6 `GET /api/activities/lottery`（M4 抽奖概览，只读）

各账号抽奖玩法概览（剩余次数 + 最近记录 + 已获奖品概要）。**只读**：不调用 `POST /activity/growth/lottery/draw`。

> **上游字段未核实占位**：Apifox 导出 growth 域存在 `lottery/summary`、`lottery/chances`、`lottery/draws`、`lottery/rewards` 四个端点，但响应 schema 均为空对象/空数组占位（`data: {}`），具体字段名未经真实响应确认。本契约以下字段为面板侧约定形状，聚合时按多候选键提取，提取不到给零值并在 `note` 标注。

**响应 200：**

```json
{
  "items": [
    {
      "uid": "u-1001",
      "nickname": "阿明",
      "chances": 2,
      "draws_total": 5,
      "recent": [
        { "prize": "积分 +10", "at": "2026-09-10 12:00" }
      ],
      "rewards_total": 1,
      "note": "上游字段未核实的占位"
    },
    {
      "uid": "u-1002",
      "nickname": "阿红",
      "chances": 0,
      "draws_total": 0,
      "recent": [],
      "rewards_total": 0,
      "note": "",
      "error": "上游抽奖概要查询失败"
    }
  ]
}
```

| 字段 | 类型 | 语义 |
|---|---|---|
| items[].uid / nickname | string | 账号标识与昵称 |
| items[].chances | number | 当前可用抽奖次数（上游 `lottery/chances` 多候选键提取，未提供为 0） |
| items[].draws_total | number | 抽奖历史总条数（上游 `lottery/draws` 的 `total` 字段，未提供为 0） |
| items[].recent | array | 最近抽奖记录（最多 5 条；prize/at 多候选键提取，未提供为空数组） |
| items[].rewards_total | number | 已获奖品总条数（上游 `lottery/rewards` 的 `total` 字段，未提供为 0） |
| items[].note | string | 占位标注：`上游字段未核实的占位`（任一字段为推导零值时给出） |
| items[].error | string（可省略） | 该账号上游查询失败原因；失败时其余字段给零值 |

**错误：** 401 未登录。单账号失败降级到该条目的 `error` 字段，整批恒 200。
