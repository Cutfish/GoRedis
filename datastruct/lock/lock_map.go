package lock

import (
	"GoRedis/lib/utils"
	"sort"
	"sync"
)

type Locks struct {
	table []*sync.RWMutex
}

func NewLocks(tableSize int) *Locks {
	table := make([]*sync.RWMutex, tableSize)
	for i := 0; i < tableSize; i++ {
		table[i] = &sync.RWMutex{}
	}
	return &Locks{
		table: table,
	}
}

func (locks *Locks) spread(hashCode uint32) uint32 {
	tableSize := uint32(len(locks.table))
	return (tableSize - 1) & uint32(hashCode)
}

func (locks *Locks) Lock(key string) {
	index := locks.spread(utils.Fnv32(key))
	mu := locks.table[index]
	mu.Lock()
}

func (locks *Locks) Unlock(key string) {
	index := locks.spread(utils.Fnv32(key))
	mu := locks.table[index]
	mu.Unlock()
}

func (lock *Locks) toLockIndices(keys []string, reverse bool) []uint32 {
	indexMap := make(map[uint32]struct{})
	for _, key := range keys {
		index := lock.spread(utils.Fnv32(key))
		indexMap[index] = struct{}{}
	}
	indices := make([]uint32, 0, len(indexMap))
	for index := range indexMap {
		indices = append(indices, index)
	}
	sort.Slice(indices, func(i, j int) bool {
		if !reverse {
			return indices[i] < indices[j]
		} else {
			return indices[i] > indices[j]
		}
	})
	return indices
}

// 允许 writeKeys 和 readKeys 中存在重复的 key
func (locks *Locks) RWLocks(writeKeys []string, readKeys []string) {
	keys := append(writeKeys, readKeys...)
	indices := locks.toLockIndices(keys, false)
	writeIndices := locks.toLockIndices(writeKeys, false)
	writeIndexSet := make(map[uint32]struct{})
	for _, idx := range writeIndices {
		writeIndexSet[idx] = struct{}{}
	}
	for _, index := range indices {
		_, w := writeIndexSet[index]
		mu := locks.table[index]
		if w {
			mu.Lock()
		} else {
			mu.RLock()
		}
	}
}

func (locks *Locks) RWUnlocks(writeKeys []string, readKeys []string) {
	keys := append(writeKeys, readKeys...)
	indices := locks.toLockIndices(keys, false)
	writeIndices := locks.toLockIndices(writeKeys, false)
	writeIndexSet := make(map[uint32]struct{})
	for _, idx := range writeIndices {
		writeIndexSet[idx] = struct{}{}
	}
	for _, index := range indices {
		_, w := writeIndexSet[index]
		mu := locks.table[index]
		if w {
			mu.Unlock()
		} else {
			mu.RUnlock()
		}
	}
}
