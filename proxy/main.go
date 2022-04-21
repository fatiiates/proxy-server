package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

type Proxy struct {
	addr          string
	services      []string
	requestNumber uint32
	pickOrigin    func() string
}

func (p *Proxy) Start() {
	originServerHandler := http.HandlerFunc(p.ReqHandler)

	p.services = append(p.services, "http://0.0.0.0:8080")
	p.services = append(p.services, "http://0.0.0.0:8081")

	if os.Getenv("LOAD_BALANCING_METHOD") == "random" {
		p.pickOrigin = p.PickByRandomService
	} else if os.Getenv("LOAD_BALANCING_METHOD") == "round_robin" {
		p.pickOrigin = p.PickServiceByRoundRobin
	} else {
		p.pickOrigin = p.PickByRandomService
	}

	log.Fatal(http.ListenAndServe(p.addr, originServerHandler))
}

func (p *Proxy) ReqHandler(rw http.ResponseWriter, req *http.Request) {

	p.Request(rw, req)
}

func (p *Proxy) Request(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())

	if req.Proto != "HTTP/1.1" {
		rw.WriteHeader(http.StatusHTTPVersionNotSupported)
		return
	}

	URL, err := url.Parse(p.pickOrigin())
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
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(rw, err)
		return
	}

	// return response to the client
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, originServerResponse.Body)
}

func (p *Proxy) PickByRandomService() string {
	// picking a service by randomly
	return p.services[rand.Intn(len(p.services))]
}

func (p *Proxy) PickServiceByRoundRobin() string {
	// picking a service by round robin method
	n := atomic.AddUint32(&p.requestNumber, 1)
	service := p.services[int(n)%len(p.services)]

	return service
}

func main() {
	p := &Proxy{
		addr: "localhost:5000",
	}
	p.Start()
}
