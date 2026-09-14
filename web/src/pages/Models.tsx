// 模型管理（M4 升级）：按账号直连上游查询可用模型 + 模型价格与限免视图（Issue #70）。
// 价格数据源为上游 credits 描述串原样透传；上游无数值单价字段，不编造。
import { useCallback, useEffect, useState } from 'react'
import type { ModelItem, ModelPricingResponse, ModelsResponse } from '../types'
import { getModels, getModelPricing } from '../services/api'
import { Head } from '../components/Head'
import { Loading } from '../components/Loading'
import { ErrorView } from '../components/ErrorView'
import { Badge } from '../components/Badge'

export function Models() {
  const [items, setItems] = useState<ModelItem[] | null>(null)
  const [pricing, setPricing] = useState<ModelPricingResponse | null>(null)
  const [error, setError] = useState('')

  const load = useCallback(() => {
    setError('')
    getModels()
      .then((d: ModelsResponse) => setItems(d.items))
      .catch(e => setError(e instanceof Error ? e.message : '模型查询失败'))
    getModelPricing()
      .then(setPricing)
      .catch(() => setPricing({ items: [], pricing_available: false, warning: '上游未提供模型价格数据' }))
  }, [])
  useEffect(load, [load])

  if (error) return <ErrorView text={error} />
  if (!items) return <Loading />
  return (
    <>
      <Head eyebrow="MODELS" title="模型管理" text="按账号直连上游查询可用模型与价格（credits 描述串透传），不依赖网关模型列表。" action={<button className="black" onClick={load}>刷新模型</button>} />
      <section className="card table">
        <table>
          <thead><tr><th>账号</th><th>模型 ID</th><th>名称</th><th>价格</th><th>能力</th><th>状态</th></tr></thead>
          <tbody>
            {items.map((m, i) => (
              <tr key={`${m.uid}-${m.id}-${i}`}>
                <td>{m.nickname || m.uid}</td>
                <td className="mono">{m.id || '—'}</td>
                <td>{m.name || '—'}</td>
                <td>
                  {m.credits ? <span className="mono">{m.credits}</span> : <span className="muted">上游未提供</span>}
                  {m.badges.map(b => (
                    <Badge key={b} ok>{b}</Badge>
                  ))}
                </td>
                <td>{m.supports_images ? <Badge ok>支持图片</Badge> : <span className="muted">—</span>}</td>
                <td>{m.error ? <Badge ok={false}>{m.error}</Badge> : <Badge ok>来自上游</Badge>}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </section>
      <section className="card table" data-testid="pricing-section">
        <h2>模型价格与限免</h2>
        {pricing === null ? (
          <Loading />
        ) : !pricing.pricing_available ? (
          <p className="muted">{pricing.warning ?? '上游未提供模型价格数据'}</p>
        ) : (
          <table>
            <thead><tr><th>账号</th><th>模型 ID</th><th>名称</th><th>价格（上游透传）</th><th>限免</th><th>备注</th></tr></thead>
            <tbody>
              {pricing.items.map((p, i) => (
                <tr key={`${p.uid}-${p.id}-${i}`}>
                  <td>{p.nickname || p.uid}</td>
                  <td className="mono">{p.id || '—'}</td>
                  <td>{p.name || '—'}</td>
                  <td>{p.credits ? <span className="mono">{p.credits}</span> : <span className="muted">上游未提供</span>}</td>
                  <td>
                    {p.free ? <Badge ok>限免</Badge> : <span className="muted">—</span>}
                    {p.badges.map(b => (
                      <Badge key={b} ok>{b}</Badge>
                    ))}
                  </td>
                  <td className="muted">{p.note || ''}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <p className="muted">价格来源为上游 credits 描述串原样透传；上游未提供数值单价与限免额度字段。</p>
      </section>
    </>
  )
}
