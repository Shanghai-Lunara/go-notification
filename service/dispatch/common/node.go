package common

import (
	"fmt"
	"sync"
	"time"
)

const (
	NodeIdle = iota
	NodeInUse
)

type Node struct {
	addr        string
	allocatedId int32
	heartTime   time.Time
	status      int
	count       int
}

type NodeHub struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

func (s *Service) NewNodeHub() *NodeHub {
	nh := &NodeHub{
		nodes: make(map[string]*Node, len(s.dao.Redis.RedisConf)),
	}
	for _, v := range s.dao.Redis.RedisConf {
		addr := fmt.Sprintf("%s:%s", v[0], v[1])
		n := &Node{
			addr:        addr,
			allocatedId: 0,
			heartTime:   time.Now(),
			status:      NodeIdle,
			count:       0,
		}
		nh.nodes[addr] = n
	}
	go s.loop()
	return nh
}

func (s *Service) loop() {
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tick.C:
			s.nodeHub.mu.Lock()
			for _, v := range s.nodeHub.nodes {
				if v.status == NodeIdle || v.allocatedId == 0 {
					continue
				}
				if t, ok := s.hub.connections[v.allocatedId]; ok {
					if t.status == ConnClosed {
						v.allocatedId = 0
						v.status = NodeIdle
					}
				} else {
					v.allocatedId = 0
					v.status = NodeIdle
				}
			}
			s.nodeHub.mu.Unlock()
		}
	}
}

func (s *Service) handleGetAllocatedNode(id int32) string {
	s.nodeHub.mu.Lock()
	defer s.nodeHub.mu.Unlock()
	var b []string
	for k, v := range s.nodeHub.nodes {
		if v.status == NodeInUse {
			continue
		}
		b = append(b, k)
	}
	for i := 0; i < len(b); i++ {
		for j := 0; j < len(b)-1-i; j++ {
			if s.nodeHub.nodes[b[j]].count >= s.nodeHub.nodes[b[j+1]].count {
				b[j], b[j+1] = b[j+1], b[j]
			}
		}
	}
	if len(b) == 0 {
		return ""
	}
	s.nodeHub.nodes[b[0]].allocatedId = id
	s.nodeHub.nodes[b[0]].status = NodeInUse
	s.nodeHub.nodes[b[0]].heartTime = time.Now()
	return b[0]
}

func (s *Service) handlePingNode(id int32, addr string) {
	s.nodeHub.mu.Lock()
	defer s.nodeHub.mu.Unlock()
	if t, ok := s.nodeHub.nodes[addr]; ok {
		if t.allocatedId == id {
			t.heartTime = time.Now()
		}
	}
}

func (s *Service) handleCompleteNode(id int32, addr string) {
	s.nodeHub.mu.Lock()
	defer s.nodeHub.mu.Unlock()
	if t, ok := s.nodeHub.nodes[addr]; ok {
		if t.allocatedId == id {
			t.allocatedId = 0
			t.status = NodeIdle
			t.count++
		}
	}
}
