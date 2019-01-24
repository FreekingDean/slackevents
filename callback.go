package slackevents

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	InnerEventTypeMessage = "message"
)

type CallbackHandler func(*Callback) error

var DefaultCallbackServer = &CallbackServer{}

type CallbackShared struct {
	TeamID      string    `json:"team_id"`
	APIAppID    string    `json:"api_app_id"`
	AuthedUsers []string  `json:"authed_users"`
	EventID     string    `json:"event_id"`
	EventTime   *UnixTime `json:"event_time"`
}

type InnerEvent struct {
	*HasType
	ParsedEvent interface{}
}

type Message struct {
	*HasType
	Text      string    `json:"text"`
	Channel   string    `json:"channel"`
	User      string    `json:"user"`
	Timestamp *UnixTime `json:"ts"`
}

type Callback struct {
	*CallbackShared
	InnerEvent *InnerEvent `json:"event"`
}

type UnixTime struct {
	time.Time
}

type CallbackServer struct {
	MessageHandler func(*Message) error
}

func (server *CallbackServer) Handler(event *Callback) error {
	switch innerEvent := event.InnerEvent.ParsedEvent.(type) {
	case *Message:
		return server.MessageHandler(innerEvent)
	}
	return fmt.Errorf("Could not handle event of type %T", event.InnerEvent)
}

func (t *UnixTime) UnmarshalJSON(data []byte) error {
	parts := strings.Split(strings.Trim(string(data), "\""), ".")
	i, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return err
	}
	nt := time.Unix(i, 0)
	t.Time = nt
	return nil
}

func (inner *InnerEvent) UnmarshalJSON(data []byte) error {
	eventType := &HasType{}
	err := json.Unmarshal(data, eventType)
	if err != nil {
		return err
	}
	inner.HasType = eventType
	switch inner.Type {
	case InnerEventTypeMessage:
		inner.ParsedEvent = &Message{}
	default:
		return fmt.Errorf("Could not parese event type %s", inner.Type)
	}
	err = json.Unmarshal(data, inner.ParsedEvent)
	return err
}
