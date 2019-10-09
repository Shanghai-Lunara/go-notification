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

func (w *Worker) Combine(info []string, pid int) (del []string, meet map[int]int, min int) {
	//fmt.Sprintf("%d:%d:%d:%d:%d", playerId, type, end_time, subid, status)
	var (
		t = make(map[int]int, 0)
		n = make(map[int]map[int]int, 0)
	)
	del = make([]string, 0)
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
				del = append(del, fmt.Sprintf("%d:%d:%d:%d:%d", pid, infoType, tmp2, infoSub, StatusAlive))
			}
		}
		if status == StatusClosed {
			del = append(del, fmt.Sprintf("%d:%d:%d:%d:%d", pid, infoType, infoTime, infoSub, StatusClosed))
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
