package services

import (
	"database/sql"
	"os"
	"terrors/internal/models"
	"testing"

	_ "github.com/lib/pq"
)

// getTestDB crée une connexion de test à la base de données
// Utilise PG_URL_TEST si disponible, sinon PG_URL
func getTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("PG_URL_TEST")
	if dbURL == "" {
		dbURL = os.Getenv("PG_URL")
	}
	if dbURL == "" {
		t.Skip("PG_URL ou PG_URL_TEST non défini, skip des tests")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("Erreur connexion DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Erreur ping DB: %v", err)
	}

	return db
}

// cleanupTestData nettoie les données de test
func cleanupTestData(t *testing.T, db *sql.DB, appID string) {
	_, err := db.Exec("DELETE FROM apps WHERE app_id = $1", appID)
	if err != nil {
		t.Logf("Erreur nettoyage: %v", err)
	}
}

func TestAppService_CreateApp(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	service := NewAppService(db)

	tests := []struct {
		name    string
		req     models.AppRequest
		wantErr bool
	}{
		{
			name: "Création app valide",
			req: models.AppRequest{
				Name:        "Test App",
				Description: "Description de test",
				Origins:     []string{"https://example.com"},
			},
			wantErr: false,
		},
		{
			name: "Création app sans nom",
			req: models.AppRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "Création app avec plusieurs origins",
			req: models.AppRequest{
				Name:    "Multi Origin App",
				Origins: []string{"https://example.com", "https://app.example.com"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, token, err := service.CreateApp(tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Vérifications sur l'app créée
				if app == nil {
					t.Fatal("CreateApp() a retourné nil")
				}
				if app.Name != tt.req.Name {
					t.Errorf("CreateApp() app.Name = %v, want %v", app.Name, tt.req.Name)
				}
				if app.AppID == "" {
					t.Error("CreateApp() app.AppID est vide")
				}
				if token == "" {
					t.Error("CreateApp() token est vide")
				}
				if len(token) < 32 {
					t.Errorf("CreateApp() token trop court: %d caractères", len(token))
				}

				// Nettoyer après le test
				defer cleanupTestData(t, db, app.AppID)
			}
		})
	}
}

func TestAppService_GetAppByAppID(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	service := NewAppService(db)

	// Créer une app de test
	app, _, err := service.CreateApp(models.AppRequest{
		Name:        "Test Get App",
		Description: "Pour tester GetAppByAppID",
	})
	if err != nil {
		t.Fatalf("Erreur création app de test: %v", err)
	}
	defer cleanupTestData(t, db, app.AppID)

	// Test récupération
	got, err := service.GetAppByAppID(app.AppID)
	if err != nil {
		t.Fatalf("GetAppByAppID() error = %v", err)
	}

	if got.AppID != app.AppID {
		t.Errorf("GetAppByAppID() appID = %v, want %v", got.AppID, app.AppID)
	}
	if got.Name != app.Name {
		t.Errorf("GetAppByAppID() name = %v, want %v", got.Name, app.Name)
	}

	// Test app inexistante
	_, err = service.GetAppByAppID("app_inexistant")
	if err == nil {
		t.Error("GetAppByAppID() devrait retourner une erreur pour app inexistante")
	}
}

func TestAppService_ValidateApp(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	service := NewAppService(db)

	// Créer une app active
	app, _, err := service.CreateApp(models.AppRequest{
		Name: "Test Validate",
	})
	if err != nil {
		t.Fatalf("Erreur création app: %v", err)
	}
	defer cleanupTestData(t, db, app.AppID)

	// Test app active
	err = service.ValidateApp(app.AppID)
	if err != nil {
		t.Errorf("ValidateApp() error = %v, devrait être valide", err)
	}

	// Désactiver l'app
	_, err = service.UpdateApp(app.AppID, nil, nil, nil, boolPtr(false))
	if err != nil {
		t.Fatalf("Erreur désactivation app: %v", err)
	}

	// Test app désactivée
	err = service.ValidateApp(app.AppID)
	if err == nil {
		t.Error("ValidateApp() devrait retourner une erreur pour app désactivée")
	}
}

func TestAppService_ValidateOrigin(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	service := NewAppService(db)

	// Créer une app avec origins
	app, _, err := service.CreateApp(models.AppRequest{
		Name:    "Test Origins",
		Origins: []string{"https://example.com", "https://app.example.com"},
	})
	if err != nil {
		t.Fatalf("Erreur création app: %v", err)
	}
	defer cleanupTestData(t, db, app.AppID)

	tests := []struct {
		name    string
		origin  string
		wantErr bool
	}{
		{"Origin autorisée 1", "https://example.com", false},
		{"Origin autorisée 2", "https://app.example.com", false},
		{"Origin non autorisée", "https://evil.com", true},
		{"Origin vide (acceptée si pas de restrictions)", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateOrigin(app.AppID, tt.origin)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrigin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
}
