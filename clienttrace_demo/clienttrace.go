package clienttrace_demo

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httptrace"
)

func Demo1() {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.Println("GetConn, hostPort=", hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			log.Printf("GotConn: GotConnInfo=%+v\n", connInfo)
		},
		PutIdleConn: func(err error) {
			log.Printf("PutIdleConn: err=%+v\n", err)
		},
		GotFirstResponseByte: func() {
			log.Println("GotFirstResponseByte")
		},
		Got100Continue: nil,
		Got1xxResponse: nil,
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Printf("DNSStart: DNSStartInfo=%+v\n", info)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Printf("DNSDone: DNSDoneInfo=%+v\n", dnsInfo)
		},
		ConnectStart: func(network, addr string) {
			log.Printf("ConnectStart: network=%s, addr=%s\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			log.Printf("ConnectDone: network=%s, addr=%s, err=%+v\n", network, addr, err)
		},
		TLSHandshakeStart: func() {
			log.Println("TLSHandshakeStart")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			log.Printf("TLSHandshakeDone: ConnectionState=%+v, err=%+v\n", state, err)
		},
		WroteHeaderField: func(key string, value []string) {
			log.Printf("WroteHeaderField: key=%s, value=%+v\n", key, value)
		},
		WroteHeaders: func() {
			log.Println("WroteHeaders")
		},
		Wait100Continue: nil,
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			log.Printf("WroteRequest: WroteRequestInfo=%s\n", info)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	_, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
}
