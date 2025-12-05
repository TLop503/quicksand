  const iframe = document.getElementById('vncWindow');

  const themeSwitch = document.getElementById('themeSwitch');

  function applyTheme(theme) {
    document.body.classList.toggle('light', theme === 'light');
    if (themeSwitch) themeSwitch.checked = (theme === 'light');
    localStorage.setItem('theme', theme);
  }

  // Load saved theme on startup
  const savedTheme = localStorage.getItem('theme') || 'dark';
  applyTheme(savedTheme);

// Listen for toggle changes
  if (themeSwitch) {
    themeSwitch.addEventListener('change', () => {
      applyTheme(themeSwitch.checked ? 'light' : 'dark');
    });
  }

  // Generic POST to your API with optional JSON body
  async function api(action, body = {}) {
    const res = await fetch(`/api/${action}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: Object.keys(body).length ? JSON.stringify(body) : null
    });
    if (!res.ok) throw new Error(`/${action} failed`);
    return res.json(); // expect { ok: true, iframeUrl?: string }
  }

  // Unsure if this is needed, but was reading that this is a good function to have
  // Poll /health until server/container is ready (adjust timeout/interval as needed)
  async function waitUntilReady({ timeoutMs = 30000, intervalMs = 1000 } = {}) {
    const start = Date.now();
    while (Date.now() - start < timeoutMs) {
      try {
        const r = await fetch('/api/health', { cache: 'no-store' });
        if (r.ok) return true;
      } catch (_) { /* ignore while container is cycling */ }
      await new Promise(r => setTimeout(r, intervalMs));
    }
    throw new Error('Timed out waiting for container readiness');
  }

  // Update iframe to a provided URL, or keep current target
  function updateIframe(url) {
    iframe.src = url || iframe.src; // replace or just refresh
  }

  // Button actions
  async function restartContainer() {
    disableButtons(true);
    try {
      const { iframeUrl } = await api('restart');
      await waitUntilReady();
      updateIframe(iframeUrl || 'http://localhost:5800');
    } 

    finally {
      disableButtons(false);
    }
  }

  async function stopContainer() {
    disableButtons(true);
    try { 
        await api('stop'); 
    } 

    finally { 
        disableButtons(false); 
    }
  }

  async function startContainer() {
    disableButtons(true);
    try {
      const { iframeUrl } = await api('start');
      await waitUntilReady();
      updateIframe(iframeUrl || 'http://localhost:5800');
    } 

    finally {
      disableButtons(false);
    }
  }

  async function swapBrowsingStylesTor() {
    disableButtons(true);
    try {
      const { iframeUrl } = await api('swap', { to: 'tor' });
      await waitUntilReady();
      updateIframe(iframeUrl || 'http://localhost:5800');
    } 

    finally {
      disableButtons(false);
    }
  }

  async function swapBrowsingStylesFirefox() {
    disableButtons(true);
    try {
      const { iframeUrl } = await api('swap', { to: 'firefox' });
      await waitUntilReady();
      updateIframe(iframeUrl || 'http://localhost:5800');
    } 

    finally {
      disableButtons(false);
    }
  }

  function disableButtons(disabled) {
    document.querySelectorAll('button').forEach(b => b.disabled = disabled);
  }