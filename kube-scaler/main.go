package main

import (
	"context"
	"fmt"
	"github.com/ChrisDevOpsOrg/kube-scaler/argorollout"
	"github.com/ChrisDevOpsOrg/kube-scaler/client"
	"github.com/ChrisDevOpsOrg/kube-scaler/config"
	"github.com/ChrisDevOpsOrg/kube-scaler/deployment"
	"log"
	"time"
)

type ResourceLister interface {
	ListResources(namespace string) ([]string, error)
}

type ResourceScaler interface {
	ScaleResource(namespace, name string, replicas int32) error
}

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	clientSet, err := client.GetClientSet(c.UseInClusterClient, ctx)
	if err != nil {
		log.Fatal("failed to create clientset: ", err)
	}

	if c.ResourceKind == "DEPLOYMENT" {
		// Create the listers
		deploymentLister := deployment.DeploymentLister{ClientSet: clientSet}

		deploymentNames, err := deploymentLister.ListResources(c.Namespace)
		if err != nil {
			log.Fatal("failed to list deployments: ", err)
		}

		// Print the retrieved resource names
		for _, name := range deploymentNames {
			log.Println("Deployment:", name)
		}

		deployment.ScaleDeployments(clientSet, c.Namespace, c.Replicas)
	} else if c.ResourceKind == "ARGO_ROLLOUT" {
		rolloutLister := argorollout.RolloutLister{ClientSet: clientSet}
		rolloutNames, err := rolloutLister.ListResources(c.Namespace)
		if err != nil {
			log.Fatal("failed to list rollouts: ", err)
		}

		for _, name := range rolloutNames {
			log.Println("Rollout:", name)
		}

		argorollout.ScaleRollouts(clientSet, c.Namespace, c.Replicas)
	} else {
		fmt.Printf("Unsupported resource kind %s\n", c.ResourceKind)
	}

}
