package service

import (
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"log"
	"mime/multipart"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

type VideoServiceImpl struct {
	UserService
	LikeService
	VideoService
}

func (videoService VideoServiceImpl) Feed(lastTime time.Time) ([]dao.TableVideo, error) {
	feed, err := dao.GetVideosByLastTime(lastTime)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64, title string) error {
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

func (videoService *VideoServiceImpl) GetVideo(videoId int64, userId int64) (Video, error) {
	//初始化video对象
	var video Video

	//从数据库中查询数据，如果查询不到数据，就直接失败返回，后续流程就不需要执行了
	data, err := dao.GetVideoByVideoId(videoId)
	if err != nil {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 失败：%v", err)
		return video, err
	} else {
		log.Printf("方法dao.GetVideoByVideoId(videoId) 成功")
	}

	//插入从数据库中查到的数据
	videoService.creatVideo(&video, &data, userId)
	return video, nil
}

func (videoService *VideoServiceImpl) creatVideo(video *Video, data *dao.TableVideo, userId int64) {
	//建立协程组，当这一组的携程全部完成后，才会结束本方法
	var wg sync.WaitGroup
	wg.Add(3)
	var err error
	video.TableVideo = *data
	//插入Author，这里需要将视频的发布者和当前登录的用户传入，才能正确获得isFollow，
	//如果出现错误，不能直接返回失败，将默认值返回，保证稳定
	go func() {
		video.Author, err = videoService.GetUserById(data.AuthorId)
		if err != nil {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 失败：%v", err)
		} else {
			log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 成功")
		}
		wg.Done()
	}()

	//插入点赞数量，同上所示，不将nil直接向上返回，数据没有就算了，给一个默认就行了
	go func() {
		video.FavoriteCount, err = videoService.FavouriteCount(data.Id)
		if err != nil {
			log.Printf("方法videoService.FavouriteCount(data.ID) 失败：%v", err)
		} else {
			log.Printf("方法videoService.FavouriteCount(data.ID) 成功")
		}
		wg.Done()
	}()

	//获取当前用户是否点赞了该视频
	go func() {
		video.IsFavorite, err = videoService.IsFavourite(video.Id, userId)
		if err != nil {
			log.Printf("方法videoService.IsFavourit(video.Id, userId) 失败：%v", err)
		} else {
			log.Printf("方法videoService.IsFavourit(video.Id, userId) 成功")
		}
		wg.Done()
	}()

	wg.Wait()
}
