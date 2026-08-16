package main

import "fmt"

type Mahasiswa struct {
	Nilai float64
}

func main() {
	// 5 variabel tipe berbeda
	var nama string = "Rangga"
	var umur int = 21
	var ipk float64 = 3.75
	var aktif bool = true
	var hobi []string = []string{"ngoding", "futsal", "nonton"}

	fmt.Println("=== Variabel ===")
	fmt.Println("Nama :", nama)
	fmt.Println("Umur :", umur)
	fmt.Println("IPK  :", ipk)
	fmt.Println("Aktif:", aktif)
	fmt.Println("Hobi :", hobi)

	// Map data mahasiswa: nama -> nilai
	nilaiMahasiswa := make(map[string]float64)

	// Menambah
	nilaiMahasiswa["Rangga"] = 88.5
	nilaiMahasiswa["Sari"] = 92.0
	nilaiMahasiswa["Budi"] = 75.0

	fmt.Println("\n=== Map Mahasiswa ===")

	// Membaca dengan pengecekan keberadaan
	if nilai, ada := nilaiMahasiswa["Sari"]; ada {
		fmt.Println("Nilai Sari ditemukan:", nilai)
	} else {
		fmt.Println("Sari belum punya nilai")
	}

	if nilai, ada := nilaiMahasiswa["Dewi"]; ada {
		fmt.Println("Nilai Dewi:", nilai)
	} else {
		fmt.Println("Dewi belum punya nilai di map")
	}

	// Menghapus
	delete(nilaiMahasiswa, "Budi")
	fmt.Println("Budi setelah dihapus, masih ada?", func() bool {
		_, ada := nilaiMahasiswa["Budi"]
		return ada
	}())

	// Menelusuri seluruh isi
	fmt.Println("\nSeluruh data mahasiswa:")
	for nama, nilai := range nilaiMahasiswa {
		fmt.Printf("- %s: %.2f\n", nama, nilai)
	}
}