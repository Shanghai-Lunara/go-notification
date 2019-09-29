package common

import (
	"context"
	"go-notification/dao"
)

type Worker struct {
	dao     *dao.Dao
	addr    string
	count   int
	players map[int]int
	ctx     context.Context
	cancel  context.CancelFunc
}

func (w *Worker) appendLoop() {

}

func (w *Worker) logicLoop() {

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
			dao:     s.dao,
			addr:    "",
			count:   0,
			players: make(map[int]int, 0),
			ctx:     s.ctx,
		}
		go w.workers[i].appendLoop()
		go w.workers[i].logicLoop()
	}
	return w
}
