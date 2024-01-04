package heartbeat

import (
	"time"

	"github.com/go-zoox/logger"
	"github.com/go-zoox/websocket/conn"
)

func (hb *heartbeat) Apply(conn conn.Conn) error {
	heartbeatCh := make(chan struct{})

	// @TODO auto listen ping + sennd pong
	conn.OnPong(func(message []byte) error {
		heartbeatCh <- struct{}{}
		return nil
	})

	// heartbeat
	go func() {
		time.After(hb.opt.Interval)
		for {
			select {
			case <-conn.Context().Done():
				logger.Debugf("[heartbeat][interval] context done => cancel")
				return
			case <-time.After(hb.opt.Interval):
				logger.Debugf("[heartbeat][interval] send heartbeat ->")
				if err := conn.Ping(nil); err != nil {
					logger.Errorf("[heartbeat][interval] fail to send heartbeat: %v", err)
					close(heartbeatCh)
					go conn.Close()
					return
				}

				select {
				case <-conn.Context().Done():
					logger.Debugf("[heartbeat][timeout] context done => cancel")
					return
				case <-time.After(hb.opt.Timeout):
					logger.Errorf("[heartbeat][timeout] fail to wait heartbeat")
					close(heartbeatCh)
					go conn.Close()
					return
				case <-heartbeatCh:
					logger.Debugf("[heartbeat][timeout] receive heartbeat <-")
				}
			}
		}
	}()

	return nil
}
