package service

import (
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"log"
	"mime/multipart"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type VideoServiceImpl struct {
	UserService
}

func (videoService VideoServiceImpl) Feed(lastTime time.Time) ([]dao.TableVideo, error) {
	feed, err := dao.GetVideosByLastTime(lastTime)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (videoService *VideoServiceImpl) Publish(c *gin.Context, data *multipart.FileHeader, userId int64, title string) error {
	file, err := data.Open()
	if err != nil {
		log.Printf("方法data.Open() 失败%v", err)
		return err
	}
	log.Printf("方法data.Open() 成功")
	//ext := filepath.Ext(data.Filename)
	//生成一个uuid作为视频的名字
	videoName := uuid.NewV4().String()
	log.Printf("生成视频名称%v", videoName)
	err = dao.VideoFTP(file, videoName)
	if err != nil {
		log.Printf("方法dao.VideoFTP(file, videoName) 失败%v", err)
		return err
	}
	log.Printf("方法dao.VideoFTP(file, videoName) 成功")
	defer file.Close()

	//在服务器上执行ffmpeg 从视频流中获取第一帧截图，并上传图片服务器，保存图片链接
	imageName := uuid.NewV4().String()
	//向队列中添加消息
	ffmpeg.Ffchan <- ffmpeg.Ffmsg{
		videoName,
		imageName,
	}
	//组装并持久化

	err = dao.Save(videoName+".mp4", imageName+".jpg", userId, title)
	if err != nil {
		log.Printf("方法dao.Save(videoName, imageName, userId) 失败%v", err)
		return err
	}
	log.Printf("方法dao.Save(videoName, imageName, userId) 成功")
	return nil
}

func (videoService VideoServiceImpl) List(userId int64) ([]dao.TableVideo, error) {
	feed, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
