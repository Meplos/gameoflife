package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Meplos/GameOfLife/board"
	"github.com/Meplos/GameOfLife/client"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var b = board.NewBoard(200, 200)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("[WS HANDLER] Incomming connxion")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS HANDLER] Cant upgrader conn: %v", err)
		return
	}
	c := client.NewClient(conn)
	defer client.UnregisterClients(c)
	client.RegisterClients(c)
	go c.ExecCommand(b)
	c.Listen()

}

func main() {
	log.Printf("hello\n")
	http.HandleFunc("/ws", handler)
	fs := http.FileServer(http.Dir("./public"))

	http.Handle("/", fs)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		b.Restart()
		for t := range ticker.C {
			log.Printf("Tick %v, Alive: %v Running:%v\n", t, b.AliveCount(), !b.IsPaused)
			b.Next()
			client.Broadcast(*b)
		}
	}()

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("gol: error %v\n", err)
	}

}
