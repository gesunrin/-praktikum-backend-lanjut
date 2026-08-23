# API Students — Praktikum Backend Lanjut Pertemuan 2

REST API sederhana untuk mengelola data mahasiswa. Data disimpan di memori (belum pakai database), jadi data akan hilang tiap server dimatikan.

## Menjalankan

```
go run .
```

Server berjalan di `http://localhost:3000`.

## Bentuk Respons

Semua endpoint mengembalikan bentuk yang sama:

```json
// Berhasil
{ "success": true, "message": "...", "data": { ... } }

// Berhasil, berupa daftar
{ "success": true, "message": "...", "data": [ ... ], "meta": { "page": 1, "limit": 10, "total": 3, "total_pages": 1 } }

// Gagal
{ "success": false, "message": "..." }

// Gagal validasi
{ "success": false, "message": "validasi gagal", "errors": { "nim": "wajib diisi" } }
```

## Kontrak Endpoint

| Metode | Endpoint | Parameter | Contoh Body Request | Status Mungkin | Contoh Response |
|---|---|---|---|---|---|
| GET | `/api/v1/students` | Query: `page` (default 1), `limit` (default 10, maks 50), `search` (nama, tidak case-sensitive), `sort` (`id`\|`nim`\|`name`\|`grade`\|`created_at`), `order` (`asc`\|`desc`), `is_active` (`true`\|`false`) | — | 200 | `{"success":true,"message":"daftar mahasiswa berhasil diambil","data":[...],"meta":{"page":1,"limit":10,"total":3,"total_pages":1}}` |
| GET | `/api/v1/students/:id` | Path: `id` (angka) | — | 200, 400, 404 | `{"success":true,"message":"mahasiswa ditemukan","data":{"id":1,"nim":"12345","name":"Rangga","grade":85,"is_active":true,"created_at":"..."}}` |
| POST | `/api/v1/students` | — | `{"nim":"12345","name":"Rangga","grade":85}` | 201, 400, 415, 422, 409 | `{"success":true,"message":"mahasiswa berhasil dibuat","data":{...}}` (header `Location: /api/v1/students/1`) |
| PUT | `/api/v1/students/:id` | Path: `id` (angka). Semua field body wajib diisi | `{"nim":"12345","name":"Rangga","grade":90,"is_active":true}` | 200, 400, 404, 415, 422 | `{"success":true,"message":"mahasiswa berhasil diganti seluruhnya","data":{...}}` |
| PATCH | `/api/v1/students/:id` | Path: `id` (angka). Kirim hanya field yang ingin diubah | `{"is_active":false}` | 200, 400, 404, 415, 422 | `{"success":true,"message":"mahasiswa berhasil diperbarui sebagian","data":{...}}` |
| DELETE | `/api/v1/students/:id` | Path: `id` (angka) | — | 204, 400, 404 | *(tanpa body)* |
| GET | `/api/v1/health` | — | — | 200 | `{"success":true,"message":"server berjalan","data":{"timestamp":"..."}}` |

## Daftar Status HTTP yang Dipakai

| Status | Kapan Terjadi |
|---|---|
| 200 OK | Pengambilan atau perubahan data berhasil |
| 201 Created | Mahasiswa baru berhasil dibuat (disertai header `Location`) |
| 204 No Content | Mahasiswa berhasil dihapus |
| 400 Bad Request | Body bukan JSON yang valid, atau `id` di path bukan angka |
| 404 Not Found | Mahasiswa dengan `id` tersebut tidak ada |
| 409 Conflict | NIM yang dikirim sudah dipakai mahasiswa lain |
| 415 Unsupported Media Type | Header `Content-Type` bukan `application/json` pada POST/PUT/PATCH |
| 422 Unprocessable Entity | Validasi isi gagal (field kosong, format salah, di luar rentang) |

## Field Struct Student

| Field | Tipe | Keterangan |
|---|---|---|
| id | int | Dibuat otomatis oleh server |
| nim | string | Unik, wajib diisi |
| name | string | Wajib diisi |
| grade | float64 | 0–100 |
| is_active | bool | Default `true` saat dibuat |
| created_at | time.Time | Diisi otomatis oleh server |

## Catatan Desain

- **Batas atas `limit` = 50**: data mahasiswa dalam satu kelas/angkatan realistis tidak akan sampai ribuan per halaman, jadi 50 sudah cukup longgar sambil tetap mencegah klien meminta seluruh data sekaligus dan membebani server.
- **Daftar putih `sort`**: hanya field yang memang ada (`id`, `nim`, `name`, `grade`, `created_at`) yang diterima; field lain otomatis jatuh ke `id` supaya klien tidak bisa menyisipkan nama field sembarangan.
- **PATCH memakai pointer** (`*string`, `*float64`, `*bool`) supaya bisa membedakan "field tidak dikirim" (nil) dengan "field dikirim bernilai kosong/false/0".