package common

import "sync"

type Player struct {
	mu    sync.RWMutex
	Pid   int
	Value int
	Delay int
}

type ListNode struct {
	LLink  *ListNode
	RLink  *ListNode
	Player *Player
}

type ListNodes struct {
	Players map[int]*ListNode
}

func InitList() *ListNodes {
	t := &ListNode{Player: &Player{Pid: 0, Value: 0, Delay: DelayDefault}}
	l := &ListNodes{
		Players: make(map[int]*ListNode, 0),
	}
	l.Players[0] = t
	return l
}

func (l *ListNodes) AppendOrModify(p *Player) {
	var m *ListNode
	if t, ok := l.Players[p.Pid]; ok {
		m = t
	} else {
		m = &ListNode{Player: p}
		l.Players[p.Pid] = m
	}
	if p.Value == 0 {
		if l.Players[p.Pid].LLink != nil && l.Players[p.Pid].RLink != nil {
			l.Players[p.Pid].LLink.RLink, l.Players[p.Pid].RLink.LLink = l.Players[p.Pid].RLink, l.Players[p.Pid].LLink
		}
		if l.Players[p.Pid].LLink != nil && l.Players[p.Pid].RLink == nil {
			l.Players[p.Pid].LLink.RLink = nil
		}
		if l.Players[p.Pid].LLink == nil && l.Players[p.Pid].RLink != nil {
			l.Players[p.Pid].RLink.LLink = nil
		}
		l.Players[p.Pid].RLink, l.Players[p.Pid].LLink = nil, nil
	} else {
		l.listRightPush(l.Players[0], m)
	}
}

func (l *ListNodes) listRightPush(t *ListNode, p *ListNode) {
	if t.Player.Value <= p.Player.Value {
		if t.RLink == nil {
			if t.Player.Pid == p.Player.Pid {
				return
			}
			if p.LLink != nil {
				p.LLink.RLink = p.RLink
			}
			if p.RLink != nil {
				p.RLink.LLink = p.LLink
			}
			t.RLink, p.LLink, p.RLink = p, t, nil
		} else {
			l.listRightPush(t.RLink, p)
		}
	} else {
		if p.LLink != nil {
			p.LLink.RLink = p.RLink
		}
		if p.RLink != nil {
			p.RLink.LLink = p.LLink
		}
		t.LLink.RLink, p.LLink = p, t.LLink
		t.LLink, p.RLink = p, t
	}
}
