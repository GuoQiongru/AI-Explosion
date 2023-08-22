package rabbitmq

import (
	"TikTok/config"
	"TikTok/dao"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

type LikeMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

func NewLikeRabbitMQ(queueName string) *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}
	cha, err := likeMQ.conn.Channel()
	likeMQ.channel = cha
	Rmq.failOnErr(err, "获取通道失败")
	return likeMQ
}

func (l *LikeMQ) Publish(message string) {

	_, err := l.channel.QueueDeclare(
		l.queueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}

	err1 := l.channel.Publish(
		l.exchange,
		l.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err1 != nil {
		panic(err)
	}

}

func (l *LikeMQ) Consumer() {

	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	messages, err1 := l.channel.Consume(
		l.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}

	forever := make(chan bool)
	switch l.queueName {
	case "like_add":
		go l.consumerLikeAdd(messages)
	case "like_del":
		go l.consumerLikeDel(messages)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

func (l *LikeMQ) consumerLikeAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		for i := 0; i < config.Attempts; i++ {
			flag := false
			var likeData dao.Like
			likeInfo, err := dao.GetLikeInfo(userId, videoId)
			if err != nil {
				log.Printf(err.Error())
				flag = true
			} else {
				if likeInfo == (dao.Like{}) {
					likeData.UserId = userId
					likeData.VideoId = videoId
					likeData.Cancel = config.IsLike
					if err := dao.InsertLike(likeData); err != nil {
						log.Printf(err.Error())
						flag = true
					}
				} else {
					if err := dao.UpdateLike(userId, videoId, config.IsLike); err != nil {
						log.Printf(err.Error())
						flag = true
					}
				}
				if flag == false {
					break
				}
			}
		}
	}
}

func (l *LikeMQ) consumerLikeDel(messages <-chan amqp.Delivery) {
	for d := range messages {
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.ParseInt(params[0], 10, 64)
		videoId, _ := strconv.ParseInt(params[1], 10, 64)
		for i := 0; i < config.Attempts; i++ {
			flag := false
			likeInfo, err := dao.GetLikeInfo(userId, videoId)
			if err != nil {
				log.Printf(err.Error())
				flag = true
			} else {
				if likeInfo == (dao.Like{}) {
					log.Printf(errors.New("can't find data,this action invalid").Error())
				} else {
					if err := dao.UpdateLike(userId, videoId, config.Unlike); err != nil {
						log.Printf(err.Error())
						flag = true
					}
				}
			}
			if flag == false {
				break
			}
		}
	}
}

var RmqLikeAdd *LikeMQ
var RmqLikeDel *LikeMQ

func InitLikeRabbitMQ() {
	RmqLikeAdd = NewLikeRabbitMQ("like_add")
	go RmqLikeAdd.Consumer()

	RmqLikeDel = NewLikeRabbitMQ("like_del")
	go RmqLikeDel.Consumer()
}
