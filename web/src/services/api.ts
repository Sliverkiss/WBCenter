// fetch 封装：同源凭据、JSON 体、错误信封解析、401 统一登出流转。
// 域方法与 docs/api-contract.md 逐端点对照，返回值直接用契约类型。
import type {
  Account,
  AccountsResponse,
  ActionResult,
  ActivitiesResponse,
  Automation,
  AutomationInput,
  AutomationsResponse,
  CreditsResponse,
  ErrorEnvelope,
  LogsResponse,
  MockSchedulerTask,
  MockSchedulerTaskInput,
  MockSchedulerTasksResponse,
  ModelPricingResponse,
  ModelsResponse,
  OAuthPollResponse,
  OAuthStartInput,
  OAuthStartResponse,
  OkResponse,
  OverviewData,
  ProbeResponse,
  SchedulerTaskDetailResponse,
  SchedulerTasksResponse,
  SessionInfo,
  StatsSummary,
} from '../types'

// 401 全局登出回调：由 App 注入（setUnauthorizedHandler），触发后跳转登录 hash。
// 模块级注入避免引入状态库。
let onUnauthorized: (() => void) | null = null
export function setUnauthorizedHandler(cb: (() => void) | null) {
  onUnauthorized = cb
}

function handleUnauthorized(url: string) {
  // 登录/会话端点自身的 401 属于业务失败（密码错误等），不触发全局登出。
  if (url.startsWith('/api/login') || url.startsWith('/api/session')) return
  window.location.hash = '#/'
  onUnauthorized?.()
}

export async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    credentials: 'same-origin',
    headers: init?.body ? { 'Content-Type': 'application/json' } : undefined,
    ...init,
  })
  const body = (await response.json().catch(() => ({}))) as Partial<ErrorEnvelope>
  if (!response.ok) {
    if (response.status === 401) handleUnauthorized(url)
    throw new Error(String(body.error || '请求失败'))
  }
  return body as T
}

export const get = <T,>(url: string) => request<T>(url)
export const post = <T,>(url: string, body: unknown = {}) =>
  request<T>(url, { method: 'POST', body: JSON.stringify(body) })
const put = <T,>(url: string, body: unknown) => request<T>(url, { method: 'PUT', body: JSON.stringify(body) })
const del = <T,>(url: string) => request<T>(url, { method: 'DELETE' })

const enc = encodeURIComponent

// ---- 会话 ----
export const getSession = () => get<SessionInfo>('/api/session')
export const login = (input: { username: string; password: string }) => post<OkResponse>('/api/login', input)
export const logout = () => post<OkResponse>('/api/logout')

// ---- 概览与账号 ----
export const getOverview = () => get<OverviewData>('/api/overview')
export const getAccounts = () => get<AccountsResponse>('/api/accounts')
export const accountAction = (uid: string, action: 'checkin' | 'travel' | 'refresh') =>
  post<ActionResult>(`/api/accounts/${enc(uid)}/actions/${action}`)

// ---- 积分 ----
export const getCredits = () => get<CreditsResponse>('/api/credits')
export const getCreditDetail = (uid: string) => get<CreditsResponse>(`/api/credits/${enc(uid)}`)

// ---- 模型与活动 ----
export const getModels = () => get<ModelsResponse>('/api/models')
export const getActivities = () => get<ActivitiesResponse>('/api/activities')

// ---- 云端定时任务（上游只读） ----
export const getSchedulerTasks = () => get<SchedulerTasksResponse>('/api/scheduler-tasks')
export const getSchedulerTaskDetail = (uid: string, taskID: string) =>
  get<SchedulerTaskDetailResponse>(`/api/scheduler-tasks/${enc(uid)}/${enc(taskID)}`)

// ---- 本地 Mock 定时任务 CRUD ----
export const getMockSchedulerTasks = () => get<MockSchedulerTasksResponse>('/api/mock-scheduler-tasks')
export const createMockSchedulerTask = (input: MockSchedulerTaskInput) => post<MockSchedulerTask>('/api/mock-scheduler-tasks', input)
export const updateMockSchedulerTask = (id: string, input: MockSchedulerTaskInput) =>
  put<MockSchedulerTask>(`/api/mock-scheduler-tasks/${enc(id)}`, input)
export const deleteMockSchedulerTask = (id: string) => del<OkResponse>(`/api/mock-scheduler-tasks/${enc(id)}`)

// ---- 本地自动化 ----
export const getAutomations = () => get<AutomationsResponse>('/api/automations')
export const updateAutomation = (id: string, input: AutomationInput) => put<Automation>(`/api/automations/${enc(id)}`, input)
export const runAutomation = (id: string) => post<ActionResult>(`/api/automations/${enc(id)}/run`)

// ---- OAuth ----
export const startOAuth = (input: OAuthStartInput) => post<OAuthStartResponse>('/api/oauth/start', input)
export const pollOAuth = (id: string) => post<OAuthPollResponse>(`/api/oauth/${enc(id)}/poll`)

// ---- 新增端点（M2+，后端实现随后续里程碑补齐，当前由 MSW 提供） ----
export const startProbe = () => post<ProbeResponse>('/api/probe')
export const getStatsSummary = () => get<StatsSummary>('/api/stats/summary')
export const getLogs = () => get<LogsResponse>('/api/logs')
export const getModelPricing = () => get<ModelPricingResponse>('/api/models/pricing')

// 供页面复用的类型导出（避免页面重复 import 路径）。
export type { Account }
