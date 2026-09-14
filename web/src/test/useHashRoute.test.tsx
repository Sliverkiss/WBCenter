// useHashRoute：hash 路由 hook，读写 location.hash，默认 overview。
import { act, renderHook } from '@testing-library/react'
import { useHashRoute } from '../hooks/useHashRoute'

test('无 hash 时默认 overview', () => {
  const { result } = renderHook(() => useHashRoute())
  expect(result.current.page).toBe('overview')
})

test('读取现有 hash', () => {
  window.location.hash = '#/accounts'
  const { result } = renderHook(() => useHashRoute())
  expect(result.current.page).toBe('accounts')
})

test('navigate 写入 hash 并触发状态更新', () => {
  const { result } = renderHook(() => useHashRoute())
  act(() => result.current.navigate('credits'))
  expect(window.location.hash).toBe('#/credits')
  expect(result.current.page).toBe('credits')
})

test('外部 hashchange（回退按钮）同步状态', () => {
  const { result } = renderHook(() => useHashRoute())
  act(() => {
    window.location.hash = '#/models'
    window.dispatchEvent(new HashChangeEvent('hashchange'))
  })
  expect(result.current.page).toBe('models')
})

test('未知 hash 回退到 overview', () => {
  window.location.hash = '#/nonexistent'
  const { result } = renderHook(() => useHashRoute())
  expect(result.current.page).toBe('overview')
})
