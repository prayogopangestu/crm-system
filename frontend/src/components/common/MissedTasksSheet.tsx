"use client"

import React, { useEffect } from "react"
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from "@/components/ui/sheet"
import { Checkbox } from "@/components/ui/Checkbox"
import { useTaskStore } from "@/hooks/useTaskStore"
import { toast } from "sonner"

const indonesianMonths = [
  "Januari", "Februari", "Maret", "April", "Mei", "Juni",
  "Juli", "Agustus", "September", "Oktober", "November", "Desember",
]

function formatDateLabel(dateStr: string) {
  const date = new Date(`${dateStr}T00:00:00`)
  if (Number.isNaN(date.getTime())) return dateStr
  return `${date.getDate()} ${indonesianMonths[date.getMonth()]} ${date.getFullYear()}`
}

interface MissedTasksSheetProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function MissedTasksSheet({ open, onOpenChange }: MissedTasksSheetProps) {
  const { tasks, isLoading, loadTasks, toggleTask } = useTaskStore()

  useEffect(() => {
    if (!open) return
    const timer = window.setTimeout(() => {
      void loadTasks()
    }, 0)
    return () => window.clearTimeout(timer)
  }, [open, loadTasks])

  const missedTasks = tasks.filter((t) => t.status === "overdue" && !t.completed)

  const handleComplete = (taskId: string, title: string) => {
    void toggleTask(taskId).then(() => {
      toast.success("Tugas diselesaikan!", {
        description: `Tugas "${title}" telah diselesaikan.`,
      })
    }).catch(() => {
      toast.error("Gagal mengubah status tugas")
    })
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent
        side="right"
        className="w-full sm:max-w-md p-0 flex flex-col bg-surface-container-lowest"
      >
        <SheetHeader className="px-5 pt-5 pb-4 border-b border-outline-variant">
          <SheetTitle className="flex items-center gap-2 text-base font-semibold">
            <span className="material-symbols-outlined text-error text-[20px]">
              notifications_active
            </span>
            Tugas Terlewat
          </SheetTitle>
          <SheetDescription className="text-xs text-on-surface-variant">
            {missedTasks.length > 0
              ? `${missedTasks.length} tugas membutuhkan tindakan segera.`
              : "Tidak ada tugas yang terlewat."}
          </SheetDescription>
        </SheetHeader>

        <div className="flex-1 overflow-y-auto px-5 py-4">
          {isLoading ? (
            <div className="flex flex-col items-center justify-center py-16 text-center">
              <span className="material-symbols-outlined text-on-surface-variant/40 text-[40px] mb-2 animate-spin">
                progress_activity
              </span>
              <p className="text-xs text-on-surface-variant font-medium">
                Memuat tugas...
              </p>
            </div>
          ) : missedTasks.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-center">
              <span className="material-symbols-outlined text-emerald-500/60 text-[44px] mb-3">
                task_alt
              </span>
              <p className="text-sm text-on-surface font-semibold">
                Semua tugas tertangani
              </p>
              <p className="text-[11px] text-on-surface-variant mt-1">
                Tidak ada tugas yang terlewat saat ini.
              </p>
            </div>
          ) : (
            <ul className="flex flex-col divide-y divide-outline-variant/40">
              {missedTasks.map((task) => (
                <li
                  key={task.id}
                  className="py-3.5 first:pt-0 last:pb-0 flex items-start gap-3 group"
                >
                  <Checkbox
                    checked={task.completed}
                    onCheckedChange={() => handleComplete(task.id, task.title)}
                    className="mt-0.5 cursor-pointer accent-red-600"
                  />
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <h4 className="font-semibold text-sm text-on-surface group-hover:text-primary transition-colors leading-snug">
                        {task.title}
                      </h4>
                      <span className="shrink-0 px-2 py-0.5 rounded-full text-[9px] font-bold tracking-wider uppercase bg-red-500/10 text-error">
                        {task.type}
                      </span>
                    </div>
                    <p className="text-xs text-on-surface-variant mt-1 truncate">
                      {task.company}
                    </p>
                    <div className="flex items-center gap-3 mt-2 text-[10px] text-on-surface-variant font-semibold">
                      <span className="flex items-center gap-1">
                        <span className="material-symbols-outlined text-[13px]">
                          event
                        </span>
                        {formatDateLabel(task.date)}
                      </span>
                      <span className="flex items-center gap-1">
                        <span className="material-symbols-outlined text-[13px]">
                          schedule
                        </span>
                        {task.time}
                      </span>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        {missedTasks.length > 0 && (
          <div className="px-5 py-3 border-t border-outline-variant bg-surface-container-low">
            <p className="text-[10px] text-on-surface-variant font-medium text-center">
              Centang tugas untuk menandainya selesai
            </p>
          </div>
        )}
      </SheetContent>
    </Sheet>
  )
}
