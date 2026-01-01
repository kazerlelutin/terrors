package services

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"terrors/internal/models"
	"time"
)

type AppService struct {
	db *sql.DB
}

func NewAppService(db *sql.DB) *AppService {
	return &AppService{db: db}
}

func (s *AppService) generateAppID() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	for i := 0; i < 10; i++ {
		b := make([]byte, length)
		for i := range b {
			randomByte := make([]byte, 1)
			if _, err := rand.Read(randomByte); err != nil {
				return "", fmt.Errorf("erreur génération aléatoire: %w", err)
			}
			b[i] = charset[randomByte[0]%byte(len(charset))]
		}
		appID := "app_" + string(b)

		var exists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM apps WHERE app_id = $1)", appID).Scan(&exists)
		if err != nil {
			return "", fmt.Errorf("erreur vérification unicité: %w", err)
		}
		if !exists {
			return appID, nil
		}
	}

	return "", fmt.Errorf("impossible de générer un app_id unique après 10 tentatives")
}

func (s *AppService) generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("erreur génération token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (s *AppService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AppService) CreateApp(req models.AppRequest) (*models.App, string, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, "", fmt.Errorf("le nom de l'application est requis")
	}

	appID, err := s.generateAppID()
	if err != nil {
		return nil, "", fmt.Errorf("erreur génération app_id: %w", err)
	}

	token, err := s.generateToken()
	if err != nil {
		return nil, "", fmt.Errorf("erreur génération token: %w", err)
	}

	tokenHash := s.hashToken(token)

	originsStr := ""
	if len(req.Origins) > 0 {
		originsStr = strings.Join(req.Origins, ",")
	}

	now := time.Now()
	var id int64
	err = s.db.QueryRow(`
		INSERT INTO apps (name, app_id, token_hash, description, origins, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, req.Name, appID, tokenHash, req.Description, originsStr, true, now, now).Scan(&id)

	if err != nil {
		return nil, "", fmt.Errorf("erreur insertion en base: %w", err)
	}

	app := &models.App{
		ID:          id,
		Name:        req.Name,
		AppID:       appID,
		TokenHash:   tokenHash,
		Description: req.Description,
		Origins:     originsStr,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return app, token, nil
}

// GetAppByAppID récupère une application par son app_id
func (s *AppService) GetAppByAppID(appID string) (*models.App, error) {
	app := &models.App{}
	err := s.db.QueryRow(`
		SELECT id, name, app_id, token_hash, description, origins, is_active, created_at, updated_at
		FROM apps
		WHERE app_id = $1
	`, appID).Scan(
		&app.ID,
		&app.Name,
		&app.AppID,
		&app.TokenHash,
		&app.Description,
		&app.Origins,
		&app.IsActive,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("application non trouvée")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur récupération app: %w", err)
	}

	return app, nil
}

// GetAppByID récupère une application par son ID numérique
func (s *AppService) GetAppByID(id int64) (*models.App, error) {
	app := &models.App{}
	err := s.db.QueryRow(`
		SELECT id, name, app_id, token_hash, description, origins, is_active, created_at, updated_at
		FROM apps
		WHERE id = $1
	`, id).Scan(
		&app.ID,
		&app.Name,
		&app.AppID,
		&app.TokenHash,
		&app.Description,
		&app.Origins,
		&app.IsActive,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("application non trouvée")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur récupération app: %w", err)
	}

	return app, nil
}

// ListApps liste toutes les applications
func (s *AppService) ListApps() ([]*models.App, error) {
	rows, err := s.db.Query(`
		SELECT id, name, app_id, token_hash, description, origins, is_active, created_at, updated_at
		FROM apps
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération apps: %w", err)
	}
	defer rows.Close()

	var apps []*models.App
	for rows.Next() {
		app := &models.App{}
		err := rows.Scan(
			&app.ID,
			&app.Name,
			&app.AppID,
			&app.TokenHash,
			&app.Description,
			&app.Origins,
			&app.IsActive,
			&app.CreatedAt,
			&app.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur scan app: %w", err)
		}
		apps = append(apps, app)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur parcours rows: %w", err)
	}

	return apps, nil
}

// UpdateApp met à jour une application
func (s *AppService) UpdateApp(appID string, name, description *string, origins *[]string, isActive *bool) (*models.App, error) {
	// Construire la requête dynamiquement
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if name != nil {
		trimmed := strings.TrimSpace(*name)
		if trimmed == "" {
			return nil, fmt.Errorf("le nom ne peut pas être vide")
		}
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, trimmed)
		argPos++
	}

	if description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *description)
		argPos++
	}

	if origins != nil {
		originsStr := strings.Join(*origins, ",")
		updates = append(updates, fmt.Sprintf("origins = $%d", argPos))
		args = append(args, originsStr)
		argPos++
	}

	if isActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *isActive)
		argPos++
	}

	if len(updates) == 0 {
		return s.GetAppByAppID(appID)
	}

	// Toujours mettre à jour updated_at
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	// Ajouter app_id à la fin pour le WHERE
	args = append(args, appID)

	query := fmt.Sprintf(`
		UPDATE apps
		SET %s
		WHERE app_id = $%d
		RETURNING id, name, app_id, token_hash, description, origins, is_active, created_at, updated_at
	`, strings.Join(updates, ", "), argPos)

	app := &models.App{}
	err := s.db.QueryRow(query, args...).Scan(
		&app.ID,
		&app.Name,
		&app.AppID,
		&app.TokenHash,
		&app.Description,
		&app.Origins,
		&app.IsActive,
		&app.CreatedAt,
		&app.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("application non trouvée")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur mise à jour app: %w", err)
	}

	return app, nil
}

// ValidateApp vérifie qu'une app existe et est active
func (s *AppService) ValidateApp(appID string) error {
	var isActive bool
	err := s.db.QueryRow(`
		SELECT is_active
		FROM apps
		WHERE app_id = $1
	`, appID).Scan(&isActive)

	if err == sql.ErrNoRows {
		return fmt.Errorf("application non trouvée")
	}
	if err != nil {
		return fmt.Errorf("erreur validation app: %w", err)
	}

	if !isActive {
		return fmt.Errorf("application désactivée")
	}

	return nil
}

// ValidateOrigin vérifie qu'une origin est autorisée pour une app
func (s *AppService) ValidateOrigin(appID, origin string) error {
	var originsStr string
	var isActive bool
	err := s.db.QueryRow(`
		SELECT origins, is_active
		FROM apps
		WHERE app_id = $1
	`, appID).Scan(&originsStr, &isActive)

	if err == sql.ErrNoRows {
		return fmt.Errorf("application non trouvée")
	}
	if err != nil {
		return fmt.Errorf("erreur validation origin: %w", err)
	}

	if !isActive {
		return fmt.Errorf("application désactivée")
	}

	// Si pas d'origins définies, on accepte tout (pour compatibilité)
	if originsStr == "" {
		return nil
	}

	// Vérifier si l'origin est dans la liste
	origins := strings.Split(originsStr, ",")
	for _, allowedOrigin := range origins {
		allowedOrigin = strings.TrimSpace(allowedOrigin)
		if allowedOrigin == origin || allowedOrigin == "*" {
			return nil
		}
	}

	return fmt.Errorf("origin non autorisée: %s", origin)
}
