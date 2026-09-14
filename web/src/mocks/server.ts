// MSW node server：Vitest 环境使用（jsdom 下拦截全局 fetch）。
import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)
