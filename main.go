
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
)

var (
    urunler = []string{"Bilgisayar", "Telefon", "Tablet"}
    mu sync.Mutex
)

func urunleriGetir(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    liste := append([]string{}, urunler...)
    mu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(liste)
}

func urunEkle(w http.ResponseWriter, r *http.Request) {
    var yeniUrun struct {
        Ad string `json:"ad"`
    }

    err := json.NewDecoder(r.Body).Decode(&yeniUrun)
    if err != nil || yeniUrun.Ad == "" {
        http.Error(w, "Gecersiz urun", http.StatusBadRequest)
        return
    }

    mu.Lock()
    urunler = append(urunler, yeniUrun.Ad)
    mu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(yeniUrun)
}

func main() {
    http.HandleFunc("GET /urunler", urunleriGetir)
    http.HandleFunc("POST /urunler", urunEkle)

    fmt.Println("Sunucu 8080 portunda baslatiliyor...")
    fmt.Println(http.ListenAndServe(":8080", nil))
}
