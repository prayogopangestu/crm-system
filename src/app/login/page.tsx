"use client"

import React, { useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Card } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"

export default function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [isLoading, setIsLoading] = useState(false)

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault()
    if (!email || !password) return

    setIsLoading(true)
    // Simulate API loading
    setTimeout(() => {
      setIsLoading(false)
      localStorage.setItem("crm_logged_in", "true")
      router.push("/")
    }, 1200)
  }

  const handleSocialLogin = (provider: string) => {
    setIsLoading(true)
    setTimeout(() => {
      setIsLoading(false)
      localStorage.setItem("crm_logged_in", "true")
      router.push("/")
    }, 1000)
  }

  return (
    <Card className="w-full max-w-md p-8 shadow-xl bg-surface-container-lowest animate-in fade-in zoom-in-95 duration-200">
      {/* Brand Logo & Heading */}
      <div className="flex flex-col items-center mb-8">
        <div className="w-12 h-12 rounded-xl bg-primary flex items-center justify-center text-on-primary font-bold shadow-md mb-3">
          <span className="material-symbols-outlined text-[28px]">domain</span>
        </div>
        <h2 className="text-2xl font-bold text-on-surface tracking-tight">CRM Enterprise</h2>
        <p className="text-xs text-on-surface-variant mt-1">Masuk untuk mengelola sales &amp; hubungan klien</p>
      </div>

      {/* Login Form */}
      <form onSubmit={handleLogin} className="space-y-4">
        <div>
          <label className="block text-xs font-semibold text-on-surface-variant mb-1">
            Alamat Email
          </label>
          <Input
            required
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="sarah.j@company.com"
          />
        </div>
        <div>
          <div className="flex justify-between items-center mb-1">
            <label className="block text-xs font-semibold text-on-surface-variant">
              Kata Sandi
            </label>
            <a href="#" className="text-[10px] text-primary hover:underline font-semibold">Lupa sandi?</a>
          </div>
          <Input
            required
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
          />
        </div>

        <Button
          type="submit"
          variant="default"
          className="w-full font-semibold py-2.5"
          disabled={isLoading}
        >
          {isLoading ? "Sedang Masuk..." : "Masuk ke Akun"}
        </Button>
      </form>

      {/* Social Divider */}
      <div className="relative my-6">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-outline-variant"></div>
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-surface-container-lowest px-2 text-on-surface-variant font-medium text-[10px]">
            Atau masuk dengan
          </span>
        </div>
      </div>

      {/* Social Providers */}
      <div className="grid grid-cols-2 gap-3">
        <Button
          type="button"
          variant="outline"
          onClick={() => handleSocialLogin("Google")}
          className="flex items-center justify-center gap-2 text-xs py-2 font-semibold cursor-pointer"
        >
          <svg className="w-4 h-4" viewBox="0 0 24 24">
            <path
              fill="#EA4335"
              d="M12 5.04c1.66 0 3.2.57 4.38 1.69l3.27-3.27C17.67 1.54 14.98 1 12 1 7.35 1 3.4 3.65 1.5 7.5l3.85 3C6.3 7.57 8.93 5.04 12 5.04z"
            />
            <path
              fill="#4285F4"
              d="M23.49 12.27c0-.81-.07-1.59-.2-2.36H12v4.51h6.46c-.29 1.48-1.14 2.73-2.42 3.57l3.77 2.92c2.2-2.03 3.48-5.02 3.48-8.64z"
            />
            <path
              fill="#FBBC05"
              d="M5.35 14.5c-.25-.76-.39-1.57-.39-2.4s.14-1.64.39-2.4L1.5 6.7C.54 8.62 0 10.74 0 12s.54 3.38 1.5 5.3l3.85-3z"
            />
            <path
              fill="#34A853"
              d="M12 23c3.24 0 5.97-1.07 7.96-2.91l-3.77-2.92c-1.12.75-2.54 1.21-4.19 1.21-3.07 0-5.7-2.53-6.65-5.46L1.5 16.38C3.4 20.35 7.35 23 12 23z"
            />
          </svg>
          Google
        </Button>
        <Button
          type="button"
          variant="outline"
          onClick={() => handleSocialLogin("Microsoft")}
          className="flex items-center justify-center gap-2 text-xs py-2 font-semibold cursor-pointer"
        >
          <svg className="w-4 h-4" viewBox="0 0 23 23">
            <path fill="#F25022" d="M0 0h11v11H0z" />
            <path fill="#7FBA00" d="M12 0h11v11H12z" />
            <path fill="#00A4EF" d="M0 12h11v11H0z" />
            <path fill="#FFB900" d="M12 12h11v11H12z" />
          </svg>
          Microsoft
        </Button>
      </div>

      {/* Redirect to Register */}
      <div className="mt-8 text-center text-xs text-on-surface-variant font-medium">
        Belum memiliki akun?{" "}
        <Link href="/daftar" className="text-primary font-bold hover:underline">
          Daftar sekarang
        </Link>
      </div>
    </Card>
  )
}
