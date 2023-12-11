package repository

import (
	"context"
	"math/big"
	"sync"
	"testing"
)

func TestNumRepo_GenUniqueNum(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mu := sync.Mutex{}
	testMap := make(map[string]struct{})
	testMap["100500"] = struct{}{}
	testMap["123456789"] = struct{}{}

	controlMap := make(map[string]struct{}, len(testMap))
	for key, value := range testMap {
		controlMap[key] = value
	}

	type fields struct {
		UniqueNums map[string]struct{}
		mu         *sync.Mutex
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   big.Int
	}{
		{
			name: "",
			fields: fields{
				UniqueNums: testMap,
				mu:         &mu,
			},
			args: args{ctx: ctx},
			want: *big.NewInt(100500),
		},
	}
	for i := 0; i < 100000; i++ {
		t.Run(tests[0].name, func(t *testing.T) {
			nr := &NumRepo{
				UniqueNums: tests[0].fields.UniqueNums,
				mu:         tests[0].fields.mu,
			}
			got := nr.GenUniqueNum(tests[0].args.ctx)
			if _, ok := controlMap[got.String()]; ok {
				t.Errorf("NumRepo.GenUniqueNum() returned not unique value: %v", got)
			}
			controlMap[got.String()] = struct{}{}
		})
	}
}
