package demo_kafka

import (
	"sync"
	"testing"
)

func TestDemoKafka_GetPartitions(t *testing.T) {
	demo := &DemoKafka{}
	demo.demoGetPartitions()
}

func TestDemoKafka_demoReader(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		demo := &DemoKafka{}
		demo.demoWrite()
	}()
	for i := 0; i < 5; i++ {
		go func() {
			defer wg.Done()
			demo := &DemoKafka{}
			demo.demoReader()
		}()
	}
	wg.Wait()
}

func TestDemoKafka_demoReaderOne(t *testing.T) {

	demo := &DemoKafka{}
	demo.demoReader()
}

func TestDemoKafka_demoFetch(t *testing.T) {
	demo := &DemoKafka{}
	demo.demoFetch()
}

func TestDemoKafka_demoWrite(t *testing.T) {
	demo := &DemoKafka{}
	for i := 0; i < 100; i++ {
		demo.demoWrite()
	}
}
