import { useEffect, useRef, useState } from 'react'
import { notificationsWsUrl } from '../api/notifications'
import { useAuth } from '../auth/AuthContext'

interface NotificationItem {
  id: number
  message: string
  channel: string
  receivedAt: number
}

export function NotificationBell() {
  const { token } = useAuth()
  const [items, setItems] = useState<NotificationItem[]>([])
  const [open, setOpen] = useState(false)
  const [connected, setConnected] = useState(false)
  const idRef = useRef(0)

  useEffect(() => {
    if (!token) return

    const ws = new WebSocket(notificationsWsUrl(token))
    ws.onopen = () => setConnected(true)
    ws.onclose = () => setConnected(false)
    ws.onerror = () => setConnected(false)
    ws.onmessage = (event: MessageEvent<string>) => {
      try {
        const data = JSON.parse(event.data) as { channel: string; message: string }
        idRef.current += 1
        setItems((prev) =>
          [{ id: idRef.current, message: data.message, channel: data.channel, receivedAt: Date.now() }, ...prev].slice(0, 20),
        )
      } catch {
        // ignore malformed payloads
      }
    }

    return () => ws.close()
  }, [token])

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((o) => !o)}
        className="relative rounded-full p-2 text-slate-600 hover:bg-slate-100 dark:text-slate-300 dark:hover:bg-slate-800"
        aria-label="Notifications"
      >
        <BellIcon />
        <span
          className={`absolute right-1.5 top-1.5 h-2 w-2 rounded-full ${connected ? 'bg-emerald-500' : 'bg-slate-400'}`}
          title={connected ? 'Connected' : 'Disconnected'}
        />
        {items.length > 0 && (
          <span className="absolute -right-1 -top-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-semibold text-white">
            {items.length}
          </span>
        )}
      </button>
      {open && (
        <div className="absolute right-0 z-10 mt-2 w-80 rounded-lg border border-slate-200 bg-white shadow-lg dark:border-slate-700 dark:bg-slate-800">
          <div className="flex items-center justify-between border-b border-slate-200 px-3 py-2 dark:border-slate-700">
            <span className="text-sm font-medium text-slate-700 dark:text-slate-200">Notifications</span>
            {items.length > 0 && (
              <button
                onClick={() => setItems([])}
                className="text-xs text-slate-500 hover:underline dark:text-slate-400"
              >
                Clear
              </button>
            )}
          </div>
          <div className="max-h-80 overflow-y-auto">
            {items.length === 0 ? (
              <p className="px-3 py-6 text-center text-sm text-slate-400">No notifications yet</p>
            ) : (
              items.map((item) => (
                <div key={item.id} className="border-b border-slate-100 px-3 py-2 text-sm last:border-0 dark:border-slate-700/60">
                  <p className="text-slate-700 dark:text-slate-200">{item.message}</p>
                  <p className="mt-0.5 text-xs text-slate-400">{new Date(item.receivedAt).toLocaleTimeString()}</p>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}

function BellIcon() {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="h-5 w-5">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M14.857 17.082a23.85 23.85 0 0 0 5.454-1.31A8.97 8.97 0 0 1 18 9.75V9A6 6 0 0 0 6 9v.75a8.97 8.97 0 0 1-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.26 24.26 0 0 1-5.714 0m5.714 0a3 3 0 1 1-5.714 0"
      />
    </svg>
  )
}
