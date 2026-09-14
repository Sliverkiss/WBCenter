// 组件原语测试：Badge/EmptyState/Loading/ErrorView/Head/Notice/StatCard。
import { render, screen, fireEvent } from '@testing-library/react'
import { Badge } from '../components/Badge'
import { EmptyState } from '../components/EmptyState'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { Head } from '../components/Head'
import { Notice } from '../components/Notice'
import { StatCard } from '../components/StatCard'

test('Badge 正常态渲染 badge 类', () => {
  render(<Badge ok>正常</Badge>)
  const el = screen.getByText('正常')
  expect(el).toHaveClass('badge')
  expect(el).not.toHaveClass('bad')
})

test('Badge ok=false 时加 bad 类（边界：坏状态）', () => {
  render(<Badge ok={false}>已过期</Badge>)
  expect(screen.getByText('已过期')).toHaveClass('badge', 'bad')
})

test('EmptyState 渲染提示文本', () => {
  render(<EmptyState text="尚无账号" />)
  expect(screen.getByText('尚无账号')).toHaveClass('empty')
})

test('Loading 渲染上游加载提示', () => {
  render(<Loading />)
  expect(screen.getByText(/正在直接查询 WorkBuddy 上游/)).toBeInTheDocument()
})

test('ErrorView 渲染错误文本并带 error card 类', () => {
  render(<ErrorView text="网络错误" />)
  expect(screen.getByText('网络错误')).toHaveClass('error', 'card')
})

test('Head 渲染 eyebrow/title/text 与可选 action', () => {
  render(<Head eyebrow="MODELS" title="模型管理" text="按账号直连上游" action={<button>刷新</button>} />)
  expect(screen.getByText('MODELS')).toHaveClass('eyebrow')
  expect(screen.getByRole('heading', { name: '模型管理' })).toBeInTheDocument()
  expect(screen.getByText('按账号直连上游')).toBeInTheDocument()
  expect(screen.getByRole('button', { name: '刷新' })).toBeInTheDocument()
})

test('Head 无 action 时不渲染多余按钮', () => {
  render(<Head eyebrow="E" title="T" text="X" />)
  expect(screen.queryByRole('button')).not.toBeInTheDocument()
})

test('Notice 渲染消息并支持关闭', () => {
  const onClose = vi.fn()
  render(<Notice text="已保存" onClose={onClose} />)
  expect(screen.getByText('已保存')).toBeInTheDocument()
  fireEvent.click(screen.getByRole('button', { name: '×' }))
  expect(onClose).toHaveBeenCalledOnce()
})

test('StatCard 渲染标签与数值', () => {
  render(<StatCard label="健康账号" value={3} />)
  expect(screen.getByText('健康账号')).toBeInTheDocument()
  expect(screen.getByText('3')).toBeInTheDocument()
})
