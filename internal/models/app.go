package models

import (
	"time"
)

// App représente une application enregistrée
type App struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	AppID       string    `json:"appId" db:"app_id"`
	TokenHash   string    `json:"-" db:"token_hash"` // Jamais exposé en JSON
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// AppRequest représente la requête pour créer une app
type AppRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AppResponse représente la réponse de l'API pour les apps
type AppResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	App            *App   `json:"app,omitempty"`
	DashboardToken string `json:"dashboardToken,omitempty"` // Affiché une seule fois
	Warning        string `json:"warning,omitempty"`
}
