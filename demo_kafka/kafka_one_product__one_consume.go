package demo_kafka

import (
	"context"
	"fmt"
	"github.com/510909033/redisx"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

var conn *kafka.Conn
var once sync.Once
var msgNum = int64(0)
var msgMutex sync.Mutex

func getDemoMessage() int64 {
	msgMutex.Lock()
	msgNum++
	msgMutex.Unlock()
	return msgNum

	client := redisx.GetTestRedisClient()
	before, err := client.IncrBy(redisx.GetTestCtx(), "kafka_key", 1)
	if err != nil {
		panic(err)
	}
	return before
}
func getDemoMessageStr() string {
	before := getDemoMessage()
	return fmt.Sprintf("%d", before)
}

func getConn() *kafka.Conn {
	once.Do(func() {
		var err error
		conn, err = kafka.Dial("tcp", getAddrs())
		if err != nil {
			panic(err.Error())
		}
	})
	//defer conn.Close()
	return conn
}

type DemoKafkaOneProductOneConsume struct {
}

func logDebug(format string, v ...interface{}) {
	log.Printf("[ DEBUG ] "+format+"\n", v...)
}

func logErr(format string, v ...interface{}) {
	log.Printf("[ ERROR ] "+format+"\n", v...)
}

func (demo *DemoKafkaOneProductOneConsume) getPartitions(topic string) []kafka.Partition {
	conn := getConn()
	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		panic(err.Error())
	}
	kafka.NewConsumerGroup()

	logDebug("topic=%s, Partition数=%d\n", topic, len(partitions))
	for k, v := range partitions {
		_ = k
		_ = v
		//logDebug("topic=%s, partition 详情：%+v", topic, v)
		logDebug("topic=%s, partition 详情：ID=%d", topic, v.ID)
	}

	return partitions
}

func (demo *DemoKafkaOneProductOneConsume) getControllerConn() *kafka.Conn {
	conn := getConn()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}

	logDebug("Controller=%+v", controller)

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	return controllerConn

	//defer controllerConn.Close()
}
func (demo *DemoKafkaOneProductOneConsume) deleteTopic(topic string) bool {
	var controllerConn = demo.getControllerConn()
	defer controllerConn.Close()

	err := controllerConn.DeleteTopics(topic)
	if err != nil {
		logErr("DeleteTopics fail, topic=%s, err=%+v", topic, err)
		return false
	}
	logDebug("删除topic成功, topic=%s", topic)
	return true
}
func (demo *DemoKafkaOneProductOneConsume) createTopic(topic string, numPartitions int) bool {
	//if demo.topicExists(topic) {
	//	return false
	//}

	conn := getConn()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}

	logDebug("Controller=%+v", controller)

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	//err = controllerConn.DeleteTopics(topicName)
	//if err != nil {
	//	log.Fatalf("DeleteTopics fail, err=%+v", err)
	//}

	config := kafka.TopicConfig{
		Topic:              topic,
		NumPartitions:      numPartitions,
		ReplicationFactor:  2,
		ReplicaAssignments: nil,
		ConfigEntries:      nil,
	}
	logDebug("CreateTopics, config=%+v", config)

	err = controllerConn.CreateTopics(config)
	if err != nil {
		logErr("CreateTopics fail, err=%+v", err)
		return false
	}
	logDebug("创建topic成功, topic=%s\n", topic)
	return true
}
func (demo *DemoKafkaOneProductOneConsume) topicExists(topic string) bool {
	partitions := demo.getPartitions(topic)

	log.Printf("Partition数=%d\n", len(partitions))

	return len(partitions) > 0
}

func (demo *DemoKafkaOneProductOneConsume) demoReader() {
	log.Println("start ", time.Now().String())
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: getAddrsList(),
		//Partition: 1,
		GroupID:     "demo_fetch_group",
		GroupTopics: nil,
		Topic:       demoGetTopic(),
		Partition:   0,
		Dialer: &kafka.Dialer{
			ClientID: fmt.Sprintf("%d", getClientID()),
			//DialFunc:        nil,
			Timeout:  time.Second * 5,
			Deadline: time.Now().Add(time.Second),
			//LocalAddr:       nil,
			//DualStack:       false,
			//FallbackDelay:   0,
			//KeepAlive:       0,
			//Resolver:        nil,
			//TLS:             nil,
			//SASLMechanism:   nil,
			//TransactionalID: "",
		},
		QueueCapacity: 0,
		//MinBytes: 10e3, // 10KB
		//MaxBytes: 10e6, // 10MB
		MinBytes:               1, // 10KB
		MaxBytes:               5, // 10MB
		MaxWait:                0,
		ReadLagInterval:        0,
		GroupBalancers:         nil,
		HeartbeatInterval:      0,
		CommitInterval:         0,
		PartitionWatchInterval: 0,
		WatchPartitionChanges:  false,
		SessionTimeout:         0,
		RebalanceTimeout:       0,
		JoinGroupBackoff:       0,
		RetentionTime:          0,
		StartOffset:            0,
		ReadBackoffMin:         0,
		ReadBackoffMax:         0,
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("Logger , "+msg, args...)
		}),
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("ErrorLogger,"+msg, args...)
		}),
		IsolationLevel: 0,
		MaxAttempts:    0,
	})
	//r.SetOffset(42)

	size := 20
	currNum := 0

	r.Lag()

	for {
		currNum++
		if currNum > size {
			break
		}
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			panic(err)
			break
		}

		fmt.Printf("message at ClientID:%s, Partition:%d, offset %d: %s = %s\n",
			r.Stats().ClientID,
			m.Partition, m.Offset, string(m.Key), string(m.Value))
		//fmt.Printf("message %+v\n", m.Time)

	}

	log.Println("sleep ... ")
	time.Sleep(time.Second * 2)

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("over")
}

func (demo *DemoKafkaOneProductOneConsume) demoFetch() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  getAddrsList(),
		Topic:    demoGetTopic(),
		MinBytes: 10e3, // 10KB
		//MinBytes: 1,    // 10KB
		MaxBytes: 10e6, // 10MB
		GroupID:  "demo_fetch_group",
	})

	size := 2
	currNum := 0

	ctx := context.Background()
	for {
		currNum++
		if currNum > size {
			break
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}

	log.Println("sleep ... ")
	time.Sleep(time.Minute)

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("over")
}

func (demo *DemoKafkaOneProductOneConsume) getTopicName() string {
	return "Topic__v6"
}

func (demo *DemoKafkaOneProductOneConsume) getGroupId() string {
	return "demo_group_v9"
}

func (demo *DemoKafkaOneProductOneConsume) demoWrite() {
	w := kafka.Writer{
		Addr:     kafka.TCP(getAddrsList()...),
		Topic:    demo.getTopicName(),
		Balancer: nil,
		//MaxAttempts:  0,
		//BatchSize:    0,
		//BatchBytes:   0,
		//BatchTimeout: 0,
		//ReadTimeout:  0,
		//WriteTimeout: 0,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		//Completion:   nil,
		//Compression:  0,
		//Logger:       nil,
		//ErrorLogger:  nil,
		//Transport:    nil,
	}
	defer func() {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}()

	for i := 0; i < 100; i++ {
		ctx := context.Background()
		err := w.WriteMessages(ctx, kafka.Message{
			//Topic:         demoGetTopic(),
			//Partition:     0,
			//Offset:        0,
			//HighWaterMark: 0,
			//Key:           nil,
			Value:   []byte(time.Now().String()),
			Headers: nil,
			//Time:          time.Time{},
		})
		if err != nil {
			panic(err)
		}
	}

	stats := w.Stats()
	log.Printf("stats=%+v\n", stats)
}

func (demo *DemoKafkaOneProductOneConsume) consumeMessageByPartitions(topic string, groupId string, partition int, num int) {
	//if groupId == "" {
	//	panic("groupId不能为空")
	//}
	//kafka.NewConsumerGroup()

	logDebug("consumeMessageByPartitions start %s", time.Now().String())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     getAddrsList(),
		GroupID:     groupId,
		GroupTopics: nil,
		Topic:       topic,
		Partition:   partition,
		Dialer: &kafka.Dialer{
			ClientID: fmt.Sprintf("%d", getClientID()),
			//DialFunc:        nil,
			Timeout: time.Second * 5,
			//Deadline: time.Now().Add(time.Second),
			//LocalAddr:       nil,
			//DualStack:       false,
			//FallbackDelay:   0,
			//KeepAlive:       0,
			//Resolver:        nil,
			//TLS:             nil,
			//SASLMechanism:   nil,
			//TransactionalID: "",
		},
		QueueCapacity: 0,
		//MinBytes: 10e3, // 10KB
		//MaxBytes: 10e6, // 10MB
		MinBytes:               1, // 10KB
		MaxBytes:               5, // 10MB
		MaxWait:                0,
		ReadLagInterval:        0,
		GroupBalancers:         []kafka.GroupBalancer{&MyGroupBalancer{}},
		HeartbeatInterval:      time.Second,
		CommitInterval:         0,
		PartitionWatchInterval: 0,
		WatchPartitionChanges:  false,
		SessionTimeout:         0,
		RebalanceTimeout:       0,
		JoinGroupBackoff:       0,
		RetentionTime:          0,
		StartOffset:            0,
		ReadBackoffMin:         0,
		ReadBackoffMax:         0,
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("Logger , "+msg, args...)
		}),
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("ErrorLogger,"+msg, args...)
		}),
		IsolationLevel: 0,
		MaxAttempts:    0,
	})
	//reader.SetOffset(42)

	currNum := 0

	//todo
	time.Sleep(time.Second * 5)

	for {
		currNum++
		if currNum > num {
			break
		}
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logErr("ReadMessage , topic=%s, groupId=%s, err=%+v", topic, groupId, err)
			time.Sleep(time.Second)
			continue
			//todo
			panic(err)
			break
		}

		fmt.Printf("ReadMessage topic=%s,GroupId=%s, ClientID=%s, Partition=%d, offset=%d, key=%s, 消息=%s\n",
			topic,
			groupId,
			reader.Stats().ClientID,
			msg.Partition,
			msg.Offset,
			string(msg.Key),
			string(msg.Value),
		)

	}

	if err := reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	logDebug("consumeMessageByPartitions over")
}

// 通过分区号发送消息
func (demo *DemoKafkaOneProductOneConsume) writeMessageByPartitions(topic string, cnt int, balancer kafka.Balancer) {
	w := kafka.Writer{
		Addr:     kafka.TCP(getAddrsList()...),
		Topic:    topic,
		Balancer: balancer,
		//MaxAttempts:  0,
		//BatchSize:    0,
		//BatchBytes:   0,
		//BatchTimeout: 0,
		//ReadTimeout:  0,
		//WriteTimeout: 0,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		Completion:   nil,
		//Compression:  0,
		//Logger:       nil,
		//ErrorLogger:  nil,
		//Transport:    nil,
	}
	defer func() {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}()

	for i := 0; i < cnt; i++ {
		ctx := context.Background()
		err := w.WriteMessages(ctx, kafka.Message{
			//Topic:         demoGetTopic(),
			//Partition: partition, //只读属性，设置了也没用
			//Offset:        0,
			//HighWaterMark: 0,
			//Key:           nil,
			Value:   []byte(getDemoMessageStr()),
			Headers: nil,
			//Time:          time.Time{},
			Key: []byte("a"),
		})
		if err != nil {
			panic(err)
		}
	}

	logDebug("writeMessageByPartitions, topic=%s, 发送次数=%d", topic, cnt)
}
