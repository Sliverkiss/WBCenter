// 页头：eyebrow + 大标题 + 描述 + 可选操作区。
import type { ReactNode } from 'react'

export function Head({ eyebrow, title, text, action }: { eyebrow: string; title: string; text: string; action?: ReactNode }) {
  return (
    <header className="head">
      <div>
        <div className="eyebrow">{eyebrow}</div>
        <h1>{title}</h1>
        <p>{text}</p>
      </div>
      {action}
    </header>
  )
}
