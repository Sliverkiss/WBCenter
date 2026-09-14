// 账号管理：OAuth 添加、凭据刷新与单账号操作。
import { useCallback, useEffect, useState } from 'react'
import type { Any } from '../types'
import { get, post } from '../services/api'
import { Head } from '../components/Head'
import { ErrorView } from '../components/ErrorView'
import { EmptyState } from '../components/EmptyState'
import { Badge } from '../components/Badge'

export function Accounts({ onNotice }: { onNotice: (s: string) => void }) {
  const [items, setItems] = useState<Any[]>([])
  const [error, setError] = useState('')
  const load = useCallback(() => void get<{ accounts: Any[] }>('/api/accounts').then(d => setItems(d.accounts)).catch(e => setError(e.message)), [])
  useEffect(load, [load])
  const action = async (uid: string, a: string) => {
    try {
      const r = await post<{ message: string }>(`/api/accounts/${encodeURIComponent(uid)}/actions/${a}`)
      onNotice(r.message)
      load()
    } catch (e) {
      onNotice(e instanceof Error ? e.message : '执行失败')
    }
  }
  if (error) return <ErrorView text={error} />
  return (
    <>
      <Head eyebrow="ACCOUNTS" title="账号管理" text="OAuth 添加、凭据刷新与单账号操作" action={<button className="black" onClick={load}>刷新</button>} />
      <section className="card table">
        <table>
          <thead><tr><th>账号</th><th>区域</th><th>Token</th><th>操作</th></tr></thead>
          <tbody>
            {items.map(a => (
              <tr key={String(a.uid)}>
                <td><b>{String(a.nickname || a.uid)}</b><small>{String(a.uid)}</small></td>
                <td>{String(a.domain || 'cn')}</td>
                <td><Badge ok={!a.expired}>{a.expired ? '已过期' : a.needs_refresh ? '即将过期' : '正常'}</Badge></td>
                <td>
                  <button onClick={() => void action(String(a.uid), 'refresh')}>刷新凭据</button>
                  <button onClick={() => void action(String(a.uid), 'checkin')}>签到</button>
                  <button onClick={() => void action(String(a.uid), 'travel')}>旅行巡检</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {items.length === 0 && <EmptyState text="尚无账号，请先通过 OAuth 添加账号。" />}
      </section>
    </>
  )
}
