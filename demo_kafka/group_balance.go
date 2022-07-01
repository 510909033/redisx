package demo_kafka

import (
	"github.com/segmentio/kafka-go"
	"sort"
)

/**
发送时 根据条件（消息相关）可以选择分区
	type BalanceFn func(msg kafka.Message, partitions ...int) int

消费时 通过实现GroupBalancer接口，可以自定义每个消费者消费哪个分区
	缺点： 无法获取到每个消费者的唯一标记， member.ID 虽然唯一，但是是系统生成的

++++

https://blog.csdn.net/java_atguigu/article/details/123920233
Kafka 要保证消息的消费顺序，可以有2种方法：
一、1个Topic（主题）只创建1个Partition(分区)，这样生产者的所有数据都发送到了一个Partition(分区)，保证了消息的消费顺序。
二、生产者在发送消息的时候指定要发送到哪个Partition(分区)。


*/

// 可以实现： 消费消息时，每个member 消费哪几个分区
type MyGroupBalancer struct{}

func (r MyGroupBalancer) ProtocolName() string {
	return "my_group_balancer_protocol_name"
}

func (r MyGroupBalancer) UserData() ([]byte, error) {
	return nil, nil
}

func (r MyGroupBalancer) AssignGroups(members []kafka.GroupMember, topicPartitions []kafka.Partition) kafka.GroupMemberAssignments {
	groupAssignments := kafka.GroupMemberAssignments{}

	membersByTopic := findMembersByTopic(members)

	logDebug("AssignGroups, membersByTopic=%+v", membersByTopic)

	for topic, members := range membersByTopic {
		partitions := findPartitions(topic, topicPartitions)

		partitionCount := len(partitions)
		memberCount := len(members)
		logDebug("AssignGroups, topic=%s, partitionCount=%d, memberCount=%d,所有分区=%+v",
			topic, partitionCount, memberCount, partitions)

		for memberIndex, member := range members {
			assignmentsByTopic, ok := groupAssignments[member.ID]
			if !ok {
				assignmentsByTopic = map[string][]int{}
				groupAssignments[member.ID] = assignmentsByTopic
			}

			minIndex := memberIndex * partitionCount / memberCount
			maxIndex := (memberIndex + 1) * partitionCount / memberCount

			for partitionIndex, partition := range partitions {
				if partitionIndex >= minIndex && partitionIndex < maxIndex {
					//todo 都指定到4分区
					//partition = 4

					assignmentsByTopic[topic] = append(assignmentsByTopic[topic], partition)
				}
			}
		}
	}

	logDebug("AssignGroups, 结果=%+v", groupAssignments)

	return groupAssignments
}

func findMembersByTopic(members []kafka.GroupMember) map[string][]kafka.GroupMember {
	membersByTopic := map[string][]kafka.GroupMember{}
	for _, member := range members {
		for _, topic := range member.Topics {
			membersByTopic[topic] = append(membersByTopic[topic], member)
		}
	}

	// normalize ordering of members to enabling grouping across topics by partitions
	//
	// Want:
	// 		C0 [T0/P0, T1/P0]
	// 		C1 [T0/P1, T1/P1]
	//
	// Not:
	// 		C0 [T0/P0, T1/P1]
	// 		C1 [T0/P1, T1/P0]
	//
	// Even though the later is still round robin, the partitions are crossed
	//
	for _, members := range membersByTopic {
		sort.Slice(members, func(i, j int) bool {
			return members[i].ID < members[j].ID
		})
	}

	return membersByTopic
}

func findPartitions(topic string, partitions []kafka.Partition) []int {
	var ids []int
	for _, partition := range partitions {
		if partition.Topic == topic {
			ids = append(ids, partition.ID)
		}
	}
	return ids
}
