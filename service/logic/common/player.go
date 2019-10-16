package common

import (
	"log"
	"strconv"
	"strings"
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
		log.Printf("UpdatePlayerValue TransferMinTime min:%d", min)
	}
	log.Printf("UpdatePlayerValue 111111 min:%d", min)
	if p.Value != min {
		log.Printf("UpdatePlayerValue 2222222 min:%d", min)
		p.Value = min
		w.listNodes.AppendOrModify(p)
	}
}

func (w *Worker) UpdatePlayerSettings(p *Player) (update bool, err error) {
	if _, delay, err := w.dao.GetPlayerSettings(p.Pid); err != nil {
		return true, err
	} else {
		if p.Value == 0 {
			log.Println(111111111)
			p.Delay = delay
		} else {
			log.Println(2222222)
			if p.Delay != delay {
				log.Println(333333)
				if p.Delay == DelayDefault {
					log.Println(44444)
					p.Delay = delay
					p.Value = w.TransferMinTime(p.Value, delay)
					return false, nil
				} else {
					log.Println(5555555)
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
	if len(tmp) > 1 {
		log.Println("RefreshOne UpdatePlayerSettings 111111")
		if up, err = w.UpdatePlayerSettings(p); err != nil {
			return err
		}
		if up == false {
			log.Println("RefreshOne UpdatePlayerValue 222222")
			w.UpdatePlayerValue(p, p.Value)
			return nil
		}
	}
	log.Println("RefreshOne start lock 1111")
	p.mu.Lock()
	defer func() {
		log.Println("RefreshOne close lock")
		p.mu.Unlock()
	}()
	if _, min, err := w.PullPlayerOne(pid, false); err != nil {
		return err
	} else {
		log.Println("RefreshOne UpdatePlayerValue 3333333")
		w.UpdatePlayerValue(p, min)
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
				//if cid == "" {
				//	return nil
				//}
				for k, v := range meet {
					_ = v
					if err = w.ApiPost(pid, k, cid); err != nil {
						return nil
					}
				}
			}
			w.UpdatePlayerValue(p, min)
		}
		return nil
	}
}

func (w *Worker) ApiPost(pid, infoType int, cid string) (err error) {
	log.Printf("api-post pid:%d type:%d cid:%s \n", pid, infoType, cid)
	return nil
}
