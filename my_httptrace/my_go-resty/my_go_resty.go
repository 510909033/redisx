package my_go_resty

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"sync"
	"time"
)

var mu sync.Mutex
var transport *http.Transport

func getDemoTransport() *http.Transport {
	if transport != nil {
		return transport
	}
	mu.Lock()
	defer mu.Unlock()
	transport = &http.Transport{
		Proxy:       nil,
		DialContext: nil,
		Dial:        nil,
		//DialTLSContext:         nil,
		DialTLS:               nil,
		TLSClientConfig:       nil,
		TLSHandshakeTimeout:   time.Second * 5,
		DisableKeepAlives:     false,
		DisableCompression:    false,
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   2,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       0,
		ResponseHeaderTimeout: 0,
		ExpectContinueTimeout: time.Second * 10,
		TLSNextProto:          nil,
		ProxyConnectHeader:    nil,
		//GetProxyConnectHeader:  nil,
		MaxResponseHeaderBytes: 0,
		WriteBufferSize:        0,
		ReadBufferSize:         0,
		ForceAttemptHTTP2:      false,
	}
	return transport
}

func DemoMyGoResty() {
	url := "http://localhost:12345/demo-php-http-server-12345.php"
	url = "http://172.16.7.228:12345/demo-php-http-server-12345.php"
	url = "http://172.20.10.40:9001/#/auth"
	url = "http://172.20.10.40:12345/demo-php-http-server-12345.php"
	url = "http://172.20.10.40:12345/"
	//url = "https://github.com/go-resty/resty"
	//url = `https://xx.cnblogs.com/semanticAnalysis/api2?app_id=2882303761517406029&token=5651740635029&timestamp=1493188096125&queries=[{"query":"打开越野跑","confidence":0.8}]&device_id=robot_aaaUG6LExb9Fsc80taeUv&device={"ip":"192.168.1.1","network":"wifi"}&session={"is_new":true,"id":"","start_timestamp":1493188096125}&request_id=afaaaa&version=2.1`
	url = `http://tool.lu/timestamp/`
	url = `http://space.babytree-inc.com/pages/viewpage.action?pageId=30636605`
	client := resty.New()
	transport := getDemoTransport()

	client.SetTransport(transport)
	client.SetTimeout(time.Second * 3)

	resp, err := client.R().
		EnableTrace().
		Get(url)

	if err == nil {
		ti := resp.Request.TraceInfo()
		dumpTrace(ti)
		//fmt.Println("  ConnTime      :", ti.ConnTime)
		//fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
		//fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
		//fmt.Println("  IsConnReused  :", ti.IsConnReused)
		//fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
		//fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)

		return
	}
	// Explore response object
	//fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode())
	//fmt.Println("  Status     :", resp.Status())
	//fmt.Println("  Proto      :", resp.Proto())
	//fmt.Println("  Time       :", resp.Time())
	//fmt.Println("  Received At:", resp.ReceivedAt())
	//fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	//fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	//fmt.Println("  DNSLookup     :", ti.DNSLookup)
	//fmt.Println("  ConnTime      :", ti.ConnTime)
	//fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	//fmt.Println("  ServerTime    :", ti.ServerTime)
	//fmt.Println("  ResponseTime  :", ti.ResponseTime)
	//fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	//fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())

}

func dumpTrace(ti resty.TraceInfo) {
	//fmt.Println("Request Trace Info:")
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
	fmt.Println("")
}
