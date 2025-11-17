package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wylu1037/polyglot-plugin-host-server/app/common"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/request"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/service"
)

type PluginController interface {
	InstallPlugin(c echo.Context) error
	ListPlugins(c echo.Context) error
	GetPlugin(c echo.Context) error
	ActivatePlugin(c echo.Context) error
	DeactivatePlugin(c echo.Context) error
	UninstallPlugin(c echo.Context) error
	CallPlugin(c echo.Context) error
}

type pluginController struct {
	service service.PluginService
}

func NewPluginController(service service.PluginService) PluginController {
	return &pluginController{
		service: service,
	}
}

func (ctrl *pluginController) InstallPlugin(c echo.Context) error {
	var req request.InstallPluginRequest
	if err := c.Bind(&req); err != nil {
		return common.ErrBadRequest.WithDetails("Invalid request body format").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return common.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	plugin, err := ctrl.service.InstallPlugin(&service.InstallPluginRequest{
		DownloadURL: req.DownloadURL,
		Name:        req.Name,
		Version:     req.Version,
		Type:        req.Type,
		Description: req.Description,
		Checksum:    req.Checksum,
		Config:      req.Config,
		Metadata:    req.Metadata,
	})

	if err != nil {
		return common.ErrPluginInstallFailed.WithInternal(err)
	}

	return c.JSON(http.StatusCreated, plugin)
}

func (ctrl *pluginController) ListPlugins(c echo.Context) error {
	filters := make(map[string]any)

	if pluginType := c.QueryParam("type"); pluginType != "" {
		filters["type"] = pluginType
	}

	if status := c.QueryParam("status"); status != "" {
		filters["status"] = status
	}

	plugins, err := ctrl.service.ListPlugins(filters)
	if err != nil {
		return common.ErrInternalServer.WithDetails("Failed to list plugins").WithInternal(err)
	}

	return c.JSON(http.StatusOK, plugins)
}

func (ctrl *pluginController) GetPlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return common.ErrBadRequest.WithDetails("Invalid plugin ID format").WithInternal(err)
	}

	plugin, err := ctrl.service.GetPluginInfo(uint(id))
	if err != nil {
		return common.ErrPluginNotFound.WithInternal(err)
	}

	return c.JSON(http.StatusOK, plugin)
}

func (ctrl *pluginController) ActivatePlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return common.ErrBadRequest.WithDetails("Invalid plugin ID format").WithInternal(err)
	}

	if err := ctrl.service.ActivatePlugin(uint(id)); err != nil {
		return common.ErrPluginActivateFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Plugin activated successfully",
	})
}

func (ctrl *pluginController) DeactivatePlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return common.ErrBadRequest.WithDetails("Invalid plugin ID format").WithInternal(err)
	}

	if err := ctrl.service.DeactivatePlugin(uint(id)); err != nil {
		return common.ErrPluginDeactivateFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Plugin deactivated successfully",
	})
}

func (ctrl *pluginController) UninstallPlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return common.ErrBadRequest.WithDetails("Invalid plugin ID format").WithInternal(err)
	}

	if err := ctrl.service.UninstallPlugin(uint(id)); err != nil {
		return common.ErrPluginUninstallFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Plugin uninstalled successfully",
	})
}

func (ctrl *pluginController) CallPlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return common.ErrBadRequest.WithDetails("Invalid plugin ID format").WithInternal(err)
	}

	var req request.CallPluginRequest
	if err := c.Bind(&req); err != nil {
		return common.ErrBadRequest.WithDetails("Invalid request body format").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return common.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	result, err := ctrl.service.CallPlugin(uint(id), &service.CallPluginRequest{
		Method: req.Method,
		Params: req.Params,
	})

	if err != nil {
		return common.ErrPluginCallFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, result)
}
