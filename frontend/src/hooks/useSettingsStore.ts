import { create } from "zustand"
import { apiRequest, MutationResponse } from "@/lib/api"
import { PipelineStage, TeamMember, TelegramIntegration, UserProfile } from "@/types/crm"

type Tab = "profil" | "tim" | "pipeline" | "integrasi"

interface SettingsState {
  activeTab: Tab
  profile: UserProfile | null
  teamMembers: TeamMember[]
  webhookEnabled: boolean
  webhookUrl: string
  stages: PipelineStage[]
  showInviteModal: boolean
  isLoading: boolean
  error: string | null
  setActiveTab: (tab: Tab) => void
  setShowInviteModal: (show: boolean) => void
  loadSettings: () => Promise<void>
  updateProfile: (profile: { firstName: string; lastName: string; email: string }) => Promise<void>
  inviteMember: (member: { name: string; email: string; role: "Admin" | "Staf Sales" }) => Promise<void>
  addStage: (stage: { name: string; color: string }) => Promise<void>
  setWebhookEnabled: (enabled: boolean) => Promise<void>
}

export const useSettingsStore = create<SettingsState>((set) => ({
  activeTab: "tim",
  profile: null,
  teamMembers: [],
  webhookEnabled: false,
  webhookUrl: "",
  stages: [],
  showInviteModal: false,
  isLoading: false,
  error: null,
  setActiveTab: (activeTab) => set({ activeTab }),
  setShowInviteModal: (showInviteModal) => set({ showInviteModal }),
  loadSettings: async () => {
    set({ isLoading: true, error: null })
    try {
      const [profile, teamMembers, stages, telegram] = await Promise.all([
        apiRequest<UserProfile>("/api/profile"),
        apiRequest<TeamMember[]>("/api/team"),
        apiRequest<PipelineStage[]>("/api/pipeline-stages"),
        apiRequest<TelegramIntegration>("/api/integrations/telegram"),
      ])
      set({
        profile,
        teamMembers,
        stages,
        webhookEnabled: telegram.enabled,
        webhookUrl: telegram.webhookUrl,
        isLoading: false,
      })
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal memuat pengaturan",
      })
    }
  },
  updateProfile: async (profileInput) => {
    set({ isLoading: true, error: null })
    try {
      const response = await apiRequest<MutationResponse<UserProfile>>("/api/profile", {
        method: "PUT",
        body: profileInput,
      })
      set({ profile: response.data, isLoading: false })
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal memperbarui profil",
      })
      throw error
    }
  },
  inviteMember: async (member) => {
    set({ isLoading: true, error: null })
    try {
      const response = await apiRequest<MutationResponse<TeamMember>>("/api/team/invite", {
        method: "POST",
        body: member,
      })
      set((state) => ({
        teamMembers: [...state.teamMembers, response.data],
        showInviteModal: false,
        isLoading: false,
      }))
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal mengundang anggota",
      })
      throw error
    }
  },
  addStage: async (stage) => {
    set({ isLoading: true, error: null })
    try {
      const response = await apiRequest<MutationResponse<PipelineStage>>("/api/pipeline-stages", {
        method: "POST",
        body: stage,
      })
      set((state) => ({
        stages: [...state.stages, response.data],
        isLoading: false,
      }))
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal menambah tahapan",
      })
      throw error
    }
  },
  setWebhookEnabled: async (webhookEnabled) => {
    set({ webhookEnabled, error: null })
    try {
      const response = await apiRequest<{ success: boolean; enabled: boolean }>("/api/integrations/telegram", {
        method: "PUT",
        body: { enabled: webhookEnabled },
      })
      set({
        webhookEnabled: response.enabled,
      })
    } catch (error) {
      set((state) => ({
        webhookEnabled: !state.webhookEnabled,
        error: error instanceof Error ? error.message : "Gagal mengubah integrasi",
      }))
      throw error
    }
  },
}))
