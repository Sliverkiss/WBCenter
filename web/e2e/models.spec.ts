// E2E 主链路 4：模型价格。
// 导航到模型管理 → 验证价格列/限免标识 → pricing 视图渲染。
import { expect, test } from '@playwright/test'

test('模型管理展示价格列与限免标识', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: '模型管理' }).click()

  await expect(page.getByRole('heading', { name: '模型管理' })).toBeVisible()

  // pricing 视图（M4）：价格与限免列渲染。
  const pricing = page.getByTestId('pricing-section')
  await expect(pricing.getByRole('heading', { name: '模型价格与限免' })).toBeVisible()
  await expect(pricing.getByText('WorkBuddy Pro')).toBeVisible()
  await expect(pricing.getByText('x0.51 credits')).toBeVisible()
  await expect(pricing.getByText('x0 credits')).toBeVisible()
  // 限免标识：WorkBuddy Lite 限时免费。
  await expect(pricing.locator('span.badge', { hasText: '限免' }).first()).toBeVisible()
  await expect(pricing.locator('span.badge', { hasText: '限时免费' })).toBeVisible()
  // 上游未提供价格时的降级展示。
  await expect(pricing.getByText('该模型上游未提供价格与限免数据')).toBeVisible()
})
