package deps

import (
	"github.com/nilsgstrabo/aspnet-webapp/internal/client"
	"github.com/nilsgstrabo/aspnet-webapp/internal/logger"
)

func NewDefaultDeps(endpoint, token, logLevel string) (*Deps, error) {
	return &Deps{
		Client: client.New(endpoint, token),
		Logger: logger.New(logLevel),
	}, nil
}
