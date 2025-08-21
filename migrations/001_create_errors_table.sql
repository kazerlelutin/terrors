-- Migration pour créer la table des erreurs
CREATE TABLE IF NOT EXISTS errors (
    id SERIAL PRIMARY KEY,
    app_id VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    stack TEXT,
    fingerprint VARCHAR(64) NOT NULL,
    url TEXT,
    type VARCHAR(50) DEFAULT 'error',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index pour optimiser les requêtes
CREATE INDEX IF NOT EXISTS idx_errors_app_id ON errors(app_id);
CREATE INDEX IF NOT EXISTS idx_errors_fingerprint ON errors(fingerprint);
CREATE INDEX IF NOT EXISTS idx_errors_created_at ON errors(created_at);
CREATE INDEX IF NOT EXISTS idx_errors_type ON errors(type);

-- Trigger pour mettre à jour updated_at automatiquement
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_errors_updated_at 
    BEFORE UPDATE ON errors 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
