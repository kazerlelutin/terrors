# 🔔 Webhooks GitHub - Repositories d'Organisation

Guide spécifique pour configurer les webhooks GitHub avec des repositories d'organisation.

## 🎯 Le problème

Quand vous essayez de créer une issue dans un repository d'organisation, vous obtenez :

```
403 Forbidden: Resource not accessible by personal access token
```

Cela signifie que votre token personnel n'a pas accès à l'organisation.

## ✅ Solution : Autoriser le token pour l'organisation

### Méthode 1 : Lors de la création du token (Recommandé)

1. **Créer un nouveau token** :

   - GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Cliquez sur "Generate new token (classic)"
   - Donnez un nom au token (ex: "Terrors - Asso-M05")
   - Cochez la permission `repo` (accès complet aux repositories)
   - **Important** : Dans la section "Organization access", sélectionnez votre organisation
   - Cliquez sur "Generate token"
   - **Copiez le token immédiatement** (il ne sera plus visible après)

2. **Utiliser le token** :
   ```json
   {
     "type": "github",
     "config": {
       "token": "ghp_votre_nouveau_token",
       "owner": "Asso-M05",
       "repo": "liquid",
       "labels": ["bug", "error"]
     }
   }
   ```

### Méthode 2 : Autoriser un token existant

Si vous avez déjà un token :

1. **Via les paramètres de l'organisation** :

   - Allez sur votre organisation GitHub (ex: `github.com/Asso-M05`)
   - Settings → Third-party access → Personal access tokens
   - Trouvez votre token dans la liste
   - Cliquez sur "Grant" ou "Approve" pour autoriser l'accès

2. **Via les paramètres personnels** :
   - GitHub → Settings → Developer settings → Personal access tokens
   - Trouvez votre token
   - Cliquez sur "Configure" ou "Edit"
   - Dans "Organization access", sélectionnez votre organisation
   - Sauvegardez

## 🔐 Alternative : GitHub App (Plus sécurisé)

Pour une solution plus professionnelle, vous pouvez créer une GitHub App :

1. **Créer une GitHub App** :

   - Organisation → Settings → Developer settings → GitHub Apps
   - New GitHub App
   - Configurez les permissions : `Issues: Write`
   - Installez l'app sur votre organisation/repository

2. **Générer un token d'installation** :
   - Utilisez l'API GitHub pour générer un token d'installation
   - Ce token a accès uniquement aux repositories autorisés

**Note** : L'implémentation GitHub App nécessite plus de code. Pour l'instant, utilisez un PAT avec accès organisation.

## 🧪 Tester l'accès

Avant de configurer le webhook, testez que votre token fonctionne :

```bash
curl -H "Authorization: token ghp_votre_token" \
     -H "Accept: application/vnd.github.v3+json" \
     https://api.github.com/repos/Asso-M05/liquid
```

Si vous obtenez les infos du repo, le token fonctionne. Sinon, vérifiez les permissions.

## 📝 Configuration dans Terrors

Une fois le token autorisé, configurez le webhook normalement :

```bash
curl -X POST "http://localhost:3000/api/webhooks?appId=app_xxxxxxxx" \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "github",
    "config": {
      "token": "ghp_votre_token_autorisé",
      "owner": "Asso-M05",
      "repo": "liquid",
      "labels": ["bug", "error"]
    }
  }'
```

## 🐛 Dépannage

### Erreur 403 persistante

1. Vérifiez que le token a la permission `repo`
2. Vérifiez que le token est autorisé pour l'organisation
3. Vérifiez que vous avez les droits sur le repository (membre de l'orga avec accès au repo)
4. Essayez de créer une issue manuellement avec le token pour tester

### Token fonctionne mais webhook échoue

- Vérifiez les logs du serveur Terrors
- Vérifiez que le token est bien stocké dans la config (pas d'espaces, caractères corrects)
- Testez la création d'issue directement avec curl

## 🔒 Sécurité

- **Ne partagez jamais votre token** : Il donne accès à vos repositories
- **Utilisez des tokens avec permissions minimales** : Seulement `repo` si possible
- **Régénérez le token si compromis** : GitHub → Settings → Developer settings → Revoke
