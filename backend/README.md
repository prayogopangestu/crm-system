# CRM Backend

Backend Go untuk CRM Enterprise dengan clean architecture, GORM/PostgreSQL, Redis,
REST, gRPC, Telegram outbox, CSV/PDF export, dan isolasi multi-tenant.

## Alur arsitektur

```text
HTTP/gRPC handler
    -> kontrak use case per fitur
    -> implementasi use case
    -> kontrak repository per fitur
    -> repository GORM
    -> PostgreSQL
```

Pembagian utamanya:

- `internal/domain`: entity dan kontrak repository. Tidak boleh mengimpor GORM.
- `internal/usecase`: aturan bisnis dan kontrak use case per fitur.
- `internal/delivery`: HTTP/gRPC; hanya bergantung pada kontrak use case.
- `internal/repository/postgres`: model persistence dan query GORM.
- `cmd/api`: composition root untuk memasangkan seluruh dependency.

Saat menambah fitur baru, buat domain model, repository interface, use case
interface/implementation, repository GORM, lalu handler. Hindari menambahkan
semua operasi ke satu interface atau handler global.

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
