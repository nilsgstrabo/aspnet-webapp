package main

import (
	"os"

	"github.com/nilsgstrabo/aspnet-webapp/internal/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
