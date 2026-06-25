"use client"

import React, { useEffect, useState } from "react"
import { usePathname, useRouter } from "next/navigation"
import { Sidebar } from "@/components/common/Sidebar"
import { Header } from "@/components/common/Header"
import { SidebarProvider, useSidebar } from "@/context/SidebarContext"
import { Toaster } from "@/components/ui/sonner"

function LayoutContent({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const router = useRouter()
  const isAuthPage = pathname === "/login" || pathname === "/daftar"
  const { isCollapsed } = useSidebar()
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null)

  useEffect(() => {
    const loggedIn = localStorage.getItem("crm_logged_in") === "true"
    setIsAuthenticated(loggedIn)

    if (!loggedIn && !isAuthPage) {
      router.push("/login")
    }
  }, [pathname, isAuthPage, router])

  // While checking authentication, show a loading spinner for protected pages
  if (isAuthenticated === null && !isAuthPage) {
    return (
      <div className="min-h-screen w-full bg-background flex flex-col items-center justify-center gap-3">
        <div className="w-10 h-10 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
        <span className="text-xs text-on-surface-variant font-medium">Memuat sistem...</span>
      </div>
    )
  }

  if (isAuthPage) {
    return <div className="min-h-screen w-full bg-background flex items-center justify-center p-4">{children}</div>
  }

  // If not logged in and not auth page, we are redirecting, so show loading
  if (!isAuthenticated) {
    return (
      <div className="min-h-screen w-full bg-background flex flex-col items-center justify-center gap-3">
        <div className="w-10 h-10 border-4 border-primary border-t-transparent rounded-full animate-spin"></div>
        <span className="text-xs text-on-surface-variant font-medium">Mengarahkan ke halaman masuk...</span>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex w-full relative">
      {/* Persistent SideNavBar */}
      <Sidebar />

      {/* Main Content Area */}
      <div 
        className={`flex-1 flex flex-col min-h-screen transition-all duration-300 ease-in-out ${
          isCollapsed 
            ? "md:ml-[70px] w-[calc(100%-70px)]" 
            : "md:ml-[260px] w-[calc(100%-260px)]"
        }`}
      >
        <Header />
        <main className="flex-1 overflow-x-hidden">{children}</main>
      </div>
      <Toaster position="top-right" richColors />
    </div>
  )
}

export function LayoutWrapper({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      <LayoutContent>{children}</LayoutContent>
    </SidebarProvider>
  )
}
