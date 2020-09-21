package common

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Shanghai-Lunara/go-notification/api"
)

func (w *Worker) PullPlayerOne(pid int, clear bool) (meet map[int]int, min int, err error) {
	var (
		info []string
		del  map[string]int
	)
	meet = make(map[int]int, 0)
	if info, err = w.dao.GetSinglePlayerList(pid); err != nil {
		return meet, min, err
	}
	del, meet, min = w.Combine(info, pid, clear)
	if len(del) > 0 {
		if err = w.dao.UpdateSinglePlayerList(pid, del); err != nil {
			return meet, min, err
		}
	}
	return meet, min, nil
}

func (w *Worker) UpdatePlayerValue(p *Player, min int) {
	if min != 0 {
		min = w.TransferMinTime(min, p.Delay)
	}
	if p.Value != min {
		p.Value = min
		w.listNodes.AppendOrModify(p)
	}
}

func (w *Worker) UpdatePlayerSettings(p *Player) (update bool, err error) {
	if _, delay, err := w.dao.GetPlayerSettings(p.Pid); err != nil {
		return true, err
	} else {
		if p.Value == 0 {
			p.Delay = delay
		} else {
			if p.Delay != delay {
				if p.Delay == DelayDefault {
					p.Delay = delay
					p.Value = w.TransferMinTime(p.Value, delay)
					return false, nil
				} else {
					p.Delay = delay
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (w *Worker) GetPlayer(pid int) (p *Player, err error) {
	if t, ok := w.listNodes.Players[pid]; ok {
		return t.Player, nil
	} else {
		p := &Player{
			Pid:   pid,
			Value: 0,
			Delay: DelayDefault,
		}
		if w.status == WorkerAlive {
			if err = w.dao.ZAdd(pid); err != nil {
				return nil, err
			}
		}
		w.listNodes.AppendOrModify(p)
		if _, err = w.UpdatePlayerSettings(p); err != nil {
			return nil, err
		}
		return p, nil
	}
}

func (w *Worker) RefreshOne(str string, update bool) (err error) {
	var (
		pid int
		p   *Player
		up  bool
	)
	tmp := strings.Split(str, ":")
	if pid, err = strconv.Atoi(tmp[0]); err != nil {
		return err
	}
	if p, err = w.GetPlayer(pid); err != nil {
		return err
	}
	log.Println("RefreshOne p:", p)
	if len(tmp) > 1 {
		if up, err = w.UpdatePlayerSettings(p); err != nil {
			return err
		}
		if up == false {
			w.UpdatePlayerValue(p, p.Value)
			return nil
		}
	}
	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
	}()
	if _, min, err := w.PullPlayerOne(pid, false); err != nil {
		return err
	} else {
		w.UpdatePlayerValue(p, min)
		log.Printf("p pid:%d v:%v\n", p.Pid, p)
		return nil
	}
}

func (w *Worker) CheckOne(pid int) (err error) {
	var (
		p *Player
	)
	if p, err = w.GetPlayer(pid); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if meet, min, err := w.PullPlayerOne(pid, true); err != nil {
		return err
	} else {
		if len(meet) > 0 {
			if cid, _, err := w.dao.GetPlayerSettings(pid); err != nil {
				return err
			} else {
				for k, v := range meet {
					_ = v
					if err = w.ApiPost(pid, k, cid); err != nil {
						return err
					}
				}
			}
			w.UpdatePlayerValue(p, min)
		} else {
			log.Printf("CheckOne no match pid:%d min:%d\n", pid, min)
			w.UpdatePlayerValue(p, min)
		}
		return nil
	}
}

func (w *Worker) ApiPost(pid, infoType int, cid string) (err error) {
	log.Printf("api-post pid:%d type:%d cid:%s \n", pid, infoType, cid)
	if cid == "" {
		return nil
	}
	m := api.NewMessage(w.id, pid, infoType, cid)
	if _, err := w.push.Send(m); err != nil {
		return errors.New(fmt.Sprintf("ApiPost push.Send err:%v", err))
	}
	return nil
}
