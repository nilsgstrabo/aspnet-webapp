package application

import (
	"github.com/spf13/cobra"

	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/common"
)

func NewCommand(getDeps common.DepsProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "application",
		Short: "Manage applications",
	}

	cmd.AddCommand(
		NewShowCommand(getDeps),
		NewListCommand(getDeps),
		NewCreateCommand(getDeps),
		NewDeleteCommand(getDeps),
	)

	return cmd
}
