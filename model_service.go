package gormx

import (
	"errors"
	"github.com/qwxingzhe/cast2"
	"github.com/spf13/cast"
	"log"
)

func CreateOne[TBase, TFormat any](orm Orm, data TBase) (info TFormat) {
	orm.Save(&data)
	info = cast2.CopyStruct(data, info)
	info = Format(info)
	return
}
func CreateOneSimp[TBase any](orm Orm, data TBase) TBase {
	orm.Save(&data)
	return data
}
func Format[T any](info T) T {
	if obj, ok := interface{}(info).(interface{ Format() T }); ok {
		info = obj.Format()
	}
	return info
}

func GetOne[TBase, TFormat any](orm Orm, where map[string]interface{}) (info TFormat) {
	var data TBase
	orm.First(&data, where)
	info = cast2.CopyStruct(data, info)
	info = Format(info)
	return
}
func GetOneSimp[TBase any](orm Orm, where map[string]interface{}) (info TBase) {
	orm.First(&info, where)
	return
}

func UpdateOneSimp[TBase any](orm Orm, data TBase) TBase {
	orm.Save(&data)
	return data
}

func UpdateOne[TFormat, TBase any](orm Orm, data TBase) (info TFormat) {
	orm.Save(&data)
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

	info2 = UpdateOne[Tf](orm, info)
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

// 查询列表
//+--------------------------------------------------------------------------------

// GetList 简单获取查询DB格式化后的列表
func GetList[TBase, TFormat any](orm Orm, c ListConfig) []TFormat {
	baseList := GetListSimp[TBase](orm, c)

	//formatEvery := false
	//if c != nil {
	//	formatEvery = c.FormatEvery
	//}
	return FormatListSimp[TBase, TFormat](baseList, c.FormatEvery)
}

func FormatListSimp[TBase, TFormat any](baseList []TBase, formatEvery bool) []TFormat {
	var list []TFormat
	for _, item := range baseList {
		var info TFormat
		info = cast2.CopyStruct(item, info)
		list = append(list, info)
	}
	if formatEvery {
		for i, tf := range list {
			list[i] = Format(tf)
		}
	}
	if len(list) == 0 {
		list = []TFormat{}
	}
	return list
}

func GetListSimp[TBase any](orm Orm, c ListConfig) []TBase {
	var baseList []TBase

	if c.FilterConfig.Order != "" {
		orm = orm.Order(c.FilterConfig.Order)
	}
	if c.FilterConfig.MaxLimit > 0 {
		orm = orm.Limit(c.FilterConfig.MaxLimit)
	}
	whereTrue := map[string]interface{}{}
	if c.Filter != nil {
		whereTrue = c.Filter
	}

	orm.Find(&baseList, whereTrue)

	return baseList
}

func Count[TBase any](orm Orm, maps map[string]interface{}) (count int64) {
	var list []TBase
	orm.Where(maps).GormDb.Find(&list).Count(&count)
	return
}

// FindPage 查询记录列表
func FindPage[TBase, TFormat any](o Orm, pageSize int, pageStr interface{}, c *ListConfig) PageResult {
	var baseList []TBase

	var total int64 = 0
	page := cast.ToInt(pageStr)
	offset := pageSize * (page - 1)

	log.Println("o.useTable : ", o.useTable, o.DefaultOrder, offset, pageSize)

	maps := map[string]interface{}{}
	if c.Filter != nil {
		maps = c.Filter
	}

	if o.useTable {
		o.Where(maps).Order(o.DefaultOrder).GormDb.Offset(offset).Limit(pageSize).Scan(&baseList)
	} else {
		if page == 0 {
			o.Where(maps).Order(o.DefaultOrder).GormDb.Find(&baseList)
		} else {
			total = o.Count(&baseList, maps)
			o.Where(maps).Order(o.DefaultOrder).GormDb.Offset(offset).Limit(pageSize).Find(&baseList)
		}
	}

	log.Println("FindPage 5")

	formatEvery := false
	if c != nil {
		formatEvery = c.FormatEvery
	}
	list := FormatListSimp[TBase, TFormat](baseList, formatEvery)

	log.Println("FindPage 6")

	return PageResult{
		List: list,
		PageInfo: PageInfo{
			PageSize:    pageSize,
			CurrentPage: page,
			Total:       total,
		}}

}
