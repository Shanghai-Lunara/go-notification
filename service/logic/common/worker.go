package common

import (
	"context"
	"go-notification/dao"
	"log"
	"time"
)

type Worker struct {
	dao       *dao.Dao
	addr      string
	count     int
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
			if res, err := w.dao.LPopOne(w.addr); err != nil {
				log.Print("appendLoop LPopOne err:", err)
				time.Sleep(time.Second * 1)
				continue
			} else {
				if res == "" {
					time.Sleep(time.Millisecond * 500)
					continue
				}

			}
		}
	}
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
				if addr == "" {
					continue
				}
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
