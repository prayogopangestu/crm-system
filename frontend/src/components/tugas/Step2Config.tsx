"use client"

import React, { useEffect, useMemo, useState } from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTaskStore as useFormStore } from "@/hooks/useTaskStore"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"
import { Checkbox } from "@/components/ui/Checkbox"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"
import { apiRequest } from "@/lib/api"
import { TeamMember, UserProfile } from "@/types/crm"
import { useLanguage } from "@/context/LanguageContext"

type Step2Input = {
  type: "Meeting" | "Call" | "Proposal" | "Other"
  time: string
  priority: "Tinggi" | "Sedang" | "Rendah"
  assignee: string
  completedDirectly: boolean
}

export function Step2Config() {
  const { formData, updateFormData, setStep } = useFormStore()
  const { t } = useLanguage()
  const [assignees, setAssignees] = useState<Array<{ id: string; name: string }>>([])

  const step2Schema = useMemo(
    () =>
      z.object({
        type: z.enum(["Meeting", "Call", "Proposal", "Other"]),
        time: z.string().min(1, t("quickCreate.validation.timeRequired")),
        priority: z.enum(["Tinggi", "Sedang", "Rendah"]),
        assignee: z.string().min(1, t("quickCreate.validation.assigneeRequired")),
        completedDirectly: z.boolean(),
      }),
    [t],
  )

  const { register, control, handleSubmit, formState: { errors } } = useForm<Step2Input>({
    resolver: zodResolver(step2Schema),
    defaultValues: {
      type: formData.type,
      time: formData.time,
      priority: formData.priority,
      assignee: formData.assignee,
      completedDirectly: formData.completedDirectly
    }
  })

  useEffect(() => {
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
  }, [])

  const onSubmit = (data: Step2Input) => {
    updateFormData(data)
    setStep(3)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      {/* Grid 1: Jenis Aktivitas, Waktu, Prioritas */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        {/* Jenis Aktivitas */}
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
          {errors.type && (
            <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.type.message}</p>
          )}
        </div>

        {/* Waktu */}
        <div>
          <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
            {t("quickCreate.time")}
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

        {/* Prioritas */}
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
          {errors.priority && (
            <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.priority.message}</p>
          )}
        </div>
      </div>

      {/* Grid 2: Penanggung Jawab & Status Selesai Langsung */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 items-center">
        {/* Penanggung Jawab */}
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

        {/* Status Selesai Langsung */}
        <div className="flex items-center gap-3 pt-5">
          <Controller
            name="completedDirectly"
            control={control}
            render={({ field }) => (
              <Checkbox
                checked={field.value}
                onCheckedChange={(checked) => field.onChange(checked === true)}
                id="task-completed-directly"
                className="cursor-pointer"
              />
            )}
          />
          <label htmlFor="task-completed-directly" className="text-xs font-semibold text-on-surface-variant cursor-pointer select-none">
            {t("quickCreate.markCompleted")}
          </label>
        </div>
      </div>

      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="ghost" onClick={() => setStep(1)} className="cursor-pointer">
          {t("common.back")}
        </Button>
        <Button type="submit" variant="default" className="flex items-center gap-1.5 cursor-pointer">
          {t("common.next")}
          <span className="material-symbols-outlined text-[18px]">arrow_forward</span>
        </Button>
      </div>
    </form>
  )
}
