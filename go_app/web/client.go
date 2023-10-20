package web

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *Hub

	Conn *websocket.Conn

	Send chan []byte
}

func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		Conn: conn,
		Send: make(chan []byte),
	}
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.broadcast <- msg
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

func (h *Handler) HandleWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(hub, conn)
	hub.register <- client

	go client.Read()
	go client.Write()
}
