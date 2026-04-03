<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
</script>

<template>
  <div class="page-wrapper">
    <AppHeader />
    <main class="container fade-up">
      <h1>Imprint / Impressum</h1>
      <p class="subtitle">Legal information and contact details</p>

      <h2>For Germany: Angaben gemäß § 5 TMG</h2>

      <h3>Diensteanbieter / Website Provider</h3>
      <p>
        Diese Webseite wird bereitgestellt von / This website is provided by:<br />
        Oliver Zimpasser
      </p>
      <p>
        Berliner Str. 39a<br />
        63110 Rodgau<br />
        Germany
      </p>

      <h3>Kontakt / Contact</h3>
      <p>
        E-Mail: <a href="mailto:oglimmer@gmail.com">oglimmer@gmail.com</a><br />
        <span class="note">
          (Angaben zur schnellen elektronischen Kontaktaufnahme und unmittelbaren Kommunikation
          gem. § 5 Abs. 1 Nr. 2 TMG)
        </span>
      </p>

      <h3>Verantwortlich i.S.d. § 18 Abs. 2 MStV</h3>
      <p>
        Verantwortlich für den Inhalt / Responsible for content:<br />
        Oliver Zimpasser, Berliner Str. 39a, 63110 Rodgau, Germany
      </p>

      <h3>EU-Streitschlichtung / Dispute Resolution</h3>
      <p>
        Die Europäische Kommission stellt eine Plattform zur Online-Streitbeilegung (OS) bereit:
        <a href="https://ec.europa.eu/consumers/odr" target="_blank" rel="noopener noreferrer">
          https://ec.europa.eu/consumers/odr</a>.
        Wir sind nicht verpflichtet und grundsätzlich nicht bereit, an Streitbeilegungsverfahren
        vor einer Verbraucherschlichtungsstelle teilzunehmen.
      </p>

      <h3>Service Information</h3>
      <p>All services are provided free of charge.</p>
    </main>
  </div>
</template>

<style scoped>
.page-wrapper {
  min-height: 100vh;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 48px 24px 80px;
}

h1 {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 400;
  margin: 0 0 8px;
  letter-spacing: -0.02em;
  font-variation-settings: 'SOFT' 100, 'WONK' 1;
}

.subtitle {
  color: var(--text-dim);
  font-size: 14px;
  margin: 0 0 32px;
}

h2 {
  font-size: 18px;
  font-weight: 600;
  margin: 32px 0 16px;
  color: var(--text);
}

h3 {
  font-size: 15px;
  font-weight: 600;
  margin: 24px 0 8px;
  color: var(--text);
}

p {
  color: var(--text-muted);
  font-size: 14px;
  line-height: 1.7;
}

p + p {
  margin-top: 8px;
}

a {
  color: var(--accent);
}

.note {
  font-size: 12px;
  color: var(--text-dim);
}
</style>
