package interaction

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
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
	actions, err := model.Actionlist()
	if err != nil {
		panic("json could not read")
	}
	idx, err := fuzzyfinder.FindMulti(
		actions,
		func(i int) string {
			return actions[i].Name
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Name: %s \nDescription: %s",
				strings.SplitAfter(actions[i].Name, "> ")[0],
				actions[i].Descriptions)
		}))
	if err != nil {
		log.Fatal(err)
	}
	return actions[idx[0]]
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

func ShowUnstructuredList(filterList []string, list []unstructured.Unstructured) string {

	idx, err := fuzzyfinder.FindMulti(
		filterList,
		func(i int) string {
			return filterList[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			var res = "NAME \t\t NAMESPACE \n"
			for _, item := range list {
				res = res + item.GetName() + "  " + item.GetNamespace() + "\n"
			}
			return res
		}))
	if err != nil {
		log.Fatal(err)
	}
	return filterList[idx[0]]
}

//func getPromtValue(info string) string {
//
//
//	return ""
//}
