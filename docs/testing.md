# 🧪 Guide de test pour Terrors

Ce guide vous explique comment tester le code de Terrors pour apprendre et vérifier que tout fonctionne.

## 📚 Types de tests

### 1. Tests unitaires (Go)

Les tests unitaires testent chaque fonction/service individuellement.

**Fichiers de test** : `*_test.go` dans chaque package

**Lancer les tests** :

```bash
# Tous les tests
go test ./...

# Tests d'un package spécifique
go test ./internal/services

# Tests avec détails
go test -v ./internal/services

# Tests avec couverture
go test -cover ./internal/services
```

**Configuration** :

Les tests nécessitent une base de données PostgreSQL. Configurez :

```bash
# Option 1: Utiliser la même DB que le serveur
export PG_URL="postgres://user:pass@localhost:5432/terrors"

# Option 2: Utiliser une DB de test séparée (recommandé)
export PG_URL_TEST="postgres://user:pass@localhost:5432/terrors_test"
```

**Exemples de tests** :

- `internal/services/app_test.go` : Tests du service App
- `internal/services/error_test.go` : Tests du service Error

### 2. Tests manuels (Script shell)

Le script `test_manual.sh` teste tous les endpoints de l'API.

**Utilisation** :

```bash
# Rendre exécutable
chmod +x test_manual.sh

# Lancer les tests
./test_manual.sh

# Avec un token admin personnalisé
ADMIN_TOKEN="mon-token" ./test_manual.sh

# Avec un webhook Discord (optionnel)
DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..." ./test_manual.sh
```

**Ce que le script teste** :

1. ✅ Endpoint Home (`/`)
2. ✅ Dashboard (`/overlook`)
3. ✅ Création d'une app
4. ✅ Liste des apps
5. ✅ Récupération d'une app
6. ✅ Capture erreur frontend (`/sadako`)
7. ✅ Capture erreur backend (`/jason`)
8. ✅ Liste des erreurs
9. ✅ Statistiques d'erreurs
10. ✅ Création webhook Discord (si URL fournie)
11. ✅ Liste des webhooks

### 3. Tests avec curl

Vous pouvez aussi tester manuellement avec curl :

```bash
# 1. Créer une app
curl -X POST http://localhost:3000/api/apps \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mon App",
    "description": "Test",
    "origins": ["https://example.com"]
  }'

# 2. Capturer une erreur
curl -X POST http://localhost:3000/sadako \
  -H "Content-Type: application/json" \
  -d '{
    "appId": "app_xxxxx",
    "message": "Test error",
    "stack": "Error: test",
    "fingerprint": "test123",
    "url": "https://example.com",
    "type": "error"
  }'

# 3. Voir les erreurs
curl -X GET "http://localhost:3000/api/errors?appId=app_xxxxx" \
  -H "Authorization: Bearer your-admin-token"
```

## 🎯 Comprendre les tests

### Structure d'un test Go

```go
func TestNomDeLaFonction(t *testing.T) {
    // Arrange : Préparer les données
    db := getTestDB(t)
    service := NewAppService(db)

    // Act : Exécuter la fonction
    result, err := service.CreateApp(models.AppRequest{
        Name: "Test",
    })

    // Assert : Vérifier le résultat
    if err != nil {
        t.Errorf("Erreur inattendue: %v", err)
    }
    if result == nil {
        t.Fatal("Résultat est nil")
    }
}
```

### Patterns de test

1. **Table-driven tests** : Plusieurs cas de test dans une boucle
2. **Setup/Teardown** : Préparer et nettoyer les données
3. **Helpers** : Fonctions utilitaires réutilisables

## 🔍 Déboguer les tests

```bash
# Voir les logs détaillés
go test -v ./internal/services

# Arrêter au premier échec
go test -failfast ./internal/services

# Exécuter un test spécifique
go test -run TestAppService_CreateApp ./internal/services

# Avec race detector (détecte les conditions de course)
go test -race ./internal/services
```

## 📊 Couverture de code

```bash
# Générer un rapport de couverture
go test -coverprofile=coverage.out ./internal/services
go tool cover -html=coverage.out

# Voir le pourcentage
go test -cover ./internal/services
```

## 🚀 Workflow recommandé

1. **Écrire le code** dans un service
2. **Écrire les tests** dans `*_test.go`
3. **Lancer les tests** : `go test ./internal/services`
4. **Vérifier la couverture** : `go test -cover`
5. **Tester manuellement** : `./test_manual.sh`
6. **Itérer** jusqu'à ce que tout fonctionne

## 💡 Conseils pour apprendre

1. **Lisez les tests existants** : Ils montrent comment utiliser les services
2. **Modifiez les tests** : Changez les valeurs et voyez ce qui se passe
3. **Ajoutez vos propres tests** : Testez des cas limites
4. **Utilisez les logs** : Ajoutez `t.Logf()` pour voir ce qui se passe
5. **Testez manuellement** : Utilisez curl pour comprendre les endpoints

## 🐛 Résoudre les problèmes

- **Tests qui échouent** : Lisez le message d'erreur, il indique ce qui ne va pas
- **DB non disponible** : Vérifiez `PG_URL` ou `PG_URL_TEST`
- **Port déjà utilisé** : Changez le port dans `.env`
- **Token invalide** : Vérifiez `ADMIN_TOKEN` dans `.env`
