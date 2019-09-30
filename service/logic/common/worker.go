package common

import (
	"context"
	"go-notification/dao"
	"time"
)

type Worker struct {
	dao       *dao.Dao
	addr      string
	count     int
	listNodes *ListNodes
	ctx       context.Context
	cancel    context.CancelFunc
}

func (w *Worker) appendLoop() {

}

func (w *Worker) logicLoop() {

}

func (s *Service) initWorker(w *Worker) {
	tick := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-s.ctx.Done():
			tick.Stop()
			return
		case <-tick.C:
			if addr, err := s.getAllocatedNode(); err != nil {
				continue
			} else {
				tick.Stop()
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
