// 类型守卫单测：正例 + 非法输入负例（对照 docs/api-contract.md 字段契约）。
import {
  isAccount,
  isCreditSummary,
  isProbeResult,
  isSessionInfo,
  isOverviewData,
  isAutomation,
} from '../types'

describe('isAccount', () => {
  const valid = {
    uid: 'u-1001',
    nickname: '阿明',
    domain: 'copilot.tencent.com',
    expires_at: 1760000000,
    expired: false,
    needs_refresh: true,
  }

  test('接受合法账号对象', () => {
    expect(isAccount(valid)).toBe(true)
  })

  test('拒绝缺少必需字段的对象', () => {
    const { uid: _uid, ...noUid } = valid
    expect(isAccount(noUid)).toBe(false)
    const { expires_at: _e, ...noExpires } = valid
    expect(isAccount(noExpires)).toBe(false)
  })

  test('拒绝字段类型错误的对象', () => {
    expect(isAccount({ ...valid, expired: 'no' })).toBe(false)
    expect(isAccount({ ...valid, expires_at: 'soon' })).toBe(false)
  })

  test('拒绝非对象输入', () => {
    expect(isAccount(null)).toBe(false)
    expect(isAccount(undefined)).toBe(false)
    expect(isAccount('u-1001')).toBe(false)
    expect(isAccount(42)).toBe(false)
    expect(isAccount([])).toBe(false)
  })
})

describe('isCreditSummary', () => {
  const valid = {
    uid: 'u-1001',
    nickname: '阿明',
    current: 1200,
    today_allocated: 100,
    today_consumed: 40,
    today_remaining: 60,
    packages: 2,
    fetched_at: '2026-09-14T08:30:00+08:00',
  }

  test('接受无 error 字段的成功对象', () => {
    expect(isCreditSummary(valid)).toBe(true)
  })

  test('接受带 error 字段的降级对象', () => {
    expect(isCreditSummary({ ...valid, error: '上游积分概要查询失败' })).toBe(true)
  })

  test('拒绝 error 非字符串的对象', () => {
    expect(isCreditSummary({ ...valid, error: 502 })).toBe(false)
  })

  test('拒绝数值字段缺失或类型错误的对象', () => {
    const { current: _c, ...noCurrent } = valid
    expect(isCreditSummary(noCurrent)).toBe(false)
    expect(isCreditSummary({ ...valid, today_remaining: '60' })).toBe(false)
  })

  test('拒绝非对象输入', () => {
    expect(isCreditSummary(null)).toBe(false)
    expect(isCreditSummary('credits')).toBe(false)
  })
})

describe('isProbeResult', () => {
  test('接受成功探测结果（含全部可选子对象）', () => {
    expect(
      isProbeResult({
        uid: 'u-1001',
        nickname: '阿明',
        ok: true,
        checkin: { checked_in: true, streak_days: 7 },
        credits: { current: 1200, today_remaining: 60 },
        travel: { status: 'traveling', location: '杭州', departed_at: '2026-09-14T06:00:00+08:00', arrive_at: '2026-09-14T18:00:00+08:00' },
      }),
    ).toBe(true)
  })

  test('接受失败探测结果（仅 error）', () => {
    expect(isProbeResult({ uid: 'u-1002', nickname: '阿红', ok: false, error: '上游会话失效（12153）' })).toBe(true)
  })

  test('拒绝缺 uid/ok 的对象', () => {
    expect(isProbeResult({ nickname: '阿明', ok: true })).toBe(false)
    expect(isProbeResult({ uid: 'u-1', nickname: '阿明' })).toBe(false)
  })

  test('拒绝子对象字段类型错误的对象', () => {
    expect(isProbeResult({ uid: 'u-1', nickname: 'n', ok: true, checkin: { checked_in: 'yes', streak_days: 7 } })).toBe(false)
    expect(isProbeResult({ uid: 'u-1', nickname: 'n', ok: true, travel: { status: 3 } })).toBe(false)
  })

  test('拒绝非对象输入', () => {
    expect(isProbeResult(undefined)).toBe(false)
    expect(isProbeResult([])).toBe(false)
  })
})

describe('isSessionInfo / isOverviewData / isAutomation', () => {
  test('isSessionInfo 正例与负例', () => {
    expect(isSessionInfo({ authenticated: true, read_only: false, using_default_password: false, timezone: 'Asia/Shanghai' })).toBe(true)
    expect(isSessionInfo({ authenticated: 'yes', read_only: false, using_default_password: false, timezone: 'Asia/Shanghai' })).toBe(false)
    expect(isSessionInfo(null)).toBe(false)
  })

  test('isOverviewData 正例与负例', () => {
    expect(
      isOverviewData({ accounts: 5, healthy: 4, expired: 1, automations: 2, warnings: [], updated_at: '2026-09-14T08:30:00+08:00' }),
    ).toBe(true)
    expect(isOverviewData({ accounts: '5', healthy: 4, expired: 1, automations: 2, warnings: [], updated_at: '' })).toBe(false)
    expect(isOverviewData({ accounts: 5, healthy: 4, expired: 1, automations: 2, warnings: 'none', updated_at: '' })).toBe(false)
  })

  test('isAutomation 正例与负例', () => {
    expect(isAutomation({ id: 'checkin', name: '每日签到', action: 'checkin', enabled: false, every_minutes: 1440 })).toBe(true)
    expect(isAutomation({ id: 'checkin', name: '每日签到', action: 'checkin', enabled: false, every_minutes: '1440' })).toBe(false)
    expect(isAutomation(undefined)).toBe(false)
  })
})
