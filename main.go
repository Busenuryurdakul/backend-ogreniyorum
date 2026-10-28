
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	urunler = []string{"Bilgisayar", "Telefon", "Tablet"}
	mu      sync.Mutex
)

// GET - Urunleri listele
func urunleriGetir(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	liste := append([]string{}, urunler...)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(liste)
}

// POST - Urun ekle
func urunEkle(w http.ResponseWriter, r *http.Request) {
	var veri struct {
		Ad string `json:"ad"`
	}

	if err := json.NewDecoder(r.Body).Decode(&veri); err != nil || veri.Ad == "" {
		http.Error(w, "Gecersiz urun", http.StatusBadRequest)
		return
	}

	mu.Lock()
	urunler = append(urunler, veri.Ad)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(veri)
}

// PUT - Urun guncelle
func urunGuncelle(w http.ResponseWriter, r *http.Request) {
	var veri struct {
		EskiAd string `json:"eskiAd"`
		YeniAd string `json:"yeniAd"`
	}

	if err := json.NewDecoder(r.Body).Decode(&veri); err != nil ||
		veri.EskiAd == "" || veri.YeniAd == "" {
		http.Error(w, "Gecersiz veri", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, urun := range urunler {
		if urun == veri.EskiAd {
			urunler[i] = veri.YeniAd
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(veri)
			return
		}
	}

	http.Error(w, "Urun bulunamadi", http.StatusNotFound)
}

// DELETE - Urun sil
func urunSil(w http.ResponseWriter, r *http.Request) {
	var veri struct {
		Ad string `json:"ad"`
	}

	if err := json.NewDecoder(r.Body).Decode(&veri); err != nil || veri.Ad == "" {
		http.Error(w, "Gecersiz veri", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, urun := range urunler {
		if urun == veri.Ad {
			urunler = append(urunler[:i], urunler[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Urun bulunamadi", http.StatusNotFound)
}

func main() {
	http.HandleFunc("GET /urunler", urunleriGetir)
	http.HandleFunc("POST /urunler", urunEkle)
	http.HandleFunc("PUT /urunler", urunGuncelle)
	http.HandleFunc("DELETE /urunler", urunSil)

	fmt.Println("Sunucu 8080 portunda baslatiliyor...")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
