"use client"

import React, { useState } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { Card } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"

export default function DaftarPage() {
  const router = useRouter()
  const [fullName, setFullName] = useState("")
  const [companyName, setCompanyName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState("")

  const handleRegister = (e: React.FormEvent) => {
    e.preventDefault()
    setError("")

    if (!fullName || !companyName || !email || !password || !confirmPassword) {
      setError("Semua kolom harus diisi.")
      return
    }

    if (password !== confirmPassword) {
      setError("Konfirmasi kata sandi tidak cocok.")
      return
    }

    if (password.length < 6) {
      setError("Kata sandi minimal harus 6 karakter.")
      return
    }

    setIsLoading(true)
    // Simulate API loading for registration
    setTimeout(() => {
      setIsLoading(false)
      localStorage.setItem("crm_logged_in", "true")
      router.push("/")
    }, 1500)
  }

  const handleSocialRegister = (provider: string) => {
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
      <div className="flex flex-col items-center mb-6">
        <div className="w-12 h-12 rounded-xl bg-primary flex items-center justify-center text-on-primary font-bold shadow-md mb-3">
          <span className="material-symbols-outlined text-[28px]">domain</span>
        </div>
        <h2 className="text-2xl font-bold text-on-surface tracking-tight">Daftar Akun Baru</h2>
        <p className="text-xs text-on-surface-variant mt-1">Mulai kelola sales &amp; hubungan klien Anda sekarang</p>
      </div>

      {error && (
        <div className="mb-4 p-3 rounded-lg bg-error-container text-on-error-container text-xs font-semibold flex items-center gap-2 animate-shake">
          <span className="material-symbols-outlined text-[18px]">error</span>
          <span>{error}</span>
        </div>
      )}

      {/* Register Form */}
      <form onSubmit={handleRegister} className="space-y-4">
        <div>
          <label className="block text-xs font-semibold text-on-surface-variant mb-1">
            Nama Lengkap
          </label>
          <Input
            required
            type="text"
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            placeholder="Sarah Jenkins"
          />
        </div>

        <div>
          <label className="block text-xs font-semibold text-on-surface-variant mb-1">
            Nama Perusahaan
          </label>
          <Input
            required
            type="text"
            value={companyName}
            onChange={(e) => setCompanyName(e.target.value)}
            placeholder="Acme Corporation"
          />
        </div>

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
          <label className="block text-xs font-semibold text-on-surface-variant mb-1">
            Kata Sandi
          </label>
          <Input
            required
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
          />
        </div>

        <div>
          <label className="block text-xs font-semibold text-on-surface-variant mb-1">
            Konfirmasi Kata Sandi
          </label>
          <Input
            required
            type="password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            placeholder="••••••••"
          />
        </div>

        <Button
          type="submit"
          variant="default"
          className="w-full font-semibold py-2.5 mt-2"
          disabled={isLoading}
        >
          {isLoading ? "Membuat Akun..." : "Daftar Akun"}
        </Button>
      </form>

      {/* Social Divider */}
      <div className="relative my-5">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-outline-variant"></div>
        </div>
        <div className="relative flex justify-center text-xs uppercase">
          <span className="bg-surface-container-lowest px-2 text-on-surface-variant font-medium text-[10px]">
            Atau daftar dengan
          </span>
        </div>
      </div>

      {/* Social Providers */}
      <div className="grid grid-cols-2 gap-3">
        <Button
          type="button"
          variant="outline"
          onClick={() => handleSocialRegister("Google")}
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
          onClick={() => handleSocialRegister("Microsoft")}
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

      {/* Redirect to Login */}
      <div className="mt-6 text-center text-xs text-on-surface-variant font-medium">
        Sudah memiliki akun?{" "}
        <Link href="/login" className="text-primary font-bold hover:underline">
          Masuk di sini
        </Link>
      </div>
    </Card>
  )
}
