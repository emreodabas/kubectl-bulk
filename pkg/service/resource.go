package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"github.com/emreodabas/kubectl-bulk/pkg/utils"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"strings"
)

var resourceList []model.Resource

const path = ".api-resource-cache.json"

func ResourceSelection(command *model.Command) error {

	list, err := GetResourceList()
	if err != nil {
		return err
	}
	command.Resource = interaction.ShowResourceList(list)
	return nil
}

func GetResourceList() ([]model.Resource, error) {
	if cacheFileExist() {
		resourceList, err := readResourcelist()

		if err != nil {
			return nil, err
		}
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
				groupV := result[res.APIResources[i].Kind].GroupVersion
				if groupV != nil {
					groupV = append(groupV, groupVersion)
					result[res.APIResources[i].Kind] =
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
					result[res.APIResources[i].Kind] =
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
	saveResourcelist(resourceList)
	return resourceList, nil
}

func saveResourcelist(resourceList []model.Resource) {
	//resources := model.Resources{Resources: resourceList}
	_, err := WriteDataToFileAsJSON(resourceList, path)
	if err != nil {
		fmt.Println(err)
	}
}

func WriteDataToFileAsJSON(data interface{}, filedir string) (int, error) {
	//write data as buffer to json encoder
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return 0, err
	}
	file, err := os.OpenFile(filedir, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return 0, err
	}
	n, err := file.Write(buffer.Bytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}

func cacheFileExist() bool {

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()

}

func readResourcelist() ([]model.Resource, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("ERROR!!")
		fmt.Println(err)
		return []model.Resource{}, err
	}

	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we initialize our Users array
	var resources []model.Resource

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &resources)
	if err != nil {
		fmt.Errorf("Marshall problem")
	}
	return resources, nil
}

func GetResource(resourceName string) (model.Resource, error) {
	//TODO could be cached
	resourceList, _ := GetResourceList()
	resourceName = strings.ToLower(resourceName)
	for i := 0; i < len(resourceList); i++ {
		if strings.ToLower(resourceList[i].Kind) == resourceName || utils.Contains(resourceList[i].ShortName, resourceName) {
			return resourceList[i], nil
		}
	}
	return model.Resource{}, fmt.Errorf(resourceName + " is not a valid resource.")
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

func SourceSelection(command *model.Command) error {
	// filter or multi selection could be ask to user
	var err error
	if command.Resource.Namespaced {
		namespaces, err := GetNamespaces()
		command.Namespace = interaction.ShowList(namespaces)
		if err != nil {

			return fmt.Errorf("Namespace list could not fetch")
		}
		err = FetchInstances(command)
	} else {
		err = FetchInstances(command)
	}
	if err != nil {
		return err
	}
	return err

}

func FetchInstances(command *model.Command) error {

	var res []unstructured.Unstructured
	var list *unstructured.UnstructuredList
	resource := command.Resource
	namespace := command.Namespace
	clientset, _, err := getClientSet()
	if err != nil {
		return fmt.Errorf("K8s client could not created ")
	}

	var next string
	options := v1.ListOptions{
		Limit:    250,
		Continue: next,
	}
	if command.Label != "" {
		options.LabelSelector = command.Label
	}

	if command.FieldSelector != "" {
		options.FieldSelector = command.FieldSelector
	}
	for i := 0; i < len(resource.GroupVersion); i++ {
		gv := resource.GroupVersion[i]

		resourceInterface := clientset.Resource(schema.GroupVersionResource{
			Group:    gv.Group,
			Version:  gv.Version,
			Resource: resource.Name,
		})
		if namespace != "all-namespaces[-A]" {
			for {
				list, err = resourceInterface.Namespace(namespace).List(options)
				if err != nil {
					return err
				}
				res = append(res, list.Items...)
				next = list.GetContinue()
				if next == "" {
					break
				}
			}
		} else {
			for {
				list, err = resourceInterface.List(options)
				if err != nil {
					return err
				}

				if list != nil {
					res = append(res, list.Items...)
					next = list.GetContinue()
					if next == "" {
						break
					}
				} else {
					break
				}
			}
		}

	}
	if list != nil {
		command.List = utils.Unique(list)
	}
	return nil
}
