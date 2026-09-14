// 运行日志页（M5）：面板后端环形缓冲快照，级别过滤 + 10s 自动刷新（卸载清理）。
import { useCallback, useEffect, useState } from 'react'
import type { LogLine } from '../types'
import { getLogs } from '../services/api'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { EmptyState } from '../components/EmptyState'

/** 自动刷新间隔 10s（M5 任务书定值）。 */
const REFRESH_INTERVAL_MS = 10_000

type LevelFilter = 'all' | 'info' | 'error'

const LEVEL_BADGE: Record<LogLine['level'], string> = {
  info: 'badge',
  warn: 'badge warn',
  error: 'badge bad',
}

export function Logs() {
  const [lines, setLines] = useState<LogLine[] | null>(null)
  const [error, setError] = useState('')
  const [filter, setFilter] = useState<LevelFilter>('all')

  const load = useCallback(
    () => void getLogs(100).then(r => setLines(r.lines)).catch(e => setError(e instanceof Error ? e.message : '日志查询失败')),
    [],
  )
  // 首屏加载 + 10s 轮询；卸载时 clear interval（验收项）。
  useEffect(() => {
    load()
    const timer = window.setInterval(load, REFRESH_INTERVAL_MS)
    return () => window.clearInterval(timer)
  }, [load])

  if (error) return <ErrorView text={error} />
  if (!lines) return <Loading />

  // warn 级在 all 下展示；过滤档只提供 all/info/error（warn 归入非 info，与任务书过滤档位一致）。
  const visible = lines.filter(l => filter === 'all' || l.level === filter)

  return (
    <>
      <Head eyebrow="LOGS" title="运行日志" text="面板运行事件环形缓冲（最近 100 条，新→旧）：登录、OAuth、探测、一键任务与自动化执行记录。敏感字段已在写入前脱敏。" />
      <section className="card">
        <div className="filter-bar">
          {(['all', 'info', 'error'] as const).map(f => (
            <button key={f} className={filter === f ? 'black' : 'link'} onClick={() => setFilter(f)}>{f}</button>
          ))}
          <span className="muted">每 10 秒自动刷新</span>
        </div>
        {visible.length === 0 && <EmptyState text="暂无运行日志。" />}
        {visible.length > 0 && (
          <table>
            <thead>
              <tr><th>时间</th><th>级别</th><th>动作</th><th>消息</th></tr>
            </thead>
            <tbody>
              {visible.map((l, i) => (
                <tr key={`${l.at}-${i}`}>
                  <td className="muted">{l.at}</td>
                  <td><span className={LEVEL_BADGE[l.level] ?? 'badge'}>{l.level}</span></td>
                  <td>{l.action}</td>
                  <td>{l.message}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </>
  )
}
