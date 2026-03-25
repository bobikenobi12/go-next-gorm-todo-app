# Phase 3: CI Pipeline (GitHub Actions)

**Goal:** Automate the Continuous Integration process including tests, building, pushing Docker images, and webhook notifications.

## Tasks:
- [x] Create `.github/workflows/ci.yaml`.
- [x] Trigger on `push` to `main` and `pull_request` to `main`.
- [x] **Backend Job:**
  - [x] Checkout code.
  - [x] Setup Go.
  - [x] Run Go tests (`go test ./...`).
  - [x] Run Go linter.
- [x] **Frontend Job:**
  - [x] Checkout code.
  - [x] Setup Node.js.
  - [x] Install dependencies (`pnpm install`).
  - [x] Run Frontend tests/linter.
- [x] **Build & Push Job (Docker):**
  - [x] Authenticate to Docker Hub (using GitHub Secrets for credentials).
  - [x] Build Backend Docker image.
  - [x] Push Backend image to Docker Hub (`username/backend:commit-sha`).
  - [x] Build Frontend Docker image.
  - [x] Push Frontend image to Docker Hub (`username/frontend:commit-sha`).
- [x] **Notifications:**
  - [x] Add a step to send a Webhook notification (e.g., to Slack/Discord) on pipeline success or failure.
- [x] **Cost Visibility:**
  - [x] Integrate Infracost to post cost estimates on Pull Requests.

## Notes
Make sure Docker Hub credentials are securely stored in GitHub Secrets.
