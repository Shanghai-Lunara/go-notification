package common

import (
	"fmt"
	"strconv"
	"strings"
)

func (s *Service) Combine(info []string, pid int) (res []string, min int) {
	//fmt.Sprintf("%d:%d:%d:%d:%d", playerId, type, end_time, subid, status)
	var (
		t = make(map[int]int, 0)
		n = make(map[int]map[int]int, 0)
	)
	res = make([]string, 0)
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
				res = append(res, fmt.Sprintf("%d:%d:%d:%d:%d", pid, infoType, tmp2, infoSub, 0))
			}
		}
		if status == 1 {
			res = append(res, fmt.Sprintf("%d:%d:%d:%d:%d", pid, infoType, infoTime, infoSub, 1))
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
	for _, v := range t {
		if min == 0 {
			min = v
		}
		if min > v {
			min = v
		}
	}
	return res, min
}
