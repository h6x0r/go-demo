package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}

	log.Println("Успешное подключение к базе данных PostgreSQL")

	DB = db
}

func Migrate() {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			firstname VARCHAR NOT NULL,
			lastname  VARCHAR NOT NULL,
			email     VARCHAR NOT NULL UNIQUE,
			age       INT NOT NULL,
			created   TIMESTAMP NOT NULL
		);
		`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка миграции (создание таблицы users): %v", err)
	}

	log.Println("Миграция успешно выполнена (таблица users проверена/создана).")
}
