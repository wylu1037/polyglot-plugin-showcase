package request

import "github.com/wylu1037/polyglot-plugin-host-server/app/database/models"

type InstallPluginRequest struct {
	DownloadURL string         `json:"downloadURL" validate:"required,url"`
	Namespace   string         `json:"namespace" validate:"required"`                // 新增：命名空间
	Name        string         `json:"name" validate:"required"`
	Version     string         `json:"version" validate:"required"`
	Type        string         `json:"type" validate:"required"`                    // 修改：从枚举改为 string
	OS          string         `json:"os" validate:"required,oneof=linux darwin windows"` // 新增：操作系统
	Arch        string         `json:"arch" validate:"required,oneof=amd64 arm64"` // 新增：架构
	Description string         `json:"description"`
	Config      models.JSONMap `json:"config"`
	Metadata    models.JSONMap `json:"metadata"`
}

type CallPluginRequest struct {
	ID     uint           `param:"id" validate:"required,gt=0"` // Plugin ID from path parameter
	Method string         `json:"method" validate:"required"`   // Method name from request body
	Params map[string]any `json:"params" validate:"required"`   // Method parameters from request body
}

type ListPluginsRequest struct {
	Namespace string `query:"namespace" validate:"omitempty"` // 新增：按命名空间过滤
	Type      string `query:"type" validate:"omitempty"`      // 修改：移除枚举限制
	Status    string `query:"status" validate:"omitempty,oneof=active inactive disabled error installing"`
	OS        string `query:"os" validate:"omitempty"`        // 新增：按 OS 过滤
	Arch      string `query:"arch" validate:"omitempty"`      // 新增：按架构过滤
}

type PluginIDRequest struct {
	ID uint `param:"id" validate:"required,gt=0"`
}
