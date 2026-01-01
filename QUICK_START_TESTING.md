# 🚀 Démarrage rapide - Tests

Guide rapide pour tester et comprendre le code de Terrors.

## ⚡ Démarrage en 3 étapes

### 1. Préparer la base de données

```bash
# Créer une base de données de test (optionnel mais recommandé)
createdb terrors_test

# Ou utiliser la même que le serveur
export PG_URL="postgres://user:pass@localhost:5432/terrors"
```

### 2. Lancer les tests unitaires

```bash
# Tous les tests
go test ./...

# Tests avec détails (recommandé pour apprendre)
go test -v ./internal/services

# Tests avec couverture
go test -cover ./internal/services
```

### 3. Tester manuellement l'API

```bash
# Démarrer le serveur dans un terminal
go run cmd/server/main.go

# Dans un autre terminal, lancer les tests manuels
# (Sur Windows, utilisez Git Bash ou WSL)
bash test_manual.sh

# Ou avec PowerShell (Windows)
# Modifiez test_manual.sh pour utiliser curl.exe
```

## 📖 Comprendre les tests

### Structure d'un test

```go
func TestNomFonction(t *testing.T) {
    // 1. ARRANGE : Préparer
    db := getTestDB(t)
    service := NewAppService(db)

    // 2. ACT : Exécuter
    result, err := service.CreateApp(...)

    // 3. ASSERT : Vérifier
    if err != nil {
        t.Errorf("Erreur: %v", err)
    }
}
```

### Exemples de tests à lire

1. **`internal/services/app_test.go`** :

   - `TestAppService_CreateApp` : Comment créer une app
   - `TestAppService_ValidateOrigin` : Comment valider les origins

2. **`internal/services/error_test.go`** :
   - `TestErrorService_SaveError` : Comment sauvegarder une erreur
   - `TestErrorService_ListErrors` : Comment lister les erreurs

## 🎯 Exercices pour apprendre

### Exercice 1 : Modifier un test

1. Ouvrez `internal/services/app_test.go`
2. Trouvez `TestAppService_CreateApp`
3. Ajoutez un nouveau cas de test dans la slice `tests`
4. Relancez : `go test -v ./internal/services -run TestAppService_CreateApp`

### Exercice 2 : Créer votre propre test

Créez `internal/services/webhook_test.go` :

```go
package services

import (
    "testing"
    "terrors/internal/models"
)

func TestWebhookService_CreateWebhook(t *testing.T) {
    db := getTestDB(t)
    defer db.Close()

    service := NewWebhookService(db)

    // Créer une app d'abord
    appService := NewAppService(db)
    app, _, _ := appService.CreateApp(models.AppRequest{
        Name: "Test Webhook App",
    })
    defer cleanupTestData(t, db, app.AppID)

    // Créer un webhook
    webhook, err := service.CreateWebhook(app.AppID, models.WebhookRequest{
        Type: "discord",
        URL:  "https://discord.com/api/webhooks/test",
    })

    if err != nil {
        t.Fatalf("Erreur: %v", err)
    }

    if webhook.Type != "discord" {
        t.Errorf("Type = %v, want 'discord'", webhook.Type)
    }
}
```

### Exercice 3 : Tester avec curl

```bash
# 1. Créer une app
curl -X POST http://localhost:3000/api/apps \
  -H "Authorization: Bearer your-token" \
  -H "Content-Type: application/json" \
  -d '{"name": "Mon Test", "origins": ["https://test.com"]}'

# Notez l'appId dans la réponse, puis :

# 2. Envoyer une erreur
curl -X POST http://localhost:3000/sadako \
  -H "Content-Type: application/json" \
  -d '{
    "appId": "app_xxxxx",
    "message": "Mon erreur de test",
    "stack": "Error: test",
    "fingerprint": "test123",
    "url": "https://test.com",
    "type": "error"
  }'

# 3. Voir les erreurs
curl -X GET "http://localhost:3000/api/errors?appId=app_xxxxx" \
  -H "Authorization: Bearer your-token"
```

## 🔍 Déboguer

### Voir les logs détaillés

```bash
go test -v ./internal/services
```

### Tester une fonction spécifique

```bash
go test -v -run TestAppService_CreateApp ./internal/services
```

### Arrêter au premier échec

```bash
go test -failfast ./internal/services
```

## 💡 Conseils

1. **Lisez les tests** : Ils sont la meilleure documentation
2. **Modifiez-les** : Changez les valeurs et voyez ce qui se passe
3. **Ajoutez des logs** : Utilisez `t.Logf("Valeur: %v", valeur)`
4. **Testez manuellement** : Utilisez curl pour comprendre les endpoints
5. **Utilisez le debugger** : VS Code peut déboguer les tests Go

## 📚 Ressources

- **Documentation Go Testing** : https://pkg.go.dev/testing
- **Guide complet** : Voir `docs/testing.md`
- **Makefile** : Utilisez `make help` pour voir les commandes disponibles

## 🐛 Problèmes courants

**"PG_URL non défini"** :

```bash
export PG_URL="postgres://user:pass@localhost:5432/terrors"
```

**"Port déjà utilisé"** :

```bash
# Changez le port dans .env
PORT=3001
```

**"Tests qui échouent"** :

- Lisez le message d'erreur, il indique ce qui ne va pas
- Vérifiez que la DB est accessible
- Vérifiez que les migrations sont appliquées

---

**Bon apprentissage ! 🎓**
