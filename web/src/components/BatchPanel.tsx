// 一键任务面板（M3）：四个批量操作按钮 + 逐账号结果表。
// 语义对照契约 §2.3 与 PR #60：禁用跳过、12153 会话失效标记、单账号失败不中断。
// 自动化所有权：面板默认只开新活动探测，接管网关任务需部署层确认——此处只做提示，不自动接管。
import { useCallback, useState } from 'react'
import type { BatchActionInput, BatchActionResult } from '../types'
import { runBatchAction } from '../services/api'
import { Badge } from './Badge'

const ACTIONS: { action: BatchActionInput['action']; label: string }[] = [
  { action: 'checkin', label: '批量签到' },
  { action: 'travel', label: '批量旅行巡检' },
  { action: 'refresh', label: '批量刷新凭据' },
]

function resultBadge(r: BatchActionResult) {
  if (r.skipped) return <Badge>已跳过</Badge>
  if (r.session_dead) return <Badge ok={false}>会话失效</Badge>
  if (r.ok) return <Badge>成功</Badge>
  return <Badge ok={false}>失败</Badge>
}

export function BatchPanel() {
  const [results, setResults] = useState<BatchActionResult[] | null>(null)
  const [pending, setPending] = useState<BatchActionInput['action'] | null>(null)
  const [error, setError] = useState('')

  const run = useCallback(async (action: BatchActionInput['action']) => {
    setPending(action)
    setError('')
    setResults(null)
    try {
      const res = await runBatchAction({ action })
      setResults(res.results)
    } catch (e) {
      setError(e instanceof Error ? e.message : '批量任务执行失败')
    } finally {
      setPending(null)
    }
  }, [])

  return (
    <section className="card">
      <h2>一键任务</h2>
      <p className="muted">
        对面板可见的全部账号批量执行动作；禁用账号自动跳过，12153 会话失效单独标记。
        自动化任务默认只开新活动探测，接管网关任务需部署层确认。
      </p>
      <div className="actions">
        {ACTIONS.map(({ action, label }) => (
          <button
            key={action}
            className="black"
            disabled={pending !== null}
            onClick={() => void run(action)}
          >
            {pending === action ? `${label}中…` : label}
          </button>
        ))}
      </div>
      {error && <div className="error">{error}</div>}
      {results !== null && (
        results.length === 0 ? (
          <p className="muted">没有可执行的账号。</p>
        ) : (
          <table className="table">
            <thead>
              <tr><th>账号</th><th>状态</th><th>说明</th></tr>
            </thead>
            <tbody>
              {results.map(r => (
                <tr key={r.uid}>
                  <td>{r.nickname || r.uid}</td>
                  <td>{resultBadge(r)}</td>
                  <td>{r.error || r.message || ''}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )
      )}
    </section>
  )
}
