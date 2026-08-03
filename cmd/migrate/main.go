package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env не загружен (используются переменные окружения): %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("Ошибка создания миграции: %v", err)
	}
	defer func() {
		if _, cerr := m.Close(); cerr != nil {
			log.Printf("Ошибка закрытия миграции: %v", cerr)
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка применения миграции: %v", err)
	}

	log.Println("Миграции успешно применены")
}
