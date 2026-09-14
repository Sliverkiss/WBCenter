// 账号池监控（M2 卡片网格版）：每账号一张卡，聚合 Token 状态、探测结果（签到/积分/猫猫）、
// 12153 会话死徽标；全量探测按钮触发 POST /api/probe，单账号操作局部刷新。
import { useCallback, useEffect, useMemo, useState } from 'react'
import type { Account, ProbeResult } from '../types'
import { accountAction, getAccounts, startProbe } from '../services/api'
import { Head } from '../components/Head'
import { ErrorView } from '../components/ErrorView'
import { EmptyState } from '../components/EmptyState'
import { Loading } from '../components/Loading'
import { Badge } from '../components/Badge'

// 区域徽标：上游 domain 反推（与 internal/upstream/client.go regionOfAccount 对齐）。
function regionLabel(domain: string): string {
  const d = domain.toLowerCase()
  return d.includes('workbuddy.ai') ? 'Global' : 'CN'
}

// Token 徽标三态：已过期 > 即将过期（needs_refresh） > 正常。
function tokenBadge(a: Account) {
  if (a.expired) return <Badge ok={false}>已过期</Badge>
  if (a.needs_refresh) return <Badge ok={false}>即将过期</Badge>
  return <Badge>正常</Badge>
}

// 猫猫旅行徽标：仅当探测结果存在时显示；idle 不显示位置。
function travelBadge(p: ProbeResult | undefined) {
  if (!p || !p.travel) return null
  const t = p.travel
  if (t.status === 'traveling') {
    return (
      <div className="account-row">
        <span className="muted">猫猫</span>
        <span>旅行中{t.location ? ` · ${t.location}` : ''}</span>
      </div>
    )
  }
  if (t.status === 'arrived') {
    return (
      <div className="account-row">
        <span className="muted">猫猫</span>
        <span>已到达{t.location ? ` · ${t.location}` : ''}</span>
      </div>
    )
  }
  return (
    <div className="account-row">
      <span className="muted">猫猫</span>
      <span>未出行</span>
    </div>
  )
}

function AccountCard({
  account,
  probe,
  onAction,
}: {
  account: Account
  probe: ProbeResult | undefined
  onAction: (uid: string, action: 'checkin' | 'travel' | 'refresh') => void
}) {
  return (
    <article className="card account-card">
      <header className="account-head">
        <div>
          <h2>{account.nickname || account.uid}</h2>
          <small className="muted">{account.uid}</small>
        </div>
        <div className="account-badges">
          <span className="badge">{regionLabel(account.domain)}</span>
          {tokenBadge(account)}
          {probe?.session_dead && <Badge ok={false}>会话失效</Badge>}
        </div>
      </header>

      {probe?.ok && probe.credits && (
        <div className="account-credits">
          <div className="account-credit-big">{probe.credits.current}</div>
          <div className="muted">当前积分 · 今日剩余 {probe.credits.today_remaining}</div>
        </div>
      )}

      {probe?.ok && probe.checkin && (
        <div className="account-row">
          <span className="muted">签到</span>
          <span>{probe.checkin.checked_in ? `已签到 · 连续 ${probe.checkin.streak_days} 天` : '未签到'}</span>
        </div>
      )}

      {travelBadge(probe)}

      {probe && !probe.ok && probe.error && (
        <div className="account-error">
          <span className="muted">探测失败</span>
          <span>{probe.error}</span>
        </div>
      )}

      <footer className="account-actions">
        <button onClick={() => onAction(account.uid, 'refresh')}>刷新凭据</button>
        <button onClick={() => onAction(account.uid, 'checkin')}>签到</button>
        <button onClick={() => onAction(account.uid, 'travel')}>旅行巡检</button>
      </footer>
    </article>
  )
}

export function Accounts({ onNotice }: { onNotice: (s: string) => void }) {
  const [items, setItems] = useState<Account[] | null>(null)
  const [error, setError] = useState('')
  const [probing, setProbing] = useState(false)
  // 探测结果按 uid 索引：渲染时按 uid 合并到卡片，与拉取顺序解耦。
  const [probeMap, setProbeMap] = useState<Record<string, ProbeResult>>({})

  const load = useCallback(
    () =>
      void getAccounts()
        .then(d => setItems(d.accounts))
        .catch(e => setError(e instanceof Error ? e.message : '加载失败')),
    [],
  )
  useEffect(load, [load])

  const runProbe = useCallback(async () => {
    if (probing) return
    setProbing(true)
    try {
      const r = await startProbe()
      const next: Record<string, ProbeResult> = {}
      for (const it of r.results) next[it.uid] = it
      setProbeMap(next)
      const ok = r.results.filter(x => x.ok).length
      onNotice(`探测完成：成功 ${ok}，失败 ${r.results.length - ok}`)
    } catch (e) {
      onNotice(e instanceof Error ? e.message : '探测失败')
    } finally {
      setProbing(false)
    }
  }, [probing, onNotice])

  const action = useCallback(
    async (uid: string, a: 'checkin' | 'travel' | 'refresh') => {
      try {
        const r = await accountAction(uid, a)
        onNotice(r.message)
        load()
      } catch (e) {
        onNotice(e instanceof Error ? e.message : '执行失败')
      }
    },
    [load, onNotice],
  )

  const actionBar = useMemo(
    () => (
      <div className="account-toolbar">
        <button className="black" onClick={() => void runProbe()} disabled={probing}>
          {probing ? '探测中…' : '全量探测'}
        </button>
        <button className="black" onClick={load}>刷新</button>
      </div>
    ),
    [runProbe, probing, load],
  )

  if (error) return <ErrorView text={error} />
  if (items === null) return <Loading />
  return (
    <>
      <Head eyebrow="ACCOUNTS" title="账号池" text="凭据状态、签到与猫猫旅行一览" action={actionBar} />
      {items.length === 0 ? (
        <EmptyState text="尚无账号，请先通过 OAuth 添加账号。" />
      ) : (
        <div className="account-grid">
          {items.map(a => (
            <AccountCard key={a.uid} account={a} probe={probeMap[a.uid]} onAction={(uid, act) => void action(uid, act)} />
          ))}
        </div>
      )}
    </>
  )
}
