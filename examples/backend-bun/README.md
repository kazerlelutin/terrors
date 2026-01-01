# Terrors Client pour Bun

Client TypeScript natif pour Bun qui utilise les APIs natives de Bun.

## Installation

Aucune dépendance nécessaire ! Bun inclut tout ce dont on a besoin :

- `fetch` natif
- `crypto.subtle` pour SHA-1
- Support TypeScript natif

## Utilisation

### Client simple

```typescript
import TerrorsClient from './terrors'

const terrors = new TerrorsClient('http://localhost:3000', 'app_xxxxxxxx')

// Capturer une erreur
await terrors.captureError(new Error('Database error'), '/api/users')
```

### Avec Bun.serve()

```typescript
import TerrorsClient from './terrors'

const terrors = new TerrorsClient('http://localhost:3000', 'app_xxxxxxxx')

Bun.serve({
  async fetch(request) {
    try {
      // Votre code qui peut générer des erreurs
      throw new Error('Something went wrong')
    } catch (error) {
      // Capturer l'erreur
      await terrors.captureError(error as Error, new URL(request.url).pathname)
      return new Response('Error', { status: 500 })
    }
  },
  error(error) {
    // Handler global pour les erreurs non capturées
    terrors.captureError(error, 'unknown').catch(console.error)
    return new Response('Internal Server Error', { status: 500 })
  },
})
```

## Exemple complet

Voir `server.ts` pour un exemple complet avec plusieurs routes.

## Lancer l'exemple

```bash
bun run server.ts
```

Le serveur démarre sur `http://localhost:3001` et toutes les erreurs sont automatiquement envoyées à Terrors.
