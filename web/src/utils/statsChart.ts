// 统计页图表 option 构造：纯函数，从 StatsSummary 响应映射到 echarts option。
// 与设计 token 对齐：薄荷绿主色，不用 echarts 默认调色板。
import type { EChartsOption } from 'echarts'
import type { StatsSummary } from '../types'

const MINT = '#54c7ad'
const INK = '#101010'

const baseAxis = {
  axisLine: { lineStyle: { color: '#e9e9e6' } },
  axisLabel: { color: '#707070' },
}

/** 账号状态分布条形图：健康/过期/会话死/禁用。 */
export function buildStatusChartOption(s: StatsSummary): EChartsOption {
  return {
    grid: { left: 8, right: 16, top: 8, bottom: 8, containLabel: true },
    xAxis: { type: 'category', data: ['健康', '已过期', '会话失效', '已禁用'], ...baseAxis },
    yAxis: { type: 'value', minInterval: 1, ...baseAxis, splitLine: { lineStyle: { color: '#e9e9e6' } } },
    series: [
      {
        type: 'bar',
        data: [s.status_healthy, s.status_expired, s.status_session_dead, s.status_disabled],
        itemStyle: {
          borderRadius: [6, 6, 0, 0],
          color: (p: { dataIndex: number }) => [MINT, '#f0a35e', '#e05d5d', '#b0b0b0'][p.dataIndex],
        },
        barMaxWidth: 42,
      },
    ],
    textStyle: { color: INK },
  }
}

/** 今日积分额度条形图：发放/已用/剩余。 */
export function buildCreditChartOption(s: StatsSummary): EChartsOption {
  return {
    grid: { left: 8, right: 16, top: 8, bottom: 8, containLabel: true },
    xAxis: { type: 'category', data: ['今日发放', '今日已用', '今日剩余'], ...baseAxis },
    yAxis: { type: 'value', ...baseAxis, splitLine: { lineStyle: { color: '#e9e9e6' } } },
    series: [
      {
        type: 'bar',
        data: [s.credits_today_allocated_total, s.credits_today_consumed_total, s.credits_today_remaining_total],
        itemStyle: { borderRadius: [6, 6, 0, 0], color: MINT },
        barMaxWidth: 42,
      },
    ],
    textStyle: { color: INK },
  }
}
