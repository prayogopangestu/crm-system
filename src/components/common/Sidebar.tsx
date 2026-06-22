"use client"

import React from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"
import { useSidebar } from "@/context/SidebarContext"

interface NavItem {
  name: string
  href: string
  icon: string
}

const navItems: NavItem[] = [
  { name: "Dashboard", href: "/", icon: "dashboard" },
  { name: "Kontak & Perusahaan", href: "/kontak", icon: "business_center" },
  { name: "Pipeline Penjualan", href: "/pipeline", icon: "view_kanban" },
  { name: "Tugas & Aktivitas", href: "/tugas", icon: "assignment" },
  { name: "Laporan", href: "/laporan", icon: "analytics" },
  { name: "Pengaturan", href: "/pengaturan", icon: "settings" }
]

export function Sidebar() {
  const pathname = usePathname()
  const { isCollapsed, isMobileOpen, closeMobile } = useSidebar()

  return (
    <>
      {/* Mobile Sidebar Overlay Backdrop */}
      {isMobileOpen && (
        <div 
          className="fixed inset-0 bg-black/60 z-40 md:hidden animate-in fade-in duration-200" 
          onClick={closeMobile}
        />
      )}

      {/* Sidebar Container */}
      <nav 
        className={`flex flex-col h-full fixed left-0 top-0 bg-[#0f172a] p-4 z-50 md:z-30 border-r border-slate-800 transition-all duration-300 ease-in-out ${
          isCollapsed ? "w-[70px] px-2.5" : "w-[260px]"
        } ${
          isMobileOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"
        }`}
      >
        {/* Brand Header */}
        <div className={`flex items-center gap-3 mb-6 px-1.5 ${isCollapsed ? "justify-center" : ""}`}>
          <div className="w-8 h-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white font-bold shrink-0 shadow-sm shadow-indigo-500/30">
            <span className="material-symbols-outlined text-[20px]">domain</span>
          </div>
          {!isCollapsed && (
            <div className="animate-in fade-in duration-300">
              <h1 className="font-headline-sm text-[16px] font-bold text-white tracking-tight leading-none">
                CRM Enterprise
              </h1>
              <p className="font-label-sm text-[10px] text-slate-400 mt-1 leading-none">
                Sistem Manajemen Sales
              </p>
            </div>
          )}
        </div>

        {/* Navigation List */}
        <div className="flex-1 overflow-y-auto space-y-1">
          {navItems.map((item) => {
            const isActive = pathname === item.href
            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={closeMobile}
                className={`flex items-center rounded-lg transition-all ${
                  isCollapsed ? "justify-center p-2.5" : "gap-3 px-3 py-2.5"
                } ${
                  isActive
                    ? "text-white font-bold bg-indigo-600 shadow-md shadow-indigo-500/20 scale-[0.98]"
                    : "text-slate-400 hover:bg-slate-800/60 hover:text-slate-100 group"
                }`}
                title={isCollapsed ? item.name : undefined}
              >
                <span
                  className={`material-symbols-outlined shrink-0 ${
                    isActive ? "text-white" : "group-hover:text-indigo-400 transition-colors"
                  }`}
                  style={isActive ? { fontVariationSettings: "'FILL' 1" } : {}}
                >
                  {item.icon}
                </span>
                {!isCollapsed && (
                  <span className="font-body-md text-sm font-medium animate-in fade-in duration-300">
                    {item.name}
                  </span>
                )}
              </Link>
            )
          })}
        </div>

        {/* Footer Actions (Help & Logout) */}
        <div className="mt-auto pt-4 border-t border-slate-800/80 space-y-2">
          <button 
            className={`w-full flex items-center justify-center rounded-lg border border-slate-800 hover:bg-slate-800/60 transition-colors text-slate-400 hover:text-slate-100 cursor-pointer text-sm font-medium ${
              isCollapsed ? "p-2" : "space-x-2 py-2 px-4"
            }`}
            title={isCollapsed ? "Bantuan" : undefined}
          >
            <span className="material-symbols-outlined text-[18px] shrink-0">help</span>
            {!isCollapsed && <span className="font-label-md animate-in fade-in duration-300">Bantuan</span>}
          </button>
          
          <Link
            href="/login"
            onClick={closeMobile}
            className={`w-full flex items-center justify-center rounded-lg bg-red-500/10 hover:bg-red-500/20 text-red-400 transition-colors cursor-pointer text-sm font-bold text-center ${
              isCollapsed ? "p-2" : "space-x-2 py-2.5 px-4"
            }`}
            title={isCollapsed ? "Keluar" : undefined}
          >
            <span className="material-symbols-outlined text-[18px] shrink-0">logout</span>
            {!isCollapsed && <span className="font-label-md animate-in fade-in duration-300">Keluar</span>}
          </Link>
        </div>
      </nav>
    </>
  )
}
