// 一键任务面板测试：按钮 pending 禁用、结果逐行渲染、跳过标记、失败文案、空结果、自动化提示。
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../mocks/server'
import { BatchPanel } from '../components/BatchPanel'

test('BatchPanel 渲染四个操作按钮', () => {
  render(<BatchPanel />)
  expect(screen.getByRole('button', { name: '批量签到' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: '批量旅行巡检' })).toBeInTheDocument()
  expect(screen.getByRole('button', { name: '批量刷新凭据' })).toBeInTheDocument()
})

test('BatchPanel 点击后 pending 期间按钮禁用', async () => {
  // 用永不 resolve 的 handler 撑开 pending 窗口。
  server.use(http.post('/api/batch-actions', () => new Promise(() => undefined)))
  render(<BatchPanel />)
  const btn = screen.getByRole('button', { name: '批量签到' })
  fireEvent.click(btn)
  await waitFor(() => expect(btn).toBeDisabled())
  expect(screen.getByRole('button', { name: '批量旅行巡检' })).toBeDisabled()
})

test('BatchPanel 结果逐行渲染：成功/跳过/12153/失败', async () => {
  render(<BatchPanel />)
  fireEvent.click(screen.getByRole('button', { name: '批量签到' }))
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  // 成功
  expect(screen.getByText('签到成功')).toBeInTheDocument()
  // 跳过
  expect(screen.getByText('已跳过')).toBeInTheDocument()
  // 12153
  expect(screen.getByText('会话失效')).toBeInTheDocument()
  // 失败
  expect(screen.getByText('上游签到接口调用失败')).toBeInTheDocument()
})

test('BatchPanel 禁用账号行渲染跳过标记且非失败态', async () => {
  render(<BatchPanel />)
  fireEvent.click(screen.getByRole('button', { name: '批量签到' }))
  await waitFor(() => expect(screen.getByText('阿红')).toBeInTheDocument())
  const row = screen.getByText('阿红').closest('tr') as HTMLElement
  expect(row.textContent).toContain('已跳过')
  expect(row.textContent).not.toContain('失败')
})

test('BatchPanel 失败行渲染错误文案', async () => {
  render(<BatchPanel />)
  fireEvent.click(screen.getByRole('button', { name: '批量签到' }))
  await waitFor(() => expect(screen.getByText('阿绿')).toBeInTheDocument())
  const row = screen.getByText('阿绿').closest('tr') as HTMLElement
  expect(row.textContent).toContain('上游签到接口调用失败')
})

test('BatchPanel 空结果渲染空态', async () => {
  server.use(http.post('/api/batch-actions', () => HttpResponse.json({ results: [] })))
  render(<BatchPanel />)
  fireEvent.click(screen.getByRole('button', { name: '批量签到' }))
  await waitFor(() => expect(screen.getByText(/没有可执行的账号/)).toBeInTheDocument())
})

test('BatchPanel 请求失败渲染错误提示', async () => {
  server.use(http.post('/api/batch-actions', () => HttpResponse.json({ error: '服务端已开启只读模式' }, { status: 403 })))
  render(<BatchPanel />)
  fireEvent.click(screen.getByRole('button', { name: '批量签到' }))
  await waitFor(() => expect(screen.getByText('服务端已开启只读模式')).toBeInTheDocument())
})

test('BatchPanel 自动化所有权提示文案', () => {
  render(<BatchPanel />)
  expect(screen.getByText(/默认只开新活动探测/)).toBeInTheDocument()
  expect(screen.getByText(/接管网关任务需部署层确认/)).toBeInTheDocument()
})
