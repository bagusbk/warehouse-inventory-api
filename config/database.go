package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SQLDB *sql.DB

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, pakai environment bawaan OS")
	}
}

func DBInit() *gorm.DB {
	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, pass, dbname, port,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal koneksi database (gorm):", err)
	}

	log.Println("✔️ Success Connection to database")

	DB = gormDB

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("❌ Gagal mengambil sql.DB:", err)
	}

	SQLDB = sqlDB

	return gormDB
}
