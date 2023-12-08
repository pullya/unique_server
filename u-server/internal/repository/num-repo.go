package repository

import (
	"context"
	"math/big"
	"sync"
)

type NumRepo struct {
	UniqueNums map[string]struct{}
	mu         *sync.RWMutex
}

func NewNumRepo() NumRepo {
	uniqueNums := make(map[string]struct{}, 0)
	return NumRepo{
		UniqueNums: uniqueNums,
		mu:         &sync.RWMutex{},
	}
}

type INumRepo interface {
	GenUniqueNum(ctx context.Context) big.Int
}

// TO_DO Доделать этот метод!!!
func (nr *NumRepo) GenUniqueNum(ctx context.Context) big.Int {
	numStr := ""

	nr.mu.RLock()
	defer nr.mu.Unlock()

	if _, ok := nr.UniqueNums[numStr]; ok {
		return *big.NewInt(1)
	}
	return *big.NewInt(1)
}
