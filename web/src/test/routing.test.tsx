// 路由集成测试：hash 切换渲染对应页面、未登录守卫、侧栏导航点击。
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { vi } from 'vitest'
import App from '../App'

function stubAll() {
  vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
    const url = String(input)
    if (url.includes('/api/session')) return new Response(JSON.stringify({ authenticated: true }), { status: 200 })
    if (url.includes('/api/overview')) return new Response(JSON.stringify({ healthy: 1, accounts: 1, expired: 0, automations: 0 }), { status: 200 })
    if (url.includes('/api/accounts')) return new Response(JSON.stringify({ accounts: [] }), { status: 200 })
    if (url.includes('/api/credits')) return new Response(JSON.stringify({ items: [] }), { status: 200 })
    if (url.includes('/api/models')) return new Response(JSON.stringify({ items: [] }), { status: 200 })
    if (url.includes('/api/automations')) return new Response(JSON.stringify({ items: [], runs: [] }), { status: 200 })
    if (url.includes('/api/scheduler-tasks')) return new Response(JSON.stringify({ items: [], note: '' }), { status: 200 })
    if (url.includes('/api/mock-scheduler-tasks')) return new Response(JSON.stringify({ items: [], note: '' }), { status: 200 })
    if (url.includes('/api/activities')) return new Response(JSON.stringify({ items: [] }), { status: 200 })
    return new Response('{}', { status: 200 })
  }))
}

test('未登录一律渲染 Login（守卫）', async () => {
  window.location.hash = '#/accounts'
  vi.stubGlobal('fetch', vi.fn(async () => new Response(JSON.stringify({ authenticated: false }), { status: 200 })))
  render(<App />)
  await waitFor(() => expect(screen.getByText('进入控制台')).toBeInTheDocument())
  expect(screen.queryByText('账号管理', { selector: 'h1' })).not.toBeInTheDocument()
})

test('hash #/credits 渲染积分管理页', async () => {
  window.location.hash = '#/credits'
  stubAll()
  render(<App />)
  await waitFor(() => expect(screen.getByRole('heading', { name: '积分管理' })).toBeInTheDocument())
})

test('点击侧栏导航切换 hash 并渲染对应页面', async () => {
  stubAll()
  render(<App />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /模型管理/ }))
  await waitFor(() => expect(screen.getByRole('heading', { name: '模型管理' })).toBeInTheDocument())
  expect(window.location.hash).toBe('#/models')
})

test('hash #/oauth 渲染添加账号页', async () => {
  window.location.hash = '#/oauth'
  stubAll()
  render(<App />)
  await waitFor(() => expect(screen.getByRole('heading', { name: '添加账号' })).toBeInTheDocument())
})

test('hash #/scheduler 渲染账号云端定时任务页', async () => {
  window.location.hash = '#/scheduler'
  stubAll()
  render(<App />)
  await waitFor(() => expect(screen.getByRole('heading', { name: '账号云端定时任务' })).toBeInTheDocument())
})
