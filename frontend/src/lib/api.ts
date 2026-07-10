const configuredApiBaseUrl = process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "")

export const API_BASE_URL =
  configuredApiBaseUrl || (process.env.NODE_ENV === "development" ? "http://localhost:8080" : "")

export const TOKEN_KEY = "crm_token"
export const USER_KEY = "crm_user"

interface ApiRequestOptions {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE"
  body?: unknown
  auth?: boolean
  headers?: HeadersInit
}

interface ErrorBody {
  error?: {
    code?: string
    message?: string
  }
  message?: string
}

export interface MutationResponse<T> {
  success: boolean
  message?: string
  data: T
  inviteUrl?: string
}

export class ApiError extends Error {
  status: number
  code?: string

  constructor(message: string, status: number, code?: string) {
    super(message)
    this.name = "ApiError"
    this.status = status
    this.code = code
  }
}

export function getStoredToken() {
  if (typeof window === "undefined") return null
  return localStorage.getItem(TOKEN_KEY)
}

export function setStoredAuth(token: string, user: unknown) {
  if (typeof window === "undefined") return
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export function clearStoredAuth() {
  if (typeof window === "undefined") return
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  localStorage.removeItem("crm_logged_in")
}

function buildUrl(path: string) {
  if (path.startsWith("http")) return path
  if (!API_BASE_URL) {
    throw new ApiError(
      "NEXT_PUBLIC_API_URL belum dikonfigurasi. Isi dengan URL backend Railway.",
      0,
      "MISSING_API_URL"
    )
  }
  return `${API_BASE_URL}${path.startsWith("/") ? path : `/${path}`}`
}

async function parseError(response: Response) {
  let body: ErrorBody | null = null
  try {
    body = await response.json()
  } catch {
    body = null
  }
  return new ApiError(
    body?.error?.message || body?.message || "Request backend gagal",
    response.status,
    body?.error?.code
  )
}

export async function apiRequest<T>(path: string, options: ApiRequestOptions = {}): Promise<T> {
  const headers = new Headers(options.headers)
  const token = getStoredToken()

  if (options.body !== undefined) {
    headers.set("Content-Type", "application/json")
  }
  if (options.auth !== false && token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const response = await fetch(buildUrl(path), {
    method: options.method || "GET",
    headers,
    body: options.body === undefined ? undefined : JSON.stringify(options.body),
    cache: "no-store",
  })

  if (!response.ok) {
    if (response.status === 401) {
      clearStoredAuth()
    }
    throw await parseError(response)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return response.json() as Promise<T>
}

export async function apiDownload(path: string, filename: string) {
  const token = getStoredToken()
  const headers = new Headers()
  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const response = await fetch(buildUrl(path), { headers })
  if (!response.ok) {
    throw await parseError(response)
  }

  const blob = await response.blob()
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

export function queryString(params: Record<string, string | number | undefined | null>) {
  const search = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== "") {
      search.set(key, String(value))
    }
  }
  const value = search.toString()
  return value ? `?${value}` : ""
}
