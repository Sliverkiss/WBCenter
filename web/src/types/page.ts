// 页面路由类型：与 API 契约无关，单独成文件避免类型层耦合 UI 结构。
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
