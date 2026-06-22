package deployment

import (
	"github.com/spf13/cobra"

	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/common"
)

func NewDeleteCommand(getDeps common.DepsProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			return common.RunNotImplemented(cmd, getDeps, "deployment", "delete")
		},
	}
}
