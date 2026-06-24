import { create } from "zustand"

interface AuthState {
  isLoading: boolean
  error: string | null
  setError: (error: string | null) => void
  setIsLoading: (loading: boolean) => void
  login: (onSuccess: () => void) => Promise<void>
  register: (onSuccess: () => void) => Promise<void>
  logout: (onSuccess: () => void) => void
}

export const useAuthStore = create<AuthState>((set) => ({
  isLoading: false,
  error: null,
  setError: (error) => set({ error }),
  setIsLoading: (isLoading) => set({ isLoading }),
  login: async (onSuccess) => {
    set({ isLoading: true, error: null })
    // Simulate API request
    await new Promise((resolve) => setTimeout(resolve, 1200))
    localStorage.setItem("crm_logged_in", "true")
    set({ isLoading: false })
    onSuccess()
  },
  register: async (onSuccess) => {
    set({ isLoading: true, error: null })
    // Simulate API request
    await new Promise((resolve) => setTimeout(resolve, 1500))
    localStorage.setItem("crm_logged_in", "true")
    set({ isLoading: false })
    onSuccess()
  },
  logout: (onSuccess) => {
    localStorage.removeItem("crm_logged_in")
    onSuccess()
  }
}))
