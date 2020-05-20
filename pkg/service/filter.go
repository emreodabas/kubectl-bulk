package service

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
)

func DoFilter(command *model.Command) error {

	switch command.Filter.Name {
	//TODO label is not working
	case "label":
		prompt := interaction.Prompt("Samples? \n env=production ,type env in  (production, development) ")
		command.Label = prompt
		FetchInstances(command)
		break
	case "multi-select":

	default:
		return fmt.Errorf("Filter option is not implemented yet.")
	}
	return nil
}
