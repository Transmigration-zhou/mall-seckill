package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"mall-seckill/datamodels"
	"mall-seckill/services"
	"sync"
)

const MQURL = "amqp://guest:guest@127.0.0.1:5672/"

// RabbitMQ rabbitMQ结构体
type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	MqUrl     string
	sync.Mutex
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, MqUrl: MQURL}
}

func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Printf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	//申请队列，如果不存在会自动创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //是否自动删除（当最后一个消费者从连接里断开后）
		false, //是否具有排他性
		false, //是否阻塞处理
		nil,
	)
	if err != nil {
		return err
	}

	r.channel.PublishWithContext(
		context.Background(),
		r.Exchange, // 此处为空
		r.QueueName,
		false, //如果为true，根据自身exchange类型和routeKey规则；无法找到符合条件的队列会把消息返还给发送者
		false, //如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return nil
}

func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //是否持久化
		false, //是否自动删除
		false, //是否具有排他性
		false, //是否阻塞处理
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	msg, err := r.channel.Consume(
		q.Name,
		"",    //用来区分多个消费者
		false, //是否自动应答
		false, //是否具有排他性
		false, //设置为true，表示不能将同一个Connection中生产者发送的消息传递给这个Connection中的消费者
		false, // 是否阻塞处理
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	//消费者流控制 防止爆库
	r.channel.Qos(
		1,     //当前消费者一次能接受的最大消息数量
		0,     //服务器传递的最大容量（以八位字节为单位）
		false, //如果设置为true 对channel可用
	)

	// 此处使用forever的意思为因为协程会始终监听消息(除非手动结束)
	// 手动结束才会进行 <-forever 有协程且一直尝试读取数据
	forever := make(chan bool)
	go func() {
		for d := range msg {
			log.Printf("Receibed a message")
			fmt.Println(string(d.Body))
			//消息逻辑处理
			message := &datamodels.Message{}
			err := json.Unmarshal(d.Body, message)
			if err != nil {
				fmt.Println(err)
			}
			//插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}
			//扣除商品数量
			err = productService.SubNumberOne(message.ProductId)
			if err != nil {
				fmt.Println(err)
			}
			d.Ack(false) //true表示确认所有未确认的消息 false表示确认当前消息
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
