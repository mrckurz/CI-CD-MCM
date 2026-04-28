# Exercise 2: Microservice Architecture, Docker & GitHub Actions

**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 24

## Learning Objectives

- Understand microservice architecture with a REST API in Go
- Containerize applications using Docker (multi-stage builds)
- Orchestrate services with Docker Compose
- Set up a basic CI pipeline with GitHub Actions

## Prerequisites

- Completed Exercise 1
- Docker Desktop installed
- Basic understanding of REST APIs

## Project Overview

The Product Catalog API has been extended with:
- **PostgreSQL storage** (`internal/store/postgres.go`) -- persistent database backend
- **Dockerfile** -- multi-stage build for minimal container image
- **docker-compose.yml** -- orchestrates API + PostgreSQL
- **GitHub Actions** (`.github/workflows/ci.yml`) -- basic CI pipeline

### Architecture

```
┌──────────────┐     ┌──────────────┐
│   Client     │────▶│   API (Go)   │
│  (curl/HTTP) │     │   Port 8080  │
└──────────────┘     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  PostgreSQL  │
                     │  Port 5432   │
                     └──────────────┘
```

### Local Development

```bash
# Run with in-memory store (no Docker needed)
go run ./cmd/api

# Run with Docker Compose (API + PostgreSQL)
docker compose up --build

# Test the API
curl http://localhost:8080/health
curl http://localhost:8080/products
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99}'
```

---

## Tasks

### Task 1: Understand the Architecture (2 Points)

1. Read the source code and understand how the API handles requests.
2. Draw a diagram (or describe in text) showing the request flow from HTTP request to database and back.
3. Explain the difference between `MemoryStore` and `PostgresStore` -- when would you use each?

**Deliverable:** Add an `ARCHITECTURE.md` file with your diagram and explanation.

---

### Task 2: Complete the GitHub Actions Workflow (6 Points)

The CI workflow (`.github/workflows/ci.yml`) has a `TODO` for a Docker build job. Your tasks:

1. **Add a `docker-build` job** that:
   - Runs after the `test` job succeeds (`needs: test`)
   - Checks out the code
   - Sets up Docker Buildx
   - Builds the Docker image with tag `product-catalog:${{ github.sha }}`
   - (Bonus) Pushes to GitHub Container Registry if on `main` branch

2. **Add a build badge** to your README showing the CI status.

**Deliverable:** Working CI pipeline (green check on your PR). Screenshot of the Actions run.

> **Heads-up: Workflow does not run? ("Get started with Actions")**
>
> If your Actions tab keeps showing the "Get started with Actions" onboarding page, two things are usually going on:
>
> 1. **GitHub lists workflows from the default branch.** If `.github/workflows/ci.yml` only exists on your feature branch (e.g. `exercise/02-microservice-docker`) and not on `main`, the Actions tab UI looks empty until the first run completes.
> 2. **The default trigger is restricted to `main`.** The shipped `ci.yml` has:
>    ```yaml
>    on:
>      push:
>        branches: [main]
>      pull_request:
>        branches: [main]
>    ```
>    A push to `exercise/02-microservice-docker` matches neither filter, so nothing fires.
>
> **Two ways to fix it — pick one:**
>
> - **Option A (matches the task — "green check on PR"):** open a Pull Request inside your own fork: base `main` ← compare `exercise/02-microservice-docker`. The existing `pull_request: branches: [main]` filter triggers, the workflow runs, and merging the PR also lands the workflow on `main` so the Actions tab is clean afterwards.
> - **Option B (loosen the trigger):** edit `.github/workflows/ci.yml` on your feature branch and change `on:` to:
>    ```yaml
>    on:
>      push:
>        branches: ["**"]
>      pull_request:
>    ```
>    Commit + push — the workflow now runs on every branch.
>
> Also make sure GitHub Actions is **enabled** on your fork (Actions tab → "I understand my workflows, go ahead and enable them"). GitHub disables Actions on forks by default.

---

### Task 3: Docker & Docker Compose (8 Points)

1. **Analyze the Dockerfile:**
   - Explain each stage of the multi-stage build. Why two stages?
   - What does `CGO_ENABLED=0` do and why is it important?
   - What is the final image size? Compare it to a single-stage build.

2. **Run the application with Docker Compose:**
   ```bash
   docker compose up --build
   ```

3. **Test all CRUD operations** using `curl` or a tool like Postman:
   - Create at least 3 products
   - List all products
   - Update a product
   - Delete a product
   - Verify the product is gone

4. **Verify data persistence:**
   - Stop and restart the containers (`docker compose down` then `up`)
   - Check if the products still exist (they should, thanks to the volume)

**Deliverable:** Document your CRUD tests and answers in `DOCKER.md`.

---

### Task 4: Add Handler Tests (8 Points)

The file `internal/handler/handler_test.go` contains a `TODO` for additional tests. Add:

1. **TestUpdateProduct** -- Create a product via POST, update it via PUT, verify the response.
2. **TestDeleteProduct** -- Create a product, delete it, verify GET returns 404.
3. **TestCreateInvalidProduct** -- POST with invalid payload (empty name), expect 400.

All tests must use `httptest.NewRecorder` (no actual HTTP server needed).

**Deliverable:** Completed test file, all tests passing (`go test -v ./internal/handler/`).

---

## API Reference

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| GET | `/health` | Health check | -- |
| GET | `/products` | List all products | -- |
| POST | `/products` | Create product | `{"name":"...","price":0.00}` |
| GET | `/products/{id}` | Get product by ID | -- |
| PUT | `/products/{id}` | Update product | `{"name":"...","price":0.00}` |
| DELETE | `/products/{id}` | Delete product | -- |

---

## Grading

| Task | Points |
|------|--------|
| Architecture Documentation | 2 |
| GitHub Actions Workflow | 6 |
| Docker & Docker Compose | 8 |
| Handler Tests | 8 |
| **Total** | **24** |

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)

