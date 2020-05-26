package service

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"github.com/emreodabas/kubectl-bulk/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func UpdateResources(command *model.Command) error {

	/*
		setAnnotations
		setLabels
		updateUnstructred Content
		 	[specs]

	*/

	selectedPreference := interaction.ShowList(updatePreference)

	switch selectedPreference {
	case "[update] [labels]":
		keys := getLabelKeys(command)
		selection := interaction.ShowList(keys)
		prompt := interaction.Prompt("Define value for label \"%s\" ", selection)
		fmt.Println(selection + "=" + prompt)
		break
	case "[add] [labels]":
		prompt := interaction.Prompt("Define your new label like env=prod ")
		fmt.Println(prompt)
		break
	case "[remove] [labels]":
		keys := getLabelKeys(command)
		selection := interaction.ShowList(keys)
		prompt := interaction.Prompt("Define what value of label will be removed (set * for all values of label) \"%s\" ", selection)
		fmt.Println(selection + "=" + prompt)

		break
	case "[upsert] [labels]":
		break
	case "[update] [annotations]":
		break
	case "[update] [specs]":
		break
	}

	return nil
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

func updateSources(command model.Command) error {
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

		list := command.List
		for i := 0; i < len(list); i++ {
			resourceInterface.Update(&list[i], options)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var updatePreference = []string{
	"[add] [labels]", "[add] [annotations]", "[add] [specs]",
	"[remove] [labels]", "[remove] [annotations]", "[remove] [specs]",
	"[update] [labels]", "[update] [annotations]", "[update] [specs]",
	"[upsert] [labels]", "[upsert] [annotations]", "[upsert] [specs]",
}
