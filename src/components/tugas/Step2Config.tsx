"use client"

import React from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTaskStore as useFormStore } from "@/hooks/useTaskStore"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"
import { Checkbox } from "@/components/ui/Checkbox"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"

const step2Schema = z.object({
  type: z.enum(["Meeting", "Call", "Proposal", "Other"]),
  time: z.string().min(1, "Waktu wajib diisi"),
  priority: z.enum(["Tinggi", "Sedang", "Rendah"]),
  assignee: z.string().min(1, "Penanggung jawab wajib dipilih"),
  completedDirectly: z.boolean()
})

type Step2Input = z.infer<typeof step2Schema>

export function Step2Config() {
  const { formData, updateFormData, setStep } = useFormStore()
  
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

  const onSubmit = (data: Step2Input) => {
    updateFormData(data)
    setStep(3) // Lanjut ke step 3
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      {/* Grid 1: Jenis Aktivitas, Waktu, Prioritas */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        {/* Jenis Aktivitas */}
        <div>
          <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
            Jenis Aktivitas
          </label>
          <Controller
            name="type"
            control={control}
            render={({ field }) => (
              <Select value={field.value} onValueChange={field.onChange}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Pilih Jenis" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Call">Panggilan Telepon (Call)</SelectItem>
                  <SelectItem value="Meeting">Pertemuan (Meeting)</SelectItem>
                  <SelectItem value="Proposal">Kirim Proposal</SelectItem>
                  <SelectItem value="Other">Lainnya</SelectItem>
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
            Waktu
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
            Prioritas
          </label>
          <Controller
            name="priority"
            control={control}
            render={({ field }) => (
              <Select value={field.value} onValueChange={field.onChange}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Pilih Prioritas" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Tinggi">Tinggi</SelectItem>
                  <SelectItem value="Sedang">Sedang</SelectItem>
                  <SelectItem value="Rendah">Rendah</SelectItem>
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
            Penanggung Jawab
          </label>
          <Controller
            name="assignee"
            control={control}
            render={({ field }) => (
              <Select value={field.value} onValueChange={field.onChange}>
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Pilih Penanggung Jawab" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Sarah Jenkins">Sarah Jenkins</SelectItem>
                  <SelectItem value="Michael Kusuma">Michael Kusuma</SelectItem>
                  <SelectItem value="Anita Larasati">Anita Larasati</SelectItem>
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
            Tandai tugas ini langsung selesai (Logged Activity)
          </label>
        </div>
      </div>

      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="ghost" onClick={() => setStep(1)} className="cursor-pointer">
          Kembali
        </Button>
        <Button type="submit" variant="default" className="flex items-center gap-1.5 cursor-pointer">
          Lanjut
          <span className="material-symbols-outlined text-[18px]">arrow_forward</span>
        </Button>
      </div>
    </form>
  )
}
