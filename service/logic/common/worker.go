package common

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/nevercase/go-notification/api"
	"github.com/nevercase/go-notification/dao"
)

const (
	WorkerAlive = iota
	WorkerLoading
	WorkerInterrupt
	WorkerClosed
)

type Worker struct {
	mu        sync.RWMutex
	wg        sync.WaitGroup
	id        int
	dao       *dao.Dao
	addr      string
	count     int
	status    int
	listNodes *ListNodes
	push      api.Push
	ctx       context.Context
}

func (w *Worker) loadCache() {
	tick := time.NewTicker(time.Second * 1)
	defer func() {
		w.wg.Done()
		tick.Stop()
		w.status = WorkerAlive
	}()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-tick.C:
			if w.status == WorkerClosed {
				return
			}
			var (
				p          []string
				err        error
				start, end int
			)
			length := 1000
			start, end = 0, 0
			for {
				start = end
				end += length
				if p, err = w.dao.ZRevRange(w.addr, start, end); err != nil {
					log.Println("err:", err)
				} else {
					for _, v := range p {
						if w.status == WorkerClosed {
							return
						}
						if err = w.RefreshOne(v, true); err != nil {
							log.Print("loadCache RefreshOne err:", err)
						}
					}
					if len(p) != length {
						log.Println("ZRevRange break")
						return
					}
					time.Sleep(time.Millisecond * 200)
				}
			}
		}
	}
}

func (w *Worker) appendLoop() {
	tick := time.NewTicker(time.Millisecond * 10)
	defer tick.Stop()
	defer w.wg.Done()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-tick.C:
			if w.status == WorkerAlive {
				var (
					err   error
					p     []string
					count int
				)
				if p, err = w.dao.LRange(w.addr, 200); err != nil {
					log.Println("appendLoop LRange err:", err)
					time.Sleep(time.Second * 1)
				} else {
					count = 0
					for _, v := range p {
						if w.status == WorkerInterrupt {
							continue
						}
						if w.status == WorkerClosed {
							if count > 0 {
								if err = w.dao.LTRIM(w.addr, count); err != nil {
									log.Print("appendLoop WorkerClosed LTRIM err:", err)
								}
							}
							return
						}
						if v != "" {
							if err = w.RefreshOne(v, true); err != nil {
								log.Print("appendLoop RefreshOne err:", err)
							}
							count++
						}
					}
					if count > 0 {
						if err = w.dao.LTRIM(w.addr, count); err != nil {
							log.Print("appendLoop LTRIM err:", err)
						}
					}
				}
			}
		}
	}
}

func (w *Worker) logicLoop() {
	tick := time.NewTicker(time.Millisecond * 10)
	tick1 := time.NewTicker(time.Second * 1)
	defer func() {
		tick.Stop()
		tick1.Stop()
		w.wg.Done()
	}()
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-tick.C:
			if w.status == WorkerLoading || w.status == WorkerInterrupt {
				log.Println("logicLoop WorkerLoading orWorkerInterrupt w.id:", w.addr)
				time.Sleep(time.Second * 1)
			} else {
				if t, ok := w.listNodes.Players[0]; ok {
					if t.RLink != nil {
						p := t.RLink.Player
						if p.Value > int(time.Now().Unix()) {
							//log.Printf("listNodes continue 1111 v:%d time:%d \n", p.Value, int(time.Now().Unix()))
							time.Sleep(time.Millisecond * 5)
						} else {
							log.Printf("listNodes match pid:%d v:%d time:%d \n", p.Pid, p.Value, int(time.Now().Unix()))
							if err := w.CheckOne(t.RLink.Player.Pid); err != nil {
								log.Printf("CheckOne pid:%d p:%v err:%v \n", p.Pid, p, err)
							}
						}
					}
				}
			}
		case <-tick1.C:
			log.Printf("logicLoop WorkerId:%d addr:%s listNodes-len: %d \n", w.id, w.addr, len(w.listNodes.Players))
			if tmp, ok := w.listNodes.Players[1]; ok {
				log.Println("logicLoop Players[1]:", tmp.Player)
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
			log.Println("initWorker loop worker_id:", w.id)
			if addr, err := s.getAllocatedNode(); err != nil {
				log.Println("getAllocatedNode err:", err)
			} else {
				log.Println("initWorker addr:", addr)
				if addr != "" {
					w.addr = addr
					w.wg.Add(3)
					go w.loadCache()
					go w.appendLoop()
					go w.logicLoop()
					return
				}
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
			id:        i,
			dao:       s.dao,
			addr:      "",
			count:     0,
			status:    WorkerLoading,
			listNodes: InitList(),
			push:      s.push,
			ctx:       s.ctx,
		}
		go s.initWorker(w.workers[i])
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
		if v.status == WorkerInterrupt && status == WorkerLoading {
			v.status = WorkerAlive
		} else {
			v.status = status
		}
	}
}
