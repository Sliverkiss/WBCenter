// services/api 封装测试：错误信封解析、401 全局登出回调、域方法 URL/方法对照契约。
import { http, HttpResponse } from 'msw'
import { server } from '../mocks/server'
import {
  request,
  setUnauthorizedHandler,
  getSession,
  login,
  getAccounts,
  getCredits,
  getCreditDetail,
  startProbe,
  getStatsSummary,
  getLogs,
  getModelPricing,
} from '../services/api'
import { isAccount, isCreditSummary, isProbeResult } from '../types'

describe('request 错误信封', () => {
  test('非 2xx 时抛出 body.error 中文文案', async () => {
    server.use(http.get('/api/test-err', () => HttpResponse.json({ error: '上游不可用' }, { status: 502 })))
    await expect(request('/api/test-err')).rejects.toThrow('上游不可用')
  })

  test('非 2xx 且无 error 字段时回退默认文案', async () => {
    server.use(http.get('/api/test-err2', () => HttpResponse.json({}, { status: 500 })))
    await expect(request('/api/test-err2')).rejects.toThrow('请求失败')
  })
})

describe('401 统一登出流转', () => {
  test('401 触发登出回调并把 hash 重置为 #/', async () => {
    const onLogout = vi.fn()
    setUnauthorizedHandler(onLogout)
    window.location.hash = '#/credits'
    server.use(http.get('/api/accounts', () => HttpResponse.json({ error: '未登录或会话已过期' }, { status: 401 })))
    await expect(getAccounts()).rejects.toThrow('未登录或会话已过期')
    expect(onLogout).toHaveBeenCalledTimes(1)
    expect(window.location.hash).toBe('#/')
    setUnauthorizedHandler(null)
  })

  test('登录端点自身的 401 不触发全局登出（避免登录失败误登出）', async () => {
    const onLogout = vi.fn()
    setUnauthorizedHandler(onLogout)
    await expect(login({ username: 'admin', password: 'wrong' })).rejects.toThrow('用户名或密码错误')
    expect(onLogout).not.toHaveBeenCalled()
    setUnauthorizedHandler(null)
  })
})

describe('域方法契约形状', () => {
  test('getSession 返回 SessionInfo', async () => {
    const s = await getSession()
    expect(s.authenticated).toBe(true)
    expect(typeof s.timezone).toBe('string')
  })

  test('getAccounts 返回的每个账号通过 isAccount 守卫', async () => {
    const d = await getAccounts()
    expect(d.accounts.length).toBeGreaterThan(0)
    expect(d.accounts.every(isAccount)).toBe(true)
  })

  test('getCredits 与 getCreditDetail 返回的条目通过 isCreditSummary 守卫', async () => {
    const all = await getCredits()
    expect(all.items.every(isCreditSummary)).toBe(true)
    const one = await getCreditDetail('u-1001')
    expect(one.items).toHaveLength(1)
    expect(one.items.every(isCreditSummary)).toBe(true)
  })

  test('startProbe 用 POST 且结果通过 isProbeResult 守卫', async () => {
    const d = await startProbe()
    expect(d.results.every(isProbeResult)).toBe(true)
    expect(d.results.some(r => !r.ok && r.error)).toBe(true)
  })

  test('getStatsSummary / getLogs / getModelPricing 返回契约字段', async () => {
    const stats = await getStatsSummary()
    expect(typeof stats.credits_current_total).toBe('number')
    const logs = await getLogs()
    expect(Array.isArray(logs.lines)).toBe(true)
    expect(logs.lines[0].message).not.toMatch(/token|cookie|password/i)
    const pricing = await getModelPricing()
    expect(pricing.items.some(i => i.free)).toBe(true)
  })
})
