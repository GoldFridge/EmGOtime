package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDatabase() {
	// Задаем переменные напрямую (замените значения на свои)
	dbHost := "192.168.31.130"
	dbPort := "5432"
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := "postgres"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
}
