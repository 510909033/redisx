package my_httptrace

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"sync"
	"time"
)

var Client *http.Client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
		Proxy:             http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second, // 连接超时
			KeepAlive: 10 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:       120 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	Timeout: 500 * time.Millisecond,
}

func GenLogId() string {
	h2 := md5.New()
	rand.Seed(time.Now().Unix())
	str := fmt.Sprintf("%d%d%d", os.Getpid(), time.Now().UnixNano(), rand.Int())

	h2.Write([]byte(str))
	uniqid := hex.EncodeToString(h2.Sum(nil))

	return uniqid
}

func DemoMyHttpTrace() {

	//启动http服务
	go listenHttp()
	time.Sleep(time.Second)

	var (
		wg           sync.WaitGroup
		maxParallel  int       = 1
		parallelChan chan bool = make(chan bool, maxParallel)
	)
	for {
		parallelChan <- true
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				<-parallelChan
			}()
			testHttp2()
		}()
	}
	wg.Wait()
}

//启动一个http服务
func listenHttp() {
	log.Println("启动http服务")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})
	_ = http.ListenAndServe(":11222", nil) // <-今天讲的就是这个ListenAndServe是如何工作的
	//https://blog.csdn.net/gophers/article/details/37815009
}

func testHttp2() {
	url := "http://localhost:11222/index.php"
	req, _ := http.NewRequest("GET", url, nil)

	timeStart := time.Now()
	consume := func() int64 {
		return time.Since(timeStart).Nanoseconds()
	}

	uniqId := GenLogId()
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.Println("GetConn id:", uniqId, consume(), hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Println("GotConn id:", uniqId, consume(), connInfo.Conn.LocalAddr())
		},

		ConnectStart: func(network, addr string) {
			log.Println("ConnectStart id:", uniqId, consume(), network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			log.Println("ConnectDone id:", uniqId, consume(), network, addr, err)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := Client.Do(req)
	if err != nil {
		log.Println("err: id", uniqId, consume(), err)
		return
	}

	log.Println("all ", consume())

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("error", string(body))
	}
	return

}
