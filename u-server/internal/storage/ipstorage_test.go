package storage

import (
	"context"
	"sync"
	"testing"

	"github.com/pullya/unique_server/u-server/internal/config"
)

func TestIpStorage_IsNewIp(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mu := sync.Mutex{}
	type fields struct {
		uniqueIps map[uint32]map[string]int
		shardFunc func(in string) uint32
		mu        *sync.Mutex
	}
	type args struct {
		ctx context.Context
		ip  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Brand new ip",
			fields: fields{
				uniqueIps: map[uint32]map[string]int{uint32(2): {"127.0.0.1": 1}, uint32(0): {"192.168.1.1": 2}},
				shardFunc: ShardKey,
				mu:        &mu,
			},
			args: args{
				ctx: ctx,
				ip:  "100.200.300.1",
			},
			want: true,
		},
		{
			name: "Not a new ip, but connection approved #1",
			fields: fields{
				uniqueIps: map[uint32]map[string]int{uint32(2): {"127.0.0.1": 1}, uint32(0): {"192.168.1.1": config.Config.MaxIpConn - 1}},
				shardFunc: ShardKey,
				mu:        &mu,
			},
			args: args{
				ctx: ctx,
				ip:  "192.168.1.1",
			},
			want: true,
		},
		{
			name: "Not a new ip, connection = MaxIpConn #2",
			fields: fields{
				uniqueIps: map[uint32]map[string]int{uint32(2): {"127.0.0.1": 1}, uint32(0): {"192.168.1.1": config.Config.MaxIpConn}},
				shardFunc: ShardKey,
				mu:        &mu,
			},
			args: args{
				ctx: ctx,
				ip:  "192.168.1.1",
			},
			want: false,
		},
		{
			name: "Not a new ip, connection refused",
			fields: fields{
				uniqueIps: map[uint32]map[string]int{uint32(2): {"127.0.0.1": 1}, uint32(0): {"192.168.1.1": config.Config.MaxIpConn + 1}},
				shardFunc: ShardKey,
				mu:        &mu,
			},
			args: args{
				ctx: ctx,
				ip:  "192.168.1.1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := &IpStorage{
				uniqueIps: tt.fields.uniqueIps,
				shardFunc: ShardKey,
				mu:        tt.fields.mu,
			}

			if got := ir.IsNewIp(tt.args.ctx, tt.args.ip); got != tt.want {
				t.Errorf("IpStorage.IsNewIp() = %v, want %v", got, tt.want)
			}
		})
	}
}
