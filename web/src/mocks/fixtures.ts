// MSW 确定性 fixtures：dev 与 test 双环境共用同一份假数据。
// 字段名与 docs/api-contract.md 逐字对齐；所有数据均为虚构，不含真实凭据。
import type {
  Account,
  ActivityItem,
  Automation,
  BatchActionResponse,
  CreditSummary,
  LogsResponse,
  LotteryResponse,
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

// 模型样本覆盖三类：有价格+限免徽标+图片支持 / 有价格无限免 / 无价格数据降级 / 上游失败。
export const models: ModelItem[] = [
  {
    uid: 'u-1001',
    nickname: '阿明',
    id: 'wb-pro',
    name: 'WorkBuddy Pro',
    credits: 'x0.51 credits',
    supports_images: true,
    description_zh: '旗舰对话模型',
    description_en: 'Flagship chat model',
    badges: [],
    raw: { id: 'wb-pro', name: 'WorkBuddy Pro', credits: 'x0.51 credits', supportsImages: true },
  },
  {
    uid: 'u-1001',
    nickname: '阿明',
    id: 'wb-lite',
    name: 'WorkBuddy Lite',
    credits: 'x0 credits',
    supports_images: false,
    description_zh: '',
    description_en: '',
    badges: ['限时免费'],
    raw: { id: 'wb-lite', name: 'WorkBuddy Lite', credits: 'x0 credits', tags: ['限时免费'] },
  },
  {
    uid: 'u-1001',
    nickname: '阿明',
    id: 'wb-old',
    name: 'WorkBuddy Old',
    credits: '',
    supports_images: false,
    description_zh: '',
    description_en: '',
    badges: [],
    raw: { id: 'wb-old', name: 'WorkBuddy Old' },
  },
  { uid: 'u-1002', nickname: '阿红', id: '', name: '', credits: '', supports_images: false, description_zh: '', description_en: '', badges: [], error: '上游模型查询失败' },
]

// 活动任务样本覆盖任务类型分组：daily（重复性） / once（单次）。
export const activities: ActivityItem[] = [
  { uid: 'u-1001', nickname: '阿明', code: 'T-DAILY', name: '每日签到', status: '可领取', reward: 100, task_type: 'daily', new: true },
  { uid: 'u-1001', nickname: '阿明', code: 'T-ONCE', name: '首次领养猫猫', status: '进行中', reward: 300, task_type: 'once', new: false },
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
    { at: '2026-09-14T08:30:00+08:00', level: 'info', action: 'run', message: '账号 u-1001 签到成功' },
    { at: '2026-09-14T08:29:10+08:00', level: 'warn', action: 'probe', message: '账号 u-1002 上游会话失效（12153）' },
    { at: '2026-09-14T08:28:00+08:00', level: 'error', action: 'oauth', message: 'OAuth 轮询失败（region=cn）: 上游 503' },
    { at: '2026-09-14T08:27:00+08:00', level: 'info', action: 'login', message: '登录成功（来源 127.0.0.1）' },
    { at: '2026-09-14T08:26:00+08:00', level: 'info', action: 'batch', message: '一键任务 checkin：4 账号，失败 1' },
  ],
}

export const modelPricing: ModelPricingResponse = {
  items: [
    { uid: 'u-1001', nickname: '阿明', id: 'wb-pro', name: 'WorkBuddy Pro', credits: 'x0.51 credits', free: false, badges: [], supports_images: true, note: '' },
    { uid: 'u-1001', nickname: '阿明', id: 'wb-lite', name: 'WorkBuddy Lite', credits: 'x0 credits', free: true, badges: ['限时免费'], supports_images: false, note: '' },
    { uid: 'u-1001', nickname: '阿明', id: 'wb-old', name: 'WorkBuddy Old', credits: '', free: false, badges: [], supports_images: false, note: '该模型上游未提供价格与限免数据' },
  ],
  pricing_available: true,
}

export const lottery: LotteryResponse = {
  items: [
    {
      uid: 'u-1001',
      nickname: '阿明',
      chances: 2,
      draws_total: 5,
      recent: [{ prize: '积分 +10', at: '2026-09-10 12:00' }],
      rewards_total: 1,
      note: '上游字段未核实的占位',
    },
    { uid: 'u-1002', nickname: '阿红', chances: 0, draws_total: 0, recent: [], rewards_total: 0, note: '', error: '上游抽奖概要查询失败' },
  ],
}
