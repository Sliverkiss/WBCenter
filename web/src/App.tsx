// 应用壳：会话检查 + hash 路由装配 + 侧栏布局。页面实现见 pages/。
import { useCallback, useEffect, useState } from 'react'
import type { Page } from './types'
import { getSession, logout, setUnauthorizedHandler } from './services/api'
import { useHashRoute } from './hooks/useHashRoute'
import { useTheme } from './hooks/useTheme'
import { useNotice } from './hooks/useNotice'
import { Notice } from './components/Notice'
import { Login } from './pages/Login'
import { Overview } from './pages/Overview'
import { Accounts } from './pages/Accounts'
import { Credits } from './pages/Credits'
import { Models } from './pages/Models'
import { Automation } from './pages/Automation'
import { Scheduler } from './pages/Scheduler'
import { Activities } from './pages/Activities'
import { OAuth } from './pages/OAuth'
import { Stats } from './pages/Stats'

const NAV: [Page, string, string][] = [
  ['overview', '概览', '⌁'], ['accounts', '账号管理', '◉'], ['credits', '积分管理', '＋'],
  ['models', '模型管理', '◇'], ['stats', '统计分析', '∿'], ['automation', '自动化管理', '◷'],
  ['scheduler', '账号云端定时任务', '☁'],
  ['activities', '活动管理', '✦'], ['oauth', '添加账号', '→'],
]

export default function App() {
  const [authed, setAuthed] = useState<boolean | null>(null)
  const { page, navigate } = useHashRoute()
  const { theme, toggle } = useTheme()
  const { notice, setNotice, clear } = useNotice()

  const check = useCallback(async () => {
    try { const s = await getSession(); setAuthed(s.authenticated) } catch { setAuthed(false) }
  }, [])
  useEffect(() => { void check() }, [check])
  // 任意 API 返回 401 时全局登出，回到登录页。
  useEffect(() => {
    setUnauthorizedHandler(() => setAuthed(false))
    return () => setUnauthorizedHandler(null)
  }, [])
  if (authed === null) return <div className="center">正在连接本地控制台…</div>
  if (!authed) return <Login onDone={check} />

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="wordmark"><b>WORK<br />BUDDY<span>›</span></b><small>LOCAL CONTROL ROOM</small></div>
        <nav>{NAV.map(([id, label, icon]) => <button key={id} className={page === id ? 'nav active' : 'nav'} onClick={() => navigate(id)}><i>{icon}</i>{label}</button>)}</nav>
        <div className="side-bottom">
          <button className="link" onClick={toggle}>{theme === 'dark' ? '切换为浅色主题' : '切换为暗色主题'}</button>
          <span>第三方本地控制台</span>
          <button className="link" onClick={() => void logout().then(() => setAuthed(false))}>退出登录</button>
        </div>
      </aside>
      <main className="main">
        {notice && <Notice text={notice} onClose={clear} />}
        {page === 'overview' && <Overview onPage={navigate} />}
        {page === 'accounts' && <Accounts onNotice={setNotice} />}
        {page === 'credits' && <Credits />}
        {page === 'models' && <Models />}
        {page === 'stats' && <Stats />}
        {page === 'automation' && <Automation onNotice={setNotice} />}
        {page === 'scheduler' && <Scheduler />}
        {page === 'activities' && <Activities />}
        {page === 'oauth' && <OAuth onNotice={setNotice} onDone={() => navigate('accounts')} />}
      </main>
    </div>
  )
}
