"use client"

import React, { useEffect } from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { Card } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/Table"
import { Avatar, AvatarFallback } from "@/components/ui/Avatar"
import { useSettingsStore } from "@/hooks/useSettingsStore"
import { toast } from "sonner"

type Tab = "profil" | "tim" | "pipeline" | "integrasi"

const profileSchema = z.object({
  firstName: z.string().min(1, "Nama depan wajib diisi"),
  lastName: z.string().min(1, "Nama belakang wajib diisi"),
  email: z.string().email("Format email tidak valid").min(1, "Email wajib diisi")
})

const inviteSchema = z.object({
  inviteName: z.string().min(2, "Nama lengkap wajib diisi"),
  inviteEmail: z.string().email("Format email tidak valid").min(1, "Email wajib diisi"),
  inviteRole: z.enum(["Admin", "Staf Sales"])
})

type ProfileInput = z.infer<typeof profileSchema>
type InviteInput = z.infer<typeof inviteSchema>

export default function PengaturanPage() {
  const {
    activeTab,
    profile,
    teamMembers,
    webhookEnabled,
    webhookUrl,
    stages,
    showInviteModal,
    isLoading,
    error,
    setActiveTab,
    setShowInviteModal,
    loadSettings,
    updateProfile,
    inviteMember,
    addStage,
    setWebhookEnabled
  } = useSettingsStore()

  const { register: registerProfile, handleSubmit: handleSubmitProfile, reset: resetProfile, formState: { errors: profileErrors } } = useForm<ProfileInput>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      email: ""
    }
  })

  const { register: registerInvite, control: controlInvite, handleSubmit: handleSubmitInvite, reset: resetInvite, formState: { errors: inviteErrors } } = useForm<InviteInput>({
    resolver: zodResolver(inviteSchema),
    defaultValues: {
      inviteName: "",
      inviteEmail: "",
      inviteRole: "Staf Sales"
    }
  })

  useEffect(() => {
    void loadSettings()
  }, [loadSettings])

  useEffect(() => {
    if (!profile) return
    resetProfile({
      firstName: profile.firstName,
      lastName: profile.lastName,
      email: profile.email,
    })
  }, [profile, resetProfile])

  const handleCopyWebhook = () => {
    navigator.clipboard.writeText(webhookUrl)
    toast.success("Webhook URL disalin ke clipboard!")
  }

  const handleInviteUser = async (data: InviteInput) => {
    try {
      await inviteMember({
        name: data.inviteName,
        email: data.inviteEmail,
        role: data.inviteRole,
      })
      toast.success("Undangan anggota berhasil dibuat")
      resetInvite()
    } catch {
      toast.error("Undangan anggota gagal dibuat")
    }
  }

  const handleUpdateProfile = async (data: ProfileInput) => {
    try {
      await updateProfile(data)
      toast.success("Profil diperbarui!")
    } catch {
      toast.error("Profil gagal diperbarui")
    }
  }

  const handleAddStage = async () => {
    const stageName = prompt("Masukkan nama tahapan penjualan baru:")
    if (!stageName) return

    try {
      await addStage({
        name: stageName,
        color: "bg-surface-variant",
      })
      toast.success("Tahapan berhasil ditambahkan")
    } catch {
      toast.error("Tahapan gagal ditambahkan")
    }
  }

  return (
    <div className="p-gutter max-w-container-max-width mx-auto w-full">
      {/* Page Header */}
      <div className="mb-stack-lg">
        <h2 className="font-headline-md text-[24px] font-semibold text-on-surface mb-2">
          Pengaturan
        </h2>
        <p className="font-body-lg text-sm text-on-surface-variant">
          Kelola konfigurasi sistem, tim, dan integrasi eksternal.
        </p>
      </div>

      {error && (
        <div className="mb-4 rounded-lg border border-error-container bg-error-container/40 px-4 py-3 text-xs font-semibold text-on-error-container">
          {error}
        </div>
      )}

      {/* Settings Tabs */}
      <div className="flex border-b border-outline-variant mb-stack-lg overflow-x-auto hide-scrollbar">
        {[
          { id: "profil", label: "Profil Saya" },
          { id: "tim", label: "Manajemen Tim" },
          { id: "pipeline", label: "Kustomisasi Pipeline" },
          { id: "integrasi", label: "Integrasi Webhook" }
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as Tab)}
            className={`px-6 py-3 font-label-md text-sm font-semibold whitespace-nowrap cursor-pointer transition-all border-b-2 ${
              activeTab === tab.id
                ? "text-primary border-primary"
                : "text-on-surface-variant border-transparent hover:text-on-surface hover:bg-surface-container-high/30"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Content Panels */}
      <div className="grid grid-cols-12 gap-gutter">
        {/* Tab 1: Profil Saya */}
        {activeTab === "profil" && (
          <Card className="col-span-12 p-6 max-w-2xl">
            <h3 className="font-headline-md text-base font-bold text-on-surface mb-4">Profil Saya</h3>
            <div className="flex items-center space-x-4 mb-6">
              <Avatar className="h-16 w-16">
                <AvatarFallback className="text-xl">{profile?.initials || "?"}</AvatarFallback>
              </Avatar>
              <div>
                <h4 className="font-bold text-base text-on-surface">{profile?.name || "Memuat profil..."}</h4>
                <p className="text-xs text-on-surface-variant">{profile?.role || "-"}</p>
              </div>
            </div>
            <form onSubmit={handleSubmitProfile(handleUpdateProfile)} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Nama Depan</label>
                  <Input
                    {...registerProfile("firstName")}
                    className={profileErrors.firstName ? "border-red-500" : ""}
                  />
                  {profileErrors.firstName && (
                    <p className="text-red-500 text-[10px] mt-1 font-medium">{profileErrors.firstName.message}</p>
                  )}
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Nama Belakang</label>
                  <Input
                    {...registerProfile("lastName")}
                    className={profileErrors.lastName ? "border-red-500" : ""}
                  />
                  {profileErrors.lastName && (
                    <p className="text-red-500 text-[10px] mt-1 font-medium">{profileErrors.lastName.message}</p>
                  )}
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Alamat Email</label>
                <Input
                  type="email"
                  {...registerProfile("email")}
                  className={profileErrors.email ? "border-red-500" : ""}
                />
                {profileErrors.email && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{profileErrors.email.message}</p>
                )}
              </div>
              <div className="flex justify-end pt-2">
                <Button type="submit" variant="default" disabled={isLoading}>
                  {isLoading ? "Menyimpan..." : "Simpan Perubahan"}
                </Button>
              </div>
            </form>
          </Card>
        )}

        {/* Tab 2: Manajemen Tim */}
        {activeTab === "tim" && (
          <Card className="col-span-12 lg:col-span-8 overflow-hidden flex flex-col">
            <div className="p-6 border-b border-outline-variant flex justify-between items-center bg-surface/50">
              <div>
                <h3 className="font-headline-md text-base font-bold text-on-surface">Anggota Tim</h3>
                <p className="text-xs text-on-surface-variant mt-1">Kelola akses dan peran pengguna dalam sistem.</p>
              </div>
              <Button
                onClick={() => setShowInviteModal(true)}
                variant="default"
                size="sm"
                className="flex items-center gap-2 font-semibold"
              >
                <span className="material-symbols-outlined text-[18px]">person_add</span>
                Undang Pengguna
              </Button>
            </div>
            <div className="flex-1 overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Pengguna</TableHead>
                    <TableHead>Peran</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Aksi</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {isLoading && teamMembers.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center py-8 text-on-surface-variant/70">
                        Memuat anggota tim...
                      </TableCell>
                    </TableRow>
                  )}
                  {teamMembers.map((member) => (
                    <TableRow key={member.id}>
                      <TableCell>
                        <div className="flex items-center space-x-3">
                          <Avatar className="h-9 w-9">
                            <AvatarFallback>{member.initials}</AvatarFallback>
                          </Avatar>
                          <div>
                            <p className="text-sm font-semibold text-on-surface">{member.name}</p>
                            <p className="text-xs text-on-surface-variant">{member.email}</p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold ${
                          member.role === "Admin"
                            ? "bg-secondary-container text-on-secondary-container"
                            : "bg-surface-container-highest text-on-surface border border-outline-variant"
                        }`}>
                          {member.role}
                        </span>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center text-xs font-semibold">
                          <div className={`h-2 w-2 rounded-full mr-2 ${
                            member.status === "Aktif" ? "bg-primary" : "bg-outline-variant"
                          }`}></div>
                          <span className="text-on-surface-variant">{member.status}</span>
                        </div>
                      </TableCell>
                      <TableCell className="text-right">
                        <button className="text-on-surface-variant hover:text-primary p-2 rounded-full hover:bg-surface-variant cursor-pointer">
                          <span className="material-symbols-outlined text-[20px]">more_vert</span>
                        </button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </Card>
        )}

        {/* Tab 3: Kustomisasi Pipeline */}
        {activeTab === "pipeline" && (
          <Card className="col-span-12 lg:col-span-6 p-6">
            <div className="flex justify-between items-start mb-4">
              <div>
                <h3 className="font-headline-sm text-base font-bold text-on-surface">Tahapan Penjualan</h3>
                <p className="text-xs text-on-surface-variant mt-0.5">Edit susunan tahapan penjualan.</p>
              </div>
              <Button onClick={handleAddStage} variant="outline" size="sm">
                Tambah Tahapan
              </Button>
            </div>
            <div className="space-y-2 mt-4">
              {stages.map((stage) => (
                <div
                  key={stage.id}
                  className="flex items-center justify-between p-3 rounded border border-outline-variant bg-surface hover:border-primary transition-colors cursor-move"
                >
                  <div className="flex items-center space-x-2">
                    <span className="material-symbols-outlined text-outline-variant text-[18px] select-none">
                      drag_indicator
                    </span>
                    <span className="font-label-md text-sm font-semibold text-on-surface">{stage.name}</span>
                  </div>
                  <span className={`w-3.5 h-3.5 rounded-full ${stage.color || "bg-primary-container"}`}></span>
                </div>
              ))}
            </div>
          </Card>
        )}

        {/* Tab 4: Webhook Integration */}
        {activeTab === "integrasi" && (
          <Card className="col-span-12 lg:col-span-6 p-6 flex flex-col">
            <div className="flex items-center space-x-3 mb-4">
              <div className="p-2 bg-surface-container rounded-lg text-on-surface-variant">
                <span className="material-symbols-outlined">api</span>
              </div>
              <div>
                <h3 className="font-headline-sm text-base font-bold text-on-surface">Integrasi Webhook</h3>
                <p className="text-xs text-on-surface-variant mt-0.5">Notifikasi Real-time</p>
              </div>
            </div>
            <div className="bg-surface border border-outline-variant rounded-lg p-4 mb-4 flex-1">
              <div className="flex justify-between items-center mb-3">
                <div className="flex items-center space-x-2">
                  <span className="material-symbols-outlined text-primary text-[20px]">send</span>
                  <span className="font-label-md text-sm font-bold text-on-surface">Telegram Bot API</span>
                </div>
                <label className="relative inline-flex items-center cursor-pointer select-none">
                  <input
                    type="checkbox"
                    checked={webhookEnabled}
                    onChange={() => {
                      setWebhookEnabled(!webhookEnabled).catch(() => {
                        toast.error("Integrasi gagal diperbarui")
                      })
                    }}
                    className="sr-only peer"
                  />
                  <div className="w-9 h-5 bg-surface-container-highest peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-outline-variant after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-primary"></div>
                </label>
              </div>
              <p className="text-xs text-on-surface-variant mb-4 leading-relaxed">
                Kirimkan notifikasi telegram secara instan setiap kali kesepakatan penjualan berubah ke tahap Won (Menang).
              </p>
              <div className="space-y-3">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Webhook URL</label>
                  <div className="flex items-center border border-outline-variant rounded-lg overflow-hidden focus-within:border-primary transition-colors bg-surface-container-lowest">
                    <input
                      className="w-full bg-transparent border-none text-xs text-on-surface py-2 px-3 focus:ring-0 outline-none opacity-70"
                      readOnly
                      type="text"
                      value={webhookUrl || "-"}
                    />
                    <button
                      type="button"
                      onClick={handleCopyWebhook}
                      className="p-2 bg-surface-container-high hover:bg-surface-variant text-on-surface-variant transition-colors border-l border-outline-variant cursor-pointer"
                    >
                      <span className="material-symbols-outlined text-[16px]">content_copy</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
            <Button
              onClick={() => toast.info("Membuka pengaturan API Telegram...")}
              className="w-full font-semibold"
              variant="secondary"
            >
              Konfigurasi Telegram
            </Button>
          </Card>
        )}
      </div>

      {/* Invite Member Modal */}
      {showInviteModal && (
        <div className="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <Card className="w-full max-w-md bg-surface-container-lowest p-6 shadow-xl relative animate-in fade-in zoom-in duration-200">
            <button
              onClick={() => {
                setShowInviteModal(false)
                resetInvite()
              }}
              className="absolute right-4 top-4 text-on-surface-variant hover:text-on-surface transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined">close</span>
            </button>
            <h3 className="font-headline-sm text-lg font-bold text-on-surface mb-4">
              Undang Anggota Tim
            </h3>
            <form onSubmit={handleSubmitInvite(handleInviteUser)} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Nama Lengkap
                </label>
                <Input
                  {...registerInvite("inviteName")}
                  placeholder="Nama Lengkap Anggota"
                  className={inviteErrors.inviteName ? "border-red-500" : ""}
                />
                {inviteErrors.inviteName && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{inviteErrors.inviteName.message}</p>
                )}
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Alamat Email
                </label>
                <Input
                  type="email"
                  {...registerInvite("inviteEmail")}
                  placeholder="name@company.com"
                  className={inviteErrors.inviteEmail ? "border-red-500" : ""}
                />
                {inviteErrors.inviteEmail && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{inviteErrors.inviteEmail.message}</p>
                )}
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  Peran Akses (Role)
                </label>
                <Controller
                  name="inviteRole"
                  control={controlInvite}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      onValueChange={field.onChange}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Pilih Peran" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="Staf Sales">Staf Sales</SelectItem>
                        <SelectItem value="Admin">Admin</SelectItem>
                      </SelectContent>
                    </Select>
                  )}
                />
                {inviteErrors.inviteRole && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{inviteErrors.inviteRole.message}</p>
                )}
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => {
                    setShowInviteModal(false)
                    resetInvite()
                  }}
                >
                  Batal
                </Button>
                <Button type="submit" variant="default" disabled={isLoading}>
                  {isLoading ? "Mengundang..." : "Undang"}
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
