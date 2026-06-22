"use client"

import React, { useState } from "react"
import { Card } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/Avatar"
import { Input } from "@/components/ui/Input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"
import { initialDeals, Deal } from "@/lib/mockData"

interface KanbanStage {
  id: Deal['stage']
  title: string
  color: string
}

const stages: KanbanStage[] = [
  { id: "lead", title: "Lead Masuk", color: "border-primary" },
  { id: "contacted", title: "Dihubungi", color: "border-secondary" },
  { id: "meeting", title: "Meeting", color: "border-tertiary-container" },
  { id: "negotiation", title: "Negosiasi", color: "border-primary-container" },
  { id: "won", title: "Deal Won", color: "border-surface-tint" }
]

export default function PipelinePage() {
  const [deals, setDeals] = useState<Deal[]>(initialDeals)
  const [showAddDealModal, setShowAddDealModal] = useState(false)

  // Add Deal Form States
  const [newTitle, setNewTitle] = useState("")
  const [newDesc, setNewDesc] = useState("")
  const [newValue, setNewValue] = useState("")
  const [newPriority, setNewPriority] = useState<Deal['priority']>("Medium")
  const [newStage, setNewStage] = useState<Deal['stage']>("lead")

  // Drag and Drop State
  const [draggingId, setDraggingId] = useState<string | null>(null)

  const handleDragStart = (e: React.DragEvent, id: string) => {
    setDraggingId(id)
    e.dataTransfer.setData("text/plain", id)
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
  }

  const handleDrop = (e: React.DragEvent, targetStage: Deal['stage']) => {
    e.preventDefault()
    const dealId = e.dataTransfer.getData("text/plain") || draggingId
    if (!dealId) return

    setDeals(prev =>
      prev.map(deal =>
        deal.id === dealId ? { ...deal, stage: targetStage } : deal
      )
    )
    setDraggingId(null)
  }

  const handleAddDeal = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newTitle || !newDesc || !newValue) return

    const newDeal: Deal = {
      id: "deal_" + Date.now().toString(),
      title: newTitle,
      company: newDesc,
      value: parseFloat(newValue.replace(/\D/g, "")) || 0,
      priority: newPriority,
      stage: newStage,
      assignee: {
        name: "Sarah",
        avatarUrl: "https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&q=80&w=100"
      }
    }

    setDeals([...deals, newDeal])
    setShowAddDealModal(false)

    // Reset forms
    setNewTitle("")
    setNewDesc("")
    setNewValue("")
    setNewPriority("Medium")
    setNewStage("lead")
  }

  // Priority color badges
  const priorityBadges: Record<Deal['priority'], string> = {
    High: "bg-error-container text-on-error-container",
    Medium: "bg-secondary-container text-on-secondary-container",
    Low: "bg-surface-variant text-on-surface-variant"
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

      {/* Kanban Board Container */}
      <div className="flex-1 overflow-x-auto overflow-y-hidden kanban-scroll pb-4 -mx-gutter px-gutter">
        <div className="flex gap-6 h-full min-w-max items-start">
          {stages.map(stage => {
            const stageDeals = deals.filter(d => d.stage === stage.id)
            const totalValue = stageDeals.reduce((sum, d) => sum + d.value, 0)

            return (
              <div
                key={stage.id}
                onDragOver={handleDragOver}
                onDrop={(e) => handleDrop(e, stage.id)}
                className={`w-[300px] flex flex-col max-h-full bg-surface-container-lowest border border-outline-variant/50 rounded-xl shrink-0 transition-all ${
                  draggingId ? "bg-surface-container-low/20" : ""
                }`}
              >
                {/* Stage Header */}
                <div className={`p-4 border-b border-outline-variant/30 flex justify-between items-center bg-surface/50 rounded-t-xl border-t-4 ${stage.color}`}>
                  <div className="flex items-center gap-2">
                    <h3 className="font-label-md text-xs font-bold text-on-surface uppercase tracking-wider">
                      {stage.title}
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
                            <AvatarFallback>{deal.assignee.name[0]}</AvatarFallback>
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
              onClick={() => setShowAddDealModal(false)}
              className="absolute right-4 top-4 text-on-surface-variant hover:text-on-surface transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined">close</span>
            </button>
            <h3 className="font-headline-sm text-lg font-bold text-on-surface mb-4">
              Tambah Deal Baru
            </h3>
            <form onSubmit={handleAddDeal} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Nama Perusahaan / Klien
                </label>
                <Input
                  required
                  value={newTitle}
                  onChange={(e) => setNewTitle(e.target.value)}
                  placeholder="Contoh: PT Telkomsel"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Deskripsi Proyek/Deal
                </label>
                <Input
                  required
                  value={newDesc}
                  onChange={(e) => setNewDesc(e.target.value)}
                  placeholder="Contoh: Pengadaan Cloud Server"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    Nilai Deal (Rupiah)
                  </label>
                  <Input
                    required
                    type="number"
                    value={newValue}
                    onChange={(e) => setNewValue(e.target.value)}
                    placeholder="Contoh: 150000000"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    Prioritas
                  </label>
                  <Select
                    value={newPriority}
                    onValueChange={(val) => setNewPriority(val as Deal['priority'])}
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
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Tahap Awal
                </label>
                <Select
                  value={newStage}
                  onValueChange={(val) => setNewStage(val as Deal['stage'])}
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
                  </SelectContent>
                </Select>
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setShowAddDealModal(false)}
                >
                  Batal
                </Button>
                <Button type="submit" variant="default">
                  Simpan Deal
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
