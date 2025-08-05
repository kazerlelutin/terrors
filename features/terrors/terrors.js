(async () => {
  const
    el = document.querySelector('script[src*="terrors.js"]'),
    url = el?.getAttribute('src'),
    origin = url?.split('/').slice(0, -2).join('/'),
    appId = el?.getAttribute('appid');

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

  async function handleError(message, stack, type = 'error') {
    const fingerprint = await computeFingerprint(message, stack);
    fetch(origin + '/sadako', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        appId,
        message,
        stack,
        fingerprint,
        url: window.location.href,
        ts: Date.now(),
        type
      })
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