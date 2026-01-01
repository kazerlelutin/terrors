package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Migration représente une migration
type Migration struct {
	ID   string
	File string
	SQL  string
}

func main() {
	fmt.Println("🎭 Terrors - Migration Tool")
	fmt.Println("==========================")

	// Charger les variables d'environnement
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Fichier .env non trouvé, utilisation des variables système")
	}

	// Récupérer la connexion DB
	dbURL := os.Getenv("PG_URL")
	if dbURL == "" {
		log.Fatal("❌ Variable PG_URL non définie")
	}

	// Se connecter à la base de données
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Erreur de connexion à la DB: %v", err)
	}
	defer db.Close()

	// Tester la connexion
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Erreur de ping de la DB: %v", err)
	}

	fmt.Println("✅ Connexion à la base de données établie")

	// Créer la table de migrations si elle n'existe pas
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("❌ Erreur lors de la création de la table migrations: %v", err)
	}

	// Charger les migrations
	migrations, err := loadMigrations()
	if err != nil {
		log.Fatalf("❌ Erreur lors du chargement des migrations: %v", err)
	}

	// Appliquer les migrations
	if err := applyMigrations(db, migrations); err != nil {
		log.Fatalf("❌ Erreur lors de l'application des migrations: %v", err)
	}

	fmt.Println("🎉 Migrations terminées avec succès !")
}

// createMigrationsTable crée la table pour suivre les migrations
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(query)
	return err
}

// loadMigrations charge les fichiers de migration
func loadMigrations() ([]Migration, error) {
	migrationsDir := "migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du dossier migrations: %w", err)
	}

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(migrationsDir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la lecture de %s: %w", file.Name(), err)
			}

			migration := Migration{
				ID:   strings.TrimSuffix(file.Name(), ".sql"),
				File: file.Name(),
				SQL:  string(content),
			}
			migrations = append(migrations, migration)
		}
	}

	// Trier les migrations par nom de fichier
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].ID < migrations[j].ID
	})

	return migrations, nil
}

// applyMigrations applique les migrations non appliquées
func applyMigrations(db *sql.DB, migrations []Migration) error {
	for _, migration := range migrations {
		// Vérifier si la migration a déjà été appliquée
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE id = $1", migration.ID).Scan(&count)
		if err != nil {
			return fmt.Errorf("erreur lors de la vérification de la migration %s: %w", migration.ID, err)
		}

		if count > 0 {
			fmt.Printf("⏭️  Migration %s déjà appliquée\n", migration.ID)
			continue
		}

		// Appliquer la migration
		fmt.Printf("🔄 Application de la migration %s...\n", migration.ID)

		// Démarrer une transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("erreur lors du début de transaction pour %s: %w", migration.ID, err)
		}

		// Exécuter le SQL de la migration
		_, err = tx.Exec(migration.SQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erreur lors de l'exécution de la migration %s: %w", migration.ID, err)
		}

		// Marquer la migration comme appliquée
		_, err = tx.Exec("INSERT INTO migrations (id) VALUES ($1)", migration.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erreur lors de l'enregistrement de la migration %s: %w", migration.ID, err)
		}

		// Valider la transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("erreur lors de la validation de la migration %s: %w", migration.ID, err)
		}

		fmt.Printf("✅ Migration %s appliquée avec succès\n", migration.ID)
	}

	return nil
}
