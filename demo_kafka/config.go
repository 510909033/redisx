package demo_kafka

import "strings"

// {"code":200,"data":[{"configId":161,"createTime":1573482215000,"environmentId":1,"fileName":"Golang_Kafka","id":9012,"key":"babytree.zk.addrs","project":"Golang_Kafka","secure":false,"updateTime":1649320687000,"value":"172.20.12.3:8281,172.20.12.55:8281,172.20.12.117:8281","version":"master"},{"configId":161,"createTime":1573482186000,"environmentId":1,"fileName":"Golang_Kafka","id":9011,"key":"babytree.addrs","project":"Golang_Kafka","secure":false,"updateTime":1649235436000,"value":"172.20.12.3:9011,172.20.12.55:9011,172.20.12.117:9011","version":"master"},{"configId":161,"createTime":1649401826000,"environmentId":1,"fileName":"Golang_Kafka","id":12898,"key":"alikafka.brokers","project":"Golang_Kafka","secure":false,"updateTime":1649401826000,"value":"172.20.12.3:9011,172.20.12.55:9011,172.20.12.117:9011","version":"master"}]}
func getAddrs() string {
	//return "172.20.12.3:9011,172.20.12.55:9011,172.20.12.117:9011"
	return "172.20.12.3:9011"
}

func getAddrsList() []string {
	return strings.Split("172.20.12.3:9011,172.20.12.55:9011,172.20.12.117:9011", ",")
}

func demoGetTopic() string {
	return "demo_wbt_v2"
	//return "pregnant_user_visitor"
}
