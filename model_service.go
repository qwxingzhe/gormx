package gormx

import (
	"errors"
	"github.com/qwxingzhe/cast2"
)

type FilterConfig struct {
	Order    string
	MaxLimit int
}
type ListConfig struct {
	GetFilter    map[string]interface{}
	FilterConfig FilterConfig
	FormatEvery  bool
}
type UpdateConfig struct {
	GetFilter map[string]interface{}
	EmptyMsg  string
}

// GetList 简单获取查询DB格式化后的列表
func GetList[TBase, TFormat any](orm Orm, c *ListConfig) []TFormat {
	baseList := GetListSimp[TBase](orm, c)

	var list []TFormat
	for _, item := range baseList {
		var info TFormat
		info = cast2.CopyStruct(item, info)
		list = append(list, info)
	}
	if c != nil && c.FormatEvery {
		for i, tf := range list {
			list[i] = Format(tf)
		}
	}
	return list
}

func GetListSimp[TBase any](orm Orm, c *ListConfig) []TBase {
	var baseList []TBase
	whereTrue := map[string]interface{}{}
	if c != nil {
		if c.FilterConfig.Order != "" {
			orm = orm.Order(c.FilterConfig.Order)
		}
		if c.FilterConfig.MaxLimit > 0 {
			orm = orm.Limit(c.FilterConfig.MaxLimit)
		}
		if c.GetFilter != nil {
			whereTrue = c.GetFilter
		}
	}

	orm.Find(&baseList, whereTrue)

	return baseList
}

func CreateOne[TBase, TFormat any](orm Orm, data TBase) (info TFormat) {
	orm.Save(&data)
	info = cast2.CopyStruct(data, info)
	info = Format(info)
	return
}

func Format[T any](info T) T {
	if obj, ok := interface{}(info).(interface{ Format() T }); ok {
		info = obj.Format()
	}
	return info
}

func UpdateOne[TBase, TFormat any](orm Orm, data TBase) (info TFormat) {
	orm.Save(&data)
	info = cast2.CopyStruct(data, info)
	info = Format(info)
	return
}

func GetOne[TBase, TFormat any](orm Orm, where map[string]interface{}) (info TFormat) {
	var data TBase
	orm.First(&data, where)
	info = cast2.CopyStruct(data, info)
	info = Format(info)
	return
}

func Update[Tm, Tf any](orm Orm, uc UpdateConfig, SetFunc func(info Tm) (Tm, error)) (info2 Tf, err error) {
	info := GetOne[Tm, Tm](orm, uc.GetFilter)
	id := cast2.StructValue(info, "ID")
	if id.(int64) == 0 {
		msg := uc.EmptyMsg
		if msg == "" {
			msg = "数据不存在或无权限修改"
		}
		err = errors.New(msg)
		return
	}
	info, err2 := SetFunc(info)
	if err2 != nil {
		err = errors.New(err2.Error())
		return
	}

	info2 = UpdateOne[Tm, Tf](orm, info)
	return
}

func Delete[Tm any](orm Orm, uc UpdateConfig) (err error) {
	info := GetOne[Tm, Tm](orm, uc.GetFilter)
	id := cast2.StructValue(info, "ID")
	if id.(int64) == 0 {
		msg := uc.EmptyMsg
		if msg == "" {
			msg = "数据不存在或无权限删除"
		}
		err = errors.New(msg)
		return
	}
	orm.Delete(&info)
	return
}
