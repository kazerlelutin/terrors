/**
 * Client Terrors pour Bun
 * Usage simple pour capturer les erreurs backend
 */

class TerrorsClient {
  private baseURL: string
  private appId: string

  constructor(baseURL: string, appId: string) {
    this.baseURL = baseURL
    this.appId = appId
  }

  /**
   * Calcule le fingerprint d'une erreur (SHA-1)
   */
  private async computeFingerprint(message: string, stack: string): Promise<string> {
    const stackLines = stack ? stack.split('\n')[1] || '' : ''
    const raw = message + '\n' + stackLines

    // Utiliser crypto.subtle pour SHA-1
    const encoder = new TextEncoder()
    const data = encoder.encode(raw)
    const hashBuffer = await crypto.subtle.digest('SHA-1', data)
    const hashArray = Array.from(new Uint8Array(hashBuffer))
    return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
  }

  /**
   * Capture une erreur et l'envoie à Terrors (endpoint /jason pour backend)
   */
  async captureError(error: Error, url: string = ''): Promise<any> {
    const message = error.message || String(error)
    const stack = error.stack || ''
    const fingerprint = await this.computeFingerprint(message, stack)

    const errorReq = {
      appId: this.appId,
      message: message,
      stack: stack,
      fingerprint: fingerprint,
      url: url,
      ts: Date.now(),
      type: 'error'
    }

    try {
      const response = await fetch(this.baseURL + '/jason', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(errorReq)
      })

      if (!response.ok) {
        const text = await response.text()
        throw new Error(`Server error: ${response.status} - ${text}`)
      }

      return await response.json()
    } catch (err) {
      console.error('❌ Erreur lors de l\'envoi à Terrors:', err)
      throw err
    }
  }

  /**
   * Wrapper pour capturer les erreurs dans une fonction async
   */
  async wrapAsync<T>(fn: () => Promise<T>): Promise<T> {
    try {
      return await fn()
    } catch (error) {
      await this.captureError(error as Error)
      throw error // Re-throw pour que l'app puisse gérer
    }
  }
}

export default TerrorsClient

