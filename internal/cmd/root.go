package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/application"
	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/common"
	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/config"
	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/deployment"
	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd/pipelinejob"
	"github.com/nilsgstrabo/aspnet-webapp/internal/deps"
)

// DepsFactory builds runtime dependencies from root flags.
type DepsFactory func(endpoint, token, logLevel string) (*deps.Deps, error)

var depsFactory DepsFactory = deps.NewDefaultDeps

func SetDepsFactoryForTest(factory DepsFactory) {
	if factory == nil {
		depsFactory = deps.NewDefaultDeps
		return
	}
	depsFactory = factory
}

func NewRootCmd() *cobra.Command {
	var (
		endpoint    string
		token       string
		logLevel    string
		runtimeDeps *deps.Deps
	)

	getDeps := common.DepsProvider(func() *deps.Deps {
		return runtimeDeps
	})

	cmd := &cobra.Command{
		Use:          "rx",
		Short:        "CLI for noun+verb operations",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_ = args
			d, err := depsFactory(endpoint, token, logLevel)
			if err != nil {
				return err
			}
			runtimeDeps = d
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&endpoint, "endpoint", os.Getenv("RX_ENDPOINT"), "API endpoint")
	cmd.PersistentFlags().StringVar(&token, "token", os.Getenv("RX_TOKEN"), "API token")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Logger verbosity")

	cmd.AddCommand(
		deployment.NewCommand(getDeps),
		application.NewCommand(getDeps),
		config.NewCommand(getDeps),
		pipelinejob.NewCommand(getDeps),
	)

	return cmd
}
