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
	var action model.Action
	var resource model.Resource
	var err error
	if len(args) == 1 {
		actionArg = args[0]
	} else if len(args) == 2 {
		actionArg, resourceArg = args[0], args[1]
	}

	if actionArg == "" {
		action = interaction.ShowActionList()
	} else {
		action, err = model.GetAction(actionArg)
		if err != nil {
			return err
		}
	}

	if resourceArg == "" {
		list, err := service.GetResourceList()
		fmt.Println("SIZE", len(list))

		if err != nil {
			return err
		}
		resource = interaction.ShowList(list)
	} else {
		resource, err = service.GetResource(resourceArg)
		if err != nil {
			return err
		}
	}

	fmt.Println("action", action, "resource", resource)
	return nil
}
