package provider

import (
	"context"
	"testing"

	"go.keploy.io/server/v3/config"
	"go.uber.org/zap"
)

func TestNewServiceProvider(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.New()

	sp := NewServiceProvider(logger, cfg)

	if sp == nil {
		t.Fatal("expected non-nil ServiceProvider")
	}
	if sp.logger != logger {
		t.Errorf("expected logger %v, got %v", logger, sp.logger)
	}
	if sp.cfg != cfg {
		t.Errorf("expected config %v, got %v", cfg, sp.cfg)
	}
}

func TestServiceProvider_GetService_InvalidCommand(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.New()
	sp := NewServiceProvider(logger, cfg)

	svc, err := sp.GetService(context.Background(), "invalid-cmd")
	if err == nil {
		t.Fatal("expected error for invalid command, got nil")
	}
	if err.Error() != "invalid command" {
		t.Errorf("expected error 'invalid command', got '%s'", err.Error())
	}
	if svc != nil {
		t.Errorf("expected nil service, got %v", svc)
	}
}
