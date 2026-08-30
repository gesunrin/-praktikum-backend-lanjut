package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"api-students-db/app/model"
)

// Sentinel error milik lapisan repository, bukan milik driver database.
var (
	ErrNotFound = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("data sudah ada")
)

// KONTRAK — tidak ada satu pun kata SQL atau mysql di sini.
type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

// Daftar putih kolom yang boleh dipakai untuk ORDER BY.
// Nama kolom tidak bisa dikirim sebagai parameter, jadi ini satu-satunya
// pertahanan dari SQL injection lewat parameter sort.
var kolomUrut = map[string]string{
	"id":         "id",
	"nim":        "nim",
	"name":       "name",
	"grade":      "grade",
	"created_at": "created_at",
}

type studentMySQLRepository struct {
	db *sql.DB
}

// NewStudentRepository mengembalikan interface, bukan struct konkret.
func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentMySQLRepository{db: db}
}

// buildFilter menyusun bagian WHERE beserta argumennya.
// Nilai dari klien selalu jadi argumen (?), tidak pernah disambung ke teks SQL.
func buildFilter(q model.ListQuery) (string, []any) {
	where := " WHERE 1 = 1"
	args := []any{}

	if q.Search != "" {
		where += " AND name LIKE ?"
		args = append(args, "%"+q.Search+"%")
	}
	if q.IsActive != nil {
		where += " AND is_active = ?"
		args = append(args, *q.IsActive)
	}
	return where, args
}

func (r *studentMySQLRepository) FindAll(
	ctx context.Context, q model.ListQuery,
) ([]model.Student, int, error) {
	where, args := buildFilter(q)

	// 1) Hitung total sebelum dipenggal, untuk keperluan meta.
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM students"+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("menghitung student: %w", err)
	}

	// 2) Ambil satu halaman saja. Penyaringan, pengurutan, dan pemenggalan
	// dikerjakan basis data, bukan oleh Go.
	arah := "ASC"
	if q.Order == "desc" {
		arah = "DESC"
	}
	sqlText := fmt.Sprintf(
		`SELECT id, nim, name, grade, is_active, created_at
		 FROM students%s
		 ORDER BY %s %s
		 LIMIT ? OFFSET ?`,
		where, kolomUrut[q.Sort], arah,
	)
	args = append(args, q.Limit, q.Offset())

	rows, err := r.db.QueryContext(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("mengambil daftar student: %w", err)
	}
	defer rows.Close()

	hasil := []model.Student{}
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("membaca baris student: %w", err)
		}
		hasil = append(hasil, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("membaca hasil query: %w", err)
	}

	return hasil, total, nil
}

func (r *studentMySQLRepository) FindByID(ctx context.Context, id int) (model.Student, error) {
	var s model.Student
	err := r.db.QueryRowContext(ctx,
		`SELECT id, nim, name, grade, is_active, created_at
		 FROM students WHERE id = ?`, id,
	).Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengambil student: %w", err)
	}
	return s, nil
}

func (r *studentMySQLRepository) Create(ctx context.Context, s model.Student) (model.Student, error) {
	// MySQL tidak punya RETURNING, jadi id diambil lewat LastInsertId,
	// dan created_at diambil dengan query tambahan sesudahnya.
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO students (nim, name, grade, is_active) VALUES (?, ?, ?, ?)`,
		s.NIM, s.Name, s.Grade, s.IsActive,
	)
	if err != nil {
		if isDuplicateError(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("menyimpan student: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Student{}, fmt.Errorf("mengambil id baru: %w", err)
	}

	return r.FindByID(ctx, int(id))
}

func (r *studentMySQLRepository) Update(ctx context.Context, s model.Student) (model.Student, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE students SET nim = ?, name = ?, grade = ?, is_active = ? WHERE id = ?`,
		s.NIM, s.Name, s.Grade, s.IsActive, s.ID,
	)
	if err != nil {
		if isDuplicateError(err) {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("memperbarui student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Student{}, fmt.Errorf("membaca hasil update: %w", err)
	}
	if rowsAffected == 0 {
		return model.Student{}, ErrNotFound
	}

	return r.FindByID(ctx, s.ID)
}

func (r *studentMySQLRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM students WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("menghapus student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("membaca hasil delete: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// isDuplicateError memeriksa apakah error berasal dari pelanggaran UNIQUE.
// Kode 1062 adalah kode resmi MySQL untuk itu (setara 23505 di PostgreSQL).
func isDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}