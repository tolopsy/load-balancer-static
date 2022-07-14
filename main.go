package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverAddresses := strings.Split(os.Getenv("SERVER_ADDRESSES"), ",")
	lbPort := os.Getenv("LOAD_BALANCER_ADDR")
	servers := initiateNewServers(serverAddresses)

	loadBalancer := newLoadBalancer(lbPort, servers)
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.serveProxy(w, r)
	}
	go loadBalancer.runHealthCheck()

	http.HandleFunc("/", handleRedirect)
	fmt.Printf("Listening at port %s\n", loadBalancer.port)
	if err := http.ListenAndServe(":" + loadBalancer.port, nil); err != nil {
		log.Fatalln(err.Error())
	}
}
