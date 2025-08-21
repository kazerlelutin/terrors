package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"terrors/internal/models"
	"time"

	"golang.org/x/exp/rand"
)

// Handlers contient tous les handlers HTTP
type Handlers struct {
	db *sql.DB
}

// NewHandlers crée une nouvelle instance de Handlers
func NewHandlers(db *sql.DB) *Handlers {
	return &Handlers{db: db}
}

// Home gère la route racine
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]string{
		"message": "Welcome to the Overlook Hotel - Error monitoring service",
		"quote":   "All work and no play makes Jack a dull boy",
		"year":    "1980",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Sadako gère l'endpoint pour recevoir les erreurs
func (h *Handlers) Sadako(w http.ResponseWriter, r *http.Request) {
	// Gérer les requêtes OPTIONS (CORS)
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Vérifier la méthode HTTP
	if r.Method != "POST" {
		http.Error(w, "Wrong turn at Albuquerque - Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lire le body de la requête
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "The body is missing - Bad request", http.StatusBadRequest)
		return
	}

	// Parser le JSON
	var errorReq models.ErrorRequest
	if err := json.Unmarshal(body, &errorReq); err != nil {
		http.Error(w, "Something's not right with this JSON - Bad request", http.StatusBadRequest)
		return
	}

	// Afficher l'erreur reçue
	fmt.Printf("🎭 The call is coming from inside the house: %+v\n", errorReq)

	// Si la DB est disponible, sauvegarder l'erreur
	if h.db != nil {
		// TODO: Implémenter la sauvegarde en base
		fmt.Println("🔪 They're here... saving to database")
	} else {
		fmt.Println("👻 Demo mode: error logged to console only")
	}

	// Réponse de succès avec référence film d'horreur
	response := models.ErrorResponse{
		Success:   true,
		Message:   "Error captured and stored in the basement",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Quote:     getRandomHorrorQuote(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}

// ServeTerrorsJS sert le script JavaScript
func (h *Handlers) ServeTerrorsJS(w http.ResponseWriter, r *http.Request) {
	// Lire le fichier JavaScript
	content, err := os.ReadFile("static/terrors.js")
	if err != nil {
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(content)
}

// getRandomHorrorQuote retourne une citation aléatoire de film d'horreur
func getRandomHorrorQuote() string {
	quotes := []string{
		// The Shining (1980)
		"Here's Johnny!",
		"All work and no play makes Jack a dull boy",
		"Come play with us, Danny. Forever... and ever... and ever",

		// Halloween (1978)
		"It's Halloween, everyone's entitled to one good scare",
		"The boogeyman is coming",

		// Friday the 13th (1980)
		"They say the lake has a bottomless depth",
		"Camp Crystal Lake has a death curse",

		// A Nightmare on Elm Street (1984)
		"One, two, Freddy's coming for you",
		"Don't fall asleep",
		"Sweet dreams",

		// The Exorcist (1973)
		"The power of Christ compels you",
		"Your mother sucks cocks in hell",

		// Carrie (1976)
		"They're all gonna laugh at you",
		"Plug it up, plug it up",

		// The Texas Chain Saw Massacre (1974)
		"Who will survive and what will be left of them?",

		// Alien (1979)
		"In space no one can hear you scream",
		"Game over, man, game over",

		// The Thing (1982)
		"Nobody trusts anybody now",
		"Things start getting weird",

		// Poltergeist (1982)
		"They're here",
		"This house is clean",

		// The Evil Dead (1981)
		"Join us",
		"Dead by dawn",

		// Hellraiser (1987)
		"We have such sights to show you",
		"Jesus wept",

		// Child's Play (1988)
		"Hi, I'm Chucky, wanna play?",
		"Don't fuck with the Chuck",

		// Scream (1996)
		"What's your favorite scary movie?",
		"Ghostface is calling",

		// The Blair Witch Project (1999)
		"We're lost",
		"I'm so scared",

		// The Ring (2002)
		"Seven days",
		"Before you die, you see the ring",
	}

	return quotes[rand.Intn(len(quotes))]
}
