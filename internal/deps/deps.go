package deps

import (
	"fmt"

	"github.com/nilsgstrabo/aspnet-webapp/internal/client"
	"github.com/nilsgstrabo/aspnet-webapp/internal/logger"
)

// Deps contains all runtime dependencies injected from root into subcommands.
type Deps struct {
	Client client.Client
	Logger logger.Logger
}

func (d *Deps) Validate() error {
	if d == nil {
		return fmt.Errorf("dependencies are nil")
	}
	if d.Client == nil {
		return fmt.Errorf("client dependency is nil")
	}
	if d.Logger == nil {
		return fmt.Errorf("logger dependency is nil")
	}
	return nil
}
