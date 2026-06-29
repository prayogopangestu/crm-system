"use client"

import React, { useState, useEffect } from "react"
import { Card } from "@/components/ui/Card"
import { Checkbox } from "@/components/ui/Checkbox"
import { cn } from "@/lib/utils"
import { apiRequest } from "@/lib/api"
import { Activity, ConversionPoint, DashboardStats, Task } from "@/types/crm"
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts"

interface AreaTooltipProps {
  active?: boolean
  payload?: Array<{ value: number }>
  label?: string
}

const CustomTooltip = ({ active, payload, label }: AreaTooltipProps) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-surface-container-lowest border border-outline-variant px-3 py-2 rounded-xl shadow-lg text-xs">
        <p className="font-bold text-on-surface mb-1">{label}</p>
        <p className="text-emerald-600 font-bold flex items-center gap-1">
          <span className="material-symbols-outlined text-[14px]">show_chart</span>
          Konversi: {payload[0].value}%
        </p>
      </div>
    )
  }
  return null
}

type TrendDir = "up" | "down" | "flat"

function parseTrend(trend?: string): TrendDir {
  if (!trend) return "flat"
  const t = trend.toLowerCase().trim()
  if (t.startsWith("+") || t.includes("naik") || t.includes("up") || t.includes("tinggi")) {
    return "up"
  }
  if (t.startsWith("-") || t.includes("turun") || t.includes("down") || t.includes("rendah")) {
    return "down"
  }
  return "flat"
}

const ACCENTS = {
  primary: {
    chip: "bg-primary-fixed text-primary",
    blob: "bg-primary-fixed/60",
    glow: "group-hover:shadow-primary/10",
  },
  emerald: {
    chip: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
    blob: "bg-emerald-400/40",
    glow: "group-hover:shadow-emerald-500/10",
  },
  amber: {
    chip: "bg-amber-500/10 text-amber-600 dark:text-amber-400",
    blob: "bg-amber-400/40",
    glow: "group-hover:shadow-amber-500/10",
  },
} as const

function StatCard({
  icon,
  label,
  value,
  trend,
  accent,
  delay = 0,
}: {
  icon: string
  label: string
  value: string
  trend?: string
  accent: (typeof ACCENTS)[keyof typeof ACCENTS]
  delay?: number
}) {
  const dir = parseTrend(trend)
  const trendStyles: Record<TrendDir, string> = {
    up: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
    down: "bg-red-500/10 text-red-600 dark:text-red-400",
    flat: "bg-surface-container-high text-on-surface-variant",
  }
  const trendIcons: Record<TrendDir, string> = {
    up: "trending_up",
    down: "trending_down",
    flat: "trending_flat",
  }

  return (
    <Card
      className="relative overflow-hidden p-5 group transition-all duration-300 hover:-translate-y-1 hover:shadow-xl hover:ring-primary/20"
      style={{ animationDelay: `${delay}ms` }}
    >
      {/* Decorative gradient blob */}
      <div
        className={cn(
          "pointer-events-none absolute -right-10 -top-10 w-36 h-36 rounded-full blur-3xl opacity-70 group-hover:opacity-100 transition-opacity duration-500",
          accent.blob
        )}
      />
      {/* Subtle corner accent line */}
      <div className="pointer-events-none absolute right-0 top-0 h-20 w-20 rounded-bl-[3rem] bg-gradient-to-br from-transparent to-foreground/[0.02]" />

      <div className="relative z-10 flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <div
            className={cn(
              "flex items-center justify-center w-11 h-11 rounded-2xl ring-1 ring-inset ring-foreground/5 transition-transform duration-300 group-hover:scale-105",
              accent.chip
            )}
          >
            <span className="material-symbols-outlined text-[22px]">{icon}</span>
          </div>
        </div>

        <div>
          <p className="text-on-surface-variant text-xs font-semibold tracking-wide uppercase mb-1.5">
            {label}
          </p>
          <h3 className="font-display-lg text-[28px] leading-none font-bold text-on-surface tracking-tight">
            {value}
          </h3>
        </div>

        <div>
          <span
            className={cn(
              "inline-flex items-center gap-1 px-2 py-1 rounded-full text-[11px] font-semibold",
              trendStyles[dir]
            )}
          >
            <span className="material-symbols-outlined text-[14px]">{trendIcons[dir]}</span>
            {trend || "-"}
          </span>
        </div>
      </div>
    </Card>
  )
}

const AVATAR_COLORS = [
  "bg-primary-fixed text-primary",
  "bg-tertiary-fixed text-tertiary dark:text-tertiary-fixed-dim",
  "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
  "bg-amber-500/10 text-amber-600 dark:text-amber-400",
  "bg-sky-500/10 text-sky-600 dark:text-sky-400",
  "bg-rose-500/10 text-rose-600 dark:text-rose-400",
]

function initialsOf(name: string): string {
  return name
    .split(" ")
    .map((w) => w[0])
    .filter(Boolean)
    .slice(0, 2)
    .join("")
    .toUpperCase()
}

function colorFor(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) hash = name.charCodeAt(i) + ((hash << 5) - hash)
  return AVATAR_COLORS[Math.abs(hash) % AVATAR_COLORS.length]
}

const PRIORITY_DOT: Record<string, string> = {
  Tinggi: "bg-red-500",
  Sedang: "bg-amber-500",
  Rendah: "bg-emerald-500",
}

export default function Dashboard() {
  const [mounted, setMounted] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [chartData, setChartData] = useState<ConversionPoint[]>([])
  const [activities, setActivities] = useState<Activity[]>([])
  const [tasks, setTasks] = useState<Task[]>([])

  useEffect(() => {
    const timer = window.setTimeout(() => setMounted(true), 0)
    return () => window.clearTimeout(timer)
  }, [])

  useEffect(() => {
    const loadDashboard = async () => {
      setIsLoading(true)
      setError(null)
      try {
        const [statsResult, chartResult, activityResult, taskResult] = await Promise.all([
          apiRequest<DashboardStats>("/api/dashboard/stats"),
          apiRequest<ConversionPoint[]>("/api/dashboard/conversion-chart"),
          apiRequest<Activity[]>("/api/dashboard/activities?limit=5"),
          apiRequest<Task[]>("/api/tasks?status=today"),
        ])
        setStats(statsResult)
        setChartData(chartResult)
        setActivities(activityResult)
        setTasks(taskResult)
      } catch (error) {
        setError(error instanceof Error ? error.message : "Gagal memuat dashboard")
      } finally {
        setIsLoading(false)
      }
    }
    const timer = window.setTimeout(() => {
      void loadDashboard()
    }, 0)
    return () => window.clearTimeout(timer)
  }, [])

  const handleToggleTask = async (taskId: string) => {
    const target = tasks.find((task) => task.id === taskId)
    if (!target) return
    const completed = !target.completed
    const previous = tasks
    setTasks(prev =>
      prev.map(task =>
        task.id === taskId ? { ...task, completed } : task
      )
    )
    try {
      await apiRequest<{ success: boolean; completed: boolean }>(`/api/tasks/${taskId}/toggle`, {
        method: "PATCH",
        body: { completed },
      })
    } catch (error) {
      setTasks(previous)
      setError(error instanceof Error ? error.message : "Gagal mengubah status tugas")
    }
  }

  const urgentCount = stats?.urgentTasksCount || 0

  const avgConversion = chartData.length
    ? Math.round(chartData.reduce((s, d) => s + d.Konversi, 0) / chartData.length)
    : 0
  const latestConversion = chartData.length ? chartData[chartData.length - 1].Konversi : 0

  return (
    <div className="p-gutter max-w-container-max-width mx-auto w-full">
      {/* Title Header */}
      <div className="mb-8">
        <h2 className="font-headline-md text-[24px] text-on-surface font-semibold">
          Dashboard Overview
        </h2>
        <p className="text-on-surface-variant font-body-md text-sm mt-1">
          Ringkasan performa penjualan dan aktivitas hari ini.
        </p>
      </div>

      {error && (
        <div className="mb-4 rounded-lg border border-error-container bg-error-container/40 px-4 py-3 text-xs font-semibold text-on-error-container">
          {error}
        </div>
      )}

      {/* Metrics Row */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <StatCard
          icon="group"
          label="Total Prospek"
          value={(stats?.totalLeads || 0).toLocaleString("id-ID")}
          trend={stats?.leadsTrend}
          accent={ACCENTS.primary}
          delay={0}
        />
        <StatCard
          icon="emoji_events"
          label="Deal Menang"
          value={String(stats?.dealWonCount || 0)}
          trend={stats?.wonTrend}
          accent={ACCENTS.emerald}
          delay={80}
        />
        <StatCard
          icon="payments"
          label="Total Pendapatan"
          value={stats?.totalRevenue || "Rp 0"}
          trend={stats?.revenueTrend}
          accent={ACCENTS.amber}
          delay={160}
        />
      </div>

      {/* Lower Section Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Left Column */}
        <div className="lg:col-span-2 flex flex-col gap-4">
          {/* Main Chart Area */}
          <Card className="p-5 flex flex-col w-full">
            <div className="flex flex-wrap justify-between items-start gap-3 mb-5">
              <div>
                <h3 className="font-headline-sm text-lg font-semibold text-on-surface">
                  Konversi Penjualan Bulanan
                </h3>
                <p className="text-on-surface-variant text-xs mt-0.5">
                  Persentase konversi prospek menjadi deal
                </p>
              </div>
              <div className="flex items-center gap-2">
                {/* Summary stat chips */}
                <div className="hidden sm:flex items-center gap-3 mr-1">
                  <div className="text-right">
                    <p className="text-[10px] uppercase tracking-wide text-on-surface-variant font-semibold">
                      Rata-rata
                    </p>
                    <p className="text-sm font-bold text-on-surface">{avgConversion}%</p>
                  </div>
                  <div className="h-8 w-px bg-outline-variant" />
                  <div className="text-right">
                    <p className="text-[10px] uppercase tracking-wide text-on-surface-variant font-semibold">
                      Bulan Ini
                    </p>
                    <p className="text-sm font-bold text-emerald-600 dark:text-emerald-400">
                      {latestConversion}%
                    </p>
                  </div>
                </div>
                <button className="text-secondary hover:bg-surface-container p-2 rounded-lg transition-colors cursor-pointer ring-1 ring-transparent hover:ring-outline-variant">
                  <span className="material-symbols-outlined text-[20px]">more_vert</span>
                </button>
              </div>
            </div>
            {/* Recharts Area Chart */}
            <div className="h-[240px] w-full">
              {mounted && !isLoading ? (
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart
                    data={chartData}
                    margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
                  >
                    <defs>
                      <linearGradient id="colorKonversi" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#10b981" stopOpacity={0.35} />
                        <stop offset="95%" stopColor="#10b981" stopOpacity={0.0} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="var(--color-outline-variant, rgba(0,0,0,0.1))" vertical={false} />
                    <XAxis
                      dataKey="name"
                      stroke="var(--color-outline, #6e7881)"
                      fontSize={11}
                      tickLine={false}
                      axisLine={false}
                    />
                    <YAxis
                      stroke="var(--color-outline, #6e7881)"
                      fontSize={11}
                      tickLine={false}
                      axisLine={false}
                      tickFormatter={(value) => `${value}%`}
                    />
                    <Tooltip content={<CustomTooltip />} />
                    <Area
                      type="monotone"
                      dataKey="Konversi"
                      stroke="#10b981"
                      strokeWidth={3}
                      fillOpacity={1}
                      fill="url(#colorKonversi)"
                      dot={false}
                      activeDot={{ r: 5, strokeWidth: 2, stroke: "#fff" }}
                    />
                  </AreaChart>
                </ResponsiveContainer>
              ) : (
                <div className="h-full w-full flex items-center justify-center text-xs text-on-surface-variant font-medium">
                  <span className="material-symbols-outlined animate-spin mr-2 text-[18px]">
                    progress_activity
                  </span>
                  Memuat grafik...
                </div>
              )}
            </div>
          </Card>

          {/* Recent Activities */}
          <Card className="p-5 flex flex-col w-full">
            <div className="flex justify-between items-center mb-4">
              <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                Aktivitas Terbaru
              </h3>
              <a href="/tugas" className="text-primary text-xs font-semibold hover:underline inline-flex items-center gap-0.5">
                Lihat Semua
                <span className="material-symbols-outlined text-[14px]">chevron_right</span>
              </a>
            </div>
            <ul className="flex flex-col">
              {activities.map((activity, idx) => (
                <li
                  key={activity.id}
                  className={cn(
                    "flex items-center gap-3 py-3 transition-colors hover:bg-surface-container-low/50 -mx-2 px-2 rounded-lg",
                    idx !== activities.length - 1 && "border-b border-outline-variant/40"
                  )}
                >
                  <div
                    className={cn(
                      "flex-shrink-0 w-9 h-9 rounded-full flex items-center justify-center text-xs font-bold",
                      colorFor(activity.user)
                    )}
                  >
                    {initialsOf(activity.user)}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs text-on-surface truncate">
                      <span className="font-semibold">{activity.user}</span>{" "}
                      <span className="text-on-surface-variant">{activity.action}</span>{" "}
                      <span className="font-medium text-primary">{activity.target}</span>
                    </p>
                  </div>
                  <div className="flex-shrink-0 flex items-center gap-1.5 text-[11px] text-outline font-medium">
                    <span className={cn("w-1.5 h-1.5 rounded-full", activity.isHighlight ? "bg-primary" : "bg-outline-variant")} />
                    {activity.time}
                  </div>
                </li>
              ))}
              {activities.length === 0 && !isLoading && (
                <li className="py-8 text-center text-xs text-on-surface-variant">
                  Belum ada aktivitas
                </li>
              )}
            </ul>
          </Card>
        </div>

        {/* Right Column */}
        <div className="lg:col-span-1 flex flex-col">
          {/* Today's Tasks */}
          <Card className="p-5 flex-1 flex flex-col h-full">
            <div className="flex justify-between items-center mb-4">
              <div className="flex items-center gap-2">
                <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                  Tugas Hari Ini
                </h3>
              </div>
              {urgentCount > 0 && (
                <span className="inline-flex items-center gap-1 bg-error-container text-on-error-container text-[11px] px-2.5 py-1 rounded-full font-semibold">
                  <span className="material-symbols-outlined text-[13px]">priority_high</span>
                  {urgentCount} Penting
                </span>
              )}
            </div>
            <ul className="flex flex-col gap-2.5 overflow-y-auto pr-1 -mr-1">
              {tasks.map(task => (
                <li
                  key={task.id}
                  className={cn(
                    "group/task flex items-start gap-3 p-3 rounded-xl border bg-surface transition-all",
                    task.completed
                      ? "border-outline-variant/30 opacity-60"
                      : "border-outline-variant/40 hover:border-primary/40 hover:bg-surface-container-low/50"
                  )}
                >
                  <Checkbox
                    checked={task.completed}
                    onChange={() => void handleToggleTask(task.id)}
                    className="mt-0.5"
                  />
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <span
                        className={cn(
                          "w-1.5 h-1.5 rounded-full flex-shrink-0",
                          PRIORITY_DOT[task.priority] || "bg-outline-variant"
                        )}
                        title={`Prioritas ${task.priority}`}
                      />
                      <p
                        className={cn(
                          "font-label-md text-sm text-on-surface truncate",
                          task.completed && "line-through text-on-surface-variant"
                        )}
                      >
                        {task.title}
                      </p>
                    </div>
                    <p className="font-body-md text-xs text-on-surface-variant mt-0.5 truncate pl-3.5">
                      {task.company}
                    </p>
                  </div>
                </li>
              ))}
              {tasks.length === 0 && !isLoading && (
                <li className="py-10 flex flex-col items-center justify-center text-center gap-2">
                  <span className="material-symbols-outlined text-[32px] text-outline-variant">
                    task_alt
                  </span>
                  <p className="text-xs text-on-surface-variant">Tidak ada tugas hari ini</p>
                </li>
              )}
            </ul>
          </Card>
        </div>
      </div>
    </div>
  )
}
