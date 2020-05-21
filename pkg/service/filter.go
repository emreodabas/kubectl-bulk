package service

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"strings"
)

func DoFilter(command *model.Command) error {

	switch command.Filter.Name {
	//TODO label is not working
	case "label":
		err := promptLabel(command)
		if err != nil {
			return err
		}
		break
	case "field-selector":
		//TODO
		err := promptFieldSelector(command)
		if err != nil {
			return err
		}
		break
	case "multi-select":
		promptMultiSelect(command)
		break
	default:
		return fmt.Errorf("Filter option is not implemented yet.")
	}
	command.SelectedFilters = append(command.SelectedFilters, command.Filter)
	return nil
}

func promptLabel(command *model.Command) error {
	for {
		prompt := interaction.Prompt("Define label selector!! \n Samples: \n -- environment=production,tier=frontend \n -- env in  (production, development) \n")
		command.Label = prompt
		err := FetchInstances(command)
		if err != nil {
			if strings.Contains(err.Error(), "exit") ||
				strings.Contains(err.Error(), "quit") {
				command.Label = ""
				return nil
			} else {
				fmt.Println("Error occured", err, "\n please specify a valid label or you could exit with write [exit] or [quit]")
			}
		} else {
			break
		}
	}
	return nil
}

func promptFieldSelector(command *model.Command) error {
	for {
		prompt := interaction.Prompt("Define field selector!! \n Samples: \n -- metadata.namespace!=default \n -- metadata.name!=test \n")
		command.FieldSelector = prompt
		err := FetchInstances(command)
		if err != nil {
			if strings.Contains(err.Error(), "exit") ||
				strings.Contains(err.Error(), "quit") {
				command.FieldSelector = ""
				return nil
			} else {
				fmt.Println("Error occured", err, "\n please specify a valid field selector or you could exit with write [exit] or [quit] ")
			}
		} else {
			break
		}
	}
	return nil
}

func promptMultiSelect(command *model.Command) error {

	var resultStr []string
	var result []unstructured.Unstructured
	if command.List == nil {
		FetchInstances(command)
	}
	var listStr = make([]string, len(command.List))

	for i, item := range command.List {
		listStr[i] = strconv.Itoa(i) + "-|-" + item.GetName() + "-|-" + item.GetNamespace()
	}

	prompt := &survey.MultiSelect{
		Message: "Which resources do you select for bulk actions? \n ** [Space] for select-deselect \n ** [Enter] for finalize selection :",
		Options: listStr,
	}
	survey.AskOne(prompt, &resultStr)

	for i := 0; i < len(resultStr); i++ {
		split := strings.Split(resultStr[i], "-|-")
		atoi, _ := strconv.Atoi(split[0])
		result = append(result, command.List[atoi])
	}
	command.List = result
	return nil
}
