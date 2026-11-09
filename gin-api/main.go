
package main

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    router.GET("/urunler", func(c *gin.Context) {
        urunler := []string{
            "Bilgisayar",
            "Telefon",
            "Tablet",
        }

        c.JSON(http.StatusOK, urunler)
    })

    log.Println("Gin sunucusu 8081 portunda baslatiliyor...")

    if err := router.Run(":8081"); err != nil {
        log.Println("Sunucu hatasi:", err)
    }
}
