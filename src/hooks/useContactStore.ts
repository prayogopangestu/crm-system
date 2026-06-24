import { create } from "zustand"
import { initialContacts, Contact } from "@/lib/mockData"

interface ContactState {
  contacts: Contact[]
  search: string
  statusFilter: string
  showAddModal: boolean
  setSearch: (search: string) => void
  setStatusFilter: (filter: string) => void
  setShowAddModal: (show: boolean) => void
  addContact: (contact: Contact) => void
}

export const useContactStore = create<ContactState>((set) => ({
  contacts: initialContacts,
  search: "",
  statusFilter: "Semua",
  showAddModal: false,
  setSearch: (search) => set({ search }),
  setStatusFilter: (statusFilter) => set({ statusFilter }),
  setShowAddModal: (showAddModal) => set({ showAddModal }),
  addContact: (contact) => set((state) => ({
    contacts: [contact, ...state.contacts],
    showAddModal: false
  }))
}))
