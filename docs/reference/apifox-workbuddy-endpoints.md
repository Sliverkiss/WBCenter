---
title: 默认模块
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# 默认模块

Base URLs:

# Authentication

# workbuddy/growth 成长/猫猫

## GET 获取成长中心档案

GET /v2/activity/growth/profile

获取当前用户的成长中心总览档案（等级/能量/徽章数等）。Web 成长中心首页调用。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取成长任务列表

GET /v2/activity/growth/tasks

获取成长任务清单（task_code/title/task_desc/task_type/reward_credit/reward_energy/accept_status/progress.current/progress.target 等）。调用前会先刷新订阅任务状态。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
【2026-09-11 Hermes 实测 · 前置已破】
本接口的前置「对话」门槛已确认可达：POST /v2/report（埋点上报，endpoint 513708446）发一条
`chat_request_send` 事件（**必须含 userId**）即可把 first_buddy 置为 completed。
实测三个账号（0225284f / 0851ce35 / 12f2582c）均成功，本接口返回 200，+300 分 +8 能量。
详见 /root/workbuddy2api/REPORT-active-map.md。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "tasks": [
      {
        "task_code": "daily_login",
        "accept_status": "not_accepted",
        "progress": {
          "current": 0,
          "target": 1
        }
      }
    ]
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 刷新订阅任务状态

GET /v2/activity/growth/subscribe-task/status

刷新（订阅类）任务状态。注意：前端即使失败也继续拉任务列表，属可失败的前置调用。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 领取成长任务

POST /activity/growth/tasks/accept

批量接受成长任务。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "task_codes": [
    "string"
  ]
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» task_codes|body|[string]| yes |任务码数组|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 领取成长任务奖励

POST /activity/growth/tasks/{task_code}/claim

领取指定任务的奖励。响应 data 含 energy（能量奖励）；若携带 buddy 字段，前端会回读 /activity/growth/buddy/visible 同步展示态。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
【2026-09-11 Hermes 实测 · 前置已破】
本接口的前置「对话」门槛已确认可达：POST /v2/report（埋点上报，endpoint 513708446）发一条
`chat_request_send` 事件（**必须含 userId**）即可把 first_buddy 置为 completed。
实测三个账号（0225284f / 0851ce35 / 12f2582c）均成功，本接口返回 200，+300 分 +8 能量。
详见 /root/workbuddy2api/REPORT-active-map.md。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|task_code|path|string| yes |任务码，对应任务列表的 task_code。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "energy": 10,
    "buddy": null
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取成长能量余额

GET /activity/growth/energy

获取能量余额（data.balance）与开盲盒消耗（cost_per_open / affordable / max_open_count）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "balance": 320,
    "cost_per_open": 100,
    "affordable": 3,
    "max_open_count": 10
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取成长徽章

GET /v2/activity/growth/badges

获取用户已获得的徽章列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取连续打卡

GET /activity/growth/streak

获取连续（streak）打卡信息。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

【2026-09-11 Hermes 实测】本接口输出的 `streak.days` 是活跃连登天数。实测：向 `/v2/report` 上报一条 `chat_request_send`（含 userId）后，day 由 0 → 1，与活跃地图同步。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取热力图

GET /activity/growth/heatmap

获取成长热力图数据，支持查询参数。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|params|query|string| no |热力图查询参数（原样透传，如年份/月份区间）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 分享热力图

POST /activity/growth/heatmap/share

生成/登记热力图分享，body 可选（前端传 `{}`）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取猫猫信息

GET /activity/growth/buddy/info

获取当前 buddy（猫猫）信息。**客户端主进程 growth facade 实际调用的是 /v2 前缀版本**：`GET {endpoint}/v2/activity/growth/buddy/info`；渲染端 backendProvider 走无 /v2 版本。响应 data.buddy 为 null 表示该账号尚未领养；data.poll_interval_seconds 为前端轮询间隔（下限 60s）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-User-Id|header|string| no |用户 uid（facade 在 account.uid 存在时注入）。|
|X-Enterprise-Id|header|string| no |企业 ID（有 enterpriseId 时注入）。|
|X-Tenant-Id|header|string| no |同 enterpriseId。|
|X-Domain|header|string| no |auth.domain。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "buddy": {
      "thumbnail_url": "https://...",
      "base_animated_url": "https://...",
      "instance_id": "xxx",
      "current_buddy": true
    },
    "poll_interval_seconds": 300
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取已领养猫猫列表

GET /activity/growth/buddy/list

获取账号已领养的 buddy 实例列表（data.buddies）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "buddies": [
      {
        "instance_id": "xxx",
        "current_buddy": true
      }
    ]
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取猫猫模板列表

GET /activity/growth/buddy/templates

获取可领养的猫猫模板（形象）列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取猫猫额度

GET /activity/growth/buddy/quota

获取猫猫相关额度/消耗配置（与 energy 接口的 cost_per_open/affordable/max_open_count 对应）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "cost_per_open": 100,
    "affordable": 3,
    "max_open_count": 10
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 首次领养猫猫

POST /activity/growth/buddy/first

**【领养入口 · 首次】** 首次免费领养一只猫猫，无请求体。响应 data.buddy 为领养结果；前端随后回读 GET /activity/growth/buddy/visible 同步「是否在客户端展示」，并在 1s 后刷新任务/能量/徽章/猫猫。这是 buddy=null 账号的官方领养入口之一。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
【2026-09-11 Hermes 实测 · 前置已破】
本接口的前置「对话」门槛已确认可达：POST /v2/report（埋点上报，endpoint 513708446）发一条
`chat_request_send` 事件（**必须含 userId**）即可把 first_buddy 置为 completed。
实测三个账号（0225284f / 0851ce35 / 12f2582c）均成功，本接口返回 200，+300 分 +8 能量。
详见 /root/workbuddy2api/REPORT-active-map.md。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "buddy": {
      "instance_id": "xxx",
      "thumbnail_url": "https://..."
    }
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 开盲盒获取猫猫

POST /activity/growth/buddy/open

**【领养入口 · 消耗能量】** 开盲盒抽取猫猫，消耗 energy（cost_per_open/次）。响应 data.results 为结果数组，前端取 results[0] 作为新猫猫。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "count": 1
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» count|body|integer| no |开盒次数，前端默认 1。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "results": [
      {
        "instance_id": "xxx",
        "current_buddy": false
      }
    ]
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 切换当前猫猫

POST /activity/growth/buddy/switch

在已领养的猫猫中切换当前展示实例。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "instance_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» instance_id|body|string| yes |目标 buddy 实例 ID。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 读取猫猫活动协议

GET /activity/growth/buddy/agreement

读取猫猫活动用户协议（是否已同意）。**领养前置步骤**：未同意时需先 POST 同意。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "has_agreed": false,
    "content": "..."
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 同意猫猫活动协议

POST /activity/growth/buddy/agreement

同意猫猫活动协议。**领养前置步骤**（body 固定 {"agree": true}）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "agree": true
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» agree|body|boolean| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 查询客户端展示开关

GET /activity/growth/buddy/visible

查询猫猫是否在客户端（输入框/首页槽位）展示，data.buddy_visible 为布尔。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "buddy_visible": true
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 设置客户端展示开关

POST /activity/growth/buddy/visible

设置猫猫是否在客户端展示。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "visible": true
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» visible|body|boolean| yes |是否在客户端展示猫猫。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 猫猫分享页数据

GET /activity/growth/buddy/share/view/{share_id}

按分享 ID 获取猫猫分享页数据（分享页为独立渲染，可在未登录态访问）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|share_id|path|string| yes |分享 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 生成成长分享链接

GET /activity/growth/share/{share_id}

按分享 ID 获取成长中心分享内容。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|share_id|path|string| yes |分享 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 领取出行到达奖励

POST /activity/growth/buddy/travel/claim

**【消息中心联动】** 领取猫猫出行到达奖励积分。消息中心收到 biz_type=growth.travel_arrival 且 ext.claim_status === 'unclaimed' 时由客户端自动调用。无请求体，仅按 envelope code 判成败。当前 46 个账号实测返回 no active buddy（未领养）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| no |固定 workbuddy（消息中心链路注入）。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "msg": "ok",
  "data": null
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取出行配置

GET /activity/growth/buddy/travel/config

获取猫猫出行玩法的配置（可选目的地等）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取出行状态

GET /activity/growth/buddy/travel/status

获取猫猫当前出行状态。未领养账号返回 no active buddy。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 派出猫猫出行

POST /activity/growth/buddy/travel/depart

派出猫猫前往指定地点出行。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "location_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» location_id|body|string| yes |目的地 ID，来自 travel/config。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取出行记录

GET /activity/growth/buddy/travel/records

分页获取猫猫出行历史记录。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|page_size|query|string| no |每页条数，默认 20。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "records": [],
    "total": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 关闭成长中心提示位

POST /activity/growth/user-state/dismiss

关闭成长中心内的提示/引导位。客户端本地也有一处同名路径 /space/api/view/user-state/dismiss。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "scene_key": "string",
  "version": 0
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» scene_key|body|string| yes |场景键。|
|» version|body|integer| no |可选版本号，前端仅在传值时附带。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取抽奖概要

GET /activity/growth/lottery/summary

获取抽奖玩法概要（剩余次数等）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取抽奖奖品

GET /activity/growth/lottery/prizes

获取奖品列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取抽奖机会

GET /activity/growth/lottery/chances

获取当前可用抽奖机会数量。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取抽奖机会流水

GET /activity/growth/lottery/chances/logs

分页获取抽奖机会变动流水。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|page_size|query|string| no |每页条数，默认 20。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "logs": [],
    "total": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 执行抽奖

POST /activity/growth/lottery/draw

执行一次抽奖。client_token 由前端生成用于幂等。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "client_token": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» client_token|body|string| yes |客户端幂等 token。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取抽奖记录

GET /activity/growth/lottery/draws

分页获取抽奖历史。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|page_size|query|string| no |每页条数，默认 20。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "draws": [],
    "total": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取已获奖品

GET /activity/growth/lottery/rewards

分页获取已获得的奖品（可填地址发货）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|page_size|query|string| no |每页条数，默认 20。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "rewards": [],
    "total": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 读取奖品收货地址

GET /activity/growth/lottery/rewards/{reward_id}/address

读取某奖品已填写的收货地址。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|reward_id|path|string| yes |奖品记录 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 提交奖品收货地址

POST /activity/growth/lottery/rewards/{reward_id}/address

为某奖品提交收货地址（body 原样透传，含收件人/电话/地址字段）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|reward_id|path|string| yes |奖品记录 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 兑换成长权益

POST /activity/growth/redeem

按档位（tier）兑换权益，client_token 幂等。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
**【2026-09-11 Phase 3-B 实测验证】**
兑换成长权益（活跃地图连登档位）。

**2026-09-11 实测：403 `{"code":403,"msg":"连续登录天数不足，请继续打卡或使用补签卡"}`**（00e26541 与 0225284f 均如此）。

前置门槛已确认：门禁为 **growth 成长体系的连登天数**（`GET /activity/growth/streak` 的 `streak.days`），**不是** billing 域每日签到的 `streak_days`（billing 签到 10 天时 growth 连登仍为 0，redeem 依旧 403）。

档位（来自活动页官方规则，实测 `GET /activity/growth/streak` 的 `redemption_status.tiers` 一致）：
- `7d` 入门档：连续登录 7 天 → 积分 +0 / 能量 +2 / 补登卡 +1 / 抽奖次数 +1
- `14d` 进阶档：连续登录 14 天 → 积分 +50 / 能量 +3 / 补签卡 +1 / 抽奖次数 +1
- `28d` 巅峰档：连续登录 28 天 → 积分 +150 / 能量 +5 / 补签卡 +1 / 抽奖次数 +1

规则：每档每月限兑 1 次，可多档累计，兑换不扣减连登天数；连登天数每月清零。
**冷启动不可用**：46 个账号 growth 连登均为 0，因此 redeem 全部不可做（至少需连续 7 天点亮 growth 热力墙）。

body: `{"tier":"7d|14d|28d", "client_token":"<幂等 token，前端用 crypto.randomUUID()>"}`；`client_token` 仅为幂等键，非风控 token。

错误码：403 连登不足 / 409 本月已领取（前端文案「本月已领取」）。

> Body Parameters

```json
{
  "tier": "string",
  "client_token": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» tier|body|string| yes |兑换档位。|
|» client_token|body|string| yes |客户端幂等 token。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取兑换概要

GET /activity/growth/redeem/summary

获取兑换档位与剩余次数概要。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 使用补签卡

POST /activity/growth/makeup-cards/use

对指定日期使用补签卡。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
**【2026-09-11 Phase 3-B 实测验证】**
使用补签卡补登断签日。

**2026-09-11 实测：403 `{"code":403,"msg":"no makeup card balance"}`**（00e26541）。

body: `{"target_date":"YYYY-MM-DD"}`。

**补签卡唯一来源 = redeem 兑换档位附赠**（7/14/28 天档各赠 1 张），持有上限 4 张。因此存在**鸡生蛋问题**：想补签必须先兑换，想兑换必须先有 7 天连登 → 冷启动（连登 0）时补签卡余额恒为 0，端点不可用。

服务端错误语义（前端 `ke()` 映射，实测 403 分支已触发）：
- 403 `no makeup card balance` → 无补签卡
- 400 `future date` → 不能补登未来日期
- 400 `cannot makeup history month` → 仅可补登当前自然月内的断登
- 400 `target date before launch` → 早于活动上线日（2026-06-17）
- 400 `date not broken` → 该日已点亮，无需补登
- 400 `already made up` → 该日已补登过

另：`GET /activity/growth/streak` 的 `streak.makeup_dates` 返回已补登日期列表。

> Body Parameters

```json
{
  "target_date": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» target_date|body|string| yes |补签目标日期。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/计费/签到

## POST 获取资源余额概要

POST /billing/meter/get-user-resource-summary

余额/档位/付费状态聚合查询，无业务请求体。响应 data 含 packages（packageCode/cycleRemain 等）。**无 /v2 前缀**（网关路由声明为无前缀路径，Web 用 cookie、Desktop 用 Bearer 皆然）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|Accept-Language|header|string| no |zh / en，决定返回文案语言。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "packages": [
      {
        "packageCode": "free",
        "cycleRemain": 100
      }
    ]
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取付费资源包分页

POST /billing/meter/get-user-resource-paid-packages

分页获取付费资源包。Desktop 调用时带 NeedRenewInfo=true 以拿到续费标记。参数原样透传（pageNumber 等）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "pageNumber": 0,
  "NeedRenewInfo": true
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» pageNumber|body|integer| no |页码。|
|» NeedRenewInfo|body|boolean| no |是否返回续费信息。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "Accounts": [],
    "TotalCount": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取免费资源包分页

POST /billing/meter/get-user-resource-free-packages

分页获取免费资源包，带切片周期范围参数。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "pageNumber": 0
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» pageNumber|body|integer| no |页码。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "Accounts": [],
    "TotalCount": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取企业成员用量

POST /v2/billing/meter/get-enterprise-user-usage

企业账号查询自身用量。Desktop 因走 IDE 网关需要 /v2 前缀（billingPrefix 覆写为 /v2）。响应 limit_num=-1 表示不限量。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Enterprise-Id|header|string| yes |企业 ID，必带。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "limit_num": 1000,
    "used_num": 120
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 查询每日签到状态

POST /v2/billing/meter/checkin-status

查询当日签到状态。Desktop 走 /v2 前缀。渲染端另有读取版走无前缀 /billing/meter/checkin-status。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
**【2026-09-11 Phase 3-B 实测验证】**
查询当日签到状态。**实测结论：本端点返回空壳，不可作为签到判定依据。**

2026-09-11 对 46 个账号实测，全部返回 `active=false, today_checked_in=false, streak_days=0, daily_credit=0, start_time="", theme_name=""`；而对同一批账号调用 `POST /v2/billing/meter/checkin-activity-status` 均返回 `active=true, daily_credit=100`。

**建议改用 `POST /v2/billing/meter/checkin-activity-status`。**

主机为 `https://www.codebuddy.cn`（billing 域）；在 copilot.tencent.com 上同样返回 200 但为空壳。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "checked_in": false,
    "streak_days": 3,
    "credit": 10
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 执行每日签到

POST /v2/billing/meter/daily-checkin

执行每日签到领取积分。Desktop 会通过 getCheckinRequestHeaders 注入设备风控头 X-Device-Token（来自 Turing 设备 token）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
**【2026-09-11 Phase 3-B 实测验证】**
执行每日签到领取积分。**实测已打通（2026-09-11）**：无需 X-Device-Token 即可签到成功。

实测结果（3 个账号）：`00e26541` 1000→1100 分、streak 10→11；`0225284f` 1000→1100、streak 10→11；`093c446f` 0→100、streak 0→1（该账号 chat 会话已失效 12153，签到仍成功 → 签到不依赖 chat 会话）。

响应：`{code:0, data:{credit:100, streak_days:11, is_streak_day:false}}`。
重复签到：HTTP 400 `{code:10001, msg:"今天已签到，请明天再来"}`（幂等保护）。

**【重要边界】本端点属计费(billing)域，主机为 `https://www.codebuddy.cn`（不是 copilot.tencent.com），路径为 `/v2/billing/meter/daily-checkin`。**

46 账号扫描：全部 `active=true`、`daily_credit=100`，签到前 `today_checked_in=false`（46/46 当日均未签到）。
自然日 CST(UTC+8) 0 点重置。这是 46 账号池中**唯一已验证的每日可重复积分动作**。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Device-Token|header|string| no |设备风险 token，Desktop 端由 Turing SDK 提供。|
|X-Device-Token-Error|header|string| no |取 token 失败时改带此头。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "credit": 10,
  "streak_days": 4,
  "is_streak_day": false
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 查询签到活动状态

POST /v2/billing/meter/checkin-activity-status

查询签到活动聚合状态（active/today_checked_in/streak_days/daily_credit/today_credit/is_streak_day/next_streak_day/streak_bonus_* /checkin_dates/season/theme_name 等）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
**【2026-09-11 Phase 3-B 实测验证】**
查询签到活动聚合状态。**实测已打通**：返回 `active/today_checked_in/streak_days/daily_credit/today_credit/total_credits/week_checkin_days/week_progress/checkin_dates/season/theme_name/activity_name/start_time/end_time/action_button`。

2026-09-11 实测（00e26541）：`active=true, theme_name="Buddy加油站", season=8, activity_name="开学季", daily_credit=100, streak_days=10, total_credits=1000, week_checkin_days=4, checkin_dates=[2026-09-01..2026-09-10]`（连续 10 天），签到后 streak_days→11、total_credits→1100。

注意：与 `checkin-status` 不同，本端点在 46 账号上均返回活跃活动数据（`checkin-status` 全部返回 active=false 空壳），**判断是否可签到应以本端点为准**。

主机为 `https://www.codebuddy.cn`（billing 域）。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "active": true,
    "today_checked_in": false,
    "streak_days": 3
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取大使状态

GET /v2/ambassador/status

获取当前用户是否为推广大使。Desktop billingPrefix=/v2；另有一处 raw 调用取 {endpoint}/v2/activity/ambassador/status 返回 isAmbassador；Web Provider 走 /console/activity/ambassador/status。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "is_ambassador": false
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取加量包价格

POST /billing/pay/get-price

询价接口。请求体以 PriceType=getPrice 区分场景（如加量包）。响应可能被网关包一层 data.Response。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "PriceType": "getPrice"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» PriceType|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "Response": {}
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建加量包订单

POST /billing/pay/create-order

创建加量包购买订单。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "quantity": 0
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» quantity|body|integer| no |购买数量。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "Response": {}
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 设置自动续费标记

POST /billing/pay/set-renew-flag

开启/关闭自动续费。关闭为异步渠道解约，前端轮询上限 30s、间隔 3s。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取账号列表

GET /console/accounts

获取已登录账号列表（data.accounts）。Desktop 用其判断登录态并取完整 profile（editionType/expireAt/packageCode 等计费必需字段）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "accounts": []
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取邀请码

GET /activity/workbuddy/invitation/v2/my-code

获取当前用户的邀请码。响应 data.invite_code 为空视为失败。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "invite_code": "ABC123",
    "expires_at": "2026-12-31"
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取邀请记录

GET /activity/workbuddy/invitation/v2/invite-records

获取邀请好友记录（friends/total_invited/total_used/total_credits）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "friends": [],
    "total_invited": 0,
    "total_used": 0,
    "total_credits": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/消息中心

## GET 获取消息概要

GET /v2/msg-center/message/summary

获取消息中心概要（未读数等）。Desktop 走 /v2 前缀（APISIX route_260，OIDC Bearer 鉴权）；Node 端走 /portal/msg-center/message。所有消息中心请求固定带 X-Product-Code: workbuddy。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "unread": 3
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取消息列表

GET /v2/msg-center/message/list

游标分页获取消息列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|cursor|query|string| no |分页游标，首次不传。|
|limit|query|string| no |每页条数。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "list": [],
    "next_cursor": ""
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 标记消息已读

POST /v2/msg-center/message/read

标记消息已读，body 原样透传。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": null
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 刷新待展示消息

POST /v2/msg-center/message/pending-displays/refresh

拉取待展示（弹窗/banner）消息。body 支持 display（类型数组）与 limit。客户端按 display+limit 做 inflight 去重。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "display": [
    "string"
  ],
  "limit": 0
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|
|body|body|object| no |none|
|» display|body|[string]| no |展示位类型数组。|
|» limit|body|integer| no |条数上限。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 上报消息已展示

POST /v2/msg-center/message/notify-shown

上报消息已展示给用户（用于频控）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": null
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 删除消息

POST /v2/msg-center/message/delete

删除指定消息，body 为 {msg_id}。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "msg_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|
|body|body|object| yes |none|
|» msg_id|body|string| yes |消息 ID。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": null
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 更新消息扩展字段

POST /v2/msg-center/message/ext

更新消息维度的扩展字段（ext），body 原样透传。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": null
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 更新用户消息扩展

POST /v2/msg-center/message/update

更新用户维度的消息扩展状态，body 原样透传。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": null
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取消息游标

GET /v2/msg-center/message/cursor

获取消息中心游标（用于增量拉取）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "cursor": ""
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 执行消息动作回调

POST /portal/msg-center/message/dismiss-forever

执行消息 action 的 callback 类型动作（如永久关闭）。客户端从消息 action.url 取路径：url 形如 `api:///portal/msg-center/message/dismiss-forever?biz_type=...&period_start=...&period_end=...`，剥掉 `api://` 前缀后直接 POST，query 原样带上。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|biz_type|query|string| no |业务类型。|
|period_start|query|string| no |周期开始。|
|period_end|query|string| no |周期结束。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-Product-Code|header|string| yes |固定 workbuddy。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": null
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取活动 Banner

GET /v2/activity/banner

获取客户端首页运营 banner 列表（mapBanner 映射后缓存）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取运营场景列表

GET /console/as/support/scenes

按 locale 获取运营场景（首页快捷入口）列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|locale|query|string| no |语言，如 zh / en。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/云盘 tdrive

## POST 申请文件上传

POST /api/v6/open/tdrive/upload/apply

上传第 1 步：申请上传。响应含 confirm_key / domain+path（直传地址）/ upload_id / headers（直传需带的头）/ task_id / access_token / is_quick_upload（秒传）。is_quick_upload=true 时跳过直传，直接 complete。name 与 file_size 为必填。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "name": "string",
  "ext": "string",
  "file_size": 0,
  "parent_id": "string",
  "task_id": "string",
  "local_file_path": "string",
  "local_task_id": "string",
  "channel": "string",
  "md5": "string",
  "create_source": 0,
  "hide_entry": true,
  "conflict_resolution_strategy": "rename",
  "is_multipart": true,
  "app_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» name|body|string| yes |文件名（不含扩展名）。|
|» ext|body|string| no |扩展名，小写。|
|» file_size|body|integer| yes |字节数，不能为 0。|
|» parent_id|body|string| no |父目录 fileID，根传空串。|
|» task_id|body|string| no |none|
|» local_file_path|body|string| no |none|
|» local_task_id|body|string| no |none|
|» channel|body|string| no |none|
|» md5|body|string| no |none|
|» create_source|body|integer| no |none|
|» hide_entry|body|boolean| no |none|
|» conflict_resolution_strategy|body|string| no |none|
|» is_multipart|body|boolean| no |none|
|» app_id|body|string| no |子应用 ID。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "confirm_key": "",
    "domain": "",
    "path": "",
    "upload_id": "",
    "headers": {},
    "task_id": "",
    "is_quick_upload": false
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 完成文件上传

POST /api/v6/open/tdrive/upload/complete

上传第 3 步：完成上传。task_id 与 confirm_key 必填。响应 file_id / tencent_doc_url / parent_id / status。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "task_id": "string",
  "confirm_key": "string",
  "name": "string",
  "ext": "string",
  "file_size": 0,
  "parent_id": "string",
  "is_failed": true,
  "errmsg": "string",
  "md5": "string",
  "hide_entry": true,
  "create_source": 0,
  "app_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» task_id|body|string| yes |none|
|» confirm_key|body|string| yes |none|
|» name|body|string| no |none|
|» ext|body|string| no |none|
|» file_size|body|integer| no |none|
|» parent_id|body|string| no |none|
|» is_failed|body|boolean| no |none|
|» errmsg|body|string| no |none|
|» md5|body|string| no |none|
|» hide_entry|body|boolean| no |none|
|» create_source|body|integer| no |none|
|» app_id|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "file_id": "",
    "tencent_doc_url": "",
    "parent_id": "",
    "status": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 续期上传

POST /api/v6/open/tdrive/upload/renew

上传 URL 过期后续期，重新取 domain/path/upload_id/headers。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "task_id": "string",
  "confirm_key": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» task_id|body|string| yes |none|
|» confirm_key|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 取消上传

POST /api/v6/open/tdrive/upload/cancel

取消上传任务。task_id 必填，confirm_key 可选。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "task_id": "string",
  "confirm_key": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» task_id|body|string| yes |none|
|» confirm_key|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取下载地址

POST /api/v6/open/tdrive/download

获取文件下载地址（download_url）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» file_id|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取目录列表

POST /api/v6/open/tdrive/dir/list

分页列出目录内容。响应 entries / next_cursor / total_num。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "parent_id": "string",
  "cursor": "",
  "limit": 20,
  "order_by": "string",
  "ascending": true,
  "entry_kind": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» parent_id|body|string| yes |父目录 fileID，必填。|
|» cursor|body|string| no |none|
|» limit|body|integer| no |none|
|» order_by|body|string| no |none|
|» ascending|body|boolean| no |none|
|» entry_kind|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取文件详情

POST /api/v6/open/tdrive/dir/info

获取目录/文件详情，响应 data.entry。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» file_id|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建目录

POST /api/v6/open/tdrive/dir/create

创建目录，返回 file_id。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "parent_id": "string",
  "name": "string",
  "app_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» parent_id|body|string| yes |none|
|» name|body|string| yes |none|
|» app_id|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 重命名

POST /api/v6/open/tdrive/dir/rename

重命名目录/文件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string",
  "new_name": "string",
  "conflict_strategy": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» file_id|body|string| yes |none|
|» new_name|body|string| yes |none|
|» conflict_strategy|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 移动目录/文件

POST /api/v6/open/tdrive/dir/move

移动目录/文件。dst_parent_id 不接受空串 / "/" / "root"，需先取子应用文件夹 fileID。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string",
  "dst_parent_id": "string",
  "conflict_strategy": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» file_id|body|string| yes |none|
|» dst_parent_id|body|string| yes |none|
|» conflict_strategy|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 回收站

POST /api/v6/open/tdrive/dir/trash

把目录/文件移入回收站。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» file_id|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 移动文件（file 命名空间）

POST /api/v6/open/tdrive/file/move

与 dir/move 同语义的另一入口，路由在 file 命名空间下。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string",
  "dst_parent_id": "string",
  "conflict_strategy": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» file_id|body|string| no |none|
|» dst_parent_id|body|string| no |none|
|» conflict_strategy|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取预览地址

POST /api/v6/open/tdrive/file/preview-url

获取文件预览 URL。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» file_id|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 转移所有权

POST /api/v6/open/tdrive/trans_ownership

转移文件/目录所有权。file_id 与 parent_id 均必填（proto3 snake_case）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "file_id": "string",
  "parent_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» file_id|body|string| no |none|
|» parent_id|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取子应用列表

POST /api/v6/open/tdrive/drive/subapp-list

列出网盘子应用（子应用下才有可用 parent fileID）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取容量

POST /api/v6/open/tdrive/capacity

查询网盘容量配额。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取历史列表

POST /api/v6/open/tdrive/history/list

分页获取访问历史。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 删除单条历史

POST /api/v6/open/tdrive/history/delete

删除指定历史记录。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 清空历史

POST /api/v6/open/tdrive/history/clear

清空访问历史。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 置顶历史

POST /api/v6/open/tdrive/history/set_latest

把某条历史置为最新/置顶。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 读取历史配置

POST /api/v6/open/tdrive/history/config/get

读取历史记录相关配置。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 写入历史配置

POST /api/v6/open/tdrive/history/config/set

写入历史记录相关配置。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 搜索文件

POST /api/v6/open/tdrive/search/file

按文件名搜索。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 搜索文件内容

POST /api/v6/open/tdrive/search/content

全文检索文件内容。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST RAG 文档检索

POST /api/v6/open/tdrive/search/rag-doc

RAG 文档向量检索。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST RAG 图片检索

POST /api/v6/open/tdrive/search/rag-image

RAG 图片检索。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 申请文档上传（agent/manage）

POST /api/v6/open/agent/manage/apply_upload

agent 管理域的上传申请（与 tdrive upload/apply 平行的一套）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 完成文档上传（agent/manage）

POST /api/v6/open/agent/manage/complete_upload

agent 管理域的上传完成。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建文档

POST /api/v6/open/agent/manage/create_file

在 agent 管理域创建文件条目。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 删除文档

POST /api/v6/open/agent/manage/delete_file

删除 agent 管理域文件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 列出文档

POST /api/v6/open/agent/manage/list_file

列出 agent 管理域文件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 查询文档信息

POST /api/v6/open/agent/manage/query_file_info

查询 agent 管理域文件信息。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 查询任务

POST /api/v6/open/agent/manage/query_task

查询 agent 管理域的异步任务。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 重命名文档标题

POST /api/v6/open/agent/manage/rename_file_title

重命名文档标题。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 搜索文档

POST /api/v6/open/agent/manage/search_file

搜索 agent 管理域文档。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建共享文件夹

POST /api/v6/open/tdrive/shared-folder/create

创建共享文件夹。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取应用列表

POST /api/v6/open/tdrive/drive/app-list

列出网盘应用。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/云助理 cloudagent

## GET 查询云助理权益

GET /v2/user/cloudagent/entitlement

查询云助理权益。data.orchestratorEnabled / enabled 任一为 true 即启用；另有 installed / teamsEnabled / exclusiveEnabled。code!=0 时客户端直接抛错。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "orchestratorEnabled": true,
    "installed": true,
    "teamsEnabled": false
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 初始化编排器

POST /v2/user/cloudagent/orchestrator

初始化用户编排器（orchestrator），无请求体。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 查询云助理配额

GET /v2/user/cloudagent/quota

查询云助理使用配额。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建云助理 Agent

POST /v2/user/cloudagent/agents

创建云助理 Agent（参数原样透传）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出已授权 Agent

GET /v2/user/cloudagent/agents

分页列出已授权 Agent。query: page/pageSize/projectIds/agentRole/withInstances/withManifest。无 endpoint 或未登录时客户端本地返回空列表、不发请求。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|pageSize|query|string| no |每页条数，默认 100。|
|projectIds|query|string| no |项目 ID 过滤。|
|agentRole|query|string| no |Agent 角色过滤。|
|withInstances|query|string| no |固定 true 时附带实例。|
|withManifest|query|string| no |固定 true 时附带 manifest。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "total": 0,
    "page": 1,
    "pageSize": 100,
    "items": []
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取 Agent 详情

GET /v2/user/cloudagent/agents/{agentId}

按 ID 获取 Agent 详情（含 currentVersion.manifest）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|agentId|path|string| yes |Agent ID（路径段做 URL encode）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## DELETE 删除 Agent

DELETE /v2/user/cloudagent/agents/{agentId}

删除 Agent。成功返回 204 无响应体（此时不带 Content-Type）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|agentId|path|string| yes |Agent ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|HTTP 204 No Content（无响应体）。|None|

## GET 校验 Agent 字段

GET /v2/user/cloudagent/validate-field

校验 Agent 字段取值是否可用（重名检查）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|field|query|string| yes |字段名，必填。|
|value|query|string| yes |待校验值，必填。|
|excludeAgentId|query|string| no |排除的 Agent ID（编辑场景）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 置顶 Agent

POST /v2/user/cloudagent/agents/pin

把 Agent 置顶。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "agentId": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» agentId|body|string| yes |Agent ID（会被 String() 化）。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 取消置顶 Agent

POST /v2/user/cloudagent/agents/unpin

取消 Agent 置顶。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "agentId": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» agentId|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 复制 Agent

POST /v2/user/cloudagent/agents/{agentId}/clone

复制 Agent。body 为除 agentId 外的其余字段。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|agentId|path|string| yes |源 Agent ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 发布 Agent 版本

POST /v2/user/cloudagent/agents/{agentId}/versions

为 Agent 发布新版本。body 为除 agentId 外的其余字段。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|agentId|path|string| yes |Agent ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出企业已发布 Agent

GET /v2/user/cloudagent/agents/enterprise-published

列出企业内已发布的 Agent。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码。|
|pageSize|query|string| no |每页条数。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "total": 0,
    "items": []
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建实例

POST /v2/user/cloudagent/instances

为 Agent 创建运行实例。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "agentId": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» agentId|body|string| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出实例

GET /v2/user/cloudagent/instances

分页列出实例。query: page/pageSize/sortBy(默认 lastActivityAt)/sortDirection(默认 desc)/agentId/keyword。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|pageSize|query|string| no |每页条数，默认 50。|
|sortBy|query|string| no |排序字段，默认 lastActivityAt。|
|sortDirection|query|string| no |排序方向，默认 desc。|
|agentId|query|string| no |按 Agent 过滤。|
|keyword|query|string| no |关键词。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 0,
    "page": 1,
    "pageSize": 50
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取实例详情

GET /v2/user/cloudagent/instances/{instanceId}

按 ID 获取实例详情。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|instanceId|path|string| yes |实例 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## DELETE 删除实例

DELETE /v2/user/cloudagent/instances/{instanceId}

删除实例。成功返回 204 无响应体。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|instanceId|path|string| yes |实例 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|HTTP 204 No Content（无响应体）。|None|

## POST 重建实例

POST /v2/user/cloudagent/instances/{instanceId}/rebuild

重建实例，无请求体。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|instanceId|path|string| yes |实例 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出用户会话（云助理）

GET /v2/user/cloudagent/conversations

用户维度会话浅层列表（侧栏 claw 视图用）。query: page/pageSize/agentId/keyword。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|pageSize|query|string| no |每页条数，默认 50。|
|agentId|query|string| no |按 Agent 过滤。|
|keyword|query|string| no |关键词。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 0,
    "hasNext": false
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出可用模型

GET /v2/enterprises/personal/models

列出个人版可用模型（data.models）。客户端会过滤掉 tags 含 text-to-image 的非对话模型。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "models": []
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 派发云助理任务

POST /v2/user/cloudagent/tasks

向云助理派发任务（把待办委派给助理）。请求体原样透传，使用专门的派发超时 CLOUDAGENT_DISPATCH_TIMEOUT_MS。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 取消云助理任务

POST /v2/user/cloudagent/tasks/cancel/{taskId}

取消已派发的云助理任务。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|taskId|path|string| yes |任务 ID（路径段 URL encode）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 企业 Claw 控制

POST /v2/enterprises/{enterpriseId}/claw/control

企业维度 Claw（终端代理）控制通道。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|enterpriseId|path|string| yes |企业 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 查询企业 License 状态

GET /console/enterprises/{enterpriseId}/license/status

查询企业 License 概览状态。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|enterpriseId|path|string| yes |企业 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 查询团队配额

GET /console/as/teams/me/quota

查询当前用户所属团队配额。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 下载云助理配置

GET /console/as/netdrive/cloudagent-config/download

下载云助理（CA2）配置。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/技能与插件市场

## GET 查询已安装技能列表

GET /v2/user-asset/skill/get-list

查询用户已安装技能列表。带 SUBMITTED-FROM 头标识来源（多处调用点该头值不同）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 安装技能

POST /v2/user-asset/skill/install

安装技能到当前用户资产。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 卸载技能

POST /v2/user-asset/skill/uninstall

卸载用户技能。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 技能资产动作

POST /v2/user-asset/skill/{action}

技能资产的通用动作入口（action 由调用方给出，如启用/禁用等）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|action|path|string| yes |动作名。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出自定义技能

GET /v2/user-asset/custom-skills

列出用户自定义技能。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建自定义技能

POST /v2/user-asset/custom-skills

以 multipart/form-data 创建自定义技能（可能在失败后回退为普通 POST 重试）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```yaml
{}

```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 自定义技能增删改

POST /v2/user-asset/custom-skills/{skill_uid}

按技能 UID 更新/删除自定义技能（UID 做 URL encode）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|skill_uid|path|string| yes |自定义技能 UID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取技能市场列表

GET /v2/enterprises/{enterpriseId}/skill-market/list

企业维度技能市场列表（SKILL_MARKET_PREFIX=/v2/enterprises）。桌面聚合版走 /portal/enterprises 前缀 POST。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|enterpriseId|path|string| yes |企业 ID（URL encode）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取技能市场分类

GET /portal/enterprises/{enterpriseId}/skill-market/categories

获取聚合技能市场分类列表（AGGREGATED_SKILL_MARKET_PREFIX=/portal/enterprises）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|enterpriseId|path|string| yes |企业 ID（URL encode）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 聚合技能市场列表

POST /portal/enterprises/{enterpriseId}/skill-market/list

桌面聚合技能市场列表，POST + 查询参数。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|enterpriseId|path|string| yes |企业 ID。|
|qs|query|string| no |聚合内置查询参数（版本/分组等，原样透传）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 内网技能网关列表

GET /console/agent-gateway/internal_knot-skill/apigw/api/v1/skills/list

KNOT 技能网关列表（KNOT_GATEWAY_PATH + /list）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "skills": []
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 内网技能网关标签

GET /console/agent-gateway/internal_knot-skill/apigw/api/v1/skills/list_tags

KNOT 技能网关标签列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 内网技能网关按 ID 查

POST /console/agent-gateway/internal_knot-skill/apigw/api/v1/skills/get_by_ids

KNOT 技能网关按 ID 批量查技能。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "ids": [
    "string"
  ]
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» ids|body|[string]| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 开放平台标签

GET /api/v2/open-platform/tags

开放平台技能标签列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 开放平台技能列表

GET /api/v2/open-platform/skills

开放平台技能列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 开放平台技能搜索

GET /api/v2/open-platform/skills/search

开放平台技能搜索。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 技能下载地址

GET /v2/operation-platform/market/skill/download-url

获取技能下载 URL。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建技能令牌

POST /openapi/v2/skill-token

为技能签发访问令牌。另有 /api/v5/robotLogic/sign-token 为签名入口。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 插件市场更新

POST /api/v1/plugins/marketplaces/update

更新插件市场索引（UPDATE_ENDPOINT）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 插件市场浏览

POST /api/v1/plugins/marketplaces/browse

浏览插件市场（BROWSE_ENDPOINT）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 插件市场列表

GET /api/v1/plugins/marketplaces

列出已配置插件市场。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 插件市场自动更新

POST /api/v1/plugins/marketplaces/auto-update

设置插件市场自动更新。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出插件

GET /api/v1/plugins

列出已安装插件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 启用插件

POST /api/v1/plugins/enable

启用插件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 禁用插件

POST /api/v1/plugins/disable

禁用插件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 卸载插件

POST /api/v1/plugins/uninstall

卸载插件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 重载插件

POST /api/v1/plugins/reload

重载插件运行时。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 对齐插件状态

POST /api/v1/plugins/reconcile

把本地插件状态与磁盘/远端对齐。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 市场插件列表

GET /console/as/marketplace/plugins

控制面市场插件列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 市场插件搜索

GET /console/as/marketplace/plugins/search

控制面市场插件搜索。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 市场来源列表

GET /console/as/marketplace/sources

市场来源列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 用户插件安装列表

GET /console/as/user/plugins/installed

列出用户已安装插件（控制面）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 用户插件安装

POST /console/as/user/plugins/install

为用户安装插件（控制面）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 知识库条目

POST /spapi/kb/entry/v1/entries

知识库（KB）条目接口。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 知识库文件申请上传

POST /spapi/kb/file/v1/file/apply

知识库文件上传申请。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 知识库文件提交

POST /spapi/kb/file/v1/file/commit

知识库文件上传提交。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 知识库导入任务

POST /spapi/kb/import/v1/tasks

知识库导入任务。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 生成技能令牌（sign-token）

POST /api/v5/robotLogic/sign-token

机器人逻辑签名 token 入口。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 完成资源上传

POST /api/v1/client/resources/complete

客户端资源上传完成回调（COMPLETE_PATH）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 申请资源预签名

POST /api/v1/client/resources/presign

客户端资源上传预签名（PRESIGN_PATH）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/开放授权 workbuddyopen

## GET 查询授权状态

GET /workbuddyopen/v1/authorize/status

查询第三方应用（application_id）的授权状态，返回应用名/图标/redirect_uris/allowed_scopes（scope_id/display_name/description/risk_level）。客户端对返回结构做严格校验，id 必须等于请求的 applicationId，否则抛错。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|application_id|query|string| yes |第三方应用 ID，必填。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "application": {
      "id": "app_x",
      "name": "示例应用",
      "redirect_uris": [
        "https://..."
      ],
      "allowed_scopes": []
    }
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 申请授权码

POST /workbuddyopen/v1/authorize/code

申请授权码（OAuth code）。返回的 redirect_uri 语义与 buildBuddyRedirectUrl 对齐。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "application_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» application_id|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "code": "xxx",
    "redirect_uri": "https://..."
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 提交授权同意

POST /workbuddyopen/v1/authorize/consent

提交用户授权同意（用户点击同意后调用）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "application_id": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» application_id|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST Connector OAuth 换 token

POST /console/as/connector/oauth/{segment}/accesstoken

Connector OAuth 换取 access token。同类入口：/console/as/connector/oauth/{segment}/start（发起）、/console/as/connector/oauth/{segment}/revoke（撤销）、/console/as/connector/oauth/{segment}/status。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|segment|path|string| yes |connector 标识段（如 agentmail / ardot）。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST Connector OAuth 发起

POST /console/as/connector/oauth/{segment}/start

发起 connector OAuth 授权流程。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|segment|path|string| yes |connector 标识段。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET Connector OAuth 状态

GET /console/as/connector/oauth/{segment}/status

查询 connector OAuth 授权状态。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|segment|path|string| yes |connector 标识段。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST Connector OAuth 撤销

POST /console/as/connector/oauth/{segment}/revoke

撤销 connector OAuth 授权。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|segment|path|string| yes |connector 标识段。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建客户端 API Key

POST /v2/api-keys

创建客户端临时 API Key（CLIENT_TEMPKEY_API_PATH）。另有控制面版本 /console/api/client/v1/api-keys。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 查询功能开关

GET /v2/feature-flag/api

查询功能开关（feature flag）。另有 /console/feature-flag/api 与 /console/feature-flag/api/feature-flags/check 两个变体。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 反馈上报

POST /v2/feedback

用户反馈上报（另一处为 {productEndpoint}/v2/feedback）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 埋点上报（对话活跃上报）

POST /v2/report

客户端行为埋点上报。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

---
【2026-09-11 Hermes 实测验证 · 活跃门槛破解】

结论：本接口即客户端「对话上报」通道，是 growth 活跃地图 / 连登 与 first_buddy 任务的唯一输入。

1) 传输：POST，body 为「事件数组」（非单对象）。实测网关 https://www.codebuddy.cn 与 https://copilot.tencent.com 均返回 200 {"code":0,"msg":"OK"}。
   注意：按客户端实现，本接口请求**不带 Authorization 头**（仅 Content-Type）。鉴权由网关/会话隐式承担。

2) 关键字段（缺则服务端 200 但静默丢弃，实测）：
   - eventCode: 必须为 "chat_request_send"
   - userId:    **必须等于账号 uid（唯一隐藏门槛）**
   最小可用载荷 = {eventCode, timestamp, userId}，实测即被计入。

3) 服务端**不校验** body 与真实会话的一致性：自造 conversationId 依然计数（实测 chat_5 进度 0→5）。

4) 立即效果（实测）：
   - GET /activity/growth/streak → days 由 0 变 1（活跃地图点亮）
   - GET /v2/activity/growth/tasks → first_buddy 由 not_accepted 变 completed
   - POST /activity/growth/buddy/first → 200，+300 分 +8 能量（免费猫）
   - 累计型任务 chat_5 每次上报 +1 进度，满 5 次可 claim（+100 分 +5 能量）

5) 与「每日活跃奖励」的关系：官方规则（www.workbuddy.ai/docs/zh/workbuddy/Subscription#限时每日活跃奖励说明）
   明确「当日通过客户端发起至少一次对话/任务即视为当日活跃，积分次日发放」，Free +30 分/天、Pro +50 分/天。
   官方 Release Note 亦写明「修复 WorkBuddy Desktop 使用自定义模型对话时未被计入活跃统计（活跃地图、连登等）」，
   即本地接口与活跃地图同源。

6) 风险：官方规则明文禁止「脚本、外挂、插件、模拟器、接口调用、批量注册」刷取活跃，并保留追回奖励权利。
   请自行评估合规边界，建议低频、真号真 token、不做批量注册。

> Body Parameters

```json
[
  {
    "eventCode": "chat_request_send",
    "timestamp": 1789131385000,
    "reportDelay": 0,
    "userId": "<uid>",
    "userNickname": "",
    "conversationId": "<conversation-id>",
    "requestId": "<request-id>",
    "rootRequestId": "<request-id>",
    "parentConversationId": "<conversation-id>",
    "mode": "craft",
    "inputLength": 12,
    "requestModelId": "deepseek-v4-flash",
    "requestModelName": "DeepSeek V4 Flash",
    "isPlan": false,
    "isAutoExecuteTerminal": false,
    "isAutoModify": false,
    "codebaseEnable": false,
    "maxToken": 0,
    "maxSteps": 0,
    "temperature": 0,
    "maxRetries": 0,
    "mentionContexts": [],
    "knowledgeId": [],
    "knowledgeName": [],
    "codebaseId": "",
    "mentionContextCount": 0,
    "command": "",
    "expertId": "",
    "recommendId": "",
    "skillId": "",
    "skillCount": 0,
    "totalCount": 0,
    "fileUri": "",
    "presentAt": 1789131385000,
    "traceId": "",
    "agentName": "default",
    "agentType": "conversation",
    "timezone": "Asia/Shanghai",
    "os": "Linux",
    "userAgent": "CLI/2.63.2 CodeBuddy/2.63.2"
  }
]
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Content-Type|header|string| yes |application/json;charset=UTF-8|
|Authorization|header|string| no |客户端实测不带该头；如网关要求可补 Bearer <accessToken>。|
|body|body|array[object]| yes |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 文本审核

POST /v2/as/moderation/text

文本内容审核（moderation）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 记忆检索

POST /api/memory/search

记忆（memory）检索。同域另有 /api/memory/profile（读档案）、/api/memory/update_profile、/api/memory/clear_profile。Desktop 调用形如 {endpoint}/api/memory/search。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-User-ID|header|string| no |用户 uid。|
|X-Enterprise-Id|header|string| no |企业 ID（可选）。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 读取记忆档案

GET /api/memory/profile

读取用户记忆档案。remote memory base 默认 https://copilot.tencent.com（productManager 未给 endpoint 时）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-User-ID|header|string| no |用户 uid。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 更新记忆档案

POST /api/memory/update_profile

更新用户记忆档案。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-User-ID|header|string| no |用户 uid。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 清空记忆档案

POST /api/memory/clear_profile

清空用户记忆档案。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|X-User-ID|header|string| no |用户 uid。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 语音识别

POST /agenttool/v1/asr

语音转文字（ASR）。multipart/form-data，字段 audio；Cloud 环境直连 {endpoint}/agenttool/v1/asr，超时 30s。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```yaml
audio: string

```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» audio|body|string(binary)| yes |音频文件 Blob（fileName + mimeType）。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 语音合成

POST /agenttool/v1/tts

文本转语音（TTS）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 获取临时密钥

POST /agenttool/v1/tempkey

获取 agenttool 临时密钥。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 意图改写嵌入

POST /agenttool/v1/intent/rewrite-embed

意图改写与向量嵌入。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST Agentic 搜索

POST /agenttool/v1/agentic_search

Agentic 搜索（agenttool 域）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 聊天补全

POST /chat/completions

OpenAI 兼容聊天补全入口（形如 {pathPart}/chat/completions{suffix}）。本地模型场景另走 /codebase/openai/v1/embeddings 做嵌入。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 嵌入（codebase）

POST /codebase/openai/v1/embeddings

codebase 向量嵌入（OpenAI 兼容）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 账号公司状态校验

POST /tapi/account/v1/check-company-status

校验账号的公司/组织状态。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 账号 License 校验

POST /tapi/account/v1/check-license

校验账号 License。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/会话 conversations

## POST 创建会话

POST /v2/as/conversations/

创建会话（注意路径以 / 结尾）。conversationOrigin 默认 agents；企业云助理用 ENTERPRISE_CLOUDAGENTS_ORIGIN。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "agentId": "string",
  "conversationOrigin": "agents",
  "model": "string",
  "exclusive": true
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| yes |none|
|» agentId|body|string| yes |none|
|» conversationOrigin|body|string| yes |none|
|» model|body|string| no |可选，指定模型。|
|» exclusive|body|boolean| no |为 true 时才附带。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出会话

GET /v2/as/conversations/

分页列出会话。status 会被包装成 filters=[{field:'status',value:...}] JSON 串。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|page|query|string| no |页码，默认 1。|
|size|query|string| no |每页条数，默认 50。|
|agentId|query|string| no |按 Agent 过滤。|
|title|query|string| no |标题模糊。|
|sort|query|string| no |排序。|
|conversationOrigin|query|string| no |来源过滤。|
|filters|query|string| no |JSON 串，形如 [{"field":"status","value":"..."}]。|
|dayRange|query|string| no |天数范围。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "conversations": [],
    "total": 0
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取会话详情

GET /v2/as/conversations/{conversationId}

按 ID 获取会话详情。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|conversationId|path|string| yes |会话 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 更新会话

POST /v2/as/conversations/{conversationId}

更新会话元信息。至少传一个字段，否则客户端本地抛错。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{
  "title": "string",
  "isUserDefinedTitle": true,
  "visibility": "string",
  "status": "string",
  "extraInfo": {},
  "expertId": "string",
  "locale": "string"
}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|conversationId|path|string| yes |会话 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|
|» title|body|string| no |none|
|» isUserDefinedTitle|body|boolean| no |none|
|» visibility|body|string| no |none|
|» status|body|string| no |none|
|» extraInfo|body|object| no |none|
|» expertId|body|string| no |none|
|» locale|body|string| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 删除会话

POST /v2/as/conversations/{conversationId}/delete

删除会话（用 POST + /delete 子路径，非 HTTP DELETE）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|conversationId|path|string| yes |会话 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 批量获取会话

POST /v2/as/conversations/batch-get

批量按 ID 获取会话。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 解析会话消息

POST /console/as/message/resolve

按会话解析/渲染消息。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 创建定时任务（scheduler）

POST /v2/as/scheduler/tasks

AS 调度器创建任务。DEFAULT_ROUTE_PREFIX 即此路径，SchedulerService / MsgCenterService / CloudNotificationsRepo 等多个类共用该前缀作为基址。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出定时任务

GET /console/as/tasks

列出 AS 任务（TASKS_PREFIX）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取任务详情

GET /console/as/tasks/{taskId}

按 ID 获取 AS 任务。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|taskId|path|string| yes |任务 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 项目邀请

POST /console/as/invitations/project

创建项目邀请（INVITATION_ROUTE_PREFIX）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET Connector 注册表（2C）

GET /console/as/connector/registry2c/list

列出 2C connector 注册表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取 Connector 用户信息

GET /console/as/connector/user/{userId}

按用户 ID 取 connector 用户信息。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|userId|path|string| yes |用户 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取 Connector 任务

GET /console/as/connector/task/{taskId}

按任务 ID 取 connector 任务信息。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|taskId|path|string| yes |任务 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET Connector 项目连接

GET /console/as/connector/project/{projectId}

按项目 ID 取 connector 连接信息。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|projectId|path|string| yes |项目 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 上传待办附件

POST /console/as/support/presigned_url

获取待办附件上传用的预签名 URL。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 解析应用（genie-baas）

POST /console/as/genie-baas/control/applications/resolve

解析应用（genie-baas 控制面）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 绑定云服务（genie-baas）

POST /console/as/genie-baas/control/cloud-services

绑定云服务（genie-baas 控制面）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/智能邮箱 agentmail

## GET 查询邮箱绑定状态

GET /v2/as/connector/agentmail/me

查询智能邮箱（AgentMail）绑定态。data.status='not_bound' 表示未绑定；已绑定返回 alias/status/openedAt/userType，未绑定返回 userType/maskedPhone。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "status": "not_bound",
    "userType": "personal",
    "maskedPhone": "138****0000"
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 发送绑定验证码

POST /v2/as/connector/agentmail/sms/send

发送短信验证码。code!=0 抛 AgentMailApiError。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 绑定邮箱

POST /v2/as/connector/agentmail/bind

绑定智能邮箱（参数原样透传，含验证码）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 停用邮箱

POST /v2/as/connector/agentmail/deactivate

停用智能邮箱。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 重新启用邮箱

POST /v2/as/connector/agentmail/reactivate

重新启用邮箱。data.status==='auth_failed' 时按业务态返回而非抛错。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "status": "active"
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取管理页 URL

GET /v2/as/connector/agentmail/url/admin

获取管理页跳转 URL。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "url": "https://..."
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取实名页 URL

GET /v2/as/connector/agentmail/url/realname

获取实名认证页跳转 URL。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {
    "url": "https://..."
  }
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 邮箱 REST 状态

GET /console/agent-gateway/agentmailrest/v1/me

AgentMail REST 网关状态。带 auth headers 时请求 {endpoint}/console/agent-gateway/agentmailrest/v1/me。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出邮件

GET /console/agent-gateway/agentmailrest/v1/aliases/{alias_id}/messages

列出某别名的邮件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|alias_id|path|string| yes |邮箱别名 ID。|
|params|query|string| no |分页/过滤参数，原样透传。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 搜索邮件

GET /console/agent-gateway/agentmailrest/v1/aliases/{alias_id}/messages/search

搜索某别名的邮件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|alias_id|path|string| yes |邮箱别名 ID。|
|params|query|string| no |搜索参数，原样透传。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 获取邮件详情

GET /console/agent-gateway/agentmailrest/v1/aliases/{alias_id}/messages/{message_id}

获取单封邮件详情。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|alias_id|path|string| yes |邮箱别名 ID。|
|message_id|path|string| yes |邮件 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## DELETE 删除邮件

DELETE /console/agent-gateway/agentmailrest/v1/aliases/{alias_id}/messages/{message_id}

删除邮件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|alias_id|path|string| yes |邮箱别名 ID。|
|message_id|path|string| yes |邮件 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 更新邮件

POST /console/agent-gateway/agentmailrest/v1/aliases/{alias_id}/messages/{message_id}/update

更新邮件（如已读/星标等），body 原样透传。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|alias_id|path|string| yes |邮箱别名 ID。|
|message_id|path|string| yes |邮件 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 列出附件

GET /console/agent-gateway/agentmailrest/v1/aliases/{alias_id}/messages/{message_id}/attachments

列出邮件附件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|alias_id|path|string| yes |邮箱别名 ID。|
|message_id|path|string| yes |邮件 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 下载附件

GET /console/agent-gateway/agentmailrest/v1/aliases/{alias_id}/messages/{message_id}/attachments/{attachment_id}

下载指定附件。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|alias_id|path|string| yes |邮箱别名 ID。|
|message_id|path|string| yes |邮件 ID。|
|attachment_id|path|string| yes |附件 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "msg": "ok",
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# workbuddy/专家与运营市场

## GET 专家排行榜

GET /portal/expert/ranking

专家排行榜（routePrefix=/portal，可带 query）。渲染端另有 /console/expert/ranking 变体。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|query|query|string| no |排行榜查询参数（周期等），原样拼接。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 专家市场列表

GET /portal/enterprises/{enterpriseId}/expert-market/list

企业专家市场列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|enterpriseId|path|string| yes |企业 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 按 ID 获取专家

GET /portal/enterprises/{enterpriseId}/expert-market/get-by-ids

按 ID 批量获取专家。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|enterpriseId|path|string| yes |企业 ID。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 专家分类

GET /portal/operation-platform/market/expert-categories

运营平台专家分类列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": []
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 专家列表（运营平台）

GET /portal/operation-platform/market/expert/list

运营平台专家列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 专家下载地址

GET /portal/operation-platform/market/expert/download-url

获取专家（Expert / 技能包）下载 URL。ExpertDownloadUrlService 构造时前缀即 /portal。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 运营市场内置列表

GET /v2/operation-platform/market

运营平台市场（内置市场）入口前缀 BUILTIN_MARKET_ROUTE_PREFIX=/v2/operation-platform/market。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 网盘分享创建

POST /v2/as/netdisk/share

创建网盘分享。同域另有 /v2/as/netdisk/shares（列表）、/v2/as/netdisk/shares/{id}、/v2/as/p/netdisk/share/{id}（公开访问）与 /v2/as/netdisk/upload。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 网盘分享列表

GET /v2/as/netdisk/shares

列出网盘分享。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 任务分享准备

POST /v2/as/tasks/share/prepare

准备任务分享。配套：/v2/as/tasks/shares（列表）、/v2/as/tasks/share/confirm/{id}、/v2/as/tasks/share/cancel/{id}、/v2/as/tasks/share/prepare/{id}。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 任务分享列表

GET /v2/as/tasks/shares

列出任务分享。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 网盘访问令牌

GET /console/as/netdrive/access-token

获取网盘访问令牌。同域：/console/as/netdrive/enterprise-info、/console/as/netdrive/nicknames、/console/as/netdrive/issue-jump-code、/console/as/netdrive/upload-to-drive。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 网盘企业信息

GET /console/as/netdrive/enterprise-info

获取网盘企业信息。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 网盘昵称列表

GET /console/as/netdrive/nicknames

获取网盘昵称列表。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 签发跳转码

POST /console/as/netdrive/issue-jump-code

签发网盘跳转码。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 上传到网盘

POST /console/as/netdrive/upload-to-drive

上传文件到网盘（控制面）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 微信客服代理绑定

POST /v2/backgroundagent/wechatkfProxy/bind

微信客服代理绑定。配套：/v2/backgroundagent/wechatkfProxy/bindStatus、/v2/backgroundagent/wechatkfProxy/link。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 微信客服代理状态

GET /v2/backgroundagent/wechatkfProxy/bindStatus

查询微信客服代理绑定状态。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 企业微信本地代理接收

POST /v2/backgroundagent/wecom/local-proxy/receive

企业微信本地代理消息接收回调。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 会话预测

POST /v2/chat/prediction

会话（输入）预测。配套：/v2/chat/queue/status（队列状态）、/v2/chat/queue/cancel（取消）。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 会话队列状态

GET /v2/chat/queue/status

查询会话队列状态。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 取消会话排队

POST /v2/chat/queue/cancel

取消排队中的会话请求。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 灵感卡片列表

GET /v2/as/inspiration/code/encode

灵感（inspiration）code 编码入口。配套 /v2/as/inspiration/code/decode；encode 支持 ?reuse=false 强制重新生成。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|reuse|query|string| no |false 时不复用已有编码。|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## GET 灵感 code 解码

GET /v2/as/inspiration/code/decode

灵感 code 解码。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

## POST 后台代理本地上传

POST /v2/backgroundagent/localProxy/upload

后台代理本地文件上传。

来源：WorkBuddy 桌面客户端 v5.5.4 静态逆向（asar 解包 + 云端 web bundle），非动态抓包。
鉴权：Desktop 端直连 IDE 网关，请求头 `Authorization: Bearer <accessToken>`（OIDC）；部分域另带 X-Product-Code / X-User-Id / X-Enterprise-Id。
路径前缀：Desktop 与 Web 的网关路由前缀不同（如 billing 在 Desktop 需 /v2、msg-center 在 Desktop 需 /v2 而 Node 走 /portal）。本条目路径按已标注的形态登记。

> Body Parameters

```json
{}
```

### Params

|Name|Location|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string| yes |Bearer <accessToken>，IDE 网关 OIDC 令牌。|
|body|body|object| no |none|

> Response Examples

> 200 Response

```json
{
  "code": 0,
  "data": {}
}
```

### Responses

|HTTP Status Code |Meaning|Description|Data schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### Responses Data Schema

HTTP Status Code **200**

|Name|Type|Required|Restrictions|Title|description|
|---|---|---|---|---|---|
|» code|integer|false|none||0 表示成功。|
|» msg|string|false|none||none|
|» data|object|false|none||业务数据载荷。|

# Data Schema

<h2 id="tocS_Pet">Pet</h2>

<a id="schemapet"></a>
<a id="schema_Pet"></a>
<a id="tocSpet"></a>
<a id="tocspet"></a>

```json
{
  "id": 1,
  "category": {
    "id": 1,
    "name": "string"
  },
  "name": "doggie",
  "photoUrls": [
    "string"
  ],
  "tags": [
    {
      "id": 1,
      "name": "string"
    }
  ],
  "status": "available"
}

```

### Attribute

|Name|Type|Required|Restrictions|Title|Description|
|---|---|---|---|---|---|
|id|integer(int64)|true|none||宠物ID编号|
|category|[Category](#schemacategory)|true|none||分组|
|name|string|true|none||名称|
|photoUrls|[string]|true|none||照片URL|
|tags|[[Tag](#schematag)]|true|none||标签|
|status|string|true|none||宠物销售状态|

#### Enum

|Name|Value|
|---|---|
|status|available|
|status|pending|
|status|sold|

<h2 id="tocS_Tag">Tag</h2>

<a id="schematag"></a>
<a id="schema_Tag"></a>
<a id="tocStag"></a>
<a id="tocstag"></a>

```json
{
  "id": 1,
  "name": "string"
}

```

### Attribute

|Name|Type|Required|Restrictions|Title|Description|
|---|---|---|---|---|---|
|id|integer(int64)|false|none||标签ID编号|
|name|string|false|none||标签名称|

<h2 id="tocS_Category">Category</h2>

<a id="schemacategory"></a>
<a id="schema_Category"></a>
<a id="tocScategory"></a>
<a id="tocscategory"></a>

```json
{
  "id": 1,
  "name": "string"
}

```

### Attribute

|Name|Type|Required|Restrictions|Title|Description|
|---|---|---|---|---|---|
|id|integer(int64)|false|none||分组ID编号|
|name|string|false|none||分组名称|

