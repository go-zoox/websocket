package event

type PayloadClose struct {
	Code    int
	Message string
}

const TypeClose = "[internal] close"
