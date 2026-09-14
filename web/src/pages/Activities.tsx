// 活动管理：直接读取成长任务，新发现任务以徽标标出。
import { useCallback, useEffect, useState } from 'react'
import type { Any } from '../types'
import { get } from '../services/api'
import { fmt } from '../utils/format'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { EmptyState } from '../components/EmptyState'
import { Badge } from '../components/Badge'

export function Activities() {
  const [items, setItems] = useState<Any[] | null>(null)
  const load = useCallback(() => void get<{ items: Any[] }>('/api/activities').then(d => setItems(d.items)), [])
  useEffect(load, [load])
  if (!items) return <Loading />
  return (
    <>
      <Head eyebrow="ACTIVITIES" title="活动管理" text="直接读取成长任务，自动化探测以任务编号去重发现新增任务。" action={<button className="black" onClick={load}>立即探测</button>} />
      <section className="card table">
        <table>
          <thead><tr><th>账号</th><th>活动任务</th><th>状态</th><th>奖励</th><th></th></tr></thead>
          <tbody>
            {items.map((x, i) => (
              <tr key={`${String(x.uid)}-${String(x.code)}-${i}`}>
                <td>{String(x.nickname || x.uid)}</td>
                <td>{String(x.name || x.code)}</td>
                <td>{String(x.status || '—')}</td>
                <td>{fmt(x.reward)}</td>
                <td>{Boolean(x.new) && <Badge ok>新发现</Badge>}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {items.length === 0 && <EmptyState text="当前没有可读取的成长任务。" />}
      </section>
    </>
  )
}
