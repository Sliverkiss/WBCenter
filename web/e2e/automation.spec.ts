// E2E 主链路 3：一键任务。
// 导航到自动化 → 点一键任务按钮 → 验证逐账号结果表（含跳过与会话失效标记）。
import { expect, test } from '@playwright/test'

test('一键任务批量签到并呈现逐账号结果表', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: '自动化管理' }).click()

  await expect(page.getByRole('heading', { name: '一键任务' })).toBeVisible()
  await page.getByRole('button', { name: '批量签到' }).click()

  // 逐账号结果表：四个账号各一行，状态徽标覆盖成功/跳过/会话失效/失败。
  const table = page.locator('table.table')
  await expect(table.getByText('阿明')).toBeVisible()
  await expect(table.getByText('阿红')).toBeVisible()
  await expect(table.getByText('阿蓝')).toBeVisible()
  await expect(table.getByText('阿绿')).toBeVisible()
  await expect(table.locator('span.badge', { hasText: '成功' })).toBeVisible()
  await expect(table.locator('span.badge', { hasText: '已跳过' })).toBeVisible()
  await expect(table.locator('span.badge', { hasText: '会话失效' })).toBeVisible()
  await expect(table.getByText('签到成功')).toBeVisible()
  await expect(table.getByText('账号已禁用，已跳过')).toBeVisible()
})
