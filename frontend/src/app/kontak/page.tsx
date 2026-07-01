"use client"

import React, { useEffect, useMemo } from "react"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { Button } from "@/components/ui/Button"
import { Card } from "@/components/ui/Card"
import { Input } from "@/components/ui/Input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/Select"
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/components/ui/Table"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/Avatar"
import { ContactStatus } from "@/types/crm"
import { useContactStore } from "@/hooks/useContactStore"
import { toast } from "sonner"
import { useLanguage } from "@/context/LanguageContext"

type ContactInput = {
  name: string
  email: string
  company: string
  role?: string
  status: "Negosiasi" | "Menang" | "Prospek Awal" | "Proposal" | "Kalah" | "Kualifikasi"
}

const CONTACT_STATUS_VALUES = ["Negosiasi", "Menang", "Prospek Awal", "Proposal", "Kalah", "Kualifikasi"] as const

export default function KontakPage() {
  const {
    contacts,
    total,
    search,
    statusFilter,
    isLoading,
    error,
    showAddModal,
    setSearch,
    setStatusFilter,
    setShowAddModal,
    loadContacts,
    addContact
  } = useContactStore()
  const { t } = useLanguage()

  const contactSchema = useMemo(
    () =>
      z.object({
        name: z.string().min(2, t("contacts.validation.nameMin")),
        email: z.string().email(t("contacts.validation.emailInvalid")).min(1, t("contacts.validation.emailRequired")),
        company: z.string().min(2, t("contacts.validation.companyMin")),
        role: z.string().optional(),
        status: z.enum(CONTACT_STATUS_VALUES),
      }),
    [t],
  )

  const { register, control, handleSubmit, reset, formState: { errors } } = useForm<ContactInput>({
    resolver: zodResolver(contactSchema),
    defaultValues: {
      name: "",
      email: "",
      company: "",
      role: "",
      status: "Prospek Awal"
    }
  })

  useEffect(() => {
    void loadContacts()
  }, [loadContacts, search, statusFilter])

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
  const statusBadges: Record<ContactStatus, string> = {
    "Negosiasi": "bg-tertiary-fixed text-on-tertiary-fixed-variant",
    "Menang": "bg-primary-fixed text-on-primary-fixed-variant",
    "Prospek Awal": "bg-surface-variant text-on-surface-variant",
    "Proposal": "bg-secondary-container text-on-secondary-container",
    "Kalah": "bg-error-container text-on-error-container",
    "Kualifikasi": "bg-surface-container-high text-on-surface"
  }

  const onSubmit = async (data: ContactInput) => {
    try {
      await addContact(data)
      toast.success(t("toast.contactSaved"))
      reset()
    } catch {
      toast.error(t("toast.contactSaveError"))
    }
  }

  return (
    <div className="p-gutter max-w-container-max-width mx-auto w-full relative">
      {/* Header Area */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-stack-lg gap-stack-md">
        <div>
          <h2 className="font-headline-md text-[24px] font-semibold text-on-surface">
            {t("contacts.title")}
          </h2>
          <p className="font-body-md text-sm text-on-surface-variant mt-1">
            {t("contacts.subtitle")}
          </p>
        </div>
        <Button
          onClick={() => setShowAddModal(true)}
          variant="default"
          className="flex items-center gap-2 rounded-lg text-sm font-semibold shadow-sm"
        >
          <span className="material-symbols-outlined text-[20px]">add</span>
          {t("contacts.addContact")}
        </Button>
      </div>

      {error && (
        <div className="mb-4 rounded-lg border border-error-container bg-error-container/40 px-4 py-3 text-xs font-semibold text-on-error-container">
          {error}
        </div>
      )}

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
            placeholder={t("contacts.searchPlaceholder")}
          />
        </div>
        <div className="w-full sm:w-48">
          <Select
            value={statusFilter}
            onValueChange={setStatusFilter}
          >
            <SelectTrigger className="w-full">
              <SelectValue placeholder={t("contacts.selectStatus")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Semua">{t("contacts.allStatus")}</SelectItem>
              <SelectItem value="Prospek Awal">{t("values.contactStatus.Prospek Awal")}</SelectItem>
              <SelectItem value="Kualifikasi">{t("values.contactStatus.Kualifikasi")}</SelectItem>
              <SelectItem value="Proposal">{t("values.contactStatus.Proposal")}</SelectItem>
              <SelectItem value="Negosiasi">{t("values.contactStatus.Negosiasi")}</SelectItem>
              <SelectItem value="Menang">{t("values.contactStatus.Menang")}</SelectItem>
              <SelectItem value="Kalah">{t("values.contactStatus.Kalah")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Data Table */}
      <Card className="overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[250px]">{t("contacts.colNameEmail")}</TableHead>
              <TableHead>{t("contacts.colCompany")}</TableHead>
              <TableHead className="hidden md:table-cell">{t("contacts.colRole")}</TableHead>
              <TableHead>{t("contacts.colStatus")}</TableHead>
              <TableHead className="hidden lg:table-cell">{t("contacts.colLastContact")}</TableHead>
              <TableHead className="text-right">{t("common.action")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading && contacts.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-on-surface-variant/70">
                  {t("contacts.loading")}
                </TableCell>
              </TableRow>
            ) : filteredContacts.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-on-surface-variant/70">
                  {t("contacts.notFound")}
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
                      {t(`values.contactStatus.${contact.status}`)}
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
            {t("contacts.showing", { shown: filteredContacts.length, total })}
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
              onClick={() => {
                setShowAddModal(false)
                reset()
              }}
              className="absolute right-4 top-4 text-on-surface-variant hover:text-on-surface transition-colors cursor-pointer"
            >
              <span className="material-symbols-outlined">close</span>
            </button>
            <h3 className="font-headline-sm text-lg font-bold text-on-surface mb-4">
              {t("contacts.addTitle")}
            </h3>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  {t("contacts.fullName")}
                </label>
                <Input
                  {...register("name")}
                  placeholder={t("contacts.fullNamePlaceholder")}
                  className={errors.name ? "border-red-500" : ""}
                />
                {errors.name && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.name.message}</p>
                )}
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  {t("contacts.email")}
                </label>
                <Input
                  type="email"
                  {...register("email")}
                  placeholder="name@company.com"
                  className={errors.email ? "border-red-500" : ""}
                />
                {errors.email && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.email.message}</p>
                )}
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    {t("contacts.company")}
                  </label>
                  <Input
                    {...register("company")}
                    placeholder={t("contacts.companyPlaceholder")}
                    className={errors.company ? "border-red-500" : ""}
                  />
                  {errors.company && (
                    <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.company.message}</p>
                  )}
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                    {t("contacts.role")}
                  </label>
                  <Input
                    {...register("role")}
                    placeholder={t("contacts.rolePlaceholder")}
                  />
                </div>
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">
                  {t("contacts.dealStatus")}
                </label>
                <Controller
                  name="status"
                  control={control}
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      onValueChange={field.onChange}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder={t("contacts.selectStatus")} />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="Prospek Awal">{t("values.contactStatus.Prospek Awal")}</SelectItem>
                        <SelectItem value="Kualifikasi">{t("values.contactStatus.Kualifikasi")}</SelectItem>
                        <SelectItem value="Proposal">{t("values.contactStatus.Proposal")}</SelectItem>
                        <SelectItem value="Negosiasi">{t("values.contactStatus.Negosiasi")}</SelectItem>
                        <SelectItem value="Menang">{t("values.contactStatus.Menang")}</SelectItem>
                        <SelectItem value="Kalah">{t("values.contactStatus.Kalah")}</SelectItem>
                      </SelectContent>
                    </Select>
                  )}
                />
                {errors.status && (
                  <p className="text-red-500 text-[10px] mt-1 font-medium">{errors.status.message}</p>
                )}
              </div>
              <div className="flex justify-end gap-2 pt-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => {
                    setShowAddModal(false)
                    reset()
                  }}
                >
                  {t("common.cancel")}
                </Button>
                <Button type="submit" variant="default" disabled={isLoading}>
                  {isLoading ? t("common.saving") : t("contacts.saveBtn")}
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
