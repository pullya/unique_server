package storage

import (
	"github.com/pullya/unique_server/u-server/internal/config"
	"github.com/spaolacci/murmur3"
)

func ShardKey(in string) uint32 {
	key := []byte(in)
	numShards := config.Config.ShardsCnt
	if numShards == 0 {
		numShards = 8
	}

	hashValue := HashShard(key)

	return hashValue % uint32(numShards)
}

func HashShard(data []byte) uint32 {
	hasher := murmur3.New32()
	hasher.Write(data)
	return hasher.Sum32()
}
