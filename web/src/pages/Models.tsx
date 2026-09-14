// 模型管理：按账号直连上游查询可用模型。
import { useCallback, useEffect, useState } from 'react'
import type { Any } from '../types'
import { get } from '../services/api'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { Badge } from '../components/Badge'

export function Models() {
  const [items, setItems] = useState<Any[] | null>(null)
  const [error, setError] = useState('')
  const load = useCallback(() => void get<{ items: Any[] }>('/api/models').then(d => setItems(d.items)).catch(e => setError(e.message)), [])
  useEffect(load, [load])
  if (error) return <ErrorView text={error} />
  if (!items) return <Loading />
  return (
    <>
      <Head eyebrow="MODELS" title="模型管理" text="按账号直连上游查询可用模型，不依赖网关模型列表。" action={<button className="black" onClick={load}>刷新模型</button>} />
      <section className="card table">
        <table>
          <thead><tr><th>账号</th><th>模型 ID</th><th>名称</th><th>状态</th></tr></thead>
          <tbody>
            {items.map((m, i) => (
              <tr key={`${String(m.uid)}-${i}`}>
                <td>{String(m.nickname || m.uid)}</td>
                <td className="mono">{String(m.id || '—')}</td>
                <td>{String(m.name || '—')}</td>
                <td>{m.error ? <Badge ok={false}>{String(m.error)}</Badge> : <Badge ok>来自上游</Badge>}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
    </>
  )
}
