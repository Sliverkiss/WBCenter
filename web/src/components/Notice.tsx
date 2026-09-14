// 顶部通知条，可关闭。
export function Notice({ text, onClose }: { text: string; onClose: () => void }) {
  return (
    <div className="notice">
      {text}
      <button onClick={onClose}>×</button>
    </div>
  )
}
