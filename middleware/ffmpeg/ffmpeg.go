package ffmpeg

import (
	"TikTok/config"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

type Ffmsg struct {
	VideoName string
	ImageName string
}

var ClientSSH *ssh.Client
var Ffchan chan Ffmsg

func InitSSH() {
	var err error
	SSHconfig := &ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            config.UserSSH,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if config.TypeSSH == "password" {
		SSHconfig.Auth = []ssh.AuthMethod{ssh.Password(config.PasswordSSH)}
	}
	addr := fmt.Sprintf("%s:%d", config.HostSSH, config.PortSSH)
	ClientSSH, err = ssh.Dial("tcp", addr, SSHconfig)
	if err != nil {
		log.Fatal("Create ssh client Failed", err)
	}
	log.Printf("Create ssh client Successfully: %v", ClientSSH)
	Ffchan = make(chan Ffmsg, config.MaxMsgCount)
	go dispatcher()
	go keepAlive()
}

func dispatcher() {
	for ffmsg := range Ffchan {
		go func(f Ffmsg) {
			err := Ffmpeg(f.VideoName, f.ImageName)
			if err != nil {
				Ffchan <- f
				log.Fatal("redispatcher")
			}
			log.Printf("The screenshot of : %v is processed successfully", f.VideoName)
		}(ffmsg)
	}
}

func Ffmpeg(videoName string, imageName string) error {
	session, err := ClientSSH.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	combo, err := session.CombinedOutput("ffmpeg -ss 00:00:01 -i /home/ftp/ftpuser/video/" + videoName + ".mp4 -frames:v 1 /home/ftp/ftpuser/images/" + imageName + ".jpg")
	if err != nil {
		log.Fatal("shell failed:", string(combo))
		return err
	}
	return nil
}

func keepAlive() {
	for {
		time.Sleep(time.Duration(config.SSHHeartbeatTime) * time.Second)
		session, _ := ClientSSH.NewSession()
		session.Close()
	}
}
