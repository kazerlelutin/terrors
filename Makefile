# Makefile pour faciliter les tests et le développement

.PHONY: test test-cover test-services test-manual help

# Variables
TEST_DB_URL ?= $(PG_URL_TEST)
ifndef TEST_DB_URL
	TEST_DB_URL = $(PG_URL)
endif

# Aide
help:
	@echo "🎭 Terrors - Commandes disponibles"
	@echo ""
	@echo "Tests:"
	@echo "  make test          - Lancer tous les tests"
	@echo "  make test-cover    - Tests avec couverture"
	@echo "  make test-services - Tests des services uniquement"
	@echo "  make test-manual   - Tests manuels (script shell)"
	@echo ""
	@echo "Développement:"
	@echo "  make run           - Lancer le serveur"
	@echo "  make migrate       - Lancer les migrations"
	@echo "  make clean         - Nettoyer les fichiers temporaires"
	@echo ""
	@echo "Variables d'environnement:"
	@echo "  PG_URL_TEST        - URL de la DB de test (optionnel)"
	@echo "  ADMIN_TOKEN        - Token admin pour les tests manuels"

# Tests
test:
	@echo "🧪 Lancement des tests..."
	@if [ -z "$(TEST_DB_URL)" ]; then \
		echo "⚠️  PG_URL ou PG_URL_TEST non défini, certains tests seront ignorés"; \
	fi
	go test -v ./...

test-cover:
	@echo "📊 Tests avec couverture..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo ""
	@echo "📄 Rapport HTML généré: coverage.out"
	@echo "   Ouvrir avec: go tool cover -html=coverage.out"

test-services:
	@echo "🧪 Tests des services..."
	go test -v ./internal/services

test-manual:
	@echo "🧪 Tests manuels..."
	@if [ -z "$(ADMIN_TOKEN)" ]; then \
		echo "⚠️  ADMIN_TOKEN non défini, utilisez: ADMIN_TOKEN=xxx make test-manual"; \
	fi
	@bash test_manual.sh

# Développement
run:
	@echo "🚀 Lancement du serveur..."
	go run cmd/server/main.go

migrate:
	@echo "🔄 Lancement des migrations..."
	go run cmd/migrate/main.go

# Nettoyage
clean:
	@echo "🧹 Nettoyage..."
	rm -f coverage.out
	rm -f tmp/*.log
	@echo "✅ Nettoyage terminé"

