"use client"

import React, { useState, useEffect } from "react"
import { Card } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/Table"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/Avatar"
import { apiDownload, apiRequest, queryString } from "@/lib/api"
import { LeaderboardEntry, LostReason, PerformanceGoal } from "@/types/crm"
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  ResponsiveContainer,
} from "recharts"
import { toast } from "sonner"

interface PieTooltipProps {
  active?: boolean
  payload?: Array<{ payload: LostReason }>
}

const CustomPieTooltip = ({ active, payload }: PieTooltipProps) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload
    return (
      <div className="bg-surface-container-lowest border border-outline-variant p-3 rounded-lg shadow-md text-xs">
        <p className="font-bold text-on-surface mb-1">{data.name}</p>
        <p className="text-primary font-semibold">
          Deals: {data.value} ({data.percentage}%)
        </p>
      </div>
    )
  }
  return null
}

export default function LaporanPage() {
  const [mounted, setMounted] = useState(false)
  const [leaderboardData, setLeaderboardData] = useState<LeaderboardEntry[]>([])
  const [pieData, setPieData] = useState<LostReason[]>([])
  const [performanceGoals, setPerformanceGoals] = useState<PerformanceGoal[]>([])
  const [error, setError] = useState<string | null>(null)
  useEffect(() => {
    const timer = window.setTimeout(() => setMounted(true), 0)
    return () => window.clearTimeout(timer)
  }, [])

  const [period, setPeriod] = useState("Bulan Ini")

  useEffect(() => {
    const loadReports = async () => {
      setError(null)
      try {
        const [leaderboard, lostReasons, goals] = await Promise.all([
          apiRequest<LeaderboardEntry[]>(`/api/reports/leaderboard${queryString({ period })}`),
          apiRequest<LostReason[]>("/api/reports/lost-reasons"),
          apiRequest<PerformanceGoal[]>("/api/reports/goals"),
        ])
        setLeaderboardData(leaderboard)
        setPieData(lostReasons)
        setPerformanceGoals(goals)
      } catch (error) {
        setError(error instanceof Error ? error.message : "Gagal memuat laporan")
      }
    }
    const timer = window.setTimeout(() => {
      void loadReports()
    }, 0)
    return () => window.clearTimeout(timer)
  }, [period])

  const handleExportCSV = async () => {
    try {
      await apiDownload("/api/reports/export/csv", "crm-report.csv")
      toast.success("Data laporan diekspor sebagai CSV")
    } catch {
      toast.error("Export CSV gagal")
    }
  }

  const handleExportPDF = async () => {
    try {
      await apiDownload("/api/reports/export/pdf", "crm-report.pdf")
      toast.success("Dokumen laporan diunduh sebagai PDF")
    } catch {
      toast.error("Export PDF gagal")
    }
  }

  const totalLostDeals = pieData.reduce((sum, item) => sum + item.value, 0)

  return (
    <div className="p-gutter max-w-container-max-width mx-auto w-full">
      {/* Header & Actions */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-stack-lg">
        <div>
          <h2 className="font-headline-md text-[24px] font-semibold text-on-surface">
            Laporan &amp; Analitik
          </h2>
          <p className="text-sm text-on-surface-variant mt-1">
            Tinjauan performa tim dan analisis penjualan bulan ini.
          </p>
        </div>
        <div className="flex gap-3">
          <Button
            onClick={handleExportCSV}
            variant="outline"
            className="flex items-center gap-2 text-sm font-semibold whitespace-nowrap cursor-pointer"
          >
            <span className="material-symbols-outlined text-[18px]">file_download</span>
            Export CSV
          </Button>
          <Button
            onClick={handleExportPDF}
            variant="default"
            className="flex items-center gap-2 text-sm font-semibold whitespace-nowrap cursor-pointer"
          >
            <span className="material-symbols-outlined text-[18px]">picture_as_pdf</span>
            Export PDF
          </Button>
        </div>
      </div>

      {error && (
        <div className="mb-4 rounded-lg border border-error-container bg-error-container/40 px-4 py-3 text-xs font-semibold text-on-error-container">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-6">
        {/* Leaderboard (Team Performance) */}
        <Card className="lg:col-span-2 p-4 flex flex-col">
          <div className="flex items-center justify-between mb-6">
            <h3 className="font-headline-sm text-base font-semibold text-on-surface">
              Performa Tim (Leaderboard)
            </h3>
            <div className="w-36">
              <Select
                value={period}
                onValueChange={setPeriod}
              >
                <SelectTrigger className="w-full text-xs">
                  <SelectValue placeholder="Pilih Periode" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Bulan Ini">Bulan Ini</SelectItem>
                  <SelectItem value="Bulan Lalu">Bulan Lalu</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <div className="flex-1">
            <ul className="flex flex-col gap-4">
              {leaderboardData.length === 0 && (
                <li className="p-4 text-center text-xs text-on-surface-variant">
                  Belum ada data performa.
                </li>
              )}
              {leaderboardData.map((rep) => (
                <li
                  key={rep.rank}
                  className="flex items-center justify-between p-3 rounded-lg bg-surface-container-low border border-surface-variant hover:shadow-sm transition-all"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-full bg-primary/10 text-primary flex items-center justify-center font-bold text-sm">
                      {rep.rank}
                    </div>
                    <Avatar className="h-10 w-10">
                      <AvatarImage src={rep.avatarUrl} alt={rep.name} />
                      <AvatarFallback>{rep.name.split(" ").map(n => n[0]).join("")}</AvatarFallback>
                    </Avatar>
                    <div>
                      <p className="font-label-md text-sm font-semibold text-on-surface">{rep.name}</p>
                      <p className="text-xs text-on-surface-variant">{rep.role}</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-label-md text-sm font-bold text-on-surface">
                      Rp {rep.amount.toLocaleString("id-ID")}
                    </p>
                    <p className={`text-xs flex items-center justify-end gap-1 font-semibold ${
                      rep.isPositive ? "text-green-600" : "text-error"
                    }`}>
                      <span className="material-symbols-outlined text-[14px]">
                        {rep.isPositive ? "trending_up" : "trending_down"}
                      </span>{" "}
                      {rep.trend}
                    </p>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </Card>

        {/* Deal Lost Reason (Recharts Pie Chart) */}
        <Card className="p-4 flex flex-col">
          <h3 className="font-headline-sm text-base font-semibold text-on-surface mb-6">
            Alasan Deal Lost
          </h3>
          <div className="flex-1 flex flex-col items-center justify-center">
            {mounted ? (
              <div className="relative w-32 h-32 flex items-center justify-center">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={pieData}
                      cx="50%"
                      cy="50%"
                      innerRadius={35}
                      outerRadius={55}
                      paddingAngle={3}
                      dataKey="value"
                    >
                      {pieData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip content={<CustomPieTooltip />} />
                  </PieChart>
                </ResponsiveContainer>
                {/* Center Text inside Donut Chart */}
                <div className="absolute flex flex-col items-center justify-center pointer-events-none">
                  <span className="font-headline-md text-xl font-bold text-on-surface leading-none">{totalLostDeals}</span>
                  <span className="text-[8px] text-on-surface-variant font-bold uppercase tracking-wider mt-1">Total Deals</span>
                </div>
              </div>
            ) : (
              <div className="w-32 h-32 flex items-center justify-center text-xs text-on-surface-variant font-medium">
                Memuat grafik...
              </div>
            )}

            {/* Legends */}
            <div className="w-full mt-6 flex flex-col gap-2.5">
              {pieData.length === 0 && (
                <div className="text-xs text-on-surface-variant text-center">Belum ada deal lost.</div>
              )}
              {pieData.map((item, idx) => (
                <div key={idx} className="flex items-center justify-between text-xs font-semibold">
                  <div className="flex items-center gap-2">
                    <div className="w-3 h-3 rounded-full" style={{ backgroundColor: item.color }}></div>
                    <span className="text-on-surface-variant">{item.name}</span>
                  </div>
                  <span className="text-on-surface font-bold">{item.percentage}%</span>
                </div>
              ))}
            </div>
          </div>
        </Card>
      </div>

      {/* Target vs Pencapaian Table */}
      <Card className="p-4">
        <h3 className="font-headline-sm text-base font-semibold text-on-surface mb-4">
          Target vs Pencapaian Bulanan
        </h3>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Bulan</TableHead>
              <TableHead>Target (Goal)</TableHead>
              <TableHead>Pencapaian (Actual)</TableHead>
              <TableHead>Status Pencapaian</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {performanceGoals.length === 0 && (
              <TableRow>
                <TableCell colSpan={4} className="text-center py-8 text-on-surface-variant/70">
                  Belum ada target performa.
                </TableCell>
              </TableRow>
            )}
            {performanceGoals.map((goal, idx) => {
              const isAchieved = goal.percentage >= 100
              return (
                <TableRow key={idx}>
                  <TableCell className="font-semibold text-sm">{goal.month}</TableCell>
                  <TableCell className="text-sm text-on-surface-variant">
                    Rp {goal.goal.toLocaleString("id-ID")}
                  </TableCell>
                  <TableCell className="text-sm font-bold text-on-surface">
                    Rp {goal.actual.toLocaleString("id-ID")}
                  </TableCell>
                  <TableCell>
                    <span className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold ${
                      isAchieved
                        ? "bg-green-100 text-green-800"
                        : "bg-red-100 text-red-800"
                    }`}>
                      {goal.status}
                    </span>
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
      </Card>
    </div>
  )
}
