package service

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DoFilter(command *model.Command) error {

	switch command.Filter.Name {
	//TODO label is not working
	case "label":
		prompt := interaction.Prompt("Write your label with equation? \n env=production \n env in  (production, development) ")
		command.Label, _ = v1.ParseToLabelSelector(prompt)
		FetchInstances(command)
		break
	case "multi-select":

	default:
		return fmt.Errorf("Filter option is not implemented yet.")
	}
	return nil
}
