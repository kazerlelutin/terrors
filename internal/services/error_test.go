package services

import (
	"terrors/internal/models"
	"testing"
)

func TestErrorService_SaveError(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	appService := NewAppService(db)
	errorService := NewErrorService(db)

	// Créer une app de test
	app, _, err := appService.CreateApp(models.AppRequest{
		Name: "Test Error App",
	})
	if err != nil {
		t.Fatalf("Erreur création app: %v", err)
	}
	defer cleanupTestData(t, db, app.AppID)

	tests := []struct {
		name    string
		req     models.ErrorRequest
		wantErr bool
	}{
		{
			name: "Erreur valide",
			req: models.ErrorRequest{
				AppID:       app.AppID,
				Message:     "Test error message",
				Stack:       "Error: test\n    at test.js:10:5",
				Fingerprint: "a1b2c3d4e5f6",
				URL:         "https://example.com/page",
				Type:        "error",
			},
			wantErr: false,
		},
		{
			name: "Erreur sans stack",
			req: models.ErrorRequest{
				AppID:       app.AppID,
				Message:     "Simple error",
				Fingerprint: "b2c3d4e5f6a1",
			},
			wantErr: false,
		},
		{
			name: "Erreur avec app inexistante",
			req: models.ErrorRequest{
				AppID:       "app_inexistant",
				Message:     "Test",
				Fingerprint: "c3d4e5f6a1b2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savedError, err := errorService.SaveError(tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if savedError == nil {
					t.Fatal("SaveError() a retourné nil")
				}
				if savedError.Message != tt.req.Message {
					t.Errorf("SaveError() savedError.Message = %v, want %v", savedError.Message, tt.req.Message)
				}
				if savedError.AppID != tt.req.AppID {
					t.Errorf("SaveError() savedError.AppID = %v, want %v", savedError.AppID, tt.req.AppID)
				}
				if savedError.Status != "new" {
					t.Errorf("SaveError() savedError.Status = %v, want 'new'", savedError.Status)
				}

				// Nettoyer l'erreur créée
				defer func() {
					db.Exec("DELETE FROM errors WHERE id = $1", savedError.ID)
				}()
			}
		})
	}
}

func TestErrorService_ListErrors(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	appService := NewAppService(db)
	errorService := NewErrorService(db)

	// Créer une app de test
	app, _, err := appService.CreateApp(models.AppRequest{
		Name: "Test List Errors",
	})
	if err != nil {
		t.Fatalf("Erreur création app: %v", err)
	}
	defer cleanupTestData(t, db, app.AppID)

	// Créer quelques erreurs de test
	errors := []models.ErrorRequest{
		{
			AppID:       app.AppID,
			Message:     "Error 1",
			Fingerprint: "fingerprint1",
			Type:        "error",
		},
		{
			AppID:       app.AppID,
			Message:     "Error 2",
			Fingerprint: "fingerprint2",
			Type:        "backend",
		},
	}

	var createdErrors []*models.Error
	for _, errReq := range errors {
		savedError, err := errorService.SaveError(errReq)
		if err != nil {
			t.Fatalf("Erreur création erreur de test: %v", err)
		}
		createdErrors = append(createdErrors, savedError)
		defer db.Exec("DELETE FROM errors WHERE id = $1", savedError.ID)
	}

	// Test liste complète
	allErrors, err := errorService.ListErrors(app.AppID, nil, nil)
	if err != nil {
		t.Fatalf("ListErrors() error = %v", err)
	}

	if len(allErrors) < len(createdErrors) {
		t.Errorf("ListErrors() a retourné %d erreurs, attendu au moins %d", len(allErrors), len(createdErrors))
	}

	// Test avec limite
	limitedErrors, err := errorService.ListErrors(app.AppID, nil, intPtr(1))
	if err != nil {
		t.Fatalf("ListErrors() avec limite error = %v", err)
	}

	if len(limitedErrors) > 1 {
		t.Errorf("ListErrors() avec limite a retourné %d erreurs, attendu max 1", len(limitedErrors))
	}
}

func TestErrorService_UpdateErrorStatus(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	appService := NewAppService(db)
	errorService := NewErrorService(db)

	// Créer une app et une erreur
	app, _, err := appService.CreateApp(models.AppRequest{
		Name: "Test Update Status",
	})
	if err != nil {
		t.Fatalf("Erreur création app: %v", err)
	}
	defer cleanupTestData(t, db, app.AppID)

	savedError, err := errorService.SaveError(models.ErrorRequest{
		AppID:       app.AppID,
		Message:     "Test error",
		Fingerprint: "test123",
	})
	if err != nil {
		t.Fatalf("Erreur création erreur: %v", err)
	}
	defer db.Exec("DELETE FROM errors WHERE id = $1", savedError.ID)

	// Test mise à jour vers "treated"
	updated, err := errorService.UpdateErrorStatus(savedError.ID, "treated")
	if err != nil {
		t.Fatalf("UpdateErrorStatus() error = %v", err)
	}

	if updated.Status != "treated" {
		t.Errorf("UpdateErrorStatus() status = %v, want 'treated'", updated.Status)
	}

	// Test statut invalide
	_, err = errorService.UpdateErrorStatus(savedError.ID, "invalid")
	if err == nil {
		t.Error("UpdateErrorStatus() devrait retourner une erreur pour statut invalide")
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
