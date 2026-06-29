export type ContactStatus =
  | "Negosiasi"
  | "Menang"
  | "Prospek Awal"
  | "Proposal"
  | "Kalah"
  | "Kualifikasi"

export interface Contact {
  id: string
  name: string
  email: string
  company: string
  role: string
  status: ContactStatus
  lastContacted: string
  initials: string
  avatarUrl?: string
}

export type DealPriority = "High" | "Medium" | "Low"
export type SystemDealStage = "lead" | "contacted" | "meeting" | "negotiation" | "won" | "lost"
export type DealStage = SystemDealStage | (string & {})

export interface Deal {
  id: string
  title: string
  company: string
  value: number
  priority: DealPriority
  stage: DealStage
  lostReason?: string
  assignee: {
    id?: string
    name: string
    avatarUrl: string
  }
}

export type TaskType = "Meeting" | "Call" | "Proposal" | "Other"
export type TaskStatus = "overdue" | "today" | "upcoming"
export type TaskPriority = "Tinggi" | "Sedang" | "Rendah"

export interface Task {
  id: string
  title: string
  company: string
  time: string
  date: string
  type: TaskType
  status: TaskStatus
  completed: boolean
  notes: string
  priority: TaskPriority
  assignee: string
  assigneeId?: string
}

export interface Activity {
  id: string
  user: string
  action: string
  target: string
  time: string
  isHighlight?: boolean
}

export interface TeamMember {
  id: string
  name: string
  email: string
  role: "Admin" | "Staf Sales"
  status: "Aktif" | "Offline" | "Menunggu" | "Dicabut"
  initials: string
  avatarUrl?: string
}

export interface PipelineStage {
  id: string
  key: DealStage
  name: string
  color: string
  position: number
  isSystem?: boolean
}

export interface PerformanceGoal {
  month: string
  goal: number
  actual: number
  status: string
  percentage: number
}

export interface LeaderboardEntry {
  rank: number
  name: string
  role: string
  amount: number
  trend: string
  isPositive: boolean
  avatarUrl: string
}

export interface LostReason {
  name: string
  value: number
  percentage: number
  color: string
}

export interface DashboardStats {
  totalLeads: number
  leadsTrend: string
  dealWonCount: number
  wonTrend: string
  totalRevenue: string
  revenueTrend: string
  urgentTasksCount: number
}

export interface ConversionPoint {
  name: string
  Konversi: number
}

export interface UserProfile {
  id: string
  firstName: string
  lastName: string
  name: string
  email: string
  role: "Admin" | "Staf Sales"
  status?: string
  avatarUrl?: string
  initials?: string
}

export interface TelegramIntegration {
  enabled: boolean
  webhookUrl: string
  chatId?: string
  hasToken: boolean
}
