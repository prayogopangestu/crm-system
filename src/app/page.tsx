"use client"

import React, { useState, useEffect } from "react"
import { Card } from "@/components/ui/Card"
import { Checkbox } from "@/components/ui/Checkbox"
import { initialTasks, initialActivities, initialDeals, initialContacts } from "@/lib/mockData"
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts"

const chartData = [
  { name: "Jan", Konversi: 40 },
  { name: "Feb", Konversi: 55 },
  { name: "Mar", Konversi: 45 },
  { name: "Apr", Konversi: 70 },
  { name: "Mei", Konversi: 90 },
  { name: "Jun", Konversi: 85 },
]

const CustomTooltip = ({ active, payload, label }: any) => {
  if (active && payload && payload.length) {
    return (
      <div className="bg-surface-container-lowest border border-outline-variant p-3 rounded-lg shadow-md text-xs">
        <p className="font-bold text-on-surface mb-1">{label}</p>
        <p className="text-emerald-600 font-bold">
          Konversi: {payload[0].value}%
        </p>
      </div>
    )
  }
  return null
}

export default function Dashboard() {
  const [mounted, setMounted] = useState(false)
  useEffect(() => {
    setMounted(true)
  }, [])

  const [tasks, setTasks] = useState(initialTasks.filter(t => t.status === "today"))
  const [activities] = useState(initialActivities)

  const handleToggleTask = (taskId: string) => {
    setTasks(prev =>
      prev.map(task =>
        task.id === taskId ? { ...task, completed: !task.completed } : task
      )
    )
  }

  // Calculate stats dynamically
  const totalLeads = initialContacts.length + 1232 // Dynamic base + list size
  const dealWonCount = initialDeals.filter(d => d.stage === "won").length + 38
  const totalRevenueText = "Rp 1.2M"

  const urgentCount = tasks.filter(t => !t.completed && t.type === "Proposal" || t.type === "Meeting").length

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

      {/* Metrics Row (Bento Grid Style) */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        {/* Card 1: Total Prospek */}
        <Card className="p-4 relative overflow-hidden group">
          <div className="absolute -right-4 -top-4 w-24 h-24 bg-primary-fixed/20 rounded-full blur-2xl group-hover:bg-primary-fixed/30 transition-colors"></div>
          <div className="flex justify-between items-start mb-3 relative z-10">
            <span className="text-on-surface-variant text-xs font-semibold">Total Prospek</span>
            <div className="w-8 h-8 rounded-full bg-secondary-container flex items-center justify-center text-primary">
              <span className="material-symbols-outlined text-[18px]">group</span>
            </div>
          </div>
          <div className="relative z-10">
            <h3 className="font-display-lg text-2xl font-bold text-on-surface">
              {totalLeads.toLocaleString("id-ID")}
            </h3>
            <div className="flex items-center gap-1 mt-2 text-primary text-xs font-semibold">
              <span className="material-symbols-outlined text-[16px]">trending_up</span>
              <span>+12% bulan ini</span>
            </div>
          </div>
        </Card>

        {/* Card 2: Deal Menang */}
        <Card className="p-4 relative overflow-hidden group">
          <div className="absolute -right-4 -top-4 w-24 h-24 bg-tertiary-fixed/20 rounded-full blur-2xl group-hover:bg-tertiary-fixed/30 transition-colors"></div>
          <div className="flex justify-between items-start mb-3 relative z-10">
            <span className="text-on-surface-variant text-xs font-semibold">Deal Menang</span>
            <div className="w-8 h-8 rounded-full bg-tertiary-fixed flex items-center justify-center text-on-tertiary-container">
              <span className="material-symbols-outlined text-[18px]">workspace_premium</span>
            </div>
          </div>
          <div className="relative z-10">
            <h3 className="font-display-lg text-2xl font-bold text-on-surface">
              {dealWonCount}
            </h3>
            <div className="flex items-center gap-1 mt-2 text-primary text-xs font-semibold">
              <span className="material-symbols-outlined text-[16px]">trending_up</span>
              <span>+5 dari minggu lalu</span>
            </div>
          </div>
        </Card>

        {/* Card 3: Total Pendapatan */}
        <Card className="p-4 relative overflow-hidden group">
          <div className="absolute -right-4 -top-4 w-24 h-24 bg-secondary-fixed/30 rounded-full blur-2xl group-hover:bg-secondary-fixed/50 transition-colors"></div>
          <div className="flex justify-between items-start mb-3 relative z-10">
            <span className="text-on-surface-variant text-xs font-semibold">Total Pendapatan</span>
            <div className="w-8 h-8 rounded-full bg-surface-variant flex items-center justify-center text-on-surface">
              <span className="material-symbols-outlined text-[18px]">payments</span>
            </div>
          </div>
          <div className="relative z-10">
            <h3 className="font-display-lg text-2xl font-bold text-on-surface">
              {totalRevenueText}
            </h3>
            <div className="flex items-center gap-1 mt-2 text-outline text-xs font-semibold">
              <span className="material-symbols-outlined text-[16px]">trending_flat</span>
              <span>Stabil</span>
            </div>
          </div>
        </Card>
      </div>

      {/* Lower Section Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {/* Left Column (Main Chart Area & Recent Activities) */}
        <div className="lg:col-span-2 flex flex-col gap-4">
          {/* Main Chart Area */}
          <Card className="p-4 flex flex-col w-full">
            <div className="flex justify-between items-center mb-6">
              <h3 className="font-headline-sm text-lg font-semibold text-on-surface">
                Konversi Penjualan Bulanan
              </h3>
              <button className="text-secondary hover:bg-surface-container p-2 rounded-md transition-colors cursor-pointer">
                <span className="material-symbols-outlined text-[20px]">more_vert</span>
              </button>
            </div>
            {/* Recharts Area Chart */}
            <div className="h-[220px] w-full bg-surface-container-low rounded-lg border border-surface-variant p-4 flex items-center justify-center">
              {mounted ? (
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart
                    data={chartData}
                    margin={{ top: 15, right: 10, left: -20, bottom: 0 }}
                  >
                    <defs>
                      <linearGradient id="colorKonversi" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#10b981" stopOpacity={0.3} />
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
                    />
                  </AreaChart>
                </ResponsiveContainer>
              ) : (
                <div className="text-xs text-on-surface-variant font-medium">Memuat grafik...</div>
              )}
            </div>
          </Card>

          {/* Recent Activities (Panjang Horisontal) */}
          <Card className="p-4 flex flex-col w-full">
            <div className="flex justify-between items-center mb-4">
              <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                Aktivitas Terbaru
              </h3>
              <a href="/tugas" className="text-primary text-xs font-semibold hover:underline">
                Lihat Semua
              </a>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-surface-variant text-xs text-on-surface-variant font-semibold">
                    <th className="pb-2">Pengguna</th>
                    <th className="pb-2">Aktivitas</th>
                    <th className="pb-2">Target</th>
                    <th className="pb-2 text-right">Waktu</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-surface-variant/40 text-xs">
                  {activities.map((activity) => (
                    <tr key={activity.id} className="hover:bg-surface-container-low/30 transition-colors">
                      <td className="py-2.5 font-semibold text-on-surface flex items-center gap-2">
                        <div className={`w-2 h-2 rounded-full ${activity.isHighlight ? "bg-primary" : "bg-outline-variant"}`}></div>
                        {activity.user}
                      </td>
                      <td className="py-2.5 text-on-surface-variant">{activity.action}</td>
                      <td className="py-2.5 font-medium text-primary">{activity.target}</td>
                      <td className="py-2.5 text-outline text-right">{activity.time}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </Card>
        </div>

        {/* Right Column (Tinggi/Tall) */}
        <div className="lg:col-span-1 flex flex-col">
          {/* Today's Tasks */}
          <Card className="p-4 flex-1 flex flex-col h-full justify-between">
            <div>
              <div className="flex justify-between items-center mb-4">
                <h3 className="font-headline-sm text-base font-semibold text-on-surface">
                  Tugas Hari Ini
                </h3>
                {urgentCount > 0 && (
                  <span className="bg-error-container text-on-error-container text-[11px] px-2.5 py-0.5 rounded-full font-semibold">
                    {urgentCount} Penting
                  </span>
                )}
              </div>
              <ul className="flex flex-col gap-3 overflow-y-auto pr-1">
                {tasks.map(task => (
                  <li
                    key={task.id}
                    className={`flex items-start gap-3 p-3 bg-surface border rounded-lg transition-all ${
                      task.completed
                        ? "border-outline-variant/35 opacity-60"
                        : "border-outline-variant/50 hover:border-primary/30"
                    }`}
                  >
                    <Checkbox
                      checked={task.completed}
                      onChange={() => handleToggleTask(task.id)}
                      className="mt-1"
                    />
                    <div className="min-w-0 flex-1">
                      <p className={`font-label-md text-sm text-on-surface ${task.completed ? "line-through text-on-surface-variant" : ""}`}>
                        {task.title}
                      </p>
                      <p className="font-body-md text-xs text-on-surface-variant mt-0.5 truncate">
                        {task.company}
                      </p>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
