# 🔔 Webhooks GitHub - Documentation

## Comment ça marche ?

Les webhooks GitHub permettent de créer automatiquement des issues GitHub lorsqu'une nouvelle erreur est capturée par Terrors.

## 📋 Prérequis

1. **Token GitHub** : Vous devez créer un Personal Access Token (PAT) avec les permissions `repo` pour créer des issues

   - Allez sur GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Créez un token avec la permission `repo`
   - **Important pour les organisations** : Le token doit être autorisé pour l'organisation
     - Lors de la création du token, cochez les organisations nécessaires dans "Organization access"
     - Ou allez dans les paramètres de l'organisation → Third-party access → Personal access tokens
     - Autorisez votre token pour l'organisation

2. **Repository GitHub** : Vous devez avoir un repository où créer les issues
   - Pour les repos d'organisation, assurez-vous que le token a accès à l'organisation

## 🚀 Configuration

### 1. Créer un webhook GitHub

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_xxxxxxxx" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "github",
    "url": "https://api.github.com/repos/owner/repo/issues",
    "config": {
      "token": "ghp_xxxxxxxxxxxxxxxxxxxx",
      "labels": ["bug", "error", "terrors"],
      "owner": "owner",
      "repo": "repo"
    }
  }'
```

### 2. Format de la configuration

Le champ `config` doit contenir :

```json
{
  "token": "ghp_xxxxxxxxxxxxxxxxxxxx", // Token GitHub (requis)
  "owner": "username", // Propriétaire du repo (requis)
  "repo": "repository-name", // Nom du repo (requis)
  "labels": ["bug", "error"] // Labels optionnels
}
```

**Note** : L'URL peut être soit :

- L'URL complète : `https://api.github.com/repos/owner/repo/issues`
- Ou juste le repo : `owner/repo` (le système complétera l'URL)

## 🔄 Fonctionnement

1. **Une erreur est capturée** (via `/sadako` ou `/jason`)
2. **L'erreur est sauvegardée** en base de données
3. **Les webhooks actifs sont déclenchés** automatiquement
4. **Pour chaque webhook GitHub** :
   - Le système récupère la config (token, owner, repo, labels)
   - Formate une issue GitHub avec :
     - **Titre** : Le message d'erreur (tronqué à 100 caractères)
     - **Body** : Détails complets (message, stack, URL, fingerprint, etc.)
     - **Labels** : Les labels configurés + le type d'erreur (frontend/backend)
   - Envoie une requête POST à l'API GitHub Issues
   - Crée l'issue dans le repository

## 📝 Format de l'issue GitHub

### Titre

```
[Error] Message d'erreur tronqué...
```

### Body

```markdown
## 🎭 Erreur capturée par Terrors

**Type** : frontend/backend
**Fingerprint** : `a1b2c3d4e5f6...`
**URL** : https://example.com/page
**Timestamp** : 2025-01-XX 12:34:56 UTC

### Message
```

Message d'erreur complet

```

### Stack Trace
```

Error: ...
at function (file.js:10:5)
...

```

### Détails
- **App ID** : app_xxxxxxxx
- **Error ID** : 123
- **Status** : new
```

## 🔒 Sécurité

- Le token GitHub est stocké dans la config (champ JSON)
- Il n'est jamais exposé dans les réponses API
- Utilisez HTTPS en production
- Régénérez le token si compromis

## ⚙️ Exemple complet

```bash
# 1. Créer le webhook
curl -X POST "http://localhost:3000/api/webhooks?appId=app_abc123" \
  -H "Authorization: Bearer admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "github",
    "url": "https://api.github.com/repos/monuser/monrepo/issues",
    "config": {
      "token": "ghp_xxxxxxxxxxxx",
      "labels": ["bug", "error"]
    }
  }'

# 2. Une erreur se produit dans votre app
# → Une issue GitHub est automatiquement créée !

# 3. Voir les webhooks configurés
curl -X GET "http://localhost:3000/api/webhooks?appId=app_abc123" \
  -H "Authorization: Bearer admin-token"
```

## 🐛 Dépannage

- **401 Unauthorized** : Vérifiez que le token est valide et a les permissions `repo`
- **403 Forbidden** :
  - **Pour les repos d'organisation** : Le token doit être autorisé pour l'organisation
  - Voir `docs/webhooks-github-orga.md` pour le guide complet
  - GitHub → Settings → Developer settings → Personal access tokens → Autoriser pour l'organisation
- **404 Not Found** : Vérifiez que l'owner/repo existe et que le token a accès
- **422 Unprocessable Entity** : Vérifiez le format de la requête (labels valides, etc.)

## 🏢 Repositories d'Organisation

Si vous utilisez un repository d'organisation et obtenez une erreur 403, consultez le guide détaillé : **`docs/webhooks-github-orga.md`**

En résumé : Le token doit être autorisé pour l'organisation lors de sa création ou via les paramètres de l'organisation.
