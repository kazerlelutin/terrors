package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"terrors/internal/api"
	"terrors/internal/database"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	_ = r

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("⚠️  Erreur lors du chargement du fichier .env: %v", err)
		} else {
			log.Println("✅ Fichier .env chargé")
		}
	} else {
		log.Println("ℹ️  Aucun fichier .env trouvé, utilisation des variables d'environnement système")
	}

	db, err := database.Init()
	if err != nil {
		log.Printf("⚠️  Base de données non disponible: %v", err)
		log.Println("Le serveur démarre en mode démo (sans persistance)")
		db = nil
	} else {
		defer db.Close()
	}

	handlers := api.NewHandlers(db)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/sadako", handlers.Sadako)
	http.HandleFunc("/cdn/terrors.js", handlers.ServeTerrorsJS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Serveur démarré sur http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
