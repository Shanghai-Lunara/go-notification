package common

import (
	"context"
	"go-notification/config"
	"net"
	"sync"
	"sync/atomic"
)

type Hub struct {
	mu          sync.RWMutex
	c           *config.Config
	autoId      *int32
	connections map[int32]*Conn
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewHub(conf *config.Config, ctx context.Context) *Hub {
	var a int32
	subCtx, cancel := context.WithCancel(ctx)
	h := &Hub{
		c:      conf,
		autoId: &a,
		ctx:    subCtx,
		cancel: cancel,
	}
	return h
}

func (h *Hub) getClientId() int32 {
	return atomic.AddInt32(h.autoId, 1)
}

func (h *Hub) close(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.connections, c.id)
}

func (s *Service) handleInit(addr net.Addr, ext string) int32 {
	return 1
}

func (s *Service) handlePing(id int32) {

}
