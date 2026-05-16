package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

var DB *gorm.DB

func Connect() {
	// 1. Muat file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Error memuat file .env")
	}

	// 2. Ambil hanya nilai password dari .env
	dbPassword := os.Getenv("DB_PASSWORD")

	// 3. Susun string DSN, gunakan %s untuk tempat masuknya dbPassword
	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=assesment_db port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbPassword)

	// 4. Lakukan koneksi ke database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	log.Println("Koneksi ke database berhasil!")
}
