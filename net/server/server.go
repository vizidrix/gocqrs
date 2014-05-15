package server

import (
	//"bytes"
	//"bufio"
	"fmt"
	//"io"
	"net"
	"time"
	"errors"
)

type TCPServer interface {
	ListenOn(port uint32, closeChan chan struct{})
	HandleConnection(connection net.Conn) error
	Stop()
}

type cqrsTCPServer struct {
	logChan chan string
	isClosed bool
	closeChan chan struct{}
	indexChan chan uint64
	connections map[uint64]net.Conn
}

func NewTCPServer() TCPServer {
	server := &cqrsTCPServer {
		logChan: make(chan string),
		isClosed: false,
		closeChan: make(chan struct{}),
		indexChan: make(chan uint64),
		connections: make(map[uint64]net.Conn),
	}
	var index uint64 = 0
	go func() {
		for {
			index++
			server.indexChan<-index
			fmt.Printf("Index consumed [ %d ]\n", index)
		}
	}()
	go func() {
		for {
			fmt.Printf("* Log:\n%s\n", <-server.logChan)
		}
	}()
	return server
}

func (svr *cqrsTCPServer) Stop() {
	if svr.isClosed {
		return
	}
	svr.isClosed = true
	close(svr.closeChan)
}

func (svr *cqrsTCPServer) Log(message string) {
	svr.logChan<-message
}

func (svr *cqrsTCPServer) ListenOn(port uint32, closeChan chan struct{}) {
	var listener net.Listener
	quitChanClosed := false
	quitChan := make(chan struct{})
	go func() { // Run a goroutine which will close the connection on cancel chan
		select {
		case <-svr.closeChan:
		case <-closeChan:
		}
		svr.Log("Cancel chan was closed for listener\n")
		if !quitChanClosed {
			quitChanClosed = true
			close(quitChan)
		}
		listener.Close()
	}()
	go func() {
		var err error
		if listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err != nil {
			return
		}
		defer listener.Close()
		defer func() { // Ensure graceful shutdown of all clients if possible
			svr.Log("Closing server connections...\n")
			for _, connection := range svr.connections {
				if connection != nil {
					connection.Close()
				}
			}
		}()
		for {
			connection, err := listener.Accept()
			if err != nil {
				select {
				case <-quitChan:
					return // Close was requested
				default:
				}
				svr.Log(fmt.Sprintf("Error accepting connection: [ %s ]\n", err))
				<-time.After(1 * time.Second)
				continue
			}
			go svr.HandleConnection(connection) // Handle connection in another goroutine
		}
	}()
	return
}

func (svr *cqrsTCPServer) HandleConnection(connection net.Conn) error {
	index := <-svr.indexChan
	defer func() { // Clean up the resource
		connection.Close()
		delete(svr.connections, index)
	}()
	buffer := make([]byte, 4096)
	connection.SetDeadline(time.Now().Add(1 * time.Second))
	size, err := connection.Read(buffer)
	if err != nil {
		svr.Log(fmt.Sprintf("Error reading [ %d ]: [ %s ]\n", index, err))
		return err
	}
	svr.Log(fmt.Sprintf("[ %d ] Read [ %d ] bytes [\n%s\n]\n", index, size, buffer))
	if buffer[0] == 'G' {
		//svr.Stop()
		svr.Log("Serving GET\n")
		connection.Write([]byte(`
<html><head></head><body>
<form action="http://localhost:8080/api/v1/things/addcommand" method="post" enctype="application/x-www-form-urlencoded">
	<input type="text"  name="blabha"></input>
	<button type="submit">Submit</button>
</form>
</body></html>
`))
		return nil
	}
	if buffer[0] == 'P' {
		svr.Log("Serving PUT or POST\n")
		connection.Write([]byte(`{"status":"received","domainid":"20",commandid":"10"}`))
		return nil
	}
	return errors.New("Invalid request\n")
}
