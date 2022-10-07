package cache

// Interface 数据库访问缓存层
type Interface interface {
	// SetCache 设置缓存
	SetCache(key string, value string) error
	// GetCache 获取缓存
	GetCache(key string) string
	// DelCache 删除缓存
	DelCache(key string) error
}
