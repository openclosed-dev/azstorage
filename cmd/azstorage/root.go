package main

import "github.com/spf13/cobra"

func newRootCommand() *cobra.Command {

	var command = &cobra.Command{
		Use:          "azstorage [command]",
		Short:        "A command line tool for Azure Storage",
		SilenceUsage: true,
	}

	command.AddCommand(newRemoveCommand())

	return command
}
