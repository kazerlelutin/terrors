/**
 * Exemple d'utilisation de Terrors avec Fastify
 */

const fastify = require('fastify')({ logger: true })
const TerrorsClient = require('./terrors')

// Initialiser le client Terrors
const terrors = new TerrorsClient('http://localhost:3000', 'app_xxxxxxxx')

// Enregistrer le plugin Terrors
fastify.register(terrors.fastifyPlugin())

// Route qui peut générer une erreur
fastify.get('/api/users/:id', async (request, reply) => {
  const { id } = request.params

  // Simuler une erreur si l'ID est invalide
  if (id === 'error') {
    throw new Error('User not found')
  }

  return { id, name: 'John Doe' }
})

// Route avec erreur async
fastify.get('/api/data', async (request, reply) => {
  // Simuler une erreur de base de données
  throw new Error('Database connection failed')
})

// Démarrer le serveur
const start = async () => {
  try {
    await fastify.listen({ port: 3001 })
    console.log('🚀 Serveur Fastify démarré sur http://localhost:3001')
    console.log('📡 Terrors configuré pour capturer les erreurs')
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}

start()

