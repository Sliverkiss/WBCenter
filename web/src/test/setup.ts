// Vitest 全局测试环境：jest-dom 断言扩展 + MSW 拦截全部 HTTP（未命中 handler 即报错，杜绝裸 fetch 漏网）。
import '@testing-library/jest-dom/vitest'
import { afterAll, afterEach, beforeAll, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import { server } from '../mocks/server'
import { resetMockState } from '../mocks/handlers'

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))

afterEach(() => {
  cleanup()
  server.resetHandlers()
  resetMockState()
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
  window.location.hash = ''
  localStorage.clear()
  document.documentElement.removeAttribute('data-theme')
})

afterAll(() => server.close())
