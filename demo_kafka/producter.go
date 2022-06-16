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

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	//for k := range m {
	//	fmt.Println(k)
	//}

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
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

}

var ClientID int64

func getClientID() int64 {
	return atomic.AddInt64(&ClientID, 1)
}
func (demo *DemoKafka) demoReader() {
	log.Println("start ", time.Now().String())
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  getAddrsList(),
		Topic:    demoGetTopic(),
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		//Partition: 1,
		GroupID: "demo_fetch_group",
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
	})
	r.SetOffset(42)

	size := 2000
	currNum := 0

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
			Value: []byte(time.Now().String()),
			//Headers:       nil,
			//Time:          time.Time{},
		})
		if err != nil {
			panic(err)
		}
	}

	stats := w.Stats()
	log.Printf("stats=%+v\n", stats)
}
