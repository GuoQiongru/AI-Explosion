package dao

import (
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "gorm.io/gorm/logger"
  "log"
  "os"
  "time"
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
  dsn := "root:root@(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
  Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: newLogger,
  })
  if err != nil {
    log.Panicln("err:", err.Error())
  }
}