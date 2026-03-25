# Exercise 4: Vulnerability Scanning & Kubernetes Deployment

**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 30

## Learning Objectives

- Integrate vulnerability scanning into the CI/CD pipeline
- Scan Docker images and Go dependencies for known vulnerabilities
- Deploy a multi-tier application to Kubernetes (Minikube)
- Understand Kubernetes concepts: Deployments, Services, Secrets, Probes

## Prerequisites

- Completed Exercise 3 (CI pipeline with quality gates)
- Docker Desktop installed
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) installed
- [kubectl](https://kubernetes.io/docs/tasks/tools/) installed
- [Trivy](https://aquasecurity.github.io/trivy/) installed (optional, for local scanning)

## What's New in This Exercise

- **Kubernetes manifests** (`k8s/`) -- Deployment, Service, Secret, PVC
- **Trivy scanning** -- container image vulnerability scanning
- **Dependency scanning** -- Go module vulnerability checks
- **Complete CD pipeline** -- from code to running in Kubernetes

---

## Tasks

### Task 1: Vulnerability Scanning -- Docker Image (8 Points)

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

4. **Add a Trivy scan job to the CI pipeline:**
   ```yaml
   trivy-scan:
     runs-on: ubuntu-latest
     needs: docker-build
     steps:
       - uses: actions/checkout@v4
       - name: Build image
         run: docker build -t product-catalog:scan .
       - name: Run Trivy vulnerability scanner
         uses: aquasecurity/trivy-action@master
         with:
           image-ref: 'product-catalog:scan'
           format: 'table'
           exit-code: '1'
           severity: 'CRITICAL,HIGH'
   ```

**Deliverable:** Trivy scan output (before and after base image optimization). Updated CI workflow.

---

### Task 2: Vulnerability Scanning -- Dependencies (6 Points)

1. **Scan Go dependencies:**
   ```bash
   # Using govulncheck (official Go vulnerability checker)
   go install golang.org/x/vuln/cmd/govulncheck@latest
   govulncheck ./...
   ```

2. **Add dependency scanning to the CI pipeline:**
   ```yaml
   - name: Run govulncheck
     run: |
       go install golang.org/x/vuln/cmd/govulncheck@latest
       govulncheck ./...
   ```

3. **If vulnerabilities are found:**
   - Update the affected dependencies (`go get -u <module>`)
   - Document the CVEs and how you resolved them

**Deliverable:** govulncheck output. Updated `go.mod` if changes were needed.

---

### Task 3: Kubernetes Deployment with Minikube (10 Points)

1. **Start Minikube:**
   ```bash
   minikube start
   ```

2. **Build the image inside Minikube's Docker daemon:**
   ```bash
   eval $(minikube docker-env)
   docker build -t product-catalog:latest .
   ```

3. **Deploy the application:**
   ```bash
   kubectl apply -f k8s/namespace.yml
   kubectl apply -f k8s/postgres-deployment.yml
   kubectl apply -f k8s/api-deployment.yml
   ```

4. **Verify the deployment:**
   ```bash
   kubectl get all -n product-catalog
   kubectl logs deployment/product-catalog-api -n product-catalog
   ```

5. **Access the API:**
   ```bash
   minikube service product-catalog-api -n product-catalog --url
   # Use the returned URL to test the API
   curl <URL>/health
   curl <URL>/products
   ```

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
eval $(minikube docker-env)           # Use Minikube's Docker daemon

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
| Vulnerability Scanning -- Docker Image | 8 |
| Vulnerability Scanning -- Dependencies | 6 |
| Kubernetes Deployment with Minikube | 10 |
| Production Readiness | 6 |
| **Total** | **30** |
