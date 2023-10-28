package serial

import (
	"bytes"
	"encoding/json"
	"main/pkg/api"
)

type EventWrapper struct {
	Type  EventType     `json:"event_type"`
	Event api.ApiUpdate `json:"event"`
}

func WrapEvent(eventType EventType, event api.ApiUpdate) *EventWrapper {
	return &EventWrapper{eventType, event}
}

func (e *EventWrapper) Serialize() ([]byte, error) {
	var b bytes.Buffer
	var data []byte
	var err error

	// could have alternate protocols
	data, err = json.Marshal(e)
	if err != nil {
		return nil, err
	}
	b.Write(data)
	return b.Bytes(), nil
}

func Deserialize(data []byte) (*EventWrapper, error) {
	out := &EventWrapper{}
	err := json.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
