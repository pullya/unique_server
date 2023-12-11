package repository

import (
	"context"
	"math/big"
	"math/rand"
	"sync"
)

type NumRepo struct {
	UniqueNums map[string]struct{}
	mu         *sync.Mutex
}

func NewNumRepo() NumRepo {
	uniqueNums := make(map[string]struct{}, 0)
	return NumRepo{
		UniqueNums: uniqueNums,
		mu:         &sync.Mutex{},
	}
}

type INumRepo interface {
	GenUniqueNum(ctx context.Context) big.Int
}

func (nr *NumRepo) GenUniqueNum(ctx context.Context) big.Int {
	newInt := rand.Int63()
	newBigInt := big.NewInt(newInt)

	nr.mu.Lock()
	defer nr.mu.Unlock()

	for {
		if _, ok := nr.UniqueNums[newBigInt.String()]; !ok {
			nr.UniqueNums[newBigInt.String()] = struct{}{}
			break
		}
		newInt = rand.Int63()
		newBigInt = big.NewInt(newInt)
	}
	return *newBigInt
}
