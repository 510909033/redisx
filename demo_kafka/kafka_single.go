package demo_kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type DemoSingle struct {
	Num int64
	sync.Mutex
}

func (s *DemoSingle) createTopic() {
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
		if p.Topic == s.getTopic() {
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

	//err = controllerConn.DeleteTopics(topicName)
	//if err != nil {
	//	log.Fatalf("DeleteTopics fail, err=%+v", err)
	//}

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:              s.getTopic(),
		NumPartitions:      1,
		ReplicationFactor:  1,
		ReplicaAssignments: nil,
		ConfigEntries:      nil,
	})
	if err != nil {
		log.Fatalf("CreateTopics fail, err=%+v", err)
	}
	log.Printf("创建topic无报错, topic=%s\n", s.getTopic())
}

func (s *DemoSingle) getTopic() string {
	return "demo_single_v2"
}
func (s *DemoSingle) getGroupId() string {
	return "group_v2"
}

func (s *DemoSingle) read() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: getAddrsList(),
		//Partition: 1,
		GroupID:     s.getGroupId(),
		GroupTopics: nil,
		Topic:       demoGetTopic(),
		//Partition:   0,
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
	var mu sync.Mutex

	for {
		mu.Lock()
		m, err := r.ReadMessage(context.Background())

		mu.Unlock()
		if err != nil {
			panic(err)
			break
		}

		fmt.Printf("message at ClientID:%s, Partition:%d, offset %d: %s = %s\n",
			r.Stats().ClientID,
			m.Partition, m.Offset, string(m.Key), string(m.Value))
		//fmt.Printf("message %+v\n", m.Time)

	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("over")
}

func (s *DemoSingle) Run() {
	w := kafka.Writer{
		Addr:         kafka.TCP(getAddrsList()...),
		Topic:        s.getTopic(),
		Balancer:     nil,
		MaxAttempts:  0,
		BatchSize:    0,
		BatchBytes:   0,
		BatchTimeout: 0,
		ReadTimeout:  0,
		WriteTimeout: 0,
		//MaxAttempts:  0,
		//BatchSize:    0,
		//BatchBytes:   0,
		//BatchTimeout: 0,
		//ReadTimeout:  0,
		//WriteTimeout: 0,
		RequiredAcks: -1,
		Async:        false,
		Completion:   nil,
		Compression:  0,
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("Logger , "+msg, args...)
		}),
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("ErrorLogger,"+msg, args...)
		}),
		Transport: nil,
	}
	defer func() {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	newCtx, cancelFunc := context.WithTimeout(ctx, time.Millisecond)
	_ = newCtx
	_ = cancelFunc

	err := w.WriteMessages(newCtx, kafka.Message{
		//Topic:         demoGetTopic(),
		//Partition:     0,
		//Offset:        0,
		//HighWaterMark: 0,
		//Key:           nil,
		Value: []byte(time.Now().Format("2006-01-02 15:04:05")),
		//Headers:       nil,
		//Time:          time.Time{},
	})
	if err != nil {
		panic(err)
	}

}
