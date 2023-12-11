package repository

import (
	"context"
	"sync"
	"testing"

	"github.com/pullya/unique_server/u-server/internal/config"
)

func TestIpRepo_IsNewIp(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mu := sync.Mutex{}
	type fields struct {
		UniqueIps map[string]int
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
				UniqueIps: map[string]int{"127.0.0.1": 1, "192.168.1.1": 2},
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
				UniqueIps: map[string]int{"127.0.0.1": 1, "192.168.1.1": config.MaxIpConnection - 1},
				mu:        &mu,
			},
			args: args{
				ctx: ctx,
				ip:  "192.168.1.1",
			},
			want: true,
		},
		{
			name: "Not a new ip, but connection approved #2",
			fields: fields{
				UniqueIps: map[string]int{"127.0.0.1": 1, "192.168.1.1": config.MaxIpConnection},
				mu:        &mu,
			},
			args: args{
				ctx: ctx,
				ip:  "192.168.1.1",
			},
			want: true,
		},
		{
			name: "Not a new ip, connection refused",
			fields: fields{
				UniqueIps: map[string]int{"127.0.0.1": 1, "192.168.1.1": config.MaxIpConnection + 1},
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
			ir := &IpRepo{
				UniqueIps: tt.fields.UniqueIps,
				mu:        tt.fields.mu,
			}
			if got := ir.IsNewIp(tt.args.ctx, tt.args.ip); got != tt.want {
				t.Errorf("IpRepo.IsNewIp() = %v, want %v", got, tt.want)
			}
		})
	}
}
