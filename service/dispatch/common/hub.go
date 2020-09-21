package common

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/Shanghai-Lunara/go-notification/config"
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
		c:           conf,
		autoId:      &a,
		connections: make(map[int32]*Conn, 0),
		ctx:         subCtx,
		cancel:      cancel,
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

func (s *Service) handleInit(addr net.Addr, lastId int32) int32 {
	s.hub.mu.Lock()
	defer s.hub.mu.Unlock()
	var (
		id int32
	)
	if lastId <= *s.hub.autoId && lastId != 0 {
		id = lastId
	} else {
		id = s.hub.getClientId()
	}
	c := s.hub.NewConn(id, addr)
	s.hub.connections[id] = c
	go s.hub.keepAlive(c)
	return c.id
}

func (s *Service) handlePing(id int32) {
	if t, ok := s.hub.connections[id]; ok {
		t.ping()
	}
}
