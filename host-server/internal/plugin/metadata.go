package plugin

import (
	"encoding/json"
	"fmt"

	"github.com/wylu1037/polyglot-plugin-host-server/internal/validator"
)

// PluginMetadata contains structured metadata about a plugin
type PluginMetadata struct {
	Name            string           `json:"name" validate:"required"`
	Version         string           `json:"version" validate:"required"`
	Type            string           `json:"type" validate:"required"` // 插件自己声明的类型
	Description     string           `json:"description"`
	Author          string           `json:"author"`
	ProtocolVersion int32            `json:"protocol_version" validate:"required,gt=0"`
	Capabilities    []string         `json:"capabilities"`
	Methods         []MethodMetadata `json:"methods" validate:"dive"`
	Dependencies    []Dependency     `json:"dependencies,omitempty" validate:"dive"`
}

// MethodMetadata describes a method provided by the plugin
type MethodMetadata struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Parameters  map[string]ParamSchema `json:"parameters" validate:"dive"`
	Returns     ReturnSchema           `json:"returns"`
}

// ParamSchema describes a method parameter
type ParamSchema struct {
	Type        string `json:"type" validate:"required"` // string, int, bool, object, array
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     any    `json:"default,omitempty"`
}

// ReturnSchema describes method return value
type ReturnSchema struct {
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
}

// Dependency describes a plugin dependency
type Dependency struct {
	Name    string `json:"name" validate:"required"`
	Version string `json:"version" validate:"required"` // 支持语义化版本范围，如 ">=1.0.0 <2.0.0"
	Type    string `json:"type" validate:"required"`    // plugin, library, service
}

func (m *PluginMetadata) Validate() error {
	if err := validator.Validate(m); err != nil {
		return err
	}
	return nil
}

// ToJSON converts metadata to JSON string
func (m *PluginMetadata) ToJSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}
	return string(data), nil
}

// FromJSON parses metadata from JSON string
func FromJSON(data string) (*PluginMetadata, error) {
	var metadata PluginMetadata
	if err := json.Unmarshal([]byte(data), &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	if err := metadata.Validate(); err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}

	return &metadata, nil
}

// PluginPath represents a structured plugin path
type PluginPath struct {
	BaseDir   string `validate:"required"`
	Namespace string `validate:"required"`
	Type      string `validate:"required"`
	Name      string `validate:"required"`
	Version   string `validate:"required"`
	OS        string `validate:"required,oneof=linux darwin windows"`
	Arch      string `validate:"required,oneof=amd64 arm64"`
}

// String returns the full path
func (p *PluginPath) String() string {
	return fmt.Sprintf("%s/%s/%s/%s/%s/%s_%s/plugin",
		p.BaseDir,
		p.Namespace,
		p.Type,
		p.Name,
		p.Version,
		p.OS,
		p.Arch,
	)
}

// Validate validates the plugin path
func (p *PluginPath) Validate() error {
	if err := validator.Validate(p); err != nil {
		return err
	}
	return nil
}

// ParsePluginPath parses a plugin path into structured components
func ParsePluginPath(basePath, fullPath string) (*PluginPath, error) {
	return nil, fmt.Errorf("not implemented yet")
}
