package common

import (
	"context"
	"go-notification/dao"
	"log"
	"sync"
	"time"
)

const (
	WorkerAlive = iota
	WorkerInterrupt
	WorkerClosed
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
				err   error
				p     []string
				count int
			)
			if w.status == WorkerInterrupt {
				time.Sleep(time.Second * 1)
				continue
			}
			if p, err = w.dao.LRange(w.addr, 200); err != nil {
				log.Println("appendLoop LRange err:", err)
				time.Sleep(time.Second * 1)
				continue
			}
			count = 0
			for _, v := range p {
				if w.status == WorkerInterrupt {
					continue
				}
				if w.status == WorkerClosed {
					if count == 0 {
						return
					}
					if err = w.dao.LTRIM(w.addr, count); err != nil {
						log.Print("appendLoop WorkerClosed LTRIM err:", err)
					}
					return
				}
				if v == "" {
					time.Sleep(time.Millisecond * 500)
					continue
				}
				if err = w.RefreshOne(v, true); err != nil {
					log.Print("appendLoop RefreshOne err:", err)
				}
				count++
				//time.Sleep(time.Millisecond * 10)
			}
			if count == 0 {
				time.Sleep(time.Second * 1)
				continue
			}
			if err = w.dao.LTRIM(w.addr, count); err != nil {
				log.Print("appendLoop LTRIM err:", err)
			}
		}
	}
}

func (w *Worker) logicLoop() {
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	tick1 := time.NewTicker(time.Second * 1)
	defer tick1.Stop()
	for {
		select {
		case <-w.ctx.Done():
			w.wg.Done()
			return
		case <-tick.C:
			if w.status == WorkerInterrupt {
				log.Println("logicLoop WorkerInterrupt w.id:", w.addr)
				time.Sleep(time.Second * 1)
				continue
			}
			if t, ok := w.listNodes.Players[0]; ok {
				if t.RLink != nil {
					p := t.RLink.Player
					if p.Value > int(time.Now().Unix()) {
						log.Printf("listNodes continue 1111 v:%d time:%d \n", p.Value, int(time.Now().Unix()))
						continue
					}
					if err := w.CheckOne(t.RLink.Player.Pid); err != nil {
						log.Printf("CheckOne pid:%d p:%v err:%v \n", p.Pid, p, err)
					}
				}
			}
		case <-tick1.C:
			log.Println("logicLoop listNodes-len:", len(w.listNodes.Players))
			if tmp, ok := w.listNodes.Players[1]; ok {
				log.Println("logicLoop Players[1]:", tmp.Player)
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
	mu      sync.RWMutex
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
			status:    WorkerAlive,
			listNodes: InitList(),
			ctx:       s.ctx,
		}
		go s.initWorker(w.workers[i], i)
	}
	return w
}

func (s *Service) CloseWorkers() {
	s.ChangeWorkerStatus(WorkerClosed)
	for _, v := range s.workers.workers {
		v.wg.Wait()
	}
}

func (s *Service) ChangeWorkerStatus(status int) {
	s.workers.mu.Lock()
	defer s.workers.mu.Unlock()
	for _, v := range s.workers.workers {
		if v.status == WorkerClosed {
			continue
		}
		v.status = status
	}
}
