// M5 OAuth 页升级测试：完整三步流（发起→授权链接→轮询状态机）+ 超时倒计时 + 区域路由 + 重试。
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { describe, expect, test, vi } from 'vitest'
import { OAuth } from '../pages/OAuth'
import { server } from '../mocks/server'

const noop = () => undefined

describe('OAuth 页（M5 完整流）', () => {
  test('发起后展示授权链接（新窗口打开）与轮询中状态', async () => {
    render(<OAuth onNotice={noop} onDone={noop} />)
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    const link = await screen.findByRole('link', { name: /打开 WorkBuddy 授权页/ })
    expect(link).toHaveAttribute('href', expect.stringContaining('copilot.tencent.com'))
    expect(link).toHaveAttribute('target', '_blank')
    // 轮询中状态可见（waiting）。
    expect(await screen.findByText(/等待授权/)).toBeInTheDocument()
  })

  test('轮询 waiting→waiting→success 后展示账号信息并回调 onDone', async () => {
    const onDone = vi.fn()
    const onNotice = vi.fn()
    render(<OAuth onNotice={onNotice} onDone={onDone} />)
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    // 默认 MSW 序列：waiting→waiting→success（u-1003 新账号）。
    expect(await screen.findByText(/授权成功/), { timeout: 15000 }).toBeInTheDocument()
    expect(screen.getByText(/新账号/)).toBeInTheDocument()
    expect(screen.getByText(/u-1003/)).toBeInTheDocument()
    await waitFor(() => expect(onDone).toHaveBeenCalled())
  }, 20000)

  test('轮询 error 终态展示错误文案与重试按钮', async () => {
    server.use(
      http.post('/api/oauth/:id/poll', () => HttpResponse.json({ status: 'error', message: '上游 503 不可用' })),
    )
    render(<OAuth onNotice={noop} onDone={noop} />)
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    expect(await screen.findByText(/授权失败/)).toBeInTheDocument()
    expect(screen.getByText(/上游 503 不可用/)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /重试/ })).toBeInTheDocument()
  })

  test('重试按钮重置流程可再次发起', async () => {
    server.use(
      http.post('/api/oauth/:id/poll', () => HttpResponse.json({ status: 'error', message: '失败' })),
    )
    render(<OAuth onNotice={noop} onDone={noop} />)
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    const retry = await screen.findByRole('button', { name: /重试/ })
    fireEvent.click(retry)
    // 重置后回到可发起状态。
    expect(await screen.findByRole('button', { name: /发起 OAuth 授权/ })).toBeInTheDocument()
  })

  test('轮询 timeout 终态展示超时文案', async () => {
    server.use(
      http.post('/api/oauth/:id/poll', () => HttpResponse.json({ status: 'timeout', message: '授权超时（5 分钟未完成），请重新发起' })),
    )
    render(<OAuth onNotice={noop} onDone={noop} />)
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    expect(await screen.findByText(/授权超时/)).toBeInTheDocument()
  })

  test('轮询中展示 5 分钟倒计时', async () => {
    server.use(
      http.post('/api/oauth/:id/poll', () => HttpResponse.json({ status: 'waiting' })),
    )
    render(<OAuth onNotice={noop} onDone={noop} />)
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    // 倒计时初始接近 5:00（mm:ss 形态）。
    expect(await screen.findByText(/剩余时间/)).toBeInTheDocument()
    expect(screen.getByText(/4:5\d|5:00/)).toBeInTheDocument()
  })

  test('global 区域授权链接指向国际版域名', async () => {
    render(<OAuth onNotice={noop} onDone={noop} />)
    fireEvent.change(screen.getByLabelText(/账号区域/), { target: { value: 'global' } })
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    const link = await screen.findByRole('link', { name: /打开 WorkBuddy 授权页/ })
    expect(link).toHaveAttribute('href', expect.stringContaining('www.workbuddy.ai'))
  })

  test('区域选择带说明文案（中国大陆/国际版）', () => {
    render(<OAuth onNotice={noop} onDone={noop} />)
    expect(screen.getByText(/中国大陆/)).toBeInTheDocument()
    expect(screen.getByText(/国际版/)).toBeInTheDocument()
    // 说明文案：两个区域的差异提示。
    expect(screen.getByText(/copilot\.tencent\.com/)).toBeInTheDocument()
    expect(screen.getByText(/workbuddy\.ai/)).toBeInTheDocument()
  })

  test('卸载时清理轮询定时器（不再发起 poll）', async () => {
    const pollSpy = vi.fn(() => HttpResponse.json({ status: 'waiting' }))
    server.use(http.post('/api/oauth/:id/poll', pollSpy))
    const { unmount } = render(<OAuth onNotice={noop} onDone={noop} />)
    fireEvent.click(screen.getByRole('button', { name: /发起 OAuth 授权/ }))
    await waitFor(() => expect(pollSpy.mock.calls.length).toBeGreaterThan(0))
    unmount()
    const callsAtUnmount = pollSpy.mock.calls.length
    await new Promise(r => setTimeout(r, 3500))
    expect(pollSpy.mock.calls.length).toBe(callsAtUnmount)
  }, 10000)
})
