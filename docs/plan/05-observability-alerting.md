# Phase 5: Observability & Alerting

**Goal:** Gain visibility into the application and infrastructure with metrics, logs, and alerts.

## Tasks:
- [x] **Logging:**
  - [x] Ensure application logs are outputted to stdout/stderr.
  - [x] Rely on GKE's native Cloud Logging integration.
- [x] **Alerting:**
  - [x] Configure CI/CD status notifications.
  - [x] Configure Alertmanager (standard in GKE or via Discord webhooks on pipeline failures).
  - [x] Send a Discord Webhook notification when deployments fail.

## Notes
GCP Operations Suite (formerly Stackdriver) is also a strong candidate here since we are using GCP, but Prometheus/Grafana is standard for K8s. We will define which route to take when starting this phase.
