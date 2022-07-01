package demo_kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type DemoKafka struct {
}

//列出topic
func (demo *DemoKafka) demoGetPartitions() {
	conn, err := kafka.Dial("tcp", getAddrs())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	//log.Printf("%#v\n", partitions)
	log.Printf("Partition数=%d\n", len(partitions))

	m := map[string]kafka.Partition{}

	for _, p := range partitions {
		if p.Topic == demoGetTopic() {
			log.Printf("hit demoGetTopic, topic=%s\n", p.Topic)
		}
		m[p.Topic] = p
	}
	//for k := range m {
	//	fmt.Println(k)
	//}

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}

	log.Printf("Controller=%+v\n", controller)

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	var topicName = "demo_wbt_topic"
	//err = controllerConn.DeleteTopics(topicName)
	//if err != nil {
	//	log.Fatalf("DeleteTopics fail, err=%+v", err)
	//}

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:              topicName,
		NumPartitions:      3,
		ReplicationFactor:  1,
		ReplicaAssignments: nil,
		ConfigEntries:      nil,
	})
	if err != nil {
		log.Fatalf("CreateTopics fail, err=%+v", err)
	}
	log.Printf("创建topic无报错, topic=%s\n", topicName)

	//ID 是不对的，都是0 。。。，但是host是对的
	log.Printf("conn.Broker().ID=%d, host=%s", conn.Broker().ID, conn.Broker().Host)
	log.Printf("controllerConn.Broker().ID=%d, host=%s", controllerConn.Broker().ID, controllerConn.Broker().Host)
	brokers, err := conn.Brokers()
	if err != nil {
		panic(err)
	}
	for _, v := range brokers {
		log.Printf("一个broker, id=%d", v.ID)
	}

	dialPartition, err := kafka.DialPartition(context.Background(), "tcp", getAddrs(), m[demoGetTopic()])
	if err != nil {
		panic(err)
	}
	message, err := dialPartition.ReadMessage(1024)
	log.Printf("msg=%s, message=%+v, err=%+v", string(message.Value), message, err)
}

var ClientID int64

func getClientID() int64 {
	return atomic.AddInt64(&ClientID, 1)
}
func (demo *DemoKafka) demoReader() {
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

func (demo *DemoKafka) demoFetch() {
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

func (demo *DemoKafka) demoWrite() {
	w := kafka.Writer{
		Addr:     kafka.TCP(getAddrsList()...),
		Topic:    demoGetTopic(),
		Balancer: nil,
		//MaxAttempts:  0,
		//BatchSize:    0,
		//BatchBytes:   0,
		//BatchTimeout: 0,
		//ReadTimeout:  0,
		//WriteTimeout: 0,
		RequiredAcks: -1,
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
