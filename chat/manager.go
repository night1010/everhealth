package chat

import (
	"sync"

	"github.com/night1010/everhealth/usecase"
)

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers    map[EventType]EventHandler
	chatUsecase usecase.ChatUsecase
}

func NewManager(chatUsecase usecase.ChatUsecase) *Manager {
	m := &Manager{
		clients:     make(ClientList),
		handlers:    make(map[EventType]EventHandler),
		chatUsecase: chatUsecase,
	}
	m.setupEventHandler()
	return m
}

func (m *Manager) AddClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		_ = client.connection.Close()
		delete(m.clients, client)
	}
}

func (m *Manager) setupEventHandler() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventSendImage] = SendMessageHandler
	m.handlers[EventSendPdf] = SendMessageHandler
	m.handlers[EventTypingSignal] = SendTypingSignalHandler
}

func (m *Manager) routeEvent(e Event, c *Client) error {
	if handler, ok := m.handlers[e.Type]; ok {
		if err := handler(e, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}
