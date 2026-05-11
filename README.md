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
   > SonarCloud only fully analyses `main` and branches matching its long-lived branch regex. The default pattern is `(branch|release)-.*`, which does **not** match `exercise/03-ci-pipeline`. Without adjustment your branch is treated as short-lived and the analysis will be incomplete or not show up as expected.
   >
   > In SonarCloud: open your project → **Branches** in the left sidebar → click the **edit icon** next to *"Long-lived branches pattern"* in the top-right → replace the default with:
   >
   > ```
   > exercise/.*
   > ```
   >
   > → **Save**.
   >
   > The branch type is set on first analysis and cannot be changed afterwards. If you already triggered a scan for `exercise/03-ci-pipeline` before adjusting the pattern, delete that branch in SonarCloud (Branches page → select the branch → Delete) and push again to force a fresh first analysis.

3. **Add a `sonarcloud` job** to the CI workflow that:
   - Runs after the `test` job (`needs: test`)
   - Checks out the code with full history (`fetch-depth: 0`)
   - Downloads the coverage artifact from the test job
   - Runs the SonarCloud scan using `SonarSource/sonarqube-scan-action@v5`
   - Passes the `SONAR_TOKEN` as an environment variable

   > **Hint:** Look at the `sonar-project.properties` file to understand what SonarCloud expects.

4. **Add the `SONAR_TOKEN` secret** to your repository settings.

5. **Review the SonarCloud dashboard:**
   - What is the code coverage percentage?
   - Are there any code smells or bugs detected?
   - What is the technical debt estimate?

**Deliverable:** Link to your SonarCloud project dashboard. Screenshot showing the quality gate result.

---

### Task 4: Code Coverage Improvement (6 Points)

1. **Check current coverage:**
```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   go tool cover -html=coverage.out -o coverage.html
```

2. **Improve coverage to at least 80%** by adding tests for uncovered code paths. Focus on:
   - Edge cases in handlers (invalid IDs, malformed JSON)
   - Error paths in the store layer
   - The `Validate()` method edge cases

3. **Add a coverage threshold check** to the CI pipeline as a step after running tests:
   - Extract the total coverage percentage from `go tool cover -func`
   - Fail the build if coverage is below 80%
   - Use `::error::` to display the error in the GitHub Actions UI

   > **Hint:** `go tool cover -func=coverage.out | grep total` gives you the total line. Use `awk` and `sed` to extract the number. Use `bc` for the comparison (works on both Linux and macOS).

4. **Upload a coverage HTML report** as a build artifact:
   - Generate an HTML report using `go tool cover -html`
   - Upload it using `actions/upload-artifact@v4` so it can be downloaded from the Actions run

**Deliverable:** Coverage report showing >= 80%. Updated tests. Coverage HTML artifact downloadable from the Actions run.

---

## Grading

| Task | Points |
|------|--------|
| Matrix Builds | 4 |
| Linting with golangci-lint | 6 |
| SonarCloud Integration | 8 |
| Code Coverage Improvement | 6 |
| **Total** | **24** |

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)
