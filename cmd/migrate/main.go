package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/config"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"os"
)

/*
=======KULLANIM=======
# Tüm migration'ları uygula
go run cmd/migrate/main.go -action up

# Son 2 migration'ı geri al
go run cmd/migrate/main.go -action down -step 2

# Migration listesini görüntüle
go run cmd/migrate/main.go -action status
*/

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Config yükleme hatası: %v", err)
		os.Exit(1)
	}

	// varsayılan olarak "up" işlemi yapılıyor => go run cmd/migrate/main.go
	var (
		action = flag.String("action", "up", "Migration action (up/down/status)")
		step   = flag.Int("step", 0, "Number of migrations (0 for all)")
	)
	flag.Parse()

	// Veritabanı bağlantısı
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DBConfig.GetDSN())))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer func(db *bun.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("Veritabanı kapatma hatası: %v\n", err)
		}
	}(db)

	ctx := context.Background()

	// Migration işlemini gerçekleştir
	switch *action {
	case "up":
		if err = migrations.Up(ctx, db, *step); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migration başarıyla tamamlandı")

	case "down":
		if err = migrations.Down(ctx, db, *step); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migration başarıyla geri alındı")

	case "status":
		status, err := migrations.Status(ctx, db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(status)

	default:
		log.Fatal("Geçersiz işlem. 'up', 'down' veya 'status' kullanın")
	}
}
