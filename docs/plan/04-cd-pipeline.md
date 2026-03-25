# Phase 4: CD Pipeline (GitHub Runner to K8s)

**Goal:** Automate the deployment of the application to the Kubernetes cluster using a GitHub Actions Runner. No ArgoCD as requested.

## Tasks:
- [x] Set up a GitHub Actions self-hosted runner (using GitHub managed runners with GCP Workload Identity Federation).
- [x] Create `.github/workflows/cd.yaml`.
- [x] Trigger on successful completion of the CI pipeline on the `main` branch.
- [x] **Deploy Job:**
  - [x] Authenticate to GCP using Workload Identity Federation (securely connected via Pulumi).
  - [x] Get GKE cluster credentials (`gcloud container clusters get-credentials`).
  - [x] Update Kubernetes manifests (`k8s/` folder) with the new Docker image tags from the CI pipeline (using `sed`).
  - [x] Apply Kubernetes manifests (`kubectl apply -f k8s/`).
  - [x] Verify deployment success (`kubectl rollout status`).
- [x] **Notifications:**
  - [x] Send webhook notification upon successful or failed deployment.

## Notes
Since we are not using ArgoCD, GitHub Actions will directly push changes to the GKE cluster using `kubectl`.
