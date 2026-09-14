// App 冒烟测试：未登录渲染登录页，已登录渲染壳与概览。会话态由 MSW handler 提供。
import { render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../mocks/server'
import App from '../App'

test('未登录时渲染登录页', async () => {
  server.use(http.get('/api/session', () => HttpResponse.json({ authenticated: false, read_only: false, using_default_password: true, timezone: 'Asia/Shanghai' })))
  render(<App />)
  await waitFor(() => expect(screen.getByText('账号工作台')).toBeInTheDocument())
  expect(screen.getByText('进入控制台')).toBeInTheDocument()
})

test('已登录时渲染侧栏导航与概览数据', async () => {
  render(<App />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  expect(screen.getByText('概览')).toBeInTheDocument()
  expect(screen.getByText('账号管理')).toBeInTheDocument()
})
