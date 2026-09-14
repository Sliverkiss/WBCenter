// 概览统计卡。
export function StatCard({ label, value }: { label: string; value: unknown }) {
  return (
    <div className="stat">
      <span>{label}</span>
      <b>{String(value)}</b>
    </div>
  )
}
