# 🔔 Webhooks - Guide d'utilisation

Terrors supporte les webhooks pour notifier automatiquement d'autres services lorsqu'une erreur est capturée.

## 🎯 Types de webhooks supportés

- **Discord** : Envoie des notifications dans un channel Discord
- **GitHub** : Crée automatiquement des issues GitHub

## 📋 Création d'un webhook

### Discord

1. **Créer un webhook Discord** :

   - Allez dans les paramètres de votre serveur Discord
   - Channels → Intégrations → Webhooks → Nouveau webhook
   - Copiez l'URL du webhook

2. **Créer le webhook dans Terrors** :

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_xxxxxxxx" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "discord",
    "url": "https://discord.com/api/webhooks/123456789/abcdefghijklmnop",
    "config": {
      "username": "Terrors Bot"
    }
  }'
```

### GitHub

1. **Créer un token GitHub** :

   - GitHub → Settings → Developer settings → Personal access tokens
   - Créez un token avec la permission `repo`

2. **Créer le webhook dans Terrors** :

**Option 1 : Sans URL (recommandé)** - L'URL est construite automatiquement :

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_xxxxxxxx" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "github",
    "config": {
      "token": "ghp_xxxxxxxxxxxxxxxxxxxx",
      "owner": "username",
      "repo": "repository-name",
      "labels": ["bug", "error", "terrors"]
    }
  }'
```

**Option 2 : Avec URL complète** :

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_xxxxxxxx" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "github",
    "url": "https://api.github.com/repos/owner/repo/issues",
    "config": {
      "token": "ghp_xxxxxxxxxxxxxxxxxxxx",
      "owner": "username",
      "repo": "repository-name",
      "labels": ["bug", "error", "terrors"]
    }
  }'
```

## 🔄 Fonctionnement automatique

Une fois configuré, les webhooks sont **automatiquement déclenchés** :

1. Une erreur est capturée (via `/sadako` ou `/jason`)
2. L'erreur est sauvegardée en base
3. Tous les webhooks actifs de l'application sont déclenchés
4. Chaque webhook envoie une notification formatée selon son type

## 📝 Format des notifications

### Discord

Les notifications Discord utilisent des **embeds** avec :

- Titre avec le type d'erreur
- Description avec le message
- Champs : Fingerprint, URL, App ID
- Stack trace dans un bloc de code
- Timestamp et Error ID

### GitHub

Les issues GitHub contiennent :

- **Titre** : `[Error] Message d'erreur`
- **Body** : Détails complets avec sections markdown
- **Labels** : Labels configurés + type d'erreur (frontend/backend)

## 🔍 Voir les webhooks configurés

```bash
curl -X GET "http://localhost:3000/api/webhooks?appId=app_xxxxxxxx" \
  -H "Authorization: Bearer your-admin-token"
```

## 🗑️ Supprimer un webhook

```bash
curl -X DELETE "http://localhost:3000/api/webhooks/123" \
  -H "Authorization: Bearer your-admin-token"
```

## ⚙️ Configuration avancée

### Discord

- `username` : Nom du bot dans Discord (optionnel)

### GitHub

- `token` : Token GitHub (requis)
- `owner` : Propriétaire du repo (requis si pas dans l'URL)
- `repo` : Nom du repo (requis si pas dans l'URL)
- `labels` : Liste de labels à ajouter (optionnel)

## 🐛 Dépannage

- **Webhook non déclenché** : Vérifiez que le webhook est `isActive: true`
- **Erreur Discord** : Vérifiez que l'URL du webhook est valide
- **Erreur GitHub** : Vérifiez le token et les permissions `repo`
