package application

import (
	"github.com/spf13/cobra"

	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/common"
)

func NewListCommand(getDeps common.DepsProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			return common.RunNotImplemented(cmd, getDeps, "application", "list")
		},
	}
}
