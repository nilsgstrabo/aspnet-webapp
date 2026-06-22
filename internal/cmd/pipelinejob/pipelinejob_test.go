package pipelinejob

import (
	"context"
	"strings"
	"testing"

	"github.com/nilsgstrabo/aspnet-webapp/internal/deps"
)

type testClient struct{}

func (t *testClient) Ping(ctx context.Context) error {
	_ = ctx
	return nil
}

type testLogger struct{}

func (t *testLogger) Infof(format string, args ...any) {
	_, _ = format, args
}

func (t *testLogger) Errorf(format string, args ...any) {
	_, _ = format, args
}

func TestPipelineJobShowSkeleton(t *testing.T) {
	cmd := NewCommand(func() *deps.Deps {
		return &deps.Deps{Client: &testClient{}, Logger: &testLogger{}}
	})
	cmd.SetArgs([]string{"show"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected not implemented error")
	}
	if !strings.Contains(err.Error(), "pipelinejob show is not implemented yet") {
		t.Fatalf("unexpected error: %v", err)
	}
}
