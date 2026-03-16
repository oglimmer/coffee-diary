# coffee diary

A web application for tracking espresso brewing sessions. Record your beans, equipment, brewing parameters, and tasting notes to refine your workflow over time.

## Features

- **Diary entries** -- log temperature, grind size, input/output weight, brew time, rating (1-5 stars), and notes
- **Coffee & sieve management** -- maintain a personal inventory of beans and filters
- **Per-user isolation** -- each account has its own data
- **Rate limiting** -- per-IP throttling (10 req/s) via Guava RateLimiter

## Tech stack

| Layer    | Technology                                      |
|----------|------------------------------------------------|
| Backend  | Spring Boot 3.4, Java 21, Spring Security, JPA |
| Frontend | Vue 3, TypeScript, Vite, Pinia, Vue Router     |
| Database | MariaDB 11, Flyway migrations                  |
| Infra    | Docker, Helm, GitHub Actions, Nginx             |

## Prerequisites

- Java 21
- Node.js 20.19+ or 22.12+
- MariaDB 11
- Docker (optional, for containerised setup)

## Running locally

### With Docker Compose

```sh
docker compose up --build
```

This starts MariaDB, the backend (port 8080), and the frontend (port 3000).

### Without Docker

**Backend**

```sh
cd backend
export DB_USERNAME=root DB_PASSWORD=secret
./mvnw spring-boot:run
```

**Frontend**

```sh
cd frontend
npm install
npm run dev
```

## Deployment

The CI pipeline (`.github/workflows/build.yml`) builds and pushes Docker images on every push to `main`. A Helm chart under `helm/coffee-diary/` deploys to Kubernetes with:

- Ingress at `coffee.oglimmer.com` (TLS via cert-manager)
- Path-based routing: `/api`, `/actuator` to backend; `/` to frontend
- SealedSecret for database credentials
- Prometheus metrics at `/actuator/prometheus`

## Project structure

```
backend/     Spring Boot REST API
frontend/    Vue 3 SPA
helm/        Kubernetes Helm chart
```

## License

Private.
