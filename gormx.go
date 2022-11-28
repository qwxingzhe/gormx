package gormx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/qwxingzhe/gormx/cache"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type (
	BaseModel struct {
		Id        int64 `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt
	}
	BaseModelWithoutDeletedAt struct {
		Id        int64 `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	BaseModelWithoutDeletedCreatedAt struct {
		Id        int64 `gorm:"primary_key"`
		UpdatedAt time.Time
	}
	BaseModelWithoutDeletedUpdatedAt struct {
		Id        int64 `gorm:"primary_key"`
		CreatedAt time.Time
	}
	BaseModelOnlyId struct {
		Id int64 `gorm:"primary_key"`
	}
	IBaseModel interface {
		Format() interface{}
	}
	Orm struct {
		// gorm实例
		GormDb *gorm.DB
		// 默认排序方式
		DefaultOrder string
		// 缓存引擎实现实例
		cache cache.Interface
	}
)

var instance Orm
var once sync.Once

// InitOrm 初始化orm单例对象
func InitOrm(mysqlConfig MysqlConfig, cache cache.Interface) Orm {
	once.Do(func() {
		var gormDb *gorm.DB
		var err error
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlConfig.DBUserName,
			mysqlConfig.DBPassword,
			mysqlConfig.DBAddress,
			mysqlConfig.DBDatabaseName)

		gormDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   mysqlConfig.DBTablePrefix,
				SingularTable: mysqlConfig.SingularTable,
			},
		})
		if err != nil {
			panic("连接数据库失败")
		}

		sqlDB, errDb := gormDb.DB()

		if errDb == nil {
			// SetMaxIdleConns 设置空闲连接池中连接的最大数量
			sqlDB.SetMaxIdleConns(mysqlConfig.DBMaxIdleConns) // 10
			// SetMaxOpenConns 设置打开数据库连接的最大数量。
			sqlDB.SetMaxOpenConns(mysqlConfig.DBMaxOpenConns) // 100
			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			sqlDB.SetConnMaxLifetime(time.Second * time.Duration(mysqlConfig.DBMaxLifeTime))
		}
		instance = Orm{
			GormDb:       gormDb,
			DefaultOrder: "id desc",
			cache:        cache,
		}
	})
	return instance
}
func GetOrm() Orm {
	if instance.GormDb == nil {
		panic("GormDb 对象未被初始化，请检查调用方式/流程是否合法")
	}
	return instance
}

// 涉及缓存的方法区 start
//+---------------------------------------------------------------------

// GetInfo 获取详情
func (o Orm) GetInfo(info interface{}, id int64) error {
	key := ""
	// 查询缓存是否存在
	if o.cache != nil {
		key = o.getCacheKey(info, id)
		cache := o.getCache(key)

		if cache != "" { // 存在，返回缓存
			json.Unmarshal([]byte(cache), info)
			return nil
		}
	}

	tx := o.GormDb.Where("id = ?", id).First(info)
	if tx.Error != nil { // 未查询到数据
		fmt.Println("First Err:", tx.Error)
		return tx.Error
	}

	// 写入缓存
	if o.cache != nil {
		err := o.setCache(key, info)
		if err != nil {
			fmt.Println("setCache Err:", err.Error())
			return err
		}
	}

	return nil
}

// Delete 删除记录
func (o Orm) Delete(info interface{}) {

	tx := o.GormDb.Delete(info)
	fmt.Println("tx::::", tx.Error)

	// 删除缓存
	o.delCacheByInfo(info)

}

// Save 保存记录
func (o Orm) Save(info interface{}) error {
	err := o.GormDb.Save(info).Error
	if err == nil {
		// 删除缓存
		o.delCacheByInfo(info)
	}
	return err
}

// 涉及缓存的方法区 end
//+---------------------------------------------------------------------

// FormatList 格式化列表
func FormatList(originList interface{}, formatFunc interface{}) []interface{} {
	formatMethod := "Format"
	if formatFunc != nil {
		formatMethod = cast.ToString(formatFunc)
	}

	fmt.Println(reflect.TypeOf(originList).Kind())
	switch reflect.TypeOf(originList).Kind() {
	case reflect.Slice, reflect.Array:
		list := reflect.ValueOf(originList)
		count := list.Len()
		result := make([]interface{}, count)
		for i := 0; i < count; i++ {
			mv := list.Index(i).MethodByName(formatMethod)
			result[i] = mv.Call(nil)[0].Interface()
		}
		return result
	}
	return nil
}

// KeyList 获取list中指定key数组
func KeyList(originList interface{}, key string) []interface{} {
	switch reflect.TypeOf(originList).Kind() {
	case reflect.Slice, reflect.Array:
		list := reflect.ValueOf(originList)
		count := list.Len()
		result := make([]interface{}, count)
		for i := 0; i < count; i++ {
			result[i] = list.Index(i).FieldByName(key).Interface()
		}
		return result
	}
	return nil
}

// Begin 开始事务
func (o Orm) Begin() Orm {
	o.GormDb = o.GormDb.Begin()
	return o
}

// Rollback 回滚事务
func (o Orm) Rollback() {
	o.GormDb.Rollback()
}

// Commit 提交事务
func (o Orm) Commit() {
	o.GormDb.Commit()
}

// GetInfoByKey 通过指定的key、value获取数据
func (o Orm) GetInfoByKey(info interface{}, key string, value interface{}) {
	o.GormDb.Where(key+" = ?", value).First(info)
}

// GetInfoByKeyLike 通过模糊匹配获取数据
func (o Orm) GetInfoByKeyLike(info interface{}, key string, value interface{}) {
	o.GormDb.Where(key+" like ?", value).First(info)
}

// First 查询一条记录
func (o Orm) First(info interface{}, maps map[string]interface{}) {
	for s, i := range maps {
		o.GormDb = o.GormDb.Where(s, i)
	}
	o.GormDb.First(info)
}

func (o Orm) Order(v interface{}) Orm {
	o.GormDb = o.GormDb.Order(v)
	return o
}
func (o Orm) Limit(limit int) Orm {
	o.GormDb = o.GormDb.Limit(limit)
	return o
}
func (o Orm) Debug() Orm {
	o.GormDb = o.GormDb.Debug()
	return o
}
func (o Orm) In(key string, values []interface{}) Orm {
	o.GormDb = o.GormDb.Where("? = (?)", key, values)
	return o
}
func (o Orm) Select(fields string) Orm {
	o.GormDb = o.GormDb.Select(fields)
	return o
}

// 监测数据是否存在，存在则异常
//func (o Orm) CheckExist(info interface{}, maps map[string]interface{},msg string) error{
//	o.First(info,maps)
//	//if info.Id
//}

func (o Orm) Unscoped() Orm {
	o.GormDb = o.GormDb.Unscoped()
	return o
}

func (o Orm) Where(maps map[string]interface{}) Orm {
	for s, i := range maps {
		o.GormDb = o.GormDb.Where(s, i)
	}
	return o
}

func (o Orm) Original(f func(*gorm.DB) *gorm.DB) Orm {
	o.GormDb = f(o.GormDb)
	return o
}

func (o Orm) SetDefaultOrder(s string) Orm {
	o.DefaultOrder = s
	return o
}

// FindPage 查询记录列表
func (o Orm) FindPage(pageSize int, page int, list interface{}, maps map[string]interface{}) PageInfo {
	var total int64 = 0
	if page == 0 {
		o.Where(maps).Order(o.DefaultOrder).GormDb.Find(list)
	} else {
		offset := pageSize * (page - 1)
		o.Where(maps).Order(o.DefaultOrder).GormDb.Offset(offset).Limit(pageSize).Find(list)
		total = o.Count(list, maps)
	}

	return PageInfo{
		PageSize:    pageSize,
		CurrentPage: page,
		Total:       total,
	}
}

// Find 查询记录列表
func (o Orm) Find(list interface{}, maps map[string]interface{}) {
	o.Where(maps).GormDb.Find(list)
}
func (o Orm) Count(list interface{}, maps map[string]interface{}) (count int64) {
	o.Where(maps).GormDb.Find(list).Count(&count)
	return
}
