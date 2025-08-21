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
	// Initialiser le générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())

	// Charger les variables d'environnement depuis .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Fichier .env non trouvé, utilisation des variables système")
	}

	// Initialiser la base de données (optionnel)
	db, err := database.Init()
	if err != nil {
		log.Printf("⚠️  Base de données non disponible: %v", err)
		log.Println("Le serveur démarre en mode démo (sans persistance)")
		db = nil
	} else {
		defer db.Close()
	}

	// Initialiser les handlers avec la DB (peut être nil)
	handlers := api.NewHandlers(db)

	// Configurer les routes
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/sadako", handlers.Sadako)
	http.HandleFunc("/cdn/terrors.js", handlers.ServeTerrorsJS)

	// Port par défaut
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Serveur démarré sur http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
