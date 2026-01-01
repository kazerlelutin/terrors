/**
 * Client Terrors pour Node.js
 * Usage simple pour capturer les erreurs backend
 */

const crypto = require('crypto');
class TerrorsClient {
  constructor(baseURL, appId) {
    this.baseURL = baseURL;
    this.appId = appId;
  }

  /**
   * Calcule le fingerprint d'une erreur
   */
  computeFingerprint(message, stack) {
    const stackLines = stack ? stack.split('\n')[1] || '' : '';
    const raw = message + '\n' + stackLines;
    const hash = crypto.createHash('sha1').update(raw).digest('hex');
    return hash;
  }

  /**
   * Capture une erreur et l'envoie à Terrors (endpoint /jason pour backend)
   */
  async captureError(error, url = '') {
    const message = error.message || String(error);
    const stack = error.stack || '';
    const fingerprint = this.computeFingerprint(message, stack);

    const errorReq = {
      appId: this.appId,
      message: message,
      stack: stack,
      fingerprint: fingerprint,
      url: url,
      ts: Date.now(),
      type: 'error'
    };

    try {
      const response = await fetch(this.baseURL + '/jason', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(errorReq)
      });

      if (!response.ok) {
        const text = await response.text();
        throw new Error(`Server error: ${response.status} - ${text}`);
      }

      return await response.json();
    } catch (err) {
      console.error('❌ Erreur lors de l\'envoi à Terrors:', err);
      throw err;
    }
  }

  /**
   * Middleware Express pour capturer les erreurs
   */
  expressMiddleware() {
    return (err, req, res, next) => {
      // Envoyer l'erreur à Terrors (en arrière-plan)
      this.captureError(err, req.originalUrl).catch(console.error);

      // Passer à l'erreur suivante
      next(err);
    };
  }

  /**
   * Plugin Fastify pour capturer les erreurs
   */
  fastifyPlugin() {
    return async (fastify, options) => {
      // Hook pour capturer les erreurs
      fastify.setErrorHandler(async (error, request, reply) => {
        // Envoyer l'erreur à Terrors (en arrière-plan, ne pas bloquer)
        this.captureError(error, request.url).catch(console.error);

        // Laisser Fastify gérer l'erreur normalement
        reply.send(error);
      });
    };
  }

  /**
   * Wrapper pour capturer les panics/erreurs dans une fonction async
   */
  async wrapAsync(fn) {
    try {
      return await fn();
    } catch (error) {
      await this.captureError(error);
      throw error; // Re-throw pour que l'app puisse gérer
    }
  }
}

// Exemple d'utilisation
if (require.main === module) {
  const client = new TerrorsClient('http://localhost:3000', 'app_xxxxxxxx');

  // Exemple 1: Capturer une erreur
  client.captureError(new Error('Database connection failed'), 'http://localhost:3000/api/users')
    .then(() => console.log('✅ Erreur envoyée'))
    .catch(console.error);

  // Exemple 2: Avec Express
  // const express = require('express');
  // const app = express();
  // app.use(client.expressMiddleware());

  // Exemple 3: Avec Fastify
  // const fastify = require('fastify')();
  // await fastify.register(client.fastifyPlugin());
  // await fastify.listen({ port: 3000 });
}

module.exports = TerrorsClient;

