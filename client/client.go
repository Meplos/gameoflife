package client

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Meplos/GameOfLife/board"
	"github.com/Meplos/GameOfLife/gol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id     uuid.UUID
	Conn   *websocket.Conn
	Send   chan board.BoardState
	Active bool
	Game   *gol.Game
}

type CmdOption struct {
	H float64 `json:"h"`
	W float64 `json:"w"`
}

type ClientCommand struct {
	Cmd     string    `json:"cmd"`
	Options CmdOption `json:"options"`
}

func NewClient(conn *websocket.Conn) *Client {
	c := Client{
		Id:     uuid.New(),
		Conn:   conn,
		Send:   make(chan board.BoardState),
		Active: false,
	}
	return &c
}

type ActiveClients struct {
	clients map[uuid.UUID]*Client
}

func (c *Client) Listen() {
	for b := range c.Send {
		log.Printf("Message receive")
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

func (client *Client) SendInitState(b *board.Board) {
	client.Conn.WriteJSON(b.ToInitialState())

}

func UnregisterClients(client *Client) {
	client.Active = false
	defer client.Conn.Close()
	delete(actives.clients, client.Id)
}

func (c *Client) ExecCommand() {
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
		case "init":
			c.InitGame(incomming.Options)
		case "pause":
			c.PauseGame()
		case "play":
			go c.StartGame()
		case "restart":
			c.RestartGame()
		}
	}

}

func (c *Client) InitGame(o CmdOption) {
	c.Game = gol.NewGame(uint(o.W), uint(o.H))
	c.Game.Init()

}
func (c *Client) PauseGame() {
	c.Game.B.Pause()
}

func (c *Client) StartGame() {
	log.Printf("[client.Start] c: %v, game: %v\n", c.Id, c.Game.ID)
	c.Game.Start(c.Send)
}

func (c *Client) RestartGame() {
	c.Game.Restart()
}
