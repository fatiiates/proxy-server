package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Proxy struct {
	addr          string
	services      []string
	requestNumber int
	sem           chan int
	pickOrigin    func() string
}

func (p *Proxy) Start() {
	originServerHandler := http.HandlerFunc(p.ReqHandler)

	if os.Getenv("LOAD_BALANCING_METHOD") == "random" {
		p.pickOrigin = p.PickByRandomService
	} else if os.Getenv("LOAD_BALANCING_METHOD") == "round_robin" {
		p.pickOrigin = p.PickServiceByRoundRobin
	} else {
		p.pickOrigin = p.PickByRandomService
	}

	p.sem = make(chan int, 10)

	log.Fatal(http.ListenAndServe(p.addr, originServerHandler))
}

func (p *Proxy) ReqHandler(rw http.ResponseWriter, req *http.Request) {
	p.sem <- 1

	go func() {
		p.Request(rw, req)
		<-p.sem
	}()
}

func (p *Proxy) Request(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("[reverse proxy server] received request at: %s\n", time.Now())

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
	return p.services[p.requestNumber%len(p.services)]
}

func main() {
	p := &Proxy{
		addr: "localhost:8081",
	}
	p.Start()
}
