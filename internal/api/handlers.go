package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"terrors/internal/models"
	"terrors/internal/services"
	"time"

	"golang.org/x/exp/rand"
)

type Handlers struct {
	db             *sql.DB
	appService     *services.AppService
	errorService   *services.ErrorService
	webhookService *services.WebhookService
}

func NewHandlers(db *sql.DB) *Handlers {
	var appService *services.AppService
	var errorService *services.ErrorService
	var webhookService *services.WebhookService
	if db != nil {
		appService = services.NewAppService(db)
		errorService = services.NewErrorService(db)
		webhookService = services.NewWebhookService(db)
	}
	return &Handlers{
		db:             db,
		appService:     appService,
		errorService:   errorService,
		webhookService: webhookService,
	}
}

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

func (h *Handlers) Sadako(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Wrong turn at Albuquerque - Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "The body is missing - Bad request", http.StatusBadRequest)
		return
	}

	var errorReq models.ErrorRequest
	if err := json.Unmarshal(body, &errorReq); err != nil {
		http.Error(w, "Something's not right with this JSON - Bad request", http.StatusBadRequest)
		return
	}

	if h.db != nil && h.appService != nil {
		if err := h.appService.ValidateApp(errorReq.AppID); err != nil {
			http.Error(w, "Invalid app", http.StatusForbidden)
			return
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			referer := r.Header.Get("Referer")
			if referer != "" {
				parts := strings.Split(referer, "/")
				if len(parts) >= 3 {
					origin = parts[0] + "//" + parts[2]
				}
			}
		}
		if origin != "" {
			if err := h.appService.ValidateOrigin(errorReq.AppID, origin); err != nil {
				http.Error(w, "Origin not allowed", http.StatusForbidden)
				return
			}
		}

		if h.errorService != nil {
			if errorReq.Type == "" || errorReq.Type == "error" {
				errorReq.Type = "frontend"
			}

			savedError, err := h.errorService.SaveError(errorReq)
			if err != nil {
				fmt.Printf("❌ Erreur lors de la sauvegarde: %v\n", err)
			} else {
				fmt.Printf("✅ Erreur frontend sauvegardée: ID=%d\n", savedError.ID)
				go h.triggerWebhooks(errorReq.AppID, savedError)
			}
		} else {
			fmt.Println("🔪 They're here... saving to database")
		}
	} else {
		fmt.Println("👻 Demo mode: error logged to console only")
	}

	fmt.Printf("🎭 The call is coming from inside the house: %+v\n", errorReq)

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

func (h *Handlers) Jason(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var errorReq models.ErrorRequest
	if err := json.Unmarshal(body, &errorReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if h.db != nil && h.appService != nil {
		if err := h.appService.ValidateApp(errorReq.AppID); err != nil {
			http.Error(w, "Invalid app", http.StatusForbidden)
			return
		}

		if h.errorService != nil {
			if errorReq.Type == "" || errorReq.Type == "error" {
				errorReq.Type = "backend"
			}

			savedError, err := h.errorService.SaveError(errorReq)
			if err != nil {
				fmt.Printf("❌ Erreur lors de la sauvegarde: %v\n", err)
			} else {
				fmt.Printf("✅ Erreur backend sauvegardée: ID=%d\n", savedError.ID)
				go h.triggerWebhooks(errorReq.AppID, savedError)
			}
		}
	} else {
		fmt.Println("👻 Demo mode: backend error logged to console only")
	}

	fmt.Printf("🔪 Jason is here... backend error captured: %+v\n", errorReq)

	response := models.ErrorResponse{
		Success:   true,
		Message:   "Backend error captured and stored in the basement",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Quote:     getRandomHorrorQuote(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) ServeTerrorsJS(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("static/terrors.js")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(content)
}

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

func (h *Handlers) Overlook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	commands := []models.Command{
		{
			Name:        "ListApps",
			Description: "Lister toutes les applications",
			Method:      "GET",
			Path:        "/api/apps",
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
		{
			Name:        "CreateApp",
			Description: "Créer une nouvelle application",
			Method:      "POST",
			Path:        "/api/apps",
			Params: map[string]interface{}{
				"name":        "string (required)",
				"description": "string (optional)",
				"origins":     "[]string (optional) - Liste des origins autorisées",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
					"Content-Type":  "application/json",
				},
				"body": map[string]interface{}{
					"name":        "Mon App",
					"description": "Description de mon app",
					"origins":     []string{"https://example.com", "https://app.example.com"},
				},
			},
		},
		{
			Name:        "GetApp",
			Description: "Récupérer une application par son app_id",
			Method:      "GET",
			Path:        "/api/apps/{appId}",
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
		{
			Name:        "UpdateApp",
			Description: "Mettre à jour une application (nom, description, origins, statut)",
			Method:      "PATCH",
			Path:        "/api/apps/{appId}",
			Params: map[string]interface{}{
				"name":        "string (optional)",
				"description": "string (optional)",
				"origins":     "[]string (optional)",
				"isActive":    "boolean (optional)",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
					"Content-Type":  "application/json",
				},
				"body": map[string]interface{}{
					"name":     "Nouveau nom",
					"isActive": false,
				},
			},
		},
		{
			Name:        "DeleteApp",
			Description: "Désactiver une application (soft delete)",
			Method:      "DELETE",
			Path:        "/api/apps/{appId}",
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
		{
			Name:        "CaptureFrontendError",
			Description: "Capturer une erreur frontend (JavaScript)",
			Method:      "POST",
			Path:        "/sadako",
			Params: map[string]interface{}{
				"appId":       "string (required)",
				"message":     "string (required)",
				"stack":       "string (optional)",
				"fingerprint": "string (required)",
				"url":         "string (optional)",
				"ts":          "int64 (optional)",
				"type":        "string (optional) - 'error', 'unhandledrejection'",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Content-Type": "application/json",
				},
				"body": map[string]interface{}{
					"appId":       "app_xxxxx",
					"message":     "Error message",
					"stack":       "Error: ...\n    at ...",
					"fingerprint": "a1b2c3d4...",
					"url":         "https://example.com/page",
					"ts":          1234567890,
					"type":        "error",
				},
			},
		},
		{
			Name:        "CaptureBackendError",
			Description: "Capturer une erreur backend (Go, Node.js, Python, etc.)",
			Method:      "POST",
			Path:        "/jason",
			Params: map[string]interface{}{
				"appId":       "string (required)",
				"message":     "string (required)",
				"stack":       "string (optional)",
				"fingerprint": "string (required)",
				"url":         "string (optional) - Endpoint/route où l'erreur s'est produite",
				"ts":          "int64 (optional)",
				"type":        "string (optional) - 'error', 'panic', 'exception'",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Content-Type": "application/json",
				},
				"body": map[string]interface{}{
					"appId":       "app_xxxxx",
					"message":     "Database connection failed",
					"stack":       "Error: connection timeout\n    at connect (db.js:42:11)",
					"fingerprint": "a1b2c3d4...",
					"url":         "/api/users",
					"ts":          1234567890,
					"type":        "error",
				},
			},
		},
		{
			Name:        "ListErrors",
			Description: "Lister les erreurs d'une application",
			Method:      "GET",
			Path:        "/api/errors?appId={appId}",
			Params: map[string]interface{}{
				"appId":  "string (required)",
				"status": "string (optional) - 'new', 'treated', 'deleted'",
				"limit":  "int (optional) - Nombre de résultats",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
		{
			Name:        "GetError",
			Description: "Récupérer une erreur par son ID",
			Method:      "GET",
			Path:        "/api/errors/{errorId}",
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
		{
			Name:        "GetErrorStats",
			Description: "Récupérer les statistiques d'erreurs d'une application",
			Method:      "GET",
			Path:        "/api/errors/stats?appId={appId}",
			Params: map[string]interface{}{
				"appId": "string (required)",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
		{
			Name:        "UpdateErrorStatus",
			Description: "Mettre à jour le statut d'une erreur",
			Method:      "PATCH",
			Path:        "/api/errors/{errorId}",
			Params: map[string]interface{}{
				"status": "string (required) - 'new', 'treated', 'deleted'",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
					"Content-Type":  "application/json",
				},
				"body": map[string]interface{}{
					"status": "treated",
				},
			},
		},
		{
			Name:        "ListWebhooks",
			Description: "Lister les webhooks d'une application",
			Method:      "GET",
			Path:        "/api/webhooks?appId={appId}",
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
		{
			Name:        "CreateWebhook",
			Description: "Créer un webhook (Discord ou GitHub). Pour GitHub, l'URL peut être omise si owner/repo sont dans config",
			Method:      "POST",
			Path:        "/api/webhooks?appId={appId}",
			Params: map[string]interface{}{
				"appId":  "string (required) - Dans query param ou body",
				"type":   "string (required) - 'discord' ou 'github'",
				"url":    "string (required pour Discord, optionnel pour GitHub)",
				"config": "object (optional) - Pour GitHub: token, owner, repo, labels",
			},
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
					"Content-Type":  "application/json",
				},
				"body_discord": map[string]interface{}{
					"type": "discord",
					"url":  "https://discord.com/api/webhooks/...",
				},
				"body_github": map[string]interface{}{
					"type": "github",
					"config": map[string]interface{}{
						"token":  "ghp_xxxxxxxxxxxx",
						"owner":  "username",
						"repo":   "repository-name",
						"labels": []string{"bug", "error"},
					},
				},
			},
		},
		{
			Name:        "DeleteWebhook",
			Description: "Supprimer un webhook",
			Method:      "DELETE",
			Path:        "/api/webhooks/{webhookId}",
			Example: map[string]interface{}{
				"headers": map[string]string{
					"Authorization": "Bearer <ADMIN_TOKEN>",
				},
			},
		},
	}

	response := models.DashboardResponse{
		Message:  "Welcome to the Overlook Hotel - Command Center",
		Quote:    getRandomHorrorQuote(),
		Commands: commands,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateApp crée une nouvelle application
func (h *Handlers) CreateApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var req models.AppRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	app, token, err := h.appService.CreateApp(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.AppResponse{
		Success:        true,
		Message:        "Application créée avec succès",
		App:            app,
		DashboardToken: token,
		Warning:        "⚠️ Conservez ce token en sécurité, il ne sera affiché qu'une seule fois !",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetApp récupère une application par son app_id
func (h *Handlers) GetApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	// Extraire app_id de l'URL (format: /api/apps/app_xxxxx)
	path := strings.TrimPrefix(r.URL.Path, "/api/apps/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "App ID required", http.StatusBadRequest)
		return
	}

	app, err := h.appService.GetAppByAppID(path)
	if err != nil {
		if strings.Contains(err.Error(), "non trouvée") {
			http.Error(w, "Application not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

// ListApps liste toutes les applications
func (h *Handlers) ListApps(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	apps, err := h.appService.ListApps()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"apps":    apps,
		"count":   len(apps),
	})
}

// UpdateApp met à jour une application
func (h *Handlers) UpdateApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" && r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/apps/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "App ID required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	if err := json.Unmarshal(body, &updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var name, description *string
	var origins *[]string
	var isActive *bool

	if n, ok := updateData["name"].(string); ok {
		name = &n
	}
	if d, ok := updateData["description"].(string); ok {
		description = &d
	}
	if o, ok := updateData["origins"].([]interface{}); ok {
		originsList := make([]string, len(o))
		for i, v := range o {
			if s, ok := v.(string); ok {
				originsList[i] = s
			}
		}
		origins = &originsList
	}
	if active, ok := updateData["isActive"].(bool); ok {
		isActive = &active
	}

	app, err := h.appService.UpdateApp(path, name, description, origins, isActive)
	if err != nil {
		if strings.Contains(err.Error(), "non trouvée") {
			http.Error(w, "Application not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Application mise à jour",
		"app":     app,
	})
}

// DeleteApp désactive une application (soft delete)
func (h *Handlers) DeleteApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/apps/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "App ID required", http.StatusBadRequest)
		return
	}

	isActive := false
	app, err := h.appService.UpdateApp(path, nil, nil, nil, &isActive)
	if err != nil {
		if strings.Contains(err.Error(), "non trouvée") {
			http.Error(w, "Application not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Application désactivée",
		"app":     app,
	})
}

// HandleApps route les requêtes vers les différents handlers d'apps
func (h *Handlers) HandleApps(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/apps")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		if r.Method == "GET" {
			h.ListApps(w, r)
		} else if r.Method == "POST" {
			h.CreateApp(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		if r.Method == "GET" {
			h.GetApp(w, r)
		} else if r.Method == "PATCH" || r.Method == "PUT" {
			h.UpdateApp(w, r)
		} else if r.Method == "DELETE" {
			h.DeleteApp(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// ========== HANDLERS ERREURS ==========

// ListErrors liste les erreurs d'une application
func (h *Handlers) ListErrors(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	appID := r.URL.Query().Get("appId")
	if appID == "" {
		http.Error(w, "appId parameter required", http.StatusBadRequest)
		return
	}

	var status *string
	if s := r.URL.Query().Get("status"); s != "" {
		status = &s
	}

	var limit *int
	if l := r.URL.Query().Get("limit"); l != "" {
		var parsed int
		if _, err := fmt.Sscanf(l, "%d", &parsed); err == nil && parsed > 0 {
			limit = &parsed
		}
	}

	errors, err := h.errorService.ListErrors(appID, status, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"errors":  errors,
		"count":   len(errors),
	})
}

func (h *Handlers) GetError(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/errors/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "Error ID required", http.StatusBadRequest)
		return
	}

	var errorID int64
	if _, err := fmt.Sscanf(path, "%d", &errorID); err != nil {
		http.Error(w, "Invalid error ID", http.StatusBadRequest)
		return
	}

	error, err := h.errorService.GetError(errorID)
	if err != nil {
		if strings.Contains(err.Error(), "non trouvée") {
			http.Error(w, "Error not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(error)
}

func (h *Handlers) UpdateErrorStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" && r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/errors/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "Error ID required", http.StatusBadRequest)
		return
	}

	var errorID int64
	if _, err := fmt.Sscanf(path, "%d", &errorID); err != nil {
		http.Error(w, "Invalid error ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	if err := json.Unmarshal(body, &updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	status, ok := updateData["status"].(string)
	if !ok || status == "" {
		http.Error(w, "status field required", http.StatusBadRequest)
		return
	}

	error, err := h.errorService.UpdateErrorStatus(errorID, status)
	if err != nil {
		if strings.Contains(err.Error(), "non trouvée") {
			http.Error(w, "Error not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Error status updated",
		"error":   error,
	})
}

func (h *Handlers) GetErrorStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	appID := r.URL.Query().Get("appId")
	if appID == "" {
		http.Error(w, "appId parameter required", http.StatusBadRequest)
		return
	}

	stats, err := h.errorService.GetErrorStats(appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"stats":   stats,
	})
}

func (h *Handlers) HandleErrors(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/errors")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		if r.Method == "GET" {
			h.ListErrors(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if path == "stats" {
		if r.Method == "GET" {
			h.GetErrorStats(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		if r.Method == "GET" {
			h.GetError(w, r)
		} else if r.Method == "PATCH" || r.Method == "PUT" {
			h.UpdateErrorStatus(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// ========== HANDLERS WEBHOOKS ==========

// CreateWebhook crée un nouveau webhook
func (h *Handlers) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var req models.WebhookRequest
	var bodyData map[string]interface{}
	if err := json.Unmarshal(body, &bodyData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	appID := r.URL.Query().Get("appId")
	if appID == "" {
		if appIDObj, ok := bodyData["appId"]; ok {
			if appIDStr, ok := appIDObj.(string); ok {
				appID = appIDStr
			}
		}
		if appID == "" {
			http.Error(w, "appId required (query param or body)", http.StatusBadRequest)
			return
		}
	}

	reqBytes, _ := json.Marshal(bodyData)
	if err := json.Unmarshal(reqBytes, &req); err != nil {
		http.Error(w, "Invalid webhook request format", http.StatusBadRequest)
		return
	}

	webhook, err := h.webhookService.CreateWebhook(appID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := models.WebhookResponse{
		Success: true,
		Message: "Webhook créé avec succès",
		Webhook: webhook,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	appID := r.URL.Query().Get("appId")
	if appID == "" {
		http.Error(w, "appId parameter required", http.StatusBadRequest)
		return
	}

	webhooks, err := h.webhookService.ListWebhooks(appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"webhooks": webhooks,
		"count":    len(webhooks),
	})
}

func (h *Handlers) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/webhooks/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "Webhook ID required", http.StatusBadRequest)
		return
	}

	var webhookID int64
	if _, err := fmt.Sscanf(path, "%d", &webhookID); err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	err := h.webhookService.DeleteWebhook(webhookID)
	if err != nil {
		if strings.Contains(err.Error(), "non trouvé") {
			http.Error(w, "Webhook not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Webhook supprimé",
	})
}

func (h *Handlers) HandleWebhooks(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/webhooks")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		if r.Method == "GET" {
			h.ListWebhooks(w, r)
		} else if r.Method == "POST" {
			h.CreateWebhook(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		if r.Method == "DELETE" {
			h.DeleteWebhook(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (h *Handlers) triggerWebhooks(appID string, error *models.Error) {
	if h.webhookService == nil {
		return
	}

	webhooks, err := h.webhookService.GetActiveWebhooks(appID)
	if err != nil {
		fmt.Printf("❌ Erreur récupération webhooks: %v\n", err)
		return
	}

	if len(webhooks) == 0 {
		return
	}

	for _, webhook := range webhooks {
		go h.sendWebhook(webhook, error)
	}
}

func (h *Handlers) sendWebhook(webhook *models.Webhook, error *models.Error) {
	switch webhook.Type {
	case "github":
		h.sendGitHubWebhook(webhook, error)
	case "discord":
		h.sendDiscordWebhook(webhook, error)
	default:
		fmt.Printf("⚠️ Type de webhook inconnu: %s\n", webhook.Type)
	}
}

func (h *Handlers) sendGitHubWebhook(webhook *models.Webhook, error *models.Error) {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(webhook.Config), &config); err != nil {
		fmt.Printf("❌ Erreur parsing config webhook GitHub: %v\n", err)
		return
	}

	token, _ := config["token"].(string)
	if token == "" {
		fmt.Printf("❌ Token GitHub manquant dans la config\n")
		return
	}

	owner, _ := config["owner"].(string)
	repo, _ := config["repo"].(string)

	if owner == "" || repo == "" {
		urlParts := strings.Split(webhook.URL, "/")
		if len(urlParts) >= 2 {
			if strings.Contains(webhook.URL, "api.github.com/repos/") {
				for i, part := range urlParts {
					if part == "repos" && i+2 < len(urlParts) {
						owner = urlParts[i+1]
						repo = strings.TrimSuffix(urlParts[i+2], "/issues")
						break
					}
				}
			} else if !strings.Contains(webhook.URL, "http") {
				parts := strings.Split(webhook.URL, "/")
				if len(parts) >= 2 {
					owner = parts[0]
					repo = parts[1]
				}
			}
		}
	}

	if owner == "" || repo == "" {
		fmt.Printf("❌ Owner/repo manquant pour webhook GitHub\n")
		return
	}

	// Construire l'URL de l'API GitHub
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)

	// Préparer les labels
	labels := []string{"[T]errors"}
	if labelsConfig, ok := config["labels"].([]interface{}); ok {
		for _, label := range labelsConfig {
			if labelStr, ok := label.(string); ok {
				labels = append(labels, labelStr)
			}
		}
	}
	labels = append(labels, error.Type)

	title := error.Message
	if len(title) > 100 {
		title = title[:97] + "..."
	}
	title = fmt.Sprintf("[Error] %s", title)

	body := fmt.Sprintf("## 🧟‍♀️ Erreur capturée par [T]errors\n\n"+
		"**Type** : %s\n"+
		"**Fingerprint** : `%s`\n"+
		"**URL** : %s\n"+
		"**Timestamp** : %s\n\n"+
		"### Message\n"+
		"```\n%s\n```\n\n"+
		"### Stack Trace\n"+
		"```\n%s\n```\n\n"+
		"### Détails\n"+
		"- **App ID** : %s\n"+
		"- **Error ID** : %d\n"+
		"- **Status** : %s\n",
		error.Type,
		error.Fingerprint,
		error.URL,
		error.CreatedAt.Format("2026-01-02 15:04:05 UTC"),
		error.Message,
		error.Stack,
		error.AppID,
		error.ID,
		error.Status,
	)

	payload := map[string]interface{}{
		"title":  title,
		"body":   body,
		"labels": labels,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("❌ Erreur sérialisation payload GitHub: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(payloadJSON)))
	if err != nil {
		fmt.Printf("❌ Erreur création requête GitHub: %v\n", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Erreur envoi webhook GitHub: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var issueResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&issueResp); err == nil {
			if issueURL, ok := issueResp["html_url"].(string); ok {
				fmt.Printf("✅ Issue GitHub créée: %s (erreur %d)\n", issueURL, error.ID)
			} else {
				fmt.Printf("✅ Issue GitHub créée pour erreur %d\n", error.ID)
			}
		} else {
			fmt.Printf("✅ Issue GitHub créée pour erreur %d\n", error.ID)
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		var errorResp map[string]interface{}
		json.Unmarshal(body, &errorResp)

		errorMsg := string(body)
		if msg, ok := errorResp["message"].(string); ok {
			errorMsg = msg
		}

		fmt.Printf("❌ Erreur création issue GitHub (status %d): %s\n", resp.StatusCode, errorMsg)

		if resp.StatusCode == 403 {
			fmt.Printf("💡 Aide: Pour un repository d'organisation, le token doit être autorisé pour l'organisation.\n")
			fmt.Printf("   → GitHub → Settings → Developer settings → Personal access tokens\n")
			fmt.Printf("   → Cochez l'organisation dans 'Organization access'\n")
			fmt.Printf("   → Ou: Organisation → Settings → Third-party access → Autoriser le token\n")
		} else if resp.StatusCode == 401 {
			fmt.Printf("💡 Aide: Token invalide ou expiré. Vérifiez votre token GitHub.\n")
		} else if resp.StatusCode == 404 {
			fmt.Printf("💡 Aide: Repository non trouvé. Vérifiez owner/repo: %s/%s\n", owner, repo)
		}
	}
}

func (h *Handlers) sendDiscordWebhook(webhook *models.Webhook, error *models.Error) {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(webhook.Config), &config); err != nil {
		fmt.Printf("❌ Erreur parsing config webhook Discord: %v\n", err)
		return
	}

	message := error.Message
	if len(message) > 500 {
		message = message[:497] + "..."
	}

	stack := error.Stack
	if len(stack) > 1000 {
		stack = stack[:997] + "..."
	}

	embed := map[string]interface{}{
		"title":       fmt.Sprintf("🎭 Erreur %s", error.Type),
		"description": message,
		"color":       15158332, // Rouge
		"fields": []map[string]interface{}{
			{
				"name":   "Fingerprint",
				"value":  fmt.Sprintf("`%s`", error.Fingerprint),
				"inline": true,
			},
			{
				"name":   "URL",
				"value":  error.URL,
				"inline": true,
			},
			{
				"name":   "App ID",
				"value":  error.AppID,
				"inline": true,
			},
			{
				"name":  "Stack Trace",
				"value": fmt.Sprintf("```\n%s\n```", stack),
			},
		},
		"timestamp": error.CreatedAt.Format(time.RFC3339),
		"footer": map[string]interface{}{
			"text": fmt.Sprintf("Error ID: %d", error.ID),
		},
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{embed},
	}

	if username, ok := config["username"].(string); ok && username != "" {
		payload["username"] = username
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("❌ Erreur sérialisation payload Discord: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhook.URL, "application/json", strings.NewReader(string(payloadJSON)))
	if err != nil {
		fmt.Printf("❌ Erreur envoi webhook Discord: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("✅ Notification Discord envoyée pour erreur %d\n", error.ID)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Erreur envoi Discord (status %d): %s\n", resp.StatusCode, string(body))
	}
}
