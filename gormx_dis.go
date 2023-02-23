package gormx

//
//// Info 获取详情
//func Info[T any](o Orm, id int64) error {
//	var info T
//	key := ""
//	// 查询缓存是否存在
//	if o.cache != nil {
//		key = o.getCacheKey(info, id)
//		cache := o.getCache(key)
//
//		if cache != "" { // 存在，返回缓存
//			json.Unmarshal([]byte(cache), info)
//			return nil
//		}
//	}
//
//	tx := o.GormDb.Where("id = ?", id).First(info)
//	if tx.Error != nil { // 未查询到数据
//		fmt.Println("First Err:", tx.Error)
//		return tx.Error
//	}
//
//	// 写入缓存
//	if o.cache != nil {
//		err := o.setCache(key, info)
//		if err != nil {
//			fmt.Println("setCache Err:", err.Error())
//			return err
//		}
//	}
//
//	return nil
//}
//
//// Delete 删除记录
//func Delete(o Orm, info interface{}) {
//	tx := o.GormDb.Delete(info)
//	fmt.Println("tx::::", tx.Error)
//
//	// 删除缓存
//	o.delCacheByInfo(info)
//}
//
//// Save 保存记录
//func Save(o Orm, info interface{}) error {
//	err := o.GormDb.Save(info).Error
//	if err == nil {
//		// 删除缓存
//		o.delCacheByInfo(info)
//	}
//	return err
//}
