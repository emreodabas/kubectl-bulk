package service

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

func createClientset() (*kubernetes.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

var resourceList []model.Resource

func GetResourceList() ([]model.Resource, error) {

	if resourceList != nil {
		return resourceList, nil
	}

	result := make(map[string]model.Resource) // resource could have multi group

	clientset, err := createClientset()
	if err != nil {
		return nil, err
	}
	resources, _ := clientset.ServerPreferredResources()

	for _, res := range resources {
		if res.APIResources != nil && len(res.APIResources) > 0 {

			for i := 0; i < len(res.APIResources); i++ {

				result[res.APIResources[i].Name] =
					model.Resource{res.APIResources[i].Name, res.APIResources[i].Namespaced, res.APIResources[i].Kind, res.APIResources[i].ShortNames, res.APIResources[i].Verbs}
			}
		}
	}

	for i := range result {
		resourceList = append(resourceList, result[i])

	}
	return resourceList, nil
}

func GetResource(resourceName string) (model.Resource, error) {
	//TODO could be cached
	resourceList, _ := GetResourceList()
	resourceName = strings.ToLower(resourceName)
	fmt.Println("SIZE", len(resourceList))
	for i := 0; i < len(resourceList); i++ {
		fmt.Println(i, "-->", resourceList[i].Name, resourceList[i].ShortName)
		if strings.ToLower(resourceList[i].Name) == resourceName || contains(resourceList[i].ShortName, resourceName) {
			fmt.Println(resourceList[i], "is selected")
			return resourceList[i], nil
		}
	}
	return model.Resource{}, fmt.Errorf(resourceName + " is not a valid resource.")
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
