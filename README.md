# Exercise 3: CI Pipeline -- SonarCloud, Matrix Builds & Linting

**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 24

## Learning Objectives

- Extend a CI pipeline with quality gates and code analysis
- Configure SonarCloud for static code analysis and coverage tracking
- Use matrix builds to test across multiple Go versions
- Integrate linting with golangci-lint
- Understand code quality metrics and technical debt

## Prerequisites

- Completed Exercise 2 (working CI pipeline with Docker build)
- SonarCloud account (free for open-source projects)
- Understanding of GitHub Actions workflow syntax

## What's New in This Exercise

- **Matrix builds** in `.github/workflows/ci.yml` -- test across multiple Go versions
- **SonarCloud configuration** (`sonar-project.properties`) -- static analysis setup
- **golangci-lint configuration** (`.golangci.yml`) -- linter rules
- **Coverage reporting** -- `go test -coverprofile`

---

## Tasks

### Task 1: Matrix Builds (4 Points)

The CI workflow already has a matrix strategy with one Go version. Your tasks:

1. **Extend the matrix** to include Go versions `1.25` and `1.26` (see the TODO in `ci.yml`).
2. **Verify** that the pipeline runs tests for both Go versions in parallel.
3. **Add an OS matrix dimension** (`ubuntu-latest`, `macos-latest`) so tests run on both platforms.

**Expected result:** 4 parallel test jobs (2 Go versions x 2 OS).

**Deliverable:** Screenshot of the GitHub Actions matrix view showing all jobs.

---

### Task 2: Linting with golangci-lint (6 Points)

1. **Add a `lint` job** to the CI workflow that:
   - Runs `golangci-lint` using the `golangci/golangci-lint-action@v4` action
   - Uses the `.golangci.yml` configuration file
   - Runs in parallel with the test matrix (does not depend on `test`)

2. **Enable additional linters** in `.golangci.yml` (see TODOs):
   - `gofmt` -- enforces standard Go formatting
   - `gocyclo` -- detects overly complex functions
   - `misspell` -- catches common typos
   - `gocritic` -- advanced Go code analysis

3. **Fix any linting issues** that are reported in the existing code.

**Deliverable:** Clean lint run (no warnings). Screenshot of the lint job passing.

---

### Task 3: SonarCloud Integration (8 Points)

1. **Create a SonarCloud project:**
   - Go to [sonarcloud.io](https://sonarcloud.io) and sign in with GitHub.
   - Import your repository as a new project.
   - Note your `projectKey` and `organization`.

2. **Configure `sonar-project.properties`:**
   - Replace `YOUR_PROJECT_KEY` and `YOUR_ORGANIZATION` with your actual values.
   - Ensure coverage reporting is configured correctly.

   > **Important -- adjust the long-lived branch pattern *before* your first CI run:**
   >
   > SonarCloud only fully analyses `
