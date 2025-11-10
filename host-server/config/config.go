package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Plugin   PluginConfig   `mapstructure:"plugin"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	Debug           bool          `mapstructure:"debug"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`           // disable, require, verify-ca, verify-full
	MaxOpenConns    int           `mapstructure:"max_open_conns"`     // Maximum number of open connections
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`     // Maximum number of idle connections
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`  // Maximum lifetime of a connection
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"` // Maximum idle time of a connection
	LogLevel        string        `mapstructure:"log_level"`          // silent, error, warn, info
}

// PluginConfig holds plugin-related configuration
type PluginConfig struct {
	Dir              string        `mapstructure:"dir"`
	Protocol         string        `mapstructure:"protocol"` // "grpc" or "netrpc"
	HandshakeTimeout time.Duration `mapstructure:"handshake_timeout"`
	StartupTimeout   time.Duration `mapstructure:"startup_timeout"`
	AutoLoad         []string      `mapstructure:"auto_load"` // Plugin names to load on startup
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, console
	Output string `mapstructure:"output"` // stdout, file path
}

// Load loads configuration from file and environment variables
// Priority: env vars > config file > defaults
func Load(configPath string) (*Config, error) {
	v := viper.New()

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	// Enable environment variable reading
	v.SetEnvPrefix("PLUGIN_HOST")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found; using defaults and env vars
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	// Validate database config
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	validSSLModes := map[string]bool{
		"disable":     true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if !validSSLModes[c.Database.SSLMode] {
		return fmt.Errorf("invalid database ssl_mode: %s", c.Database.SSLMode)
	}

	// Validate plugin protocol
	if c.Plugin.Protocol != "grpc" && c.Plugin.Protocol != "netrpc" {
		return fmt.Errorf("invalid plugin protocol: %s (must be 'grpc' or 'netrpc')", c.Plugin.Protocol)
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("invalid log level: %s", c.Log.Level)
	}

	return nil
}

// GetServerAddr returns the server address in "host:port" format
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// GetDatabaseDSN returns the PostgreSQL connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Database,
		c.Database.SSLMode,
	)
}
