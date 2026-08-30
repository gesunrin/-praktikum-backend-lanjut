CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade DECIMAL(5,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- NIM wajib unik. MySQL secara default pakai collation yang tidak
-- membedakan huruf besar/kecil (case-insensitive), jadi UNIQUE biasa
-- di kolom nim sudah cukup, tidak perlu dibungkus LOWER() seperti di PostgreSQL.
CREATE UNIQUE INDEX students_nim_key ON students (nim);

-- Index tambahan untuk mempercepat pencarian & pengurutan berdasarkan nama.
CREATE INDEX students_name_idx ON students (name);