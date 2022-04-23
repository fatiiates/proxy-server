package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Proxy Proxy `yaml:"proxy,omitempty"`
}

type Proxy struct {
	Method   string    `yaml:"method,omitempty"`
	Listen   Host      `yaml:"listen,omitempty"`
	Services []Service `yaml:"services,omitempty"`
}

type Host struct {
	Address string `yaml:"address,omitempty"`
	Port    int    `yaml:"port,omitempty"`
}

type Service struct {
	Name   string `yaml:"name,omitempty"`
	Domain string `yaml:"domain,omitempty"`
	Hosts  []Host `yaml:"hosts,omitempty"`
}

func LoadConfigurations(p *ReverseProxy, fn string) {
	file, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal([]byte(file), &config)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	p.services = make(map[string]struct {
		Host          []string
		requestNumber int
	})

	// validation listen address and port
	if config.Proxy.Listen.Address == "" {
		log.Fatalf("there is no 'listen' address")
	}

	if config.Proxy.Listen.Port == 0 {
		log.Fatalf("there is no 'port' for the listen address")
	}

	p.addr = config.Proxy.Listen.Address + ":" + fmt.Sprint(config.Proxy.Listen.Port)

	// validation for the services
	for i, service := range config.Proxy.Services {
		if service.Domain == "" {
			log.Fatalf("there is no 'domain' in the %d. service", i+1)
		}

		if service.Name == "" {
			log.Fatalf("there is no domain 'name' in the %d. service", i+1)
		}

		hosts := []string{}
		// validation for the service hosts
		for j, host := range service.Hosts {
			if host.Address == "" {
				log.Fatalf("there is no 'address' in the %d. services' %d. host", i+1, j+1)
			}

			if host.Port == 0 {
				log.Fatalf("there is no 'port' in the %d. services' %d. host", i+1, j+1)
			}

			hosts = append(hosts, ConvertHostToString(host))
		}

		if _, ok := p.services[service.Domain]; ok {
			log.Fatalf("services domains can not be same string. %d and %d have the same domain name.", p.services[service.Domain].requestNumber+1, i+1)
		}

		p.services[service.Domain] = struct {
			Host          []string
			requestNumber int
		}{
			Host:          hosts,
			requestNumber: i,
		}
		p.requestNumber = append(p.requestNumber, 0)
	}

	if config.Proxy.Method == "random" {
		p.pickOrigin = p.PickByRandomService
	} else if config.Proxy.Method == "round-robin" {
		p.pickOrigin = p.PickServiceByRoundRobin
	} else {
		p.pickOrigin = p.PickByRandomService
	}
}

func ConvertHostToString(h Host) string {
	URL, err := url.Parse(h.Address)
	if err != nil {
		panic(err)
	}
	address := h.Address
	if URL.Scheme == "" {
		address = "http://" + address
	}
	address = address + ":" + fmt.Sprint(h.Port)
	return address
}
