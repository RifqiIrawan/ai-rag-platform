import { apiPost } from './client'

export interface LoginResponse {
  token: string
}

export interface RegisterResponse {
  id: string
  email: string
}

export function register(email: string, password: string): Promise<RegisterResponse> {
  return apiPost<RegisterResponse>('/api/v1/auth/register', { email, password })
}

export function login(email: string, password: string): Promise<LoginResponse> {
  return apiPost<LoginResponse>('/api/v1/auth/login', { email, password })
}
