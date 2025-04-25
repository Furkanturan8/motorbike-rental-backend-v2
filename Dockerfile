# Build aşaması
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Gerekli paketleri yükle
RUN apk add --no-cache gcc musl-dev git

# Go modüllerini kopyala ve indir
COPY go.mod go.sum ./
RUN go mod download

# Tüm kaynak kodlarını kopyala
COPY .. .

# Uygulamayı derle
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/server/main.go

# Çalışma aşaması
FROM alpine:latest

WORKDIR /app

# SSL sertifikaları ve temel araçlar
RUN apk --no-cache add ca-certificates tzdata

# Builder aşamasından gerekli dosyaları kopyala
COPY --from=builder /app/main .

# Log dizini oluştur ve izinleri ayarla
RUN mkdir -p /app/logs

# Port bilgisi (opsiyonel)
EXPOSE 3005

# Uygulamayı çalıştır
CMD ["./main"]