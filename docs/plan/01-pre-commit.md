# Phase 1: Pre-commit Hooks

**Goal:** Ensure code quality and prevent secrets from being committed to the Git repository.

## Tasks:
- [x] Install `pre-commit` tool locally.
- [x] Create `.pre-commit-config.yaml` in the root of the repository.
- [x] Add hook: `trailing-whitespace`.
- [x] Add hook: `end-of-file-fixer`.
- [x] Add hook: `check-yaml`.
- [x] Add hook: `detect-private-key` to prevent committing secrets.
- [x] Add hook: `gitleaks` or similar secret scanner.
- [x] Add hook: Go linter (`golangci-lint`) for backend.
- [x] Add hook: ESLint/Prettier for frontend.
- [x] Document in `README.md` how developers should install the hooks locally (`pre-commit install`).

## Notes
These hooks will run locally before every commit to ensure no bad code or secrets are leaked into the repository.
