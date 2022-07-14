package main

import (
	"fmt"
	"net/http"
	"sync"
)

type loadBalancer struct {
	port            string
	roundRobinCount int
	mu              sync.Mutex
	servers         []*server
}

func newLoadBalancer(port string, servers []*server) *loadBalancer {
	return &loadBalancer{
		port:            port,
		roundRobinCount: 0,
		mu:              sync.Mutex{},
		servers:         servers,
	}
}

func (l *loadBalancer) getNextAvailableServer() *server {
	serversCount := len(l.servers)
	l.mu.Lock()
	selectedServer := l.servers[l.roundRobinCount%serversCount]

	for !selectedServer.getIsAlive() {
		l.roundRobinCount++
		selectedServer = l.servers[l.roundRobinCount%serversCount]
	}
	l.roundRobinCount++
	l.mu.Unlock()
	return selectedServer
}

func (l *loadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := l.getNextAvailableServer()
	r.Header.Del("X-Forwarded-For") // to prevent IP spoofing

	fmt.Printf("Proxying request to %s: %s\n", targetServer.addr, r.RequestURI)
	targetServer.serve(w, r)
}
