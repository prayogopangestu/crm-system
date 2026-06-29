import { create } from "zustand"
import { apiRequest, MutationResponse, queryString } from "@/lib/api"
import { Contact, ContactStatus } from "@/types/crm"

interface ContactListResponse {
  data: Contact[]
  total: number
  page: number
}

interface ContactInput {
  name: string
  email: string
  company: string
  role?: string
  status: ContactStatus
}

interface ContactState {
  contacts: Contact[]
  total: number
  page: number
  search: string
  statusFilter: string
  isLoading: boolean
  error: string | null
  showAddModal: boolean
  setSearch: (search: string) => void
  setStatusFilter: (filter: string) => void
  setShowAddModal: (show: boolean) => void
  loadContacts: () => Promise<void>
  addContact: (contact: ContactInput) => Promise<void>
}

export const useContactStore = create<ContactState>((set, get) => ({
  contacts: [],
  total: 0,
  page: 1,
  search: "",
  statusFilter: "Semua",
  isLoading: false,
  error: null,
  showAddModal: false,
  setSearch: (search) => set({ search }),
  setStatusFilter: (statusFilter) => set({ statusFilter }),
  setShowAddModal: (showAddModal) => set({ showAddModal }),
  loadContacts: async () => {
    const { search, statusFilter } = get()
    set({ isLoading: true, error: null })
    try {
      const response = await apiRequest<ContactListResponse>(
        `/api/contacts${queryString({
          search,
          status: statusFilter === "Semua" ? "" : statusFilter,
          page: 1,
          limit: 100,
        })}`
      )
      set({
        contacts: response.data,
        total: response.total,
        page: response.page,
        isLoading: false,
      })
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal memuat kontak",
      })
    }
  },
  addContact: async (contact) => {
    set({ isLoading: true, error: null })
    try {
      const response = await apiRequest<MutationResponse<Contact>>("/api/contacts", {
        method: "POST",
        body: { ...contact, role: contact.role || "Staff" },
      })
      set((state) => ({
        contacts: [response.data, ...state.contacts],
        total: state.total + 1,
        isLoading: false,
        showAddModal: false,
      }))
    } catch (error) {
      set({
        isLoading: false,
        error: error instanceof Error ? error.message : "Gagal menyimpan kontak",
      })
      throw error
    }
  },
}))
