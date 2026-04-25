# Continuous Delivery in Agile Software Development -- Exercises

This repository contains four progressive exercises for the Master course **Continuous Delivery in Agile Software Development**.

## Overview

| Exercise | Topic | Branch |
|----------|-------|--------|
| 1 | Git Basics: PRs, Interactive Rebase, Unit Tests | `exercise/01-git-basics` |
| 2 | Microservice Architecture, Docker & GitHub Actions | `exercise/02-microservice-docker` |
| 3 | CI Pipeline: SonarCloud, Matrix Builds, Linting | `exercise/03-ci-pipeline` |
| 4 | Vulnerability Scanning & Kubernetes Deployment | `exercise/04-security-k8s` |

## Technology Stack

- **Language:** Go 1.24+
- **Web Framework:** Gorilla Mux
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **CI/CD:** GitHub Actions
- **Code Quality:** SonarCloud, golangci-lint
- **Security:** Trivy, govulncheck
- **Deployment:** Kubernetes (Minikube)

## Project: Product Catalog API

Throughout the four exercises you will build and evolve a **Product Catalog API** -- a RESTful web service for managing products (create, read, update, delete). The API is written in Go and grows in complexity with each exercise.

### What the Application Does

The Product Catalog API exposes the following HTTP endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/products` | List all products |
| POST | `/products` | Create a new product |
| GET | `/products/{id}` | Get a product by ID |
| PUT | `/products/{id}` | Update a product |
| DELETE | `/products/{id}` | Delete a product |

A product has three fields: `id`, `name`, and `price`.

### Project Structure

```
cmd/api/main.go                # Application entry point -- starts the HTTP server
internal/
  model/product.go             # Product data model and validation
  store/
    memory.go                  # In-memory store (Exercise 1-2)
    postgres.go                # PostgreSQL store (from Exercise 2)
  handler/handler.go           # HTTP request handlers (routing, JSON encoding)
Dockerfile                     # Multi-stage Docker build (from Exercise 2)
docker-compose.yml             # Orchestrates API + PostgreSQL (from Exercise 2)
.github/workflows/ci.yml       # CI/CD pipeline (from Exercise 2, extended in 3-4)
k8s/                           # Kubernetes manifests (Exercise 4)
```

### What You Build in Each Exercise

| Exercise | What You Do |
|----------|-------------|
| **1 -- Git Basics** | Fork the repo, write unit tests for the in-memory store, create your first Pull Request, and practice interactive rebase to clean up commit history. |
| **2 -- Microservice & Docker** | Understand the microservice architecture, complete a GitHub Actions CI pipeline with a Docker build job, analyze the Dockerfile and Docker Compose setup, and add HTTP handler tests. |
| **3 -- CI Pipeline** | Extend the pipeline with matrix builds (multiple Go versions and OS), integrate golangci-lint for code quality, set up SonarCloud for static analysis, and improve test coverage to ≥ 80%. |
| **4 -- Security & K8s** | Scan the Docker image with Trivy, scan Go dependencies with govulncheck, deploy the application to a local Kubernetes cluster (Minikube), and configure production-readiness features (probes, resource limits). |

By the end of the course, you will have a fully containerized Go microservice with a complete CI/CD pipeline including automated testing, linting, security scanning, and Kubernetes deployment.

## Prerequisites

- Go 1.24+ installed
- Git 2.30+
- GitHub Account
- Docker Desktop (from Exercise 2)
- Minikube (Exercise 4)

## Getting Started

1. **Fork** this repository on GitHub (click the "Fork" button in the top right corner). **Uncheck** "Copy the `main` branch only" so that all exercise branches are included in your fork.
2. **Clone** your fork:

```bash
git clone https://github.com/<your-username>/CI-CD-MCM.git
cd CI-CD-MCM
```

3. Switch to the respective exercise branch:

```bash
git checkout exercise/01-git-basics
```

> **Important:** Do not clone the original repository directly — always work on your own fork so you can push changes and create Pull Requests.

Each exercise branch contains a detailed `README.md` with instructions.

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)
