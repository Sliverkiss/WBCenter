// 添加账号（M5 完整流）：选择区域发起 OAuth → 新窗口完成授权 → 面板轮询状态机
// （waiting/success/error/timeout + 5 分钟倒计时）→ 成功展示账号 / 失败可重试。
import { useCallback, useEffect, useRef, useState } from 'react'
import type { OAuthPollResponse, OAuthStartInput } from '../types'
import { pollOAuth, startOAuth } from '../services/api'
import { Head } from '../components/Head'

/** 轮询间隔 3s（与 PR #60 语义一致）；总上限 5 分钟（buddy-oauth.ts:482）。 */
const POLL_INTERVAL_MS = 3_000
const FLOW_TIMEOUT_MS = 5 * 60_000

/**
 * 单次 poll 并返回下一相位（终态）或 null（继续 waiting）。
 * 独立导出便于针对状态机本身做无定时器单测。
 */
export async function pollOnce(id: string): Promise<FlowPhase | null> {
  const res: OAuthPollResponse = await pollOAuth(id)
  if (res.status === 'success') return { kind: 'success', uid: res.uid, nickname: res.nickname }
  if (res.status === 'error') return { kind: 'error', message: res.message }
  if (res.status === 'timeout') return { kind: 'timeout', message: res.message }
  return null
}

export type FlowPhase =
  | { kind: 'idle' }
  | { kind: 'waiting'; id: string; url: string; deadline: number }
  | { kind: 'success'; uid: string; nickname: string }
  | { kind: 'error'; message: string }
  | { kind: 'timeout'; message: string }

const REGION_TEXT: Record<OAuthStartInput['region'], string> = {
  cn: '中国大陆（copilot.tencent.com）：国内版 WorkBuddy 账号',
  global: '国际版（www.workbuddy.ai）：海外版 WorkBuddy 账号',
}

function fmtCountdown(ms: number): string {
  const total = Math.max(0, Math.ceil(ms / 1000))
  const m = Math.floor(total / 60)
  const s = total % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

export function OAuth({ onNotice, onDone }: { onNotice: (s: string) => void; onDone: () => void }) {
  const [region, setRegion] = useState<OAuthStartInput['region']>('cn')
  const [phase, setPhase] = useState<FlowPhase>({ kind: 'idle' })
  const [starting, setStarting] = useState(false)
  const [remainMs, setRemainMs] = useState(FLOW_TIMEOUT_MS)
  const doneRef = useRef(false)

  const start = async () => {
    setStarting(true)
    try {
      const res = await startOAuth({ region })
      doneRef.current = false
      setPhase({ kind: 'waiting', id: res.id, url: res.url, deadline: Date.now() + FLOW_TIMEOUT_MS })
    } catch (e) {
      onNotice(e instanceof Error ? e.message : '发起失败')
    } finally {
      setStarting(false)
    }
  }

  const reset = useCallback(() => setPhase({ kind: 'idle' }), [])

  // 轮询状态机：waiting 期间每 3s 一次 poll；success/error/timeout 终态停表。
  // id/deadline 经依赖数组显式传递，避免闭包捕获过期 phase。
  const waitingId = phase.kind === 'waiting' ? phase.id : null
  const waitingDeadline = phase.kind === 'waiting' ? phase.deadline : null
  useEffect(() => {
    if (waitingId === null || waitingDeadline === null) return
    let active = true
    const tick = async () => {
      if (!active) return
      if (Date.now() >= waitingDeadline) {
        setPhase({ kind: 'timeout', message: '授权超时（5 分钟未完成），请重新发起' })
        return
      }
      setRemainMs(waitingDeadline - Date.now())
      try {
        const next = await pollOnce(waitingId)
        if (!active || next === null) return // null=waiting，继续下一轮
        if (next.kind === 'success' && !doneRef.current) {
          doneRef.current = true
          onNotice(`账号 ${next.nickname || next.uid} 已保存`)
          onDone()
        }
        setPhase(next)
      } catch {
        // 瞬时网络错误不打断轮询（上游 502 等由下一次 poll 重试）。
      }
    }
    void tick()
    const timer = window.setInterval(() => void tick(), POLL_INTERVAL_MS)
    return () => {
      active = false
      window.clearInterval(timer)
    }
  }, [waitingId, waitingDeadline, onDone, onNotice])

  return (
    <>
      <Head eyebrow="OAUTH" title="添加账号" text="选择账号区域后完成 WorkBuddy 网页授权；凭据只写入共享 auths 目录。轮询上限 5 分钟，超时需重新发起。" />
      <section className="card oauth">
        <label>账号区域
          <select value={region} onChange={e => setRegion(e.target.value as OAuthStartInput['region'])} disabled={phase.kind === 'waiting'}>
            <option value="cn">中国大陆</option>
            <option value="global">国际版</option>
          </select>
        </label>
        <p className="muted">{REGION_TEXT[region]}</p>

        {phase.kind === 'idle' && (
          <button className="black" disabled={starting} onClick={() => void start()}>
            {starting ? '正在获取授权链接…' : '发起 OAuth 授权'}
          </button>
        )}

        {phase.kind === 'waiting' && (
          <div className="oauth-link">
            <p>请在新窗口完成授权，控制台会自动轮询结果。</p>
            <a className="black" href={phase.url} target="_blank" rel="noreferrer">打开 WorkBuddy 授权页</a>
            <p className="muted">
              <span className="badge">等待授权</span> 剩余时间 {fmtCountdown(remainMs)}
            </p>
            <button className="link" onClick={reset}>取消并重新发起</button>
          </div>
        )}

        {phase.kind === 'success' && (
          <div className="oauth-link">
            <p><span className="badge">授权成功</span> 账号 {phase.nickname}（{phase.uid}）已加入账号池。</p>
            <button className="link" onClick={reset}>继续添加下一个账号</button>
          </div>
        )}

        {phase.kind === 'error' && (
          <div className="oauth-link">
            <p><span className="badge bad">授权失败</span> {phase.message}</p>
            <button className="black" onClick={reset}>重试</button>
          </div>
        )}

        {phase.kind === 'timeout' && (
          <div className="oauth-link">
            <p><span className="badge bad">授权超时</span> {phase.message}</p>
            <button className="black" onClick={reset}>重试</button>
          </div>
        )}
      </section>
    </>
  )
}
