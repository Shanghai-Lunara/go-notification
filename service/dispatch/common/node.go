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
		}
		nh.nodes[addr] = n
	}
	return nh
}

func (s *Service) handleGetAllocatedNode() {
	
}

func (s *Service) handlePingNode(addr string) {

}

func (s *Service) handleCompleteNode(addr string) {

}