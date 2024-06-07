package cs

import "testing"

func TestEvent(t *testing.T) {
	// encode and decode
	before := Event{
		ID:   "1",
		Type: "type",
		Payload: EventPayload{
			Data: "payload",
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

	var v string
	if err := after.Payload.Decode(&v); err != nil {
		t.Fatal(err)
	}

	if after.Payload.Data != before.Payload.Data {
		t.Fatalf("expect: %v, got: %v", before.Payload.Data, after.Payload.Data)
	}
}
