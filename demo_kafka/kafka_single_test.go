package demo_kafka

import "testing"

func TestDemoSingle_Run(t *testing.T) {

	s := &DemoSingle{}
	s.Run()
}

func TestDemoSingle_createTopic(t *testing.T) {

	s := &DemoSingle{}
	s.createTopic()
}
