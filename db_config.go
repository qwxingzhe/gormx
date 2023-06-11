package gormx

func getTestMysqlConfig() MysqlConfig {
	return MysqlConfig{
		DBAddress:      "124.223.82.122:3306",
		DBUserName:     "test",
		DBPassword:     "111111",
		DBDatabaseName: "test",
	}
}
