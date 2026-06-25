# Backend Modules

Setiap folder di direktori ini mewakili satu kemampuan bisnis dan dapat dibaca
secara mandiri.

Pola file standar:

```text
model.go             entity, request, dan response
service.go           aturan bisnis serta repository interface
repository_gorm.go   akses PostgreSQL menggunakan GORM
http.go              HTTP handler dan pendaftaran route
*_test.go            unit atau integration test modul
```

Urutan membaca modul:

```text
http.go -> service.go -> repository_gorm.go -> model.go
```

Aturan dependency:

- Modul boleh memakai `shared` dan `platform`.
- Modul tidak boleh mengimpor `server`.
- `server` hanya memasangkan route dan middleware.
- `cmd/api` adalah composition root yang membuat repository, service, dan handler.
- Hindari membuat interface global yang berisi operasi dari banyak modul.
