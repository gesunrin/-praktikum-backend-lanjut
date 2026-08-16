package main

import "fmt"

type Student struct {
	ID       int
	Name     string
	Grade    float64
	IsActive bool
}

// GetInfo pakai pointer receiver: konsisten dengan method lain di struct ini
// dan menghindari copy struct yang tidak perlu tiap kali dipanggil
func (s *Student) GetInfo() string {
	status := "tidak aktif"
	if s.IsActive {
		status = "aktif"
	}
	return fmt.Sprintf("ID:%d | %s | Nilai:%.2f | Status:%s", s.ID, s.Name, s.Grade, status)
}

// UpdateGrade WAJIB pointer receiver karena mengubah data asli struct
func (s *Student) UpdateGrade(grade float64) {
	s.Grade = grade
}

// Activate & Deactivate WAJIB pointer receiver, sama-sama mengubah field IsActive
func (s *Student) Activate() {
	s.IsActive = true
}

func (s *Student) Deactivate() {
	s.IsActive = false
}

func main() {
	student := &Student{
		ID:       1,
		Name:     "Rangga",
		Grade:    0,
		IsActive: false,
	}

	fmt.Println("Sebelum:", student.GetInfo())

	student.Activate()
	student.UpdateGrade(87.5)
	fmt.Println("Sesudah aktivasi & update nilai:", student.GetInfo())

	student.UpdateGrade(91.0)
	fmt.Println("Setelah update nilai lagi      :", student.GetInfo())

	student.Deactivate()
	fmt.Println("Setelah nonaktif               :", student.GetInfo())
}