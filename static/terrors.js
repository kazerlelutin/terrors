(async () => {
  const
    el = document.querySelector('script[src*="terrors.js"]'),
    url = el?.getAttribute('src'),
    origin = url?.split('/').slice(0, -2).join('/'),
    appId = el?.getAttribute('data-app-id');

  if (!appId) {
    console.error('App ID for terrors.js not found');
    return;
  } else {
    console.log('Terrors.js loaded');
  }

  async function computeFingerprint(message, stack) {
    try {
      const topFrame = stack.split('\n')[1] || '';
      const raw = message + '\n' + topFrame;
      const buf = await crypto.subtle.digest('SHA-1', new TextEncoder().encode(raw));
      return Array.from(new Uint8Array(buf))
        .map(b => b.toString(16).padStart(2, '0'))
        .join('');
    } catch (error) {
      return 'unknown';
    }
  }

  // Cache pour éviter les envois multiples de la même erreur rapidement
  const sentErrors = new Map();
  const DEBOUNCE_TIME = 5000; // 5 secondes

  async function handleError(message, stack, type = 'error') {
    const fingerprint = await computeFingerprint(message, stack);
    const cacheKey = `${fingerprint}-${type}`;

    // Vérifier si on a déjà envoyé cette erreur récemment
    const lastSent = sentErrors.get(cacheKey);
    const now = Date.now();

    if (lastSent && (now - lastSent) < DEBOUNCE_TIME) {
      console.log('Error déjà envoyée récemment, ignorée:', message);
      return;
    }

    // Marquer comme envoyée
    sentErrors.set(cacheKey, now);

    // Nettoyer le cache périodiquement (garder seulement les 100 dernières)
    if (sentErrors.size > 100) {
      const oldestKey = sentErrors.keys().next().value;
      sentErrors.delete(oldestKey);
    }

    fetch(origin + '/sadako', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        appId,
        message,
        stack,
        fingerprint,
        url: window.location.href,
        ts: now,
        type
      })
    }).catch(err => {
      console.error('Erreur envoi à Terrors:', err);
      // En cas d'erreur, on retire du cache pour pouvoir réessayer
      sentErrors.delete(cacheKey);
    });
  }

  window.addEventListener('error', event => {
    console.log('Error captured:', event.message);
    handleError(event.message, event.error?.stack || '', 'error');
  });

  window.addEventListener('unhandledrejection', event => {
    console.log('Promise rejected:', event.reason);
    const reason = event.reason;
    handleError(
      reason?.message || String(reason),
      reason?.stack || '',
      'unhandledrejection'
    );
    event.preventDefault();
  });
})();
