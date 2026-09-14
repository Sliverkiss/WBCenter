// Playwright E2E 配置：通过真实 dev server（vite + MSW）跑关键用户流，不依赖后端。
// 四绿之外单独执行：npm run e2e。CI 模式下禁用并行、失败重试 1 次。
import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  // mock 会话标志是 dev server 模块级状态，跨页面共享：E2E 必须单 worker 串行。
  workers: 1,
  reporter: process.env.CI ? 'line' : 'list',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  // 自动拉起 vite dev server；dev 环境默认启用 MSW（见 src/main.tsx）。
  // 注意用 localhost 而非 127.0.0.1：本环境 vite 默认只绑定 ::1。
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 60_000,
  },
})
