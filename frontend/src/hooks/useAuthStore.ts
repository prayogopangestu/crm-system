import { create } from "zustand"
import { apiRequest, clearStoredAuth, setStoredAuth, TOKEN_KEY, USER_KEY } from "@/lib/api"
import { UserProfile } from "@/types/crm"

interface AuthState {
  isLoading: boolean
  error: string | null
  token: string | null
  user: Pick<UserProfile, "id" | "name" | "role"> | null
  setError: (error: string | null) => void
  setIsLoading: (loading: boolean) => void
  hydrate: () => void
  login: (input: { email: string; password: string }, onSuccess: () => void) => Promise<void>
  register: (
    input: { fullName: string; companyName: string; email: string; password: string },
    onSuccess: () => void
  ) => Promise<void>
  logout: (onSuccess: () => void) => void
}

export const useAuthStore = create<AuthState>((set) => ({
  isLoading: false,
  error: null,
  token: null,
  user: null,
  setError: (error) => set({ error }),
  setIsLoading: (isLoading) => set({ isLoading }),
  hydrate: () => {
    if (typeof window === "undefined") return
    const token = localStorage.getItem(TOKEN_KEY)
    const rawUser = localStorage.getItem(USER_KEY)
    let user: AuthState["user"] = null
    if (rawUser) {
      try {
        user = JSON.parse(rawUser)
      } catch {
        user = null
      }
    }
    set({ token, user })
  },
  login: async (input, onSuccess) => {
    set({ isLoading: true, error: null })
    try {
      const result = await apiRequest<{ token: string; user: AuthState["user"] }>("/api/auth/login", {
        method: "POST",
        body: input,
        auth: false,
      })
      setStoredAuth(result.token, result.user)
      set({ token: result.token, user: result.user, isLoading: false })
      onSuccess()
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Login gagal",
      })
    }
  },
  register: async (input, onSuccess) => {
    set({ isLoading: true, error: null })
    try {
      await apiRequest<{ success: boolean; message: string }>("/api/auth/register", {
        method: "POST",
        body: input,
        auth: false,
      })
      const result = await apiRequest<{ token: string; user: AuthState["user"] }>("/api/auth/login", {
        method: "POST",
        body: { email: input.email, password: input.password },
        auth: false,
      })
      setStoredAuth(result.token, result.user)
      set({ token: result.token, user: result.user, isLoading: false })
      onSuccess()
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Registrasi gagal",
      })
    }
  },
  logout: (onSuccess) => {
    clearStoredAuth()
    set({ token: null, user: null })
    onSuccess()
  }
}))
