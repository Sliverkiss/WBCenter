// 过渡类型：上游响应结构在 M1 契约终稿前先用宽松记录类型。
export type Any = Record<string, unknown>

export type Page =
  | 'overview'
  | 'accounts'
  | 'credits'
  | 'models'
  | 'automation'
  | 'scheduler'
  | 'activities'
  | 'oauth'

export const PAGES: Page[] = [
  'overview',
  'accounts',
  'credits',
  'models',
  'automation',
  'scheduler',
  'activities',
  'oauth',
]
