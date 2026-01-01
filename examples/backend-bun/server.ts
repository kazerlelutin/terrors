/**
 * Exemple d'utilisation de Terrors avec Bun.serve()
 */

import TerrorsClient from './terrors'

// Initialiser le client Terrors
const terrors = new TerrorsClient('http://localhost:3000', 'app_xxxxxxxx')

// Fonction helper pour gérer les erreurs
async function handleError(error: Error, request: Request): Promise<Response> {
  // Envoyer l'erreur à Terrors (en arrière-plan, ne pas bloquer)
  terrors.captureError(error, new URL(request.url).pathname).catch(console.error)

  // Retourner une réponse d'erreur
  return new Response(
    JSON.stringify({
      error: error.message,
      status: 500
    }),
    {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    }
  )
}

// Serveur Bun avec routes
Bun.serve({
  port: 3001,
  async fetch(request: Request) {
    try {
      const url = new URL(request.url)
      const path = url.pathname

      // Route qui fonctionne
      if (path === '/api/users' && request.method === 'GET') {
        return new Response(
          JSON.stringify([
            { id: 1, name: 'John Doe' },
            { id: 2, name: 'Jane Doe' }
          ]),
          {
            headers: { 'Content-Type': 'application/json' }
          }
        )
      }

      // Route qui génère une erreur
      if (path === '/api/users/error' && request.method === 'GET') {
        throw new Error('User not found')
      }

      // Route avec erreur async
      if (path === '/api/data' && request.method === 'GET') {
        // Simuler une erreur de base de données
        throw new Error('Database connection failed')
      }

      // Route avec paramètre
      if (path.startsWith('/api/users/') && request.method === 'GET') {
        const id = path.split('/').pop()
        if (id === 'error') {
          throw new Error(`User ${id} not found`)
        }
        return new Response(
          JSON.stringify({ id, name: 'John Doe' }),
          {
            headers: { 'Content-Type': 'application/json' }
          }
        )
      }

      // Route 404
      return new Response('Not Found', { status: 404 })
    } catch (error) {
      // Capturer et gérer l'erreur
      return handleError(error as Error, request)
    }
  },
  error(error: Error) {
    // Handler global pour les erreurs non capturées
    console.error('Erreur non capturée:', error)
    terrors.captureError(error, 'unknown').catch(console.error)
    return new Response('Internal Server Error', { status: 500 })
  }
})

console.log('🚀 Serveur Bun démarré sur http://localhost:3001')
console.log('📡 Terrors configuré pour capturer les erreurs')
console.log('\nRoutes disponibles:')
console.log('  GET /api/users - Liste des utilisateurs')
console.log('  GET /api/users/:id - Détails d\'un utilisateur')
console.log('  GET /api/users/error - Génère une erreur')
console.log('  GET /api/data - Génère une erreur de DB')

