package main

import "fmt"

// swap menukar nilai dua integer lewat pointer
func swap(a, b *int) {
	*a, *b = *b, *a
}

// updateSlice menambahkan item baru ke slice lewat pointer
// (perlu pointer karena append bisa mengembalikan alamat memori baru
// jika kapasitas slice sudah penuh, dan itu harus tersimpan balik ke variabel asli)
func updateSlice(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

// contoh pass by value vs pass by pointer
func tambahByValue(x int) {
	x = x + 100
}

func tambahByPointer(x *int) {
	*x = *x + 100
}

func main() {
	fmt.Println("=== Swap dengan Pointer ===")
	a, b := 5, 10
	fmt.Println("Sebelum swap:", a, b)
	swap(&a, &b)
	fmt.Println("Sesudah swap :", a, b)

	fmt.Println("\n=== Update Slice dengan Pointer ===")
	daftarBuah := []string{"apel", "jeruk"}
	fmt.Println("Sebelum:", daftarBuah)
	updateSlice(&daftarBuah, "mangga")
	fmt.Println("Sesudah:", daftarBuah)

	fmt.Println("\n=== Pass by Value vs Pass by Pointer ===")
	angka := 50
	tambahByValue(angka)
	fmt.Println("Setelah tambahByValue  :", angka) // tetap 50

	tambahByPointer(&angka)
	fmt.Println("Setelah tambahByPointer:", angka) // jadi 150

	// Penjelasan singkat:
	// tambahByValue menerima SALINAN dari angka, jadi perubahan di dalam
	// function hanya berlaku pada salinan itu dan hilang setelah function selesai.
	// tambahByPointer menerima ALAMAT dari angka, jadi *x = *x + 100 langsung
	// mengubah nilai di alamat memori aslinya, sehingga variabel angka ikut berubah.
}