package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	StatusAlive = iota
	StatusClosed
)

func (w *Worker) Combine(info []string, pid int, clear bool) (del map[string]int, meet map[int]int, min int) {
	//fmt.Sprintf("%d:%d:%d:%d:%d", playerId, type, end_time, subid, status)
	var (
		t = make(map[int]int, 0)
		n = make(map[int]map[int]int, 0)
	)
	del = make(map[string]int, 0)
	meet = make(map[int]int, 0)
	for _, v := range info {
		a := strings.Split(v, ":")
		infoType, err := strconv.Atoi(a[1])
		if err != nil {
			continue
		}
		infoTime, err := strconv.Atoi(a[2])
		if err != nil {
			continue
		}
		infoSub, err := strconv.Atoi(a[3])
		if err != nil {
			continue
		}
		status, err := strconv.Atoi(a[4])
		if err != nil {
			continue
		}
		if tmp1, ok := n[infoType]; ok {
			if tmp2, ok := tmp1[infoSub]; ok {
				str := fmt.Sprintf("%d:%d:%d:%d:%d", pid, infoType, tmp2, infoSub, StatusAlive)
				if _, ok := del[str]; ok {
					del[str]++
				} else {
					del[str] = 1
				}
			}
		}
		if status == StatusClosed {
			str := fmt.Sprintf("%d:%d:%d:%d:%d", pid, infoType, infoTime, infoSub, StatusClosed)
			if _, ok := del[str]; ok {
				del[str]++
			} else {
				del[str] = 1
			}
			delete(n[infoType], infoSub)
		} else {
			if _, ok := n[infoType]; !ok {
				n[infoType] = make(map[int]int, 0)
			}
			n[infoType][infoSub] = infoTime
		}
	}
	for infoType, v := range n {
		for _, endTime := range v {
			if m, ok := t[infoType]; ok {
				if endTime > m {
					t[infoType] = endTime
				}
			} else {
				t[infoType] = endTime
			}
		}
	}
	now := int(time.Now().Unix())
	if clear == true {
		for infoType, v := range n {
			for sub, endTime := range v {
				if endTime <= now {
					str := fmt.Sprintf("%d:%d:%d:%d:%d", pid, infoType, endTime, sub, StatusAlive)
					if _, ok := del[str]; ok {
						del[str]++
					} else {
						del[str] = 1
					}
				}
			}
		}
	}
	for k, v := range t {
		if min == 0 {
			min = v
		}
		if min > v {
			min = v
		}
		if v <= now {
			meet[k] = v
		}
	}
	return del, meet, min
}

const (
	DelayDefault = iota
	DelayAlive
)

func (w *Worker) TransferMinTime(min, delay int) int {
	if delay != DelayAlive {
		return min
	}
	t := int64(min)
	tm := time.Unix(t, 0)
	l := time.Now().Location()
	t22 := time.Date(tm.Year(), tm.Month(), tm.Day(), 22, 0, 0, 0, l)
	r1 := t22.Unix()
	if t >= r1 {
		return int(r1) + 10*3600
	}
	t8 := time.Date(tm.Year(), tm.Month(), tm.Day(), 8, 0, 0, 0, l)
	r2 := t8.Unix()
	if t <= r2 {
		return int(r2)
	}
	return min
}
