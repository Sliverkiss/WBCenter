// API 契约类型层：与 docs/api-contract.md 逐字段对齐，字段名照抄后端 Go JSON tag（snake_case 不转 camelCase）。
// 时间字段为 RFC 3339 字符串。带 //go 注释的 interface 标明 Go 侧来源结构体。
export * from './page'

/** 过渡宽松类型：未完成逐页迁移的页面暂用，M2+ 逐页替换为契约类型。 */
export type Any = Record<string, unknown>

// ---- 通用 ----

/** 错误信封：所有非 2xx 响应统一为 {"error": "中文错误"}。 */
export interface ErrorEnvelope {
  error: string
}

// ---- 会话与登录 ----

/** GET /api/session */
export interface SessionInfo {
  authenticated: boolean
  read_only: boolean
  using_default_password: boolean
  timezone: string
}

/** POST /api/login 与 /api/logout 的成功响应。 */
export interface OkResponse {
  ok: boolean
}

/** 动作类端点的成功响应（automationRun / accountAction）。 */
export interface ActionResult {
  ok: boolean
  message: string
}

// ---- 概览 ----

/** GET /api/overview */
export interface OverviewData {
  accounts: number
  healthy: number
  expired: number
  automations: number
  warnings: string[]
  updated_at: string
}

// ---- 账号 ----

/** go: control.Account */
export interface Account {
  uid: string
  nickname: string
  domain: string
  /** Token 过期时间（Unix 秒）。 */
  expires_at: number
  expired: boolean
  needs_refresh: boolean
}

/** GET /api/accounts */
export interface AccountsResponse {
  accounts: Account[]
  warnings: string[]
}

// ---- 积分 ----

/** go: control.CreditSummary；error 仅在部分/全部失败时出现。 */
export interface CreditSummary {
  uid: string
  nickname: string
  current: number
  today_allocated: number
  today_consumed: number
  today_remaining: number
  packages: number
  fetched_at: string
  error?: string
}

/** GET /api/credits 与 /api/credits/{uid} */
export interface CreditsResponse {
  items: CreditSummary[]
}

// ---- 模型 ----

/** go: control.Model */
export interface ModelItem {
  uid: string
  nickname: string
  id: string
  name: string
  raw?: Record<string, unknown>
  error?: string
}

/** GET /api/models */
export interface ModelsResponse {
  items: ModelItem[]
}

// ---- 活动（成长任务） ----

/** go: control.Activity */
export interface ActivityItem {
  uid: string
  nickname: string
  code: string
  name: string
  status: string
  reward: number
  new: boolean
}

/** GET /api/activities */
export interface ActivitiesResponse {
  items: ActivityItem[]
}

// ---- 云端定时任务（上游只读） ----

/** go: control.SchedulerTask */
export interface SchedulerTask {
  uid: string
  nickname: string
  id: string
  name: string
  status: string
  updated?: string
  raw?: Record<string, unknown>
  error?: string
}

/** GET /api/scheduler-tasks */
export interface SchedulerTasksResponse {
  items: SchedulerTask[]
  write_enabled: boolean
  note: string
}

/** GET /api/scheduler-tasks/{uid}/{taskID}（脱敏透传，形状由上游决定）。 */
export interface SchedulerTaskDetailResponse {
  item: Record<string, unknown>
}

// ---- 本地 Mock 定时任务 ----

/** go: control.MockSchedulerTaskView（内嵌 MockSchedulerTask + nickname）。 */
export interface MockSchedulerTask {
  id: string
  account_uid: string
  name: string
  cron: string
  prompt: string
  enabled: boolean
  created_at: string
  updated_at: string
  nickname: string
}

/** POST/PUT /api/mock-scheduler-tasks 请求体。 */
export interface MockSchedulerTaskInput {
  account_uid: string
  name: string
  cron: string
  prompt: string
  enabled: boolean
}

/** GET /api/mock-scheduler-tasks */
export interface MockSchedulerTasksResponse {
  items: MockSchedulerTask[]
  mode: 'mock'
  note: string
}

// ---- 本地自动化 ----

/** go: control.Automation；last_run_at/last_result/next_run_at 未运行或未启用时缺省。 */
export interface Automation {
  id: string
  name: string
  action: string
  enabled: boolean
  every_minutes: number
  last_run_at?: string
  last_result?: string
  next_run_at?: string
}

/** go: control.RunRecord */
export interface RunRecord {
  at: string
  action: string
  ok: boolean
  message: string
}

/** GET /api/automations */
export interface AutomationsResponse {
  items: Automation[]
  runs: RunRecord[]
}

/** PUT /api/automations/{id} 请求体。 */
export interface AutomationInput {
  enabled: boolean
  every_minutes: number
}

// ---- OAuth ----

/** POST /api/oauth/start 请求体。 */
export interface OAuthStartInput {
  region: 'cn' | 'global'
}

/** POST /api/oauth/start 响应。 */
export interface OAuthStartResponse {
  id: string
  url: string
}

/** POST /api/oauth/{id}/poll 等待中响应。 */
export interface OAuthPollPending {
  status: 'pending'
}

/** POST /api/oauth/{id}/poll 成功响应。 */
export interface OAuthPollSuccess {
  status: 'success'
  uid: string
  nickname: string
}

export type OAuthPollResponse = OAuthPollPending | OAuthPollSuccess

// ---- 新增端点（M2+，后端实现随后续里程碑补齐） ----

/** POST /api/probe 单账号探测结果。 */
export interface ProbeResult {
  uid: string
  nickname: string
  ok: boolean
  checkin?: {
    checked_in: boolean
    streak_days: number
  }
  credits?: {
    current: number
    today_remaining: number
  }
  travel?: {
    status: 'idle' | 'traveling' | 'arrived'
    location?: string
    departed_at?: string
    arrive_at?: string
  }
  error?: string
}

/** POST /api/probe 响应。 */
export interface ProbeResponse {
  results: ProbeResult[]
}

/** GET /api/stats/summary */
export interface StatsSummary {
  accounts_total: number
  credits_current_total: number
  credits_today_allocated_total: number
  credits_today_consumed_total: number
  checkin_done_today: number
  checkin_pending_today: number
  generated_at: string
}

/** GET /api/logs 单行。 */
export interface LogLine {
  at: string
  level: 'info' | 'warn' | 'error'
  message: string
}

/** GET /api/logs */
export interface LogsResponse {
  lines: LogLine[]
}

/** GET /api/models/pricing 条目。⚠️ 上游字段未核实，M4 实测定稿。 */
export interface ModelPricing {
  id: string
  name: string
  free: boolean
  price: number
  price_unit: string
  free_quota: number
  note: string
}

/** GET /api/models/pricing；上游无价格数据时 items 为空并带 warning。 */
export interface ModelPricingResponse {
  items: ModelPricing[]
  warning?: string
}

// ---- 类型守卫（zod-free 手写窄化） ----

const isRecord = (x: unknown): x is Record<string, unknown> => typeof x === 'object' && x !== null && !Array.isArray(x)
const isStr = (x: unknown): x is string => typeof x === 'string'
const isNum = (x: unknown): x is number => typeof x === 'number'
const isBool = (x: unknown): x is boolean => typeof x === 'boolean'
const isOptStr = (x: unknown): boolean => x === undefined || isStr(x)

export function isAccount(x: unknown): x is Account {
  return (
    isRecord(x) &&
    isStr(x.uid) &&
    isStr(x.nickname) &&
    isStr(x.domain) &&
    isNum(x.expires_at) &&
    isBool(x.expired) &&
    isBool(x.needs_refresh)
  )
}

export function isCreditSummary(x: unknown): x is CreditSummary {
  return (
    isRecord(x) &&
    isStr(x.uid) &&
    isStr(x.nickname) &&
    isNum(x.current) &&
    isNum(x.today_allocated) &&
    isNum(x.today_consumed) &&
    isNum(x.today_remaining) &&
    isNum(x.packages) &&
    isStr(x.fetched_at) &&
    isOptStr(x.error)
  )
}

export function isProbeResult(x: unknown): x is ProbeResult {
  if (!isRecord(x) || !isStr(x.uid) || !isStr(x.nickname) || !isBool(x.ok) || !isOptStr(x.error)) return false
  if (x.checkin !== undefined) {
    const c = x.checkin
    if (!isRecord(c) || !isBool(c.checked_in) || !isNum(c.streak_days)) return false
  }
  if (x.credits !== undefined) {
    const c = x.credits
    if (!isRecord(c) || !isNum(c.current) || !isNum(c.today_remaining)) return false
  }
  if (x.travel !== undefined) {
    const t = x.travel
    if (!isRecord(t) || !isStr(t.status) || !isOptStr(t.location) || !isOptStr(t.departed_at) || !isOptStr(t.arrive_at)) return false
  }
  return true
}

export function isSessionInfo(x: unknown): x is SessionInfo {
  return (
    isRecord(x) &&
    isBool(x.authenticated) &&
    isBool(x.read_only) &&
    isBool(x.using_default_password) &&
    isStr(x.timezone)
  )
}

export function isOverviewData(x: unknown): x is OverviewData {
  return (
    isRecord(x) &&
    isNum(x.accounts) &&
    isNum(x.healthy) &&
    isNum(x.expired) &&
    isNum(x.automations) &&
    Array.isArray(x.warnings) &&
    (x.warnings as unknown[]).every(isStr) &&
    isStr(x.updated_at)
  )
}

export function isAutomation(x: unknown): x is Automation {
  return (
    isRecord(x) &&
    isStr(x.id) &&
    isStr(x.name) &&
    isStr(x.action) &&
    isBool(x.enabled) &&
    isNum(x.every_minutes) &&
    isOptStr(x.last_run_at) &&
    isOptStr(x.last_result) &&
    isOptStr(x.next_run_at)
  )
}
