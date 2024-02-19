package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/night1010/everhealth/entity"
)

type EventType string

type EventHandler func(e Event, c *Client) error

type Event struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	EventSendMessage  EventType = "send_message"
	EventNewMessage   EventType = "new_message"
	EventSendImage    EventType = "image_message"
	EventSendPdf      EventType = "pdf_message"
	EventTypingSignal EventType = "typing_signal"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    uint   `json:"from"`
}

type TypingSignalEvent struct {
	UserId   uint `json:"user_id"`
	IsTyping bool `json:"is_typing"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

func SendMessageHandler(event Event, c *Client) error {
	var chatEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = chatEvent.Message
	broadMessage.From = chatEvent.From

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = event.Type

	for client := range c.manager.clients {
		if client.chatRoom == c.chatRoom {
			client.egress <- outgoingEvent
		}
	}

	chatMessage := &entity.Chat{
		TelemedicineId: c.chatRoom,
		UserId:         c.userId,
		ChatTime:       broadMessage.Sent,
		Message:        chatEvent.Message,
		MessageType:    entity.MessageTypeText,
	}

	_, _ = c.manager.chatUsecase.AddChatMessage(context.Background(), chatMessage)

	return nil
}

func SendTypingSignalHandler(event Event, c *Client) error {
	var typingEvent TypingSignalEvent
	if err := json.Unmarshal(event.Payload, &typingEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	for client := range c.manager.clients {
		if client.chatRoom == c.chatRoom {
			if client != c {
				client.egress <- event
			}
		}
	}

	return nil
}
