package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware/redis"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

type CommentServiceImpl struct {
	UserService
	CommentService
}

// CountFromVideoId
// 1、使用video id 查询Comment数量
func (c CommentServiceImpl) CountFromVideoId(videoId int64) (int64, error) {
	//先在缓存中查
	cnt, err := redis.RdbVCid.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
	if err != nil { //若查询缓存出错，则打印log
		//return 0, err
		log.Println("count from redis error:", err)
	}
	//1.缓存中查到了数量，则返回数量值-1（去除0值）
	if cnt != 0 {
		return cnt - 1, nil
	}
	//2.缓存中查不到则去数据库查
	cntDao, err1 := dao.Count(videoId)
	if err1 != nil {
		log.Println("comment count dao err:", err1)
		return 0, nil
	}
	//将评论id切片存入redis-第一次存储 V-C set 值：
	go func() {
		//查询评论id list
		cList, _ := dao.CommentIdList(videoId)
		//先在redis中存储一个-1值，防止脏读
		_, _err := redis.RdbVCid.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), config.DefaultRedisValue).Result()
		if _err != nil { //若存储redis失败，则直接返回
			log.Println("redis save one vId - cId 0 failed")
			return
		}
		//设置key值过期时间
		_, err := redis.RdbVCid.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			log.Println("redis save one vId - cId expire failed")
		}
		//评论id循环存入redis
		for _, commentId := range cList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), commentId)
		}
		log.Println("count comment save ids in redis")
	}()
	//返回结果
	return cntDao, nil
}

// Send
// 2、发表评论
func (c CommentServiceImpl) Send(comment dao.Comment) (CommentInfo, error) {
	log.Println("CommentService-Send: running") //函数已运行
	//数据准备
	var commentInfo dao.Comment
	commentInfo.VideoId = comment.VideoId         //评论视频id传入
	commentInfo.UserId = comment.UserId           //评论用户id传入
	commentInfo.CommentText = comment.CommentText //评论内容传入
	commentInfo.Cancel = config.ValidComment      //评论状态，0，有效
	commentInfo.CreateDate = comment.CreateDate   //评论时间

	//1.评论信息存储：
	commentRtn, err := dao.InsertComment(commentInfo)
	if err != nil {
		return CommentInfo{}, err
	}
	//2.查询用户信息
	impl := UserServiceImpl{}
	userData, err2 := impl.GetUserById(comment.UserId)
	if err2 != nil {
		return CommentInfo{}, err2
	}
	//3.拼接
	commentData := CommentInfo{
		Id:         commentRtn.Id,
		UserInfo:   userData,
		Content:    commentRtn.CommentText,
		CreateDate: commentRtn.CreateDate.Format(config.DateTime),
	}
	//将此发表的评论id存入redis
	go func() {
		insertRedisVideoCommentId(strconv.Itoa(int(comment.VideoId)), strconv.Itoa(int(commentRtn.Id)))
		log.Println("send comment save in redis")
	}()
	//返回结果
	return commentData, nil
}

// DelComment
// 3、删除评论，传入评论id
func (c CommentServiceImpl) DelComment(commentId int64) error {
	log.Println("CommentService-DelComment: running") //函数已运行
	//1.先查询redis，若有则删除，返回客户端-再go协程删除数据库；无则在数据库中删除，返回客户端。
	n, err := redis.RdbCVid.Exists(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		log.Println(err)
	}
	if n > 0 { //在缓存中有此值，则找出来删除，然后返回
		vid, err1 := redis.RdbCVid.Get(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err1 != nil { //没找到，返回err
			log.Println("redis find CV err:", err1)
		}
		//删除，两个redis都要删除
		del1, err2 := redis.RdbCVid.Del(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err2 != nil {
			log.Println(err2)
		}
		del2, err3 := redis.RdbVCid.SRem(redis.Ctx, vid, strconv.FormatInt(commentId, 10)).Result()
		if err3 != nil {
			log.Println(err3)
		}
		log.Println("del comment in Redis success:", del1, del2) //del1、del2代表删除了几条数据

	}
	//不在内存中，则直接走数据库删除
	return dao.DeleteComment(commentId)
}

// GetList
// 4、查看评论列表-返回评论list
func (c CommentServiceImpl) GetList(videoId int64, userId int64) ([]CommentInfo, error) {
	log.Println("CommentService-GetList: running") //函数已运行
	//1.先查询评论列表信息
	commentList, err := dao.GetCommentList(videoId)
	if err != nil {
		log.Println("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
		return nil, err
	}
	//当前有0条评论
	if commentList == nil {
		return nil, nil
	}

	//提前定义好切片长度
	commentInfoList := make([]CommentInfo, len(commentList))

	wg := &sync.WaitGroup{}
	wg.Add(len(commentList))
	idx := 0
	for _, comment := range commentList {
		//2.调用方法组装评论信息，再append
		var commentData CommentInfo
		//将评论信息进行组装，添加想要的信息,插入从数据库中查到的数据
		go func(comment dao.Comment) {
			oneComment(&commentData, &comment, userId)
			//3.组装list
			//commentInfoList = append(commentInfoList, commentData)
			commentInfoList[idx] = commentData
			idx = idx + 1
			wg.Done()
		}(comment)
	}
	wg.Wait()
	//评论排序-按照主键排序
	sort.Sort(CommentSlice(commentInfoList))
	//------------------------法二结束----------------------------

	//协程查询redis中是否有此记录，无则将评论id切片存入redis
	go func() {
		//1.先在缓存中查此视频是否已有评论列表
		cnt, err1 := redis.RdbVCid.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
		if err1 != nil { //若查询缓存出错，则打印log
			//return 0, err
			log.Println("count from redis error:", err)
		}
		//2.缓存中查到了数量大于0，则说明数据正常，不用更新缓存
		if cnt > 0 {
			return
		}
		//3.缓存中数据不正确，更新缓存：
		//先在redis中存储一个-1 值，防止脏读
		_, _err := redis.RdbVCid.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), config.DefaultRedisValue).Result()
		if _err != nil { //若存储redis失败，则直接返回
			log.Println("redis save one vId - cId 0 failed")
			return
		}
		//设置key值过期时间
		_, err2 := redis.RdbVCid.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			time.Duration(config.OneMonth)*time.Second).Result()
		if err2 != nil {
			log.Println("redis save one vId - cId expire failed")
		}
		//将评论id循环存入redis
		for _, _comment := range commentInfoList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), strconv.Itoa(int(_comment.Id)))
		}
		log.Println("comment list save ids in redis")
	}()

	log.Println("CommentService-GetList: return list success") //函数执行成功，返回正确信息
	return commentInfoList, nil
}

// 在redis中存储video_id对应的comment_id 、 comment_id对应的video_id
func insertRedisVideoCommentId(videoId string, commentId string) {
	//在redis-RdbVCid中存储video_id对应的comment_id
	_, err := redis.RdbVCid.SAdd(redis.Ctx, videoId, commentId).Result()
	if err != nil { //若存储redis失败-有err，则直接删除key
		log.Println("redis save send: vId - cId failed, key deleted")
		redis.RdbVCid.Del(redis.Ctx, videoId)
		return
	}
	//在redis-RdbCVid中存储comment_id对应的video_id
	_, err = redis.RdbCVid.Set(redis.Ctx, commentId, videoId, 0).Result()
	if err != nil {
		log.Println("redis save one cId - vId failed")
	}
}

// 此函数用于给一个评论赋值：评论信息+用户信息 填充
func oneComment(comment *CommentInfo, com *dao.Comment, userId int64) {
	var wg sync.WaitGroup
	wg.Add(1)
	//根据评论用户id和当前用户id，查询评论用户信息
	impl := UserServiceImpl{}
	var err error
	comment.Id = com.Id
	comment.Content = com.CommentText
	comment.CreateDate = com.CreateDate.Format(config.DateTime)
	comment.UserInfo, err = impl.GetUserById(com.UserId)
	if err != nil {
		log.Println("CommentService-GetList: GetUserById return err: " + err.Error()) //函数返回提示错误信息
	}
	wg.Done()
	wg.Wait()
}

// CommentSlice 此变量以及以下三个函数都是做排序-准备工作
type CommentSlice []CommentInfo

func (a CommentSlice) Len() int { //重写Len()方法
	return len(a)
}
func (a CommentSlice) Swap(i, j int) { //重写Swap()方法
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool { //重写Less()方法
	return a[i].Id > a[j].Id
}
