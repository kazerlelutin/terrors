package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Init initialise la connexion à la base de données
func Init() (*sql.DB, error) {
	// Vérifier si on a une URL PostgreSQL complète
	pgURL := os.Getenv("PG_URL")
	if pgURL != "" {
		db, err := sql.Open("postgres", pgURL)
		if err != nil {
			return nil, fmt.Errorf("erreur d'ouverture de la DB avec PG_URL: %w", err)
		}

		// Tester la connexion
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("erreur de ping de la DB: %w", err)
		}

		// Configuration du pool de connexions
		db.SetMaxOpenConns(20)
		db.SetMaxIdleConns(10)

		fmt.Println("Connexion à PostgreSQL établie via PG_URL")
		return db, nil
	}

	url := getEnv("PG_URL", "")

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("erreur d'ouverture de la DB: %w", err)
	}

	// Tester la connexion
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erreur de ping de la DB: %w", err)
	}

	// Configuration du pool de connexions
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	fmt.Println("Connexion à PostgreSQL établie")
	return db, nil
}

// getEnv récupère une variable d'environnement avec une valeur par défaut
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
