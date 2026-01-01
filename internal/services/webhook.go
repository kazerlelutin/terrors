package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"terrors/internal/models"
	"time"
)

type WebhookService struct {
	db *sql.DB
}

func NewWebhookService(db *sql.DB) *WebhookService {
	return &WebhookService{db: db}
}

func (s *WebhookService) CreateWebhook(appID string, req models.WebhookRequest) (*models.Webhook, error) {
	if req.Type != "discord" && req.Type != "github" {
		return nil, fmt.Errorf("type de webhook invalide: %s (doit être 'discord' ou 'github')", req.Type)
	}

	if req.URL == "" {
		if req.Type == "github" {
			owner, _ := req.Config["owner"].(string)
			repo, _ := req.Config["repo"].(string)
			if owner != "" && repo != "" {
				req.URL = fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
			} else {
				return nil, fmt.Errorf("pour un webhook GitHub, fournissez soit 'url', soit 'owner' et 'repo' dans config")
			}
		} else {
			return nil, fmt.Errorf("l'URL du webhook est requise")
		}
	}

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM apps WHERE app_id = $1)", appID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("erreur vérification app: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("application non trouvée")
	}

	configJSON := "{}"
	if len(req.Config) > 0 {
		configBytes, err := json.Marshal(req.Config)
		if err != nil {
			return nil, fmt.Errorf("erreur sérialisation config: %w", err)
		}
		configJSON = string(configBytes)
	}

	now := time.Now()
	var id int64
	err = s.db.QueryRow(`
		INSERT INTO webhooks (app_id, type, url, config, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, appID, req.Type, req.URL, configJSON, true, now, now).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("erreur insertion webhook: %w", err)
	}

	webhook := &models.Webhook{
		ID:        id,
		AppID:     appID,
		Type:      req.Type,
		URL:       req.URL,
		Config:    configJSON,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return webhook, nil
}

func (s *WebhookService) GetWebhook(id int64) (*models.Webhook, error) {
	webhook := &models.Webhook{}
	err := s.db.QueryRow(`
		SELECT id, app_id, type, url, config, is_active, created_at, updated_at
		FROM webhooks
		WHERE id = $1
	`, id).Scan(
		&webhook.ID,
		&webhook.AppID,
		&webhook.Type,
		&webhook.URL,
		&webhook.Config,
		&webhook.IsActive,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("webhook non trouvé")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur récupération webhook: %w", err)
	}

	return webhook, nil
}

func (s *WebhookService) ListWebhooks(appID string) ([]*models.Webhook, error) {
	rows, err := s.db.Query(`
		SELECT id, app_id, type, url, config, is_active, created_at, updated_at
		FROM webhooks
		WHERE app_id = $1
		ORDER BY created_at DESC
	`, appID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération webhooks: %w", err)
	}
	defer rows.Close()

	var webhooks []*models.Webhook
	for rows.Next() {
		webhook := &models.Webhook{}
		err := rows.Scan(
			&webhook.ID,
			&webhook.AppID,
			&webhook.Type,
			&webhook.URL,
			&webhook.Config,
			&webhook.IsActive,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur scan webhook: %w", err)
		}
		webhooks = append(webhooks, webhook)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur parcours rows: %w", err)
	}

	return webhooks, nil
}

func (s *WebhookService) UpdateWebhook(id int64, url *string, config *map[string]interface{}, isActive *bool) (*models.Webhook, error) {
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if url != nil {
		if *url == "" {
			return nil, fmt.Errorf("l'URL ne peut pas être vide")
		}
		updates = append(updates, fmt.Sprintf("url = $%d", argPos))
		args = append(args, *url)
		argPos++
	}

	if config != nil {
		configBytes, err := json.Marshal(*config)
		if err != nil {
			return nil, fmt.Errorf("erreur sérialisation config: %w", err)
		}
		updates = append(updates, fmt.Sprintf("config = $%d", argPos))
		args = append(args, string(configBytes))
		argPos++
	}

	if isActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *isActive)
		argPos++
	}

	if len(updates) == 0 {
		return s.GetWebhook(id)
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	args = append(args, id)

	setClause := ""
	for i, update := range updates {
		if i > 0 {
			setClause += ", "
		}
		setClause += update
	}

	query := fmt.Sprintf(`
		UPDATE webhooks
		SET %s
		WHERE id = $%d
		RETURNING id, app_id, type, url, config, is_active, created_at, updated_at
	`, setClause, argPos)

	webhook := &models.Webhook{}
	err := s.db.QueryRow(query, args...).Scan(
		&webhook.ID,
		&webhook.AppID,
		&webhook.Type,
		&webhook.URL,
		&webhook.Config,
		&webhook.IsActive,
		&webhook.CreatedAt,
		&webhook.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("webhook non trouvé")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur mise à jour webhook: %w", err)
	}

	return webhook, nil
}

func (s *WebhookService) DeleteWebhook(id int64) error {
	result, err := s.db.Exec("DELETE FROM webhooks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("erreur suppression webhook: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur vérification suppression: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("webhook non trouvé")
	}

	return nil
}

func (s *WebhookService) GetActiveWebhooks(appID string) ([]*models.Webhook, error) {
	rows, err := s.db.Query(`
		SELECT id, app_id, type, url, config, is_active, created_at, updated_at
		FROM webhooks
		WHERE app_id = $1 AND is_active = TRUE
		ORDER BY created_at DESC
	`, appID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération webhooks actifs: %w", err)
	}
	defer rows.Close()

	var webhooks []*models.Webhook
	for rows.Next() {
		webhook := &models.Webhook{}
		err := rows.Scan(
			&webhook.ID,
			&webhook.AppID,
			&webhook.Type,
			&webhook.URL,
			&webhook.Config,
			&webhook.IsActive,
			&webhook.CreatedAt,
			&webhook.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur scan webhook: %w", err)
		}
		webhooks = append(webhooks, webhook)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur parcours rows: %w", err)
	}

	return webhooks, nil
}
