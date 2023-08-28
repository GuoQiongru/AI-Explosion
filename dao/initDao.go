package dao

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func Init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	var err error
	dsn := "root:123456@(47.113.148.197:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	Db.AutoMigrate(&TableUser{})
	Db.AutoMigrate(&TableVideo{})
	Db.AutoMigrate(&Like{})
	Db.AutoMigrate(&Comment{})

	if err != nil {
		log.Panicln("err:", err.Error())
	}
}
