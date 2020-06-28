package service

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"github.com/emreodabas/kubectl-bulk/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"strings"
)

func DoFilter(command *model.Command) error {

	switch command.Filter.Name {
	//TODO label is not working
	case "label":
		err := promptLabelSelector(command)
		if err != nil {
			return err
		}
	case "field-selector":
		//TODO
		err := promptFieldSelector(command)
		if err != nil {
			return err
		}
	case "grep":
		err := promptGrepSelector(command)
		if err != nil {
			return err
		}
	case "multi-select":
		multiResourceSelect(command)
	case "none":
		fmt.Println("No Filter selected")
	}
	command.SelectedFilters = append(command.SelectedFilters, command.Filter)
	return nil
}

func promptLabelSelector(command *model.Command) error {
	for {
		prompt := interaction.Prompt("Define label selector!! \n Samples: \n -- environment=production,tier=frontend \n -- env in  (production, development) \n")
		command.LabelFilter = prompt
		err := FetchInstances(command)
		if err != nil {
			if strings.Contains(err.Error(), "exit") ||
				strings.Contains(err.Error(), "quit") {
				command.LabelFilter = ""
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

func promptGrepSelector(command *model.Command) error {
	for {
		prompt := interaction.Prompt(" Define text value for searching. add -i for ignoring case")
		command.GrepFilter = append(command.GrepFilter, prompt)
		err := FetchInstances(command)
		if err != nil {
			if strings.Contains(err.Error(), "exit") ||
				strings.Contains(err.Error(), "quit") {
				command.GrepFilter = utils.RemoveItem(command.GrepFilter, len(command.GrepFilter)-1)
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

func multiResourceSelect(command *model.Command) error {

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
		Message:  "Which resources do you select for bulk actions? \n ** [Space] for select-deselect \n ** [Enter] for finalize selection :",
		Options:  listStr,
		PageSize: 20,
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

func Filter(command *model.Command) error {
	var err error
	command.Filter = interaction.ShowFilterList(model.FilterList)
	err = DoFilter(command)
	if err != nil {
		return err
	}
	var selection = []string{"action time", "more filter"}

	if interaction.ShowUnstructuredList(command.List, selection) == "more filter" {
		fmt.Println("CAREFUL--> if you select previous FILTERS it will be OVERRIDED.")
		err = Filter(command)
		if err != nil {
			return err
		}
	}
	return nil
}
