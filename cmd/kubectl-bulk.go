package main

import (
	"os"

	"github.com/emreodabas/kubectl-bulk/pkg/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {

	bulkCommand := cmd.NewBulkCommand(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := bulkCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
