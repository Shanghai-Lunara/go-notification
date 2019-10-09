package common

import (
	"context"
	"go-notification/dao"
	"log"
	"sync"
	"time"
)

type Worker struct {
	mu        sync.RWMutex
	dao       *dao.Dao
	addr      string
	count     int
	status    int
	listNodes *ListNodes
	ctx       context.Context
}

func (w *Worker) appendLoop() {
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-tick.C:
			var (
				pStr string
				err  error
			)
			if pStr, err = w.dao.LPopOne(w.addr); err != nil {
				log.Print("appendLoop LPopOne err:", err)
				time.Sleep(time.Second * 1)
				continue
			}
			if pStr == "" {
				time.Sleep(time.Millisecond * 500)
				continue
			}
			if err = w.RefreshOne(pStr); err != nil {
				log.Print("appendLoop RefreshOne err:", err)
			}
		}
	}
}

func (w *Worker) logicLoop() {
	tick := time.NewTicker(time.Millisecond * 1)
	defer tick.Stop()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-tick.C:
			if t, ok := w.listNodes.Players[0]; ok {
				if t.RLink != nil {
					t.RLink.Player.mu.Lock()
					if t.RLink.Player.Value > int(time.Now().Unix()) {
						t.RLink.Player.mu.Unlock()
						continue
					}
				}
			}
		}
	}
}

func (s *Service) initWorker(w *Worker) {
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tick.C:
			if addr, err := s.getAllocatedNode(); err != nil {
				continue
			} else {
				if addr == "" {
					continue
				}
				w.addr = addr
				go w.appendLoop()
				go w.logicLoop()
				break
			}
		}
	}
}

type Workers struct {
	workers map[int]*Worker
}

func (s *Service) NewWorkers() *Workers {
	w := &Workers{
		workers: make(map[int]*Worker, s.c.Logic.WorkerNum),
	}
	for i := 0; i < s.c.Logic.WorkerNum; i++ {
		w.workers[i] = &Worker{
			dao:       s.dao,
			addr:      "",
			count:     0,
			listNodes: InitList(),
			ctx:       s.ctx,
		}
		go s.initWorker(w.workers[i])
	}
	return w
}
