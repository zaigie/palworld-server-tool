package main

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zaigie/palworld-server-tool/internal/config"
)

func TestStartupPortOverridePrecedence(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		envValue  string
		envSet    bool
		wantPort  int
		wantIsSet bool
	}{
		{
			name: "config database remains the fallback",
		},
		{
			name:      "environment overrides config database",
			envValue:  "18080",
			envSet:    true,
			wantPort:  18080,
			wantIsSet: true,
		},
		{
			name:      "command line overrides environment",
			args:      []string{"--port", "28080"},
			envValue:  "18080",
			envSet:    true,
			wantPort:  28080,
			wantIsSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookupEnv := func(key string) (string, bool) {
				if key != startupPortEnvironment {
					t.Fatalf("environment key = %q, want %q", key, startupPortEnvironment)
				}
				return tt.envValue, tt.envSet
			}

			got, err := startupPortOverride(tt.args, lookupEnv, io.Discard)
			if err != nil {
				t.Fatalf("resolve startup port: %v", err)
			}
			if (got != nil) != tt.wantIsSet {
				t.Fatalf("override present = %t, want %t", got != nil, tt.wantIsSet)
			}
			if got != nil && got.Port != tt.wantPort {
				t.Fatalf("override port = %d, want %d", got.Port, tt.wantPort)
			}
			if got != nil && tt.args != nil && got.Source != "command_line" {
				t.Fatalf("override source = %q, want command_line", got.Source)
			}
			if got != nil && tt.args == nil && got.Source != "environment" {
				t.Fatalf("override source = %q, want environment", got.Source)
			}
		})
	}
}

func TestStartupPortOverrideRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		envValue string
		envSet   bool
		wantErr  string
	}{
		{
			name:    "command line port below range",
			args:    []string{"--port", "0"},
			wantErr: "port must be between 1 and 65535",
		},
		{
			name:    "command line port above range",
			args:    []string{"--port=65536"},
			wantErr: "port must be between 1 and 65535",
		},
		{
			name:     "environment port is not a number",
			envValue: "not-a-port",
			envSet:   true,
			wantErr:  startupPortEnvironment,
		},
		{
			name:    "unexpected positional argument",
			args:    []string{"18080"},
			wantErr: "unexpected argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lookupEnv := func(string) (string, bool) {
				return tt.envValue, tt.envSet
			}

			_, err := startupPortOverride(tt.args, lookupEnv, io.Discard)
			if err == nil {
				t.Fatalf("resolve startup port unexpectedly succeeded")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func TestStartupPortOverrideIgnoresInvalidEnvironmentWhenFlagIsSet(t *testing.T) {
	lookupEnv := func(string) (string, bool) {
		return "not-a-port", true
	}

	got, err := startupPortOverride([]string{"--port", "18080"}, lookupEnv, io.Discard)
	if err != nil {
		t.Fatalf("resolve startup port: %v", err)
	}
	if got == nil || got.Port != 18080 || got.Source != "command_line" {
		t.Fatalf("override port = %v, want 18080", got)
	}
}

func TestApplyStartupPortOverridePersistsEffectivePort(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "config.db")
	store, err := config.Open(databasePath)
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}

	settings, err := applyStartupPortOverride(store, &startupPort{
		Port:   18080,
		Source: config.WebPortOverrideEnvironment,
	})
	if err != nil {
		t.Fatalf("apply startup port override: %v", err)
	}
	if settings.Web.Port != 18080 {
		t.Fatalf("effective web port = %d, want 18080", settings.Web.Port)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("close config store: %v", err)
	}

	reopened, err := config.Open(databasePath)
	if err != nil {
		t.Fatalf("reopen config store: %v", err)
	}
	persistedWeb := reopened.Config().Web
	if persistedWeb.Port != 18080 || persistedWeb.PortSource != config.WebPortOverrideEnvironment {
		t.Fatalf("persisted web settings = %#v, want port 18080 from environment", persistedWeb)
	}

	settings, err = applyStartupPortOverride(reopened, nil)
	if err != nil {
		t.Fatalf("clear startup port override: %v", err)
	}
	if settings.Web.Port != 18080 || settings.Web.PortSource != config.WebPortOverrideNone {
		t.Fatalf("cleared startup web settings = %#v, want retained port without source", settings.Web)
	}
	if err := reopened.Close(); err != nil {
		t.Fatalf("close config store after clearing source: %v", err)
	}

	verified, err := config.Open(databasePath)
	if err != nil {
		t.Fatalf("reopen config store after clearing source: %v", err)
	}
	defer verified.Close()
	verifiedWeb := verified.Config().Web
	if verifiedWeb.Port != 18080 || verifiedWeb.PortSource != config.WebPortOverrideNone {
		t.Fatalf("verified web settings = %#v, want retained port with persisted empty source", verifiedWeb)
	}
}

func TestApplyStartupPortOverridePreservesExistingCustomPort(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "config.db")
	store, err := config.Open(databasePath)
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	existing := store.Config()
	existing.Web.Port = 19090
	if err := store.Update(existing, ""); err != nil {
		t.Fatalf("save existing custom port: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("close existing config store: %v", err)
	}

	reopened, err := config.Open(databasePath)
	if err != nil {
		t.Fatalf("reopen existing config store: %v", err)
	}
	defer reopened.Close()
	settings, err := applyStartupPortOverride(reopened, nil)
	if err != nil {
		t.Fatalf("apply startup settings without override: %v", err)
	}
	if settings.Web.Port != 19090 || settings.Web.PortSource != config.WebPortOverrideNone {
		t.Fatalf("migrated web settings = %#v, want existing port 19090 without source", settings.Web)
	}
}
