// M5 日志查看页测试：列表渲染 / 级别过滤 / 10s 轮询自动刷新 / 卸载清理。
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { afterEach, describe, expect, test, vi } from 'vitest'
import { Logs } from '../pages/Logs'
import { server } from '../mocks/server'

describe('日志查看页（M5）', () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  test('渲染最近日志列表（时间/级别/动作/消息）', async () => {
    render(<Logs />)
    // fixtures 提供 5 条样本（含 info/warn/error）。
    expect(await screen.findByText('账号 u-1001 签到成功')).toBeInTheDocument()
    expect(screen.getByText(/上游会话失效（12153）/)).toBeInTheDocument()
    expect(screen.getByText(/OAuth 轮询失败/)).toBeInTheDocument()
    // 动作列可见。
    expect(screen.getAllByText(/login|oauth|probe|run|batch/).length).toBeGreaterThan(0)
  })

  test('级别过滤 all/info/error', async () => {
    render(<Logs />)
    await screen.findByText('账号 u-1001 签到成功')

    // 过滤 error：只剩 error 级一条。
    fireEvent.click(screen.getByRole('button', { name: 'error' }))
    expect(screen.getByText(/OAuth 轮询失败/)).toBeInTheDocument()
    expect(screen.queryByText('账号 u-1001 签到成功')).not.toBeInTheDocument()

    // 过滤 info：error 消失，info 回来。
    fireEvent.click(screen.getByRole('button', { name: 'info' }))
    expect(screen.getByText('账号 u-1001 签到成功')).toBeInTheDocument()
    expect(screen.queryByText(/OAuth 轮询失败/)).not.toBeInTheDocument()

    // all 恢复全部 5 条。
    fireEvent.click(screen.getByRole('button', { name: 'all' }))
    expect(screen.getByText(/OAuth 轮询失败/)).toBeInTheDocument()
    expect(screen.getByText(/上游会话失效（12153）/)).toBeInTheDocument()
  })

  test('10 秒轮询自动刷新（fake timers 推进后重新请求）', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true })
    let calls = 0
    server.use(
      http.get('/api/logs', () => {
        calls++
        return HttpResponse.json({ lines: [{ at: '2026-09-14T08:30:00+08:00', level: 'info', action: 'run', message: `第 ${calls} 批` }] })
      }),
    )
    render(<Logs />)
    await screen.findByText('第 1 批')
    const before = calls
    await vi.advanceTimersByTimeAsync(10500)
    await waitFor(() => expect(calls).toBeGreaterThan(before))
  })

  test('卸载时清理轮询 interval（不再请求）', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true })
    let calls = 0
    server.use(
      http.get('/api/logs', () => {
        calls++
        return HttpResponse.json({ lines: [] })
      }),
    )
    const { unmount } = render(<Logs />)
    await waitFor(() => expect(calls).toBe(1))
    unmount()
    const atUnmount = calls
    await vi.advanceTimersByTimeAsync(25000)
    expect(calls).toBe(atUnmount)
  })

  test('请求失败展示错误视图', async () => {
    server.use(http.get('/api/logs', () => HttpResponse.json({ error: '服务端错误' }, { status: 500 })))
    render(<Logs />)
    expect(await screen.findByText(/服务端错误/)).toBeInTheDocument()
  })

  test('空日志展示空态提示', async () => {
    server.use(http.get('/api/logs', () => HttpResponse.json({ lines: [] })))
    render(<Logs />)
    expect(await screen.findByText(/暂无运行日志/)).toBeInTheDocument()
  })

  test('调用 /api/logs 携带 limit 参数', async () => {
    let seenUrl = ''
    server.use(
      http.get('/api/logs', ({ request }) => {
        seenUrl = request.url
        return HttpResponse.json({ lines: [] })
      }),
    )
    render(<Logs />)
    await screen.findByText(/暂无运行日志/)
    expect(seenUrl).toContain('limit=100')
  })
})
