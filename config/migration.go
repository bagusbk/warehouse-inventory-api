package config

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration() {
	driver, err := postgres.WithInstance(SQLDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal membuat PostgreSQL driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi migrasi: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Gagal menjalankan migrasi: %v", err)
	}

	log.Println("✔️ Migrasi sukses (Up)")
}
