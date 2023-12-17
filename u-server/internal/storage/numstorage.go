package storage

import (
	"context"
	"math/big"
	"math/rand"
	"sync"
)

type NumStorage struct {
	uniqueNums map[uint32]map[string]struct{}
	shardFunc  func(data string) uint32
	mu         *sync.Mutex
}

func NewNumStorage() NumStorage {
	uniqueNums := make(map[uint32]map[string]struct{}, 0)
	return NumStorage{
		uniqueNums: uniqueNums,
		shardFunc:  ShardKey,
		mu:         &sync.Mutex{},
	}
}

type NumStorageer interface {
	GenUniqueNum(ctx context.Context) big.Int
}

func (ns *NumStorage) GenUniqueNum(ctx context.Context) big.Int {
	newInt := rand.Int63()
	newBigInt := big.NewInt(newInt)

	ns.mu.Lock()
	defer ns.mu.Unlock()

	for {
		shardKey := ns.shardFunc(newBigInt.String())
		if _, ok := ns.uniqueNums[shardKey]; !ok {
			ns.uniqueNums[shardKey] = make(map[string]struct{})
			ns.uniqueNums[shardKey][newBigInt.String()] = struct{}{}
			break
		}
		if _, ok := ns.uniqueNums[shardKey][newBigInt.String()]; !ok {
			ns.uniqueNums[shardKey][newBigInt.String()] = struct{}{}
			break
		}
		newInt = rand.Int63()
		newBigInt = big.NewInt(newInt)
	}
	return *newBigInt
}
