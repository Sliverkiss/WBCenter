// 自动化管理：面板本地定时任务的配置与触发。
import { useCallback, useEffect, useState } from 'react'
import type { Any } from '../types'
import { get, post, request } from '../services/api'
import { displayTime } from '../utils/format'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { Badge } from '../components/Badge'
import { BatchPanel } from '../components/BatchPanel'

export function Automation({ onNotice }: { onNotice: (s: string) => void }) {
  const [items, setItems] = useState<Any[] | null>(null)
  const [runs, setRuns] = useState<Any[]>([])
  const [error, setError] = useState('')
  const [intervals, setIntervals] = useState<Record<string, string>>({})
  const load = useCallback(() => void get<{ items: Any[]; runs: Any[] }>('/api/automations')
    .then(d => {
      setItems(d.items)
      setRuns(d.runs)
      setIntervals(Object.fromEntries(d.items.map(a => [String(a.id), String(a.every_minutes)])))
    })
    .catch(e => setError(e.message)), [])
  useEffect(load, [load])
  const every = (a: Any) => Number(intervals[String(a.id)] || a.every_minutes)
  const save = async (a: Any, enabled: boolean) => {
    try {
      await request(`/api/automations/${encodeURIComponent(String(a.id))}`, { method: 'PUT', body: JSON.stringify({ enabled, every_minutes: every(a) }) })
      onNotice(`${String(a.name)}配置已保存`)
      load()
    } catch (e) {
      onNotice(e instanceof Error ? e.message : '保存失败')
    }
  }
  const run = async (a: Any) => {
    try {
      const r = await post<{ message: string }>(`/api/automations/${encodeURIComponent(String(a.id))}/run`)
      onNotice(r.message)
      load()
    } catch (e) {
      onNotice(e instanceof Error ? e.message : '执行失败')
    }
  }
  if (error) return <ErrorView text={error} />
  if (!items) return <Loading />
  return (
    <>
      <Head eyebrow="AUTOMATION" title="自动化管理" text="这些是真实持久化的面板定时任务。默认仅开启新活动探测，避免与网关既有任务双跑。" />
      <BatchPanel />
      <section className="card">
        <div className="automation-list">
          {items.map(a => (
            <div className="automation" key={String(a.id)}>
              <div>
                <h2>{String(a.name)}</h2>
                <p>上次：{displayTime(a.last_run_at)}</p>
                <small>{String(a.last_result || '等待执行')}</small>
              </div>
              <div className="actions">
                <label className="interval">
                  每 <input aria-label={`${String(a.name)}执行间隔`} type="number" min="5" max="10080" value={intervals[String(a.id)] || String(a.every_minutes)} onChange={e => setIntervals(v => ({ ...v, [String(a.id)]: e.target.value }))} /> 分钟
                </label>
                <button onClick={() => void save(a, Boolean(a.enabled))}>保存配置</button>
                <button className={a.enabled ? 'switch on' : 'switch'} aria-label={`切换${String(a.name)}`} onClick={() => void save(a, !a.enabled)}><span /></button>
                <button onClick={() => void run(a)}>立即运行</button>
              </div>
            </div>
          ))}
        </div>
      </section>
      {runs.length > 0 && (
        <section className="card">
          <h2>最近运行</h2>
          {runs.slice(0, 8).map((r, i) => <p className="run" key={i}><Badge ok={Boolean(r.ok)}>{r.ok ? '成功' : '失败'}</Badge> {String(r.action)} · {String(r.message)}</p>)}
        </section>
      )}
    </>
  )
}
