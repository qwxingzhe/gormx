package gormx

type MysqlConfig struct {
	DBAddress      string `yaml:"DBAddress" json:"DBAddress"`
	DBUserName     string `yaml:"DBUserName" json:"DBUserName"`
	DBPassword     string `yaml:"DBPassword" json:"DBPassword"`
	DBDatabaseName string `yaml:"DBDatabaseName" json:"DBDatabaseName"`
	DBTablePrefix  string `yaml:"DBTablePrefix" json:"DBTablePrefix"`
	DBMaxOpenConns int    `yaml:"DBMaxOpenConns" json:"DBMaxOpenConns"`
	DBMaxIdleConns int    `yaml:"DBMaxIdleConns" json:"DBMaxIdleConns"`
	DBMaxLifeTime  int64  `yaml:"DBMaxLifeTime" json:"DBMaxLifeTime"`
	SingularTable  bool   `yaml:"SingularTable" json:"SingularTable"`
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

type FilterConfig struct {
	Order    string
	MaxLimit int
}
type ListConfig struct {
	Filter       map[string]interface{}
	FilterConfig FilterConfig
	FormatEvery  bool
}
type UpdateConfig struct {
	GetFilter map[string]interface{}
	EmptyMsg  string
}
type PageResult struct {
	List     interface{} `json:"list"`
	PageInfo PageInfo    `json:"page_info"`
}
