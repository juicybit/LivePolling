package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

// Read constantly listens for Messages on the Clients connection
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		mt, p, err := c.Conn.ReadMessage()
		if err != nil {
		  log.Println(err)
		  return
		}
		msg := &Message{}
		if err := c.Conn.ReadJSON(msg); err != nil {
			log.Println("error reading json: ", msg, err)
			return
		}
		fmt.Printf("%s > %+v\n", c.ID, msg)
		message := Message{Type: string(mt), Body: string(p)}
		// c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
