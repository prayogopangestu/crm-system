# Analisis Implementasi Backend Go CRM

Dokumen ini adalah sumber keputusan implementasi backend untuk kontrak
`frontend/endpoint.md`. Backend menggunakan clean architecture:

```text
delivery (HTTP/gRPC) -> usecase -> domain interface -> repository/driver
```

## 1. Runtime dan library

| Kebutuhan | Pilihan | Alasan |
|---|---|---|
| Runtime | Go 1.25 | Versi Go yang tersedia di workspace |
| HTTP router | `go-chi/chi/v5` | Ringan, kompatibel `net/http`, route grouping dan middleware |
| PostgreSQL | `jackc/pgx/v5/pgxpool` | Driver PostgreSQL native dan connection pool |
| Migration | Runner SQL internal (marker Goose-compatible) | Tidak memerlukan CLI global dan tetap mendukung up/down/status |
| Redis | `redis/go-redis/v9` | Client resmi Redis untuk cache dan rate limit |
| JWT | `golang-jwt/jwt/v5` | Validasi ketat JWT HS256 |
| Validation | Validasi usecase eksplisit | Menjaga aturan bisnis tetap independen dari transport |
| Password | `golang.org/x/crypto/bcrypt` | Password hashing cost 12 |
| ID | `google/uuid` | UUID untuk seluruh primary key |
| gRPC | `google.golang.org/grpc` + protobuf | Kontrak UserService |
| Config | `gopkg.in/yaml.v3` | Config YAML dengan override environment |
| PDF | `go-pdf/fpdf` | Export laporan PDF tanpa service eksternal |
| Test | `testing`, `httptest`, `bufconn` | Unit, handler, Telegram, dan gRPC tanpa dependency test tambahan |
| Logging | `log/slog` | Structured logging dari standard library |

Tool lokal yang belum tersedia saat analisis: `protoc`, `buf`,
`golangci-lint`, `mockgen`, dan `sqlc`. Protobuf digenerasikan melalui
Dockerized Buf. Migration dapat dijalankan melalui binary `cmd/migrate`,
sehingga instalasi CLI migration global tidak wajib.

## 2. Multi-tenancy dan keamanan

- Registrasi membuat `organizations` dan user pertama ber-role `Admin` dalam
  satu transaksi.
- JWT berumur 24 jam dan memuat `sub` (user ID), `organization_id`, `role`,
  `iat`, dan `exp`.
- Semua repository tenant menerima `organizationID`; query tidak boleh
  mengambil resource hanya berdasarkan ID.
- Admin mengelola team, pipeline stage, dan Telegram. Admin dan Staf Sales
  dapat mengelola contact, deal, dan task.
- Password di-hash bcrypt cost 12. Token undangan acak 32 byte hanya disimpan
  dalam bentuk SHA-256 dan berlaku 72 jam.
- Token Telegram dienkripsi AES-256-GCM menggunakan
  `APP_ENCRYPTION_KEY` (base64 32 byte), dan tidak pernah dikembalikan.
- Login/register dibatasi 5 request/menit/IP melalui Redis. Jika Redis tidak
  tersedia, request tetap berjalan dan warning dicatat.
- Error JSON:

```json
{
  "error": {
    "code": "validation_error",
    "message": "request tidak valid",
    "fields": {"email": "harus berupa email yang valid"},
    "requestId": "..."
  }
}
```

## 3. Model data dan indeks

| Tabel | Data utama | Indeks/aturan penting |
|---|---|---|
| `organizations` | name, timestamps | name |
| `users` | organization, name, email, password hash, role, status | unique lower(email), organization/status |
| `user_invitations` | organization, email, role, token hash, expiry | unique token hash, organization/email |
| `contacts` | organization, owner, name, email, company, role, status | tenant + search/status, soft delete |
| `pipeline_stages` | organization, key, name, color, position, system flag | unique tenant/key, ordered position |
| `deals` | organization, assignee, title, company, value, priority, stage key, lost reason | tenant/stage, soft delete |
| `tasks` | organization, assignee, title, company, due date/time, type, priority, notes, completed | tenant/date/status, soft delete |
| `activities` | organization, actor, action, target, highlight | tenant/created_at |
| `performance_goals` | organization, month, goal | unique tenant/month |
| `notifications` | organization, user nullable, title, message, read_at | tenant/user/read |
| `telegram_integrations` | organization, encrypted token, chat ID, enabled | one row per tenant |
| `outbox_events` | organization, event type, JSON payload, attempts, next attempt | pending/next attempt |

Nilai uang disimpan sebagai `BIGINT` Rupiah. Waktu disimpan `TIMESTAMPTZ`.
Contact, deal, dan task menggunakan `deleted_at`.

## 4. Kontrak endpoint

Response sukses mempertahankan bentuk contoh `frontend/endpoint.md`.
Endpoint selain auth, health, readiness, dan accept-invite memerlukan Bearer
token.

| Method dan path | Role | Status sukses | Catatan |
|---|---|---:|---|
| `POST /api/auth/login` | Publik | 200 | JWT dan ringkasan user |
| `POST /api/auth/register` | Publik | 201 | Mendukung `name` atau `fullName`, wajib `companyName` |
| `POST /api/auth/accept-invite` | Publik | 200 | Token + password, mengaktifkan user |
| `GET /api/profile` | Login | 200 | Cache Redis 5 menit |
| `PUT /api/profile` | Login | 200 | Invalidasi cache profile |
| `GET /api/dashboard/stats` | Login | 200 | Cache 60 detik |
| `GET /api/dashboard/conversion-chart` | Login | 200 | Persentase won per bulan |
| `GET /api/dashboard/activities` | Login | 200 | `limit` default 5, maksimum 100 |
| `GET /api/contacts` | Login | 200 | search, status, page, limit |
| `POST /api/contacts` | Login | 201 | Membuat activity |
| `PUT /api/contacts/:id` | Login | 200 | Tenant-scoped |
| `DELETE /api/contacts/:id` | Login | 200 | Soft delete |
| `GET /api/deals` | Login | 200 | Deal aktif termasuk won/lost |
| `POST /api/deals` | Login | 201 | Stage harus tersedia pada tenant |
| `PATCH /api/deals/:id/stage` | Login | 200 | `lostReason` wajib untuk `lost` |
| `PUT /api/deals/:id` | Login | 200 | Update lengkap |
| `DELETE /api/deals/:id` | Login | 200 | Soft delete |
| `GET /api/tasks` | Login | 200 | filter `date` atau status relatif |
| `POST /api/tasks` | Login | 201 | Status relatif dihitung dari `date` |
| `PUT /api/tasks/:id` | Login | 200 | Partial-compatible update |
| `PATCH /api/tasks/:id/toggle` | Login | 200 | Mengubah `completed` |
| `DELETE /api/tasks/:id` | Login | 200 | Soft delete |
| `GET /api/reports/leaderboard` | Login | 200 | Period bulan ini/lalu |
| `GET /api/reports/lost-reasons` | Login | 200 | Agregasi deal stage `lost` |
| `GET /api/reports/goals` | Login | 200 | Goal vs deal won |
| `GET /api/reports/export/csv` | Login | 200 | `text/csv` attachment |
| `GET /api/reports/export/pdf` | Login | 200 | `application/pdf` attachment |
| `GET /api/team` | Admin | 200 | Termasuk status `Menunggu` |
| `POST /api/team/invite` | Admin | 201 | Mengembalikan `inviteUrl` sekali |
| `DELETE /api/team/:id` | Admin | 200 | Revoke access, histori dipertahankan |
| `GET /api/pipeline-stages` | Login | 200 | Urut berdasarkan position |
| `POST /api/pipeline-stages` | Admin | 201 | Key dibuat dari slug nama |
| `PUT /api/pipeline-stages/reorder` | Admin | 200 | Transaksi seluruh urutan |
| `DELETE /api/pipeline-stages/:id` | Admin | 200 | Ditolak jika sedang dipakai |
| `GET /api/integrations/telegram` | Admin | 200 | Tidak mengembalikan bot token |
| `PUT /api/integrations/telegram` | Admin | 200 | enabled, botToken, chatId |
| `POST /api/integrations/telegram/test` | Admin | 200 | Memanggil Bot API `sendMessage` |
| `GET /api/search` | Login | 200 | contact/task/deal, cache 30 detik |
| `GET /api/notifications` | Login | 200 | Notification tenant/global user |
| `PATCH /api/notifications/:id/read` | Login | 200 | Tenant dan user scoped |
| `PATCH /api/notifications/read-all` | Login | 200 | Semua notification user aktif |
| `GET /healthz` | Publik | 200 | Process liveness |
| `GET /readyz` | Publik | 200/503 | PostgreSQL wajib, Redis informasional |

## 5. Transaksi dan efek samping

- Register: organization + admin + default stages + initial goals.
- Invite: pending user + invitation.
- Deal won: update deal + activity + notifications + Telegram outbox.
- Deal lost: update deal + lost reason + activity.
- Pipeline reorder: seluruh posisi dalam satu transaksi.
- Telegram worker mengunci event dengan `FOR UPDATE SKIP LOCKED`, melakukan
  retry exponential, dan menandai processed tanpa mengubah transaksi CRM.

## 6. Cache

| Key | TTL | Invalidasi |
|---|---:|---|
| `crm:{org}:profile:{user}` | 5 menit | update profile/team revoke |
| `crm:{org}:dashboard:*` | 60 detik | mutasi contact/deal/task |
| `crm:{org}:reports:*` | 5 menit | deal/goal berubah |
| `crm:{org}:search:{hash}` | 30 detik | mutasi contact/deal/task |
| `crm:rate:{route}:{ip}:{minute}` | 2 menit | expiry otomatis |

## 7. Environment

`APP_ENV`, `HTTP_ADDR`, `GRPC_ADDR`, `DATABASE_URL`, `REDIS_URL`,
`JWT_SECRET`, `JWT_TTL`, `APP_ENCRYPTION_KEY`, `CORS_ALLOWED_ORIGINS`,
`APP_BASE_URL`, `APP_TIMEZONE`, `TELEGRAM_WORKER_INTERVAL`,
`TELEGRAM_WORKER_BATCH_SIZE`, `LOG_LEVEL`.

## 8. Ketidaksesuaian frontend

- Frontend masih memakai mock/Zustand dan belum memanggil API.
- Form registrasi menggunakan `fullName`; backend menerima alias tersebut.
- Task mock awal tidak mempunyai date, priority, assignee, dan notes; backend
  selalu mengembalikan semua field.
- Pipeline frontend memiliki `contacted`, tetapi setting awal tidak. Backend
  menyatukan enam stage default termasuk `lost`.
- Frontend belum memiliki UI `lostReason`, status team `Menunggu`, aktivasi
  invite, dan kredensial Telegram.
- Nilai `lastContacted`, activity time, dan notification time dikirim sebagai
  teks ramah UI agar kompatibel, sementara timestamp ISO tetap tersedia pada
  field tambahan `createdAt`/`updatedAt`.
