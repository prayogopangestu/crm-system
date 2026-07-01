"use client";

import React, { useEffect, useMemo } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useSidebar } from "@/context/SidebarContext";
import {
  Home,
  BarChart3,
  FileText,
  Bell,
  User,
  Settings,
  LogOut,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useAuthStore } from "@/hooks/useAuthStore";
import { useTaskStore } from "@/hooks/useTaskStore";
import { useLanguage } from "@/context/LanguageContext";

interface NavItem {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  badge?: number;
}

export function Sidebar() {
  const pathname = usePathname();
  const { isCollapsed, isMobileOpen, closeMobile } = useSidebar();
  const { logout } = useAuthStore();
  const { tasks, loadTasks } = useTaskStore();
  const { t } = useLanguage();

  useEffect(() => {
    void loadTasks();
  }, [loadTasks]);

  const todayTasksCount = useMemo(
    () => tasks.filter((t) => t.status === "today" && !t.completed).length,
    [tasks],
  );

  const navItems: NavItem[] = [
    { name: t("nav.dashboard"), href: "/", icon: Home },
    { name: t("nav.contacts"), href: "/kontak", icon: User },
    { name: t("nav.pipeline"), href: "/pipeline", icon: FileText },
    {
      name: t("nav.tasks"),
      href: "/tugas",
      icon: Bell,
      badge: todayTasksCount,
    },
    { name: t("nav.reports"), href: "/laporan", icon: BarChart3 },
    { name: t("nav.settings"), href: "/pengaturan", icon: Settings },
  ];

  const handleLogout = () => {
    logout(() => undefined);
    closeMobile();
  };

  return (
    <>
      {/* Mobile Sidebar Overlay Backdrop */}
      {isMobileOpen && (
        <div
          className="fixed inset-0 bg-black/40 z-40 md:hidden animate-in fade-in duration-200"
          onClick={closeMobile}
        />
      )}

      {/* Sidebar Container */}
      <nav
        className={cn(
          "flex flex-col h-full fixed left-0 top-0 bg-white dark:bg-slate-900 border-r border-slate-100 dark:border-slate-800 transition-all duration-300 ease-in-out z-50 md:z-30",
          isCollapsed ? "w-[70px] p-3" : "w-[260px] p-5",
          isMobileOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0",
        )}
      >
        {/* Brand Header */}
        <div
          className={cn(
            "flex items-center mb-6 px-1.5 transition-all duration-300 ease-in-out overflow-hidden",
            isCollapsed ? "justify-center gap-0" : "gap-3",
          )}
        >
          <div className="w-8 h-8 rounded-lg bg-indigo-600 flex items-center justify-center text-white font-bold shrink-0 shadow-sm shadow-indigo-500/30">
            <span className="material-symbols-outlined text-[20px]">
              domain
            </span>
          </div>
          <div
            className={cn(
              "flex flex-col transition-all duration-300 ease-in-out overflow-hidden whitespace-nowrap",
              isCollapsed ? "w-0 opacity-0 ml-0" : "w-auto opacity-100",
            )}
          >
            <h1 className="font-headline-sm text-[15px] font-bold text-slate-800 dark:text-white tracking-tight leading-none">
              {t("common.appName")}
            </h1>
            <p className="font-label-sm text-[10px] text-slate-500 dark:text-slate-400 mt-1 leading-none">
              {t("common.appTagline")}
            </p>
          </div>
        </div>

        {/* Navigation List */}
        <div className="flex-1 overflow-y-auto space-y-1.5 -mx-1 px-1">
          {navItems.map((item) => {
            const isActive =
              pathname === item.href ||
              (item.href !== "/" && pathname.startsWith(item.href));

            return (
              <Link
                key={item.href}
                href={item.href}
                onClick={closeMobile}
                className={cn(
                  "flex items-center justify-between rounded-xl transition-all duration-300 group relative overflow-hidden",
                  isCollapsed ? "px-2 py-2.5 justify-center" : "px-3.5 py-3",
                  isActive
                    ? "bg-[#eff6ff] dark:bg-blue-950/40 text-blue-600 dark:text-blue-400 font-medium"
                    : "text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-slate-100 hover:bg-slate-50/80 dark:hover:bg-slate-800/40",
                )}
                title={isCollapsed ? item.name : undefined}
              >
                <div className="flex items-center gap-3.5 min-w-0">
                  <item.icon
                    className={cn(
                      "h-5 w-5 shrink-0 transition-colors duration-300",
                      isActive
                        ? "text-blue-600 dark:text-blue-400"
                        : "text-slate-400 dark:text-slate-500 group-hover:text-slate-700 dark:group-hover:text-slate-350",
                    )}
                  />
                  <span
                    className={cn(
                      "text-sm font-medium tracking-wide transition-all duration-300 ease-in-out overflow-hidden whitespace-nowrap",
                      isCollapsed ? "w-0 opacity-0 ml-0" : "w-auto opacity-100",
                    )}
                  >
                    {item.name}
                  </span>
                </div>

                {/* Badges */}
                {item.badge && (
                  <div className="relative shrink-0 flex items-center justify-center min-w-[20px] h-5 ml-1">
                    {/* Small Dot Badge (Collapsed) */}
                    <span
                      className={cn(
                        "absolute w-2 h-2 bg-blue-600 dark:bg-blue-500 rounded-full border-2 border-white dark:border-slate-900 transition-all duration-300 shadow-sm",
                        isCollapsed
                          ? "scale-100 opacity-100"
                          : "scale-0 opacity-0",
                      )}
                    />

                    {/* Full Pill Badge (Expanded) */}
                    <span
                      className={cn(
                        "text-[10px] font-bold px-2 py-0.5 rounded-full text-center transition-all duration-300 whitespace-nowrap overflow-hidden block",
                        isCollapsed
                          ? "scale-0 opacity-0 w-0"
                          : "scale-100 opacity-100 w-auto",
                        isActive
                          ? "bg-blue-100/60 dark:bg-blue-950/80 text-blue-700 dark:text-blue-300"
                          : "bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400",
                      )}
                    >
                      {item.badge}
                    </span>
                  </div>
                )}
              </Link>
            );
          })}
        </div>

        {/* Logout Button */}
        <div className="pt-3 mt-3 border-t border-slate-100 dark:border-slate-800/60">
          <Link
            href="/login"
            onClick={handleLogout}
            className={cn(
              "w-full flex items-center text-red-500 hover:text-red-650 hover:bg-red-50/40 dark:hover:bg-red-950/20 rounded-xl transition-all duration-300 cursor-pointer font-medium text-sm overflow-hidden",
              isCollapsed
                ? "justify-center px-2.5 py-2.5"
                : "gap-3.5 px-3.5 py-2.5",
            )}
            title={isCollapsed ? t("nav.logout") : undefined}
          >
            <LogOut className="h-5 w-5 shrink-0" />
            <span
              className={cn(
                "transition-all duration-300 ease-in-out overflow-hidden whitespace-nowrap",
                isCollapsed ? "w-0 opacity-0 ml-0" : "w-auto opacity-100",
              )}
            >
              {t("nav.logout")}
            </span>
          </Link>
        </div>
      </nav>
    </>
  );
}
