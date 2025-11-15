package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/request"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/response"
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
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
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
		return c.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.Response{
		Success: true,
		Message: "Plugin installed successfully",
		Data:    plugin,
	})
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
		return c.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    plugins,
	})
}

func (ctrl *pluginController) GetPlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Invalid plugin ID",
		})
	}

	plugin, err := ctrl.service.GetPluginInfo(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    plugin,
	})
}

func (ctrl *pluginController) ActivatePlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Invalid plugin ID",
		})
	}

	if err := ctrl.service.ActivatePlugin(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Plugin activated successfully",
	})
}

func (ctrl *pluginController) DeactivatePlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Invalid plugin ID",
		})
	}

	if err := ctrl.service.DeactivatePlugin(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Plugin deactivated successfully",
	})
}

func (ctrl *pluginController) UninstallPlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Invalid plugin ID",
		})
	}

	if err := ctrl.service.UninstallPlugin(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Success: true,
		Message: "Plugin uninstalled successfully",
	})
}

func (ctrl *pluginController) CallPlugin(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Invalid plugin ID",
		})
	}

	var req request.CallPluginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	result, err := ctrl.service.CallPlugin(uint(id), &service.CallPluginRequest{
		Method: req.Method,
		Params: req.Params,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Response{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.Response{
		Success: true,
		Data:    result,
	})
}
