package service

import (
	"TikTok/dao"
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

	err = dao.Save(videoName+".mp4", userId, title)
	if err != nil {
		log.Printf("方法dao.Save(videoName, imageName, userId) 失败%v", err)
		return err
	}
	log.Printf("方法dao.Save(videoName, imageName, userId) 成功")
	return nil
	/**
	//获取文件名称
	fmt.Println(file.Filename)
	//文件大小
	fmt.Println(file.Size)
	//获取文件的后缀名
	extstring := path.Ext(file.Filename)
	fmt.Println(extstring)
	//根据当前时间鹾生成一个新的文件名
	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt, 10)
	//新的文件名
	fileName := fileNameStr + extstring
	//保存上传文件
	filePath := filepath.Join("upload", "/", fileName)
	err := c.SaveUploadedFile(file, filePath)
	if err != nil {
		return err
	}

	err = dao.Save(fileName, userId, title)
	if err != nil {
		log.Printf("方法dao.Save(videoName, imageName, userId) 失败%v", err)

		return err
	}
	return nil
	**/
}

func (videoService VideoServiceImpl) List(userId int64) ([]dao.TableVideo, error) {
	feed, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
