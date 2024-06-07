package cs

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Event struct {
	ID      string       `json:"id"`
	Type    string       `json:"type"`
	Payload EventPayload `json:"payload"`
}

func (e *Event) Decode(raw []byte) error {
	parts := strings.SplitN(string(raw), ",", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid event: %s", string(raw))
	}

	e.ID = parts[0]
	e.Type = parts[1]
	e.Payload = EventPayload{
		Raw: parts[2],
	}

	return nil
}

func (e *Event) Encode() ([]byte, error) {
	if _, err := e.Payload.Encode(); err != nil {
		return nil, fmt.Errorf("failed to encode event payload in event encode: %s", err)
	}

	return []byte(strings.Join([]string{
		e.ID,
		e.Type,
		e.Payload.Raw,
	}, ",")), nil
}

type EventPayload struct {
	Raw  string
	Data any
}

func (ep *EventPayload) Decode(data any) error {
	if err := json.Unmarshal([]byte(ep.Raw), data); err != nil {
		return fmt.Errorf("failed to decode event payload(%s): %s", ep.Raw, err)
	}

	ep.Data = data

	return nil
}

func (ep *EventPayload) Encode() (bytes []byte, err error) {
	if bytes, err = json.Marshal(ep.Data); err != nil {
		return nil, fmt.Errorf("failed to encode event payload(%v): %s", ep.Data, err)
	} else {
		ep.Raw = string(bytes)
		return bytes, nil
	}
}
