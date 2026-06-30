"use client"

import React, { useEffect, useState } from "react"
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

const indonesianMonths = [
  "Januari", "Februari", "Maret", "April", "Mei", "Juni",
  "Juli", "Agustus", "September", "Oktober", "November", "Desember",
]

function formatDate(date: Date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, "0")
  const day = String(date.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

const quickCreateSchema = z.object({
  title: z.string().min(3, "Judul tugas minimal 3 karakter"),
  company: z.string().min(2, "Nama klien/perusahaan minimal 2 karakter"),
  date: z
    .string()
    .min(1, "Tanggal wajib diisi")
    .regex(/^\d{4}-\d{2}-\d{2}$/, "Format tanggal harus YYYY-MM-DD"),
  time: z
    .string()
    .min(1, "Waktu wajib diisi")
    .regex(/^\d{2}:\d{2}$/, "Format waktu harus HH:MM"),
  type: z.enum(["Meeting", "Call", "Proposal", "Other"]),
  priority: z.enum(["Tinggi", "Sedang", "Rendah"]),
  assignee: z.string().min(1, "Penanggung jawab wajib dipilih"),
  notes: z.string().min(5, "Catatan minimal 5 karakter"),
  completed: z.boolean(),
})

type QuickCreateInput = z.infer<typeof quickCreateSchema>

interface QuickCreateDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function QuickCreateDialog({ open, onOpenChange }: QuickCreateDialogProps) {
  const { addTask, isLoading } = useTaskStore()
  const [assignees, setAssignees] = useState<Array<{ id: string; name: string }>>([])

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
      toast.success("Tugas berhasil ditambahkan!", {
        description: `Tugas: "${data.title}" | Klien: ${data.company} untuk tanggal ${dateObj.getDate()} ${indonesianMonths[dateObj.getMonth()]} ${dateObj.getFullYear()}`,
      })
      reset()
      handleClose()
    } catch {
      toast.error("Tugas gagal disimpan")
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg max-h-[calc(100vh-2rem)] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <span className="material-symbols-outlined text-primary text-[20px]">add_task</span>
            Quick Create Tugas
          </DialogTitle>
          <DialogDescription>
            Buat tugas baru secara instan. Lengkapi data di bawah lalu simpan.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
          {/* Judul & Klien */}
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
                {...register("company")}
                placeholder="Contoh: PT Telkomsel"
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
                Tanggal <span className="text-red-500">*</span>
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
                Waktu <span className="text-red-500">*</span>
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
            </div>
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
            </div>
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
              Catatan Detail <span className="text-red-500">*</span>
            </label>
            <textarea
              {...register("notes")}
              className={`w-full rounded-lg border bg-surface-container-lowest py-2 px-3 text-sm focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none resize-none min-h-[72px] ${
                errors.notes ? "border-red-500" : "border-outline-variant"
              }`}
              placeholder="Tulis instruksi detail, nomor telepon, ringkasan agenda..."
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
              Tandai tugas ini langsung selesai (Logged Activity)
            </label>
          </div>

          <DialogFooter>
            <DialogClose asChild>
              <Button type="button" variant="ghost" className="cursor-pointer">
                Batal
              </Button>
            </DialogClose>
            <Button
              type="submit"
              variant="default"
              disabled={isLoading}
              className="flex items-center gap-1.5 cursor-pointer"
            >
              <span className="material-symbols-outlined text-[18px]">save</span>
              {isLoading ? "Menyimpan..." : "Simpan Tugas"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
