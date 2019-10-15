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

func (w *Worker) UpdatePlayerValue(pid, min int) {
	p := w.GetPlayer(pid)
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

func (w *Worker) UpdatePlayerSettings(pid int) (err error) {
	if _, delay, err := w.dao.GetPlayerSettings(pid); err != nil {
		return err
	} else {
		p := w.GetPlayer(pid)
		if p.Value == 0 {
			log.Println(111111111)
			p.Delay = delay
			return nil
		} else {
			log.Println(2222222)
			if p.Delay != delay {
				log.Println(333333)
				if p.Delay == DelayDefault {
					log.Println(44444)
					p.Delay = delay
					p.Value = w.TransferMinTime(p.Value, delay)
				} else {
					log.Println(5555555)
					p.Delay = delay
					if err = w.RefreshOne(strconv.Itoa(pid), false); err != nil {
						log.Println(666666)
						return err
					}
				}
			}
		}
	}
	return nil
}

func (w *Worker) GetPlayer(pid int) (p *Player) {
	if t, ok := w.listNodes.Players[pid]; ok {
		return t.Player
	} else {
		p := &Player{
			Pid:   pid,
			Value: 0,
			Delay: DelayDefault,
		}
		w.listNodes.AppendOrModify(p)
		return p
	}
}

func (w *Worker) RefreshOne(str string, update bool) (err error) {
	var (
		pid int
		p   *Player
	)
	tmp := strings.Split(str, ":")
	if pid, err = strconv.Atoi(tmp[0]); err != nil {
		return err
	}
	if len(tmp) > 1 {
		return w.UpdatePlayerSettings(pid)
	}
	p = w.GetPlayer(pid)
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, min, err := w.PullPlayerOne(pid, false); err != nil {
		return err
	} else {
		if update == true {
			if err = w.UpdatePlayerSettings(pid); err != nil {
				return err
			}
		}
		w.UpdatePlayerValue(pid, min)
		return nil
	}
}

func (w *Worker) CheckOne(pid int) (err error) {
	var (
		p *Player
	)
	p = w.GetPlayer(pid)
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
			w.UpdatePlayerValue(pid, min)
		}
		return nil
	}
}

func (w *Worker) ApiPost(pid, infoType int, cid string) (err error) {
	log.Printf("api-post pid:%d type:%d cid:%s \n", pid, infoType, cid)
	return nil
}
