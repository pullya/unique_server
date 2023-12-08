package repository

import (
	"context"
	"errors"
	"sync"
)

type IpRepo struct {
	UniqueIps map[string]struct{}
	mu        *sync.RWMutex
}

func NewIpRepo() IpRepo {
	uniqueIps := make(map[string]struct{}, 0)
	return IpRepo{
		UniqueIps: uniqueIps,
		mu:        &sync.RWMutex{},
	}
}

type IIpRepo interface {
	IsNewIp(ctx context.Context, ip string) bool
	AddIp(ctx context.Context, ip string) error
}

func (ir *IpRepo) IsNewIp(ctx context.Context, ip string) bool {
	ir.mu.RLock()
	defer ir.mu.RUnlock()

	if _, ok := ir.UniqueIps[ip]; ok {
		return false
	}
	return true
}

func (ir *IpRepo) AddIp(ctx context.Context, ip string) error {
	ir.mu.Lock()
	defer ir.mu.Unlock()

	if _, ok := ir.UniqueIps[ip]; ok {
		return errors.New("ip is not unique")
	}

	ir.UniqueIps[ip] = struct{}{}
	return nil
}
