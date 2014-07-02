package webclientsockets

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/vizidrix/gocqrs/net/webclients"
)

func HandleWebClientSockets(webclientsessions *webclients.WebClientSessionsView, subscriptionchan chan WebClientConnection) func(*websocket.Conn) {
	connservice := NewConnectionService(subscriptionchan)

	go func() {
		for {
			ManageConnections(&connservice)
		}
	}()

	return func(conn *websocket.Conn) {
		defer func() { conn.Close() }()
		session := conn.Request().FormValue("session")

		webclient, err := webclientsessions.GetBySession(session)

		if err != nil {
			fmt.Printf("\nError validating session: %v", err)
			//	webclienterr := err.NewError("invalid_session")
			//	websocket.JSON.Send(conn, webclienterr)
			return
		} else {
			connection := NewConnection(session, webclient)
			connservice.addChan <- &connection

			go func() {
				for {
					if active := HandleWebClientEvent(&connservice, &connection, conn); !active {
						return
					}
				}
			}()

			go func() {
				for {
					if active := HandleWebClientMessage(&connservice, &connection, conn); !active {
						return
					}
				}
			}()

			<-connection.exitChan
		}
	}
}

func ManageConnections(connservice *ConnectionService) {
	select {
	case connection := <-connservice.addChan:
		AddConnection(connservice, connection)
	case connection := <-connservice.removeChan:
		RemoveConnection(connservice, connection)
	}
}

func AddConnection(connservice *ConnectionService, connection *ConnectionMemento) {
	//fmt.Printf("\nRegistering ConnectionMemento: %d", connection.webclient)
	if conn, active := connservice.connections[connection.webclient]; active {
		RemoveConnection(connservice, conn)
	}
	connservice.connections[connection.webclient] = connection
	connservice.subscriptionChan <- connection
}

func RemoveConnection(connservice *ConnectionService, connection *ConnectionMemento) {
	select {
	case <-connection.exitChan:
	default:
		//fmt.Printf("\nClosing out webclient for session: %d", connection.webclient)
		close(connection.exitChan)
		delete(connservice.connections, connection.webclient)
	}
}

func HandleWebClientEvent(connservice *ConnectionService, connection *ConnectionMemento, conn *websocket.Conn) bool {
	select {
	case event := <-connection.eventChan:
		if err := websocket.JSON.Send(conn, event); err != nil {
			fmt.Printf("\nError sending to WebClient:\n\t%v", err)
			connservice.removeChan <- connection
			return false
		}
		return true
	case <-connection.exitChan:
		return false
	}
}

func HandleWebClientMessage(connservice *ConnectionService, connection *ConnectionMemento, conn *websocket.Conn) bool {
	var message []byte
	if err := websocket.JSON.Receive(conn, &message); err != nil {
		fmt.Printf("\nError receiving from WebClient:\n\t%v", err)
		connservice.removeChan <- connection
		return false
	}
	select {
	case connection.messageChan <- message:
		return true
	case <-connection.exitChan:
		return false
	}
}

/*
func HandleWebClientSockets(webclientsessions *webclients.WebClientSessionsView, subscriptionchan chan *ConnectionMemento) func(*websocket.Conn) {
	connections := make(map[uint64]*ConnectionMemento)
	addchan := make(chan *ConnectionMemento, 1)
	removechan := make(chan *ConnectionMemento, 1)

	go func() {
		for {
			select {
			case connection := <-addchan:
				fmt.Printf("\nRegistering ConnectionMemento: %d", connection.webclient)
				if _, active := connections[connection.webclient]; active {
					removechan <- connection
				} else {
					connections[connection.webclient] = connection
					subscriptionchan <- connection
				}
				fmt.Printf(("\nNew ConnectionMemento: %d"), connection.webclient)
			case connection := <-removechan:
				select {
				case <-connection.exitChan:
				default:
					fmt.Printf("\nClosing out webclient for session: %d", connection.webclient)
					close(connection.exitChan)
					delete(connections, connection.webclient)
				}
			}
		}
	}()

	return func(conn *websocket.Conn) {
		defer func() { conn.Close() }()
		session := conn.Request().FormValue("session")

		webclient, err := webclientsessions.GetBySession(session)

		if err != nil {
			fmt.Printf("\nError validating session: %v", err)
			//	webclienterr := err.NewError("invalid_session")
			//	websocket.JSON.Send(conn, webclienterr)
			return
		} else {
			connection := NewConnectionMemento(session, webclient)
			addchan <- &connection

			//			fmt.Printf("\nNew connection %s", session)
			//			fmt.Printf("\nConnectionMemento %s connecting webclient infrastructure...", sessionid)

			go func() {
				//				defer func() { fmt.Println("Ending webclient event stream") }()
				for {
					select {
					case event := <-connection.eventChan:
						if err := websocket.JSON.Send(conn, event); err != nil {
							fmt.Printf("\nError sending to WebClient:\n\t%v", err)
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
						fmt.Printf("\nReceived %+v from WebClient", message)
						fmt.Printf("\nError receiving from WebClient:\n\t%v", err)
						removechan <- &connection
						return
					}
					fmt.Printf("\nMessage received from webclient %d: %v", connection.webclient, message)
					select {
					case connection.messageChan <- message:
						fmt.Printf("\nMessage from webclient %d passed to connection handler", connection.webclient)
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
