# CRM Backend

Backend Go untuk CRM Enterprise dengan clean architecture, PostgreSQL, Redis,
REST, gRPC, Telegram outbox, CSV/PDF export, dan isolasi multi-tenant.

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

