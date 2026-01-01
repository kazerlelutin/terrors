# 📚 Exemples d'intégration Terrors

Ce dossier contient des exemples pour intégrer Terrors dans différents backends.

## 🎯 Principe

Terrors a deux endpoints distincts :

- **`/sadako`** : Pour les erreurs **frontend** (JavaScript dans le navigateur)
- **`/jason`** : Pour les erreurs **backend** (Go, Node.js, Python, etc.)

Les deux acceptent le même format JSON :

- `appId` : L'ID de votre application
- `message` : Le message d'erreur
- `stack` : La stack trace
- `fingerprint` : Hash SHA-1 pour la déduplication
- `url` : L'URL où l'erreur s'est produite
- `ts` : Timestamp
- `type` : Type d'erreur (`error`, `panic`, etc.)

## 🚀 Exemples disponibles

### Go

Voir `backend-go/main.go` pour un exemple complet avec :

- Client Terrors réutilisable
- Capture d'erreurs
- Protection contre les panics
- Middleware HTTP

**Utilisation :**

```go
client := NewTerrorsClient("http://localhost:3000", "app_xxxxxxxx")

// Capturer une erreur
err := fmt.Errorf("database error")
client.CaptureError(err, "http://localhost:8080/api/users")

// Protection contre panic
defer client.CapturePanic("http://localhost:8080/api/process")
```

### Node.js

Voir `backend-nodejs/terrors.js` pour un client Node.js avec :

- Client simple
- Middleware Express
- Plugin Fastify
- Wrapper pour fonctions async

**Utilisation :**

```javascript
const TerrorsClient = require('./terrors')

const client = new TerrorsClient('http://localhost:3000', 'app_xxxxxxxx')

// Capturer une erreur
await client.captureError(new Error('Database error'), '/api/users')

// Avec Express
app.use(client.expressMiddleware())

// Avec Fastify
await fastify.register(client.fastifyPlugin())
```

**Exemple complet Fastify :** Voir `backend-nodejs/fastify-example.js`

### Bun

Voir `backend-bun/` pour un client Bun avec :

- Client TypeScript natif
- Intégration avec `Bun.serve()`
- Utilise les APIs natives de Bun (fetch, crypto)

**Utilisation :**

```typescript
import TerrorsClient from './terrors'

const terrors = new TerrorsClient('http://localhost:3000', 'app_xxxxxxxx')

// Capturer une erreur
await terrors.captureError(new Error('Database error'), '/api/users')

// Avec Bun.serve()
Bun.serve({
  async fetch(request) {
    try {
      // Votre code
      throw new Error('Something went wrong')
    } catch (error) {
      await terrors.captureError(error as Error, new URL(request.url).pathname)
      return new Response('Error', { status: 500 })
    }
  },
  error(error) {
    terrors.captureError(error, 'unknown').catch(console.error)
    return new Response('Internal Server Error', { status: 500 })
  },
})
```

**Exemple complet :** Voir `backend-bun/server.ts`

## 🔧 Intégration simple (n'importe quel langage)

Si vous voulez juste envoyer une erreur sans package, voici un exemple curl pour backend :

```bash
curl -X POST http://localhost:3000/jason \
  -H "Content-Type: application/json" \
  -d '{
    "appId": "app_xxxxxxxx",
    "message": "Database connection failed",
    "stack": "Error: connection timeout\n    at connect (db.js:42:11)",
    "fingerprint": "a1b2c3d4e5f6...",
    "url": "http://localhost:8080/api/users",
    "ts": 1234567890,
    "type": "error"
  }'
```

## 💡 Calcul du fingerprint

Le fingerprint est un hash SHA-1 de `message + première ligne du stack` :

```javascript
// JavaScript
const crypto = require('crypto')
const message = 'Error message'
const stack = 'Error: ...\n    at function (file.js:10:5)'
const topFrame = stack.split('\n')[1] || ''
const raw = message + '\n' + topFrame
const fingerprint = crypto.createHash('sha1').update(raw).digest('hex')
```

```go
// Go
import (
    "crypto/sha1"
    "encoding/hex"
)

stackLines := strings.Split(stack, "\n")
topFrame := ""
if len(stackLines) > 1 {
    topFrame = stackLines[1]
}
raw := message + "\n" + topFrame
hash := sha1.Sum([]byte(raw))
fingerprint := hex.EncodeToString(hash[:])
```

## 🎯 Bonnes pratiques

1. **Ne bloquez pas votre application** : Envoyez les erreurs en arrière-plan (goroutine, async, etc.)
2. **Gérez les erreurs d'envoi** : Si Terrors est down, votre app ne doit pas crasher
3. **Utilisez le fingerprint** : Pour éviter les doublons, calculez-le correctement
4. **Incluez le contexte** : URL, user ID, request ID si disponible

## 🔒 Sécurité

- L'endpoint `/sadako` est public mais valide :
  - Que l'`appId` existe et est actif
  - Que l'origin est autorisée (si configurée)
- Pour la gestion (créer apps, voir erreurs), utilisez le `ADMIN_TOKEN`
