package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/iam"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "gcp")
		project := conf.Get("project")
		region := conf.Get("region")
		if region == "" {
			region = "europe-west1"
		}

		// Use a specific zone to save costs instead of a multi-zone region.
		zone := region + "-b"

		// Create a VPC network for the cluster
		network, err := compute.NewNetwork(ctx, "dani-proekt-vpc", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Create a subnet for the cluster
		subnet, err := compute.NewSubnetwork(ctx, "dani-proekt-subnet", &compute.SubnetworkArgs{
			Network:     network.ID(),
			IpCidrRange: pulumi.String("10.0.0.0/16"),
			Region:      pulumi.String(region),
		})
		if err != nil {
			return err
		}

		// Create a Service Account for the GKE nodes
		gkeSA, err := serviceaccount.Account(ctx, "gke-nodepool-sa", &serviceaccount.AccountArgs{
			AccountId:   pulumi.String("gke-nodepool-sa"),
			DisplayName: pulumi.String("GKE Node Service Account"),
		})
		if err != nil {
			return err
		}

		// Create the GKE cluster - explicitly use Zone to prevent 3x nodes
		cluster, err := container.NewCluster(ctx, "dani-proekt-cluster", &container.ClusterArgs{
			Location:              pulumi.String(zone),
			Network:               network.ID(),
			Subnetwork:            subnet.ID(),
			RemoveDefaultNodePool: pulumi.Bool(true),
			InitialNodeCount:      pulumi.Int(1),
			DeletionProtection:    pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Create a Node Pool with auto-scaling and spot instances
		_, err = container.NewNodePool(ctx, "primary-node-pool", &container.NodePoolArgs{
			Cluster:  cluster.Name,
			Location: cluster.Location,
			Autoscaling: &container.NodePoolAutoscalingArgs{
				MinNodeCount: pulumi.Int(1),
				MaxNodeCount: pulumi.Int(2), // Reduced for lower costs
			},
			NodeConfig: &container.NodePoolNodeConfigArgs{
				MachineType:    pulumi.String("e2-small"), // Smaller, cheaper instance
				Spot:           pulumi.Bool(true),         // Spot VMs are ~90% cheaper
				ServiceAccount: gkeSA.Email,
				OauthScopes: pulumi.StringArray{
					pulumi.String("https://www.googleapis.com/auth/cloud-platform"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Provision a GCP Artifact Registry for Docker images
		repo, err := artifactregistry.NewRepository(ctx, "dani-proekt-repo", &artifactregistry.RepositoryArgs{
			Location:     pulumi.String(region),
			RepositoryId: pulumi.String("dani-proekt-repo"),
			Description:  pulumi.String("Docker repository for the project"),
			Format:       pulumi.String("DOCKER"),
		})
		if err != nil {
			return err
		}

		// ---- Setup CI/CD (GitHub Actions) Access ----

		// 1. Create a Service Account for GitHub Actions
		ghSA, err := serviceaccount.Account(ctx, "github-actions-sa", &serviceaccount.AccountArgs{
			AccountId:   pulumi.String("github-actions-sa"),
			DisplayName: pulumi.String("GitHub Actions Service Account"),
		})
		if err != nil {
			return err
		}

		// 2. Bind necessary roles to GitHub Actions SA
		// Artifact Registry Writer
		_, err = projects.NewIAMBinding(ctx, "gh-sa-ar-writer", &projects.IAMBindingArgs{
			Project: pulumi.String(project),
			Role:    pulumi.String("roles/artifactregistry.writer"),
			Members: pulumi.StringArray{
				pulumi.Sprintf("serviceAccount:%s", ghSA.Email),
			},
		})
		if err != nil {
			return err
		}

		// Container Developer (deploying to GKE)
		_, err = projects.NewIAMBinding(ctx, "gh-sa-gke-dev", &projects.IAMBindingArgs{
			Project: pulumi.String(project),
			Role:    pulumi.String("roles/container.developer"),
			Members: pulumi.StringArray{
				pulumi.Sprintf("serviceAccount:%s", ghSA.Email),
			},
		})
		if err != nil {
			return err
		}

		// 3. Create Workload Identity Pool
		wiPool, err := iam.NewWorkloadIdentityPool(ctx, "github-actions-pool", &iam.WorkloadIdentityPoolArgs{
			WorkloadIdentityPoolId: pulumi.String("github-actions-pool"),
			DisplayName:            pulumi.String("GitHub Actions Pool"),
			Description:            pulumi.String("Identity pool for automated CI/CD"),
		})
		if err != nil {
			return err
		}

		// 4. Create Workload Identity Provider for GitHub
		wiProvider, err := iam.NewWorkloadIdentityPoolProvider(ctx, "github-provider", &iam.WorkloadIdentityPoolProviderArgs{
			WorkloadIdentityPoolId:         wiPool.WorkloadIdentityPoolId,
			WorkloadIdentityPoolProviderId: pulumi.String("github-provider"),
			DisplayName:                    pulumi.String("GitHub Actions Provider"),
			AttributeMapping: pulumi.StringMap{
				"google.subject":       pulumi.String("assertion.sub"),
				"attribute.actor":      pulumi.String("assertion.actor"),
				"attribute.repository": pulumi.String("assertion.repository"),
			},
			Oidc: &iam.WorkloadIdentityPoolProviderOidcArgs{
				IssuerUri: pulumi.String("https://token.actions.githubusercontent.com"),
			},
		})
		if err != nil {
			return err
		}

		// 5. Allow any GitHub repo authenticating via this pool to impersonate the SA
		_, err = serviceaccount.NewIAMBinding(ctx, "gh-sa-wif-binding", &serviceaccount.IAMBindingArgs{
			ServiceAccountId: ghSA.Name,
			Role:             pulumi.String("roles/iam.workloadIdentityUser"),
			Members: pulumi.StringArray{
				pulumi.Sprintf("principalSet://iam.googleapis.com/%s/*", wiPool.Name),
			},
		})
		if err != nil {
			return err
		}

		// Export Outputs
		ctx.Export("clusterName", cluster.Name)
		ctx.Export("clusterLocation", cluster.Location)
		ctx.Export("registryUrl", pulumi.Sprintf("%s-docker.pkg.dev/%s/%s", region, project, repo.RepositoryId))
		ctx.Export("githubActionsServiceAccount", ghSA.Email)
		ctx.Export("workloadIdentityProvider", wiProvider.Name)

		return nil
	})
}
