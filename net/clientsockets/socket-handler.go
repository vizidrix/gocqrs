package clientsockets

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/vizidrix/gocqrs/net/clients"
)

func HandleClientSockets(clientsessions *clients.ClientSessionsView, subscriptionchan chan *Connection) func(*websocket.Conn) {
	addchan := make(chan *Connection, 1)
	removechan := make(chan *Connection, 1)

	go func() {
		ManageConnections(addchan, removechan, subscriptionchan)
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

			go func() {
				for {
					HandleClientEvent(conn, &connection, removechan)
				}
			}()

			go func() {
				for {
					HandleClientMessage(conn, &connection, removechan)
				}
			}()

			<-connection.exitChan
		}
	}
}

func ManageConnections(addchan chan *Connection, removechan chan *Connection, subscriptionchan chan *Connection) {
	connections := make(map[uint64]*Connection)
	for {
		select {
		case connection := <-addchan:
			fmt.Printf("\nRegistering Connection: %d", connection.client)
			if _, active := connections[connection.client]; active {
				removechan <- connection
			} else {
				connections[connection.client] = connection
				subscriptionchan <- connection
			}
			fmt.Printf(("\nNew Connection: %d"), connection.client)
		case connection := <-removechan:
			select {
			case <-connection.exitChan:
			default:
				fmt.Printf("\nClosing out client for session: %d", connection.client)
				close(connection.exitChan)
				delete(connections, connection.client)
			}
		}
	}
}

func HandleClientEvent(conn *websocket.Conn, connection *Connection, removechan chan *Connection) {
	select {
	case event := <-connection.eventChan:
		if err := websocket.JSON.Send(conn, event); err != nil {
			fmt.Printf("\nError sending to Client:\n\t%v", err)
			removechan <- connection
			return
		}
	case <-connection.exitChan:
		return
	}
}

func HandleClientMessage(conn *websocket.Conn, connection *Connection, removechan chan *Connection) {
	var message []byte
	if err := websocket.JSON.Receive(conn, &message); err != nil {
		fmt.Printf("\nError receiving from Client:\n\t%v", err)
		removechan <- connection
		return
	}
	select {
	case connection.messageChan <- message:

	case <-connection.exitChan:
		return
	}
}

/*
func HandleClientSockets(clientsessions *clients.ClientSessionsView, subscriptionchan chan *Connection) func(*websocket.Conn) {
	connections := make(map[uint64]*Connection)
	addchan := make(chan *Connection, 1)
	removechan := make(chan *Connection, 1)

	go func() {
		for {
			select {
			case connection := <-addchan:
				fmt.Printf("\nRegistering Connection: %d", connection.client)
				if _, active := connections[connection.client]; active {
					removechan <- connection
				} else {
					connections[connection.client] = connection
					subscriptionchan <- connection
				}
				fmt.Printf(("\nNew Connection: %d"), connection.client)
			case connection := <-removechan:
				select {
				case <-connection.exitChan:
				default:
					fmt.Printf("\nClosing out client for session: %d", connection.client)
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
					fmt.Printf("\nMessage received from client %d: %v", connection.client, message)
					select {
					case connection.messageChan <- message:
						fmt.Printf("\nMessage from client %d passed to connection handler", connection.client)
					case <-connection.exitChan:
						return
					}
				}
			}()

			<-connection.exitChan
		}
	}
}
*/
