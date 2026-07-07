import { useCallback, useEffect, useRef, useState, type ChangeEvent } from 'react'
import { listDocuments, uploadDocument, type DocumentSummary, type DocumentStatus } from '../api/documents'

const STATUS_STYLES: Record<DocumentStatus, string> = {
  uploaded: 'bg-slate-100 text-slate-600 dark:bg-slate-700 dark:text-slate-300',
  processing: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300',
  indexed: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300',
  failed: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300',
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

export function DocumentsPanel() {
  const [documents, setDocuments] = useState<DocumentSummary[]>([])
  const [loading, setLoading] = useState(true)
  const [uploading, setUploading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const refresh = useCallback(async () => {
    try {
      const res = await listDocuments()
      setDocuments(res.documents)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load documents')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    refresh()
    // Poll so processing/uploaded documents flip to indexed/failed without a manual refresh.
    const interval = setInterval(refresh, 3000)
    return () => clearInterval(interval)
  }, [refresh])

  const handleFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    setError(null)
    try {
      await uploadDocument(file)
      await refresh()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setUploading(false)
      if (fileInputRef.current) fileInputRef.current.value = ''
    }
  }

  return (
    <section className="flex h-full flex-col rounded-xl border border-slate-200 bg-white dark:border-slate-700 dark:bg-slate-800/50">
      <div className="flex items-center justify-between border-b border-slate-200 px-4 py-3 dark:border-slate-700">
        <h2 className="text-sm font-semibold text-slate-700 dark:text-slate-200">Documents</h2>
        <label className="cursor-pointer rounded-md bg-indigo-600 px-3 py-1.5 text-xs font-medium text-white hover:bg-indigo-500 has-disabled:opacity-50">
          {uploading ? 'Uploading…' : 'Upload'}
          <input ref={fileInputRef} type="file" className="hidden" onChange={handleFileChange} disabled={uploading} />
        </label>
      </div>
      {error && <p className="px-4 py-2 text-xs text-red-600 dark:text-red-400">{error}</p>}
      <div className="flex-1 overflow-y-auto p-2">
        {loading ? (
          <p className="px-2 py-4 text-sm text-slate-400">Loading…</p>
        ) : documents.length === 0 ? (
          <p className="px-2 py-4 text-sm text-slate-400">No documents uploaded yet.</p>
        ) : (
          <ul className="space-y-1">
            {documents.map((doc) => (
              <li
                key={doc.id}
                className="flex items-center justify-between rounded-lg px-2 py-2 hover:bg-slate-50 dark:hover:bg-slate-700/40"
              >
                <div className="min-w-0">
                  <p className="truncate text-sm text-slate-700 dark:text-slate-200">{doc.filename}</p>
                  <p className="text-xs text-slate-400">{formatBytes(doc.size_bytes)}</p>
                </div>
                <span
                  className={`ml-2 shrink-0 rounded-full px-2 py-0.5 text-[11px] font-medium ${STATUS_STYLES[doc.status] ?? STATUS_STYLES.uploaded}`}
                >
                  {doc.status}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  )
}
