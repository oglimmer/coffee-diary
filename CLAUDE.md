# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Coffee Diary is a full-stack web application for tracking espresso brewing sessions. Users log brewing parameters (temperature, grind size, weight, time), manage bean and equipment inventory, and track tasting notes with ratings. Deployed to Kubernetes at coffee.oglimmer.com.

## Build & Development Commands

### Backend (Go)
```bash
cd backend
go run ./cmd/server                       # Run dev server (see internal/config/config.go for env vars)
go test ./...                             # Run all tests
go build -o server ./cmd/server           # Build binary
```

### Frontend (Vue 3, TypeScript, Vite)
```bash
cd frontend
npm install
npm run dev                               # Dev server (proxies /api to localhost:8080)
npm run build                             # Production build
npm run type-check                        # Vue TSC type checking
npm run lint                              # Oxlint + ESLint with auto-fix
```

### Full Stack (Docker Compose)
```bash
docker compose up --build                 # MariaDB (3306) + Backend (8080) + Frontend (3000)
```

## Architecture

### Backend
- **Language:** Go with OIDC authentication (Keycloak)
- **Auth:** Session-based with OIDC. Basic auth for actuator endpoints, session auth for `/api/**`.
- **Data isolation:** All entities scoped by `user_id`.
- **Database:** MariaDB 11. SQL migrations in `backend/migrations/`.
- **Config:** Environment variables, defaults in `backend/internal/config/config.go`.

### Frontend
- **Framework:** Vue 3 Composition API + TypeScript + Pinia stores + Vue Router
- **API layer:** Native Fetch API in `frontend/src/services/`. All calls prefixed with `/api`. Automatic 401 handling redirects to login.
- **State:** `useAuthStore` (session), `useAppInfoStore` (build info). Route guards check `meta: { auth: true }` / `meta: { guest: true }`.
- **Linting:** Oxlint (fast) + ESLint. Config in `frontend/eslint.config.ts`.

### API Endpoints
- `POST/GET /api/auth/{register,login,logout,me}` — Authentication
- `GET/POST/PUT/DELETE /api/diary-entries[/{id}]` — Diary entries (paginated, filterable by coffee/sieve/date/rating)
- `GET/POST/DELETE /api/coffees[/{id}]` — Coffee beans
- `GET/POST/DELETE /api/sieves[/{id}]` — Sieves/filters
- `GET /actuator/{health,info,prometheus,metrics}` — Monitoring

### Testing
- **Backend tests:** Run with `go test ./...` from `backend/`.
- **No frontend tests** currently exist.

### Deployment
- **CI:** GitHub Actions builds Docker images on push to main, pushes to registry.oglimmer.com.
- **Production:** Kubernetes via Helm chart in `helm/coffee-diary/`. Nginx frontend proxies `/api` and `/actuator` to backend.
- **Build script:** `./oglimmer.sh` handles build/push/deploy with platform and registry options.
