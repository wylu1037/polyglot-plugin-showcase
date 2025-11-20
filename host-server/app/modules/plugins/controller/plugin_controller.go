package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/request"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/service"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/errors"
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

// InstallPlugin godoc
// @Summary      Install a new plugin
// @Description  Install a plugin from a download URL
// @Tags         Plugins
// @Accept       json
// @Produce      json
// @Param        request body request.InstallPluginRequest true "Plugin installation request"
// @Success      201 {object} models.Plugin
// @Failure      400 {object} errors.AppError
// @Failure      500 {object} errors.AppError
// @Router       /api/plugins/install [post]
func (ctrl *pluginController) InstallPlugin(c echo.Context) error {
	var req request.InstallPluginRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid request body format").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	plugin, err := ctrl.service.InstallPlugin(&req)
	if err != nil {
		return errors.ErrPluginInstallFailed.WithInternal(err)
	}

	return c.JSON(http.StatusCreated, plugin)
}

// ListPlugins godoc
// @Summary      List all plugins
// @Description  Get a list of all installed plugins with optional filters
// @Tags         Plugins
// @Accept       json
// @Produce      json
// @Param        namespace query string false "Filter by namespace"
// @Param        type      query string false "Filter by plugin type"
// @Param        status    query string false "Filter by plugin status" Enums(active, inactive, disabled, error, installing)
// @Param        os        query string false "Filter by operating system" Enums(linux, darwin, windows)
// @Param        arch      query string false "Filter by architecture" Enums(amd64, arm64)
// @Success      200 {array} models.Plugin
// @Failure      400 {object} errors.AppError
// @Failure      500 {object} errors.AppError
// @Router       /api/plugins [get]
func (ctrl *pluginController) ListPlugins(c echo.Context) error {
	var req request.ListPluginsRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid query parameters").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	plugins, err := ctrl.service.ListPlugins(&req)
	if err != nil {
		return errors.ErrInternalServer.WithDetails("Failed to list plugins").WithInternal(err)
	}

	return c.JSON(http.StatusOK, plugins)
}

// GetPlugin godoc
// @Summary      Get plugin details
// @Description  Get detailed information about a specific plugin by ID
// @Tags         Plugins
// @Accept       json
// @Produce      json
// @Param        id path int true "Plugin ID" minimum(1)
// @Success      200 {object} models.Plugin
// @Failure      400 {object} errors.AppError
// @Failure      404 {object} errors.AppError
// @Router       /api/plugins/{id} [get]
func (ctrl *pluginController) GetPlugin(c echo.Context) error {
	var req request.PluginIDRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid plugin ID").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	plugin, err := ctrl.service.GetPluginInfo(req.ID)
	if err != nil {
		return errors.ErrPluginNotFound.WithInternal(err)
	}

	return c.JSON(http.StatusOK, plugin)
}

// ActivatePlugin godoc
// @Summary      Activate a plugin
// @Description  Activate a previously installed plugin
// @Tags         Plugins
// @Accept       json
// @Produce      json
// @Param        id path int true "Plugin ID" minimum(1)
// @Success      200 {object} map[string]string
// @Failure      400 {object} errors.AppError
// @Failure      500 {object} errors.AppError
// @Router       /api/plugins/{id}/activate [post]
func (ctrl *pluginController) ActivatePlugin(c echo.Context) error {
	var req request.PluginIDRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid plugin ID").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	if err := ctrl.service.ActivatePlugin(req.ID); err != nil {
		return errors.ErrPluginActivateFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Plugin activated successfully",
	})
}

// DeactivatePlugin godoc
// @Summary      Deactivate a plugin
// @Description  Deactivate an active plugin
// @Tags         Plugins
// @Accept       json
// @Produce      json
// @Param        id path int true "Plugin ID" minimum(1)
// @Success      200 {object} map[string]string
// @Failure      400 {object} errors.AppError
// @Failure      500 {object} errors.AppError
// @Router       /api/plugins/{id}/deactivate [post]
func (ctrl *pluginController) DeactivatePlugin(c echo.Context) error {
	var req request.PluginIDRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid plugin ID").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	if err := ctrl.service.DeactivatePlugin(req.ID); err != nil {
		return errors.ErrPluginDeactivateFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Plugin deactivated successfully",
	})
}

// UninstallPlugin godoc
// @Summary      Uninstall a plugin
// @Description  Remove a plugin from the system
// @Tags         Plugins
// @Accept       json
// @Produce      json
// @Param        id path int true "Plugin ID" minimum(1)
// @Success      200 {object} map[string]string
// @Failure      400 {object} errors.AppError
// @Failure      500 {object} errors.AppError
// @Router       /api/plugins/{id} [delete]
func (ctrl *pluginController) UninstallPlugin(c echo.Context) error {
	var req request.PluginIDRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid plugin ID").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	if err := ctrl.service.UninstallPlugin(req.ID); err != nil {
		return errors.ErrPluginUninstallFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Plugin uninstalled successfully",
	})
}

// CallPlugin godoc
// @Summary      Call a plugin method
// @Description  Execute a specific method on an active plugin
// @Tags         Plugins
// @Accept       json
// @Produce      json
// @Param        id path int true "Plugin ID" minimum(1)
// @Param        request body request.CallPluginRequest true "Plugin call request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} errors.AppError
// @Failure      500 {object} errors.AppError
// @Router       /api/plugins/{id}/call [post]
func (ctrl *pluginController) CallPlugin(c echo.Context) error {
	var req request.CallPluginRequest
	if err := c.Bind(&req); err != nil {
		return errors.ErrBadRequest.WithDetails("Invalid request format").WithInternal(err)
	}

	if err := c.Validate(&req); err != nil {
		return errors.ErrValidationFailed.WithDetails(err.Error()).WithInternal(err)
	}

	result, err := ctrl.service.CallPlugin(req.ID, &req)
	if err != nil {
		return errors.ErrPluginCallFailed.WithInternal(err)
	}

	return c.JSON(http.StatusOK, result)
}
