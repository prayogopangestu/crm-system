"use client"

import React from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTaskStore as useFormStore } from "@/hooks/useTaskStore"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"

const step1Schema = z.object({
  title: z.string().min(3, "Judul tugas minimal harus 3 karakter"),
  relatedTo: z.string().min(2, "Nama perusahaan/klien wajib diisi")
})

type Step1Input = z.infer<typeof step1Schema>

export function Step1Basic() {
  const { formData, updateFormData, setStep, resetForm } = useFormStore()
  
  const { register, handleSubmit, formState: { errors } } = useForm<Step1Input>({
    resolver: zodResolver(step1Schema),
    defaultValues: {
      title: formData.title,
      relatedTo: formData.relatedTo
    }
  })

  const onSubmit = (data: Step1Input) => {
    updateFormData(data)
    setStep(2) // Lanjut ke step 2
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
            Judul Tugas <span className="text-red-500">*</span>
          </label>
          <Input
            {...register("title")}
            placeholder="Contoh: Telepon Budi Wijaya"
            className={errors.title ? "border-red-500" : ""}
          />
          {errors.title && (
            <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.title.message}</p>
          )}
        </div>
        <div>
          <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
            Terkait Dengan (Klien/Perusahaan) <span className="text-red-500">*</span>
          </label>
          <Input
            {...register("relatedTo")}
            placeholder="Contoh: PT Telkomsel"
            className={errors.relatedTo ? "border-red-500" : ""}
          />
          {errors.relatedTo && (
            <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.relatedTo.message}</p>
          )}
        </div>
      </div>

      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="ghost" onClick={resetForm} className="cursor-pointer">
          Batal
        </Button>
        <Button type="submit" variant="default" className="flex items-center gap-1.5 cursor-pointer">
          Lanjut
          <span className="material-symbols-outlined text-[18px]">arrow_forward</span>
        </Button>
      </div>
    </form>
  )
}
