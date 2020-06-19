package service

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"github.com/emreodabas/kubectl-bulk/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func UpdateResources(command *model.Command) error {

	/*
		setAnnotations
		setLabels
		updateUnstructred Content
		 	[specs]

	*/
	var removedItems = []string{}
	var definedValues = ""
	changedValues := make(map[string]string)
	selectedPreference := interaction.ShowLists(updatePreference)
	switch selectedPreference {
	case "[update] [labels]":
		keys := getLabelKeys(command)
		selection := interaction.ShowList(keys)
		definedValues = interaction.Prompt("Define value for label \"%s\" ", selection)
		changedValues[selection] = definedValues
	case "[add] [labels]":
		definedValues = interaction.Prompt("Define your new label like env=prod,app=nginx ")
		valueToMap(definedValues, changedValues)

	case "[remove] [labels]":
		keys := getLabelKeys(command)
		prompt := &survey.MultiSelect{
			Message:  "Which labels do you want to remove ?",
			Options:  keys,
			PageSize: 20,
		}
		survey.AskOne(prompt, &removedItems)

	case "[add] [annotations]":
		definedValues = interaction.Prompt("Define your new annotations like key1=value,key2=value ")
		valueToMap(definedValues, changedValues)
	case "[update] [annotations]":
		keys := getAnnotations(command)
		selection := interaction.ShowList(keys)
		definedValues = interaction.Prompt("Define value for annotations \"%s\" ", selection)
		changedValues[selection] = definedValues
	case "[remove] [annotations]":
		keys := getAnnotations(command)
		prompt := &survey.MultiSelect{
			Message:  "Which labels do you want to remove ?",
			Options:  keys,
			PageSize: 20,
		}
		survey.AskOne(prompt, &removedItems)
	case "[add] [specs]":
		selection := getOneLevelOfSpecs(command)
		definedValues = interaction.Prompt("Define spec  \"%s\" ", selection)
	case "[update] [specs]":
		keys := getSpecs(command)
		selection := interaction.ShowList(keys)
		definedValues = interaction.Prompt("Define spec  \"%s\" ", selection)
		changedValues[selection] = definedValues
		break
	case "[remove] [specs]":
		break

	}
	updateResources(command, selectedPreference, changedValues, removedItems)

	return nil
}

func valueToMap(definedValues string, changedLabels map[string]string) {
	equations := strings.Split(definedValues, ",")
	for _, item := range equations {
		label := strings.Split(item, "=")
		if len(label) == 2 {
			changedLabels[label[0]] = label[1]
		}
	}
}

func getLabelKeys(command *model.Command) []string {
	var list = command.List
	var labels []string
	for i := 0; i < len(list); i++ {
		for k, _ := range list[i].GetLabels() {
			if !utils.Contains(labels, k) {
				labels = append(labels, k)
			}
		}
	}
	return labels
}

func getAnnotations(command *model.Command) []string {
	var list = command.List
	var labels []string
	for i := 0; i < len(list); i++ {
		for k, _ := range list[i].GetAnnotations() {
			if !utils.Contains(labels, k) {
				labels = append(labels, k)
			}
		}
	}
	return labels
}

func getSpecs(command *model.Command) []string {
	var list = command.List
	specTree := make(map[string]string)
	for i := 0; i < len(list); i++ {
		content := list[i].UnstructuredContent()
		specs := content["spec"]
		m := specs.(map[string]interface{})
		appendSpecToList(m, specTree, "")
	}

	keys := utils.Keys(specTree)
	sort.Strings(keys)
	return keys
}

func getOneLevelOfSpecs(command *model.Command) string {
	var list = command.List
	specTree := make(map[string]interface{})
	for i := 0; i < len(list); i++ {
		content := list[i].UnstructuredContent()
		specs := content["spec"]
		m := specs.(map[string]interface{})
		for k, v := range m {
			specTree[k] = v
		}
	}
	return selectSpecField(specTree)

}

func selectSpecField(obj interface{}) string {

	keys := make(map[string]interface{})
	var keyValues []string
	switch typ := obj.(type) {
	case interface{}:
		vv := reflect.ValueOf(typ)
		if vv.Kind() == reflect.Map {
			for _, maps := range vv.MapKeys() {
				keys[maps.String()] = vv.MapIndex(maps).Interface()
			}
			for k, _ := range keys {
				keyValues = append(keyValues, k)
			}
			sort.Strings(keyValues)
			selection := interaction.ShowList(keyValues)
			selectSpecField(keys[selection])
		}
	case string:
		return fmt.Sprint(typ)
	default:
		return fmt.Sprint(typ)
	}

	return fmt.Sprint(obj)
}
func appendSpecToList(obj interface{}, specTree map[string]string, key string) {
	//tree
	//var isMap = false
	//for k, v := range obj {

	switch typ := obj.(type) {
	case map[string]interface{}:
		//TODO tree structured spec could be better to show select and update
		for k, v := range typ {
			key = k + "∟--" + key
			appendSpecToList(v, specTree, key)
		}
	case string:
		specTree[key] = typ
	case bool:
		specTree[key] = strconv.FormatBool(typ)
	case int:
		specTree[key] = strconv.Itoa(typ)
	case interface{}:
		vv := reflect.ValueOf(typ)
		if vv.Kind() == reflect.Map {
			for _, maps := range vv.MapKeys() {
				appendSpecToList(maps, specTree, key)
			}
		}
	}
}

func updateResources(command *model.Command, actionType string, values map[string]string, removedValues []string) error {
	clientset, _, err := getClientSet()
	resource := command.Resource

	if err != nil {
		return err
	}
	for i := 0; i < len(resource.GroupVersion); i++ {
		gv := resource.GroupVersion[i]
		meta := v1.TypeMeta{
			Kind:       resource.Kind,
			APIVersion: gv.Version,
		}
		options := v1.UpdateOptions{
			TypeMeta:     meta,
			DryRun:       nil,
			FieldManager: "",
		}
		resourceInterface := clientset.Resource(schema.GroupVersionResource{
			Group:    gv.Group,
			Version:  gv.Version,
			Resource: resource.Name,
		})
		var valueList map[string]string
		list := command.List
		for i := 0; i < len(list); i++ {

			if strings.Contains(actionType, "label") {
				valueList = list[i].GetLabels()
			} else if strings.Contains(actionType, "annotation") {
				valueList = list[i].GetAnnotations()
			}

			if strings.Contains(actionType, "add") {
				for k, v := range values {
					valueList[k] = v
				}
			} else if strings.Contains(actionType, "update") {
				for k, v := range values {
					if valueList[k] != "" {
						valueList[k] = v
					}
				}
			} else if strings.Contains(actionType, "remove") {
				for i := 0; i < len(removedValues); i++ {
					delete(valueList, removedValues[i])
				}
			}
			if strings.Contains(actionType, "label") {
				list[i].SetLabels(valueList)
			} else if strings.Contains(actionType, "annotations") {
				list[i].SetAnnotations(valueList)
			}
			if resource.Namespaced {
				_, err := resourceInterface.Namespace(command.Namespace).Update(&list[i], options)
				if err != nil {
					return err
				}
			} else {
				_, err := resourceInterface.Update(&list[i], options)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

var updatePreference = [][]string{
	{"[add] [labels]", "these feature add labels to selected resource\n IF LABEL EXIST IT WILL BE UPDATED"},
	{"[add] [annotations]", "these feature add annotations to selected resource\n IF ANNOTATION EXIST IT WILL BE UPDATED"},
	{"[add] [specs]", "these feature add spec to selected resource\n IF SPEC EXIST IT WILL BE UPDATED"},
	{"[remove] [labels]", "these feature remove labels of selected resource\n IF LABEL NOT EXIST NO CHANGE"},
	{"[remove] [annotations]", "these feature remove annotations of selected resource\n IF ANNOTATION NOT EXIST NO CHANGE"},
	{"[remove] [specs]", "these feature remove specs of selected resource\n IF SPECS NOT EXIST NO CHANGE"},
	{"[update] [labels]", "these feature update labels of selected resource\n IF LABEL NOT EXIST NO CHANGE"},
	{"[update] [annotations]", "these feature update annotation of selected resource\n IF ANNOTATION NOT EXIST NO CHANGE"},
	{"[update] [specs]", "these feature update specs of selected resource\n IF LABEL NOT EXIST NO CHANGE"},
	//{"[upsert] [labels]", "these feature update labels of selected resource\n IF LABEL NOT EXIST IT WIL BE ADDED"},
	//{"[upsert] [annotations]", "these feature update annotation of selected resource\n IF ANNOTATION NOT EXIST IT WIL BE ADDED"},
	//{"[upsert] [specs]", "these feature update spec of selected resource\n IF SPEC NOT EXIST IT WIL BE ADDED"},
}
