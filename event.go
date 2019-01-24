package slackevents

import (
	"encoding/json"
	"fmt"
)

const (
	EventTypeURLVerification = "url_verification"
	EventTypeCallback        = "event_callback"
)

type HasType struct {
	Type string `json:"type"`
}

type Event struct {
	*HasType
	ParsedEvent interface{}
}

func (event *Event) UnmarshalJSON(data []byte) error {
	eventType := &HasType{}
	err := json.Unmarshal(data, eventType)
	if err != nil {
		return err
	}
	event.HasType = eventType
	switch event.Type {
	case EventTypeURLVerification:
		event.ParsedEvent = &URLVerification{}
	case EventTypeCallback:
		event.ParsedEvent = &Callback{}
	default:
		return fmt.Errorf("Could not parese event type %s", event.Type)
	}
	err = json.Unmarshal(data, event.ParsedEvent)
	return err
}
