// 状态徽标：ok=false 时变为坏状态样式。
import type { ReactNode } from 'react'

export function Badge({ ok, children }: { ok?: boolean; children: ReactNode }) {
  return <span className={ok === false ? 'badge bad' : 'badge'}>{children}</span>
}
