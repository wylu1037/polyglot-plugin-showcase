package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
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

// PluginStore 插件存储表模型
type PluginStore struct {
	ID          uint   `gorm:"primarykey" json:"id"`                                                      // 主键
	Name        string `gorm:"type:varchar(100);uniqueIndex;not null;index:idx_name_version" json:"name"` // 插件名称 (唯一)
	DisplayName string `gorm:"type:varchar(200);not null" json:"display_name"`                            // 显示名称
	Description string `gorm:"type:text" json:"description"`                                              // 插件描述
	Version     string `gorm:"type:varchar(50);not null;index:idx_name_version" json:"version"`           // 插件版本
	Author      string `gorm:"type:varchar(100)" json:"author"`                                           // 作者

	// 插件类型和协议
	Type     PluginType     `gorm:"type:varchar(50);not null;index:idx_type_status" json:"type"` // 插件类型
	Protocol PluginProtocol `gorm:"type:varchar(20);not null" json:"protocol"`                   // 通信协议

	// 文件信息
	BinaryPath string `gorm:"type:varchar(500);not null" json:"binary_path"` // 二进制文件路径
	Checksum   string `gorm:"type:varchar(64)" json:"checksum"`              // SHA256 校验和
	FileSize   int64  `gorm:"default:0" json:"file_size"`                    // 文件大小(字节)

	// 状态信息
	Status    PluginStatus `gorm:"type:varchar(20);not null;default:'inactive';index:idx_type_status,idx_enabled_status" json:"status"` // 插件状态
	Enabled   bool         `gorm:"default:false;index:idx_enabled_status" json:"enabled"`                                               // 是否启用
	AutoStart bool         `gorm:"default:false" json:"auto_start"`                                                                     // 是否自动启动

	// 运行时信息
	ProcessID    *int       `gorm:"default:null" json:"process_id,omitempty"` // 进程 ID (运行时)
	LastStartAt  *time.Time `json:"last_start_at,omitempty"`                  // 最后启动时间
	LastStopAt   *time.Time `json:"last_stop_at,omitempty"`                   // 最后停止时间
	RestartCount int        `gorm:"default:0" json:"restart_count"`           // 重启次数
	ErrorMessage string     `gorm:"type:text" json:"error_message,omitempty"` // 错误信息

	// 配置信息
	Config       JSONMap `gorm:"type:jsonb" json:"config,omitempty"`       // 插件配置 (JSON)
	Capabilities JSONMap `gorm:"type:jsonb" json:"capabilities,omitempty"` // 插件能力 (JSON)
	Metadata     JSONMap `gorm:"type:jsonb" json:"metadata,omitempty"`     // 元数据 (JSON)

	// 协议版本
	ProtocolVersion int    `gorm:"default:1" json:"protocol_version"`                   // 协议版本
	HandshakeCookie string `gorm:"type:varchar(200)" json:"handshake_cookie,omitempty"` // 握手 Cookie

	// 性能统计
	TotalCalls    int64      `gorm:"default:0" json:"total_calls"`     // 总调用次数
	FailedCalls   int64      `gorm:"default:0" json:"failed_calls"`    // 失败调用次数
	AvgResponseMs float64    `gorm:"default:0" json:"avg_response_ms"` // 平均响应时间(毫秒)
	LastCallAt    *time.Time `json:"last_call_at,omitempty"`           // 最后调用时间

	// 依赖和标签
	Dependencies []string `gorm:"type:text[];default:'{}'" json:"dependencies,omitempty"` // 依赖的其他插件
	Tags         []string `gorm:"type:text[];default:'{}'" json:"tags,omitempty"`         // 标签

	// 安装信息
	InstallSource string         `gorm:"type:varchar(200)" json:"install_source,omitempty"` // 安装来源 (local/registry/url)
	InstallAt     *time.Time     `json:"install_at,omitempty"`                              // 安装时间
	InstalledBy   string         `gorm:"type:varchar(100)" json:"installed_by,omitempty"`   // 安装者
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (PluginStore) TableName() string {
	return "plugin_store"
}

// BeforeCreate GORM 钩子 - 创建前
func (p *PluginStore) BeforeCreate(tx *gorm.DB) error {
	p.Status = lo.Ternary(p.Status == "", PluginStatusInactive, p.Status)
	p.Protocol = lo.Ternary(p.Protocol == "", PluginProtocolGRPC, p.Protocol)
	p.ProtocolVersion = lo.Ternary(p.ProtocolVersion == 0, 1, p.ProtocolVersion)
	return nil
}

// IsRunning 检查插件是否正在运行
func (p *PluginStore) IsRunning() bool {
	return p.Status == PluginStatusActive && p.ProcessID != nil && *p.ProcessID > 0
}

// CanStart 检查插件是否可以启动
func (p *PluginStore) CanStart() bool {
	return p.Enabled &&
		p.Status != PluginStatusError &&
		p.Status != PluginStatusInstalling &&
		!p.IsRunning()
}

// MarkAsRunning 标记为运行中
func (p *PluginStore) MarkAsRunning(pid int) {
	now := time.Now()
	p.Status = PluginStatusActive
	p.ProcessID = &pid
	p.LastStartAt = &now
	p.ErrorMessage = ""
}

// MarkAsStopped 标记为已停止
func (p *PluginStore) MarkAsStopped() {
	now := time.Now()
	p.Status = PluginStatusInactive
	p.ProcessID = nil
	p.LastStopAt = &now
}

func (p *PluginStore) MarkAsError(errMsg string) {
	p.Status = PluginStatusError
	p.ErrorMessage = errMsg
	p.ProcessID = nil
}

func (p *PluginStore) IncrementCalls(success bool, responseTimeMs float64) {
	now := time.Now()
	p.TotalCalls++
	p.LastCallAt = &now

	if !success {
		p.FailedCalls++
	}

	p.AvgResponseMs = lo.Ternary(p.AvgResponseMs == 0, responseTimeMs, (p.AvgResponseMs*0.9 + responseTimeMs*0.1))
}
