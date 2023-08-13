package service

import (
	"TikTok/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"path"
	"path/filepath"
	"strconv"
	"time"
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

func (videoService *VideoServiceImpl) Publish(c *gin.Context, file *multipart.FileHeader, userId int64, title string) error {
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
}

func (videoService VideoServiceImpl) List(userId int64) ([]dao.TableVideo, error) {
	feed, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		return nil, err
	}

	return feed, nil
}
