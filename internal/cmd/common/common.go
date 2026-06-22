package common

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nilsgstrabo/aspnet-webapp/internal/deps"
)

// DepsProvider provides runtime dependencies created by root command.
type DepsProvider func() *deps.Deps

func RequireDeps(getDeps DepsProvider) (*deps.Deps, error) {
	if getDeps == nil {
		return nil, fmt.Errorf("dependency provider is nil")
	}
	d := getDeps()
	if err := d.Validate(); err != nil {
		return nil, err
	}
	return d, nil
}

func RunNotImplemented(cmd *cobra.Command, getDeps DepsProvider, noun, verb string) error {
	d, err := RequireDeps(getDeps)
	if err != nil {
		return err
	}
	_ = d.Client.Ping(cmd.Context())
	d.Logger.Infof("invoked command: %s %s", noun, verb)
	return fmt.Errorf("%s %s is not implemented yet", noun, verb)
}
