# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Coffee Diary is a full-stack web application for tracking espresso brewing sessions. Users log brewing parameters (temperature, grind size, weight, time), manage bean and equipment inventory, and track tasting notes with ratings. Deployed to Kubernetes at coffee.oglimmer.com.

## Build & Development Commands

### Backend (Spring Boot 3.4, Java 21, Maven)
```bash
cd backend
./mvnw spring-boot:run                    # Run dev server (needs DB_USERNAME, DB_PASSWORD env vars)
./mvnw test                               # Unit tests (excludes *IT.java)
./mvnw verify                             # Unit + integration tests (uses TestContainers/MariaDB)
./mvnw package -DskipTests                # Build JAR without tests
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
- **Framework:** Spring Boot 3.4 with Spring Security, Spring Data JPA, Flyway
- **Auth:** Session-based (HTTP session, not JWT). CSRF disabled. Two security filter chains: basic auth for actuator endpoints, session auth for `/api/**`.
- **Data isolation:** All entities scoped by `user_id`. Service methods enforce ownership via `userDetails.getId()`.
- **Database:** MariaDB 11. Flyway migrations in `backend/src/main/resources/db/migration/`.
- **Code style:** Lombok for boilerplate, SLF4J logging.
- **Rate limiting:** Per-IP via Guava RateLimiter filter (configurable, default 10 req/s).

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
- **Backend unit tests:** Mockito-based, in `backend/src/test/java/`. Run with `./mvnw test`.
- **Backend integration tests:** `ApiIT.java` using TestContainers (MariaDB 11). Run with `./mvnw verify`.
- **No frontend tests** currently exist.

### Deployment
- **CI:** GitHub Actions builds Docker images on push to main, pushes to registry.oglimmer.com.
- **Production:** Kubernetes via Helm chart in `helm/coffee-diary/`. Nginx frontend proxies `/api` and `/actuator` to backend.
- **Build script:** `./oglimmer.sh` handles build/push/deploy with platform and registry options.
