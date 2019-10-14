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
	wg        sync.WaitGroup
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
			w.wg.Done()
			return
		case <-tick.C:
			var (
				pStr string
				err  error
			)
			if pStr, err = w.dao.LPopOne(w.addr); err != nil {
				log.Println("appendLoop LPopOne err:", err)
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
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	for {
		select {
		case <-w.ctx.Done():
			w.wg.Done()
			return
		case <-tick.C:
			log.Println("listNodes-len:", len(w.listNodes.Players))
			if t, ok := w.listNodes.Players[0]; ok {
				if t.RLink != nil {
					p := t.RLink.Player
					if p.Value > int(time.Now().Unix()) {
						log.Printf("listNodes continue 1111 v:%d time:%d \n", p.Value, int(time.Now().Unix()))
						continue
					}
					log.Println("listNodes continue 22222 p-Player:", p)
					if err := w.CheckOne(t.RLink.Player.Pid); err != nil {
						log.Printf("CheckOne pid:%d p:%v err:%v \n", p.Pid, p, err)
					}
					log.Println("listNodes continue 33333")
				}
			}
		}
	}
}

func (s *Service) initWorker(w *Worker, id int) {
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tick.C:
			log.Println("initWorker loop worker_id:", id)
			if addr, err := s.getAllocatedNode(); err != nil {
				continue
			} else {
				log.Println("initWorker addr:", addr)
				if addr == "" {
					continue
				}
				w.addr = addr
				w.wg.Add(2)
				go w.appendLoop()
				go w.logicLoop()
				return
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
		go s.initWorker(w.workers[i], i)
	}
	return w
}

func (s *Service) CloseWorkers() {
	for _, v := range s.workers.workers {
		v.wg.Wait()
	}
}
