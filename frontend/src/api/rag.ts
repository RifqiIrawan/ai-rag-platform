import { apiPost } from './client'

export interface RagSource {
  score: number
  text?: string
}

export interface RagQueryResponse {
  query: string
  answer: string
  sources: RagSource[]
}

export function queryRag(query: string): Promise<RagQueryResponse> {
  return apiPost<RagQueryResponse>('/api/v1/rag/query', { query })
}
