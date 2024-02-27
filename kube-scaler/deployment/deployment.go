package deployment

import (
	"context"
	"fmt"
	"github.com/ChrisDevOpsOrg/kube-scaler/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentLister struct {
	ClientSet *client.ClientSets
}

func (dl DeploymentLister) ListResources(namespace string) ([]string, error) {
	deploymentList, err := dl.ClientSet.KubeClient.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	deployments := make([]string, len(deploymentList.Items))
	for i, deployment := range deploymentList.Items {
		deployments[i] = deployment.Name
	}

	return deployments, nil
}

func ScaleDeployments(clientSet *client.ClientSets, namespace string, replicas int32) error {
	deploymentList, err := clientSet.KubeClient.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, deployment := range deploymentList.Items {
		deployment.Spec.Replicas = &replicas
		fmt.Printf("start to scale deployment %s to %d", deployment.Name, replicas)
		_, err := clientSet.KubeClient.AppsV1().Deployments(namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
