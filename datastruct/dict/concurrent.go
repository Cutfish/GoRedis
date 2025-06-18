package dict

import (
	"math"
	"sync"
	"sync/atomic"
)

type ConcurrentDict struct {
	table []*Shard
	count int32
}

type Shard struct {
	m  map[string]interface{}
	mu sync.RWMutex
}

// 将一个 param 转换为一个大于或等于 param 的最小的 2 的幂次方数。
func computeCapacity(param int) (size int) {
	if param <= 16 {
		return 16
	}
	n := param - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	} else {
		return int(n + 1)
	}
}

func NewConcurrent(shardCount int) *ConcurrentDict {
	shardCount = computeCapacity(shardCount)
	table := make([]*Shard, shardCount)
	for i := 0; i < shardCount; i++ {
		table[i] = &Shard{
			m: make(map[string]interface{}),
		}
	}
	d := &ConcurrentDict{
		count: 0,
		table: table,
	}
	return d
}

// 定位shard, 当n为2的整数幂时 (n - 1) & h 就相当于 h % n
func (dict *ConcurrentDict) spread(hashCode uint32) uint32 {
	if dict == nil {
		panic("dict is nil")
	}
	tableSize := uint32(len(dict.table))
	return (tableSize - 1) & uint32(hashCode)
}

func (dict *ConcurrentDict) getShard(index uint32) *Shard {
	if dict == nil {
		panic("dict is nil")
	}
	return dict.table[index]
}

// Get方法
func (dict *ConcurrentDict) Get(key string) (val interface{}, exists bool) {
	if dict == nil {
		panic("dict is nil")
	}

	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	shard := dict.getShard(index)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	val, exists = shard.m[key]
	return
}

// 键值对的数量
func (dict *ConcurrentDict) Len() int {
	if dict == nil {
		panic("dict is nil")
	}
	return int(atomic.LoadInt32(&dict.count))
}

// Put方法，之前若存在这个键就返回 0 否则返回 1
func (dict *ConcurrentDict) Put(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	shard := dict.getShard(index)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	if _, ok := shard.m[key]; ok {
		shard.m[key] = val
		return 0
	} else {
		dict.addCount()
		shard.m[key] = val
		return 1
	}
}

// // Remove removes the key and return the number of deleted key-value
func (dict *ConcurrentDict) Remove(key string) (val interface{}, result int) {
	if dict == nil {
		panic("dict is nil")
	}
	hashCode := fnv32(key)
	index := dict.spread(hashCode)
	shard := dict.getShard(index)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if val, ok := shard.m[key]; ok {
		delete(shard.m, key)
		dict.decreaseCount()
		return val, 1
	}
	return nil, 0
}

func (dict *ConcurrentDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}

func (dict *ConcurrentDict) decreaseCount() int32 {
	return atomic.AddInt32(&dict.count, -1)
}

// 哈希算法采用FNV算法
const prime32 = uint32(16777619)

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32 // 可能会溢出，保留低32位
		hash ^= uint32(key[i])
	}
	return hash
}
