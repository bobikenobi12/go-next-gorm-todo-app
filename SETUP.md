# Setup & Configuration Guide

This guide explains how to set up, provision, and deploy the **Dani Proekt** (Student Exercise) using a professional DevOps stack.

## 1. Prerequisites
- **Google Cloud Account**: A project created in GCP.
- **Pulumi CLI**: Installed and authenticated (`pulumi login`).
- **Infracost CLI**: Installed for cost estimation.
- **Docker Hub Account**: For hosting container images.
- **Pre-commit**: Installed locally (`brew install pre-commit`).

## 2. Local Infrastructure (IaC)
The infrastructure is defined in `iac/main.go` using Pulumi. It is optimized for **zero-to-low cost**:
- **Zonal GKE Cluster**: Single-zone deployment to avoid multi-zone fees.
- **Spot Instances**: Using `e2-small` Spot VMs (up to 90% discount).
- **Workload Identity Federation (WIF)**: Securely connects GitHub Actions to GCP without static JSON keys.

### To provision:
1. `cd iac`
2. `pulumi stack init dev`
3. `pulumi config set gcp:project YOUR_PROJECT_ID`
4. `pulumi preview` (Review the cost with Infracost first)
5. `pulumi up`

## 3. CI/CD Pipeline (GitHub Actions)
The pipeline is defined in `.github/workflows/ci.yaml`.

### Required GitHub Secrets:
To make the pipeline work, you MUST add these secrets to your GitHub Repository (**Settings > Secrets and variables > Actions**):

| Secret Name | Description |
| ----------- | ----------- |
| `DOCKERHUB_USERNAME` | Your Docker Hub username |
| `DOCKERHUB_TOKEN` | Your Docker Hub Personal Access Token |
| `INFRACOST_API_KEY` | Your Infracost API key (get it from `infracost auth login`) |
| `DISCORD_WEBHOOK` | (Optional) Discord/Slack webhook URL for notifications |
| `GOOGLE_WIF_PROVIDER` | The output from `pulumi stack output workloadIdentityProvider` |
| `GOOGLE_WIF_SA` | The output from `pulumi stack output githubActionsServiceAccount` |

## 4. Security & Code Quality
### Pre-commit Hooks
We use `pre-commit` to ensure no secrets or bad code enter the repo.
- **Secret Detection**: `gitleaks` and `detect-private-key`.
- **Linting**: `golangci-lint` (Backend) and `eslint` (Frontend).

Install them once:
```bash
pre-commit install
```

## 5. Cost Transparency (Infracost)
Every Pull Request will trigger an Infracost check. You can also run it locally:
```bash
cd iac
pulumi preview --json > plan.json
infracost breakdown --path plan.json
```
Current estimate for this setup: **~$7.50/month** (GKE management fee + ~$1.50 for spot nodes).
