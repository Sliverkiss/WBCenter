// 轻量 hash 路由：#/accounts 形式，读 location.hash，监听 hashchange。
// 30 行内手写实现，不引 react-router。
import { useCallback, useEffect, useState } from 'react'
import { PAGES, type Page } from '../types'

function readHash(): Page {
  const raw = window.location.hash.replace(/^#\/?/, '')
  return (PAGES as string[]).includes(raw) ? (raw as Page) : 'overview'
}

export function useHashRoute(): { page: Page; navigate: (p: Page) => void } {
  const [page, setPage] = useState<Page>(readHash)

  useEffect(() => {
    const onChange = () => setPage(readHash())
    window.addEventListener('hashchange', onChange)
    return () => window.removeEventListener('hashchange', onChange)
  }, [])

  const navigate = useCallback((p: Page) => {
    window.location.hash = `#/${p}`
    // hash 不变时（重复点击当前页）不依赖事件，直接同步一次
    setPage(p)
  }, [])

  return { page, navigate }
}
