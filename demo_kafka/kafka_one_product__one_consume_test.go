package demo_kafka

import (
	"github.com/segmentio/kafka-go"
	"math/rand"
	"testing"
	"time"
)

func TestDeleteTopic(t *testing.T) {
	demo := &DemoKafkaOneProductOneConsume{}
	demo.deleteTopic("Topic__v3")
	demo.deleteTopic("Topic__v4")
	demo.deleteTopic("Topic__v5")
	demo.deleteTopic("Topic__v6")
}

func TestCreateTopic(t *testing.T) {
	demo := &DemoKafkaOneProductOneConsume{}
	topic := demo.getTopicName()
	//创建之前好像不能执行 exist命令，否则会创建默认的3个分区的
	//if demo.topicExists(topic) {
	//	return
	//}
	demo.deleteTopic(topic)
	demo.createTopic(topic, 5)
}

func TestDemoKafkaOneProductOneConsume_createTopic_v1(t *testing.T) {
	demo := &DemoKafkaOneProductOneConsume{}
	topic := demo.getTopicName()

	//if demo.topicExists(topic) {
	//	time.Sleep(time.Second)
	//	deleteTopic := demo.deleteTopic(topic)
	//	assert.True(t, deleteTopic)
	//}

	ticker := time.NewTicker(time.Second * 13)
	timer := time.NewTimer(time.Second * 30)
	times := 0
	for {
		times++
		select {
		case val := <-ticker.C:
			logDebug("尝试createTopic, 第%d次，time=%s", times, val)
			ticker.Stop()
			ticker.Stop()
			if !demo.createTopic(topic, 12) {
				ticker.Reset(time.Second * 3)
			}

			goto over
		case val := <-timer.C:
			logDebug("触发timer， 停止, 第%d次，time=%s", times, val)
			timer.Stop()
			timer.Stop()
			goto over
		}
	}
over:
}

func TestDemoKafkaOneProductOneConsume_writeMessageByPartitions(t *testing.T) {
	demo := &DemoKafkaOneProductOneConsume{}
	topic := demo.getTopicName()

	//demo.writeMessageByPartitions(topic, 10, &kafka.Hash{})
	var fn BalanceFn
	demo.writeMessageByPartitions(topic, 10000, fn)
}

// 可以实现：发送消息时，选择发送到哪个分区
type BalanceFn func(msg kafka.Message, partitions ...int) int

func (fn BalanceFn) Balance(msg kafka.Message, partitions ...int) (partition int) {

	return 4
	logDebug("hit BalanceFn")
	rand.Seed(time.Now().UnixNano())
	return partitions[rand.Intn(len(partitions))]
}

func TestDemoKafkaOneProductOneConsume_consumeMessageByPartitions(t *testing.T) {
	demo := &DemoKafkaOneProductOneConsume{}
	topic := demo.getTopicName()
	groupId := demo.getGroupId()
	num := 5
	//groupId := ""
	for i := 0; i < 20; i++ {
		go func() {
			demo.consumeMessageByPartitions(topic, groupId, 0, num)
		}()
	}
	time.Sleep(time.Second * 15000)
}

func TestDemoKafkaOneProductOneConsume_ONE_Partition_Three_Consumer(t *testing.T) {
	demo := &DemoKafkaOneProductOneConsume{}
	topic := "ONE_Partition_Three_Consumer"
	groupId := "one_three_group"

	demo.createTopic(topic, 1)
	time.Sleep(time.Second * 2)

	demo.writeMessageByPartitions(topic, 5, nil)

	time.Sleep(time.Second * 2)

	num := 5
	for i := 0; i < 5; i++ {
		go func() {
			demo.consumeMessageByPartitions(topic, groupId, 0, num)
		}()
	}
	time.Sleep(time.Second * 15000)
}

func TestDemoKafkaOneProductOneConsume_writeByPartitions(t *testing.T) {
	demo := &DemoKafkaOneProductOneConsume{}
	topic := "ONE_Partition_Three_Consumer"

	demo.writeByPartitions(topic)
}
