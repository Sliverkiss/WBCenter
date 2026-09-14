// useEcharts：echarts 按需引入封装。
// - 按需子包（echarts/core + BarChart + GridComponent + CanvasRenderer），不引全量包；
// - 依赖注入 init（测试传桩，生产默认 echarts.init）；
// - 卸载时 dispose + 移除 resize 监听（轮询清理验收项）。
import { useEffect, useRef } from 'react'
import * as echarts from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { EChartsCoreOption } from 'echarts/core'

echarts.use([BarChart, GridComponent, CanvasRenderer])

/** 与 echarts.init 返回形状对齐的最小接口，测试桩只需实现这三个方法。 */
export interface EchartsLike {
  setOption(option: EChartsCoreOption): void
  resize(): void
  dispose(): void
}

export type EchartsInit = (el: HTMLElement) => EchartsLike

const defaultInit: EchartsInit = el => echarts.init(el) as unknown as EchartsLike

/** jsdom 等无 canvas 环境探测：echarts 无法初始化时返回 null（图表区降级为空白容器，页面其余功能不受影响）。 */
function safeInit(el: HTMLElement, init: EchartsInit): EchartsLike | null {
  if (typeof HTMLCanvasElement !== 'undefined') {
    const probe = document.createElement('canvas')
    if (!probe.getContext('2d')) return null
  }
  try {
    return init(el)
  } catch {
    return null
  }
}

/**
 * useEcharts 挂载 echarts 实例到 ref 容器。
 * option 为 null 时不初始化（数据未就绪）；option 首次非 null 时 init + setOption，
 * 之后 option 变化只 setOption。组件卸载时 dispose 实例并移除 resize 监听。
 * init 仅在首次初始化时读取（约定调用方传稳定引用——页面传模块级常量或测试桩）。
 */
export function useEcharts(option: EChartsCoreOption | null, init?: EchartsInit) {
  const ref = useRef<HTMLDivElement | null>(null)
  const chartRef = useRef<EchartsLike | null>(null)

  useEffect(() => {
    if (!option) return
    const el = ref.current
    if (!el) return
    // 首次非 null：init 并挂 resize 监听；cleanup 在 option 回到 null 或卸载时执行。
    let chart = chartRef.current
    if (!chart) {
      // 注入桩时跳过 canvas 环境探测（测试无 canvas 也要走桩以验证生命周期）。
      chart = init ? init(el) : safeInit(el, defaultInit)
      if (!chart) return
      chartRef.current = chart
    }
    chart.setOption(option)
    const onResize = () => chart.resize()
    window.addEventListener('resize', onResize)
    return () => {
      window.removeEventListener('resize', onResize)
      chart.dispose()
      chartRef.current = null
    }
    // init 约定为稳定引用（模块级常量或测试桩），不参与依赖。
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [option])

  return ref
}
