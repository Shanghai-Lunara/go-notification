package common

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInitList(t *testing.T) {
	t1 := &ListNode{Player: &Player{Pid: 0, Value: 0}}
	l := &ListNodes{
		Players: make(map[int]*ListNode, 0),
	}
	l.Players[0] = t1
	tests := []struct {
		name string
		want *ListNodes
	}{
		{name: "case1", want: l},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListNodes_appendOrModify(t *testing.T) {
	type fields struct {
		Players map[int]*ListNode
	}
	type args struct {
		p *Player
	}
	t1 := &ListNode{Player: &Player{Pid: 0, Value: 0}}
	l := fields{
		Players: make(map[int]*ListNode, 0),
	}
	l.Players[0] = t1
	a1 := &Player{Pid: 1, Value: 1}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "TestListNodes_appendOrModify",
			fields: l,
			args: args{
				p: a1,
			},
		},
		{
			name:   "TestListNodes_appendOrModify",
			fields: l,
			args: args{
				p: &Player{Pid: 2, Value: 2},
			},
		},
		{
			name:   "TestListNodes_appendOrModify",
			fields: l,
			args: args{
				p: &Player{Pid: 3, Value: 3},
			},
		},
		{
			name:   "TestListNodes_appendOrModify",
			fields: l,
			args: args{
				p: &Player{Pid: 4, Value: 4},
			},
		},
		{
			name:   "TestListNodes_appendOrModify",
			fields: l,
			args: args{
				p: &Player{Pid: 12, Value: 4},
			},
		},
		{
			name:   "TestListNodes_appendOrModify",
			fields: l,
			args: args{
				p: &Player{Pid: 5, Value: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &ListNodes{
				Players: tt.fields.Players,
			}
			l.AppendOrModify(tt.args.p)
		})
	}
	m := l.Players[0]
	for {
		fmt.Printf("Pid: %d v: %d \n", m.Player.Pid, m.Player.Value)
		if m.RLink != nil {
			m = m.RLink
		} else {
			break
		}
	}
}

func getNew(a *Player) *Player {
	a.Value = 10
	return a
}

//func TestListNodes_listRightPush(t *testing.T) {
//	type fields struct {
//		Players map[int]*ListNode
//	}
//	type args struct {
//		t *ListNode
//		p *ListNode
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		{},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			l := &ListNodes{
//				Players: tt.fields.Players,
//			}
//			l.listRightPush(tt.args.t, tt.args.p)
//		})
//	}
//}
