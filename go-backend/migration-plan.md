# Migration Plan: Coffee Diary Spring Boot → Go

## 1. Application Overview

Coffee Diary backend is a Spring Boot 4.0 REST API for tracking espresso brewing sessions. The Go migration simplifies auth by switching to OIDC/SSO:
- **OIDC authentication** via an external identity provider (replaces custom register/login)
- MariaDB 11 with Flyway migrations
- Basic auth for actuator endpoints only
- Spring Data JPA with Specifications for dynamic query filtering

## 2. REST API Inventory

### Authentication (`/api/auth`) — OIDC-based

| Method | Path | Request Body | Response Body | Status Codes |
|--------|------|-------------|--------------|-------------|
| GET | `/api/auth/login` | none | Redirects to OIDC provider | 302 Redirect |
| GET | `/api/auth/callback` | none (query: `code`, `state`) | Sets session cookie, redirects to frontend | 302 Redirect |
| POST | `/api/auth/logout` | none | Clears session | 204 No Content |
| GET | `/api/auth/me` | none | `{"id": number, "username": string}` | 200 OK, 401 (not authenticated) |

The OIDC flow: frontend redirects to `/api/auth/login` → Go redirects to IdP → user authenticates → IdP redirects to `/api/auth/callback` → Go validates the ID token, creates/finds the user in the DB (auto-provisioning by OIDC `sub`/`preferred_username`), sets a session cookie, and redirects to the frontend.

### Coffees (`/api/coffees`) — all require authenticated session

| Method | Path | Request Body | Response Body | Status Codes |
|--------|------|-------------|--------------|-------------|
| GET | `/api/coffees` | none | `[{"id": number, "name": string}]` | 200 OK |
| POST | `/api/coffees` | `{"name": string}` | `{"id": number, "name": string}` | 201 Created, 400 (validation) |
| DELETE | `/api/coffees/{id}` | none | none | 204 No Content, 400 (not found), 403 (wrong user) |

### Sieves (`/api/sieves`) — all require authenticated session

| Method | Path | Request Body | Response Body | Status Codes |
|--------|------|-------------|--------------|-------------|
| GET | `/api/sieves` | none | `[{"id": number, "name": string}]` | 200 OK |
| POST | `/api/sieves` | `{"name": string}` | `{"id": number, "name": string}` | 201 Created, 400 (validation) |
| DELETE | `/api/sieves/{id}` | none | none | 204 No Content, 400 (not found), 403 (wrong user) |

### Diary Entries (`/api/diary-entries`) — all require authenticated session

| Method | Path | Query Params | Request Body | Response Body | Status Codes |
|--------|------|-------------|-------------|--------------|-------------|
| GET | `/api/diary-entries` | `coffeeId`, `sieveId`, `dateFrom`, `dateTo`, `ratingMin`, `page` (0-based, default 0), `size` (default 20), `sort` (default `dateTime,asc`) | none | `PageResponse<DiaryEntryResponse>` | 200 OK |
| GET | `/api/diary-entries/{id}` | none | none | `DiaryEntryResponse` | 200 OK, 400 (not found), 403 (wrong user) |
| POST | `/api/diary-entries` | none | `DiaryEntryRequest` | `DiaryEntryResponse` | 201 Created, 400 (validation) |
| PUT | `/api/diary-entries/{id}` | none | `DiaryEntryRequest` | `DiaryEntryResponse` | 200 OK, 400 (not found/validation), 403 (wrong user) |
| DELETE | `/api/diary-entries/{id}` | none | none | 204 No Content, 400 (not found), 403 (wrong user) |

**DiaryEntryRequest:**
```json
{
  "dateTime": "2024-01-01T10:00:00",  // required
  "sieveId": 1,                        // optional
  "temperature": 93,                   // optional, default 93
  "coffeeId": 1,                       // optional
  "grindSize": 5.0,                    // optional
  "inputWeight": 18.0,                 // optional
  "outputWeight": 36.0,               // optional
  "timeSeconds": 25,                   // optional
  "rating": 4,                         // optional, 1-5
  "notes": "string"                    // optional
}
```

**DiaryEntryResponse:**
```json
{
  "id": 1, "userId": 1, "dateTime": "2024-01-01T10:00:00",
  "sieveId": 1, "sieveName": "IMS", "temperature": 93,
  "coffeeId": 1, "coffeeName": "Ethiopian", "grindSize": 5.0,
  "inputWeight": 18.0, "outputWeight": 36.0, "timeSeconds": 25,
  "rating": 4, "notes": "notes"
}
```

**PageResponse:**
```json
{
  "content": [...],
  "totalElements": 100,
  "totalPages": 5,
  "number": 0,
  "size": 20
}
```

### Actuator endpoints

| Method | Path | Auth | Status |
|--------|------|------|--------|
| GET | `/actuator/health` | none | 200 |
| GET | `/actuator/info` | none | 200 |
| GET | `/actuator/prometheus` | basic auth (ACTUATOR_USERNAME/ACTUATOR_PASSWORD) | 200 |
| GET | `/actuator/metrics` | basic auth | 200 |

### Error Response Shape

All error responses follow this shape:
```json
{"status": 400, "error": "Bad Request", "message": "..."}
```

Validation errors add an `errors` field:
```json
{"status": 400, "error": "Validation Failed", "message": "Invalid request body", "errors": {"field": "message"}}
```

## 3. Authentication & Authorization Strategy

- **OIDC/SSO** — delegates authentication to an external identity provider
- Go handles the OIDC Authorization Code flow: login redirect, callback token exchange, ID token validation
- After successful OIDC callback, a lightweight session cookie is set containing the user's DB ID
- `gorilla/sessions` cookie store for session management — only stores user ID, minimal data
- User auto-provisioning: on first OIDC login, a user record is created in the DB using the OIDC `sub` claim as identifier and `preferred_username` (or `email`) as display name
- `/api/auth/login`, `/api/auth/callback` are public; all other `/api/**` require authenticated session
- No custom password handling — passwords are managed entirely by the IdP
- CSRF is disabled
- Security headers: X-Content-Type-Options, X-Frame-Options: DENY, HSTS, CSP (handled by reverse proxy in prod, but set in Go for parity)

## 4. Business Logic Components

### AuthService
- `handleCallback(idToken)` → find or create user by OIDC `sub` claim, return user ID for session
- No registration or password logic — handled by the IdP

### CoffeeService
- `findAllByUser(userId)` → list coffees for user
- `create(userId, name)` → create coffee, return CoffeeResponse
- `delete(userId, coffeeId)` → verify ownership, delete

### SieveService
- Same pattern as CoffeeService

### DiaryEntryService
- `findAll(userId, filters, page, size, sort)` → dynamic query with filters, return PageResponse
- `findById(userId, entryId)` → verify ownership
- `create(userId, request)` → validate sieve/coffee ownership, default temp=93
- `update(userId, entryId, request)` → verify ownership, validate sieve/coffee, clear associations if null
- `delete(userId, entryId)` → verify ownership

## 5. Database Layer

**Database:** MariaDB 11, driver: `github.com/go-sql-driver/mysql`

**Schema** (after all migrations):
- `users` (id BIGINT PK AUTO_INCREMENT, username VARCHAR UNIQUE, oidc_sub VARCHAR UNIQUE, created_at TIMESTAMP) — `password` column dropped, `oidc_sub` added via new migration
- `sieves` (id BIGINT PK AUTO_INCREMENT, name VARCHAR, user_id BIGINT FK→users)
- `coffees` (id BIGINT PK AUTO_INCREMENT, name VARCHAR, user_id BIGINT FK→users)
- `diary_entries` (id BIGINT PK AUTO_INCREMENT, user_id BIGINT FK, date_time DATETIME, sieve_id BIGINT FK nullable, temperature INT default 93, coffee_id BIGINT FK nullable, grind_size DOUBLE, input_weight DOUBLE, output_weight DOUBLE, time_seconds INT, rating INT, notes TEXT)

**Migrations:** Copy V1-V4 SQL files. Use `golang-migrate/migrate` with mysql driver.

**Pagination:** Spring uses 0-based page numbers. The Go implementation must preserve this (page=0 is first page).

**Sorting:** Default sort for diary entries is `date_time ASC`. Support `sort` query param in format `field,direction`.

## 6. Static & Template Asset Hosting

No static assets or templates served by the backend. Thymeleaf dependency exists but `check-template-location: false`. Frontend is a separate Vue app. **No action needed.**

## 7. Environment Variable Configuration Map

| Spring property | Go env var | Default | Purpose |
|----------------|-----------|---------|---------|
The Go app reads env vars that align with the existing Helm sealed-secret and deployment config:

| Env var | Default | Source in K8s | Purpose |
|---------|---------|---------------|---------|
| `DB_HOST` | `localhost` | configmap (`database.host`) | MariaDB host |
| `DB_PORT` | `3306` | configmap (`database.port`) | MariaDB port |
| `DB_NAME` | `coffeediary` | configmap (`database.name`) | Database name |
| `DB_USER` | `app` | sealed-secret `DB_USER` | Database username |
| `DB_PASSWORD` | `app` | sealed-secret `DB_PASSWORD` | Database password |
| `SERVER_PORT` | `8080` | — | HTTP listen port |
| `ACTUATOR_USERNAME` | `actuator` | sealed-secret `ACTUATOR_USERNAME` | Basic auth user for metrics |
| `ACTUATOR_PASSWORD` | `changeme` | sealed-secret `ACTUATOR_PASSWORD` | Basic auth pass for metrics |
| `APP_NAME` | `coffee-diary-backend` | — | App info |
| `APP_VERSION` | `0.0.1-SNAPSHOT` | — | App version |
| `SESSION_SECRET` | (required) | sealed-secret (new) | Cookie session encryption key |
| `OIDC_ISSUER_URL` | `https://id.oglimmer.de/realms/oglimmer/` | configmap or env | Keycloak OIDC issuer URL |
| `OIDC_CLIENT_ID` | `coffee-diary` | configmap or env | OIDC client ID |
| `OIDC_CLIENT_SECRET` | (required) | sealed-secret (new) | OIDC client secret — never commit |
| `OIDC_REDIRECT_URL` | `https://coffee.oglimmer.com/api/auth/callback` | configmap or env | OIDC callback URL |
| `FRONTEND_URL` | `https://coffee.oglimmer.com` | configmap or env | Where to redirect after OIDC login |

**Helm changes needed:** The sealed-secret needs `SESSION_SECRET` and `OIDC_CLIENT_SECRET` added. The deployment template needs updated env vars (remove `SPRING_*`, `JAVA_OPTS`; add OIDC + session vars). Backend resource limits can be reduced significantly (no JVM).

## 8. Go Project Layout

```
go-backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   ├── auth.go          # OIDC login redirect, callback, logout, /me
│   │   ├── coffee.go
│   │   ├── sieve.go
│   │   ├── diary_entry.go
│   │   ├── actuator.go
│   │   └── middleware.go    # Session auth middleware, security headers
│   ├── service/
│   │   ├── auth.go          # OIDC user provisioning (find-or-create)
│   │   ├── coffee.go
│   │   ├── sieve.go
│   │   └── diary_entry.go
│   ├── repository/
│   │   ├── user.go
│   │   ├── coffee.go
│   │   ├── sieve.go
│   │   └── diary_entry.go
│   ├── domain/
│   │   └── model.go
│   └── errors/
│       └── errors.go
├── migrations/
│   ├── 1_init.up.sql
│   ├── 1_init.down.sql
│   ├── 2_grind_size_to_double.up.sql
│   ├── 2_grind_size_to_double.down.sql
│   ├── 3_spring_session_tables.up.sql
│   ├── 3_spring_session_tables.down.sql
│   ├── 4_drop_spring_session_tables.up.sql
│   └── 4_drop_spring_session_tables.down.sql
├── Dockerfile
├── go.mod
└── go.sum
```

## 9. Dependencies (Go modules)

| Module | Purpose |
|--------|---------|
| `github.com/go-chi/chi/v5` | HTTP router |
| `github.com/go-sql-driver/mysql` | MariaDB driver |
| `github.com/golang-migrate/migrate/v4` | DB migrations |
| `github.com/gorilla/sessions` | Session management (cookie store) |
| `github.com/coreos/go-oidc/v3` | OIDC ID token verification |
| `golang.org/x/oauth2` | OAuth2 Authorization Code flow |
| `github.com/stretchr/testify` | Test assertions |
| `github.com/prometheus/client_golang` | Prometheus metrics |

## 10. Testing Strategy

- **Unit tests:** Mock repository interfaces, test service logic (ownership checks, validation, defaults)
- **Integration tests:** Use `testcontainers-go` with MariaDB to test full HTTP flow (mock OIDC session → CRUD → cross-user isolation)
- **HTTP tests:** `httptest.NewRecorder` + router for handler-level testing (inject session directly for auth, no need to mock full OIDC flow in unit tests)
- Mirror existing test coverage: user provisioning, coffee, sieve, diary entry services + full API integration test

## 11. Dockerfile Plan

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
USER appuser
EXPOSE 8080
CMD ["./server"]
```

## 12. Known Risks & Edge Cases

1. **Auth migration:** Switching from custom username/password to OIDC means all existing users must re-authenticate via the IdP. A new migration (V5) adds `oidc_sub` column and drops `password`. Existing user rows will need their `oidc_sub` populated on first OIDC login (matched by username).
2. **Pagination sort parsing:** Spring's `sort=dateTime,desc` format needs custom parsing in Go. Must map `dateTime` → `date_time` column name.
3. **DateTime serialization:** Spring Boot serializes `LocalDateTime` as `"2024-01-01T10:00:00"` (ISO without timezone). Go's `time.Time` includes timezone by default. Must use custom JSON marshaling to match.
4. **Null handling:** Java `Long`/`Integer` nullable fields become `*int64`/`*int` pointers in Go. JSON serialization must output `null` not `0`.
5. **Database `SPRING_SESSION` tables:** V3 creates them, V4 drops them. These are no-ops for Go but must be included for migration compatibility with existing databases.
6. **OIDC provider dependency:** The Go app requires a running OIDC provider (e.g. Keycloak, Auth0, Authentik). This must be set up before deployment. For local dev, a test/mock OIDC provider can be used.
