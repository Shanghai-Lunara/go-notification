package common

import (
	"context"
	"go-notification/dao"
	"push/config"
	"reflect"
	"testing"
)

func TestService_Combine(t *testing.T) {
	type fields struct {
		c         *config.Config
		dao       *dao.Dao
		rpcClient *RpcClient
		workers   *Workers
		ctx       context.Context
		cancel    context.CancelFunc
	}
	type args struct {
		info []string
		pid  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []string
		wantMin int
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
			wantRes: []string{
				"1001:3:100006:1:0",
				"1001:3:100006:1:1",
				"1001:1:100011:0:0",
			},
			wantMin: 100003,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				//c:         tt.fields.c,
				dao:       tt.fields.dao,
				rpcClient: tt.fields.rpcClient,
				workers:   tt.fields.workers,
				ctx:       tt.fields.ctx,
				cancel:    tt.fields.cancel,
			}
			gotRes, gotMin := s.Combine(tt.args.info, tt.args.pid)
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Service.Combine() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
			if gotMin != tt.wantMin {
				t.Errorf("Service.Combine() gotMin = %v, want %v", gotMin, tt.wantMin)
			}
		})
	}
}
