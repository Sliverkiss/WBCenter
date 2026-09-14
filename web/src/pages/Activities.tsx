// 活动管理（M4 升级）：成长任务按类型分组 + 奖励展示 + 抽奖概览（只读）。
// lottery 数据源为上游 growth 域四只读端点；上游 schema 未核实，零值带占位标注。
import { useCallback, useEffect, useMemo, useState } from 'react'
import type { ActivitiesResponse, ActivityItem, LotteryOverview, LotteryResponse } from '../types'
import { getActivities, getLottery } from '../services/api'
import { fmt } from '../utils/format'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { EmptyState } from '../components/EmptyState'
import { Badge } from '../components/Badge'

/** 任务类型分组标题：daily/weekly/monthly 视为重复性，once/其他为单次。 */
function groupOf(taskType: string): 'repeat' | 'once' {
  const t = taskType.toLowerCase()
  if (t === 'daily' || t === 'weekly' || t === 'monthly' || t === 'repeat' || t === 'recurring') return 'repeat'
  return 'once'
}

export function Activities() {
  const [items, setItems] = useState<ActivityItem[] | null>(null)
  const [lottery, setLottery] = useState<LotteryOverview[] | null>(null)

  const load = useCallback(() => {
    getActivities().then((d: ActivitiesResponse) => setItems(d.items))
    getLottery()
      .then((d: LotteryResponse) => setLottery(d.items))
      .catch(() => setLottery([]))
  }, [])
  useEffect(load, [load])

  const groups = useMemo(() => {
    const repeat: ActivityItem[] = []
    const once: ActivityItem[] = []
    for (const x of items ?? []) {
      if (groupOf(x.task_type) === 'repeat') repeat.push(x)
      else once.push(x)
    }
    return { repeat, once }
  }, [items])

  if (!items) return <Loading />

  const renderRows = (list: ActivityItem[]) =>
    list.map((x, i) => (
      <tr key={`${x.uid}-${x.code}-${i}`}>
        <td>{x.nickname || x.uid}</td>
        <td>{x.name || x.code}</td>
        <td>{x.status || '—'}</td>
        <td>{fmt(x.reward)}</td>
        <td>{x.new && <Badge ok>新发现</Badge>}</td>
      </tr>
    ))

  return (
    <>
      <Head eyebrow="ACTIVITIES" title="活动管理" text="直接读取成长任务，自动化探测以任务编号去重发现新增任务。" action={<button className="black" onClick={load}>立即探测</button>} />
      <section className="card table">
        <h2>重复性任务</h2>
        <table>
          <thead><tr><th>账号</th><th>活动任务</th><th>状态</th><th>奖励</th><th></th></tr></thead>
          <tbody>{renderRows(groups.repeat)}</tbody>
        </table>
        {groups.repeat.length === 0 && <EmptyState text="当前没有重复性任务。" />}
      </section>
      <section className="card table">
        <h2>单次任务</h2>
        <table>
          <thead><tr><th>账号</th><th>活动任务</th><th>状态</th><th>奖励</th><th></th></tr></thead>
          <tbody>{renderRows(groups.once)}</tbody>
        </table>
        {groups.once.length === 0 && items.length === 0 && <EmptyState text="当前没有可读取的成长任务。" />}
        {groups.once.length === 0 && items.length > 0 && <EmptyState text="当前没有单次任务。" />}
      </section>
      <section className="card table" data-testid="lottery-section">
        <h2>抽奖概览</h2>
        <p className="muted">只读视图：抽奖机会与历史来自上游 growth 域；面板不提供「立即抽奖」操作。</p>
        {lottery === null ? (
          <Loading />
        ) : (
          <table>
            <thead><tr><th>账号</th><th>剩余机会</th><th>累计抽奖</th><th>已获奖品</th><th>最近记录</th><th>状态</th></tr></thead>
            <tbody>
              {lottery.map((l, i) => (
                <tr key={`${l.uid}-${i}`}>
                  <td>{l.nickname || l.uid}</td>
                  <td>{l.chances}</td>
                  <td>{l.draws_total}</td>
                  <td>{l.rewards_total}</td>
                  <td>
                    {l.recent.length === 0
                      ? <span className="muted">—</span>
                      : l.recent.map((r, j) => <div key={j}>{r.prize}<span className="muted"> {r.at}</span></div>)}
                  </td>
                  <td>
                    {l.error
                      ? <Badge ok={false}>{l.error}</Badge>
                      : l.note
                        ? <span className="muted">{l.note}</span>
                        : <Badge ok>来自上游</Badge>}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </>
  )
}
