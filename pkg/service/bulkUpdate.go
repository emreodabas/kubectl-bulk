package service

import (
	"encoding/json"
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

func BulkUpdateResources(command *model.Command) error {

	var items = []string{}
	var definedValues = ""
	values := make(map[string]string)
	selectedPreference := interaction.ShowLists(updatePreference)
	switch selectedPreference {
	case "[update] [labels]":
		keys := getLabelKeys(command)
		selection := interaction.ShowList(keys)
		definedValues = interaction.Prompt("Define value for label \"%s\" ", selection)
		values[selection] = definedValues
	case "[add] [labels]":
		definedValues = interaction.Prompt("Define your new label like env=prod,app=nginx ")
		valueToMap(definedValues, values)

	case "[remove] [labels]":
		keys := getLabelKeys(command)
		prompt := &survey.MultiSelect{
			Message:  "Which labels do you want to remove ?",
			Options:  keys,
			PageSize: 20,
		}
		survey.AskOne(prompt, &items)

	case "[add] [annotations]":
		definedValues = interaction.Prompt("Define your new annotations like key1=value,key2=value ")
		valueToMap(definedValues, values)
	case "[update] [annotations]":
		keys := getAnnotations(command)
		selection := interaction.ShowList(keys)
		definedValues = interaction.Prompt("Define value for annotations \"%s\" ", selection)
		values[selection] = definedValues
	case "[remove] [annotations]":
		keys := getAnnotations(command)
		prompt := &survey.MultiSelect{
			Message:  "Which labels do you want to remove ?",
			Options:  keys,
			PageSize: 20,
		}
		survey.AskOne(prompt, &items)
	case "[add] [specs]":
		addIterateLevelOfSpecs(command)
		definedValues = interaction.Prompt("Define spec  \"%s\" \n use comma for multi value \n use \"\" for strings -> name=\"abc\" \n for numerics -> replicas=3", strings.Join(command.SelectedSpec, "."))
		valueToMap(definedValues, values)
	case "[update] [specs]":
		isObject := getOneLevelOfSpecs(command)
		if isObject {
			definedValues = interaction.Prompt("Define spec \"%s\" \n use comma for multi value \n use \"\" for strings -> name=\"abc\" \n for numerics -> replicas=3", strings.Join(command.SelectedSpec, "."))
		} else {
			definedValues = interaction.Prompt("Define spec value for \"%s\" ", strings.Join(command.SelectedSpec, "."))
		}
		values["value"] = definedValues
	case "[remove] [specs]":
		getOneLevelOfSpecs(command)
	}
	err := updateResources(command, selectedPreference, values, items)

	if err != nil {
		return err
	}
	return nil
}

func valueToMap(definedValues string, changedLabels map[string]string) {
	equations := strings.Split(definedValues, ",")
	for _, item := range equations {
		label := strings.Split(item, "=")
		if len(label) == 2 {
			changedLabels[label[0]] = label[1]
		} else {
			label := strings.Split(item, ":")
			if len(label) == 2 {
				changedLabels[label[0]] = label[1]
			} else {
				fmt.Errorf("Parsing problem")
			}
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

func getOneLevelOfSpecs(command *model.Command) bool {
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

	command.SelectedSpec = append(command.SelectedSpec, "spec")
	return updateSpecFieldSelection(specTree, command)
}

func updateSpecFieldSelection(obj interface{}, command *model.Command) bool {
	var ret = false
	keys := make(map[string]interface{})
	var keyValues []string
	var selection string
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
			if len(keyValues) > 0 {
				sort.Strings(keyValues)
				selection = interaction.ShowList(keyValues)
			} else if len(keyValues) == 0 {
				return true
			}
			command.SelectedSpec = append(command.SelectedSpec, selection)
			ret = updateSpecFieldSelection(keys[selection], command)
		} else if vv.Kind() == reflect.Slice {
			command.SelectedSpec = append(command.SelectedSpec, "[0]")
			ret = updateSpecFieldSelection(vv.Index(0).Interface(), command)
		}
	}
	return ret
}

func addIterateLevelOfSpecs(command *model.Command) {
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

	command.SelectedSpec = append(command.SelectedSpec, "spec")
	addSpecFieldSelection(specTree, command)
}

func addSpecFieldSelection(obj interface{}, command *model.Command) {

	keys := make(map[string]interface{})
	var keyValues []string
	var selection string
	switch typ := obj.(type) {
	case interface{}:
		vv := reflect.ValueOf(typ)
		if vv.Kind() == reflect.Map {
			for _, maps := range vv.MapKeys() {
				if hasChild(vv.MapIndex(maps).Interface()) {
					keys[maps.String()] = vv.MapIndex(maps).Interface()
				}
			}
			keys["<Add Here>"] = ""
			for k, _ := range keys {
				keyValues = append(keyValues, k)
			}
			sort.Strings(keyValues)
			selection = interaction.ShowList(keyValues)
			if selection == "<Add Here>" {
				return
			}
			command.SelectedSpec = append(command.SelectedSpec, selection)
			addSpecFieldSelection(keys[selection], command)
		} else if vv.Kind() == reflect.Slice {
			command.SelectedSpec = append(command.SelectedSpec, "[0]")
			addSpecFieldSelection(vv.Index(0).Interface(), command)
		}
	}

}

func hasChild(obj interface{}) bool {
	switch typ := obj.(type) {
	case interface{}:
		vv := reflect.ValueOf(typ)
		if vv.Kind() == reflect.Map {
			return true
		} else if vv.Kind() == reflect.Slice {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func addSpecField(item interface{}, selection []string, values map[string]string) (interface{}, error) {

	if len(selection) > 0 {
		switch typ := item.(type) {
		case interface{}:
			vv := reflect.ValueOf(typ)
			if vv.Kind() == reflect.Map {
				child := vv.MapIndex(reflect.ValueOf(selection[0])).Interface()
				resp, err := addSpecField(child, selection[1:], values)
				if err != nil {
					return nil, err
				}
				vv.SetMapIndex(reflect.ValueOf(selection[0]), reflect.ValueOf(resp))
			} else if vv.Kind() == reflect.Slice {
				resp, err := addSpecField(vv.Index(0).Interface(), selection[1:], values)
				if err != nil {
					return nil, err
				}
				vv.Index(0).Set(reflect.ValueOf(resp))
			}
		}
	} else {
		switch typ := item.(type) {
		//append to existing values
		case interface{}:
			vv := reflect.ValueOf(typ)
			if vv.Kind() == reflect.Map {
				for k, v := range values {
					vv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
				}
				item = typ
			}

		// add new values
		default:
			mapInterface := make(map[string]interface{})
			for k, v := range values {
				if strings.Contains(v, "\"") {
					mapInterface[k] = v
					//TODO Object mapping need to be revised
				} else if strings.Contains(v, "{") {
					strings.ReplaceAll(v, "{", "")
					strings.ReplaceAll(v, "}", "")
					if strings.Trim(v, "") == "" {
						var a interface{}
						mapInterface[k] = a
					} else {
						var i map[string]interface{}
						json.Unmarshal([]byte(v), &i)
						mapInterface[k] = i
					}
				} else {
					mapInterface[k], _ = strconv.Atoi(v)
				}
			}
			item = mapInterface
		}
	}
	return item, nil
}

func updateSpecField(item interface{}, selection []string, value string) (interface{}, error) {

	if len(selection) > 0 {
		switch typ := item.(type) {
		case interface{}:
			vv := reflect.ValueOf(typ)
			if vv.Kind() == reflect.Map {
				child := vv.MapIndex(reflect.ValueOf(selection[0])).Interface()
				resp, err := updateSpecField(child, selection[1:], value)
				if err != nil {
					return nil, err
				}
				vv.SetMapIndex(reflect.ValueOf(selection[0]), reflect.ValueOf(resp))
			} else if vv.Kind() == reflect.Slice {
				resp, err := updateSpecField(vv.Index(0).Interface(), selection[1:], value)
				if err != nil {
					return nil, err
				}
				vv.Index(0).Set(reflect.ValueOf(resp))
			}
		}
	} else {
		var err error
		switch item.(type) {
		case string:
			item = value
		case int, int64:
			item, err = strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
		case interface{}:
			mapValue := make(map[string]string)
			mapInterface := make(map[string]interface{})
			valueToMap(value, mapValue)
			for k, v := range mapValue {
				if strings.Contains(v, "\"") {
					mapInterface[k] = v
					//TODO Object mapping need to be revised
				} else if strings.Contains(v, "{") {
					strings.ReplaceAll(v, "{", "")
					strings.ReplaceAll(v, "}", "")
					if strings.Trim(v, "") == "" {
						var a interface{}
						mapInterface[k] = a
					} else {
						var i map[string]interface{}
						json.Unmarshal([]byte(v), &i)
						mapInterface[k] = i
					}
				} else {
					mapInterface[k], _ = strconv.Atoi(v)
				}
			}
			item = mapInterface
		default:
			item = value
		}

	}
	return item, nil
}

func removeSpecField(item interface{}, selection []string) (interface{}, error) {

	if len(selection) > 0 {
		switch typ := item.(type) {
		case interface{}:
			vv := reflect.ValueOf(typ)
			if vv.Kind() == reflect.Map {
				if len(selection) == 1 {
					vv.SetMapIndex(reflect.ValueOf(selection[0]), reflect.Value{})
					return vv, nil
				} else {
					child := vv.MapIndex(reflect.ValueOf(selection[0])).Interface()
					resp, err := removeSpecField(child, selection[1:])
					if err != nil {
						return nil, err
					}
					vv.SetMapIndex(reflect.ValueOf(selection[0]), reflect.ValueOf(resp))
				}
			} else if vv.Kind() == reflect.Slice {
				resp, err := removeSpecField(vv.Index(0).Interface(), selection[1:])
				if err != nil {
					return nil, err
				}
				vv.Index(0).Set(reflect.ValueOf(resp))
			}
		}
	}
	return item, nil
}

func appendSpecToList(obj interface{}, specTree map[string]string, key string) {
	//tree
	//var isMap = false
	//for k, v := range obj {

	switch typ := obj.(type) {
	case map[string]interface{}:
		//TODO iterative selection is better for ux
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

func updateResources(command *model.Command, actionType string, values map[string]string, items []string) error {
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

			if strings.Contains(actionType, "spec") {
				if strings.Contains(actionType, "add") {
					fmt.Println("Adding Spec ->", strings.Join(command.SelectedSpec, "."), "to ", command.Resource.Name, " ", list[i].GetName())
					field, err := addSpecField(list[i].UnstructuredContent(), command.SelectedSpec, values)
					if err != nil {
						return err
					}
					list[i].SetUnstructuredContent(field.(map[string]interface{}))
				} else if strings.Contains(actionType, "update") {
					fmt.Println("Updating Spec ->", strings.Join(command.SelectedSpec, "."), "to ", command.Resource.Name, " ", list[i].GetName())
					field, err := updateSpecField(list[i].UnstructuredContent(), command.SelectedSpec, values["value"])
					if err != nil {
						return err
					}
					list[i].SetUnstructuredContent(field.(map[string]interface{}))
				} else if strings.Contains(actionType, "remove") {
					fmt.Println("Removing Spec ->", strings.Join(command.SelectedSpec, "."), "to ", command.Resource.Name, " ", list[i].GetName())
					field, err := removeSpecField(list[i].UnstructuredContent(), command.SelectedSpec)
					if err != nil {
						return err
					}
					list[i].SetUnstructuredContent(field.(map[string]interface{}))
				}

			} else {
				if strings.Contains(actionType, "label") {
					valueList = list[i].GetLabels()
				} else if strings.Contains(actionType, "annotation") {
					valueList = list[i].GetAnnotations()
				}
				if valueList == nil {
					valueList = make(map[string]string)
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
					for i := 0; i < len(items); i++ {
						delete(valueList, items[i])
					}
				}
				if strings.Contains(actionType, "label") {
					list[i].SetLabels(valueList)
				} else if strings.Contains(actionType, "annotations") {
					list[i].SetAnnotations(valueList)
				}
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
}
