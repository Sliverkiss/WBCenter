// 主题切换测试：默认浅色、点击切换 data-theme、localStorage 持久化。
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { vi } from 'vitest'
import App from '../App'

function stubAuthed() {
  vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL) => {
    const url = String(input)
    if (url.includes('/api/session')) return new Response(JSON.stringify({ authenticated: true }), { status: 200 })
    if (url.includes('/api/overview')) return new Response(JSON.stringify({ healthy: 0, accounts: 0, expired: 0, automations: 0 }), { status: 200 })
    return new Response('{}', { status: 200 })
  }))
}

test('默认渲染为浅色主题（无 data-theme 或为 light）', async () => {
  stubAuthed()
  render(<App />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  const theme = document.documentElement.getAttribute('data-theme')
  expect(theme === null || theme === 'light').toBe(true)
})

test('点击主题切换按钮后 html data-theme 变为 dark 并写入 localStorage', async () => {
  stubAuthed()
  render(<App />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  fireEvent.click(screen.getByRole('button', { name: /切换主题|暗色|深色/ }))
  await waitFor(() => expect(document.documentElement.getAttribute('data-theme')).toBe('dark'))
  expect(localStorage.getItem('wb-theme')).toBe('dark')
})

test('localStorage 已有 dark 时初始即为暗色', async () => {
  localStorage.setItem('wb-theme', 'dark')
  stubAuthed()
  render(<App />)
  await waitFor(() => expect(screen.getByText('健康账号')).toBeInTheDocument())
  expect(document.documentElement.getAttribute('data-theme')).toBe('dark')
})
