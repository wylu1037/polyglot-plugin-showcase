package service

func NewPluginStoreService() PluginStoreService {
	return &pluginStoreService{}
}

type PluginStoreService interface{}

type pluginStoreService struct{}
