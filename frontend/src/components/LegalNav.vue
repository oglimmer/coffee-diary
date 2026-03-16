<template>
  <nav class="legal-nav">
    <RouterLink to="/privacy">Privacy Policy</RouterLink>
    <span class="sep">&middot;</span>
    <RouterLink to="/terms">Terms &amp; Conditions</RouterLink>
    <span class="sep">&middot;</span>
    <RouterLink to="/imprint">Imprint</RouterLink>
  </nav>
</template>

<style scoped>
.legal-nav {
  margin-top: 56px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 12px;
  letter-spacing: 0.02em;
}

.legal-nav a {
  color: var(--text-dim);
  text-decoration: none;
  transition: color 0.2s ease;
}

.legal-nav a:hover {
  color: var(--accent);
}

.legal-nav a.router-link-active {
  color: var(--text-muted);
  pointer-events: none;
}

.sep {
  color: var(--text-dim);
  opacity: 0.4;
}
</style>
