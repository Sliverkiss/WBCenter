// 数字与时间展示工具。

// 数字本地化，非法值显示占位符。
export function fmt(v: unknown): string {
  const n = Number(v)
  return Number.isFinite(n) ? n.toLocaleString('zh-CN', { maximumFractionDigits: 2 }) : '—'
}

// 零值时间（0001-...）视为未执行。
export function displayTime(v: unknown): string {
  const raw = String(v || '')
  return raw && !raw.startsWith('0001-') ? new Date(raw).toLocaleString() : '未执行'
}
