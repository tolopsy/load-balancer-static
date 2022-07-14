package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type server struct {
	addr    string
	isAlive bool
	mu      sync.RWMutex
	proxy   *httputil.ReverseProxy
}

func newServer(addr string) *server {
	serverURL, err := url.Parse(addr)
	if err != nil {
		log.Fatalln("Error while parsing address: ", err.Error())
	}

	return &server{
		addr:    addr,
		isAlive: false,
		mu:      sync.RWMutex{},
		proxy:   httputil.NewSingleHostReverseProxy(serverURL),
	}
}

func initiateNewServers(addrs []string) []*server {
	serverCount := len(addrs)
	servers := make([]*server, 0, serverCount)
	for _, addr := range addrs {
		servers = append(servers, newServer(addr))
	}
	return servers
}

func (s *server) setIsAlive(isAlive bool) {
	s.mu.Lock()
	s.isAlive = isAlive
	s.mu.Unlock()
}

func (s *server) getIsAlive() bool {
	s.mu.RLock()
	isAlive := s.isAlive
	s.mu.RUnlock()
	return isAlive
}

func (s *server) serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}
