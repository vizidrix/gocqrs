package clientsockets

import (
	"github.com/vizidrix/gocqrs/cqrs"
)

func HandleClientSockets(clients *ClientView, subscriptionchan chan *ClientConn) func(*websocket.Conn) {
	connections := make(map[uint64]*ClientConn)
	addchan := make(chan *ClientConn)
	removechan := make(chan *ClientConn)

	go func() {
		for {
			select {
			case clientconn := <-addchan:
				connections[clientconn.Client] = clientconn
				subscriptionchan <- clientconn
			case clientconn := <-removechan:
				select {
				case <-clientconn.ExitChan:
				default:
					fmt.Printf("\nClosing out client for session: %s", clientconn.Client)
					close(clientconn.ExitChan)
					delete(connections, clientconn.Client)
				}
			}
		}
	}()

	return func(conn *websocket.Conn) {
		defer func() { conn.Close() }()
		session := conn.Request().FormValue("session")

		client, err := clients.GetBySession(session)
		//		fmt.Println("Starting new client...")
		if err != nil {
			clienterr := NewClientError("invalid_session")
			//			fmt.Printf("\nError validating session: %v", err)
			websocket.JSON.Send(conn, clienterr)
			return
		} else {
			cliententry := NewClientConn(session, client)
			addchan <- cliententry

			//			fmt.Printf("\nNew connection %s", session)
			//			fmt.Printf("\nConnection %s connecting client infrastructure...", sessionid)

			go func() {
				//				defer func() { fmt.Println("Ending client event stream") }()
				for {
					select {
					case event := <-cliententry.EventChan:
						if err := websocket.JSON.Send(conn, event); err != nil {
							fmt.Printf("\nError sending to Client:\n\t%v", err)
							removechan <- cliententry
							return
						}
					case <-cliententry.ExitChan:
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
						removechan <- cliententry
						return
					}
					select {
					case cliententry.MessageChan <- message:

					case <-cliententry.ExitChan:
						return
					}
				}
			}()

			<-cliententry.ExitChan
		}
	}
}
