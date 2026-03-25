# Phase 6: Configuration & Secrets Management

**Goal:** Securely manage configuration and secrets without hardcoding them in the repository.

## Tasks:
- [x] Choose the Secret Management tool: GitHub Secrets (for CI/CD) and GCP Secret Manager (for runtime).
- [x] Store Pulumi state, GCP Workload Identity Provider details, and Docker Hub tokens in GitHub Secrets.
- [x] Provision GCP Workload Identity Federation (WIF) via Pulumi (zero-key security).
- [x] Document secret setup in `SETUP.md`.

## Notes
Never commit raw passwords or tokens. All runtime secrets should be injected securely at runtime.
