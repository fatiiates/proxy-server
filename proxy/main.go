package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

type ReverseProxy struct {
	addr     string
	services map[string]struct {
		Host          []string
		requestNumber int
	}
	requestNumber []uint32
	pickOrigin    func(*http.Request) string
}

var (
	CONFIG_FILE = "/etc/proxy/config.yaml"
)

func (p *ReverseProxy) Start() {
	originServerHandler := http.HandlerFunc(p.ReqHandler)

	if _, err := os.Stat(CONFIG_FILE); err != nil {
		log.Fatalf("There is no configuration file on the %s", CONFIG_FILE)
	}

	LoadConfigurations(p, CONFIG_FILE)

	log.Fatal(http.ListenAndServe(p.addr, originServerHandler))
}

func (p *ReverseProxy) ReqHandler(rw http.ResponseWriter, req *http.Request) {

	p.Request(rw, req)
}

func (p *ReverseProxy) Request(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("[Reverse Proxy server] received request at: %s\n", time.Now())

	if req.Proto == "HTTP/2.0" {
		rw.WriteHeader(http.StatusHTTPVersionNotSupported)
		return
	}
	o := p.pickOrigin(req)
	if o == "" {
		rw.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(rw, "service is not found")
		return
	}
	URL, err := url.Parse(o)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(rw, err)
		return
	}

	// set req Host, URL and Request URI to forward a request to the origin server
	req.Host = URL.Host
	req.URL.Host = URL.Host
	req.URL.Scheme = URL.Scheme
	req.RequestURI = ""

	// save the response from the origin server
	originServerResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(rw, err)
		return
	}

	// return response to the client
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, originServerResponse.Body)
}

func (p *ReverseProxy) PickByRandomService(req *http.Request) string {
	// picking a service by randomly
	host, _, _ := net.SplitHostPort(req.Host)
	if _, ok := p.services[host]; !ok {
		return ""
	}
	return p.services[host].Host[rand.Intn(len(p.services[host].Host))]
}

func (p *ReverseProxy) PickServiceByRoundRobin(req *http.Request) string {
	// picking a service by round robin method
	host, _, _ := net.SplitHostPort(req.Host)
	if _, ok := p.services[host]; !ok {
		return ""
	}

	n := atomic.AddUint32(&p.requestNumber[p.services[host].requestNumber], 1)
	return p.services[host].Host[int(n)%len(p.services[host].Host)]
}

func main() {
	p := &ReverseProxy{}
	p.Start()
}
