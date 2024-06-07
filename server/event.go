package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-zoox/websocket/conn"
)

type Event struct {
	Type    string       `json:"type"`
	Payload EventPayload `json:"payload"`
}

func (e *Event) Decode(raw []byte) error {
	parts := strings.SplitN(string(raw), ",", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid event: %s", string(raw))
	}

	e.Type = parts[0]
	e.Payload = EventPayload{
		Raw: parts[2],
	}

	return nil
}

type EventPayload struct {
	Raw string
}

func (ep *EventPayload) Decode(data any) error {
	if err := json.Unmarshal([]byte(ep.Raw), data); err != nil {
		return fmt.Errorf("failed to parse event payload(%s): %s", ep.Raw, err)
	}

	return nil
}

func (s *server) Event(name string, fn func(conn conn.Conn, payload *EventPayload) error) {
	s.cbs.events[name] = append(s.cbs.events[name], fn)
}
