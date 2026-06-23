# 🚀 API Endpoints Specification - CRM Enterprise System

Dokumen ini mendokumentasikan seluruh API endpoint yang dibutuhkan untuk menghubungkan antarmuka frontend **CRM Enterprise System** ke backend. Spesifikasi ini dirancang berdasarkan struktur modul, fitur, dan data model yang ada pada aplikasi saat ini.

---

## 📌 Daftar Isi
1. [Autentikasi & Profil Pengguna](#1-autentikasi--profil-pengguna)
2. [Dashboard Overview](#2-dashboard-overview)
3. [Manajemen Kontak (Contacts & Companies)](#3-manajemen-kontak-contacts--companies)
4. [Pipeline Penjualan (Kanban Board)](#4-pipeline-penjualan-kanban-board)
5. [Jadwal Aktivitas & Tugas (Tasks & Calendar)](#5-jadwal-aktivitas--tugas-tasks--calendar)
6. [Laporan & Analitik (Reports & Analytics)](#6-laporan--analitik-reports--analytics)
7. [Pengaturan Sistem & Manajemen Tim](#7-pengaturan-sistem--manajemen-tim)
8. [Integrasi Webhook Telegram](#8-integrasi-webhook-telegram)
9. [Pencarian Global (Global Search)](#9-pencarian-global-global-search)
10. [Notifikasi Sistem (System Notifications)](#10-notifikasi-sistem-system-notifications)

---

## 1. Autentikasi & Profil Pengguna
Mengelola sesi login, registrasi, dan informasi profil admin/sales rep yang aktif.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `POST` | `/api/auth/login` | Login ke sistem | `{ "email": "...", "password": "..." }` | `{ "token": "JWT_TOKEN", "user": { "id": "1", "name": "Sarah", "role": "Admin" } }` |
| `POST` | `/api/auth/register` | Registrasi akun baru | `{ "name": "...", "email": "...", "password": "..." }` | `{ "success": true, "message": "Registrasi berhasil" }` |
| `GET` | `/api/profile` | Ambil data profil user aktif | *None* (Header Bearer Token) | `{ "id": "1", "firstName": "Sarah", "lastName": "Jenkins", "email": "sarah.j@crm.com", "role": "Admin", "avatarUrl": "..." }` |
| `PUT` | `/api/profile` | Update profil aktif | `{ "firstName": "...", "lastName": "...", "email": "..." }` | `{ "success": true, "message": "Profil diperbarui" }` |

---

## 2. Dashboard Overview
Menampilkan statistik cepat (KPI) di halaman utama, grafik tren konversi, serta feed aktivitas terbaru.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/dashboard/stats` | Ambil ringkasan KPI dashboard | *None* | `{ "totalLeads": 1240, "leadsTrend": "+12% bulan ini", "dealWonCount": 42, "wonTrend": "+5 dari minggu lalu", "totalRevenue": "Rp 1.2M", "revenueTrend": "Stabil", "urgentTasksCount": 3 }` |
| `GET` | `/api/dashboard/conversion-chart`| Ambil data grafik konversi bulanan | *None* | `[ { "name": "Jan", "Konversi": 40 }, { "name": "Feb", "Konversi": 55 }, ... ]` |
| `GET` | `/api/dashboard/activities` | Ambil feed aktivitas terbaru (log sistem) | `?limit=5` | `[ { "id": "a1", "user": "Andi", "action": "menambahkan kontak baru", "target": "Budi Wijaya", "time": "10 menit yang lalu", "isHighlight": false } ]` |

---

## 3. Manajemen Kontak (Contacts & Companies)
Mengelola basis data klien B2B, status prospek, dan detail kontak perusahaan.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/contacts` | Ambil daftar kontak (dengan filter & search) | `?search=telkomsel&status=Negosiasi&page=1` | `{ "data": [ { "id": "1", "name": "Budi Wijaya", "email": "budi.w@telkomsel.co.id", "company": "PT Telkomsel", "role": "VP Sales", "status": "Negosiasi", "lastContacted": "Hari ini, 10:30", "initials": "BW" } ], "total": 8, "page": 1 }` |
| `POST` | `/api/contacts` | Tambahkan kontak baru | `{ "name": "...", "email": "...", "company": "...", "role": "...", "status": "..." }` | `{ "success": true, "data": { "id": "9", "name": "...", ... } }` |
| `PUT` | `/api/contacts/:id` | Update detail kontak | `{ "name": "...", "email": "...", "company": "...", "role": "...", "status": "..." }` | `{ "success": true, "data": { ... } }` |
| `DELETE` | `/api/contacts/:id` | Hapus data kontak | *None* | `{ "success": true, "message": "Kontak berhasil dihapus" }` |

---

## 4. Pipeline Penjualan (Kanban Board)
Melacak status deal proyek penjualan dengan visualisasi tahapan (stages) dari Lead hingga Deal Won.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/deals` | Ambil seluruh deal aktif | *None* | `[ { "id": "d1", "title": "PT Maju Mundur", "company": "Implementasi CRM", "value": 45000000, "priority": "High", "stage": "lead", "assignee": { "name": "Andi", "avatarUrl": "..." } } ]` |
| `POST` | `/api/deals` | Tambahkan deal proyek baru | `{ "title": "...", "company": "...", "value": 15000000, "priority": "Medium", "stage": "lead" }` | `{ "success": true, "data": { "id": "d5", "title": "...", ... } }` |
| `PATCH` | `/api/deals/:id/stage` | Update tahap deal (untuk drag & drop) | `{ "stage": "negotiation" }` | `{ "success": true, "message": "Tahap deal berhasil diperbarui" }` |
| `PUT` | `/api/deals/:id` | Edit detail deal secara lengkap | `{ "title": "...", "company": "...", "value": 120000000, "priority": "High", "stage": "meeting" }` | `{ "success": true, "data": { ... } }` |
| `DELETE` | `/api/deals/:id` | Hapus data deal proyek | *None* | `{ "success": true, "message": "Deal berhasil dihapus" }` |

---

## 5. Jadwal Aktivitas & Tugas (Tasks & Calendar)
Melacak agenda harian sales rep, deadline proposal, meeting, dan log aktivitas klien.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/tasks` | Ambil daftar tugas berdasarkan filter tanggal | `?date=2026-06-23` atau `?status=today` | `[ { "id": "t1", "title": "Telepon Budi Wijaya", "company": "PT Telkomsel", "time": "14:00", "type": "Call", "status": "today", "completed": false, "notes": "...", "priority": "Tinggi", "assignee": "Sarah" } ]` |
| `POST` | `/api/tasks` | Buat tugas baru | `{ "title": "...", "company": "...", "time": "...", "type": "Call", "date": "2026-06-23", "priority": "Sedang", "assignee": "Sarah", "notes": "...", "completed": false }` | `{ "success": true, "data": { "id": "t10", ... } }` |
| `PUT` | `/api/tasks/:id` | Edit data tugas / detail catatan | `{ "title": "...", "notes": "...", "priority": "..." }` | `{ "success": true, "data": { ... } }` |
| `PATCH` | `/api/tasks/:id/toggle` | Tandai tugas Selesai / Belum Selesai | `{ "completed": true }` | `{ "success": true, "completed": true }` |
| `DELETE` | `/api/tasks/:id` | Hapus tugas | *None* | `{ "success": true, "message": "Tugas dihapus" }` |

---

## 6. Laporan & Analitik (Reports & Analytics)
Menyediakan data performa penjualan individu, visualisasi pie chart alasan deal lost, dan ekspor laporan.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/reports/leaderboard` | Ambil peringkat performa sales rep | `?period=Bulan Ini` atau `?period=Bulan Lalu` | `[ { "rank": 1, "name": "Budi Santoso", "role": "Senior Sales Rep", "amount": 450000000, "trend": "+12%", "isPositive": true, "avatarUrl": "..." } ]` |
| `GET` | `/api/reports/lost-reasons` | Ambil data statistik penyebab Deal Lost | *None* | `[ { "name": "Harga Terlalu Tinggi", "value": 56, "percentage": 45, "color": "#6366f1" }, ... ]` |
| `GET` | `/api/reports/goals` | Ambil data target vs pencapaian bulanan | *None* | `[ { "month": "Oktober", "goal": 1000000000, "actual": 1150000000, "status": "Tercapai (115%)", "percentage": 115 } ]` |
| `GET` | `/api/reports/export/csv` | Download file CSV laporan penjualan | *None* | File Stream (Content-Type: `text/csv`) |
| `GET` | `/api/reports/export/pdf` | Download file PDF laporan penjualan | *None* | File Stream (Content-Type: `application/pdf`) |

---

## 7. Pengaturan Sistem & Manajemen Tim
Mengelola penambahan anggota tim baru, hak akses role, serta kustomisasi tahapan pipeline.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/team` | Ambil seluruh anggota tim | *None* | `[ { "id": "m1", "name": "Sarah Jenkins", "email": "sarah.j@crm.com", "role": "Admin", "status": "Aktif", "initials": "SJ" } ]` |
| `POST` | `/api/team/invite` | Undang anggota tim baru | `{ "name": "...", "email": "...", "role": "Staf Sales" }` | `{ "success": true, "data": { "id": "m4", "name": "...", "status": "Aktif", ... } }` |
| `DELETE` | `/api/team/:id` | Hapus akses anggota tim | *None* | `{ "success": true, "message": "Anggota tim dihapus" }` |
| `GET` | `/api/pipeline-stages` | Ambil daftar tahapan kustom | *None* | `[ { "id": "s1", "name": "Prospek Baru", "color": "bg-primary-container" } ]` |
| `POST` | `/api/pipeline-stages` | Tambah tahapan baru kustom | `{ "name": "...", "color": "..." }` | `{ "success": true, "data": { "id": "s5", "name": "..." } }` |
| `PUT` | `/api/pipeline-stages/reorder` | Simpan urutan baru tahapan penjualan | `{ "stagesOrder": ["s1", "s3", "s2", "s4"] }` | `{ "success": true }` |
| `DELETE` | `/api/pipeline-stages/:id` | Hapus tahapan kustom | *None* | `{ "success": true, "message": "Tahapan dihapus" }` |

---

## 8. Integrasi Webhook Telegram
Mengirimkan notifikasi pesan instan ke admin/channel Telegram ketika terjadi perubahan deal.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/integrations/telegram` | Cek status integrasi & webhook URL | *None* | `{ "enabled": true, "webhookUrl": "https://api.telegram.org/bot12345/..." }` |
| `PUT` | `/api/integrations/telegram` | Aktifkan/nonaktifkan integrasi | `{ "enabled": false }` | `{ "success": true, "enabled": false }` |
| `POST` | `/api/integrations/telegram/test` | Kirim notifikasi uji coba ke bot Telegram | *None* | `{ "success": true, "message": "Pesan uji coba terkirim" }` |

---

## 9. Pencarian Global (Global Search)
Mendukung bilah pencarian pada Header untuk menemukan kontak, tugas, atau deal di seluruh sistem secara langsung.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/search` | Mencari data kontak, tugas, dan deal secara global | `?q=query_text` | `{ "contacts": [ { "id": "1", "name": "Budi Wijaya", ... } ], "tasks": [ { "id": "t1", "title": "Telepon Budi", ... } ], "deals": [ { "id": "d1", "title": "PT Maju Mundur", ... } ] }` |

---

## 10. Notifikasi Sistem (System Notifications)
Mengelola notifikasi sistem real-time untuk pengguna yang ditampilkan pada menu notifikasi di Header.

| Method | Endpoint | Deskripsi | Request Body / Query Params | Format Response (JSON) |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/notifications` | Mengambil semua notifikasi aktif milik user | *None* | `[ { "id": "n1", "title": "Deal Won!", "message": "Deal #1042 berhasil dimenangkan oleh Andi", "time": "2 jam yang lalu", "read": false } ]` |
| `PATCH` | `/api/notifications/:id/read` | Menandai satu notifikasi sebagai telah dibaca | *None* | `{ "success": true }` |
| `PATCH` | `/api/notifications/read-all` | Menandai semua notifikasi sebagai telah dibaca | *None* | `{ "success": true }` |
