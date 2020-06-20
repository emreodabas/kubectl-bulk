package cmd

import (
	"fmt"
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

	// kubectl history and run could be an option
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
	// add cache invalidator command
	var actionArg, resourceArg string
	var command model.Command
	var err error
	if len(args) == 1 {
		actionArg = args[0]
	} else if len(args) == 2 {
		actionArg, resourceArg = args[0], args[1]
	}
	if resourceArg == "" {
		service.ResourceSelection(&command)
	} else {
		command.Resource, err = service.GetResource(resourceArg)
		if err != nil {
			return err
		}
	}
	service.SourceSelection(&command)

	err = service.Filter(&command)
	if err != nil {
		return err
	}

	service.ActionSelection(actionArg, &command)

	switch command.Action.Name {
	case "GET":
		fmt.Println("GET")
		break
	case "UPDATE":
		fmt.Println("UPDATE")
		service.UpdateResources(&command)
		break
	case "DELETE":
		fmt.Println("DELETE")
		break
	}
	return nil
}
