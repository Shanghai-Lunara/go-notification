package common

type Player struct {
	Pid   int
	Value int
	C     int
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
	t := &ListNode{Player: &Player{Pid: 0, Value: 0}}
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
	l.listRightPush(l.Players[0], m)
}

func (l *ListNodes) listRightPush(t *ListNode, p *ListNode) {
	if t.Player.Value <= p.Player.Value {
		if t.RLink == nil {
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
