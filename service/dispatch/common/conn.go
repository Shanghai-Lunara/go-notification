package common

import (
	"context"
	"net"
	"time"
)

const (
	ConnAlive = iota
	ConnClosed
)

type Conn struct {
	id           int32
	registerTime time.Time
	heartTime    time.Time
	status       int
	addr         net.Addr
	ctx          context.Context
	cancel       context.CancelFunc
}

func (h *Hub) NewConn(id int32, addr net.Addr) *Conn {
	ctx, cancel := context.WithCancel(h.ctx)
	c := &Conn{
		id:           id,
		registerTime: time.Now(),
		heartTime:    time.Now(),
		status:       ConnAlive,
		addr:         addr,
		ctx:          ctx,
		cancel:       cancel,
	}
	return c
}

func (h *Hub) keepAlive(c *Conn) {
	tick := time.NewTicker(time.Second * time.Duration(h.c.Dispatch.HeartBeatInternal))
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if time.Now().Sub(c.heartTime) > 0 {
				c.status = ConnClosed
				h.close(c)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Conn) ping() {
	c.heartTime = time.Now()
}
