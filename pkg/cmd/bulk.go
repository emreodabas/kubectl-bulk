package cmd

import (
	"fmt"
	"github.com/emreodabas/kubectl-bulk/pkg/interaction"
	"github.com/emreodabas/kubectl-bulk/pkg/model"
	"github.com/emreodabas/kubectl-bulk/pkg/service"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	bulkExample = "kubectl bulk\n" +
		"kubectl bulk node \n" +
		"kubectl bulk pod --all-namespaces \n" +
		"kubectl bulk daemonsets"
)

func NewBulkCommand(streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kubectl bulk GET|UPDATE|DELETE|LIST|REMOVE| <any resource>",
		SilenceUsage: true,
		Short:        "Bulk operations on any Kubernetes resources",
		Example:      bulkExample,
		Args:         cobra.MaximumNArgs(2),
		RunE:         run,
	}

	return cmd
}

func run(_ *cobra.Command, args []string) error {
	//TODO logging and test need to be done
	var actionArg, resourceArg string
	var command model.Command
	var err error
	if len(args) == 1 {
		actionArg = args[0]
	} else if len(args) == 2 {
		actionArg, resourceArg = args[0], args[1]
	}

	if resourceArg == "" {
		list, err := service.GetResourceList()
		fmt.Println("SIZE", len(list))

		if err != nil {
			return err
		}
		command.Resource = interaction.ShowResourceList(list)
	} else {
		command.Resource, err = service.GetResource(resourceArg)
		if err != nil {
			return err
		}
	}
	//fmt.Println("action", command.Action.Name, "resource", command.Resource.Name)

	sourceSelection(&command)

	err = Filter(&command)
	if err != nil {
		return err
	}
	//action time
	if actionArg == "" {
		command.Action = interaction.ShowActionList()
	} else {
		command.Action, err = model.GetAction(actionArg)
		if err != nil {
			return err
		}
	}

	return nil
}

func sourceSelection(command *model.Command) error {
	// filter or multi selection could be ask to user
	var err error
	if command.Resource.Namespaced {
		namespaces, err := service.GetNamespaces()
		command.Namespace = interaction.ShowList(namespaces)
		if err != nil {

			return fmt.Errorf("Namespace list could not fetch")
		}
		err = service.FetchInstances(command)

	} else {
		err = service.FetchInstances(command)
	}
	if err != nil {
		return err
	}
	return err

}

func Filter(command *model.Command) error {
	filterlist, err := model.Filterlist(command.Resource.Verbs)
	if err != nil {
		return fmt.Errorf("filter list could not fetched", err)
	}

	command.Filter = interaction.ShowFilterList(filterlist)
	err = service.DoFilter(command)
	if err != nil {
		return err
	}
	var selection = []string{"more filter", "action time"}

	if interaction.ShowUnstructuredList(command.List, selection) == "more filter" {
		fmt.Println("CAREFUL--> if you select previous FILTERS it will be OVERRIDED.")
		err = Filter(command)
		if err != nil {
			return err
		}
	}
	return nil
}
