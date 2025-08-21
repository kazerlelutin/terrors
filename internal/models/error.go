package models

import (
	"time"
)

// Error représente une erreur capturée
type Error struct {
	ID          int64     `json:"id" db:"id"`
	AppID       string    `json:"appId" db:"app_id"`
	Message     string    `json:"message" db:"message"`
	Stack       string    `json:"stack" db:"stack"`
	Fingerprint string    `json:"fingerprint" db:"fingerprint"`
	URL         string    `json:"url" db:"url"`
	Type        string    `json:"type" db:"type"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// ErrorRequest représente la requête reçue du client
type ErrorRequest struct {
	AppID       string `json:"appId"`
	Message     string `json:"message"`
	Stack       string `json:"stack"`
	Fingerprint string `json:"fingerprint"`
	URL         string `json:"url"`
	Timestamp   int64  `json:"ts"`
	Type        string `json:"type"`
}

// ErrorResponse représente la réponse de l'API
type ErrorResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Quote     string `json:"quote,omitempty"`
}
