"use client"

import React, { useEffect } from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { Card } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/Avatar"
import { Input } from "@/components/ui/Input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"
import { DealPriority, DealStage, PipelineStage } from "@/types/crm"
import { usePipelineStore } from "@/hooks/usePipelineStore"
import { toast } from "sonner"

const dealSchema = z.object({
  title: z.string().min(2, "Nama perusahaan/klien wajib diisi"),
  company: z.string().min(2, "Deskripsi proyek/deal wajib diisi"),
  value: z.string().min(1, "Nilai deal wajib diisi").refine((val) => !isNaN(parseFloat(val)) && parseFloat(val) > 0, "Nilai deal harus berupa angka positif"),
  priority: z.enum(["High", "Medium", "Low"]),
  stage: z.enum(["lead", "contacted", "meeting", "negotiation", "won", "lost"])
})

type DealInput = z.infer<typeof dealSchema>

export default function PipelinePage() {
  const {
    deals,
    stages,
    showAddDealModal,
    draggingId,
    isLoading,
    error,
    setShowAddDealModal,
    setDraggingId,
    loadDeals,
    loadStages,
    addDeal,
    updateDealStage
  } = usePipelineStore()

  const { register, control, handleSubmit, reset, formState: { errors } } = useForm<DealInput>({
    resolver: zodResolver(dealSchema),
    defaultValues: {
      title: "",
      company: "",
      value: "",
      priority: "Medium",
      stage: "lead"
    }
  })

  useEffect(() => {
    void loadStages()
    void loadDeals()
  }, [loadDeals, loadStages])

  const handleDragStart = (e: React.DragEvent, id: string) => {
    setDraggingId(id)
    e.dataTransfer.setData("text/plain", id)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const handleDrop = (e: React.DragEvent, targetStage: DealStage) => {
    e.preventDefault()
    const dealId = e.dataTransfer.getData("text/plain") || draggingId
    if (!dealId) return

    updateDealStage(dealId, targetStage).catch(() => {
      toast.error("Deal gagal dipindahkan")
    })
    setDraggingId(null)
  }

  const onSubmit = async (data: DealInput) => {
    try {
      await addDeal({
        title: data.title,
        company: data.company,
        value: parseFloat(data.value) || 0,
        priority: data.priority,
        stage: data.stage,
      })
      toast.success("Deal berhasil disimpan")
      reset()
    } catch {
      toast.error("Deal gagal disimpan")
    }
  }

  // Priority color badges
  const priorityBadges: Record<DealPriority, string> = {
    High: "bg-error-container text-on-error-container",
    Medium: "bg-secondary-container text-on-secondary-container",
    Low: "bg-surface-variant text-on-surface-variant"
  }

  const borderClass = (stage: PipelineStage) => {
    const colors: Record<string, string> = {
      lead: "border-primary",
      contacted: "border-secondary",
      meeting: "border-tertiary-container",
      negotiation: "border-primary-container",
      won: "border-surface-tint",
      lost: "border-error",
    }
    return colors[stage.key] || "border-outline-variant"
  }

  return (
    <div className="p-gutter max-w-[100%] mx-auto w-full h-[calc(100vh-64px)] flex flex-col overflow-hidden relative">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-stack-lg shrink-0">
        <div>
          <h2 className="font-headline-md text-[24px] font-semibold text-on-surface">
            Pipeline Penjualan
          </h2>
          <p className="font-body-md text-sm text-on-surface-variant mt-1">
            Kelola dan pantau progres deal tim Anda dengan antarmuka seret-dan-lepas (*drag &amp; drop*).
          </p>
        </div>
        <div className="flex items-center gap-3 w-full sm:w-auto">
          <Button
            onClick={() => setShowAddDealModal(true)}
            variant="default"
            className="flex items-center gap-2 rounded-lg text-sm font-semibold shadow-sm ml-auto"
          >
            <span className="material-symbols-outlined text-[20px]">add</span>
            Tambah Deal
          </Button>
        </div>
      </div>

      {error && (
        <div className="mb-4 rounded-lg border border-error-container bg-error-container/40 px-4 py-3 text-xs font-semibold text-on-error-container">
          {error}
        </div>
      )}

      {/* Kanban Board Container */}
      <div className="flex-1 overflow-x-auto overflow-y-hidden kanban-scroll pb-4 -mx-gutter px-gutter">
        <div className="flex gap-6 h-full min-w-max items-start">
          {isLoading && deals.length === 0 && stages.length === 0 && (
            <div className="text-sm text-on-surface-variant">Memuat pipeline dari backend...</div>
          )}
          {stages.map(stage => {
            const stageDeals = deals.filter(d => d.stage === stage.key)
            const totalValue = stageDeals.reduce((sum, d) => sum + d.value, 0)

            return (
              <div
                key={stage.id}
                onDragOver={handleDragOver}
                onDrop={(e) => handleDrop(e, stage.key)}
                className={`w-[300px] flex flex-col max-h-full bg-surface-container-lowest border border-outline-variant/50 rounded-xl shrink-0 transition-all ${
                  draggingId ? "bg-surface-container-low/20" : ""
                }`}
              >
                {/* Stage Header */}
                <div className={`p-4 border-b border-outline-variant/30 flex justify-between items-center bg-surface/50 rounded-t-xl border-t-4 ${borderClass(stage)}`}>
                  <div className="flex items-center gap-2">
                    <h3 className="font-label-md text-xs font-bold text-on-surface uppercase tracking-wider">
                      {stage.name}
                    </h3>
                    <span className="bg-surface-variant text-on-surface-variant px-2 py-0.5 rounded text-xs font-bold">
                      {stageDeals.length}
                    </span>
                  </div>
                  <span className="text-xs font-semibold text-primary">
                    Rp {(totalValue / 1000000).toFixed(1)}jt
                  </span>
                </div>

                {/* Card List */}
                <div className="p-3 flex-1 overflow-y-auto space-y-3 kanban-scroll bg-surface-container-low/30 min-h-[300px]">
                  {stageDeals.length === 0 ? (
                    <div className="border-2 border-dashed border-outline-variant/40 rounded-lg p-6 flex flex-col items-center justify-center text-center opacity-50 h-32">
                      <span className="material-symbols-outlined text-outline text-xl mb-1">
                        inbox
                      </span>
                      <p className="text-xs font-medium text-on-surface-variant">Tarik deal ke sini</p>
                    </div>
                  ) : (
                    stageDeals.map(deal => (
                      <div
                        key={deal.id}
                        draggable
                        onDragStart={(e) => handleDragStart(e, deal.id)}
                        className="bg-surface-container-lowest border border-outline-variant rounded-lg p-4 shadow-sm hover:shadow-md transition-all drag-handle group relative cursor-grab active:cursor-grabbing border-l-4 border-l-primary/60"
                      >
                        <div className="flex justify-between items-start mb-2">
                          <span className={`px-2 py-0.5 rounded text-[10px] font-bold uppercase ${priorityBadges[deal.priority]}`}>
                            {deal.priority}
                          </span>
                          <button className="text-on-surface-variant opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
                            <span className="material-symbols-outlined text-[16px]">more_vert</span>
                          </button>
                        </div>
                        <h4 className="font-semibold text-sm text-on-surface leading-tight mb-1">
                          {deal.title}
                        </h4>
                        <p className="text-xs text-on-surface-variant mb-4">
                          {deal.company}
                        </p>
                        <div className="flex justify-between items-center pt-2 border-t border-outline-variant/30">
                          <span className="text-xs font-bold text-primary">
                            Rp {deal.value.toLocaleString("id-ID")}
                          </span>
                        <Avatar className="h-6 w-6">
                          <AvatarImage src={deal.assignee.avatarUrl} alt={deal.assignee.name} />
                            <AvatarFallback>{deal.assignee.name?.[0] || "?"}</AvatarFallback>
                          </Avatar>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </div>
            )
          })}
        </div>
      </div>

      {/* Add Deal Modal */}
      {showAddDealModal && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <Card className="w-full max-w-md bg-surface-container-lowest p-6 shadow-xl relative animate-in fade-in zoom-in duration-200">
            <button
              onClick={() => {
                setShowAddDealModal(false)
                reset()
              }}
              className="absolute right-4 top-4 text-on-surface-variant hover:text-on-surface transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined">close</span>
            </button>
            <h3 className="font-headline-sm text-lg font-bold text-on-surface mb-4">
              Tambah Deal Baru
            </h3>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Nama Perusahaan / Klien
                </label>
                <Input
                  {...register("title")}
                  placeholder="Contoh: PT Telkomsel"
                  className={errors.title ? "border-red-500" : ""}
                />
                {errors.title && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.title.message}</p>
                )}
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Deskripsi Proyek/Deal
                </label>
                <Input
                  {...register("company")}
                  placeholder="Contoh: Pengadaan Cloud Server"
                  className={errors.company ? "border-red-500" : ""}
                />
                {errors.company && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.company.message}</p>
                )}
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    Nilai Deal (Rupiah)
                  </label>
                  <Input
                    type="number"
                    {...register("value")}
                    placeholder="Contoh: 150000000"
                    className={errors.value ? "border-red-500" : ""}
                  />
                  {errors.value && (
                    <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.value.message}</p>
                  )}
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    Prioritas
                  </label>
                  <Controller
                    name="priority"
                    control={control}
                    render={({ field }) => (
                      <Select
                        value={field.value}
                        onValueChange={field.onChange}
                      >
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Prioritas" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="Low">Low</SelectItem>
                          <SelectItem value="Medium">Medium</SelectItem>
                          <SelectItem value="High">High</SelectItem>
                        </SelectContent>
                      </Select>
                    )}
                  />
                  {errors.priority && (
                    <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.priority.message}</p>
                  )}
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Tahap Awal
                </label>
                <Controller
                  name="stage"
                  control={control}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      onValueChange={field.onChange}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Pilih Tahap" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="lead">Lead Masuk</SelectItem>
                        <SelectItem value="contacted">Dihubungi</SelectItem>
                        <SelectItem value="meeting">Meeting</SelectItem>
                        <SelectItem value="negotiation">Negosiasi</SelectItem>
                        <SelectItem value="won">Deal Won</SelectItem>
                        <SelectItem value="lost">Deal Lost</SelectItem>
                      </SelectContent>
                    </Select>
                  )}
                />
                {errors.stage && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.stage.message}</p>
                )}
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => {
                    setShowAddDealModal(false)
                    reset()
                  }}
                >
                  Batal
                </Button>
                <Button type="submit" variant="default" disabled={isLoading}>
                  {isLoading ? "Menyimpan..." : "Simpan Deal"}
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
