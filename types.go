package gormx

type MysqlConfig struct {
	DBAddress      string
	DBUserName     string
	DBPassword     string
	DBDatabaseName string
	DBTablePrefix  string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBMaxLifeTime  int64
	SingularTable  bool
}
type RedisConfig struct {
	Addr     string
	DB       int
	Password string
}
type PageInfo struct {
	// 当前页面
	CurrentPage int `json:"current_page"`
	// 页面记录条数
	PageSize int `json:"page_size"`
	// 总记录数
	Total int64 `json:"total"`
}
