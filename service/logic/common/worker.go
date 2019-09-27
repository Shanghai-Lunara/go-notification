package common

import (
	"context"
	"go-notification/config"
)

type Worker struct {
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

func NewWorkers(conf *config.Config, ctx context.Context) *Workers {
	w := &Workers{
		workers: make(map[int]*Worker, conf.Logic.WorkerNum),
	}
	for i := 0; i < conf.Logic.WorkerNum; i++ {
		w.workers[i] = &Worker{
			addr:    "",
			count:   0,
			players: make(map[int]int, 0),
			ctx:     ctx,
		}
	}
	return w
}
