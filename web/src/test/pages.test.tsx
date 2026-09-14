// 页面冒烟测试：渲染 + 关键数据出现在 DOM。全部走 MSW handlers（与 dev 环境同一份 fixtures），
// 错误分支用 server.use 覆盖单个 handler，不再 stubGlobal fetch。
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { server } from '../mocks/server'
import { Overview } from '../pages/Overview'
import { Accounts } from '../pages/Accounts'
import { Credits } from '../pages/Credits'
import { Models } from '../pages/Models'
import { Automation } from '../pages/Automation'
import { Scheduler } from '../pages/Scheduler'
import { Activities } from '../pages/Activities'
import { OAuth } from '../pages/OAuth'
import { Login } from '../pages/Login'

test('Overview 渲染统计卡与快速入口', async () => {
  render(<Overview onPage={() => undefined} />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  expect(screen.getByText('账号总数')).toBeInTheDocument()
  expect(screen.getByText('已过期')).toBeInTheDocument()
  expect(screen.getByText('查看今日积分汇总 →')).toBeInTheDocument()
})

test('Overview 渲染会话死账号统计卡（fixture session_dead=1）', async () => {
  render(<Overview onPage={() => undefined} />)
  await waitFor(() => expect(screen.getByText('会话失效账号')).toBeInTheDocument())
  // fixture overview.session_dead = 1
  const card = screen.getByText('会话失效账号').closest('.stat') as HTMLElement
  expect(card.textContent).toContain('1')
})

test('Overview 快捷入口：去账号池探测跳转到 accounts 页', async () => {
  const pages: string[] = []
  render(<Overview onPage={p => pages.push(p)} />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  fireEvent.click(screen.getByText('全量探测账号池 →'))
  expect(pages).toContain('accounts')
})

test('Overview 上游错误时渲染 ErrorView', async () => {
  server.use(http.get('/api/overview', () => HttpResponse.json({ error: '上游不可用' }, { status: 502 })))
  render(<Overview onPage={() => undefined} />)
  await waitFor(() => expect(screen.getByText('上游不可用')).toBeInTheDocument())
})

test('Accounts 渲染账号行与操作按钮', async () => {
  render(<Accounts onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  expect(screen.getAllByText('刷新凭据').length).toBeGreaterThan(0)
  expect(screen.getAllByText('签到').length).toBeGreaterThan(0)
  expect(screen.getAllByText('旅行巡检').length).toBeGreaterThan(0)
})

test('Accounts 空列表渲染空态', async () => {
  server.use(http.get('/api/accounts', () => HttpResponse.json({ accounts: [], warnings: [] })))
  render(<Accounts onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText(/尚无账号/)).toBeInTheDocument())
})

test('Credits 渲染积分表格与数据', async () => {
  render(<Credits />)
  await waitFor(() => expect(screen.getByText('阿明')).toBeInTheDocument())
  expect(screen.getByText('今日发放额度')).toBeInTheDocument()
})

test('Models 渲染模型列表与状态徽标', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getByText('wb-pro')).toBeInTheDocument())
  expect(screen.getAllByText('来自上游').length).toBeGreaterThan(0)
})

test('Models 上游错误行渲染坏徽标', async () => {
  render(<Models />)
  await waitFor(() => expect(screen.getByText('上游模型查询失败')).toBeInTheDocument())
})

test('Automation 渲染任务列表与开关', async () => {
  render(<Automation onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText('新活动探测')).toBeInTheDocument())
  expect(screen.getByRole('button', { name: '切换新活动探测' })).toBeInTheDocument()
})

test('Automation 渲染最近运行记录', async () => {
  render(<Automation onNotice={() => undefined} />)
  await waitFor(() => expect(screen.getByText((_c, el) => el?.className === 'run' && el.textContent?.includes('已探测 1 条活动任务') === true)).toBeInTheDocument())
  expect(screen.getByText('成功')).toBeInTheDocument()
})

test('Scheduler 渲染 Mock 表单与上游只读表', async () => {
  render(<Scheduler />)
  await waitFor(() => expect(screen.getByDisplayValue('每日摘要')).toBeInTheDocument())
  expect(screen.getByText('创建 Mock 云端任务')).toBeInTheDocument()
  expect(screen.getByText('上游真实任务（只读）')).toBeInTheDocument()
})

test('Activities 渲染成长任务表', async () => {
  render(<Activities />)
  await waitFor(() => expect(screen.getByText('每日签到')).toBeInTheDocument())
  expect(screen.getByText('新发现')).toBeInTheDocument()
})

test('Activities 空列表渲染空态', async () => {
  server.use(http.get('/api/activities', () => HttpResponse.json({ items: [] })))
  render(<Activities />)
  await waitFor(() => expect(screen.getByText(/当前没有可读取的成长任务/)).toBeInTheDocument())
})

test('OAuth 渲染区域选择与发起按钮', async () => {
  render(<OAuth onNotice={() => undefined} onDone={() => undefined} />)
  expect(screen.getByText('中国大陆')).toBeInTheDocument()
  expect(screen.getByRole('button', { name: '发起 OAuth 授权' })).toBeInTheDocument()
})

test('Login 渲染表单并可提交登录', async () => {
  const onDone = vi.fn().mockResolvedValue(undefined)
  render(<Login onDone={onDone} />)
  expect(screen.getByText('账号工作台')).toBeInTheDocument()
  fireEvent.change(screen.getByLabelText(/密码/), { target: { value: 'workbuddy' } })
  fireEvent.click(screen.getByRole('button', { name: '进入控制台' }))
  await waitFor(() => expect(onDone).toHaveBeenCalled())
})

test('Login 失败时渲染错误信息', async () => {
  render(<Login onDone={vi.fn()} />)
  fireEvent.change(screen.getByLabelText(/密码/), { target: { value: 'wrong' } })
  fireEvent.click(screen.getByRole('button', { name: '进入控制台' }))
  await waitFor(() => expect(screen.getByText('用户名或密码错误')).toBeInTheDocument())
})
