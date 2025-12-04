package client

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Meplos/GameOfLife/board"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id     uuid.UUID
	Conn   *websocket.Conn
	Send   chan board.Board
	Active bool
}

type ClientCommand struct {
	Cmd string `json:"cmd"`
}

func NewClient(conn *websocket.Conn) *Client {
	c := Client{
		Id:     uuid.New(),
		Conn:   conn,
		Send:   make(chan board.Board),
		Active: false,
	}
	return &c
}

type ActiveClients struct {
	clients map[uuid.UUID]*Client
}

func (c *Client) Listen() {
	for b := range c.Send {
		c.Conn.WriteJSON(b)
	}
}

var actives = ActiveClients{
	clients: make(map[uuid.UUID]*Client, 0),
}

func RegisterClients(client *Client) {
	log.Printf("[RegisterClient] ID: %v", client.Id.String())
	actives.clients[client.Id] = client
	client.Active = true
}
func UnregisterClients(client *Client) {
	client.Active = false
	defer client.Conn.Close()
	delete(actives.clients, client.Id)
}

func Broadcast(b board.Board) {
	for _, c := range actives.clients {
		c.Send <- b
	}
}

func (c *Client) ExecCommand(b *board.Board) {
	for {
		var incomming ClientCommand
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("client exec: %v", err)
			if errors.Is(err, websocket.ErrCloseSent) {
				UnregisterClients(c)
			}
			return
		}
		log.Printf("Client.Receive %s", msg)
		json.Unmarshal(msg, &incomming)

		switch incomming.Cmd {
		case "pause":
			b.Pause()
		case "play":
			b.Play()
		case "restart":
			b.Restart()
		}
	}

}
