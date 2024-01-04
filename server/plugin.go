package server

import (
	"fmt"

	"github.com/go-zoox/websocket/server/plugin"
)

func (s *server) Plugin(plugin plugin.Plugin) error {
	if _, exist := s.plugins[plugin.Name()]; exist {
		return fmt.Errorf("plugin %s already exist", plugin.Name())
	}

	s.plugins[plugin.Name()] = plugin
	return nil
}
