package server

import (
	"github.com/go-zoox/logger"
	"github.com/go-zoox/safe"
	"github.com/go-zoox/websocket"
	"github.com/go-zoox/websocket/conn"
	"github.com/go-zoox/websocket/extension/event/entity"
)

type Server interface {
	On(event string, fn func(conn conn.Conn, payload entity.EventPayload, callback func(error, entity.EventPayload)))
}

type server struct {
	ws websocket.Server
	//
	cbs map[string][]func(conn conn.Conn, payload entity.EventPayload, callback func(error, entity.EventPayload))
}

func New(core websocket.Server) Server {
	cbs := make(map[string][]func(conn conn.Conn, payload entity.EventPayload, callback func(err error, ep entity.EventPayload)))

	core.OnTextMessage(func(c conn.Conn, message []byte) error {
		eventRequest := &entity.Event{}
		if err := eventRequest.Decode(message); err != nil {
			return err
		}

		if fns, ok := cbs[eventRequest.Type]; ok {
			for _, fn := range fns {
				go (func(fn func(c conn.Conn, payload entity.EventPayload, callback func(err error, payload entity.EventPayload))) {
					err := safe.Do(func() error {
						fn(c, eventRequest.Payload, func(err error, payload entity.EventPayload) {
							eventResponse := &entity.Event{
								ID:      eventRequest.ID,
								Type:    eventRequest.Type,
								Payload: payload,
							}
							if err != nil {
								eventResponse.Error = err.Error()
							}

							if er, errx := eventResponse.Encode(); errx != nil {
								logger.Errorf("[event][id: %s] failed to encode event: %v", eventRequest.ID, errx)
							} else {
								err := c.WriteTextMessage(er)
								if err != nil {
									logger.Errorf("[event][id: %s] failed to write text message: %v", eventRequest.ID, err)
								}
							}
						})

						return nil
					})
					if err != nil {
						logger.Errorf("[event][id: %s] failed to handle event(type: %s): %v", eventRequest.ID, eventRequest.Type, err)
					}
				})(fn)
			}
		} else if !ok {
			logger.Errorf("supported event type: %s", eventRequest.Type)
		}

		return nil
	})

	return &server{
		ws:  core,
		cbs: cbs,
	}
}

func (s *server) On(event string, fn func(conn conn.Conn, payload entity.EventPayload, callback func(error, entity.EventPayload))) {
	s.cbs[event] = append(s.cbs[event], fn)
}
