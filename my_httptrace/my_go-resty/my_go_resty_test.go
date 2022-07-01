package my_go_resty

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDemoMyGoRestyOne(t *testing.T) {
	/*
		=== RUN   TestDemoMyGoRestyOne
		  DNSLookup     : 32.0005ms
		  ConnTime      : 217.9168ms
		  TCPConnTime   : 38.0668ms
		  TLSHandshake  : 144.916ms
		  ServerTime    : 47.999ms
		  ResponseTime  : 0s
		  TotalTime     : 265.9158ms
		  IsConnReused  : false
		  IsConnWasIdle : false
		  ConnIdleTime  : 0s
		  RequestAttempt: 1
		  RemoteAddr    : 118.31.180.41:443

		  DNSLookup     : 0s
		  ConnTime      : 0s
		  TCPConnTime   : 0s
		  TLSHandshake  : 0s
		  ServerTime    : 43.792ms
		  ResponseTime  : 0s
		  TotalTime     : 43.792ms
		  IsConnReused  : true
		  IsConnWasIdle : true
		  ConnIdleTime  : 0s
		  RequestAttempt: 1
		  RemoteAddr    : 118.31.180.41:443
	*/
	DemoMyGoResty()
	//DemoMyGoResty()
}

func TestDemoMyGoResty(t *testing.T) {
	i := int64(0)
	var wg sync.WaitGroup
	var ch = make(chan int64, 3)
	panic("aa")
	for {
		wg.Add(1)
		//log.Println(i)
		ch <- i
		go func() {
			defer func() {
				<-ch
				wg.Done()
			}()
			DemoMyGoResty()
			time.Sleep(time.Second)
		}()
		atomic.AddInt64(&i, 1)
	}
	wg.Wait()
}
