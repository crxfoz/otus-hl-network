package delivery

import (
	"encoding/json"
	"time"

	"github.com/crxfoz/otus-hl-network/internal/domain"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TODO: add handler to encapsulate logic around websocket

type Client struct {
	userName string
	hub      *wsHub
	conn     *websocket.Conn
	send     chan []byte
	onSend   func(domain.Message)
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("client unexpectedly closed")
			}
			break
		}

		msgPayload := domain.Message{
			UserName:  c.userName,
			Body:      string(message),
			CreatedAt: time.Now().UTC().UnixMilli(),
		}

		m, err := json.Marshal(msgPayload)
		if err != nil {
			logrus.WithError(err).Error("could not marshal message")
			continue
		}

		logrus.WithField("payload", string(m)).Debug("client broadcasting message")

		c.hub.broadcast <- m

		if c.onSend != nil {
			c.onSend(msgPayload)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			logrus.WithField("payload", string(message)).Debug("client sending message")

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type wsHub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	closed chan struct{}
}

func (h *wsHub) CloseNotify() <-chan struct{} {
	return h.closed
}

func (h *wsHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

			if len(h.clients) == 0 {
				h.closed <- struct{}{}
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}
