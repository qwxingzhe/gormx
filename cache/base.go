package cache

type CacheInterface interface {
	SetCache(key string, value string) error
	GetCache(key string) string
	DelCache(key string) error
}
