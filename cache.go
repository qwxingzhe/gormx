package gormx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// 缓存操作区
//+---------------------------------------------------------------------
// 获取缓存key
func (o Orm) getCacheKey(info interface{}, id int64) string {
	refH := reflect.TypeOf(info)
	tableName := strings.ToLower(refH.String())
	tableName = strings.ReplaceAll(tableName, "*", "")

	return fmt.Sprintf("cache:%s:%d", tableName, id)
}

// 删除缓存
func (o Orm) delCacheByInfo(info interface{}) error {
	if o.cache == nil {
		return nil
	}
	refH := reflect.ValueOf(info)
	var id int64
	if refH.Kind() == reflect.Struct { // 结构体
		id = refH.FieldByName("Id").Int()
	} else if refH.Kind() == reflect.Ptr { // 指针类型
		id = reflect.Indirect(refH).FieldByName("Id").Int()
	}
	key := o.getCacheKey(info, id)
	return o.delCache(key)
}

// 新建缓存
func (o Orm) setCache(key string, value interface{}) error {
	if o.cache == nil {
		return nil
	}
	str1, err := json.Marshal(&value)
	str := string(str1)
	if err != nil {
		return err
	}
	return o.cache.SetCache(key, str)
}

// 查询缓存
func (o Orm) getCache(key string) string {
	if o.cache == nil {
		return ""
	}
	return o.cache.GetCache(key)
}

// 删除缓存
func (o Orm) delCache(key string) error {
	if o.cache == nil {
		return nil
	}
	return o.cache.DelCache(key)
}
