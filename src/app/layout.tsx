import type { Metadata } from "next"
import { Inter } from "next/font/google"
import "./globals.css"
import { Sidebar } from "@/components/common/Sidebar"
import { Header } from "@/components/common/Header"

import { LayoutWrapper } from "@/components/common/LayoutWrapper"

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-sans",
})

export const metadata: Metadata = {
  title: "CRM Enterprise - Sistem Manajemen Sales",
  description: "Dashboard CRM untuk performa penjualan, pipeline, kontak, dan laporan analitik.",
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="id" className={`${inter.variable} h-full antialiased`}>
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="font-sans bg-background text-on-background min-h-screen">
        <LayoutWrapper>{children}</LayoutWrapper>
      </body>
    </html>
  )
}
