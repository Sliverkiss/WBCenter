// 页面级错误展示。
export function ErrorView({ text }: { text: string }) {
  return <div className="error card">{text}</div>
}
