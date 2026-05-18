# Exercise 4: Vulnerability Scanning & Kubernetes Deployment

**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 24

## Learning Objectives

- Integrate vulnerability scanning into the CI/CD pipeline
- Scan Docker images and Go dependencies for known vulnerabilities
- Deploy a multi-tier application to Kubernetes (Minikube)
- Understand Kubernetes concepts: Deployments, Services, Secrets, Probes

## Prerequisites

- Familiarity with CI pipelines (covered in Exercise 3) — this branch is **self-contained**, so you can work on it without having finished Exercise 3 first.
- Docker Desktop installed
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) installed
- [kubectl](https://kubernetes.io/docs/tasks/tools/) installed
- [Trivy](https://aquasecurity.github.io/trivy/) installed (optional, for local scanning)

> **Note on `ci.yml`:** The pipeline contains two leftover TODOs from Exercise 3 (`Add lint job`, `Add sonarcloud job`). Those are **not graded in Exercise 4** — leave them as-is or solve them for your own benefit. Exercise 4 only requires the two new jobs marked `vulnerability-scan` and `trivy-scan` (see Task 1 and Task 2).

## What's New in This Exercise

- **Kubernetes manifests** (`k8s/`) -- Deployment, Service, Secret, PVC
- **Trivy scanning** -- container image vulnerability scanning
- **Dependency scanning** -- Go module vulnerability checks
- **Complete CD pipeline** -- from code to running in Kubernetes

---

## Tasks

### Task 1: Vulnerability Scanning -- Docker Image (6 Points)

1. **Build the Docker image locally:**
   ```bash
   docker build -t product-catalog:latest .
   ```

2. **Scan the image with Trivy:**
   ```bash
   trivy image product-catalog:latest
   ```

3. **Analyze the results:**
   - How many vulnerabilities were found? Categorize by severity (CRITICAL, HIGH, MEDIUM, LOW).
   - Which base image contributes the most vulnerabilities?
   - Can you reduce vulnerabilities by changing the base image? Try switching to `scratch` or `distroless`.

4. **Add a Trivy scan job to the CI pipeline** (see the TODO in `ci.yml`) that:
   - Runs after the `docker-build` job
   - Scans the built Docker image using `aquasecurity/trivy-action@master`
   - Fails the build if CRITICAL or HIGH vulnerabilities are found
   - Outputs results in `table` format

   > **Hint:** The Trivy action needs `image-ref`, `format`, `exit-code`, and `severity` parameters.

5. **Upload the Trivy scan results as a build artifact:**
   - Generate a JSON report (use `format: 'json'` and `output` parameter)
   - Upload it using `actions/upload-artifact@v4`
   - Use `if: always()` so the report is uploaded even if the scan finds vulnerabilities

**Deliverable:** Trivy scan output (before and after base image optimization). Updated CI workflow. Trivy JSON report downloadable as artifact from the Actions run.

---

### Task 2: Vulnerability Scanning -- Dependencies (4 Points)

1. **Scan Go dependencies:**
   ```bash
   # Using govulncheck (official Go vulnerability checker)
   go install golang.org/x/vuln/cmd/govulncheck@latest
   govulncheck ./...
   ```

2. **Add a `vulnerability-scan` job to the CI pipeline** (see the TODO in `ci.yml`) that:
   - Runs after the `test` job
   - Installs `govulncheck` and runs it against the codebase
   - Fails if known vulnerabilities are found

   > **Hint:** Use `go install golang.org/x/vuln/cmd/govulncheck@latest` to install the tool.

3. **If vulnerabilities are found:**
   - Update the affected dependencies (`go get -u <module>`)
   - Document the CVEs and how you resolved them

**Deliverable:** govulncheck output. Updated `go.mod` if changes were needed.

---

### Task 3: Kubernetes Deployment with Minikube (8 Points)

1. **Start Minikube:**
   ```bash
   minikube start
   ```

2. **Build the image and load it into Minikube:**
   ```bash
   docker build -t product-catalog:latest .
   minikube image load product-catalog:latest
   ```

   > **Why not `eval $(minikube docker-env)`?** Modern Docker (24+) uses BuildKit/Buildx, which builds on the host's buildx builder and ignores the docker-env switch. `minikube image load` works reliably regardless of your Docker version.

3. **Deploy the application:**
   ```bash
   kubectl apply -f k8s/
   ```

4. **Verify the deployment:**
   ```bash
   kubectl get all -n product-catalog
   kubectl logs deployment/product-catalog-api -n product-catalog
   ```

   > **Expect a few API restarts at first.** The API pods start before PostgreSQL is fully ready and exit when they can't connect to the database. Kubernetes restarts them; after ~1-2 minutes both pods should reach `1/1 Running`. If they're stuck in `CrashLoopBackOff` after 3 minutes, run `kubectl logs <pod> -n product-catalog --previous` to investigate.

5. **Access the API.** Easiest cross-platform option: `kubectl port-forward` in a second terminal.
   ```bash
   kubectl port-forward svc/product-catalog-api 8080:8080 -n product-catalog
   # In a third terminal:
   curl http://localhost:8080/health
   curl http://localhost:8080/products
   ```

   Alternative: `minikube service product-catalog-api -n product-catalog --url`. On macOS with the Docker driver this command **holds the terminal open** as it maintains a tunnel — use a second terminal for `curl`, or stick with `port-forward`.

6. **Test CRUD operations** against the Kubernetes-deployed API.

**Deliverable:** Screenshots of:
- `kubectl get all -n product-catalog` output
- Successful API calls to the Kubernetes-hosted service
- Pod logs showing healthy operation

---

### Task 4: Production Readiness (6 Points)

1. **Scaling:** Scale the API deployment to 3 replicas and verify all pods are running:
   ```bash
   kubectl scale deployment product-catalog-api --replicas=3 -n product-catalog
   kubectl get pods -n product-catalog
   ```

2. **Health Checks:** The Kubernetes manifests include `readinessProbe` and `livenessProbe`. Explain:
   - What is the difference between a readiness and a liveness probe?
   - What happens if the readiness probe fails? What about the liveness probe?
   - Why are different `initialDelaySeconds` values used?

3. **Resource Limits:** The API deployment specifies CPU and memory limits. Explain:
   - What happens if a pod exceeds its memory limit?
   - What happens if it exceeds its CPU limit?
   - Why are requests and limits both specified?

**Deliverable:** Add a `K8S.md` file with your answers and screenshots.

---

## Kubernetes Manifest Overview

| File | Contents |
|------|----------|
| `k8s/namespace.yml` | Namespace `product-catalog` |
| `k8s/postgres-deployment.yml` | PostgreSQL Deployment, Service, Secret, PVC |
| `k8s/api-deployment.yml` | API Deployment (2 replicas), NodePort Service |

---

## Useful Commands

```bash
# Minikube
minikube start / stop / delete
minikube dashboard                    # Open Kubernetes dashboard
minikube image load <image>:<tag>     # Load a host-built image into Minikube

# kubectl
kubectl get pods -n product-catalog
kubectl describe pod <name> -n product-catalog
kubectl logs <pod-name> -n product-catalog
kubectl exec -it <pod-name> -n product-catalog -- /bin/sh
kubectl port-forward svc/product-catalog-api 8080:8080 -n product-catalog

# Trivy
trivy image <image>
trivy fs .                            # Scan filesystem/dependencies
```

---

## Grading

| Task | Points |
|------|--------|
| Vulnerability Scanning -- Docker Image | 6 |
| Vulnerability Scanning -- Dependencies | 4 |
| Kubernetes Deployment with Minikube | 8 |
| Production Readiness | 6 |
| **Total** | **24** |

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)

