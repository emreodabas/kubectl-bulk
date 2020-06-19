package service

import (
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientSet *kubernetes.Clientset
var dinamic dynamic.Interface

func getClientSet() (dynamic.Interface, *kubernetes.Clientset, error) {
	if clientSet != nil && dinamic != nil {
		return dinamic, clientSet, nil
	}
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("Kubernetes Client could not configured")
	}

	clientSet, err = kubernetes.NewForConfig(config)
	dinamic, err = dynamic.NewForConfig(config)

	if err != nil {
		return nil, nil, fmt.Errorf("Kubernetes Client could not configured")
	}

	return dinamic, clientSet, nil
}
