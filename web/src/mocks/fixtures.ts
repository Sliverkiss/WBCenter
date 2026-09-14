// MSW 确定性 fixtures：dev 与 test 双环境共用同一份假数据。
// 字段名与 docs/api-contract.md 逐字对齐；所有数据均为虚构，不含真实凭据。
import type {
  Account,
  ActivityItem,
  Automation,
  BatchActionResponse,
  CreditSummary,
  LogsResponse,
  MockSchedulerTask,
  ModelItem,
  ModelPricingResponse,
  OverviewData,
  ProbeResponse,
  RunRecord,
  SchedulerTask,
  SessionInfo,
  StatsSummary,
} from '../types'

export const sessionInfo: SessionInfo = {
  authenticated: true,
  read_only: false,
  using_default_password: false,
  timezone: 'Asia/Shanghai',
}

// 四类账号样本：健康 / Token 即将过期（7 天预警） / 12153 会话死 / 上游失败。
export const accounts: Account[] = [
  { uid: 'u-1001', nickname: '阿明', domain: 'copilot.tencent.com', expires_at: 1760000000, expired: false, needs_refresh: false },
  { uid: 'u-1002', nickname: '阿红', domain: 'copilot.tencent.com', expires_at: 1758500000, expired: false, needs_refresh: true },
  { uid: 'u-1003', nickname: '阿蓝', domain: 'www.workbuddy.ai', expires_at: 1700000000, expired: true, needs_refresh: true },
  { uid: 'u-1004', nickname: '阿绿', domain: 'copilot.tencent.com', expires_at: 1760000000, expired: false, needs_refresh: false },
]

export const overview: OverviewData = {
  accounts: 4,
  healthy: 3,
  expired: 1,
  automations: 1,
  session_dead: 1,
  warnings: [],
  updated_at: '2026-09-14T08:30:00+08:00',
}

export const credits: CreditSummary[] = [
  {
    uid: 'u-1001',
    nickname: '阿明',
    current: 1200,
    today_allocated: 100,
    today_consumed: 40,
    today_remaining: 60,
    packages: 2,
    fetched_at: '2026-09-14T08:30:00+08:00',
  },
  {
    uid: 'u-1002',
    nickname: '阿红',
    current: 0,
    today_allocated: 0,
    today_consumed: 0,
    today_remaining: 0,
    packages: 0,
    fetched_at: '2026-09-14T08:30:00+08:00',
    error: '上游积分概要查询失败',
  },
]

export const models: ModelItem[] = [
  { uid: 'u-1001', nickname: '阿明', id: 'wb-pro', name: 'WorkBuddy Pro', raw: { id: 'wb-pro', name: 'WorkBuddy Pro' } },
  { uid: 'u-1002', nickname: '阿红', id: '', name: '', error: '上游模型查询失败' },
]

export const activities: ActivityItem[] = [
  { uid: 'u-1001', nickname: '阿明', code: 'T-DAILY', name: '每日签到', status: '可领取', reward: 100, new: true },
]

export const schedulerTasks: SchedulerTask[] = [
  { uid: 'u-1001', nickname: '阿明', id: 'task-1', name: '每日摘要', status: 'active', updated: '2026-09-01 10:00', raw: {} },
  { uid: 'u-1002', nickname: '阿红', id: '', name: '', status: '', error: '上游拒绝授权，请刷新凭据' },
]

export const mockSchedulerTasks: MockSchedulerTask[] = [
  {
    id: 'mock-1726000000000000000',
    account_uid: 'u-1001',
    name: '每日摘要',
    cron: '0 9 * * *',
    prompt: '汇总今日日程',
    enabled: true,
    created_at: '2026-09-14T08:00:00+08:00',
    updated_at: '2026-09-14T08:00:00+08:00',
    nickname: '阿明',
  },
]

export const automations: Automation[] = [
  { id: 'activity-probe', name: '新活动探测', action: 'activity_probe', enabled: true, every_minutes: 30, next_run_at: '2026-09-14T09:00:00+08:00' },
  { id: 'checkin', name: '每日签到', action: 'checkin', enabled: false, every_minutes: 1440 },
  { id: 'travel', name: '旅行巡检', action: 'travel', enabled: false, every_minutes: 120 },
  { id: 'keepalive', name: '账号保活', action: 'refresh', enabled: false, every_minutes: 720 },
]

export const runs: RunRecord[] = [
  { at: '2026-09-14T08:00:00+08:00', action: 'activity_probe', ok: true, message: '已探测 1 条活动任务' },
]

// 探测结果覆盖四类样本：健康（旅行中）、Token 临期（已签到）、12153 会话死、上游失败。
export const probeResponse: ProbeResponse = {
  results: [
    {
      uid: 'u-1001',
      nickname: '阿明',
      ok: true,
      checkin: { checked_in: true, streak_days: 7 },
      credits: { current: 1200, today_remaining: 60 },
      travel: { status: 'traveling', location: '古镇客栈', departed_at: '2026-09-14T06:00:00+08:00', arrive_at: '2026-09-14T18:00:00+08:00' },
    },
    {
      uid: 'u-1002',
      nickname: '阿红',
      ok: true,
      checkin: { checked_in: false, streak_days: 0 },
      credits: { current: 350, today_remaining: 0 },
      travel: { status: 'idle' },
    },
    { uid: 'u-1003', nickname: '阿蓝', ok: false, session_dead: true, error: '上游会话失效（12153）' },
    { uid: 'u-1004', nickname: '阿绿', ok: false, error: '上游积分概要查询失败' },
  ],
}

export const statsSummary: StatsSummary = {
  accounts_total: 4,
  credits_current_total: 1550,
  credits_today_allocated_total: 100,
  credits_today_consumed_total: 40,
  credits_today_remaining_total: 60,
  checkin_done_today: 1,
  checkin_pending_today: 3,
  status_healthy: 3,
  status_expired: 1,
  status_session_dead: 1,
  status_disabled: 1,
  automation_runs_today_ok: 4,
  automation_runs_today_failed: 1,
  generated_at: '2026-09-14T08:30:00+08:00',
}

// 一键批量任务四类结果样本：成功 / 禁用跳过 / 12153 会话死 / 上游失败。
export const batchActionResponse: BatchActionResponse = {
  results: [
    { uid: 'u-1001', nickname: '阿明', ok: true, message: '签到成功' },
    { uid: 'u-1002', nickname: '阿红', ok: true, skipped: true, message: '账号已禁用，已跳过' },
    { uid: 'u-1003', nickname: '阿蓝', ok: false, session_dead: true, error: '上游会话失效（12153）' },
    { uid: 'u-1004', nickname: '阿绿', ok: false, error: '上游签到接口调用失败' },
  ],
}

export const logs: LogsResponse = {
  lines: [
    { at: '2026-09-14T08:30:00+08:00', level: 'info', message: '账号 u-1001 签到成功' },
    { at: '2026-09-14T08:29:10+08:00', level: 'warn', message: '账号 u-1002 上游会话失效（12153）' },
  ],
}

export const modelPricing: ModelPricingResponse = {
  items: [
    { id: 'wb-pro', name: 'WorkBuddy Pro', free: false, price: 9.9, price_unit: 'CNY/天', free_quota: 0, note: '上游字段未核实，M4 实测定稿' },
    { id: 'wb-lite', name: 'WorkBuddy Lite', free: true, price: 0, price_unit: '', free_quota: 100, note: '限免：每日 100 次' },
  ],
}
