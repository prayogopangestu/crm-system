"use client"

import React from "react"
import { usePathname } from "next/navigation"
import { Sidebar } from "@/components/common/Sidebar"
import { Header } from "@/components/common/Header"
import { SidebarProvider, useSidebar } from "@/context/SidebarContext"
import { Toaster } from "@/components/ui/sonner"

function LayoutContent({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isAuthPage = pathname === "/login" || pathname === "/daftar"
  const { isCollapsed } = useSidebar()

  if (isAuthPage) {
    return <div className="min-h-screen w-full bg-background flex items-center justify-center p-4">{children}</div>
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
