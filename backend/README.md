# CRM Backend

Backend Go untuk CRM Enterprise dengan modular architecture, GORM/PostgreSQL, Redis,
REST, gRPC, Telegram outbox, CSV/PDF export, dan isolasi multi-tenant.

## Struktur modular

```text
internal/
├── modules/
│   ├── user/
│   ├── contact/
│   ├── deal/
│   ├── task/
│   ├── analytics/
│   ├── pipeline/
│   ├── integration/
│   ├── notification/
│   └── search/
├── platform/
│   └── postgresx/
├── server/
│   ├── httpserver/
│   └── grpcserver/
└── shared/
    └── httpx/
```

Setiap folder modul mengikuti pola yang sama:

```text
modules/deal/
├── model.go             # entity dan request model
├── service.go           # aturan bisnis dan repository interface
├── repository_gorm.go   # implementasi persistence GORM
├── http.go              # handler dan route HTTP
└── service_test.go
```

Alur membaca satu fitur:

```text
http.go -> service.go -> repository_gorm.go -> PostgreSQL
```

Komponen lintas modul ditempatkan di `shared`, koneksi dan helper database di
`platform`, sedangkan bootstrap HTTP/gRPC berada di `server`. `cmd/api` hanya
bertugas memasangkan dependency.

## Menjalankan lokal

```bash
cp configs/.env.example .env
docker compose -f deployments/docker-compose.yml up --build
```

REST tersedia pada `http://localhost:8080`, gRPC pada `localhost:9090`, dan
OpenAPI berada di `api/swagger.yaml`.

Tanpa Docker:

```bash
go run ./cmd/migrate -direction up
go run ./cmd/api
```

Seeder demo dijalankan secara eksplisit:

```bash
SEED_ADMIN_EMAIL=admin@example.com SEED_ADMIN_PASSWORD=Admin123! go run ./cmd/seed
```

Lihat `endpoint-golang.md` untuk keputusan kontrak, model data, cache, RBAC,
dan catatan integrasi frontend.
