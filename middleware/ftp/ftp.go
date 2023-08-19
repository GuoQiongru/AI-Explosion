package ftp

import (
	"TikTok/config"
	"log"
	"time"

	"github.com/dutchcoders/goftp"
)

var MyFTP *goftp.FTP

func InitFTP() {

	var err error
	MyFTP, err = goftp.Connect(config.ConConfig)
	if err != nil {
		log.Printf("获取到FTP链接失败！！！")
	}
	log.Printf("获取到FTP链接成功%v：", MyFTP)

	err = MyFTP.Login(config.FtpUser, config.FtpPsw)
	if err != nil {
		log.Printf("Login FTP Fail")
	}
	log.Printf("Login FTP Successed")

	go keepAlive()
}

func keepAlive() {
	time.Sleep(time.Duration(config.HeartbeatTime) * time.Second)
	MyFTP.Noop()
}
