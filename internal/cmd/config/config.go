package config

import (
	"github.com/spf13/cobra"

	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/common"
)

func NewCommand(getDeps common.DepsProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configs",
	}

	cmd.AddCommand(
		NewShowCommand(getDeps),
		NewListCommand(getDeps),
		NewCreateCommand(getDeps),
		NewDeleteCommand(getDeps),
	)

	return cmd
}
