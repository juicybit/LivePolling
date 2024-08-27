package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool {return true},
}

// Upgrades the connection to a Websocket Connection
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error){
  return upgrader.Upgrade(w, r, nil)
}

