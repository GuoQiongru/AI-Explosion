package dao

import (
	"log"
)

type TableUser struct {
	Id       int64
	Name     string
	Password string
}

func (tableUser TableUser) TableName() string {
	return "users"
}

func GetTableUserList() ([]TableUser, error) {
	tableUsers := []TableUser{}
	if err := Db.Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

func GetTableUserByUsername(name string) (TableUser, error) {
	tableUser := TableUser{}
	if err := Db.Where("name = ?", name).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

func GetTableUserById(id int64) (TableUser, error) {
	tableUser := TableUser{}
	if err := Db.Where("id = ?", id).First(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return tableUser, err
	}
	return tableUser, nil
}

func InsertTableUser(tableUser *TableUser) bool {
	if err := Db.Create(&tableUser).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
