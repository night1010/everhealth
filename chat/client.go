package chat

import (
	"encoding/json"
	"time"

	"github.com/night1010/everhealth/logger"
	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event
	chatRoom   uint
	userId     uint
}

func NewClient(conn *websocket.Conn, manager *Manager, chatRoom uint, userId uint) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		chatRoom:   chatRoom,
		userId:     userId,
	}
}

func (c *Client) ReadMessages() {
	defer func() {
		c.manager.RemoveClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		logger.Log.Errorf("failed to set read deadline: %v", err)
		return
	}

	c.connection.SetPongHandler(func(data string) error {
		return c.connection.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Log.Errorf("error reading message: %v", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			logger.Log.Errorf("failed to unmarshal: %s", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			logger.Log.Errorf("Error handling message: %s", err)
		}
	}
}

func (c *Client) WriteMessages() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		c.manager.RemoveClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					logger.Log.Errorf("connection closed: %s", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				logger.Log.Errorf("failed to marshal json: %s", err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.Log.Errorf("failed to write message: %s", err)
			}
		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Log.Errorf("failed to ping: %s", err)
			}
		}
	}
}
