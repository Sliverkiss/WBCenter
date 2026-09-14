// E2E 主链路 2：账号池。
// 导航到账号管理 → 验证卡片网格渲染 → 点全量探测 → 验证逐账号探测结果。
import { expect, test } from '@playwright/test'

test('账号池卡片网格渲染并完成全量探测', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: '账号管理' }).click()

  // 卡片网格渲染：四个 fixture 账号昵称全部可见。
  await expect(page.getByRole('heading', { name: '账号池' })).toBeVisible()
  for (const name of ['阿明', '阿红', '阿蓝', '阿绿']) {
    await expect(page.getByRole('heading', { name, exact: true })).toBeVisible()
  }
  // 区域与过期徽标：阿蓝是 Global 且已过期。
  await expect(page.getByText('Global')).toBeVisible()
  await expect(page.getByText('已过期').first()).toBeVisible()

  // 全量探测：逐账号结果合并到卡片。
  await page.getByRole('button', { name: '全量探测' }).click()
  await expect(page.getByText(/探测完成：成功 \d+，失败 \d+/)).toBeVisible()
  // 阿明探测成功：显示当前积分；阿蓝 12153 会话失效徽标。
  await expect(page.getByText('会话失效').first()).toBeVisible()
  await expect(page.getByText(/当前积分 · 今日剩余/).first()).toBeVisible()
})
