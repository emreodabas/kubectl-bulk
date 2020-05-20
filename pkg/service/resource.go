package service

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"github.com/emreodabas/kubectl-bulk/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
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

var resourceList []model.Resource

func GetResourceList() ([]model.Resource, error) {

	if resourceList != nil {
		return resourceList, nil
	}

	result := make(map[string]model.Resource) // resource could have multi group

	_, clientset, err := getClientSet()
	if err != nil {
		return nil, err
	}
	resources, _ := clientset.ServerPreferredResources()

	for _, res := range resources {
		if res.APIResources != nil && len(res.APIResources) > 0 {
			groupVersion, _ := schema.ParseGroupVersion(res.GroupVersion)
			for i := 0; i < len(res.APIResources); i++ {
				groupV := result[res.APIResources[i].Name].GroupVersion
				if groupV != nil {
					groupV = append(groupV, groupVersion)
					result[res.APIResources[i].Name] =
						model.Resource{res.APIResources[i].Name,
							res.APIResources[i].Namespaced,
							res.APIResources[i].Kind,
							res.APIResources[i].ShortNames,
							res.APIResources[i].Verbs,
							res.GroupVersionKind(),
							groupV,
						}
				} else {
					var groupV = []schema.GroupVersion{groupVersion}
					result[res.APIResources[i].Name] =
						model.Resource{res.APIResources[i].Name,
							res.APIResources[i].Namespaced,
							res.APIResources[i].Kind,
							res.APIResources[i].ShortNames,
							res.APIResources[i].Verbs,
							res.GroupVersionKind(),
							groupV,
						}

				}

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
	for i := 0; i < len(resourceList); i++ {
		if strings.ToLower(resourceList[i].Name) == resourceName || utils.Contains(resourceList[i].ShortName, resourceName) {
			return resourceList[i], nil
		}
	}
	return model.Resource{}, fmt.Errorf(resourceName + " is not a valid resource.")
}

func FetchInstances(command *model.Command) error {

	dyn, _, err := getClientSet()
	var res []unstructured.Unstructured
	var list *unstructured.UnstructuredList
	resource := command.Resource
	namespace := command.Namespace
	if err != nil {
		return fmt.Errorf("K8s client could not created ")
	}
	for i := 0; i < len(resource.GroupVersion); i++ {
		gv := resource.GroupVersion[i]

		resourceInterface := dyn.Resource(schema.GroupVersionResource{
			Group:    gv.Group,
			Version:  gv.Version,
			Resource: resource.Name,
		})
		if namespace != "" || namespace != "all-namespaces[-A]" {
			list, err = resourceInterface.Namespace(namespace).List(v1.ListOptions{
				Limit:    250,
				Continue: "",
			})

		} else {
			list, err = resourceInterface.List(v1.ListOptions{
				Limit:    250,
				Continue: "",
			})

		}
		if err != nil {
			return fmt.Errorf("someting goes wrong while fetching ", resource.Name)
		}

		res = append(res, list.Items...)
	}
	command.List = res
	return nil
}

func GetNamespaces() ([]string, error) {
	var res []string
	res = append(res, "all-namespaces[-A]")
	_, clientset, err := getClientSet()
	if err != nil {
		return nil, fmt.Errorf("K8s client could not created")
	}
	var next string

	for {
		list, _ := clientset.CoreV1().Namespaces().List(v1.ListOptions{
			Limit:    250,
			Continue: next,
		})
		for i := 0; i < len(list.Items); i++ {
			res = append(res, list.Items[i].Name)
		}

		next = list.GetContinue()
		if next == "" {
			break
		}
	}
	return res, nil
}
