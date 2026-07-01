"use client"

import React, { useEffect, useMemo, useState } from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { toast } from "sonner"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"
import { Checkbox } from "@/components/ui/Checkbox"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/Select"
import { useTaskStore } from "@/hooks/useTaskStore"
import { apiRequest } from "@/lib/api"
import { TeamMember, UserProfile } from "@/types/crm"
import { useLanguage } from "@/context/LanguageContext"

function formatDate(date: Date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, "0")
  const day = String(date.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

type QuickCreateInput = {
  title: string
  company: string
  date: string
  time: string
  type: "Meeting" | "Call" | "Proposal" | "Other"
  priority: "Tinggi" | "Sedang" | "Rendah"
  assignee: string
  notes: string
  completed: boolean
}

function buildSchema(t: (k: string) => string) {
  return z.object({
    title: z.string().min(3, t("quickCreate.validation.titleMin")),
    company: z.string().min(2, t("quickCreate.validation.companyMin")),
    date: z
      .string()
      .min(1, t("quickCreate.validation.dateRequired"))
      .regex(/^\d{4}-\d{2}-\d{2}$/, t("quickCreate.validation.dateFormat")),
    time: z
      .string()
      .min(1, t("quickCreate.validation.timeRequired"))
      .regex(/^\d{2}:\d{2}$/, t("quickCreate.validation.timeFormat")),
    type: z.enum(["Meeting", "Call", "Proposal", "Other"]),
    priority: z.enum(["Tinggi", "Sedang", "Rendah"]),
    assignee: z.string().min(1, t("quickCreate.validation.assigneeRequired")),
    notes: z.string().min(5, t("quickCreate.validation.notesMin")),
    completed: z.boolean(),
  })
}

interface QuickCreateDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function QuickCreateDialog({ open, onOpenChange }: QuickCreateDialogProps) {
  const { addTask, isLoading } = useTaskStore()
  const { t, locale } = useLanguage()
  const [assignees, setAssignees] = useState<Array<{ id: string; name: string }>>([])

  const quickCreateSchema = useMemo(() => buildSchema(t), [t])

  const { register, handleSubmit, control, reset, formState: { errors } } = useForm<QuickCreateInput>({
    resolver: zodResolver(quickCreateSchema),
    defaultValues: {
      title: "",
      company: "",
      date: formatDate(new Date()),
      time: "12:00",
      type: "Call",
      priority: "Sedang",
      assignee: "",
      notes: "",
      completed: false,
    },
  })

  useEffect(() => {
    if (!open) return
    const loadAssignees = async () => {
      try {
        const members = await apiRequest<TeamMember[]>("/api/team")
        setAssignees(members.map((member) => ({ id: member.id, name: member.name })))
      } catch {
        try {
          const profile = await apiRequest<UserProfile>("/api/profile")
          setAssignees([{ id: profile.id, name: profile.name }])
        } catch {
          setAssignees([])
        }
      }
    }
    const timer = window.setTimeout(() => {
      void loadAssignees()
    }, 0)
    return () => window.clearTimeout(timer)
  }, [open])

  const handleClose = () => {
    onOpenChange(false)
  }

  const onSubmit = async (data: QuickCreateInput) => {
    try {
      await addTask({
        title: data.title,
        company: data.company,
        time: data.time,
        date: data.date,
        type: data.type,
        priority: data.priority,
        assignee: data.assignee,
        notes: data.notes,
        completed: data.completed,
      })

      const dateObj = new Date(`${data.date}T00:00:00`)
      const months = locale === "id"
        ? ["Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"]
        : ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
      toast.success(t("toast.taskAdded"), {
        description: t("toast.taskAddedDesc", { title: data.title, company: data.company, date: `${dateObj.getDate()} ${months[dateObj.getMonth()]} ${dateObj.getFullYear()}` }),
      })
      reset()
      handleClose()
    } catch {
      toast.error(t("toast.taskSaveError"))
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg max-h-[calc(100vh-2rem)] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <span className="material-symbols-outlined text-primary text-[20px]">add_task</span>
            {t("quickCreate.title")}
          </DialogTitle>
          <DialogDescription>
            {t("quickCreate.description")}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
          {/* Judul & Klien */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
                {t("quickCreate.taskTitle")} <span className="text-red-500">*</span>
              </label>
              <Input
                {...register("title")}
                placeholder={t("quickCreate.taskTitlePlaceholder")}
                className={errors.title ? "border-red-500" : ""}
              />
              {errors.title && (
                <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.title.message}</p>
              )}
            </div>
            <div>
              <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
                {t("quickCreate.relatedTo")} <span className="text-red-500">*</span>
              </label>
              <Input
                {...register("company")}
                placeholder={t("quickCreate.relatedToPlaceholder")}
                className={errors.company ? "border-red-500" : ""}
              />
              {errors.company && (
                <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.company.message}</p>
              )}
            </div>
          </div>

          {/* Tanggal & Waktu */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
                {t("quickCreate.date")} <span className="text-red-500">*</span>
              </label>
              <Input
                type="date"
                {...register("date")}
                className={errors.date ? "border-red-500" : ""}
              />
              {errors.date && (
                <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.date.message}</p>
              )}
            </div>
            <div>
              <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
                {t("quickCreate.time")} <span className="text-red-500">*</span>
              </label>
              <Input
                type="time"
                {...register("time")}
                className={errors.time ? "border-red-500" : ""}
              />
              {errors.time && (
                <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.time.message}</p>
              )}
            </div>
          </div>

          {/* Jenis, Prioritas, Penanggung Jawab */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div>
              <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
                {t("quickCreate.activityType")}
              </label>
              <Controller
                name="type"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder={t("quickCreate.selectType")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Call">{t("values.taskType.Call")}</SelectItem>
                      <SelectItem value="Meeting">{t("values.taskType.Meeting")}</SelectItem>
                      <SelectItem value="Proposal">{t("values.taskType.Proposal")}</SelectItem>
                      <SelectItem value="Other">{t("values.taskType.Other")}</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
            <div>
              <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
                {t("quickCreate.priority")}
              </label>
              <Controller
                name="priority"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder={t("quickCreate.selectPriority")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="Tinggi">{t("values.priority.Tinggi")}</SelectItem>
                      <SelectItem value="Sedang">{t("values.priority.Sedang")}</SelectItem>
                      <SelectItem value="Rendah">{t("values.priority.Rendah")}</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>
            <div>
              <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
                {t("quickCreate.assignee")}
              </label>
              <Controller
                name="assignee"
                control={control}
                render={({ field }) => (
                  <Select value={field.value} onValueChange={field.onChange}>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder={t("quickCreate.selectAssignee")} />
                    </SelectTrigger>
                    <SelectContent>
                      {assignees.map((assignee) => (
                        <SelectItem key={assignee.id} value={assignee.name}>
                          {assignee.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              {errors.assignee && (
                <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.assignee.message}</p>
              )}
            </div>
          </div>

          {/* Catatan */}
          <div>
            <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
              {t("quickCreate.detailNotes")} <span className="text-red-500">*</span>
            </label>
            <textarea
              {...register("notes")}
              className={`w-full rounded-lg border bg-surface-container-lowest py-2 px-3 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none resize-none min-h-[72px] ${
                errors.notes ? "border-red-500" : "border-outline-variant"
              }`}
              placeholder={t("quickCreate.notesPlaceholder")}
              rows={3}
            />
            {errors.notes && (
              <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.notes.message}</p>
            )}
          </div>

          {/* Status Selesai Langsung */}
          <div className="flex items-center gap-3">
            <Controller
              name="completed"
              control={control}
              render={({ field }) => (
                <Checkbox
                  checked={field.value}
                  onCheckedChange={(checked) => field.onChange(checked === true)}
                  id="quick-create-completed"
                  className="cursor-pointer"
                />
              )}
            />
            <label htmlFor="quick-create-completed" className="text-xs font-semibold text-on-surface-variant cursor-pointer select-none">
              {t("quickCreate.markCompleted")}
            </label>
          </div>

          <DialogFooter>
            <DialogClose asChild>
              <Button type="button" variant="ghost" className="cursor-pointer">
                {t("common.cancel")}
              </Button>
            </DialogClose>
            <Button
              type="submit"
              variant="default"
              disabled={isLoading}
              className="flex items-center gap-1.5 cursor-pointer"
            >
              <span className="material-symbols-outlined text-[18px]">save</span>
              {isLoading ? t("common.saving") : t("quickCreate.saveBtn")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
