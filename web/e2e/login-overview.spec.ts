// E2E 主链路 1：登录 → 概览。
// 打开登录页 → 输入 admin/workbuddy → 进入概览 → 验证统计卡渲染。
// 注意：mock 会话标志是 dev server 模块级状态，跨页面共享且 reload 后持久；
// 登录用例必须串行执行（不能与其它 spec 并行），用例结束恢复默认已登录态。
import { expect, test, type Page } from '@playwright/test'

// 打开页面并把 MSW 会话切到未登录态，等待登录页渲染。
// 会话态存 localStorage（见 mocks/handlers.ts），reload 后仍生效。
async function forceLoggedOut(page: Page) {
  await page.goto('/')
  await page.waitForFunction(() => typeof window.__mswSetAuthed === 'function')
  await page.evaluate(() => window.__mswSetAuthed?.(false))
  await page.reload()
  await expect(page.getByRole('button', { name: '进入控制台' })).toBeVisible({ timeout: 15_000 })
}

test('登录后进入概览并渲染统计卡', async ({ page }) => {
  await forceLoggedOut(page)
  await expect(page.getByRole('heading', { name: '账号工作台' })).toBeVisible()

  await page.getByLabel('密码').fill('workbuddy')
  await page.getByRole('button', { name: '进入控制台' }).click()

  // 概览五张统计卡与侧栏导航可见。
  await expect(page.getByText('健康账号')).toBeVisible()
  await expect(page.getByText('账号总数')).toBeVisible()
  await expect(page.getByText('已过期')).toBeVisible()
  await expect(page.getByText('会话失效账号')).toBeVisible()
  await expect(page.getByText('已启用自动化', { exact: true })).toBeVisible()
  await expect(page.getByRole('button', { name: '账号管理' })).toBeVisible()
})

test('错误密码停留在登录页并显示错误', async ({ page }) => {
  await forceLoggedOut(page)
  await page.getByLabel('密码').fill('wrong-password')
  await page.getByRole('button', { name: '进入控制台' }).click()
  await expect(page.getByText('用户名或密码错误')).toBeVisible()
  // 未进入控制台：概览统计卡不可见。
  await expect(page.getByText('健康账号')).not.toBeVisible()
})

// 恢复默认已登录态：serial 模式下后续用例与其它 spec 的复用页面不受影响。
test.afterEach(async ({ page }) => {
  await page.evaluate(() => window.__mswSetAuthed?.(true)).catch(() => {})
})
