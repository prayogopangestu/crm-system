"use client"

import React, { useEffect, useState } from "react"
import { Button } from "@/components/ui/Button"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/Avatar"
import { useSidebar } from "@/context/SidebarContext"
import { useTheme } from "next-themes"
import { Sun, Moon } from "lucide-react"
import { toast } from "sonner"

export function Header() {
  const { toggleCollapse, toggleMobile, isCollapsed } = useSidebar()
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  const handleToggle = () => {
    if (window.innerWidth < 768) {
      toggleMobile()
    } else {
      toggleCollapse()
    }
  }

  const handleQuickCreate = () => {
    toast.info("Fitur Quick Create akan membuka form pembuatan kontak/peluang baru secara instan!")
  }

  return (
    <header className="flex items-center justify-between px-4 md:px-6 w-full h-header-height bg-surface-container-lowest border-b border-outline-variant z-20 sticky top-0">
      {/* Toggle Sidebar Button */}
      <button
        onClick={handleToggle}
        className="p-2 -ml-2 mr-3 rounded-lg text-on-surface-variant hover:bg-surface-container transition-colors cursor-pointer flex items-center justify-center shrink-0"
        title={isCollapsed ? "Buka Sidebar" : "Tutup Sidebar"}
      >
        <span className="material-symbols-outlined">
          {isCollapsed ? "menu" : "menu_open"}
        </span>
      </button>

      {/* Search Input Bar */}
      <div className="flex items-center flex-1 max-w-md">
        <div className="relative w-full">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px] pointer-events-none">
            search
          </span>
          <input
            className="w-full pl-10 pr-4 py-2 bg-surface-container-low border border-outline-variant rounded-lg focus:ring-2 focus:ring-primary/20 focus:border-primary font-body-md text-sm text-on-surface placeholder:text-outline transition-all outline-none"
            placeholder="Cari kontak, tugas, atau peluang..."
            type="text"
          />
        </div>
      </div>

      {/* Action Buttons & Profile */}
      <div className="flex items-center gap-stack-md ml-auto">
        {/* Quick Create CTA */}
        <Button
          onClick={handleQuickCreate}
          variant="default"
          size="sm"
          className="hidden sm:flex items-center gap-2 rounded-full whitespace-nowrap shadow-sm text-xs font-semibold"
        >
          <span className="material-symbols-outlined text-[18px]">add</span>
          Quick Create
        </Button>

        {/* Notifications Icon Button */}
        <button className="p-2 rounded-full text-on-surface-variant hover:bg-surface-container transition-colors relative cursor-pointer">
          <span className="material-symbols-outlined">notifications</span>
          <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-error rounded-full border border-surface-container-lowest"></span>
        </button>

        {/* Theme Toggle Button */}
        <button
          onClick={() => {
            document.documentElement.classList.add("theme-transition")
            setTheme(theme === "dark" ? "light" : "dark")
            setTimeout(() => {
              document.documentElement.classList.remove("theme-transition")
            }, 300)
          }}
          className="p-2 rounded-full text-on-surface-variant hover:bg-surface-container transition-colors relative cursor-pointer flex items-center justify-center"
          title={mounted && theme === "dark" ? "Ubah ke mode terang" : "Ubah ke mode gelap"}
        >
          {mounted && theme === "dark" ? (
            <Sun className="h-5 w-5 text-amber-500 animate-in duration-300" />
          ) : (
            <Moon className="h-5 w-5 text-slate-700 dark:text-slate-350 animate-in duration-300" />
          )}
        </button>

        {/* Vertical Divider */}
        <div className="h-8 w-[1px] bg-outline-variant mx-1 hidden sm:block"></div>

        {/* User Card */}
        <div className="flex items-center space-x-3 border-l border-outline-variant pl-4 cursor-pointer group select-none">
          <Avatar className="h-8 w-8 ring-2 ring-transparent group-hover:ring-primary transition-all">
            <AvatarImage
              src="https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&q=80&w=100"
              alt="Sarah Jenkins"
            />
            <AvatarFallback>SJ</AvatarFallback>
          </Avatar>
          <div className="hidden lg:block text-left">
            <p className="font-label-md text-xs font-semibold text-on-surface group-hover:text-primary transition-colors leading-none">
              Sarah Jenkins
            </p>
            <p className="font-label-sm text-[10px] text-on-surface-variant leading-none mt-1">
              Admin
            </p>
          </div>
        </div>
      </div>
    </header>
  )
}
