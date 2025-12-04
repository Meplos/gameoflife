package main

import (
	"log"
	"net/http"

	"github.com/Meplos/GameOfLife/client"
	"github.com/gorilla/websocket"
)

const FPS = 60

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("[WS HANDLER] Incomming connexion")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS HANDLER] Cant upgrader conn: %v", err)
		return
	}
	c := client.NewClient(conn)
	defer client.UnregisterClients(c)
	client.RegisterClients(c)
	go c.ExecCommand()
	c.Listen()

}

func main() {
	log.Printf("hello\n")
	http.HandleFunc("/ws", handler)
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("gol: error %v\n", err)
	}

}
