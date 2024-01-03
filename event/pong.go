package event

type PayloadPong struct {
	Message []byte
}

const TypePong = "[internal] pong"
