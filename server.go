package main

import (
	"log"
	"net/http"

	"github.com/thomasschuiki/LivePolling/server/websocket"
)

func clientHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("browsing: ", r.URL.Path)
	p := "." + r.URL.Path
	if p == "./" {
		p = "static/client.html"
	}
	http.ServeFile(w, r, p)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("browsing: ", r.URL.Path)
	p := "." + r.URL.Path
	if p == "./admin/" {
		p = "static/admin.html"
	}
	http.ServeFile(w, r, p)
}

func wsHandler(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	client := &websocket.Client{Conn: ws, Pool: pool}
	pool.Register <- client
	client.Read()
}

func main() {
	pool := websocket.NewPool()
	go pool.Start()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", clientHandler)
	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(pool, w, r)
	})

	log.Println("Starting Server on Port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	log.Println("Done listening on Port 8080")
}
