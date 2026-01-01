package models

import (
	"time"
)

type App struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	AppID       string    `json:"appId" db:"app_id"`
	TokenHash   string    `json:"-" db:"token_hash"` // Jamais exposé en JSON
	Description string    `json:"description" db:"description"`
	Origins     string    `json:"origins" db:"origins"` // Liste d'origins séparées par des virgules
	IsActive    bool      `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type AppRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Origins     []string `json:"origins"` // Liste d'origins autorisées
}

type AppResponse struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	App            *App   `json:"app,omitempty"`
	DashboardToken string `json:"dashboardToken,omitempty"` // Displayed once
	Warning        string `json:"warning,omitempty"`
}
