package cs

import (
	"encoding/json"

	"github.com/go-zoox/core-utils/object"
	"github.com/go-zoox/tag"
	"github.com/go-zoox/tag/datasource"
)

type Event struct {
	ID      string       `json:"id"`
	Type    string       `json:"type"`
	Payload EventPayload `json:"payload"`
	Error   string       `json:"error"`
}

func (e *Event) Decode(raw []byte) error {
	return json.Unmarshal(raw, e)
}

func (e *Event) Encode() ([]byte, error) {
	return json.Marshal(e)
}

type EventPayload map[string]any

func (ep EventPayload) Get(key string) any {
	return object.Get(ep, key)
}

func (ep EventPayload) Bind(v interface{}) error {
	return tag.New("json", datasource.NewMapDataSource(ep)).Decode(v)
}
