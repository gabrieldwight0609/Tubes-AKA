package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Player struct {
	Rank  int
	Name  string
	Score int
}

type Result struct {
	N          int
	RekTime    float64
	IterTime   float64
	TopPlayers []Player
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] >= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func QuickSortRecursive(arr []int, low, high int) {
	if low < high {
		pi := partition(arr, low, high)

		QuickSortRecursive(arr, low, pi-1)
		QuickSortRecursive(arr, pi+1, high)
	}
}

func QuickSortIterative(arr []int, low, high int) {
	stackSize := high - low + 1
	stack := make([]int, stackSize)
	top := -1

	top++
	stack[top] = low
	top++
	stack[top] = high

	for top >= 0 {
		currentH := stack[top]
		top--
		currentL := stack[top]
		top--

		p := partition(arr, currentL, currentH)

		if p-1 > currentL {
			top++
			stack[top] = currentL
			top++
			stack[top] = p - 1
		}
		if p+1 < currentH {
			top++
			stack[top] = p + 1
			top++
			stack[top] = currentH
		}
	}
}

var currentSeed int64 = 42
var history []Result

// ... (Struktur Player, Result, dan fungsi partition, quickSort tetap sama) ...

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Cek apakah yang diklik adalah tombol "Reset"
		if r.FormValue("action") == "reset" {
			history = nil                       // Hapus semua riwayat
			currentSeed = time.Now().UnixNano() // Ganti seed berdasarkan waktu saat ini (acak baru)
		} else {
			// Jika tombol "Jalankan" diklik
			nStr := r.FormValue("n_size")
			n, _ := strconv.Atoi(nStr)

			// Gunakan seed yang sedang aktif (agar data tetap sama untuk n berbeda)
			rand.Seed(currentSeed)

			arr1 := make([]int, n)
			arr2 := make([]int, n)
			for i := 0; i < n; i++ {
				val := rand.Intn(100000)
				arr1[i], arr2[i] = val, val
			}

			s1 := time.Now()
			QuickSortRecursive(arr1, 0, n-1)
			d1 := time.Since(s1).Seconds() * 1000

			s2 := time.Now()
			QuickSortIterative(arr2, 0, n-1)
			d2 := time.Since(s2).Seconds() * 1000

			// Ambil Top 10 dari seed yang sama
			rand.Seed(currentSeed) // Reset seed lagi agar nama player juga konsisten
			var top10 []Player
			limit := 10
			if n < 10 {
				limit = n
			}
			for i := 0; i < limit; i++ {
				top10 = append(top10, Player{
					Rank: i + 1, Name: fmt.Sprintf("Player_%d", rand.Intn(9999)), Score: arr1[i],
				})
			}

			history = append([]Result{{N: n, RekTime: d1, IterTime: d2, TopPlayers: top10}}, history...)
		}
	}

	tmpl, _ := template.ParseFiles("TubesAka.html")
	tmpl.Execute(w, history)
}

func main() {
	http.HandleFunc("/styling.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "styling.css")
	})
	http.HandleFunc("/", handleIndex)
	println("Aplikasi Tubes AKA berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
