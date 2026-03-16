<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
import LegalNav from '@/components/LegalNav.vue'
</script>

<template>
  <div class="page-wrapper">
    <AppHeader />
    <main class="container fade-up">
      <h1>Privacy Policy</h1>
      <p class="last-updated">Last Updated: March 16, 2026</p>

      <p>
        This Privacy Policy explains how oglimmer.com / Oliver Zimpasser ("we", "our", "us") processes
        personal data when you use Coffee Diary (the "Platform"). It applies to all visitors and registered users.
      </p>

      <h2>1. Controller</h2>
      <p>
        Controller under Art. 4(7) GDPR:<br />
        Oliver Zimpasser, Berliner Str. 39a, 63110 Rodgau, Germany<br />
        Email: <a href="mailto:oglimmer@gmail.com">oglimmer@gmail.com</a>
      </p>

      <h2>2. Data We Process</h2>
      <p>We process the following categories of data depending on how you use the Platform:</p>
      <ul>
        <li>Account data (e.g., email address, login data).</li>
        <li>User content and metadata you upload (e.g., diary entries and associated metadata).</li>
        <li>Usage data (e.g., pages viewed, feature interactions).</li>
        <li>Technical data (e.g., IP address, browser/user-agent, device information).</li>
        <li>Identifiers stored in cookies/local storage (e.g., consent choice, visitor ID upon consent).</li>
      </ul>

      <h2>3. Purposes and Legal Bases</h2>
      <ul>
        <li>
          <strong>Service provision and account management (Art. 6(1)(b) GDPR):</strong>
          operating the Platform, authentication, transactional communications, user content display.
        </li>
        <li>
          <strong>Security, abuse prevention and troubleshooting (Art. 6(1)(f) GDPR):</strong>
          maintaining service integrity, preventing misuse, logging limited technical data such as IP address and user-agent.
        </li>
        <li>
          <strong>Usage statistics/analytics (Art. 6(1)(a) GDPR for cookies; Art. 6(1)(f) GDPR otherwise):</strong>
          we generate aggregated statistics. If you consent, we set a <code>visitor_id</code> cookie.
          If you decline, we still process minimal technical data (e.g., IP address and user-agent) on our servers
          to produce basic, aggregated statistics and to secure the service, without setting analytics cookies.
        </li>
        <li>
          <strong>Legal obligations (Art. 6(1)(c) GDPR):</strong>
          complying with statutory requirements.
        </li>
      </ul>

      <h2>4. Cookies</h2>
      <p>We use cookies only as described below:</p>
      <ul>
        <li><strong>cookie_consent</strong> (essential): stores your consent choice for up to 12 months.</li>
        <li>
          <strong>visitor_id</strong> (analytics, optional): set for about 90 days only if you accept analytics cookies;
          used to count unique visits.
        </li>
      </ul>
      <p>You can change your choice any time by clearing your browser cookies.</p>

      <h2>5. Analytics and Server Logs</h2>
      <p>
        We do not use third-party analytics. Analytics events are handled on our servers.
        If you accept analytics cookies, we use a <code>visitor_id</code> cookie.
        If you decline, we do not set analytics cookies; we may still process IP address and user-agent
        to derive a short-lived pseudonymous identifier for aggregated statistics and service security.
        You can object to processing based on legitimate interests (see Section 9).
      </p>

      <h2>6. Recipients and Transfers</h2>
      <p>
        We do not sell your data. Where we use service providers (e.g., hosting or email delivery),
        they act as processors under Art. 28 GDPR and process data only on our instructions.
        We do not transfer data outside the EU unless adequate safeguards are in place.
      </p>

      <h2>7. Retention</h2>
      <ul>
        <li><strong>Account data:</strong> retained for the life of your account; deleted upon account deletion.</li>
        <li><strong>User content:</strong> retained until you remove it or delete your account.</li>
        <li><strong>Analytics and usage data:</strong> retained for a limited period appropriate for aggregated reporting.</li>
        <li><strong>Server logs and security data:</strong> retained for a short period unless needed to investigate issues.</li>
      </ul>

      <h2>8. Security</h2>
      <p>
        We implement appropriate technical and organizational measures to protect your data.
        No online service can guarantee absolute security.
      </p>

      <h2>9. Your Rights (EU/EEA)</h2>
      <p>
        You have the rights of access, rectification, erasure, restriction, portability, and objection
        (Art. 15–21 GDPR). Where processing is based on consent, you may withdraw consent at any time
        with effect for the future. You also have the right to lodge a complaint with a supervisory
        authority in Germany.
      </p>

      <h2>10. Children</h2>
      <p>
        The Platform is not intended for children under 16. Do not create an account or upload content
        if you are under 16.
      </p>

      <h2>11. Changes to this Policy</h2>
      <p>
        We may update this Privacy Policy to reflect changes in our practices or legal requirements.
        The "Last Updated" date reflects the effective date of the current version.
      </p>

      <h2>12. Contact</h2>
      <p>
        For privacy questions or requests, contact:
        <a href="mailto:oglimmer@gmail.com">oglimmer@gmail.com</a>.
      </p>
      <LegalNav />
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

.last-updated {
  color: var(--text-dim);
  font-size: 13px;
  margin: 0 0 32px;
}

h2 {
  font-size: 18px;
  font-weight: 600;
  margin: 32px 0 12px;
  color: var(--text);
}

p, li {
  color: var(--text-muted);
  font-size: 14px;
  line-height: 1.7;
}

ul {
  padding-left: 20px;
  margin: 8px 0 16px;
}

li {
  margin-bottom: 6px;
}

a {
  color: var(--accent);
}

code {
  background: var(--surface);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
}
</style>
