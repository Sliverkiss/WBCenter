// App 冒烟测试：未登录渲染登录页，已登录渲染壳与概览。
import { render, screen, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import App from '../App'

function stubSession(authenticated: boolean) {
  vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
    const url = String(input)
    if (url.includes('/api/session')) {
      return new Response(JSON.stringify({ authenticated }), { status: 200 })
    }
    if (url.includes('/api/overview')) {
      return new Response(JSON.stringify({ healthy: 1, accounts: 2, expired: 0, automations: 1 }), { status: 200 })
    }
    return new Response('{}', { status: 200 })
  }))
}

test('未登录时渲染登录页', async () => {
  stubSession(false)
  render(<App />)
  await waitFor(() => expect(screen.getByText('账号工作台')).toBeInTheDocument())
  expect(screen.getByText('进入控制台')).toBeInTheDocument()
})

test('已登录时渲染侧栏导航与概览数据', async () => {
  stubSession(true)
  render(<App />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  expect(screen.getByText('概览')).toBeInTheDocument()
  expect(screen.getByText('账号管理')).toBeInTheDocument()
})
