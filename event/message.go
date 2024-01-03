package event

type PayloadMessage struct {
	Type    int
	Message []byte
}

const TypeMessage = "[internal] message"
