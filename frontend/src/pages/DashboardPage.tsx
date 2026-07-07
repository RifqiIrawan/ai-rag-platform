import { ChatPanel } from '../components/ChatPanel'
import { DocumentsPanel } from '../components/DocumentsPanel'
import { NotificationBell } from '../components/NotificationBell'
import { useAuth } from '../auth/AuthContext'

export function DashboardPage() {
  const { logout, userId } = useAuth()

  return (
    <div className="flex h-full flex-col bg-slate-50 dark:bg-slate-900">
      <header className="flex items-center justify-between border-b border-slate-200 bg-white px-6 py-3 dark:border-slate-700 dark:bg-slate-800">
        <h1 className="text-base font-semibold text-slate-800 dark:text-slate-100">ai-rag-platform</h1>
        <div className="flex items-center gap-3">
          <span className="hidden text-xs text-slate-400 sm:inline">{userId}</span>
          <NotificationBell />
          <button
            onClick={logout}
            className="rounded-md border border-slate-300 px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50 dark:border-slate-600 dark:text-slate-300 dark:hover:bg-slate-700"
          >
            Log out
          </button>
        </div>
      </header>
      <main className="grid flex-1 gap-4 overflow-hidden p-4 md:grid-cols-[320px_1fr]">
        <div className="min-h-0">
          <DocumentsPanel />
        </div>
        <div className="min-h-0">
          <ChatPanel />
        </div>
      </main>
    </div>
  )
}
