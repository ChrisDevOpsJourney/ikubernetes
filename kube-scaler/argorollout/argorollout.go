package argorollout

import (
	"context"
	"github.com/ChrisDevOpsOrg/kube-scaler/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RolloutLister struct {
	ClientSet *client.ClientSets
}

func (rl RolloutLister) ListResources(namespace string) ([]string, error) {
	rolloutList, err := rl.ClientSet.ArgoClient.ArgoprojV1alpha1().Rollouts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	rollouts := make([]string, len(rolloutList.Items))
	for i, rollout := range rolloutList.Items {
		rollouts[i] = rollout.Name
	}

	return rollouts, nil
}

func ScaleRollouts(clientSet *client.ClientSets, namespace string, replicas int32) error {
	rolloutList, err := clientSet.ArgoClient.ArgoprojV1alpha1().Rollouts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, rollout := range rolloutList.Items {
		rollout.Spec.Replicas = &replicas
		_, err := clientSet.ArgoClient.ArgoprojV1alpha1().Rollouts(namespace).Update(context.TODO(), &rollout, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
