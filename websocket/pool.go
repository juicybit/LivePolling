package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// Start will constantly liste to all messages on any channel and act accordingly
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
            client.ID = fmt.Sprintf("id-%d",len(pool.Clients) + 1)
			fmt.Println("registering client:", client.ID)
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: "info", Body: "New User Joined - " + client.ID})
			}
		case client := <-pool.Unregister:
			fmt.Println("unregistering client:", client.ID)
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: "info", Body: "User Disconnected - " + client.ID})
			}
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println("error broadcasting:", err)
					return
				}
			}
		}
	}
}
