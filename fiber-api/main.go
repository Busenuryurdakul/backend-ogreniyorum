
package main

import (
	"log"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v3"
)

// Başlangıç ürünlerimiz
var urunler = []string{
	"Bilgisayar",
	"Telefon",
	"Tablet",
}

// Eş zamanlı isteklerde listeyi korur
var mu sync.Mutex

func main() {
	// Fiber uygulamasını oluştur
	app := fiber.New()

	// GET - Ürünleri listele
	app.Get("/urunler", func(c fiber.Ctx) error {
		mu.Lock()
		liste := append([]string{}, urunler...)
		mu.Unlock()

		return c.JSON(liste)
	})

	// POST - Yeni ürün ekle
	app.Post("/urunler", func(c fiber.Ctx) error {
		var yeniUrun struct {
			Ad string `json:"ad"`
		}

		if err := c.Bind().Body(&yeniUrun); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"hata": "Gecersiz JSON verisi",
			})
		}

		yeniUrun.Ad = strings.TrimSpace(yeniUrun.Ad)

		if yeniUrun.Ad == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"hata": "Urun adi bos olamaz",
			})
		}

		mu.Lock()
		urunler = append(urunler, yeniUrun.Ad)
		mu.Unlock()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"mesaj": "Urun basariyla eklendi",
			"urun":  yeniUrun.Ad,
		})
	})

	// PUT - Mevcut ürünü güncelle
	app.Put("/urunler", func(c fiber.Ctx) error {
		var veri struct {
			EskiAd string `json:"eskiAd"`
			YeniAd string `json:"yeniAd"`
		}

		if err := c.Bind().Body(&veri); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"hata": "Gecersiz JSON verisi",
			})
		}

		veri.EskiAd = strings.TrimSpace(veri.EskiAd)
		veri.YeniAd = strings.TrimSpace(veri.YeniAd)

		if veri.EskiAd == "" || veri.YeniAd == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"hata": "Urun adlari bos olamaz",
			})
		}

		mu.Lock()
		defer mu.Unlock()

		for i, urun := range urunler {
			if urun == veri.EskiAd {
				urunler[i] = veri.YeniAd

				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"mesaj":  "Urun basariyla guncellendi",
					"eskiAd": veri.EskiAd,
					"yeniAd": veri.YeniAd,
				})
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"hata": "Urun bulunamadi",
		})
	})

	// DELETE - Mevcut ürünü sil
	app.Delete("/urunler", func(c fiber.Ctx) error {
		var veri struct {
			Ad string `json:"ad"`
		}

		// JSON verisini oku
		if err := c.Bind().Body(&veri); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"hata": "Gecersiz JSON verisi",
			})
		}

		veri.Ad = strings.TrimSpace(veri.Ad)

		if veri.Ad == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"hata": "Urun adi bos olamaz",
			})
		}

		mu.Lock()
		defer mu.Unlock()

		// Silinecek ürünü bul
		for i, urun := range urunler {
			if urun == veri.Ad {
				// Ürünü listeden çıkar
				urunler = append(
					urunler[:i],
					urunler[i+1:]...,
				)

				// Başarılı silme: 204 No Content
				return c.SendStatus(fiber.StatusNoContent)
			}
		}

		// Ürün bulunamadı
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"hata": "Urun bulunamadi",
		})
	})

	log.Println("Fiber sunucusu 8082 portunda baslatiliyor...")

	if err := app.Listen(":8082"); err != nil {
		log.Fatal("Sunucu hatasi:", err)
	}
}
