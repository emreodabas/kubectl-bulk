package interaction

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/ktr0731/go-fuzzyfinder"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
	"strconv"
	"strings"
)

func Confirm(promptValue string, args ...interface{}) bool {
	for {
		switch Prompt(promptValue, args...) {
		case "Yes", "yes", "y", "Y":
			return true
		case "No", "no", "n", "N":
			return false
		}
	}
}

func Prompt(prompt string, args ...interface{}) string {
	var s string
	fmt.Printf(prompt+": ", args...)
	_, err := fmt.Scanln(&s)
	if err != nil {
		fmt.Println("Prompt value could not read")
	}
	return s
}

func ShowActionList() model.Action {
	idx, err := fuzzyfinder.FindMulti(
		model.Actionlist,
		func(i int) string {
			return model.Actionlist[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s \nDescription: %s",
				strings.SplitAfter(model.Actionlist[i].Name, "> ")[0],
				model.Actionlist[i].Descriptions)
		}))
	if err != nil {
		log.Fatal(err)
	}
	return model.Actionlist[idx[0]]
}

func ShowResourceList(list []model.Resource) model.Resource {
	idx, err := fuzzyfinder.FindMulti(
		list,
		func(i int) string {
			return list[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s \nKind: %s \nShortNames: %s \nNamespaced: %s \nVerbs: %s ",
				strings.SplitAfter(list[i].Name, "> ")[0],
				list[i].Kind, list[i].ShortName, strconv.FormatBool(list[i].Namespaced), list[i].Verbs)
		}))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return list[idx[0]]
}

func ShowFilterList(list []model.Filter) model.Filter {
	idx, err := fuzzyfinder.FindMulti(
		list,
		func(i int) string {
			return list[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s \nDesc: %s  \nSample: %s ",
				strings.SplitAfter(list[i].Name, "> ")[0],
				list[i].Description, list[i].Sample)
		}))
	if err != nil {
		log.Fatal(err)
	}
	return list[idx[0]]
}

func ShowList(list []string) string {
	idx, err := fuzzyfinder.FindMulti(
		list,
		func(i int) string {
			return list[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s ",
				strings.SplitAfter(list[i], "> ")[0])
		}))
	if err != nil {
		log.Fatal(err)
	}
	return list[idx[0]]
}

func ShowUnstructuredList(list []unstructured.Unstructured, selectList []string) string {

	table := uitable.New()
	table.AddRow("NAME", "NAMESPACE")
	for _, item := range list {
		table.AddRow(item.GetName(), item.GetNamespace())
	}
	if len(list) == 0 {
		table.AddRow(" No Result \n")
	}
	fmt.Println("\n\n")
	fmt.Fprintln(color.Output, table)
	result := ""
	prompt := &survey.Select{
		Message: "Next Step:",
		Options: selectList,
	}
	survey.AskOne(prompt, &result)
	return result
}
