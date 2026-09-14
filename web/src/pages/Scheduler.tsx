// 账号云端定时任务：上游只读 + 本地 Mock CRUD（上游未授权时的演练骨架）。
import { useCallback, useEffect, useState } from 'react'
import type { Any } from '../types'
import { get, post, request } from '../services/api'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { EmptyState } from '../components/EmptyState'
import { Badge } from '../components/Badge'

export function Scheduler() {
  const [items, setItems] = useState<Any[] | null>(null)
  const [mock, setMock] = useState<Any[] | null>(null)
  const [accounts, setAccounts] = useState<Any[]>([])
  const [note, setNote] = useState('')
  const [detail, setDetail] = useState<Any | null>(null)
  const [error, setError] = useState('')
  const [draft, setDraft] = useState({ account_uid: '', name: '', cron: '0 9 * * *', prompt: '', enabled: true })
  const [edits, setEdits] = useState<Record<string, Any>>({})

  const load = useCallback(() => void Promise.all([
    get<{ items: Any[]; note: string }>('/api/scheduler-tasks'),
    get<{ items: Any[]; note: string }>('/api/mock-scheduler-tasks'),
    get<{ accounts: Any[] }>('/api/accounts'),
  ]).then(([up, local, acc]) => {
    setItems(up.items)
    setMock(local.items)
    setNote(`${up.note} ${local.note}`)
    setAccounts(acc.accounts)
    setDraft(v => v.account_uid ? v : (acc.accounts[0] ? { ...v, account_uid: String(acc.accounts[0].uid) } : v))
    setEdits(Object.fromEntries(local.items.map(t => [String(t.id), t])))
  }).catch(e => setError(e.message)), [])
  useEffect(load, [load])

  const inspect = (x: Any) => void get<{ item: Any }>(`/api/scheduler-tasks/${encodeURIComponent(String(x.uid))}/${encodeURIComponent(String(x.id))}`).then(d => setDetail(d.item)).catch(e => setError(e.message))
  const create = async () => {
    try {
      await post('/api/mock-scheduler-tasks', draft)
      setDraft(v => ({ ...v, name: '', prompt: '' }))
      load()
    } catch (e) {
      setError(e instanceof Error ? e.message : '创建 Mock 任务失败')
    }
  }
  const update = async (id: string) => {
    try {
      const task = edits[id]
      await request(`/api/mock-scheduler-tasks/${encodeURIComponent(id)}`, { method: 'PUT', body: JSON.stringify(task) })
      load()
    } catch (e) {
      setError(e instanceof Error ? e.message : '保存 Mock 任务失败')
    }
  }
  const remove = async (id: string) => {
    try {
      await request(`/api/mock-scheduler-tasks/${encodeURIComponent(id)}`, { method: 'DELETE' })
      load()
    } catch (e) {
      setError(e instanceof Error ? e.message : '删除 Mock 任务失败')
    }
  }

  if (error) return <ErrorView text={error} />
  if (!items || !mock) return <Loading />
  return (
    <>
      <Head eyebrow="ACCOUNT SCHEDULER · MOCK" title="账号云端定时任务" text="上游 scheduler 未授权时，使用本地 Mock 骨架演练任务管理；不会向 WorkBuddy 创建或修改任何任务。" action={<button className="black" onClick={load}>刷新任务</button>} />
      <div className="notice neutral">{note}</div>
      <section className="card mock-form">
        <h2>创建 Mock 云端任务</h2>
        <div className="mock-fields">
          <label>账号
            <select value={draft.account_uid} onChange={e => setDraft(v => ({ ...v, account_uid: e.target.value }))}>
              {accounts.map(a => <option key={String(a.uid)} value={String(a.uid)}>{String(a.nickname || a.uid)}</option>)}
            </select>
          </label>
          <label>任务名称<input value={draft.name} placeholder="例如：每日工作摘要" onChange={e => setDraft(v => ({ ...v, name: e.target.value }))} /></label>
          <label>调度表达式<input value={draft.cron} onChange={e => setDraft(v => ({ ...v, cron: e.target.value }))} /></label>
          <label>任务内容<textarea value={draft.prompt} placeholder="Mock 内容，仅本地保存" onChange={e => setDraft(v => ({ ...v, prompt: e.target.value }))} /></label>
          <label className="check"><input type="checkbox" checked={draft.enabled} onChange={e => setDraft(v => ({ ...v, enabled: e.target.checked }))} /> 启用</label>
        </div>
        <button className="black" disabled={!draft.account_uid || !draft.name.trim()} onClick={() => void create()}>创建 Mock 任务</button>
      </section>
      <section className="card table">
        <h2>本地 Mock 任务</h2>
        <p className="muted">仅用于先完成页面与交互骨架，数据保存在控制台自己的 state.json。</p>
        <table>
          <thead><tr><th>账号</th><th>名称 / 调度</th><th>内容</th><th>状态</th><th>操作</th></tr></thead>
          <tbody>
            {mock.map(t => {
              const e = edits[String(t.id)] || t
              return (
                <tr key={String(t.id)}>
                  <td>{String(t.nickname || t.account_uid)}</td>
                  <td>
                    <input value={String(e.name || '')} onChange={x => setEdits(v => ({ ...v, [String(t.id)]: { ...e, name: x.target.value } }))} />
                    <input value={String(e.cron || '')} onChange={x => setEdits(v => ({ ...v, [String(t.id)]: { ...e, cron: x.target.value } }))} />
                  </td>
                  <td><textarea value={String(e.prompt || '')} onChange={x => setEdits(v => ({ ...v, [String(t.id)]: { ...e, prompt: x.target.value } }))} /></td>
                  <td><label className="check"><input type="checkbox" checked={Boolean(e.enabled)} onChange={x => setEdits(v => ({ ...v, [String(t.id)]: { ...e, enabled: x.target.checked } }))} /> {e.enabled ? '启用' : '停用'}</label></td>
                  <td>
                    <button onClick={() => void update(String(t.id))}>保存</button>
                    <button onClick={() => void remove(String(t.id))}>删除</button>
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {mock.length === 0 && <EmptyState text="尚无 Mock 云端任务。请先创建一条本地骨架任务。" />}
      </section>
      <section className="card table">
        <h2>上游真实任务（只读）</h2>
        <table>
          <thead><tr><th>账号</th><th>任务</th><th>状态</th><th>更新时间</th><th></th></tr></thead>
          <tbody>
            {items.map((x, i) => (
              <tr key={`${String(x.uid)}-${i}`}>
                <td>{String(x.nickname || x.uid)}</td>
                <td>{String(x.name || x.id || '—')}</td>
                <td>{x.error ? <Badge ok={false}>{String(x.error)}</Badge> : <Badge ok={String(x.status).toLowerCase() !== 'failed'}>{String(x.status || '未知')}</Badge>}</td>
                <td>{String(x.updated || '—')}</td>
                <td>{Boolean(x.id) && <button onClick={() => inspect(x)}>查看详情</button>}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {items.length === 0 && <EmptyState text="上游暂未返回任何账号云端定时任务。" />}
      </section>
      {detail && (
        <div className="drawer">
          <button className="close" onClick={() => setDetail(null)}>×</button>
          <h2>账号云端任务详情</h2>
          <pre>{JSON.stringify(detail, null, 2)}</pre>
        </div>
      )}
    </>
  )
}
