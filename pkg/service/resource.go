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

	var result []model.Resource
	clientset, err := createClientset()
	if err != nil {
		return nil, err
	}
	resources, _ := clientset.ServerPreferredResources()

	for _, res := range resources {
		if res.APIResources != nil && len(res.APIResources) > 0 {

			result = append(result,
				model.Resource{res.APIResources[0].Name, res.APIResources[0].Namespaced, res.APIResources[0].Kind, res.APIResources[0].ShortNames, res.APIResources[0].Verbs})
		}
		//else {
		//	for i := 0; i < len(res.GroupVersion); i++ {
		//		resources = append(resources, )
		//	}
		//}
	}
	resourceList = result
	return result, nil
}

func GetResource(resourceName string) (model.Resource, error) {
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
