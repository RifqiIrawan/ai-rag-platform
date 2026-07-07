import { apiGet, apiUpload } from './client'

export type DocumentStatus = 'uploaded' | 'processing' | 'indexed' | 'failed'

export interface DocumentSummary {
  id: string
  filename: string
  content_type: string
  size_bytes: number
  status: DocumentStatus
  created_at: string
}

export interface UploadResponse {
  id: string
  filename: string
  object_key: string
}

export function listDocuments(): Promise<{ documents: DocumentSummary[] }> {
  return apiGet('/api/v1/documents')
}

export function uploadDocument(file: File): Promise<UploadResponse> {
  return apiUpload('/api/v1/documents', file)
}
