// 统计页测试：数据转换单测 + 页面渲染 + 轮询清理。
// echarts 在 jsdom 中无法真实渲染 canvas，useEcharts 通过依赖注入桩验证 init/dispose 生命周期。
import { render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { vi } from 'vitest'
import { server } from '../mocks/server'
import { Stats } from '../pages/Stats'
import { buildStatusChartOption, buildCreditChartOption } from '../utils/statsChart'
import { statsSummary } from '../mocks/fixtures'
import type { StatsSummary } from '../types'

// ---- 数据转换单测 ----

test('buildStatusChartOption 把四类状态转成条形图 series', () => {
  const opt = buildStatusChartOption(statsSummary)
  const series = opt.series as { data: number[] }[]
  expect(series[0].data).toEqual([3, 1, 1, 1]) // 健康/过期/会话死/禁用
  const xAxis = opt.xAxis as { data: string[] }
  expect(xAxis.data).toEqual(['健康', '已过期', '会话失效', '已禁用'])
})

test('buildStatusChartOption 空数据全零', () => {
  const empty: StatsSummary = { ...statsSummary, status_healthy: 0, status_expired: 0, status_session_dead: 0, status_disabled: 0 }
  const opt = buildStatusChartOption(empty)
  const series = opt.series as { data: number[] }[]
  expect(series[0].data).toEqual([0, 0, 0, 0])
})

test('buildCreditChartOption 把积分三项转成条形图 series', () => {
  const opt = buildCreditChartOption(statsSummary)
  const series = opt.series as { data: number[] }[]
  expect(series[0].data).toEqual([100, 40, 60]) // 发放/已用/剩余
  const xAxis = opt.xAxis as { data: string[] }
  expect(xAxis.data).toEqual(['今日发放', '今日已用', '今日剩余'])
})

// ---- 页面渲染 ----

test('Stats 渲染汇总卡：总积分/今日发放/已用/剩余', async () => {
  render(<Stats />)
  await waitFor(() => expect(screen.getByText('当前总积分')).toBeInTheDocument())
  expect(screen.getByText('今日发放额度')).toBeInTheDocument()
  expect(screen.getByText('今日已用额度')).toBeInTheDocument()
  expect(screen.getByText('今日剩余额度')).toBeInTheDocument()
  // fixture: current=1550, alloc=100, consumed=40, remaining=60（fmt 千分位）
  expect(screen.getByText((1550).toLocaleString('zh-CN'))).toBeInTheDocument()
})

test('Stats 渲染账号状态分布与签到统计', async () => {
  render(<Stats />)
  await waitFor(() => expect(screen.getByText('账号状态分布')).toBeInTheDocument())
  // fixture: checkin_done=1 pending=3
  expect(screen.getByText(/已签到/)).toBeInTheDocument()
})

test('Stats 渲染自动化今日运行统计', async () => {
  render(<Stats />)
  await waitFor(() => expect(screen.getByText('自动化任务今日运行')).toBeInTheDocument())
  // fixture: ok=4 failed=1
  expect(screen.getByText(/成功 4/)).toBeInTheDocument()
  expect(screen.getByText(/失败 1/)).toBeInTheDocument()
})

test('Stats 渲染历史趋势降级标注（M5 后补）', async () => {
  render(<Stats />)
  await waitFor(() => expect(screen.getByText(/历史趋势待 M5 日志系统后补/)).toBeInTheDocument())
})

test('Stats 上游错误时渲染 ErrorView', async () => {
  server.use(http.get('/api/stats/summary', () => HttpResponse.json({ error: '统计不可用' }, { status: 502 })))
  render(<Stats />)
  await waitFor(() => expect(screen.getByText('统计不可用')).toBeInTheDocument())
})

// ---- 轮询清理 ----

test('Stats 卸载时 dispose echarts 实例并 clear interval', async () => {
  const dispose = vi.fn()
  const setOption = vi.fn()
  const resize = vi.fn()
  const initSpy = vi.fn(() => ({ setOption, dispose, resize }))
  const clearIntervalSpy = vi.spyOn(window, 'clearInterval')

  const { unmount } = render(<Stats echartsInit={initSpy} pollIntervalMs={5000} />)
  await waitFor(() => expect(screen.getByText('当前总积分')).toBeInTheDocument())
  expect(initSpy).toHaveBeenCalled()
  unmount()
  expect(dispose).toHaveBeenCalled()
  expect(clearIntervalSpy).toHaveBeenCalled()
  clearIntervalSpy.mockRestore()
})

test('Stats 轮询间隔到期后重新拉取数据', async () => {
  vi.useFakeTimers({ shouldAdvanceTime: true })
  let calls = 0
  server.use(http.get('/api/stats/summary', () => {
    calls++
    return HttpResponse.json(statsSummary)
  }))
  const initSpy = vi.fn(() => ({ setOption: vi.fn(), dispose: vi.fn(), resize: vi.fn() }))
  render(<Stats echartsInit={initSpy} pollIntervalMs={1000} />)
  await waitFor(() => expect(screen.getByText('当前总积分')).toBeInTheDocument())
  const before = calls
  await vi.advanceTimersByTimeAsync(1100)
  expect(calls).toBeGreaterThan(before)
  vi.useRealTimers()
})
