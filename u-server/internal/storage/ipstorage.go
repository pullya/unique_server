package storage

import (
	"context"
	"sync"

	"github.com/pullya/unique_server/u-server/internal/config"
)

type IpStorage struct {
	uniqueIps map[uint32]map[string]int
	shardFunc func(data string) uint32
	mu        *sync.Mutex
}

func NewIpStorage() IpStorage {
	uniqueIps := make(map[uint32]map[string]int, 0)
	return IpStorage{
		uniqueIps: uniqueIps,
		shardFunc: ShardKey,
		mu:        &sync.Mutex{},
	}
}

type IpStorageer interface {
	IsNewIp(ctx context.Context, ip string) bool
}

func (is *IpStorage) IsNewIp(ctx context.Context, ip string) bool {
	is.mu.Lock()
	defer is.mu.Unlock()

	shardKey := is.shardFunc(ip)
	if _, ok := is.uniqueIps[shardKey]; !ok {
		is.uniqueIps[shardKey] = make(map[string]int)
		is.uniqueIps[shardKey][ip]++
		return true
	}

	cnt, ok := is.uniqueIps[shardKey][ip]
	if ok && cnt >= config.Config.MaxIpConn {
		return false
	}

	is.uniqueIps[shardKey][ip]++
	return true
}
