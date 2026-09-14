// Vitest 全局测试环境：jest-dom 断言扩展 + fetch 默认拦截。
import '@testing-library/jest-dom/vitest'
import { afterEach, vi } from 'vitest'
import { cleanup } from '@testing-library/react'

afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
  window.location.hash = ''
  localStorage.clear()
  document.documentElement.removeAttribute('data-theme')
})
