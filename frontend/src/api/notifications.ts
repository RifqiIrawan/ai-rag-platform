import { apiBaseUrl, apiPost } from './client'

export function publishNotification(message: string, targetUserId?: string): Promise<{ status: string; channel: string }> {
  return apiPost('/api/v1/notifications', { message, target_user_id: targetUserId })
}

// Browsers can't set custom headers during a WebSocket handshake, so the
// gateway's auth middleware also accepts the JWT as a ?token= query param.
export function notificationsWsUrl(token: string): string {
  const wsBase = apiBaseUrl().replace(/^http/, 'ws')
  return `${wsBase}/api/v1/notifications/ws?token=${encodeURIComponent(token)}`
}
