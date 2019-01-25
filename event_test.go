package slackevents_test

import (
	"encoding/json"
	"testing"

	"github.com/FreekingDean/slackevents"
	"github.com/stretchr/testify/assert"
)

func TestEventUnmarshalJSON(t *testing.T) {
	urlVerification := `{"type":"url_verification","challenge":"test"}`
	event := &slackevents.Event{}
	err := json.Unmarshal([]byte(urlVerification), event)
	assert.NoError(t, err)
	assert.Equal(t, slackevents.EventTypeURLVerification, event.Type)
	assert.IsType(t, &slackevents.URLVerification{}, event.ParsedEvent)
	assert.Equal(
		t,
		"test",
		event.ParsedEvent.(*slackevents.URLVerification).Challenge,
	)

	callbackEvent := `{"type":"event_callback","team_id":"TEAM ID","api_app_id":"API APP ID","authed_users":[],"event_id":"EVENT ID","event_time":"1001.199","event":{"type":"message"}}`
	event = &slackevents.Event{}
	err = json.Unmarshal([]byte(callbackEvent), event)
	assert.NoError(t, err)
	assert.Equal(t, slackevents.EventTypeCallback, event.Type)
	assert.IsType(t, &slackevents.Callback{}, event.ParsedEvent)
	callback := event.ParsedEvent.(*slackevents.Callback)
	assert.Equal(
		t,
		"API APP ID",
		callback.APIAppID,
	)
	assert.Equal(t, 1001, int(callback.EventTime.Unix()))
	assert.Equal(t, "message", callback.InnerEvent.Type)
	_, err = json.Marshal(event)
	assert.NoError(t, err)
}
