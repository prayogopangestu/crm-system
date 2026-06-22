"use client"

import React, { useState } from "react"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/Table"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/Avatar"
import { initialContacts, Contact } from "@/lib/mockData"

export default function KontakPage() {
  const [contacts, setContacts] = useState<Contact[]>(initialContacts)
  const [search, setSearch] = useState("")
  const [statusFilter, setStatusFilter] = useState("Semua")
  const [showAddModal, setShowAddModal] = useState(false)

  // Form states
  const [newName, setNewName] = useState("")
  const [newEmail, setNewEmail] = useState("")
  const [newCompany, setNewCompany] = useState("")
  const [newRole, setNewRole] = useState("")
  const [newStatus, setNewStatus] = useState<Contact['status']>("Prospek Awal")

  // Filter contacts
  const filteredContacts = contacts.filter((c) => {
    const matchesSearch =
      c.name.toLowerCase().includes(search.toLowerCase()) ||
      c.email.toLowerCase().includes(search.toLowerCase()) ||
      c.company.toLowerCase().includes(search.toLowerCase())
    
    const matchesStatus = statusFilter === "Semua" || c.status === statusFilter

    return matchesSearch && matchesStatus
  })

  // Deal status colors
  const statusBadges: Record<Contact['status'], string> = {
    "Negosiasi": "bg-tertiary-fixed text-on-tertiary-fixed-variant",
    "Menang": "bg-primary-fixed text-on-primary-fixed-variant",
    "Prospek Awal": "bg-surface-variant text-on-surface-variant",
    "Proposal": "bg-secondary-container text-on-secondary-container",
    "Kalah": "bg-error-container text-on-error-container",
    "Kualifikasi": "bg-surface-container-high text-on-surface"
  }

  const handleAddContact = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newName || !newEmail || !newCompany) return

    const initials = newName
      .split(" ")
      .map((n) => n[0])
      .join("")
      .substring(0, 2)
      .toUpperCase()

    const newContact: Contact = {
      id: Date.now().toString(),
      name: newName,
      email: newEmail,
      company: newCompany,
      role: newRole || "Staff",
      status: newStatus,
      lastContacted: "Baru saja",
      initials
    }

    setContacts([newContact, ...contacts])
    setShowAddModal(false)

    // Reset fields
    setNewName("")
    setNewEmail("")
    setNewCompany("")
    setNewRole("")
    setNewStatus("Prospek Awal")
  }

  return (
    <div className="p-gutter max-w-container-max-width mx-auto w-full relative">
      {/* Header Area */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-stack-lg gap-stack-md">
        <div>
          <h2 className="font-headline-md text-[24px] font-semibold text-on-surface">
            Kontak &amp; Perusahaan
          </h2>
          <p className="font-body-md text-sm text-on-surface-variant mt-1">
            Kelola data klien dan prospek B2B Anda.
          </p>
        </div>
        <Button
          onClick={() => setShowAddModal(true)}
          variant="default"
          className="flex items-center gap-2 rounded-lg text-sm font-semibold shadow-sm"
        >
          <span className="material-symbols-outlined text-[20px]">add</span>
          Tambah Kontak
        </Button>
      </div>

      {/* Filters & Search Row */}
      <div className="flex flex-col sm:flex-row gap-4 mb-6">
        <div className="relative flex-1 max-w-md">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px] pointer-events-none">
            search
          </span>
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-10"
            placeholder="Cari kontak, email, atau perusahaan..."
          />
        </div>
        <div className="w-full sm:w-48">
          <Select
            value={statusFilter}
            onValueChange={setStatusFilter}
          >
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Pilih Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Semua">Semua Status</SelectItem>
              <SelectItem value="Prospek Awal">Prospek Awal</SelectItem>
              <SelectItem value="Kualifikasi">Kualifikasi</SelectItem>
              <SelectItem value="Proposal">Proposal</SelectItem>
              <SelectItem value="Negosiasi">Negosiasi</SelectItem>
              <SelectItem value="Menang">Menang</SelectItem>
              <SelectItem value="Kalah">Kalah</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Data Table */}
      <Card className="overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[250px]">Nama &amp; Email</TableHead>
              <TableHead>Perusahaan</TableHead>
              <TableHead className="hidden md:table-cell">Jabatan</TableHead>
              <TableHead>Status Deal</TableHead>
              <TableHead className="hidden lg:table-cell">Terakhir Dihubungi</TableHead>
              <TableHead className="text-right">Aksi</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredContacts.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-on-surface-variant/70">
                  Tidak ada kontak ditemukan.
                </TableCell>
              </TableRow>
            ) : (
              filteredContacts.map((contact) => (
                <TableRow key={contact.id}>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={contact.avatarUrl} alt={contact.name} />
                        <AvatarFallback>{contact.initials}</AvatarFallback>
                      </Avatar>
                      <div className="min-w-0">
                        <div className="font-medium text-sm text-on-surface truncate">{contact.name}</div>
                        <div className="text-on-surface-variant text-xs truncate">{contact.email}</div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell className="text-sm truncate max-w-[150px]">
                    {contact.company}
                  </TableCell>
                  <TableCell className="hidden md:table-cell text-sm text-on-surface-variant truncate max-w-[150px]">
                    {contact.role}
                  </TableCell>
                  <TableCell>
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold ${statusBadges[contact.status]}`}>
                      {contact.status}
                    </span>
                  </TableCell>
                  <TableCell className="hidden lg:table-cell text-xs text-on-surface-variant">
                    {contact.lastContacted}
                  </TableCell>
                  <TableCell className="text-right">
                    <button className="text-on-surface-variant hover:text-primary transition-colors p-1 cursor-pointer">
                      <span className="material-symbols-outlined text-[20px]">more_vert</span>
                    </button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {/* Pagination Footer */}
        <div className="bg-surface border-t border-outline-variant px-4 py-3 flex items-center justify-between">
          <span className="text-xs text-on-surface-variant font-medium">
            Menampilkan {filteredContacts.length} dari {contacts.length} kontak
          </span>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" className="h-8 w-8 p-0" disabled>
              <span className="material-symbols-outlined text-[18px]">chevron_left</span>
            </Button>
            <Button variant="outline" size="sm" className="h-8 w-8 p-0" disabled>
              <span className="material-symbols-outlined text-[18px]">chevron_right</span>
            </Button>
          </div>
        </div>
      </Card>

      {/* Add Contact Modal Backdrop */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <Card className="w-full max-w-md bg-surface-container-lowest p-6 shadow-xl relative animate-in fade-in zoom-in duration-200">
            <button
              onClick={() => setShowAddModal(false)}
              className="absolute right-4 top-4 text-on-surface-variant hover:text-on-surface transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined">close</span>
            </button>
            <h3 className="font-headline-sm text-lg font-bold text-on-surface mb-4">
              Tambah Kontak Baru
            </h3>
            <form onSubmit={handleAddContact} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Nama Lengkap
                </label>
                <Input
                  required
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  placeholder="Nama Lengkap Klien"
                />
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Email
                </label>
                <Input
                  required
                  type="email"
                  value={newEmail}
                  onChange={(e) => setNewEmail(e.target.value)}
                  placeholder="name@company.com"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    Perusahaan
                  </label>
                  <Input
                    required
                    value={newCompany}
                    onChange={(e) => setNewCompany(e.target.value)}
                    placeholder="Nama Perusahaan"
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    Jabatan
                  </label>
                  <Input
                    value={newRole}
                    onChange={(e) => setNewRole(e.target.value)}
                    placeholder="VP Sales, Manager..."
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Status Kesepakatan (Deal)
                </label>
                <Select
                  value={newStatus}
                  onValueChange={(val) => setNewStatus(val as Contact['status'])}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Pilih Status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="Prospek Awal">Prospek Awal</SelectItem>
                    <SelectItem value="Kualifikasi">Kualifikasi</SelectItem>
                    <SelectItem value="Proposal">Proposal</SelectItem>
                    <SelectItem value="Negosiasi">Negosiasi</SelectItem>
                    <SelectItem value="Menang">Menang</SelectItem>
                    <SelectItem value="Kalah">Kalah</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setShowAddModal(false)}
                >
                  Batal
                </Button>
                <Button type="submit" variant="default">
                  Simpan Kontak
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
