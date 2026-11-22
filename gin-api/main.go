
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// Ortak ürün listesi
var urunler = []string{
	"Bilgisayar",
	"Telefon",
	"Tablet",
}

// Eş zamanlı isteklerde veriyi korur
var mu sync.Mutex

func main() {

	// Gin uygulamasını oluştur
	// Performans testi için istek loglaması kapalı
	router := gin.New()

	// Beklenmeyen hatalarda sunucuyu korur
	router.Use(gin.Recovery())

	// GET - Ürünleri listele
	router.GET("/urunler", func(c *gin.Context) {
		mu.Lock()
		liste := append([]string{}, urunler...)
		mu.Unlock()

		c.JSON(http.StatusOK, liste)
	})

	// POST - Yeni ürün ekle
	router.POST("/urunler", func(c *gin.Context) {
		var yeniUrun struct {
			Ad string `json:"ad"`
		}

		if err := c.ShouldBindJSON(&yeniUrun); err != nil ||
			yeniUrun.Ad == "" {

			c.JSON(http.StatusBadRequest, gin.H{
				"hata": "Gecersiz urun bilgisi",
			})
			return
		}

		mu.Lock()
		urunler = append(urunler, yeniUrun.Ad)
		mu.Unlock()

		c.JSON(http.StatusCreated, gin.H{
			"mesaj": "Urun basariyla eklendi",
			"urun":  yeniUrun.Ad,
		})
	})

	// PUT - Mevcut ürünü güncelle
	router.PUT("/urunler", func(c *gin.Context) {
		var veri struct {
			EskiAd string `json:"eskiAd"`
			YeniAd string `json:"yeniAd"`
		}

		if err := c.ShouldBindJSON(&veri); err != nil ||
			veri.EskiAd == "" || veri.YeniAd == "" {

			c.JSON(http.StatusBadRequest, gin.H{
				"hata": "Gecersiz urun bilgisi",
			})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		for i, urun := range urunler {
			if urun == veri.EskiAd {
				urunler[i] = veri.YeniAd

				c.JSON(http.StatusOK, gin.H{
					"mesaj":  "Urun basariyla guncellendi",
					"eskiAd": veri.EskiAd,
					"yeniAd": veri.YeniAd,
				})
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"hata": "Urun bulunamadi",
		})
	})

	// DELETE - Mevcut ürünü sil
	router.DELETE("/urunler", func(c *gin.Context) {
		var veri struct {
			Ad string `json:"ad"`
		}

		if err := c.ShouldBindJSON(&veri); err != nil ||
			veri.Ad == "" {

			c.JSON(http.StatusBadRequest, gin.H{
				"hata": "Gecersiz urun bilgisi",
			})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		for i, urun := range urunler {
			if urun == veri.Ad {

				urunler = append(
					urunler[:i],
					urunler[i+1:]...,
				)

				c.Status(http.StatusNoContent)
				return
			}
		}

		c.JSON(http.StatusNotFound, gin.H{
			"hata": "Urun bulunamadi",
		})
	})

	// Sunucuyu başlat
	log.Println("Gin sunucusu 8081 portunda baslatiliyor...")

	if err := router.Run(":8081"); err != nil {
		log.Fatal("Sunucu hatasi:", err)
	}
}
