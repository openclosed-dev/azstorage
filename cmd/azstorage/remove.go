package main

import (
	"github.com/openclosed-dev/azstorage/internal/blob"
	"github.com/spf13/cobra"
)

type removeCommand struct {
	cobra.Command
	accountName   string
	containerName string
	listFile      string
	walkers       int
	processors    int
}

func (c *removeCommand) run(cmd *cobra.Command, args []string) error {
	return blob.RemoveBlobsInList(
		c.accountName, c.containerName, c.listFile, c.walkers, c.processors)
}

func newRemoveCommand() *cobra.Command {

	var command = &removeCommand{
		Command: cobra.Command{
			Use:   "remove",
			Short: "Removes the blobs listed in the specified file",
		},
		listFile: "",
	}

	command.RunE = runWithLog(command.run)

	command.PersistentFlags().StringVar(&command.accountName, "account", "",
		"the name of the storage account")

	command.PersistentFlags().StringVar(&command.containerName, "container", "",
		"the name of the blob container")

	command.PersistentFlags().StringVar(&command.listFile, "list-file", "",
		"path to the file containing directories to delete")

	command.PersistentFlags().IntVar(&command.walkers, "walkers", 4,
		"the number of concurrent directory walkers")

	command.PersistentFlags().IntVar(&command.processors, "processors", 16,
		"the number of concurrent processors deleting found blobs")

	command.MarkPersistentFlagRequired("account")
	command.MarkPersistentFlagRequired("container")
	command.MarkPersistentFlagRequired("list-file")

	return &command.Command
}
