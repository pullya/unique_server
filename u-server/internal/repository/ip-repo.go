package repository

import (
	"context"
	"sync"

	"github.com/pullya/unique_server/u-server/internal/config"
)

type IpRepo struct {
	UniqueIps map[string]int
	mu        *sync.Mutex
}

func NewIpRepo() IpRepo {
	uniqueIps := make(map[string]int, 0)
	return IpRepo{
		UniqueIps: uniqueIps,
		mu:        &sync.Mutex{},
	}
}

type IIpRepo interface {
	IsNewIp(ctx context.Context, ip string) bool
}

func (ir *IpRepo) IsNewIp(ctx context.Context, ip string) bool {
	ir.mu.Lock()
	defer ir.mu.Unlock()

	cnt, ok := ir.UniqueIps[ip]
	if ok && cnt > config.MaxIpConnection {
		return false
	}

	ir.UniqueIps[ip]++
	return true
}
