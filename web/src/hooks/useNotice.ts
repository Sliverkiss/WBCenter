// 顶部通知条状态 hook。
import { useCallback, useState } from 'react'

export function useNotice(): { notice: string; setNotice: (s: string) => void; clear: () => void } {
  const [notice, setNotice] = useState('')
  const clear = useCallback(() => setNotice(''), [])
  return { notice, setNotice, clear }
}
