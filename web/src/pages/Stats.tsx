// 统计页（M3）：积分汇总卡 + 账号状态分布图 + 签到与自动化今日统计。
// echarts 按需引入（echarts/core 子包，见 hooks/useEcharts）；轮询刷新，卸载时清理。
import { useCallback, useEffect, useMemo, useState } from 'react'
import type { StatsSummary } from '../types'
import { getStatsSummary } from '../services/api'
import { fmt } from '../utils/format'
import { useEcharts, type EchartsInit } from '../hooks/useEcharts'
import { buildStatusChartOption, buildCreditChartOption } from '../utils/statsChart'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { StatCard } from '../components/StatCard'

/** 默认轮询间隔 30s；测试可注入更短值。 */
const DEFAULT_POLL_MS = 30_000

export function Stats({
  echartsInit,
  pollIntervalMs = DEFAULT_POLL_MS,
}: {
  echartsInit?: EchartsInit
  pollIntervalMs?: number
}) {
  const [data, setData] = useState<StatsSummary | null>(null)
  const [error, setError] = useState('')

  const load = useCallback(
    () => void getStatsSummary().then(setData).catch(e => setError(e instanceof Error ? e.message : '统计查询失败')),
    [],
  )
  // 首屏加载 + 轮询刷新；卸载时 clear interval（验收项）。
  useEffect(() => {
    load()
    const timer = window.setInterval(load, pollIntervalMs)
    return () => window.clearInterval(timer)
  }, [load, pollIntervalMs])

  const statusOption = useMemo(() => (data ? buildStatusChartOption(data) : null), [data])
  const creditOption = useMemo(() => (data ? buildCreditChartOption(data) : null), [data])
  const statusRef = useEcharts(statusOption, echartsInit)
  const creditRef = useEcharts(creditOption, echartsInit)

  if (error) return <ErrorView text={error} />
  if (!data) return <Loading />
  return (
    <>
      <Head eyebrow="STATS" title="统计分析" text="Token 与积分统计聚合：积分与签到计数来自最近探测快照与实时聚合，自动化统计来自本地运行记录。" />
      <section className="grid-stats">
        <StatCard label="当前总积分" value={fmt(data.credits_current_total)} />
        <StatCard label="今日发放额度" value={fmt(data.credits_today_allocated_total)} />
        <StatCard label="今日已用额度" value={fmt(data.credits_today_consumed_total)} />
        <StatCard label="今日剩余额度" value={fmt(data.credits_today_remaining_total)} />
      </section>
      <section className="card">
        <h2>账号状态分布</h2>
        <div ref={statusRef} className="chart" style={{ height: 220 }} />
        <p className="muted">
          账号总数 {data.accounts_total} · 健康 {data.status_healthy} · 已过期 {data.status_expired} · 会话失效 {data.status_session_dead} · 已禁用 {data.status_disabled}
        </p>
      </section>
      <section className="card">
        <h2>今日积分额度</h2>
        <div ref={creditRef} className="chart" style={{ height: 220 }} />
      </section>
      <section className="card">
        <h2>今日签到</h2>
        <p>
          已签到 <b>{data.checkin_done_today}</b> · 未签到 <b>{data.checkin_pending_today}</b>
        </p>
        <h2>自动化任务今日运行</h2>
        <p>
          成功 {data.automation_runs_today_ok} · 失败 {data.automation_runs_today_failed}
        </p>
      </section>
      <p className="muted">历史趋势待 M5 日志系统后补——当前展示为最新快照，每 {Math.round(pollIntervalMs / 1000)} 秒自动刷新。</p>
    </>
  )
}
