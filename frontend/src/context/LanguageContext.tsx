"use client"

import React, { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react"
import {
  DEFAULT_LOCALE,
  LOCALE_STORAGE_KEY,
  SUPPORTED_LOCALES,
  translations,
  type Locale,
} from "@/i18n/translations"

interface LanguageContextValue {
  locale: Locale
  setLocale: (locale: Locale) => void
  toggleLocale: () => void
  t: (key: string, vars?: Record<string, string | number>) => string
  tList: (key: string) => string[]
}

const LanguageContext = createContext<LanguageContextValue | undefined>(undefined)

function getNestedValue(obj: unknown, path: string): unknown {
  return path.split(".").reduce<unknown>((acc, key) => {
    if (acc && typeof acc === "object") {
      return (acc as Record<string, unknown>)[key]
    }
    return undefined
  }, obj)
}

function interpolate(template: string, vars?: Record<string, string | number>): string {
  if (!vars) return template
  return template.replace(/\{(\w+)\}/g, (_, name) =>
    vars[name] !== undefined ? String(vars[name]) : `{${name}}`,
  )
}

function isBrowserLocale(value: unknown): value is Locale {
  return value === "id" || value === "en"
}

export function LanguageProvider({ children }: { children: React.ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(DEFAULT_LOCALE)

  useEffect(() => {
    const applyStored = () => {
      const stored = window.localStorage.getItem(LOCALE_STORAGE_KEY)
      if (isBrowserLocale(stored)) {
        setLocaleState(stored)
      }
    }
    const timer = window.setTimeout(applyStored, 0)
    return () => window.clearTimeout(timer)
  }, [])

  useEffect(() => {
    document.documentElement.lang = locale
  }, [locale])

  const setLocale = useCallback((next: Locale) => {
    setLocaleState(next)
    window.localStorage.setItem(LOCALE_STORAGE_KEY, next)
  }, [])

  const toggleLocale = useCallback(() => {
    setLocaleState((prev) => {
      const next: Locale = prev === "id" ? "en" : "id"
      window.localStorage.setItem(LOCALE_STORAGE_KEY, next)
      return next
    })
  }, [])

  const t = useCallback(
    (key: string, vars?: Record<string, string | number>) => {
      const value = getNestedValue(translations[locale], key)
      if (typeof value === "string") {
        return interpolate(value, vars)
      }
      const fallback = getNestedValue(translations[DEFAULT_LOCALE], key)
      if (typeof fallback === "string") {
        return interpolate(fallback, vars)
      }
      return key
    },
    [locale],
  )

  const tList = useCallback(
    (key: string) => {
      const value = getNestedValue(translations[locale], key)
      if (Array.isArray(value)) {
        return value.map(String)
      }
      const fallback = getNestedValue(translations[DEFAULT_LOCALE], key)
      if (Array.isArray(fallback)) {
        return fallback.map(String)
      }
      return []
    },
    [locale],
  )

  const value = useMemo<LanguageContextValue>(
    () => ({ locale, setLocale, toggleLocale, t, tList }),
    [locale, setLocale, toggleLocale, t, tList],
  )

  return <LanguageContext.Provider value={value}>{children}</LanguageContext.Provider>
}

export function useLanguage() {
  const context = useContext(LanguageContext)
  if (!context) {
    throw new Error("useLanguage must be used within a LanguageProvider")
  }
  return context
}

export { SUPPORTED_LOCALES }
