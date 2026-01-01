package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Init() (*sql.DB, error) {
	pgURL := os.Getenv("PG_URL")

	db, err := sql.Open("postgres", pgURL)
	if err != nil {
		return nil, fmt.Errorf("erreur d'ouverture de la DB avec PG_URL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erreur de ping de la DB: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	fmt.Println("Connexion à PostgreSQL établie via PG_URL")
	return db, nil
}
