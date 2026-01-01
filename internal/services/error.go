package services

import (
	"database/sql"
	"fmt"
	"terrors/internal/models"
	"time"
)

type ErrorService struct {
	db *sql.DB
}

func NewErrorService(db *sql.DB) *ErrorService {
	return &ErrorService{db: db}
}

func (s *ErrorService) SaveError(req models.ErrorRequest) (*models.Error, error) {
	var isActive bool
	err := s.db.QueryRow(`
		SELECT is_active
		FROM apps
		WHERE app_id = $1
	`, req.AppID).Scan(&isActive)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("application non trouvée")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur vérification app: %w", err)
	}
	if !isActive {
		return nil, fmt.Errorf("application désactivée")
	}

	var existingID int64
	var existingCreatedAt time.Time
	err = s.db.QueryRow(`
		SELECT id, created_at
		FROM errors
		WHERE app_id = $1 
		  AND fingerprint = $2 
		  AND status != 'deleted'
		  AND created_at > NOW() - INTERVAL '5 minutes'
		ORDER BY created_at DESC
		LIMIT 1
	`, req.AppID, req.Fingerprint).Scan(&existingID, &existingCreatedAt)

	now := time.Now()

	if err == nil {
		existingError, err := s.GetError(existingID)
		if err == nil {
			s.db.Exec(`
				UPDATE errors 
				SET updated_at = $1 
				WHERE id = $2
			`, now, existingID)
			return existingError, nil
		}
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("erreur vérification doublon: %w", err)
	}
	var id int64
	err = s.db.QueryRow(`
		INSERT INTO errors (app_id, message, stack, fingerprint, url, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`, req.AppID, req.Message, req.Stack, req.Fingerprint, req.URL, req.Type, "new", now, now).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("erreur insertion erreur: %w", err)
	}

	error := &models.Error{
		ID:          id,
		AppID:       req.AppID,
		Message:     req.Message,
		Stack:       req.Stack,
		Fingerprint: req.Fingerprint,
		URL:         req.URL,
		Type:        req.Type,
		Status:      "new",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return error, nil
}

func (s *ErrorService) GetError(id int64) (*models.Error, error) {
	error := &models.Error{}
	err := s.db.QueryRow(`
		SELECT id, app_id, message, stack, fingerprint, url, type, status, created_at, updated_at
		FROM errors
		WHERE id = $1
	`, id).Scan(
		&error.ID,
		&error.AppID,
		&error.Message,
		&error.Stack,
		&error.Fingerprint,
		&error.URL,
		&error.Type,
		&error.Status,
		&error.CreatedAt,
		&error.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("erreur non trouvée")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur récupération: %w", err)
	}

	return error, nil
}

func (s *ErrorService) ListErrors(appID string, status *string, limit *int) ([]*models.Error, error) {
	query := `
		SELECT id, app_id, message, stack, fingerprint, url, type, status, created_at, updated_at
		FROM errors
		WHERE app_id = $1
	`
	args := []interface{}{appID}
	argPos := 2

	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *status)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	if limit != nil && *limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, *limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération erreurs: %w", err)
	}
	defer rows.Close()

	var errors []*models.Error
	for rows.Next() {
		error := &models.Error{}
		err := rows.Scan(
			&error.ID,
			&error.AppID,
			&error.Message,
			&error.Stack,
			&error.Fingerprint,
			&error.URL,
			&error.Type,
			&error.Status,
			&error.CreatedAt,
			&error.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur scan erreur: %w", err)
		}
		errors = append(errors, error)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erreur parcours rows: %w", err)
	}

	return errors, nil
}

func (s *ErrorService) UpdateErrorStatus(id int64, status string) (*models.Error, error) {
	if status != "new" && status != "treated" && status != "deleted" {
		return nil, fmt.Errorf("statut invalide: %s (doit être 'new', 'treated' ou 'deleted')", status)
	}

	error := &models.Error{}
	err := s.db.QueryRow(`
		UPDATE errors
		SET status = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, app_id, message, stack, fingerprint, url, type, status, created_at, updated_at
	`, status, time.Now(), id).Scan(
		&error.ID,
		&error.AppID,
		&error.Message,
		&error.Stack,
		&error.Fingerprint,
		&error.URL,
		&error.Type,
		&error.Status,
		&error.CreatedAt,
		&error.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("erreur non trouvée")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur mise à jour: %w", err)
	}

	return error, nil
}

func (s *ErrorService) GetErrorStats(appID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var total int
	err := s.db.QueryRow(`
		SELECT COUNT(*)
		FROM errors
		WHERE app_id = $1
	`, appID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération total: %w", err)
	}
	stats["total"] = total

	statusCounts := make(map[string]int)
	rows, err := s.db.Query(`
		SELECT status, COUNT(*)
		FROM errors
		WHERE app_id = $1
		GROUP BY status
	`, appID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération stats par statut: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("erreur scan statut: %w", err)
		}
		statusCounts[status] = count
	}
	stats["byStatus"] = statusCounts

	typeCounts := make(map[string]int)
	rows, err = s.db.Query(`
		SELECT type, COUNT(*)
		FROM errors
		WHERE app_id = $1
		GROUP BY type
	`, appID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération stats par type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var errorType string
		var count int
		if err := rows.Scan(&errorType, &count); err != nil {
			return nil, fmt.Errorf("erreur scan type: %w", err)
		}
		typeCounts[errorType] = count
	}
	stats["byType"] = typeCounts

	return stats, nil
}
