// dev 环境 MSW worker。handlers 里 /api/session 用模块级 authed 标志模拟会话；
// worker.use 的运行时覆盖在 resetHandlers 后即失效，无法跨 reload 维持未登录态。
// 因此 E2E 通过 window.__mswSetAuthed 直接翻转该标志——与 test 环境 resetMockState 同一机制。
import { setupWorker } from 'msw/browser'
import { handlers, setAuthed } from './handlers'

export const worker = setupWorker(...handlers)

// 暴露给 E2E：Playwright 通过 page.evaluate 翻转 mock 会话标志，
// 避免 /api/logout 请求在 worker 就绪前穿透到不存在的后端代理。
declare global {
  interface Window {
    __mswSetAuthed?: (authed: boolean) => void
  }
}
window.__mswSetAuthed = setAuthed
