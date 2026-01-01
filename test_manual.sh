#!/bin/bash

# Script de test manuel pour Terrors
# Permet de tester les endpoints sans avoir besoin de tests unitaires

BASE_URL="http://localhost:3000"
ADMIN_TOKEN="${ADMIN_TOKEN:-your-admin-token-here}"

echo "🧪 Tests manuels pour Terrors"
echo "=============================="
echo ""

# Couleurs pour la sortie
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour tester un endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    local expected_status=$5

    echo -n "Test: $name ... "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$url")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_status" ]; then
        echo -e "${GREEN}✓${NC} (HTTP $http_code)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        echo -e "${RED}✗${NC} (HTTP $http_code, attendu $expected_status)"
        echo "$body"
    fi
    echo ""
}

# 1. Test Home
test_endpoint "Home" "GET" "/" "" "200"

# 2. Test Dashboard
test_endpoint "Dashboard" "GET" "/overlook" "" "200"

# 3. Créer une app
APP_RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "Test App",
        "description": "App de test",
        "origins": ["https://example.com"]
    }' \
    "$BASE_URL/api/apps")

APP_ID=$(echo "$APP_RESPONSE" | jq -r '.app.appId // empty')
if [ -z "$APP_ID" ]; then
    echo -e "${RED}✗${NC} Impossible de créer une app pour les tests"
    exit 1
fi

echo -e "${GREEN}✓${NC} App créée: $APP_ID"
echo ""

# 4. Lister les apps
test_endpoint "List Apps" "GET" "/api/apps" "" "200"

# 5. Récupérer l'app
test_endpoint "Get App" "GET" "/api/apps/$APP_ID" "" "200"

# 6. Envoyer une erreur frontend
test_endpoint "Capture Frontend Error" "POST" "/sadako" \
    "{
        \"appId\": \"$APP_ID\",
        \"message\": \"Test error from script\",
        \"stack\": \"Error: test\\n    at test.js:10:5\",
        \"fingerprint\": \"test123\",
        \"url\": \"https://example.com/test\",
        \"type\": \"error\"
    }" "200"

# 7. Envoyer une erreur backend
test_endpoint "Capture Backend Error" "POST" "/jason" \
    "{
        \"appId\": \"$APP_ID\",
        \"message\": \"Backend test error\",
        \"stack\": \"Error: backend\\n    at server.go:42\",
        \"fingerprint\": \"backend123\",
        \"url\": \"/api/test\",
        \"type\": \"error\"
    }" "200"

# 8. Lister les erreurs
test_endpoint "List Errors" "GET" "/api/errors?appId=$APP_ID" "" "200"

# 9. Statistiques
test_endpoint "Error Stats" "GET" "/api/errors/stats?appId=$APP_ID" "" "200"

# 10. Créer un webhook Discord (si URL fournie)
if [ -n "$DISCORD_WEBHOOK_URL" ]; then
    test_endpoint "Create Discord Webhook" "POST" "/api/webhooks?appId=$APP_ID" \
        "{
            \"type\": \"discord\",
            \"url\": \"$DISCORD_WEBHOOK_URL\"
        }" "201"
fi

# 11. Lister les webhooks
test_endpoint "List Webhooks" "GET" "/api/webhooks?appId=$APP_ID" "" "200"

echo ""
echo "✅ Tests terminés !"
echo ""
echo "Pour nettoyer, supprimez l'app:"
echo "curl -X DELETE \"$BASE_URL/api/apps/$APP_ID\" -H \"Authorization: Bearer $ADMIN_TOKEN\""

