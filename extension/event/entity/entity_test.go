package entity

import (
	"testing"
)

func TestEvent(t *testing.T) {
	// encode and decode
	before := Event{
		ID:   "1",
		Type: "type",
		Payload: map[string]any{
			"key": "value",
		},
	}

	b, err := before.Encode()
	if err != nil {
		t.Fatal(err)
	}

	var after Event
	if err := after.Decode(b); err != nil {
		t.Fatal(err)
	}

	if after.ID != before.ID {
		t.Fatalf("expect: %s, got: %s", before.ID, after.ID)
	}

	if after.Type != before.Type {
		t.Fatalf("expect: %s, got: %s", before.Type, after.Type)
	}

	if after.Payload["key"] != before.Payload["key"] {
		t.Fatalf("expect: %v, got: %v", before.Payload, after.Payload)
	}
}
