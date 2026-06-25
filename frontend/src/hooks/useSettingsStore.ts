import { create } from "zustand"
import { initialTeamMembers, pipelineStages, TeamMember } from "@/lib/mockData"

type Tab = "profil" | "tim" | "pipeline" | "integrasi"

interface SettingsState {
  activeTab: Tab
  teamMembers: TeamMember[]
  webhookEnabled: boolean
  stages: typeof pipelineStages
  showInviteModal: boolean
  setActiveTab: (tab: Tab) => void
  setShowInviteModal: (show: boolean) => void
  inviteMember: (member: TeamMember) => void
  addStage: (stage: { id: string; name: string; color: string }) => void
  setWebhookEnabled: (enabled: boolean) => void
}

export const useSettingsStore = create<SettingsState>((set) => ({
  activeTab: "tim",
  teamMembers: initialTeamMembers,
  webhookEnabled: true,
  stages: pipelineStages,
  showInviteModal: false,
  setActiveTab: (activeTab) => set({ activeTab }),
  setShowInviteModal: (showInviteModal) => set({ showInviteModal }),
  inviteMember: (member) => set((state) => ({
    teamMembers: [...state.teamMembers, member],
    showInviteModal: false
  })),
  addStage: (stage) => set((state) => ({
    stages: [...state.stages, stage]
  })),
  setWebhookEnabled: (webhookEnabled) => set({ webhookEnabled })
}))
