import { useRef, useState, type FormEvent } from 'react'
import { queryRag, type RagSource } from '../api/rag'

interface ChatMessage {
  id: number
  role: 'user' | 'assistant'
  text: string
  sources?: RagSource[]
}

export function ChatPanel() {
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const idRef = useRef(0)

  const nextId = () => {
    idRef.current += 1
    return idRef.current
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    const query = input.trim()
    if (!query || loading) return

    setMessages((prev) => [...prev, { id: nextId(), role: 'user', text: query }])
    setInput('')
    setLoading(true)
    setError(null)

    try {
      const res = await queryRag(query)
      setMessages((prev) => [...prev, { id: nextId(), role: 'assistant', text: res.answer, sources: res.sources }])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Query failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="flex h-full flex-col rounded-xl border border-slate-200 bg-white dark:border-slate-700 dark:bg-slate-800/50">
      <div className="border-b border-slate-200 px-4 py-3 dark:border-slate-700">
        <h2 className="text-sm font-semibold text-slate-700 dark:text-slate-200">Ask your documents</h2>
      </div>
      <div className="flex-1 space-y-3 overflow-y-auto p-4">
        {messages.length === 0 && (
          <p className="text-sm text-slate-400">Upload a document, then ask a question about it.</p>
        )}
        {messages.map((m) => (
          <div key={m.id} className={`flex ${m.role === 'user' ? 'justify-end' : 'justify-start'}`}>
            <div
              className={`max-w-[80%] rounded-2xl px-3 py-2 text-sm ${
                m.role === 'user'
                  ? 'bg-indigo-600 text-white'
                  : 'bg-slate-100 text-slate-700 dark:bg-slate-700 dark:text-slate-100'
              }`}
            >
              <p className="whitespace-pre-wrap">{m.text}</p>
              {m.sources && m.sources.length > 0 && (
                <details className="mt-2 text-xs opacity-80">
                  <summary className="cursor-pointer">{m.sources.length} source(s)</summary>
                  <ul className="mt-1 space-y-1">
                    {m.sources.map((s, i) => (
                      <li key={i} className="rounded bg-black/5 p-1.5 dark:bg-white/10">
                        <span className="font-mono">{s.score.toFixed(2)}</span> — {s.text ?? '(no text)'}
                      </li>
                    ))}
                  </ul>
                </details>
              )}
            </div>
          </div>
        ))}
        {loading && <p className="text-sm text-slate-400">Thinking…</p>}
      </div>
      {error && <p className="px-4 pb-2 text-xs text-red-600 dark:text-red-400">{error}</p>}
      <form onSubmit={handleSubmit} className="flex gap-2 border-t border-slate-200 p-3 dark:border-slate-700">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Ask a question…"
          className="flex-1 rounded-lg border border-slate-300 bg-transparent px-3 py-2 text-sm text-slate-800 outline-none focus:border-indigo-500 dark:border-slate-600 dark:text-slate-100"
        />
        <button
          type="submit"
          disabled={loading || !input.trim()}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-500 disabled:opacity-50"
        >
          Send
        </button>
      </form>
    </section>
  )
}
