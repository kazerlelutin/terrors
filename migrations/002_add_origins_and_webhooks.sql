-- Migration 002: Ajout des origins et structure pour webhooks

-- Ajouter le champ origins (liste d'origins autorisées, séparées par des virgules)
ALTER TABLE apps ADD COLUMN IF NOT EXISTS origins TEXT DEFAULT '';

-- Table pour les webhooks
CREATE TABLE IF NOT EXISTS webhooks (
    id SERIAL PRIMARY KEY,
    app_id VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'discord', 'github'
    url TEXT NOT NULL,
    config JSONB, -- Configuration spécifique (channel Discord, repo GitHub, etc.)
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (app_id) REFERENCES apps(app_id) ON DELETE CASCADE
);

-- Index pour les webhooks
CREATE INDEX IF NOT EXISTS idx_webhooks_app_id ON webhooks(app_id);
CREATE INDEX IF NOT EXISTS idx_webhooks_type ON webhooks(type);
CREATE INDEX IF NOT EXISTS idx_webhooks_active ON webhooks(is_active);

-- Contraintes
ALTER TABLE webhooks ADD CONSTRAINT webhooks_type_valid CHECK (type IN ('discord', 'github'));
ALTER TABLE webhooks ADD CONSTRAINT webhooks_url_not_empty CHECK (url != '');

-- Trigger pour updated_at sur webhooks
CREATE TRIGGER update_webhooks_updated_at 
    BEFORE UPDATE ON webhooks 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

