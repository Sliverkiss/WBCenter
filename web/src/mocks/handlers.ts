// MSW handlers：覆盖 docs/api-contract.md 全部端点（22 存量 + 4 新增），dev 与 test 双环境共用。
// 未登录态通过 /api/login 切换模块级标志模拟（与后端会话语义对齐的最小模型）。
import { http, HttpResponse } from 'msw'
import type { MockSchedulerTask, MockSchedulerTaskInput } from '../types'
import {
  accounts,
  activities,
  automations,
  batchActionResponse,
  credits,
  logs,
  lottery,
  mockSchedulerTasks,
  modelPricing,
  models,
  overview,
  probeResponse,
  runs,
  schedulerTasks,
  sessionInfo,
  statsSummary,
} from './fixtures'

const json = (body: object, status = 200) => HttpResponse.json(body, { status })
const err = (status: number, error: string) => json({ error }, status)

// 模拟会话状态：login 成功后置 true；resetMockState 在每用例间复位。
let authed = true
export function resetMockState() {
  authed = true
  mockTasks = mockSchedulerTasks.map(t => ({ ...t }))
}

const unauthorized = () => err(401, '未登录或会话已过期')

let mockTasks: MockSchedulerTask[] = mockSchedulerTasks.map(t => ({ ...t }))

export const handlers = [
  // ---- 会话 ----
  http.get('/api/session', () => json({ ...sessionInfo, authenticated: authed })),
  http.post('/api/login', async ({ request }) => {
    const body = (await request.json()) as { username?: string; password?: string }
    if (body.username === 'admin' && body.password === 'workbuddy') {
      authed = true
      return json({ ok: true })
    }
    return err(401, '用户名或密码错误')
  }),
  http.post('/api/logout', () => {
    authed = false
    return json({ ok: true })
  }),

  // ---- 概览与账号 ----
  http.get('/api/overview', () => (authed ? json(overview) : unauthorized())),
  http.get('/api/accounts', () => (authed ? json({ accounts, warnings: [] }) : unauthorized())),
  http.get('/api/credits', () => (authed ? json({ items: credits }) : unauthorized())),
  http.get('/api/credits/:uid', ({ params }) => {
    if (!authed) return unauthorized()
    return json({ items: credits.filter(c => c.uid === params.uid) })
  }),

  // ---- 模型与活动 ----
  http.get('/api/models', () => (authed ? json({ items: models }) : unauthorized())),
  http.get('/api/activities', () => (authed ? json({ items: activities }) : unauthorized())),

  // ---- 云端定时任务（只读） ----
  http.get('/api/scheduler-tasks', () =>
    authed
      ? json({
          items: schedulerTasks,
          write_enabled: false,
          note: '这是各 WorkBuddy 账号的云端定时任务。创建接口请求体尚未经过真实认证验证，当前保持只读。',
        })
      : unauthorized(),
  ),
  http.get('/api/scheduler-tasks/:uid/:taskID', ({ params }) => {
    if (!authed) return unauthorized()
    const task = schedulerTasks.find(t => t.uid === params.uid && t.id === params.taskID)
    if (!task) return err(502, '账号云端定时任务详情读取不可用')
    return json({ item: { id: task.id, name: task.name, status: task.status } })
  }),

  // ---- 本地 Mock 定时任务 CRUD ----
  http.get('/api/mock-scheduler-tasks', () =>
    authed
      ? json({ items: mockTasks, mode: 'mock', note: '本地 Mock 骨架：仅写入控制台 data/state.json，不会提交到 WorkBuddy 上游。' })
      : unauthorized(),
  ),
  http.post('/api/mock-scheduler-tasks', async ({ request }) => {
    if (!authed) return unauthorized()
    const input = (await request.json()) as MockSchedulerTaskInput
    if (!input.name || input.name.length > 100) return err(400, '任务名称需为 1 至 100 个字符')
    const now = '2026-09-14T09:00:00+08:00'
    const nickname = accounts.find(a => a.uid === input.account_uid)?.nickname ?? ''
    const task: MockSchedulerTask = { id: `mock-${mockTasks.length + 1}`, ...input, created_at: now, updated_at: now, nickname }
    mockTasks = [task, ...mockTasks]
    return json(task, 201)
  }),
  http.put('/api/mock-scheduler-tasks/:id', async ({ params, request }) => {
    if (!authed) return unauthorized()
    const input = (await request.json()) as MockSchedulerTaskInput
    const task = mockTasks.find(t => t.id === params.id)
    if (!task) return err(404, 'Mock 任务不存在')
    const updated: MockSchedulerTask = { ...task, name: input.name, cron: input.cron, prompt: input.prompt, enabled: input.enabled, updated_at: '2026-09-14T09:00:00+08:00' }
    mockTasks = mockTasks.map(t => (t.id === params.id ? updated : t))
    return json(updated)
  }),
  http.delete('/api/mock-scheduler-tasks/:id', ({ params }) => {
    if (!authed) return unauthorized()
    if (!mockTasks.some(t => t.id === params.id)) return err(404, 'Mock 任务不存在')
    mockTasks = mockTasks.filter(t => t.id !== params.id)
    return json({ ok: true })
  }),

  // ---- 本地自动化 ----
  http.get('/api/automations', () => (authed ? json({ items: automations, runs }) : unauthorized())),
  http.put('/api/automations/:id', async ({ params, request }) => {
    if (!authed) return unauthorized()
    const a = automations.find(x => x.id === params.id)
    if (!a) return err(404, '自动化任务不存在')
    const input = (await request.json()) as { enabled?: boolean; every_minutes?: number }
    return json({ ...a, enabled: Boolean(input.enabled), every_minutes: Number(input.every_minutes ?? a.every_minutes) })
  }),
  http.post('/api/automations/:id/run', ({ params }) => {
    if (!authed) return unauthorized()
    if (!automations.some(x => x.id === params.id)) return err(404, '自动化任务不存在')
    return json({ ok: true, message: '成功 1，失败 0' })
  }),

  // ---- 账号动作 ----
  http.post('/api/accounts/:uid/actions/:action', ({ params }) => {
    if (!authed) return unauthorized()
    if (!['checkin', 'travel', 'refresh'].includes(String(params.action))) return err(400, '不支持的动作')
    return json({ ok: true, message: '成功 1，失败 0' })
  }),

  // ---- OAuth ----
  http.post('/api/oauth/start', async ({ request }) => {
    const body = (await request.json()) as { region?: string }
    if (body.region !== 'cn' && body.region !== 'global') return err(502, '未知区域')
    return json({ id: 'mock-flow-1', url: 'https://copilot.tencent.com/oauth/authorize?state=mock' })
  }),
  http.post('/api/oauth/:id/poll', () => json({ status: 'success', uid: 'u-1003', nickname: '新账号' })),

  // ---- 新增端点（M2+ 契约，后端未实现，mock 先行） ----
  http.post('/api/probe', () => (authed ? json(probeResponse) : unauthorized())),
  http.get('/api/stats/summary', () => (authed ? json(statsSummary) : unauthorized())),
  http.post('/api/batch-actions', async ({ request }) => {
    if (!authed) return unauthorized()
    const body = (await request.json()) as { action?: string; uids?: string[] }
    if (!body.action || !['checkin', 'travel', 'refresh'].includes(body.action)) return err(400, '不支持的动作')
    // uids 指定子集时按子集过滤样本；缺省返回全部四类样本。
    const results = Array.isArray(body.uids) && body.uids.length > 0
      ? batchActionResponse.results.filter(r => body.uids?.includes(r.uid))
      : batchActionResponse.results
    return json({ results })
  }),
  http.get('/api/logs', () => (authed ? json(logs) : unauthorized())),
  http.get('/api/models/pricing', () => (authed ? json(modelPricing) : unauthorized())),
  http.get('/api/activities/lottery', () => (authed ? json(lottery) : unauthorized())),
]
