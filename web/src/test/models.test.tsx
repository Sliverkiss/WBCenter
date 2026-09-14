// M4 前端 RED：模型页升级（价格列/限免徽标/图片支持/降级文案）+ 活动页升级（任务分组/奖励/lottery 概要）。
// 全部走 MSW handlers（fixtures 覆盖三类模型样本：有价格无限免 / 有价格+限免徽标 / 无价格数据）。
import { render, screen, waitFor, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../mocks/server'
import { Models } from '../pages/Models'
import { Activities } from '../pages/Activities'

// ---- 模型页升级 ----

test('Models 渲染价格列（credits 描述串原样透出）', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getAllByText('wb-pro').length).toBeGreaterThan(0))
  expect(screen.getAllByText('x0.51 credits').length).toBeGreaterThan(0)
  expect(screen.getAllByText('x0 credits').length).toBeGreaterThan(0)
})

test('Models 渲染限免徽标（fixture wb-lite badges=限时免费）', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getAllByText('wb-lite').length).toBeGreaterThan(0))
  expect(screen.getAllByText('限时免费').length).toBeGreaterThan(0)
})

test('Models 渲染图片支持标识（wb-pro supports_images=true）', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getAllByText('wb-pro').length).toBeGreaterThan(0))
  // wb-pro 支持图片，其余两个模型不支持
  expect(screen.getAllByText('支持图片').length).toBeGreaterThan(0)
})

test('Models 无价格数据模型显示「上游未提供」降级文案', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getAllByText('wb-old').length).toBeGreaterThan(0))
  expect(screen.getAllByText(/上游未提供/).length).toBeGreaterThan(0)
})

test('Models 价格视图展示 pricing 汇总（credits + 限免徽标 + 账号归属）', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getByText('模型价格与限免')).toBeInTheDocument())
  const pricing = screen.getByTestId('pricing-section')
  expect(within(pricing).getAllByText('x0.51 credits').length).toBeGreaterThan(0)
  expect(within(pricing).getAllByText('阿明').length).toBeGreaterThan(0)
})

test('Models pricing_available=false 时展示降级提示而非空表', async () => {
  server.use(
    http.get('/api/models/pricing', () =>
      HttpResponse.json({ items: [], pricing_available: false, warning: '上游未提供模型价格数据' }),
    ),
  )
  render(<Models />)
  await waitFor(() => expect(screen.getByText('上游未提供模型价格数据')).toBeInTheDocument())
})

test('Models 无价格数据模型在 pricing 视图中带 note 降级文案', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getByTestId('pricing-section')).toBeInTheDocument())
  const pricing = screen.getByTestId('pricing-section')
  expect(within(pricing).getByText('该模型上游未提供价格与限免数据')).toBeInTheDocument()
})

// ---- 活动页升级 ----

test('Activities 按任务类型分组渲染（daily 重复性 / once 单次）', async () => {
  render(<Activities />)
  await waitFor(() => expect(screen.getByText('每日签到')).toBeInTheDocument())
  expect(screen.getByText('首次领养猫猫')).toBeInTheDocument()
  // 分组标题出现
  expect(screen.getByText(/重复性任务/)).toBeInTheDocument()
  expect(screen.getByText(/单次任务/)).toBeInTheDocument()
})

test('Activities 渲染奖励额度（fixture 100/300）', async () => {
  render(<Activities />)
  await waitFor(() => expect(screen.getByText('每日签到')).toBeInTheDocument())
  expect(screen.getByText('100')).toBeInTheDocument()
  expect(screen.getByText('300')).toBeInTheDocument()
})

test('Activities 渲染 lottery 概览（chances/draws_total/rewards_total/recent）', async () => {
  render(<Activities />)
  await waitFor(() => expect(screen.getByText('抽奖概览')).toBeInTheDocument())
  const lot = screen.getByTestId('lottery-section')
  expect(within(lot).getByText('阿明')).toBeInTheDocument()
  expect(within(lot).getByText(/剩余机会/)).toBeInTheDocument()
  expect(within(lot).getByText('积分 +10')).toBeInTheDocument()
})

test('Activities lottery 失败账号渲染降级徽标，不拖垮其他账号', async () => {
  render(<Activities />)
  await waitFor(() => expect(screen.getByTestId('lottery-section')).toBeInTheDocument())
  const lot = screen.getByTestId('lottery-section')
  expect(within(lot).getByText('上游抽奖概要查询失败')).toBeInTheDocument()
  // 健康账号仍渲染
  expect(within(lot).getByText('阿明')).toBeInTheDocument()
})

test('Activities lottery 渲染只读说明（不提供「立即抽奖」按钮）', async () => {
  render(<Activities />)
  await waitFor(() => expect(screen.getByText('抽奖概览')).toBeInTheDocument())
  expect(screen.queryByRole('button', { name: /抽奖/ })).not.toBeInTheDocument()
  expect(screen.getByText(/只读/)).toBeInTheDocument()
})
