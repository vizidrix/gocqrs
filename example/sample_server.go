package main

import (
	"github.com/vizidrix/gocqrs/net/server"
	"fmt"
)

func main() {
	server := server.NewTCPServer()
	closeChan := make(chan struct{})
	server.ListenOn(8080, closeChan)
	fmt.Printf("Server running...\n")
	<-closeChan
}