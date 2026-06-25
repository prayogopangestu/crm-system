import { create } from "zustand"
import { initialDeals, Deal } from "@/lib/mockData"

interface PipelineState {
  deals: Deal[]
  showAddDealModal: boolean
  draggingId: string | null
  setShowAddDealModal: (show: boolean) => void
  setDraggingId: (id: string | null) => void
  addDeal: (deal: Deal) => void
  updateDealStage: (dealId: string, stage: Deal['stage']) => void
}

export const usePipelineStore = create<PipelineState>((set) => ({
  deals: initialDeals,
  showAddDealModal: false,
  draggingId: null,
  setShowAddDealModal: (showAddDealModal) => set({ showAddDealModal }),
  setDraggingId: (draggingId) => set({ draggingId }),
  addDeal: (deal) => set((state) => ({
    deals: [...state.deals, deal],
    showAddDealModal: false
  })),
  updateDealStage: (dealId, stage) => set((state) => ({
    deals: state.deals.map((deal) =>
      deal.id === dealId ? { ...deal, stage } : deal
    )
  }))
}))
