package zap_logger

import (
	"runtime"
	"sync"
	"testing"
)

func testCaller() {

	for i := 0; i < 100000; i++ {
		funcName, file, line, ok := runtime.Caller(0)
		if ok {
			_ = funcName
			_ = file
			_ = line
			_ = ok

		} else {
			panic(i)
		}
	}

}

func TestConcurrentRuntimeCaller(t *testing.T) {

	t.Logf("NumCPU: %d, GOMAXPROCS: %d\n", runtime.NumCPU(), runtime.GOMAXPROCS(-1))

	t.Logf("start")
	wg := sync.WaitGroup{}

	for i := 0; i < 2000; i++ {
		wg.Add(1)
		go func() {
			t.Log("start test")
			defer wg.Done()
			testCaller()
		}()
	}
	wg.Wait()
	t.Logf("end")

}
