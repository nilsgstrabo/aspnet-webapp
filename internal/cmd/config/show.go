package config

import (
	"github.com/spf13/cobra"

	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/common"
)

func NewShowCommand(getDeps common.DepsProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show config",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			return common.RunNotImplemented(cmd, getDeps, "config", "show")
		},
	}
}
