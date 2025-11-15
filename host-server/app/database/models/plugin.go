package models

import (
	"database/sql/driver"
	"encoding/json"
)

type PluginStatus string

const (
	PluginStatusActive     PluginStatus = "active"     // 激活
	PluginStatusInactive   PluginStatus = "inactive"   // 未激活
	PluginStatusDisabled   PluginStatus = "disabled"   // 已禁用
	PluginStatusError      PluginStatus = "error"      // 错误状态
	PluginStatusInstalling PluginStatus = "installing" // 安装中
)

// PluginType 插件类型
type PluginType string

const (
	PluginTypeDesensitization PluginType = "desensitization" // 数据脱敏
	PluginTypeEncryption      PluginType = "encryption"      // 加密
	PluginTypeValidation      PluginType = "validation"      // 验证
	PluginTypeTransform       PluginType = "transform"       // 数据转换
	PluginTypeCustom          PluginType = "custom"          // 自定义
)

type PluginProtocol string

const (
	PluginProtocolGRPC   PluginProtocol = "grpc"    // gRPC 协议
	PluginProtocolNetRPC PluginProtocol = "net-rpc" // net/rpc 协议
)

// JSONMap 用于存储 JSON 格式的配置
type JSONMap map[string]any

// Scan 实现 sql.Scanner 接口
func (j *JSONMap) Scan(value any) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Plugin 插件表模型
type Plugin struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	Name            string         `gorm:"type:varchar(100);not null;uniqueIndex:idx_name_version" json:"name"`
	Version         string         `gorm:"type:varchar(50);not null;uniqueIndex:idx_name_version" json:"version"`
	Type            PluginType     `gorm:"type:varchar(50);not null;index" json:"type"`
	Description     string         `gorm:"type:text" json:"description"`
	Status          PluginStatus   `gorm:"type:varchar(20);not null;default:'inactive';index" json:"status"`
	BinaryPath      string         `gorm:"type:varchar(500);not null" json:"binary_path"`
	DownloadURL     string         `gorm:"type:varchar(500)" json:"download_url"`
	Protocol        PluginProtocol `gorm:"type:varchar(20);not null;default:'grpc'" json:"protocol"`
	ProtocolVersion int            `gorm:"not null;default:1" json:"protocol_version"`
	Checksum        string         `gorm:"type:varchar(64)" json:"checksum"`
	ChecksumType    string         `gorm:"type:varchar(20);default:'sha256'" json:"checksum_type"`
	Config          JSONMap        `gorm:"type:jsonb" json:"config"`
	Metadata        JSONMap        `gorm:"type:jsonb" json:"metadata"`
	CreatedAt       int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       int64          `gorm:"autoUpdateTime" json:"updated_at"`
	LastUsedAt      *int64         `json:"last_used_at"`
}

func (Plugin) TableName() string {
	return "plugins"
}
