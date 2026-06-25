"use client"

import React from "react"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTaskStore as useFormStore, TaskFormData } from "@/hooks/useTaskStore"
import { Button } from "@/components/ui/Button"

const step3Schema = z.object({
  notes: z.string().min(5, "Catatan detail minimal harus 5 karakter")
})

type Step3Input = z.infer<typeof step3Schema>

interface Step3SummaryProps {
  onSave: (data: TaskFormData) => void
  onCancel: () => void
}

export function Step3Summary({ onSave, onCancel }: Step3SummaryProps) {
  const { formData, updateFormData, setStep } = useFormStore()

  const { register, handleSubmit, formState: { errors } } = useForm<Step3Input>({
    resolver: zodResolver(step3Schema),
    defaultValues: {
      notes: formData.notes
    }
  })

  const onSubmit = (data: Step3Input) => {
    // Update notes in store
    updateFormData(data)
    
    // Call the parent save function with consolidated store data
    const finalData: TaskFormData = {
      ...formData,
      notes: data.notes
    }
    onSave(finalData)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
      {/* Summary Display Card */}
      <div className="bg-surface-container-low border border-outline-variant/40 rounded-lg p-4 text-xs space-y-2 animate-in fade-in duration-200">
        <h4 className="font-semibold text-on-surface text-sm pb-1.5 border-b border-outline-variant/30">
          Ringkasan Tugas
        </h4>
        <div className="grid grid-cols-2 gap-y-2 gap-x-4 pt-1">
          <div>
            <span className="text-on-surface-variant block font-medium">Judul Tugas:</span>
            <span className="text-on-surface font-semibold">{formData.title || "-"}</span>
          </div>
          <div>
            <span className="text-on-surface-variant block font-medium">Terkait Dengan:</span>
            <span className="text-on-surface font-semibold">{formData.relatedTo || "-"}</span>
          </div>
          <div>
            <span className="text-on-surface-variant block font-medium">Jenis Aktivitas:</span>
            <span className="text-on-surface font-semibold">{formData.type}</span>
          </div>
          <div>
            <span className="text-on-surface-variant block font-medium">Waktu:</span>
            <span className="text-on-surface font-semibold">{formData.time}</span>
          </div>
          <div>
            <span className="text-on-surface-variant block font-medium">Prioritas:</span>
            <span className="text-on-surface font-semibold">{formData.priority}</span>
          </div>
          <div>
            <span className="text-on-surface-variant block font-medium">Penanggung Jawab:</span>
            <span className="text-on-surface font-semibold">{formData.assignee}</span>
          </div>
          <div className="col-span-2">
            <span className="text-on-surface-variant block font-medium">Status Selesai Langsung:</span>
            <span className="text-on-surface font-semibold">
              {formData.completedDirectly ? "Ya (Langsung Selesai)" : "Tidak (Belum Selesai)"}
            </span>
          </div>
        </div>
      </div>

      {/* Catatan Detail */}
      <div>
        <label className="block font-label-sm text-xs font-semibold text-on-surface-variant mb-1">
          Catatan Detail <span className="text-red-500">*</span>
        </label>
        <textarea
          {...register("notes")}
          className={`w-full rounded-lg border bg-surface-container-lowest py-2 px-3 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none resize-none min-h-[80px] ${
            errors.notes ? "border-red-500" : "border-outline-variant"
          }`}
          placeholder="Tulis instruksi detail, nomor telepon, ringkasan agenda..."
          rows={3}
        />
        {errors.notes && (
          <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.notes.message}</p>
        )}
      </div>

      {/* Action Buttons */}
      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="ghost" onClick={onCancel} className="cursor-pointer">
          Batal
        </Button>
        <Button type="button" variant="outline" onClick={() => setStep(2)} className="cursor-pointer">
          Kembali
        </Button>
        <Button type="submit" variant="default" className="flex items-center gap-1.5 cursor-pointer">
          <span className="material-symbols-outlined text-[18px]">save</span>
          Simpan Tugas
        </Button>
      </div>
    </form>
  )
}
