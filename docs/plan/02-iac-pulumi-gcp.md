# Phase 2: Infrastructure as Code (IaC) with Pulumi & GCP

**Goal:** Automate the provisioning of all infrastructure on Google Cloud Platform using Pulumi, specifically focusing on a Kubernetes cluster (GKE).

## Tasks:
- [x] Initialize a new Pulumi project in an `iac/` or `pulumi/` directory.
- [x] Configure Pulumi to use the Google Cloud Platform (GCP) provider.
- [x] Create a Virtual Private Cloud (VPC) network and subnets for the cluster.
- [x] Provision a Google Kubernetes Engine (GKE) cluster.
- [x] Configure Node Pools with appropriate machine types and auto-scaling.
- [x] Provision a GCP Artifact Registry for Docker images.
- [x] Configure least-privilege Service Accounts for the GKE nodes and GitHub Actions runner.
- [x] Set up state management for Pulumi (using local/managed Pulumi backend instead of GCP Cloud Storage for zero cost).
- [ ] Test `pulumi up` to verify infrastructure creation.

## Notes
All infrastructure changes must be done via Pulumi. No manual clicks in the GCP Console.
