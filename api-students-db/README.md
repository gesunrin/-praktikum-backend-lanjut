# API Students DB — Praktikum Backend Lanjut Pertemuan 3

Lanjutan dari `api-students` (Pertemuan 2). Data mahasiswa sekarang tersimpan permanen di **MySQL**, bukan lagi di memori, dan diakses lewat pola **repository**.

> Catatan: modul aslinya mencontohkan PostgreSQL + driver `pgx`. Proyek ini memakai **MySQL** + `database/sql` + driver `go-sql-driver/mysql` sebagai penyesuaian, dengan konsep yang identik (connection pool, parameterized query, repository pattern, translate error jadi status HTTP).

## Menyiapkan Basis Data dari Nol

1. Pastikan MySQL sudah berjalan (di proyek ini dijalankan lewat Laragon).
2. Buat database:
   ```
   mysql -u root -p -e "CREATE DATABASE praktikum_backend;"
   ```
3. Jalankan migrasi:
   ```
   mysql -u root -p praktikum_backend < migrations/001_create_students.sql
   ```
4. Cek tabel sudah terbentuk:
   ```
   mysql -u root -p praktikum_backend -e "DESCRIBE students;"
   ```

## Variabel Environment

Salin `.env.example` menjadi `.env`, lalu isi:

| Variabel | Keterangan | Contoh |
|---|---|---|
| APP_PORT | Port server Fiber | 3000 |
| DB_HOST | Host MySQL | localhost |
| DB_PORT | Port MySQL | 3306 |
| DB_USER | Username MySQL | root |
| DB_PASSWORD | Password MySQL | (isi sendiri) |
| DB_NAME | Nama database | praktikum_backend |
| DB_MAX_CONNS | Maksimum koneksi dalam connection pool | 10 |

**`.env` tidak ikut ter-commit** (sudah masuk `.gitignore`). Rekan yang meng-clone repo ini wajib membuat `.env` sendiri dari `.env.example`.

## Menjalankan

```
go mod tidy
go run .
```

## Skema Tabel `students`

```sql
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade DECIMAL(5,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE UNIQUE INDEX students_nim_key ON students (nim);
CREATE INDEX students_name_idx ON students (name);
```

- **`nim` dibuat UNIQUE di level database**, bukan hanya divalidasi di kode Go. Ini mencegah dua permintaan yang datang nyaris bersamaan sama-sama lolos pengecekan sebelum salah satunya sempat tersimpan — masalah race condition yang tidak bisa dicegah kalau keunikan hanya dijaga dengan perulangan cek di Go.
- **Index tambahan pada `name`** (`students_name_idx`) mempercepat pencarian dan pengurutan berdasarkan nama, yang merupakan operasi paling sering dipakai di endpoint daftar.

## Kontrak Endpoint

Sama seperti Pertemuan 2 (lihat detail lengkap di laporan), dengan tambahan: seluruh operasi kini tersimpan permanen dan `/api/v1/health` ikut memeriksa koneksi database.

| Metode | Endpoint | Status Mungkin |
|---|---|---|
| GET | `/api/v1/health` | 200, 503 |
| GET | `/api/v1/students` | 200 |
| GET | `/api/v1/students/:id` | 200, 400, 404 |
| POST | `/api/v1/students` | 201, 400, 415, 422, 409 |
| PUT | `/api/v1/students/:id` | 200, 400, 404, 415, 422, 409 |
| PATCH | `/api/v1/students/:id` | 200, 400, 404, 415, 422 |
| DELETE | `/api/v1/students/:id` | 204, 400, 404 |

## Arsitektur

```
handler → repository (interface) → model
```

- `app/model` — struct data, tidak mengenal apa pun di luar dirinya.
- `app/repository` — interface `StudentRepository` (kontrak) + implementasi MySQL. Tidak ada satu pun penyebutan `fiber` di paket ini.
- `handler.go` — menerjemahkan request HTTP jadi pemanggilan repository, dan error repository jadi status HTTP.