package gormx

import (
	"fmt"
	"testing"
)

func initOrm() Orm {
	//dbconfig := MysqlConfig{
	//	DBAddress:      "127.0.0.1:3306",
	//	DBUserName:     "root",
	//	DBPassword:     "123456",
	//	DBDatabaseName: "test",
	//}
	dbconfig := getTestMysqlConfig()

	return InitOrm(dbconfig, nil)
}

type User struct {
	// model.BaseModel
	BaseModel
	Name   string `form:"name"`
	Gender int    `form:"gender"`
}

func (u User) Format() interface{} {
	return u.Name
}

func TestSave(t *testing.T) {
	orm := initOrm()
	user := User{
		Name:   "张三",
		Gender: 1,
	}
	orm.Save(&user)
}
func TestGetInfo(t *testing.T) {
	orm := initOrm()
	user := User{}
	orm.GetInfo(&user, 1)
	fmt.Println(user.Name)
}
func TestDelete(t *testing.T) {
	orm := initOrm()
	user := User{
		Name:   "小七",
		Gender: 1,
	}
	orm.Save(&user)
	orm.Delete(&user)
}
func TestFind(t *testing.T) {
	orm := initOrm()
	var list []User
	orm.Find(&list, map[string]interface{}{
		"gender": 1,
	})
	fmt.Println(list)
}

func TestFormatList(t *testing.T) {
	orm := initOrm()
	var list []User
	orm.Find(&list, map[string]interface{}{
		"gender": 1,
	})
	list2 := FormatList(list, nil)
	fmt.Println(list2)
}
func TestKeyList(t *testing.T) {
	orm := initOrm()
	var list []User
	orm.Find(&list, map[string]interface{}{
		"gender": 1,
	})
	list2 := KeyList(list, "Id")
	fmt.Println(list2)
}
func TestFindPage(t *testing.T) {
	orm := initOrm()
	var list []User
	orm.FindPage(2, 1, &list, map[string]interface{}{
		"gender": 1,
	})
	fmt.Println(list)
}
