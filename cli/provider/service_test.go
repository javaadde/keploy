package provider

import (
	"context"
	"strings"
	"testing"

	"go.keploy.io/server/v3/config"
	"go.uber.org/zap"
)

func TestNewServiceProvider(t *testing.T) {
	logger := zap.NewNop()
	cfg := config.New()

	sp := NewServiceProvider(logger, cfg)

	if sp == nil {
		t.Fatal("expected NewServiceProvider to return non-nil ServiceProvider")
	}
	if sp.logger != logger {
		t.Errorf("got logger %v, want %v", sp.logger, logger)
	}
	if sp.cfg != cfg {
		t.Errorf("got config %v, want %v", sp.cfg, cfg)
	}
}

func TestServiceProvider_GetService(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name        string
		cmd         string
		disableTele bool
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid command",
			cmd:         "invalid-cmd",
			disableTele: false,
			wantErr:     true,
			errContains: "invalid command",
		},
		{
			name:        "empty command",
			cmd:         "",
			disableTele: true,
			wantErr:     true,
			errContains: "invalid command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.New()
			cfg.DisableTele = tt.disableTele

			sp := NewServiceProvider(logger, cfg)
			svc, err := sp.GetService(context.Background(), tt.cmd)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetService(%q) expected error, got nil", tt.cmd)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetService(%q) error = %v, want error containing %q", tt.cmd, err, tt.errContains)
				}
				if svc != nil {
					t.Errorf("GetService(%q) expected nil service on error, got %v", tt.cmd, svc)
				}
			} else {
				if err != nil {
					t.Fatalf("GetService(%q) unexpected error: %v", tt.cmd, err)
				}
				if svc == nil {
					t.Errorf("GetService(%q) expected non-nil service, got nil", tt.cmd)
				}
			}
		})
	}
}
