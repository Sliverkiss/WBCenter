// 概览页：账号池状态、今日积分与已启用自动化。
import { useCallback, useEffect, useState } from 'react'
import type { OverviewData, Page } from '../types'
import { getOverview } from '../services/api'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { StatCard } from '../components/StatCard'

export function Overview({ onPage }: { onPage: (p: Page) => void }) {
  const [data, setData] = useState<OverviewData | null>(null)
  const [error, setError] = useState('')
  const load = useCallback(() => void getOverview().then(setData).catch(e => setError(e.message)), [])
  useEffect(load, [load])
  if (error) return <ErrorView text={error} />
  if (!data) return <Loading />
  const cards: [string, number][] = [['健康账号', data.healthy], ['账号总数', data.accounts], ['已过期', data.expired], ['已启用自动化', data.automations]]
  return (
    <>
      <Head eyebrow="LOCAL CONTROL ROOM" title="账号工作台" text="账号池状态、今日积分与已启用自动化" action={<button className="black" onClick={load}>刷新状态</button>} />
      <div className="stats">{cards.map(([label, value]) => <StatCard key={label} label={label} value={value} />)}</div>
      <div className="grid two">
        <section className="card">
          <h2>快速入口</h2>
          <div className="quick">
            <button onClick={() => onPage('credits')}>查看今日积分汇总 →</button>
            <button onClick={() => onPage('automation')}>管理本地自动化 →</button>
            <button onClick={() => onPage('scheduler')}>查看账号云端定时任务 →</button>
            <button onClick={() => onPage('activities')}>探测新活动任务 →</button>
          </div>
        </section>
        <section className="card">
          <h2>运行边界</h2>
          <p className="muted">此面板不会读取或改写 workbuddy2api 的配置，也不会调用 Docker 或重启网关。凭据仅保存在共享的 auths 目录。</p>
        </section>
      </div>
    </>
  )
}
