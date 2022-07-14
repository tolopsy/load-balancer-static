package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"time"
)

const (
	OK   string = "ok"
	DEAD string = "dead"
)

func isAliveOrNot(u *url.URL) bool {
	conn, err := net.DialTimeout("tcp", u.Host, time.Second*30)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (l *loadBalancer) healthCheck() {
	for _, server := range l.servers {
		pingURL, err := url.ParseRequestURI(server.addr)
		if err != nil {
			log.Printf("Error parsing server address (%s): %s", server.addr, err.Error())
			continue
		}

		isAlive := isAliveOrNot(pingURL)
		server.setIsAlive(isAlive)
		status := OK
		if !isAlive {
			status = DEAD
		}
		log.Printf("Status of '%s': %s", server.addr, status)
	}
}

func (l *loadBalancer) runHealthCheck() {
	fmt.Println("Starting health check")
	l.healthCheck()
	for {
		<-time.After(time.Second * 30)
		l.healthCheck()
	}
}
