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
		log.Printf("Connect FTP Fail")
	}
	log.Printf("Connect FTP Successed: %v：", MyFTP)

	err = MyFTP.Login(config.FtpUser, config.FtpPsw)
	if err != nil {
		log.Printf("Login FTP Fail")
	}
	log.Printf("Login FTP Successed")

	go keepAlive()
}

func keepAlive() {
	for {
		time.Sleep(time.Duration(config.HeartbeatTime) * time.Second)
		err := MyFTP.Noop()
		if err != nil {
			log.Printf("FTP NOOP Fail")
		}
		log.Printf("FTP NOOP Successed")
	}
}
