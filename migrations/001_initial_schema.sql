-- Migration initiale pour Terrors
-- Création des tables de base

-- Table des applications
CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    app_id VARCHAR(100) NOT NULL UNIQUE,
    token_hash VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table des erreurs
CREATE TABLE IF NOT EXISTS errors (
    id SERIAL PRIMARY KEY,
    app_id VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    stack TEXT,
    fingerprint VARCHAR(64) NOT NULL,
    url TEXT,
    type VARCHAR(50) DEFAULT 'error',
    status VARCHAR(20) DEFAULT 'new', -- 'new', 'treated', 'deleted'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index pour optimiser les requêtes
CREATE INDEX IF NOT EXISTS idx_apps_app_id ON apps(app_id);
CREATE INDEX IF NOT EXISTS idx_apps_token_hash ON apps(token_hash);
CREATE INDEX IF NOT EXISTS idx_apps_active ON apps(is_active);

CREATE INDEX IF NOT EXISTS idx_errors_app_id ON errors(app_id);
CREATE INDEX IF NOT EXISTS idx_errors_fingerprint ON errors(fingerprint);
CREATE INDEX IF NOT EXISTS idx_errors_created_at ON errors(created_at);
CREATE INDEX IF NOT EXISTS idx_errors_status ON errors(status);
CREATE INDEX IF NOT EXISTS idx_errors_type ON errors(type);

-- Contraintes
ALTER TABLE apps ADD CONSTRAINT apps_name_not_empty CHECK (name != '');
ALTER TABLE apps ADD CONSTRAINT apps_app_id_not_empty CHECK (app_id != '');
ALTER TABLE apps ADD CONSTRAINT apps_token_hash_not_empty CHECK (token_hash != '');

ALTER TABLE errors ADD CONSTRAINT errors_app_id_not_empty CHECK (app_id != '');
ALTER TABLE errors ADD CONSTRAINT errors_message_not_empty CHECK (message != '');
ALTER TABLE errors ADD CONSTRAINT errors_status_valid CHECK (status IN ('new', 'treated', 'deleted'));

-- Trigger pour mettre à jour updated_at automatiquement
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_apps_updated_at 
    BEFORE UPDATE ON apps 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_errors_updated_at 
    BEFORE UPDATE ON errors 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
