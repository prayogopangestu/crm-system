"use client";

import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/Button";
import { Avatar, AvatarFallback } from "@/components/ui/Avatar";
import { useSidebar } from "@/context/SidebarContext";
import { useTheme } from "next-themes";
import { Sun, Moon } from "lucide-react";
import { useAuthStore } from "@/hooks/useAuthStore";
import { QuickCreateDialog } from "@/components/common/QuickCreateDialog";
import { MissedTasksSheet } from "@/components/common/MissedTasksSheet";

export function Header() {
  const { toggleCollapse, toggleMobile, isCollapsed } = useSidebar();
  const { theme, setTheme } = useTheme();
  const { user } = useAuthStore();
  const [mounted, setMounted] = useState(false);
  const [quickCreateOpen, setQuickCreateOpen] = useState(false);
  const [notificationsOpen, setNotificationsOpen] = useState(false);

  const name = user?.name || "User CRM";
  const initials = name
    .split(" ")
    .map((part) => part[0])
    .join("")
    .substring(0, 2)
    .toUpperCase();

  useEffect(() => {
    const timer = window.setTimeout(() => setMounted(true), 0);
    return () => window.clearTimeout(timer);
  }, []);

  const handleToggle = () => {
    if (window.innerWidth < 768) {
      toggleMobile();
    } else {
      toggleCollapse();
    }
  };

  const handleQuickCreate = () => {
    setQuickCreateOpen(true);
  };

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
        <button
          onClick={() => setNotificationsOpen(true)}
          className="p-2 rounded-full text-on-surface-variant hover:bg-surface-container transition-colors relative cursor-pointer"
          title="Tugas Terlewat"
          aria-label="Tugas Terlewat"
        >
          <span className="material-symbols-outlined">notifications</span>
          <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-error rounded-full border border-surface-container-lowest"></span>
        </button>

        {/* Theme Toggle Button */}
        <button
          onClick={() => {
            document.documentElement.classList.add("theme-transition");
            setTheme(theme === "dark" ? "light" : "dark");
            setTimeout(() => {
              document.documentElement.classList.remove("theme-transition");
            }, 300);
          }}
          className="p-2 rounded-full text-on-surface-variant hover:bg-surface-container transition-colors relative cursor-pointer flex items-center justify-center"
          title={
            mounted && theme === "dark"
              ? "Ubah ke mode terang"
              : "Ubah ke mode gelap"
          }
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
            <AvatarFallback>{initials}</AvatarFallback>
          </Avatar>
          <div className="hidden lg:block text-left">
            <p className="font-label-md text-xs font-semibold text-on-surface group-hover:text-primary transition-colors leading-none">
              {name}
            </p>
            <p className="font-label-sm text-[10px] text-on-surface-variant leading-none mt-1">
              {user?.role || "-"}
            </p>
          </div>
        </div>
      </div>

      <QuickCreateDialog
        open={quickCreateOpen}
        onOpenChange={setQuickCreateOpen}
      />
      <MissedTasksSheet
        open={notificationsOpen}
        onOpenChange={setNotificationsOpen}
      />
    </header>
  );
}
