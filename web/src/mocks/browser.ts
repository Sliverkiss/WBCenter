// MSW browser worker：dev 环境使用（不进生产 bundle，main.tsx 仅 DEV 时动态引入）。
import { setupWorker } from 'msw/browser'
import { handlers } from './handlers'

export const worker = setupWorker(...handlers)
