// 路由集成测试：hash 切换渲染对应页面、未登录守卫、侧栏导航点击。全部走 MSW。
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../mocks/server'
import App from '../App'

test('未登录一律渲染 Login（守卫）', async () => {
  window.location.hash = '#/accounts'
  server.use(http.get('/api/session', () => HttpResponse.json({ authenticated: false, read_only: false, using_default_password: false, timezone: 'Asia/Shanghai' })))
  render(<App />)
  await waitFor(() => expect(screen.getByText('进入控制台')).toBeInTheDocument())
  expect(screen.queryByText('账号管理', { selector: 'h1' })).not.toBeInTheDocument()
})

test('hash #/credits 渲染积分管理页', async () => {
  window.location.hash = '#/credits'
  render(<App />)
  await waitFor(() => expect(screen.getByRole('heading', { name: '积分管理' })).toBeInTheDocument())
})

test('点击侧栏导航切换 hash 并渲染对应页面', async () => {
  render(<App />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /模型管理/ }))
  await waitFor(() => expect(screen.getByRole('heading', { name: '模型管理' })).toBeInTheDocument())
  expect(window.location.hash).toBe('#/models')
})

test('hash #/oauth 渲染添加账号页', async () => {
  window.location.hash = '#/oauth'
  render(<App />)
  await waitFor(() => expect(screen.getByRole('heading', { name: '添加账号' })).toBeInTheDocument())
})

test('hash #/scheduler 渲染账号云端定时任务页', async () => {
  window.location.hash = '#/scheduler'
  render(<App />)
  await waitFor(() => expect(screen.getByRole('heading', { name: '账号云端定时任务' })).toBeInTheDocument())
})
