package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check server defaults
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected server host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected server port 8080, got %d", cfg.Server.Port)
	}

	// Check plugin defaults
	if cfg.Plugin.Protocol != "netrpc" {
		t.Errorf("Expected plugin protocol 'netrpc', got '%s'", cfg.Plugin.Protocol)
	}

	// Check log defaults
	if cfg.Log.Level != "info" {
		t.Errorf("Expected log level 'info', got '%s'", cfg.Log.Level)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	// Set environment variable
	os.Setenv("PLUGIN_HOST_SERVER_PORT", "9090")
	defer os.Unsetenv("PLUGIN_HOST_SERVER_PORT")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Expected server port 9090 from env var, got %d", cfg.Server.Port)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 99999},
		Plugin: PluginConfig{Protocol: "netrpc"},
		Log:    LogConfig{Level: "info"},
	}

	if err := cfg.Validate(); err == nil {
		t.Error("Expected validation error for invalid port, got nil")
	}
}

func TestValidate_InvalidProtocol(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Port: 8080},
		Plugin: PluginConfig{Protocol: "invalid"},
		Log:    LogConfig{Level: "info"},
	}

	if err := cfg.Validate(); err == nil {
		t.Error("Expected validation error for invalid protocol, got nil")
	}
}

func TestGetServerAddr(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 3000,
		},
	}

	expected := "localhost:3000"
	if addr := cfg.GetServerAddr(); addr != expected {
		t.Errorf("Expected server addr '%s', got '%s'", expected, addr)
	}
}

func TestTimeoutDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("Expected read timeout 30s, got %v", cfg.Server.ReadTimeout)
	}

	if cfg.Plugin.HandshakeTimeout != 10*time.Second {
		t.Errorf("Expected handshake timeout 10s, got %v", cfg.Plugin.HandshakeTimeout)
	}
}
