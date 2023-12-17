package storage

import (
	"context"
	"math/big"
	"sync"
	"testing"
)

func TestNumStorage_GenUniqueNum(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mu := sync.Mutex{}

	testMap := make(map[uint32]map[string]struct{})
	controlMap := make(map[uint32]map[string]struct{})

	type fields struct {
		uniqueNums map[uint32]map[string]struct{}
		shardFunc  func(in string) uint32
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
			name: "Test #1",
			fields: fields{
				uniqueNums: testMap,
				shardFunc:  ShardKey,
				mu:         &mu,
			},
			args: args{ctx: ctx},
			want: *big.NewInt(100500),
		},
	}
	for i := 0; i < 100000; i++ {
		t.Run(tests[0].name, func(t *testing.T) {
			ns := &NumStorage{
				uniqueNums: tests[0].fields.uniqueNums,
				shardFunc:  ShardKey,
				mu:         tests[0].fields.mu,
			}
			got := ns.GenUniqueNum(tests[0].args.ctx)
			sk := ShardKey(got.String())
			shard := controlMap[sk]
			if _, ok := shard[got.String()]; ok {
				t.Errorf("NumStorage.GenUniqueNum() returned not unique value: %s, sk: %d", got.String(), sk)
			}
			if _, ok := controlMap[sk]; !ok {
				controlMap[sk] = make(map[string]struct{})
			}
			controlMap[sk][got.String()] = struct{}{}
		})
	}
}
