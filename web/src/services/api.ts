// fetch 封装：同源凭据、JSON 体、错误信封解析。
// URL 与行为对照 internal/control/server.go 的 22 条路由，逐字保留。
import type { Any } from '../types'

export async function request<T>(url: string, init?: RequestInit): Promise<T> {
  const response = await fetch(url, { credentials: 'same-origin', headers: init?.body ? { 'Content-Type': 'application/json' } : undefined, ...init })
  const body = (await response.json().catch(() => ({}))) as Any
  if (!response.ok) throw new Error(String(body.error || '请求失败'))
  return body as T
}

export const get = <T,>(url: string) => request<T>(url)
export const post = <T,>(url: string, body: unknown = {}) => request<T>(url, { method: 'POST', body: JSON.stringify(body) })
