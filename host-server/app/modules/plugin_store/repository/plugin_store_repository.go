package repository

func NewPluginStoreRepository() PluginStoreRepository {
	return &pluginStoreRepository{}
}

type PluginStoreRepository interface{}

type pluginStoreRepository struct{}
