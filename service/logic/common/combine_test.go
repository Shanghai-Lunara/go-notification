package common

import (
	"context"
	"go-notification/dao"
	"reflect"
	"testing"
)

func TestWorker_Combine(t *testing.T) {
	type fields struct {
		dao       *dao.Dao
		addr      string
		count     int
		status    int
		listNodes *ListNodes
		ctx       context.Context
	}
	type args struct {
		info []string
		pid  int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantDel  []string
		wantMeet map[int]int
		wantMin  int
	}{
		{
			name:   "case1",
			fields: fields{},
			args: args{
				info: []string{
					"1001:1:100011:0:0",
					"1001:2:100003:0:0",
					"1001:3:100006:1:0",
					"1001:3:100020:2:0",
					"1001:3:100006:1:1",
					"1001:3:100030:3:0",
					"1001:4:100016:0:0",
					"1001:5:100032:0:0",
					"1001:6:100009:0:0",
					"1001:1:100014:0:0",
				},
				pid: 1001,
			},
			wantDel: []string{
				"1001:3:100006:1:0",
				"1001:3:100006:1:1",
				"1001:1:100011:0:0",
			},
			wantMeet: map[int]int{
				1: 100014,
				2: 100003,
				3: 100030,
				4: 100016,
				5: 100032,
				6: 100009,
			},
			wantMin: 100003,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				dao:       tt.fields.dao,
				addr:      tt.fields.addr,
				count:     tt.fields.count,
				status:    tt.fields.status,
				listNodes: tt.fields.listNodes,
				ctx:       tt.fields.ctx,
			}
			gotDel, gotMeet, gotMin := w.Combine(tt.args.info, tt.args.pid)
			if !reflect.DeepEqual(gotDel, tt.wantDel) {
				t.Errorf("Worker.Combine() gotDel = %v, want %v", gotDel, tt.wantDel)
			}
			if !reflect.DeepEqual(gotMeet, tt.wantMeet) {
				t.Errorf("Worker.Combine() gotMeet = %v, want %v", gotMeet, tt.wantMeet)
			}
			if gotMin != tt.wantMin {
				t.Errorf("Worker.Combine() gotMin = %v, want %v", gotMin, tt.wantMin)
			}
		})
	}
}
