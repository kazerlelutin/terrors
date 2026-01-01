package models

import (
	"time"
)

type Webhook struct {
	ID        int64     `json:"id" db:"id"`
	AppID     string    `json:"appId" db:"app_id"`
	Type      string    `json:"type" db:"type"` // 'discord', 'github'
	URL       string    `json:"url" db:"url"`
	Config    string    `json:"config" db:"config"` // JSON string
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type WebhookRequest struct {
	Type   string                 `json:"type"` // 'discord', 'github'
	URL    string                 `json:"url"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type WebhookResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Webhook *Webhook `json:"webhook,omitempty"`
}
