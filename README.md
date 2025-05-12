# Motosiklet Kiralama Sistemi Backend

Bu proje, motosiklet kiralama işlemlerini yönetmek için geliştirilmiş bir REST API servisidir. Go dilinde Fiber framework'ü kullanılarak geliştirilmiştir.

## Özellikler

- 🔐 Kullanıcı kimlik doğrulama ve yetkilendirme
- 🏍️ Motosiklet yönetimi (ekleme, silme, güncelleme, listeleme)
- 🚦 Sürüş yönetimi (başlatma, bitirme, süre ve ücret hesaplama)
- 📱 Bluetooth bağlantı yönetimi
- 📊 Prometheus ile metrik izleme
- 🔄 Redis önbellek desteği
- 📧 E-posta bildirimleri
- 🔒 JWT tabanlı güvenlik
- ⚡ Rate limiting ve CORS koruması

## API Endpoints

### Kimlik Doğrulama (`/api/v1/auth`)
- `POST /register` - Yeni kullanıcı kaydı
- `POST /login` - Kullanıcı girişi
- `POST /refresh` - Token yenileme
- `POST /forgot-password` - Şifre sıfırlama talebi
- `POST /reset-password` - Şifre sıfırlama
- `POST /logout` - Çıkış yapma

### Kullanıcı İşlemleri (`/api/v1/users`)
- `GET /me` - Kullanıcı profili görüntüleme
- `PUT /me` - Kullanıcı profili güncelleme

#### Admin İşlemleri
- `POST /` - Yeni kullanıcı oluşturma
- `GET /` - Tüm kullanıcıları listeleme
- `GET /:id` - Kullanıcı detayı görüntüleme
- `PUT /:id` - Kullanıcı güncelleme
- `DELETE /:id` - Kullanıcı silme

### Sürüş İşlemleri (`/api/v1/rides`)
- `POST /` - Yeni sürüş başlatma
- `GET /me` - Kullanıcının sürüşlerini listeleme
- `PUT /finish/:id` - Sürüşü bitirme
- `POST /photo/:id` - Sürüş fotoğrafı ekleme

#### Admin İşlemleri
- `GET /` - Tüm sürüşleri listeleme
- `GET /user/:userID` - Kullanıcının sürüşlerini listeleme
- `GET /bike/:motorbikeID` - Motosikletin sürüşlerini listeleme
- `GET /:id` - Sürüş detayı görüntüleme
- `PUT /:id` - Sürüş güncelleme
- `DELETE /:id` - Sürüş silme

### Motosiklet İşlemleri (`/api/v1/motorbike`)
- `GET /` - Tüm motosikletleri listeleme
- `GET /available` - Müsait motosikletleri listeleme
- `GET /:id` - Motosiklet detayı görüntüleme

#### Admin İşlemleri
- `POST /` - Yeni motosiklet ekleme
- `PUT /:id` - Motosiklet güncelleme
- `DELETE /:id` - Motosiklet silme
- `GET /maintenance` - Bakımdaki motosikletleri listeleme
- `GET /rented-motorbikes` - Kiralık motosikletleri listeleme
- `GET /motorbike-photos/:id` - Motosiklet fotoğraflarını görüntüleme

### Bluetooth İşlemleri (`/api/v1/bluetooth`)
#### Admin İşlemleri
- `POST /` - Yeni bluetooth bağlantısı ekleme
- `PUT /:id` - Bluetooth bağlantısı güncelleme
- `DELETE /:id` - Bluetooth bağlantısı silme
- `GET /` - Tüm bluetooth bağlantılarını listeleme
- `GET /:id` - Bluetooth bağlantı detayı görüntüleme

## Teknik Detaylar

- **Framework**: Fiber
- **Veritabanı**: PostgreSQL
- **ORM**: Bun
- **Önbellek**: Redis
- **Monitoring**: Prometheus
- **Güvenlik**: JWT
- **Rate Limiting**: 30 saniyede 10 istek
- **CORS**: Localhost:63342, 3005, 5173 için açık

## Başlangıç

1. Gerekli bağımlılıkları yükleyin:
```bash
go mod download
```

2. Veritabanı bağlantısını yapılandırın:
- PostgreSQL veritabanı oluşturun
- `.env` dosyasında bağlantı bilgilerini güncelleyin

3. Uygulamayı başlatın:
```bash
go run cmd/server/main.go
```

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır. 