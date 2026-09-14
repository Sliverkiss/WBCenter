// 积分管理：今日额度发放/已用/剩余（上游免费套餐当日切片，非本地账本）。
import { useCallback, useEffect, useState } from 'react'
import type { Any } from '../types'
import { get } from '../services/api'
import { fmt } from '../utils/format'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'

export function Credits() {
  const [items, setItems] = useState<Any[] | null>(null)
  const [error, setError] = useState('')
  const [detail, setDetail] = useState<Any | null>(null)
  const load = useCallback(() => void get<{ items: Any[] }>('/api/credits').then(d => setItems(d.items)).catch(e => setError(e.message)), [])
  useEffect(load, [load])
  if (error) return <ErrorView text={error} />
  if (!items) return <Loading />
  return (
    <>
      <Head eyebrow="CREDITS" title="积分管理" text="默认展示今日额度发放、已用、剩余；不建立本地积分账本。" action={<button className="black" onClick={load}>查询上游</button>} />
      <section className="card table">
        <table>
          <thead><tr><th>账号</th><th>当前积分</th><th>今日发放额度</th><th>今日已用</th><th>今日剩余</th><th></th></tr></thead>
          <tbody>
            {items.map(x => (
              <tr key={String(x.uid)}>
                <td>{String(x.nickname || x.uid)}</td>
                <td><b>{fmt(x.current)}</b></td>
                <td>{fmt(x.today_allocated)}</td>
                <td>{fmt(x.today_consumed)}</td>
                <td>{fmt(x.today_remaining)}</td>
                <td><button onClick={() => void get<{ items: Any[] }>(`/api/credits/${encodeURIComponent(String(x.uid))}`).then(d => setDetail(d.items[0] || null))}>查看明细</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
      <p className="muted">“今日发放/已用/剩余”来自上游免费套餐的当日切片，并非虚构的积分交易流水。</p>
      {detail && (
        <div className="drawer">
          <button className="close" onClick={() => setDetail(null)}>×</button>
          <h2>{String(detail.nickname)} 的今日额度明细</h2>
          <dl>
            <dt>当前积分</dt><dd>{fmt(detail.current)}</dd>
            <dt>套餐数</dt><dd>{fmt(detail.packages)}</dd>
            <dt>查询时间</dt><dd>{new Date(String(detail.fetched_at)).toLocaleString()}</dd>
          </dl>
          {Boolean(detail.error) && <div className="error">{String(detail.error)}</div>}
        </div>
      )}
    </>
  )
}
