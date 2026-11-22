# Go Backend Öğreniyorum

Go programlama diliyle geliştirilen, üç farklı HTTP yaklaşımını karşılaştıran bir REST API projesi.

Bu projede Go'nun standart `net/http` paketi, Gin Framework ve Fiber v3 kullanılarak ürün yönetimi API'leri geliştirilmiştir.

## Kullanılan Teknolojiler

- Go 1.26
- Go `net/http`
- Gin v1.12.0
- Fiber v3.7.0
- Docker
- Postman
- hey (HTTP yük testi aracı)
- Git ve GitHub

## Proje Yapısı

```
backend-ogreniyorum/
├── main.go
├── gin-api/
│   └── main.go
├── fiber-api/
│   └── main.go
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## REST API İşlemleri

Üç uygulamada da aşağıdaki HTTP metotları bulunmaktadır:


| HTTP Metodu | Endpoint   | Açıklama               |
| ----------- | ---------- | ---------------------- |
| GET         | `/urunler` | Ürünleri listeler      |
| POST        | `/urunler` | Yeni ürün ekler        |
| PUT         | `/urunler` | Mevcut ürünü günceller |
| DELETE      | `/urunler` | Ürün siler             |


Ürünler uygulama belleğinde tutulur. Sunucu yeniden başlatıldığında ürün listesi başlangıç değerlerine döner.

## Docker ile Çalıştırma

Komutlar proje ana klasöründe, Windows PowerShell üzerinden çalıştırılmalıdır. Docker Desktop'ın açık olması gerekir.

### 1. Go net/http — Port 8080

```
docker run --rm -it --name backend-go-api -p 8080:8080 -v "${PWD}:/app" -w /app golang:1.26 go run .
```

### 2. Gin Framework — Port 8081

```
docker run --rm -it --name backend-gin-api -p 8081:8081 -e GIN_MODE=release -v "${PWD}:/app" -w /app golang:1.26 go run ./gin-api
```

### 3. Fiber v3 — Port 8082

```
docker run --rm -it --name backend-fiber-api -p 8082:8082 -v "${PWD}:/app" -w /app golang:1.26 go run ./fiber-api
```

Sunucular `Ctrl + C` ile durdurulabilir.

## Postman ile API Testleri

Her uygulamanın kendi portu üzerinden istek gönderilebilir.

**GET — Ürünleri listeleme**

```
GET http://localhost:8082/urunler
```

**POST — Ürün ekleme**

```
POST http://localhost:8082/urunler
Content-Type: application/json
```

```
{
  "ad": "Mekanik Klavye"
}
```

**PUT — Ürün güncelleme**

```
PUT http://localhost:8082/urunler
Content-Type: application/json
```

```
{
  "eskiAd": "Mekanik Klavye",
  "yeniAd": "Kablosuz Klavye"
}
```

**DELETE — Ürün silme**

```
DELETE http://localhost:8082/urunler
Content-Type: application/json
```

```
{
  "ad": "Kablosuz Klavye"
}
```

Yukarıdaki örnekler Fiber içindir. Gin için `8081`, Go `net/http` için `8080` portu kullanılır. Go `net/http` uygulamasında PUT ve DELETE URL veya veri yapısı farklı olabileceğinden, bu iki metodun testi ilgili uygulamanın kodundaki endpoint biçimine göre yapılmalıdır.

## Performans Testleri

Performans testleri Docker içinde çalışan `hey` aracıyla gerçekleştirilmiştir.

**Test koşulları:**

- HTTP metodu: GET
- Endpoint: `/urunler`
- Test süresi: 15 saniye
- Eşzamanlı istek sayısı: 10

### Ölçülen Sonuçlar


| Uygulama                      | İstek/sn (RPS) | Ortalama yanıt süresi | Başarılı istek |
| ----------------------------- | -------------- | --------------------- | -------------- |
| Go net/http                   | 39.304,93      | 0,3 ms                | 589.594        |
| Gin (Debug, loglama açık)     | 21.045,19      | 0,5 ms                | 315.717        |
| Gin (Release, loglama açık)   | 21.139,70      | 0,5 ms                | 317.134        |
| Gin (Release, loglama kapalı) | 37.472,11      | 0,3 ms                | 562.124        |
| Fiber v3                      | 44.749,30      | 0,2 ms                | 671.298        |


Bu ilk ölçümlerde Fiber v3 en yüksek istek hızına ulaşmıştır. Gin'in istek loglaması kapatıldığında performansında belirgin iyileşme gözlemlenmiştir.

**Not:** Sonuçlar belirli bir bilgisayarda ve Docker ortamında yapılan tekil ölçümlere aittir. Kesin ve genellenebilir bir framework sıralaması olarak değerlendirilmemelidir. Daha güvenilir sonuçlar için tekrar testleri ve eşit yapılandırma koşulları gereklidir.

## Proje Amacı

Bu projeyle:

- Go ile REST API geliştirme
- Gin ve Fiber framework'lerini kullanma
- HTTP metotlarını uygulama
- JSON verisi işleme
- Docker ile uygulama çalıştırma
- Postman ile API testleri
- HTTP performans ölçümü
- Git ve GitHub ile sürüm kontrolü

konularında uygulamalı deneyim kazanılması amaçlanmıştır.