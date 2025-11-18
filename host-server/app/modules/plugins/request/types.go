package request

import "github.com/wylu1037/polyglot-plugin-host-server/app/database/models"

type InstallPluginRequest struct {
	DownloadURL string            `json:"downloadURL" validate:"required,url"`
	Name        string            `json:"name" validate:"required"`
	Version     string            `json:"version" validate:"required"`
	Type        models.PluginType `json:"type" validate:"required"`
	Description string            `json:"description"`
	Config      models.JSONMap    `json:"config"`
	Metadata    models.JSONMap    `json:"metadata"`
}

type CallPluginRequest struct {
	Method string         `json:"method" validate:"required"`
	Params map[string]any `json:"params" validate:"required"`
}

type ListPluginsRequest struct {
	Type   string `query:"type" validate:"omitempty,oneof=grpc http"`
	Status string `query:"status" validate:"omitempty,oneof=active inactive"`
}

type PluginIDRequest struct {
	ID uint `param:"id" validate:"required,gt=0"`
}
