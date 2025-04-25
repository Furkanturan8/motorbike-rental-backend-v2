package migrations

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

// Migration bir migration'ı temsil eder
type Migration struct {
	Version string
	Up      string
	Down    string
}

var Migrations = []Migration{}

// Up belirtilen sayıda migration'ı yukarı doğru çalıştırır
func Up(ctx context.Context, db *bun.DB, step int) error {
	// Çalıştırılacak migration sayısını belirle
	migrationsToRun := Migrations
	if step > 0 && step < len(Migrations) {
		migrationsToRun = Migrations[:step]
	}

	// Her migration için
	for _, migration := range migrationsToRun {
		// Transaction başlat
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("transaction başlatılamadı: %v", err)
		}

		// Migration'ı uygula
		if _, err = tx.ExecContext(ctx, migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration uygulanamadı (version %s): %v", migration.Version, err)
		}

		// Transaction'ı commit et
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("transaction commit edilemedi: %v", err)
		}

		fmt.Printf("Migration uygulandı: %s\n", migration.Version)
	}

	return nil
}

// Down belirtilen sayıda migration'ı geri alır
func Down(ctx context.Context, db *bun.DB, step int) error {
	// Geri alınacak migration sayısını belirle
	totalMigrations := len(Migrations)
	startIndex := totalMigrations - 1
	endIndex := 0

	if step > 0 {
		endIndex = totalMigrations - step
		if endIndex < 0 {
			endIndex = 0
		}
	}

	// Migration'ları sondan başa doğru geri al
	for i := startIndex; i >= endIndex; i-- {
		migration := Migrations[i]

		// Transaction başlat
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("transaction başlatılamadı: %v", err)
		}

		// Migration'ı geri al
		if _, err = tx.ExecContext(ctx, migration.Down); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration geri alınamadı (version %s): %v", migration.Version, err)
		}

		// Transaction'ı commit et
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("transaction commit edilemedi: %v", err)
		}

		fmt.Printf("Migration geri alındı: %s\n", migration.Version)
	}

	return nil
}

// Status mevcut migration'ların listesini gösterir
func Status(ctx context.Context, db *bun.DB) (string, error) {
	status := "Mevcut Migration'lar:\n\n"
	for _, m := range Migrations {
		status += fmt.Sprintf("- %s\n", m.Version)
	}
	return status, nil
}
