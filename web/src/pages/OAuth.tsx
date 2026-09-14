// 添加账号：选择区域发起 OAuth，轮询授权结果。
import { useEffect, useState } from 'react'
import { post } from '../services/api'
import { Head } from '../components/Head'

export function OAuth({ onNotice, onDone }: { onNotice: (s: string) => void; onDone: () => void }) {
  const [region, setRegion] = useState('cn')
  const [flow, setFlow] = useState<{ id: string; url: string } | null>(null)
  const [starting, setStarting] = useState(false)
  const start = async () => {
    setStarting(true)
    try {
      setFlow(await post('/api/oauth/start', { region }) as { id: string; url: string })
    } catch (e) {
      onNotice(e instanceof Error ? e.message : '发起失败')
    } finally {
      setStarting(false)
    }
  }
  useEffect(() => {
    if (!flow) return
    const timer = window.setInterval(() => void post<{ status: string; nickname?: string }>(`/api/oauth/${encodeURIComponent(flow.id)}/poll`)
      .then(r => {
        if (r.status === 'success') {
          window.clearInterval(timer)
          onNotice(`账号 ${r.nickname || ''} 已保存`)
          onDone()
        }
      })
      .catch(() => undefined), 3000)
    return () => window.clearInterval(timer)
  }, [flow, onDone, onNotice])
  return (
    <>
      <Head eyebrow="OAUTH" title="添加账号" text="选择账号区域后完成 WorkBuddy 网页授权；凭据只写入共享 auths 目录。" />
      <section className="card oauth">
        <label>账号区域
          <select value={region} onChange={e => setRegion(e.target.value)}>
            <option value="cn">中国大陆</option>
            <option value="global">国际版</option>
          </select>
        </label>
        <button className="black" disabled={starting} onClick={() => void start()}>{starting ? '正在获取授权链接…' : '发起 OAuth 授权'}</button>
        {flow && (
          <div className="oauth-link">
            <p>请在新窗口完成授权，控制台会自动轮询结果。</p>
            <a className="black" href={flow.url} target="_blank" rel="noreferrer">打开 WorkBuddy 授权页</a>
          </div>
        )}
      </section>
    </>
  )
}
