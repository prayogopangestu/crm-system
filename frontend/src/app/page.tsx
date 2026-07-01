"use client"

import React, { useState, useEffect, useMemo } from "react"
import { Card } from "@/components/ui/Card"
import { Checkbox } from "@/components/ui/Checkbox"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/Avatar"
import { cn } from "@/lib/utils"
import { apiRequest, queryString } from "@/lib/api"
import {
  Activity,
  ConversionPoint,
  DashboardStats,
  Deal,
  DealPriority,
  LeaderboardEntry,
  PipelineStage,
  Task,
} from "@/types/crm"
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts"
import { useLanguage } from "@/context/LanguageContext"

interface AreaTooltipProps {
  active?: boolean
  payload?: Array<{ value: number }>
  label?: string
}

const CustomTooltip = ({ active, payload, label }: AreaTooltipProps) => {
  const { t } = useLanguage()
  if (active && payload && payload.length) {
    return (
      <div className="bg-surface-container-lowest border border-outline-variant px-3 py-2 rounded-xl shadow-lg text-xs">
        <p className="font-bold text-on-surface mb-1">{label}</p>
        <p className="text-emerald-600 font-bold flex items-center gap-1">
          <span className="material-symbols-outlined text-[14px]">show_chart</span>
          {t("dashboard.conversionLabel", { value: payload[0].value })}
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
    chip: "bg-primary-fixed text-on-primary-fixed",
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
  sky: {
    chip: "bg-sky-500/10 text-sky-600 dark:text-sky-400",
    blob: "bg-sky-400/40",
    glow: "group-hover:shadow-sky-500/10",
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

function MiniStat({
  icon,
  label,
  value,
  sub,
  accent = "text-primary",
}: {
  icon: string
  label: string
  value: string
  sub?: string
  accent?: string
}) {
  return (
    <div className="flex items-center gap-3 rounded-xl border border-outline-variant/50 bg-surface px-3 py-2.5">
      <span
        className={cn(
          "flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-surface-container-low",
          accent
        )}
      >
        <span className="material-symbols-outlined text-[18px]">{icon}</span>
      </span>
      <div className="min-w-0">
        <p className="text-[10px] font-semibold uppercase tracking-wide text-on-surface-variant">
          {label}
        </p>
        <p className="truncate text-sm font-bold text-on-surface">{value}</p>
      </div>
      {sub && (
        <span className="ml-auto flex-shrink-0 text-[11px] font-semibold text-on-surface-variant">
          {sub}
        </span>
      )}
    </div>
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

const PRIORITY_BADGE: Record<DealPriority, string> = {
  High: "bg-error-container text-on-error-container",
  Medium: "bg-secondary-container text-on-secondary-container",
  Low: "bg-surface-variant text-on-surface-variant",
}

const STAGE_META: Record<string, { name: string; color: string }> = {
  lead: { name: "Lead Masuk", color: "#6366f1" },
  contacted: { name: "Dihubungi", color: "#0ea5e9" },
  meeting: { name: "Meeting", color: "#14b8a6" },
  negotiation: { name: "Negosiasi", color: "#f59e0b" },
  won: { name: "Deal Menang", color: "#10b981" },
  lost: { name: "Deal Hilang", color: "#ef4444" },
}

function formatCompact(n: number): string {
  if (n >= 1_000_000_000) return `Rp ${(n / 1_000_000_000).toFixed(1)}M`
  if (n >= 1_000_000) return `Rp ${(n / 1_000_000).toFixed(1)}jt`
  if (n >= 1_000) return `Rp ${(n / 1_000).toFixed(0)}rb`
  return `Rp ${n.toLocaleString("id-ID")}`
}

export default function Dashboard() {
  const { t } = useLanguage()
  const [mounted, setMounted] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [chartData, setChartData] = useState<ConversionPoint[]>([])
  const [activities, setActivities] = useState<Activity[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [deals, setDeals] = useState<Deal[]>([])
  const [stages, setStages] = useState<PipelineStage[]>([])
  const [leaders, setLeaders] = useState<LeaderboardEntry[]>([])

  useEffect(() => {
    const timer = window.setTimeout(() => setMounted(true), 0)
    return () => window.clearTimeout(timer)
  }, [])

  useEffect(() => {
    const loadDashboard = async () => {
      setIsLoading(true)
      setError(null)
      try {
        const [statsResult, chartResult, activityResult, taskResult, dealResult, stageResult, leaderResult] = await Promise.all([
          apiRequest<DashboardStats>("/api/dashboard/stats"),
          apiRequest<ConversionPoint[]>("/api/dashboard/conversion-chart"),
          apiRequest<Activity[]>("/api/dashboard/activities?limit=6"),
          apiRequest<Task[]>("/api/tasks?status=today"),
          apiRequest<Deal[]>("/api/deals"),
          apiRequest<PipelineStage[]>("/api/pipeline-stages"),
          apiRequest<LeaderboardEntry[]>(`/api/reports/leaderboard${queryString({ period: "Bulan Ini" })}`),
        ])
        setStats(statsResult)
        setChartData(chartResult)
        setActivities(activityResult)
        setTasks(taskResult)
        setDeals(dealResult)
        setStages(stageResult)
        setLeaders(leaderResult)
      } catch (error) {
        setError(error instanceof Error ? error.message : t("dashboard.loadError"))
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
    setError(error instanceof Error ? error.message : t("toast.taskToggleError"))
    }
  }

  const urgentCount = stats?.urgentTasksCount || 0

  const avgConversion = chartData.length
    ? Math.round(chartData.reduce((s, d) => s + d.Konversi, 0) / chartData.length)
    : 0
  const latestConversion = chartData.length ? chartData[chartData.length - 1].Konversi : 0

  // Derived deal metrics
  const activeDeals = useMemo(
    () => deals.filter((d) => d.stage !== "won" && d.stage !== "lost"),
    [deals]
  )
  const wonDeals = useMemo(() => deals.filter((d) => d.stage === "won"), [deals])
  const lostDeals = useMemo(() => deals.filter((d) => d.stage === "lost"), [deals])
  const pipelineValue = useMemo(
    () => activeDeals.reduce((sum, d) => sum + d.value, 0),
    [activeDeals]
  )
  const avgDealSize = activeDeals.length ? Math.round(pipelineValue / activeDeals.length) : 0
  const closedDeals = wonDeals.length + lostDeals.length
  const winRate = closedDeals ? Math.round((wonDeals.length / closedDeals) * 100) : 0

  // Top deals by value (active first, then any)
  const hotDeals = useMemo(() => {
    const sorted = [...deals].sort((a, b) => b.value - a.value)
    return sorted.slice(0, 5)
  }, [deals])

  // Pipeline distribution by stage (exclude won/lost for the funnel)
  const stageDistribution = useMemo(() => {
    return stages
      .filter((s) => s.key !== "won" && s.key !== "lost")
      .map((stage) => {
        const items = deals.filter((d) => d.stage === stage.key)
        const total = items.reduce((sum, d) => sum + d.value, 0)
        return {
          key: stage.key,
          name: stage.name || t(`values.stage.${stage.key}`) || stage.key,
          color: STAGE_META[stage.key]?.color || "#6366f1",
          count: items.length,
          total,
        }
      })
  }, [stages, deals, t])
  const maxStageTotal = stageDistribution.reduce((m, s) => Math.max(m, s.total), 0)

  return (
    <div className="p-gutter max-w-container-max-width mx-auto w-full">
      {/* Title Header */}
      <div className="mb-6 flex flex-wrap items-end justify-between gap-3">
        <div>
          <h2 className="font-headline-md text-[24px] text-on-surface font-semibold">
            {t("dashboard.title")}
          </h2>
          <p className="text-on-surface-variant font-body-md text-sm mt-1">
            {t("dashboard.subtitle")}
          </p>
        </div>
        <a
          href="/laporan"
          className="inline-flex items-center gap-1 rounded-lg border border-outline-variant bg-surface px-3 py-2 text-xs font-semibold text-on-surface-variant transition-colors hover:bg-surface-container-low hover:text-on-surface"
        >
          <span className="material-symbols-outlined text-[16px]">insights</span>
          {t("dashboard.fullAnalytics")}
        </a>
      </div>

      {error && (
        <div className="mb-4 rounded-lg border border-error-container bg-error-container/40 px-4 py-3 text-xs font-semibold text-on-error-container">
          {error}
        </div>
      )}

      {/* Metrics Row */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
        <StatCard
          icon="person_add"
          label={t("dashboard.statTotalLeads")}
          value={(stats?.totalLeads || 0).toLocaleString("id-ID")}
          trend={stats?.leadsTrend}
          accent={ACCENTS.primary}
          delay={0}
        />
        <StatCard
          icon="emoji_events"
          label={t("dashboard.statDealWon")}
          value={String(stats?.dealWonCount || 0)}
          trend={stats?.wonTrend}
          accent={ACCENTS.emerald}
          delay={80}
        />
        <StatCard
          icon="payments"
          label={t("dashboard.statTotalRevenue")}
          value={stats?.totalRevenue || "Rp 0"}
          trend={stats?.revenueTrend}
          accent={ACCENTS.amber}
          delay={160}
        />
        <StatCard
          icon="percent"
          label={t("dashboard.statConversion")}
          value={`${avgConversion}%`}
          trend={`${t("dashboard.thisMonth")} ${latestConversion}%`}
          accent={ACCENTS.sky}
          delay={240}
        />
      </div>

      {/* Secondary KPI strip */}
      <div className="mb-6 grid grid-cols-2 lg:grid-cols-4 gap-3">
        <MiniStat
          icon="trending_up"
          label={t("dashboard.statActiveDeals")}
          value={String(activeDeals.length)}
          sub={`${wonDeals.length} ${t("dashboard.won")}`}
          accent="text-primary"
        />
        <MiniStat
          icon="account_balance_wallet"
          label={t("dashboard.statPipelineValue")}
          value={formatCompact(pipelineValue)}
          accent="text-emerald-600 dark:text-emerald-400"
        />
        <MiniStat
          icon="sell"
          label={t("dashboard.statAvgDeal")}
          value={formatCompact(avgDealSize)}
          accent="text-amber-600 dark:text-amber-400"
        />
        <MiniStat
          icon="workspace_premium"
          label={t("dashboard.statWinRate")}
          value={`${winRate}%`}
          sub={`${lostDeals.length} ${t("dashboard.lost")}`}
          accent="text-sky-600 dark:text-sky-400"
        />
      </div>

      {/* Main Content Row */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
        {/* Main Chart Area */}
        <Card className="lg:col-span-2 p-5 flex flex-col w-full">
          <div className="flex flex-wrap justify-between items-start gap-3 mb-5">
            <div>
              <h3 className="font-headline-sm text-lg font-semibold text-on-surface">
                {t("dashboard.conversionTitle")}
              </h3>
              <p className="text-on-surface-variant text-xs mt-0.5">
                {t("dashboard.conversionSubtitle")}
              </p>
            </div>
            <div className="flex items-center gap-2">
              {/* Summary stat chips */}
              <div className="hidden sm:flex items-center gap-3 mr-1">
                <div className="text-right">
                  <p className="text-[10px] uppercase tracking-wide text-on-surface-variant font-semibold">
                    {t("dashboard.average")}
                  </p>
                  <p className="text-sm font-bold text-on-surface">{avgConversion}%</p>
                </div>
                <div className="h-8 w-px bg-outline-variant" />
                <div className="text-right">
                  <p className="text-[10px] uppercase tracking-wide text-on-surface-variant font-semibold">
                    {t("dashboard.thisMonth")}
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
                {t("dashboard.loadingChart")}
              </div>
            )}
          </div>
        </Card>

        {/* Today's Tasks */}
        <Card className="p-5 flex-1 flex flex-col h-full">
          <div className="flex justify-between items-center mb-4">
            <div className="flex items-center gap-2">
              <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                {t("dashboard.todayTasks")}
              </h3>
            </div>
            {urgentCount > 0 && (
              <span className="inline-flex items-center gap-1 bg-error-container text-on-error-container text-[11px] px-2.5 py-1 rounded-full font-semibold">
                <span className="material-symbols-outlined text-[13px]">priority_high</span>
                {t("dashboard.important", { count: urgentCount })}
              </span>
            )}
          </div>
          <ul className="flex flex-col gap-2.5 overflow-y-auto pr-1 -mr-1 max-h-[290px]">
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
                      title={t("dashboard.priorityTitle", { priority: t(`values.priority.${task.priority}`) })}
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
                <p className="text-xs text-on-surface-variant">{t("dashboard.noTasksToday")}</p>
              </li>
            )}
          </ul>
        </Card>
      </div>

      {/* Hot Deals + Top Performers */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
        {/* Hot Deals */}
        <Card className="lg:col-span-2 p-5 flex flex-col w-full">
          <div className="flex justify-between items-center mb-4">
            <div>
              <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                {t("dashboard.hotDealsTitle")}
              </h3>
              <p className="text-on-surface-variant text-xs mt-0.5">
                {t("dashboard.hotDealsSubtitle")}
              </p>
            </div>
            <a href="/pipeline" className="text-primary text-xs font-semibold hover:underline inline-flex items-center gap-0.5">
              {t("nav.pipeline")}
              <span className="material-symbols-outlined text-[14px]">chevron_right</span>
            </a>
          </div>
          <div className="flex flex-col">
            {hotDeals.length === 0 && !isLoading && (
              <div className="py-8 flex flex-col items-center justify-center gap-2 text-center">
                <span className="material-symbols-outlined text-[32px] text-outline-variant">trending_up</span>
                <p className="text-xs text-on-surface-variant">{t("dashboard.noDeals")}</p>
              </div>
            )}
            {hotDeals.map((deal, idx) => {
              const stageMeta = STAGE_META[deal.stage]
              return (
                <div
                  key={deal.id}
                  className={cn(
                    "flex items-center gap-3 py-3 transition-colors hover:bg-surface-container-low/50 -mx-2 px-2 rounded-lg",
                    idx !== hotDeals.length - 1 && "border-b border-outline-variant/40"
                  )}
                >
                  <div className="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-surface-container-low text-xs font-bold text-on-surface-variant">
                    {idx + 1}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-semibold text-on-surface">{deal.title}</p>
                    <p className="truncate text-xs text-on-surface-variant">{deal.company}</p>
                  </div>
                  <span
                    className={cn(
                      "hidden sm:inline-flex flex-shrink-0 items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold uppercase",
                      PRIORITY_BADGE[deal.priority]
                    )}
                  >
                    {deal.priority}
                  </span>
                  {stageMeta && (
                    <span className="hidden md:inline-flex flex-shrink-0 items-center gap-1 text-[11px] font-semibold text-on-surface-variant">
                      <span className="h-2 w-2 rounded-full" style={{ backgroundColor: stageMeta.color }} />
                      {t(`values.stage.${deal.stage}`)}
                    </span>
                  )}
                  <p className="flex-shrink-0 text-sm font-bold text-primary">
                    Rp {deal.value.toLocaleString("id-ID")}
                  </p>
                </div>
              )
            })}
          </div>
        </Card>

        {/* Top Performers */}
        <Card className="p-5 flex flex-col w-full">
          <div className="flex justify-between items-center mb-4">
            <div>
              <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                {t("dashboard.topPerformer")}
              </h3>
              <p className="text-on-surface-variant text-xs mt-0.5">{t("dashboard.thisMonth")}</p>
            </div>
            <a href="/laporan" className="text-primary text-xs font-semibold hover:underline inline-flex items-center gap-0.5">
              {t("dashboard.detail")}
              <span className="material-symbols-outlined text-[14px]">chevron_right</span>
            </a>
          </div>
          <ul className="flex flex-col gap-3">
            {leaders.length === 0 && !isLoading && (
              <li className="py-8 flex flex-col items-center justify-center gap-2 text-center">
                <span className="material-symbols-outlined text-[32px] text-outline-variant">leaderboard</span>
                <p className="text-xs text-on-surface-variant">{t("dashboard.noPerformance")}</p>
              </li>
            )}
            {leaders.slice(0, 5).map((rep) => (
              <li
                key={rep.rank}
                className="flex items-center gap-3 rounded-xl border border-outline-variant/40 bg-surface px-3 py-2.5"
              >
                <span
                  className={cn(
                    "flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full text-xs font-bold",
                    rep.rank === 1
                      ? "bg-amber-500/15 text-amber-600 dark:text-amber-400"
                      : rep.rank === 2
                        ? "bg-surface-container-high text-on-surface-variant"
                        : rep.rank === 3
                          ? "bg-orange-500/15 text-orange-600 dark:text-orange-400"
                          : "bg-surface-container-high text-on-surface-variant"
                  )}
                >
                  {rep.rank}
                </span>
                <Avatar className="h-8 w-8 flex-shrink-0">
                  <AvatarImage src={rep.avatarUrl} alt={rep.name} />
                  <AvatarFallback>{initialsOf(rep.name)}</AvatarFallback>
                </Avatar>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-semibold text-on-surface">{rep.name}</p>
                  <p className="truncate text-[11px] text-on-surface-variant">{rep.role}</p>
                </div>
                <div className="flex-shrink-0 text-right">
                  <p className="text-xs font-bold text-on-surface">{formatCompact(rep.amount)}</p>
                  <p
                    className={cn(
                      "flex items-center justify-end gap-0.5 text-[10px] font-semibold",
                      rep.isPositive ? "text-emerald-600 dark:text-emerald-400" : "text-red-600 dark:text-red-400"
                    )}
                  >
                    <span className="material-symbols-outlined text-[12px]">
                      {rep.isPositive ? "trending_up" : "trending_down"}
                    </span>
                    {rep.trend}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        </Card>
      </div>

      {/* Recent Activities + Pipeline Distribution */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Recent Activities */}
        <Card className="lg:col-span-2 p-5 flex flex-col w-full">
          <div className="flex justify-between items-center mb-4">
            <h3 className="font-headline-sm text-base font-semibold text-on-surface">
              {t("dashboard.recentActivity")}
            </h3>
            <a href="/tugas" className="text-primary text-xs font-semibold hover:underline inline-flex items-center gap-0.5">
              {t("dashboard.viewAll")}
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
                {t("dashboard.noActivity")}
              </li>
            )}
          </ul>
        </Card>

        {/* Pipeline Distribution */}
        <Card className="p-5 flex flex-col w-full">
          <div className="flex justify-between items-center mb-4">
            <div>
              <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                {t("dashboard.pipelineDistribution")}
              </h3>
              <p className="text-on-surface-variant text-xs mt-0.5">
                {t("dashboard.pipelineDistributionSubtitle", { count: activeDeals.length, value: formatCompact(pipelineValue) })}
              </p>
            </div>
          </div>
          <div className="flex flex-col gap-3">
            {stageDistribution.length === 0 && !isLoading && (
              <div className="py-8 flex flex-col items-center justify-center gap-2 text-center">
                <span className="material-symbols-outlined text-[32px] text-outline-variant">funnel</span>
                <p className="text-xs text-on-surface-variant">{t("dashboard.noPipelineData")}</p>
              </div>
            )}
            {stageDistribution.map((stage) => {
              const width = maxStageTotal > 0 ? Math.max((stage.total / maxStageTotal) * 100, 6) : 0
              return (
                <div key={stage.key}>
                  <div className="mb-1 flex items-center justify-between text-xs">
                    <span className="flex items-center gap-1.5 font-semibold text-on-surface">
                      <span className="h-2 w-2 rounded-full" style={{ backgroundColor: stage.color }} />
                      {stage.name}
                    </span>
                    <span className="text-on-surface-variant">
                      <span className="font-bold text-on-surface">{stage.count}</span> · {formatCompact(stage.total)}
                    </span>
                  </div>
                  <div className="h-2 w-full overflow-hidden rounded-full bg-surface-container-high">
                    <div
                      className="h-full rounded-full transition-all duration-500"
                      style={{ width: `${width}%`, backgroundColor: stage.color }}
                    />
                  </div>
                </div>
              )
            })}
          </div>
        </Card>
      </div>
    </div>
  )
}
