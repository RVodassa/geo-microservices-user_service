package sql

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}

	dbUser, dbPassword, dbHost, dbName, dbPort :=
		os.Getenv("db_user"),
		os.Getenv("db_password"),
		os.Getenv("db_host"),
		os.Getenv("db_name"),
		os.Getenv("db_port")

	// Формирование строки подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Открытие соединения
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение с базой данных: %w", err)
	}

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("проверка соединения не удалась: %w", err)
	}

	// Создание таблицы
	err = CreateUserTable(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	log.Println("База данных готова к работе")
	return db, nil
}

func CreateUserTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        login TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    )`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
