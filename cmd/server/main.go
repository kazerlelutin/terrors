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
			log.Printf("Error loading .env file: %v", err)
		} else {
			log.Println(".env file loaded")
		}
	} else {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := database.Init()
	if err != nil {
		log.Printf("Database not available: %v", err)
		log.Println("Server starts in demo mode (no persistence)")
		db = nil
	} else {
		defer db.Close()
	}

	handlers := api.NewHandlers(db)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/sadako", handlers.Sadako) // Frontend errors
	http.HandleFunc("/jason", handlers.Jason)   // Backend errors
	http.HandleFunc("/cdn/terrors.js", handlers.ServeTerrorsJS)
	http.HandleFunc("/overlook", handlers.Overlook) // Dashboard CQRS
	http.HandleFunc("/api/apps", api.RequireAdminToken(handlers.HandleApps))
	http.HandleFunc("/api/apps/", api.RequireAdminToken(handlers.HandleApps))
	http.HandleFunc("/api/errors", api.RequireAdminToken(handlers.HandleErrors))
	http.HandleFunc("/api/errors/", api.RequireAdminToken(handlers.HandleErrors))
	http.HandleFunc("/api/webhooks", api.RequireAdminToken(handlers.HandleWebhooks))
	http.HandleFunc("/api/webhooks/", api.RequireAdminToken(handlers.HandleWebhooks))

	port := os.Getenv("PORT")
	if port == "" {
		port = "4004"
	}

	log.Printf("🚀 Serveur démarré sur http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
