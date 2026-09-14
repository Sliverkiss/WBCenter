// 登录页：面板本地口令校验。
import { FormEvent, useState } from 'react'
import { login } from '../services/api'

export function Login({ onDone }: { onDone: () => Promise<void> }) {
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const submit = async (e: FormEvent) => {
    e.preventDefault()
    try {
      await login({ username, password })
      await onDone()
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
    }
  }
  return (
    <div className="login">
      <div className="login-card">
        <div className="eyebrow">LOCAL CONTROL ROOM</div>
        <h1>账号工作台</h1>
        <p>只读取账号凭据并直接连接 WorkBuddy 上游。</p>
        <form onSubmit={submit}>
          <label>用户名<input value={username} onChange={e => setUsername(e.target.value)} autoComplete="username" /></label>
          <label>密码<input value={password} onChange={e => setPassword(e.target.value)} type="password" autoFocus autoComplete="current-password" /></label>
          {error && <div className="error">{error}</div>}
          <button className="black" type="submit">进入控制台</button>
        </form>
      </div>
    </div>
  )
}
