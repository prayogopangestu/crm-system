import { create } from "zustand"
import { apiRequest, MutationResponse } from "@/lib/api"
import { Deal, DealPriority, DealStage, PipelineStage } from "@/types/crm"

interface DealInput {
  title: string
  company: string
  value: number
  priority: DealPriority
  stage: DealStage
}

interface PipelineState {
  deals: Deal[]
  stages: PipelineStage[]
  showAddDealModal: boolean
  draggingId: string | null
  isLoading: boolean
  error: string | null
  setShowAddDealModal: (show: boolean) => void
  setDraggingId: (id: string | null) => void
  loadDeals: () => Promise<void>
  loadStages: () => Promise<void>
  addDeal: (deal: DealInput) => Promise<void>
  updateDealStage: (dealId: string, stage: DealStage) => Promise<void>
}

export const usePipelineStore = create<PipelineState>((set, get) => ({
  deals: [],
  stages: [],
  showAddDealModal: false,
  draggingId: null,
  isLoading: false,
  error: null,
  setShowAddDealModal: (showAddDealModal) => set({ showAddDealModal }),
  setDraggingId: (draggingId) => set({ draggingId }),
  loadDeals: async () => {
    set({ isLoading: true, error: null })
    try {
      const deals = await apiRequest<Deal[]>("/api/deals")
      set({ deals, isLoading: false })
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal memuat deal",
      })
    }
  },
  loadStages: async () => {
    try {
      const stages = await apiRequest<PipelineStage[]>("/api/pipeline-stages")
      set({ stages })
    } catch (error) {
      set({ error: error instanceof Error ? error.message : "Gagal memuat tahap pipeline" })
    }
  },
  addDeal: async (deal) => {
    set({ isLoading: true, error: null })
    try {
      const response = await apiRequest<MutationResponse<Deal>>("/api/deals", {
        method: "POST",
        body: deal,
      })
      set((state) => ({
        deals: [...state.deals, response.data],
        isLoading: false,
        showAddDealModal: false,
      }))
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal menyimpan deal",
      })
      throw error
    }
  },
  updateDealStage: async (dealId, stage) => {
    const previous = get().deals
    set({
      deals: previous.map((deal) => (deal.id === dealId ? { ...deal, stage } : deal)),
      error: null,
    })
    try {
      await apiRequest<{ success: boolean; message: string }>(`/api/deals/${dealId}/stage`, {
        method: "PATCH",
        body: { stage },
      })
    } catch (error) {
      set({
        deals: previous,
        error: error instanceof Error ? error.message : "Gagal memindahkan deal",
      })
      throw error
    }
  },
}))
