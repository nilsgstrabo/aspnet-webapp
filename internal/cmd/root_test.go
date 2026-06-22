package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/nilsgstrabo/aspnet-webapp/internal/deps"
)

type rootTestClient struct{}

func (t *rootTestClient) Ping(ctx context.Context) error {
	_ = ctx
	return nil
}

type rootTestLogger struct{}

func (t *rootTestLogger) Infof(format string, args ...any) {
	_, _ = format, args
}

func (t *rootTestLogger) Errorf(format string, args ...any) {
	_, _ = format, args
}

func TestRootHasExpectedNouns(t *testing.T) {
	r := NewRootCmd()

	commands := r.Commands()
	wanted := map[string]bool{
		"deployment":  false,
		"application": false,
		"config":      false,
		"pipelinejob": false,
	}

	for _, c := range commands {
		if _, ok := wanted[c.Name()]; ok {
			wanted[c.Name()] = true
		}
	}

	for noun, found := range wanted {
		if !found {
			t.Fatalf("missing noun command: %s", noun)
		}
	}
}

func TestRootUsesInjectedFactory(t *testing.T) {
	defer SetDepsFactoryForTest(nil)
	SetDepsFactoryForTest(func(endpoint, token, logLevel string) (*deps.Deps, error) {
		_, _, _ = endpoint, token, logLevel
		return &deps.Deps{Client: &rootTestClient{}, Logger: &rootTestLogger{}}, nil
	})

	r := NewRootCmd()
	r.SetArgs([]string{"deployment", "show"})

	err := r.Execute()
	if err == nil {
		t.Fatal("expected not implemented error")
	}
	if !strings.Contains(err.Error(), "deployment show is not implemented yet") {
		t.Fatalf("unexpected error: %v", err)
	}
}
