import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react'
import { login as apiLogin, register as apiRegister } from '../api/auth'
import { getToken, setToken as persistToken } from '../api/client'
import { decodeJwt, isExpired } from './jwt'

interface AuthContextValue {
  token: string | null
  userId: string | null
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

function readValidToken(): string | null {
  const token = getToken()
  if (!token) return null
  const payload = decodeJwt(token)
  if (!payload || isExpired(payload)) {
    persistToken(null)
    return null
  }
  return token
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => readValidToken())

  useEffect(() => {
    persistToken(token)
  }, [token])

  const login = useCallback(async (email: string, password: string) => {
    const res = await apiLogin(email, password)
    setTokenState(res.token)
  }, [])

  const register = useCallback(
    async (email: string, password: string) => {
      await apiRegister(email, password)
      await login(email, password)
    },
    [login],
  )

  const logout = useCallback(() => setTokenState(null), [])

  const userId = useMemo(() => {
    if (!token) return null
    return decodeJwt(token)?.sub ?? null
  }, [token])

  const value = useMemo(
    () => ({ token, userId, login, register, logout }),
    [token, userId, login, register, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
