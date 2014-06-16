package clientsockets

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/vizidrix/gocqrs/net/clients"
)

func HandleClientSockets(clientsessions *clients.ClientSessionsView, subscriptionchan chan *Connection) func(*websocket.Conn) {
	connections := make(map[uint64]*Connection)
	addchan := make(chan *Connection)
	removechan := make(chan *Connection)

	go func() {
		for {
			select {
			case connection := <-addchan:
				if _, active := connections[connection.client]; active {
					removechan <- connection
				} else {
					connections[connection.client] = connection
					subscriptionchan <- connection
				}

			case connection := <-removechan:
				select {
				case <-connection.exitChan:
				default:
					fmt.Printf("\nClosing out client for session: %s", connection.client)
					close(connection.exitChan)
					delete(connections, connection.client)
				}
			}
		}
	}()

	return func(conn *websocket.Conn) {
		defer func() { conn.Close() }()
		session := conn.Request().FormValue("session")

		client, err := clientsessions.GetBySession(session)

		if err != nil {
			fmt.Printf("\nError validating session: %v", err)
			//	clienterr := err.NewError("invalid_session")
			//	websocket.JSON.Send(conn, clienterr)
			return
		} else {
			connection := NewConnection(session, client)
			addchan <- &connection

			//			fmt.Printf("\nNew connection %s", session)
			//			fmt.Printf("\nConnection %s connecting client infrastructure...", sessionid)

			go func() {
				//				defer func() { fmt.Println("Ending client event stream") }()
				for {
					select {
					case event := <-connection.eventChan:
						if err := websocket.JSON.Send(conn, event); err != nil {
							fmt.Printf("\nError sending to Client:\n\t%v", err)
							removechan <- &connection
							return
						}
					case <-connection.exitChan:
						return
					}
				}
			}()

			go func() {
				for {
					var message []byte
					if err := websocket.JSON.Receive(conn, &message); err != nil {
						fmt.Printf("\nReceived %+v from Client", message)
						fmt.Printf("\nError receiving from Client:\n\t%v", err)
						removechan <- &connection
						return
					}
					select {
					case connection.messageChan <- message:
						fmt.Printf("\n%v", message)
					case <-connection.exitChan:
						return
					}
				}
			}()

			<-connection.exitChan
		}
	}
}
