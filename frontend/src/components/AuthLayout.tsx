import type { ReactNode } from 'react'

export function AuthLayout({ title, children, footer }: { title: string; children: ReactNode; footer?: ReactNode }) {
  return (
    <div className="flex h-full min-h-full items-center justify-center bg-slate-50 px-4 dark:bg-slate-900">
      <div className="w-full max-w-sm rounded-2xl border border-slate-200 bg-white p-6 shadow-sm dark:border-slate-700 dark:bg-slate-800">
        <h1 className="mb-6 text-xl font-semibold text-slate-800 dark:text-slate-100">{title}</h1>
        {children}
        {footer && <p className="mt-4 text-center text-sm text-slate-500 dark:text-slate-400">{footer}</p>}
      </div>
    </div>
  )
}
