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
	w.mu.Lock()
	defer w.mu.Unlock()
	if p, ok := w.listNodes.Players[pid]; ok {
		if min != 0 {
			min = w.TransferMinTime(min, p.Player.Delay)
		}
		if p.Player.Value != min {
			p.Player.Value = min
			w.listNodes.AppendOrModify(p.Player)
		}
	} else {
		delay := DelayDefault
		if _, t, err := w.dao.GetPlayerSettings(pid); err != nil {
			log.Println("UpdatePlayerListNodes GetPlayerSettings err:", err)
		} else {
			delay = t
		}
		if min != 0 {
			min = w.TransferMinTime(min, delay)
		}
		p := &Player{
			Pid:   pid,
			Value: min,
			Delay: delay,
		}
		w.listNodes.AppendOrModify(p)
	}
}

func (w *Worker) UpdatePlayerSettings(pid int) (err error) {
	if _, delay, err := w.dao.GetPlayerSettings(pid); err != nil {
		return err
	} else {
		if p, ok := w.listNodes.Players[pid]; ok {
			p.Player.Delay = delay
		} else {
			if err = w.RefreshOne(strconv.Itoa(pid)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Worker) RefreshOne(str string) (err error) {
	var (
		pid int
	)
	tmp := strings.Split(str, ",")
	if pid, err = strconv.Atoi(tmp[0]); err != nil {
		return err
	}
	if len(tmp) > 1 {
		return w.UpdatePlayerSettings(pid)
	}
	if _, min, err := w.PullPlayerOne(pid, false); err != nil {
		return err
	} else {
		w.UpdatePlayerValue(pid, min)
		return nil
	}
}

func (w *Worker) CheckOne(pid int) (err error) {
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
