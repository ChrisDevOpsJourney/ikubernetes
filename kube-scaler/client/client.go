package client

import (
	"context"
	"github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type ClientSets struct {
	KubeClient *kubernetes.Clientset
	ArgoClient *versioned.Clientset
	ctx        context.Context
}

func GetClientSet(useInCluster bool, ctx context.Context) (*ClientSets, error) {
	cs := &ClientSets{ctx: ctx}

	var config *rest.Config
	var err error
	if useInCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	cs.KubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cs.ArgoClient, err = versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return cs, nil
}
