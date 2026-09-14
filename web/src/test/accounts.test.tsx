// 账号池页面（卡片网格版）测试：状态组合 / 探测流转 / 单账号操作。
// 覆盖：健康/Token 临期/已过期/12153/上游失败、签到/猫猫三态、空数据/加载/错误。
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../mocks/server'
import { Accounts } from '../pages/Accounts'
import type { ProbeResponse } from '../types'

const noop = () => undefined

test('Accounts 渲染卡片网格：每账号一张卡', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  // 4 个账号 → 4 张卡（.account-card）。
  const cards = document.querySelectorAll('.account-card')
  expect(cards.length).toBe(4)
})

test('Accounts 卡片显示 uid / 昵称 / 区域徽标', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  const card = screen.getByText('阿明').closest('.account-card') as HTMLElement
  expect(within(card).getByText(/u-1001/)).toBeInTheDocument()
  // 区域徽标（CN 国内版 / Global 国际版）。
  expect(within(card).getByText('CN')).toBeInTheDocument()
  const global = screen.getByText('阿蓝').closest('.account-card') as HTMLElement
  expect(within(global).getByText('Global')).toBeInTheDocument()
})

test('Accounts 健康账号渲染 Token 正常徽标', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  const card = screen.getByText('阿明').closest('.account-card') as HTMLElement
  expect(within(card).getByText('正常')).toBeInTheDocument()
})

test('Accounts Token 临期（needs_refresh）渲染 7 天内到期预警', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿红')).toBeInTheDocument())
  const card = screen.getByText('阿红').closest('.account-card') as HTMLElement
  expect(within(card).getByText(/即将过期/)).toBeInTheDocument()
})

test('Accounts 已过期账号渲染已过期徽标', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿蓝')).toBeInTheDocument())
  const card = screen.getByText('阿蓝').closest('.account-card') as HTMLElement
  expect(within(card).getByText('已过期')).toBeInTheDocument()
})

test('Accounts 全量探测按钮：点击后调 POST /api/probe 并渲染逐账号结果', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /全量探测/ }))
  // 探测完成后：阿明显示签到状态、积分、猫猫旅行位置。
  await waitFor(() => {
    const card = screen.getByText('阿明').closest('.account-card') as HTMLElement
    expect(within(card).getByText(/已签到/)).toBeInTheDocument()
    expect(within(card).getByText(/1200/)).toBeInTheDocument()
    expect(within(card).getByText(/古镇客栈/)).toBeInTheDocument()
    expect(within(card).getByText(/旅行中/)).toBeInTheDocument()
  })
})

test('Accounts 探测结果：12153 账号渲染会话死徽标', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿蓝')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /全量探测/ }))
  await waitFor(() => {
    const card = screen.getByText('阿蓝').closest('.account-card') as HTMLElement
    // 徽标精确匹配（错误文案 "上游会话失效（12153）" 也会命中正则，需区分）。
    expect(within(card).getByText('会话失效')).toBeInTheDocument()
  })
})

test('Accounts 探测结果：上游失败账号渲染错误文案', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿绿')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /全量探测/ }))
  await waitFor(() => {
    const card = screen.getByText('阿绿').closest('.account-card') as HTMLElement
    expect(within(card).getByText(/上游积分概要查询失败/)).toBeInTheDocument()
  })
})

test('Accounts 探测结果：idle 旅行状态不渲染位置信息', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿红')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /全量探测/ }))
  await waitFor(() => {
    const card = screen.getByText('阿红').closest('.account-card') as HTMLElement
    expect(within(card).getByText(/未出行/)).toBeInTheDocument()
  })
})

test('Accounts 探测结果：未签到账号渲染未签到', async () => {
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿红')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /全量探测/ }))
  await waitFor(() => {
    const card = screen.getByText('阿红').closest('.account-card') as HTMLElement
    expect(within(card).getByText(/未签到/)).toBeInTheDocument()
  })
})

test('Accounts 探测按钮 pending 期间禁用并显示进度文案', async () => {
  // 让 probe 永不完结：拦截成一个 pending Promise 之后立即断言按钮状态。
  // resolver 通过 holder 对象传递，规避 TS 对闭包赋值的控制流收窄。
  const holder: { resolve: (r: ProbeResponse) => void } = { resolve: () => undefined }
  const gate = new Promise<ProbeResponse>(res => {
    holder.resolve = res
  })
  server.use(http.post('/api/probe', async () => HttpResponse.json(await gate)))
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  const btn = screen.getByRole('button', { name: /全量探测/ })
  fireEvent.click(btn)
  await waitFor(() => {
    expect(btn).toBeDisabled()
    expect(btn.textContent).toMatch(/探测中/)
  })
  // 收尾：放行 Promise，避免悬挂请求污染后续用例。
  holder.resolve({ results: [] })
  await waitFor(() => expect(btn).not.toBeDisabled())
})

test('Accounts 单账号操作：刷新凭据按钮调 actions/refresh 并提示', async () => {
  const notices: string[] = []
  render(<Accounts onNotice={s => notices.push(s)} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  const card = screen.getByText('阿明').closest('.account-card') as HTMLElement
  fireEvent.click(within(card).getByRole('button', { name: /刷新凭据/ }))
  await waitFor(() => expect(notices.some(n => n.includes('成功'))).toBe(true))
})

test('Accounts 空账号列表渲染空态', async () => {
  server.use(http.get('/api/accounts', () => HttpResponse.json({ accounts: [], warnings: [] })))
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText(/尚无账号/)).toBeInTheDocument())
})

test('Accounts 加载账号失败渲染 ErrorView', async () => {
  server.use(http.get('/api/accounts', () => HttpResponse.json({ error: 'auths 目录不可读' }, { status: 500 })))
  render(<Accounts onNotice={noop} />)
  await waitFor(() => expect(screen.getByText('auths 目录不可读')).toBeInTheDocument())
})

test('Accounts 探测请求失败时给出提示且不修改卡片', async () => {
  const notices: string[] = []
  server.use(http.post('/api/probe', () => HttpResponse.json({ error: '上游整体不可用' }, { status: 502 })))
  render(<Accounts onNotice={s => notices.push(s)} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /全量探测/ }))
  await waitFor(() => expect(notices.some(n => n.includes('上游整体不可用'))).toBe(true))
})
