// 页面冒烟测试：渲染 + 关键数据出现在 DOM（fetch 用 vi.stubGlobal mock）。
import { render, screen, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import { Overview } from '../pages/Overview'
import { Accounts } from '../pages/Accounts'
import { Credits } from '../pages/Credits'
import { Models } from '../pages/Models'
import { Automation } from '../pages/Automation'
import { Scheduler } from '../pages/Scheduler'
import { Activities } from '../pages/Activities'
import { OAuth } from '../pages/OAuth'
import { Login } from '../pages/Login'

function stubFetch(handler: (url: string) => unknown) {
  vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
    const body = handler(String(input))
    return new Response(JSON.stringify(body ?? {}), { status: 200 })
  }))
}

test('Overview 渲染统计卡与快速入口', async () => {
  stubFetch(url => url.includes('/api/overview') ? { healthy: 2, accounts: 5, expired: 1, automations: 3 } : {})
  render(<Overview onPage={() => undefined} />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  expect(screen.getByText('2')).toBeInTheDocument()
  expect(screen.getByText('5')).toBeInTheDocument()
  expect(screen.getByText('查看今日积分汇总 →')).toBeInTheDocument()
})

test('Overview 上游错误时渲染 ErrorView', async () => {
  vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({ error: '上游不可用' }), { status: 502 })))
  render(<Overview onPage={() => undefined} />)
  await waitFor(() => expect(screen.getByText('上游不可用')).toBeInTheDocument())
})

test('Accounts 渲染账号行与操作按钮', async () => {
  stubFetch(url => url.includes('/api/accounts') ? { accounts: [{ uid: 'u1', nickname: '阿明', domain: 'cn', expired: false, needs_refresh: false }] } : {})
  render(<Accounts onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  expect(screen.getByText('刷新凭据')).toBeInTheDocument()
  expect(screen.getByText('签到')).toBeInTheDocument()
  expect(screen.getByText('旅行巡检')).toBeInTheDocument()
})

test('Accounts 空列表渲染空态', async () => {
  stubFetch(() => ({ accounts: [] }))
  render(<Accounts onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText(/尚无账号/)).toBeInTheDocument())
})

test('Credits 渲染积分表格与数据', async () => {
  stubFetch(url => url.includes('/api/credits') ? { items: [{ uid: 'u1', nickname: '阿明', current: 1200, today_allocated: 100, today_consumed: 40, today_remaining: 60 }] } : {})
  render(<Credits />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  expect(screen.getByText('今日发放额度')).toBeInTheDocument()
})

test('Models 渲染模型列表与状态徽标', async () => {
  stubFetch(url => url.includes('/api/models') ? { items: [{ uid: 'u1', nickname: '阿明', id: 'wb-pro', name: 'WorkBuddy Pro' }] } : {})
  render(<Models />)
  await waitFor(() => expect(screen.getByText('wb-pro')).toBeInTheDocument())
  expect(screen.getByText('来自上游')).toBeInTheDocument()
})

test('Models 上游错误行渲染坏徽标', async () => {
  stubFetch(url => url.includes('/api/models') ? { items: [{ uid: 'u1', error: 'session dead' }] } : {})
  render(<Models />)
  await waitFor(() => expect(screen.getByText('session dead')).toBeInTheDocument())
})

test('Automation 渲染任务列表与开关', async () => {
  stubFetch(url => url.includes('/api/automations') ? {
    items: [{ id: 'a1', name: '新活动探测', enabled: true, every_minutes: 60, last_run_at: '', last_result: '' }],
    runs: [],
  } : {})
  render(<Automation onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText('新活动探测')).toBeInTheDocument())
  expect(screen.getByRole('button', { name: '切换新活动探测' })).toBeInTheDocument()
})

test('Automation 渲染最近运行记录', async () => {
  stubFetch(url => url.includes('/api/automations') ? {
    items: [{ id: 'a1', name: '签到', enabled: false, every_minutes: 120 }],
    runs: [{ ok: true, action: '签到', message: '全部成功' }],
  } : {})
  render(<Automation onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText('全部成功')).toBeInTheDocument())
  expect(screen.getByText('成功')).toBeInTheDocument()
})

test('Scheduler 渲染 Mock 表单与上游只读表', async () => {
  stubFetch(url => {
    if (url.includes('/api/scheduler-tasks')) return { items: [], note: '上游未授权' }
    if (url.includes('/api/mock-scheduler-tasks')) return { items: [{ id: 'm1', account_uid: 'u1', nickname: '阿明', name: '每日摘要', cron: '0 9 * * *', prompt: 'hi', enabled: true }], note: '本地 Mock' }
    if (url.includes('/api/accounts')) return { accounts: [{ uid: 'u1', nickname: '阿明' }] }
    return {}
  })
  render(<Scheduler />)
  await waitFor(() => expect(screen.getByText('每日摘要')).toBeInTheDocument())
  expect(screen.getByText('创建 Mock 云端任务')).toBeInTheDocument()
  expect(screen.getByText('上游真实任务（只读）')).toBeInTheDocument()
})

test('Activities 渲染成长任务表', async () => {
  stubFetch(url => url.includes('/api/activities') ? { items: [{ uid: 'u1', nickname: '阿明', code: 'T1', name: '每日签到', status: '可领取', reward: 100, new: true }] } : {})
  render(<Activities />)
  await waitFor(() => expect(screen.getByText('每日签到')).toBeInTheDocument())
  expect(screen.getByText('新发现')).toBeInTheDocument()
})

test('Activities 空列表渲染空态', async () => {
  stubFetch(() => ({ items: [] }))
  render(<Activities />)
  await waitFor(() => expect(screen.getByText(/当前没有可读取的成长任务/)).toBeInTheDocument())
})

test('OAuth 渲染区域选择与发起按钮', async () => {
  stubFetch(() => ({}))
  render(<OAuth onNotice={() => undefined} onDone={() => undefined} />)
  expect(screen.getByText('中国大陆')).toBeInTheDocument()
  expect(screen.getByRole('button', { name: '发起 OAuth 授权' })).toBeInTheDocument()
})

test('Login 渲染表单并可提交登录', async () => {
  const calls: string[] = []
  vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
    calls.push(String(input))
    return new Response(JSON.stringify({ ok: true }), { status: 200 })
  }))
  const onDone = vi.fn().mockResolvedValue(undefined)
  render(<Login onDone={onDone} />)
  expect(screen.getByText('账号工作台')).toBeInTheDocument()
  const { fireEvent } = await import('@testing-library/react')
  fireEvent.click(screen.getByRole('button', { name: '进入控制台' }))
  await waitFor(() => expect(onDone).toHaveBeenCalled())
  expect(calls.some(u => u.includes('/api/login'))).toBe(true)
})

test('Login 失败时渲染错误信息', async () => {
  vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({ error: '密码错误' }), { status: 401 })))
  render(<Login onDone={vi.fn()} />)
  const { fireEvent } = await import('@testing-library/react')
  fireEvent.click(screen.getByRole('button', { name: '进入控制台' }))
  await waitFor(() => expect(screen.getByText('密码错误')).toBeInTheDocument())
})
