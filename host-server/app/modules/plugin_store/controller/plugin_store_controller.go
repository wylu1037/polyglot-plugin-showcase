package controller

func NewPluginStoreController() PluginStoreController {
	return &pluginStoreController{}
}

type PluginStoreController interface{}

type pluginStoreController struct{}
