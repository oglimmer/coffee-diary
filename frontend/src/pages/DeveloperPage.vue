<script setup lang="ts">
import AppHeader from '@/components/AppHeader.vue'
</script>

<template>
  <div class="page-wrapper">
    <AppHeader />
    <main class="container fade-up">
      <h1>Developer Documentation</h1>
      <p class="subtitle">REST API reference for Coffee Diary</p>

      <section>
        <h2>Overview</h2>
        <p>
          Coffee Diary exposes a JSON REST API under the <code>/api</code> prefix. All protected
          endpoints require an authenticated session (via OIDC). Unauthenticated requests to
          protected endpoints return <code>401</code>. All request and response bodies use
          <code>application/json</code>.
        </p>
      </section>

      <section>
        <h2>Authentication</h2>
        <p>
          Authentication is handled via OpenID Connect (Keycloak). The browser flow redirects the
          user to the identity provider and establishes a server-side session on callback.
        </p>
        <table>
          <thead>
            <tr>
              <th>Method</th>
              <th>Path</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/auth/login</code></td>
              <td>Redirects to the OIDC provider for login</td>
            </tr>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/auth/callback</code></td>
              <td>OIDC callback &mdash; exchanges code for session</td>
            </tr>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/auth/logout</code></td>
              <td>Destroys the session and redirects to the OIDC end-session endpoint</td>
            </tr>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/auth/me</code></td>
              <td>Returns the current authenticated user</td>
            </tr>
          </tbody>
        </table>

        <h3>GET /api/auth/me</h3>
        <pre><code>{
  "id": 1,
  "username": "oliver"
}</code></pre>
      </section>

      <section>
        <h2>Diary Entries</h2>
        <p>Core resource for espresso brewing logs. All endpoints require authentication. Entries are scoped to the authenticated user.</p>

        <table>
          <thead>
            <tr>
              <th>Method</th>
              <th>Path</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/diary-entries</code></td>
              <td>List entries (paginated, filterable)</td>
            </tr>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/diary-entries/{id}</code></td>
              <td>Get a single entry</td>
            </tr>
            <tr>
              <td><span class="method post">POST</span></td>
              <td><code>/api/diary-entries</code></td>
              <td>Create a new entry</td>
            </tr>
            <tr>
              <td><span class="method put">PUT</span></td>
              <td><code>/api/diary-entries/{id}</code></td>
              <td>Update an entry</td>
            </tr>
            <tr>
              <td><span class="method delete">DELETE</span></td>
              <td><code>/api/diary-entries/{id}</code></td>
              <td>Delete an entry</td>
            </tr>
          </tbody>
        </table>

        <h3>Query Parameters (GET /api/diary-entries)</h3>
        <table>
          <thead>
            <tr>
              <th>Param</th>
              <th>Type</th>
              <th>Default</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>page</code></td>
              <td>int</td>
              <td>0</td>
              <td>Page number (zero-based)</td>
            </tr>
            <tr>
              <td><code>size</code></td>
              <td>int</td>
              <td>20</td>
              <td>Page size</td>
            </tr>
            <tr>
              <td><code>sort</code></td>
              <td>string</td>
              <td>dateTime,asc</td>
              <td>Sort field and direction, e.g. <code>dateTime,desc</code></td>
            </tr>
            <tr>
              <td><code>coffeeId</code></td>
              <td>int</td>
              <td>&mdash;</td>
              <td>Filter by coffee ID</td>
            </tr>
            <tr>
              <td><code>sieveId</code></td>
              <td>int</td>
              <td>&mdash;</td>
              <td>Filter by sieve ID</td>
            </tr>
            <tr>
              <td><code>dateFrom</code></td>
              <td>string</td>
              <td>&mdash;</td>
              <td>ISO datetime, e.g. <code>2025-01-01T00:00:00</code></td>
            </tr>
            <tr>
              <td><code>dateTo</code></td>
              <td>string</td>
              <td>&mdash;</td>
              <td>ISO datetime upper bound</td>
            </tr>
            <tr>
              <td><code>ratingMin</code></td>
              <td>int</td>
              <td>&mdash;</td>
              <td>Minimum rating (1&ndash;5)</td>
            </tr>
          </tbody>
        </table>

        <h3>Paginated Response</h3>
        <pre><code>{
  "content": [ /* array of diary entries */ ],
  "totalElements": 42,
  "totalPages": 3,
  "number": 0,
  "size": 20
}</code></pre>

        <h3>Diary Entry Object</h3>
        <pre><code>{
  "id": 1,
  "userId": 1,
  "dateTime": "2025-06-15T08:30:00",
  "coffeeId": 3,
  "coffeeName": "Ethiopia Yirgacheffe",
  "sieveId": 2,
  "sieveName": "IMS Nanotech",
  "temperature": 93,
  "grindSize": 2.5,
  "inputWeight": 18.0,
  "outputWeight": 36.0,
  "timeSeconds": 28,
  "rating": 4,
  "notes": "Bright acidity, clean finish"
}</code></pre>

        <h3>Create / Update Request Body</h3>
        <pre><code>{
  "dateTime": "2025-06-15T08:30:00",
  "coffeeId": 3,
  "sieveId": 2,
  "temperature": 93,
  "grindSize": 2.5,
  "inputWeight": 18.0,
  "outputWeight": 36.0,
  "timeSeconds": 28,
  "rating": 4,
  "notes": "Bright acidity, clean finish"
}</code></pre>
        <p class="note">All fields except <code>dateTime</code> are optional and nullable.</p>
      </section>

      <section>
        <h2>Coffees</h2>
        <p>Manage the coffee bean inventory. Scoped to the authenticated user.</p>

        <table>
          <thead>
            <tr>
              <th>Method</th>
              <th>Path</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/coffees</code></td>
              <td>List all coffees</td>
            </tr>
            <tr>
              <td><span class="method post">POST</span></td>
              <td><code>/api/coffees</code></td>
              <td>Create a coffee</td>
            </tr>
            <tr>
              <td><span class="method delete">DELETE</span></td>
              <td><code>/api/coffees/{id}</code></td>
              <td>Delete a coffee</td>
            </tr>
          </tbody>
        </table>

        <h3>Coffee Object</h3>
        <pre><code>{
  "id": 3,
  "name": "Ethiopia Yirgacheffe"
}</code></pre>

        <h3>Create Request Body</h3>
        <pre><code>{
  "name": "Ethiopia Yirgacheffe"
}</code></pre>
      </section>

      <section>
        <h2>Sieves</h2>
        <p>Manage sieve/filter basket inventory. Scoped to the authenticated user.</p>

        <table>
          <thead>
            <tr>
              <th>Method</th>
              <th>Path</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/api/sieves</code></td>
              <td>List all sieves</td>
            </tr>
            <tr>
              <td><span class="method post">POST</span></td>
              <td><code>/api/sieves</code></td>
              <td>Create a sieve</td>
            </tr>
            <tr>
              <td><span class="method delete">DELETE</span></td>
              <td><code>/api/sieves/{id}</code></td>
              <td>Delete a sieve</td>
            </tr>
          </tbody>
        </table>

        <h3>Sieve Object</h3>
        <pre><code>{
  "id": 2,
  "name": "IMS Nanotech"
}</code></pre>

        <h3>Create Request Body</h3>
        <pre><code>{
  "name": "IMS Nanotech"
}</code></pre>
      </section>

      <section>
        <h2>Actuator / Monitoring</h2>
        <p>Health and metrics endpoints. These are public and do not require authentication.</p>

        <table>
          <thead>
            <tr>
              <th>Method</th>
              <th>Path</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/actuator/health</code></td>
              <td>Health check</td>
            </tr>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/actuator/info</code></td>
              <td>Build and version info</td>
            </tr>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/actuator/prometheus</code></td>
              <td>Prometheus-formatted metrics</td>
            </tr>
            <tr>
              <td><span class="method get">GET</span></td>
              <td><code>/actuator/metrics</code></td>
              <td>JSON metrics summary</td>
            </tr>
          </tbody>
        </table>
      </section>

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

section {
  margin-bottom: 40px;
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

.note {
  font-size: 12px;
  color: var(--text-dim);
  margin-top: 8px;
}

code {
  font-family: var(--font-mono, 'SF Mono', 'Fira Code', 'Fira Mono', monospace);
  font-size: 13px;
  background: var(--surface-raised, rgba(255, 255, 255, 0.04));
  padding: 2px 6px;
  border-radius: 4px;
}

pre {
  background: var(--surface-raised, rgba(255, 255, 255, 0.04));
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  margin: 12px 0;
}

pre code {
  background: none;
  padding: 0;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-muted);
}

table {
  width: 100%;
  border-collapse: collapse;
  margin: 12px 0;
  font-size: 14px;
}

th {
  text-align: left;
  font-weight: 600;
  color: var(--text);
  padding: 10px 12px;
  border-bottom: 2px solid var(--border);
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

td {
  padding: 10px 12px;
  border-bottom: 1px solid var(--border);
  color: var(--text-muted);
  vertical-align: top;
}

tr:last-child td {
  border-bottom: none;
}

.method {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.04em;
  padding: 2px 8px;
  border-radius: 4px;
  font-family: var(--font-mono, monospace);
}

.method.get {
  background: rgba(74, 222, 128, 0.12);
  color: #4ade80;
}

.method.post {
  background: rgba(96, 165, 250, 0.12);
  color: #60a5fa;
}

.method.put {
  background: rgba(251, 191, 36, 0.12);
  color: #fbbf24;
}

.method.delete {
  background: rgba(248, 113, 113, 0.12);
  color: #f87171;
}
</style>
