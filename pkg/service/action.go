package service

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"strings"
)

func ActionSelection(arg string, command *model.Command) error {
	var err error
	//action time
	if arg == "" {
		command.Action = interaction.ShowActionList()
	} else {
		command.Action, err = getAction(arg)
		fmt.Println("---- SELECTED ACTION ----\n---- ", command.Action.Name, " ----\n-----------------------------------------")
		if err != nil {
			return err
		}
	}

	return nil
}

func getAction(action string) (model.Action, error) {
	var actionMessage string
	action = strings.ToLower(action)
	for i := 0; i < len(model.Actionlist); i++ {
		actionMessage += model.Actionlist[i].Name + " "
		if strings.ToLower(model.Actionlist[i].Name) == action {
			return model.Actionlist[i], nil
		}
	}
	return model.Action{}, fmt.Errorf(action + " is not a valid action. Choose wisely one of them -> " + actionMessage)
}
