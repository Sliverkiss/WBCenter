import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './styles/global.css'

// dev 环境默认启用 MSW（页面全部走 mock）；VITE_USE_MSW=0 时回退到真实后端代理。
async function enableMocking() {
  if (!import.meta.env.DEV || import.meta.env.VITE_USE_MSW === '0') return
  const { worker } = await import('./mocks/browser')
  await worker.start({ onUnhandledRequest: 'bypass' })
}

void enableMocking().then(() => {
  ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>,
  )
})
