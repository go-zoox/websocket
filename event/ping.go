package event

type PayloadPing struct {
	Message []byte
}

const TypePing = "[internal] ping"
