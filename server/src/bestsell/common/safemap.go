package common

import (
	"sync"
)

var SHARD_COUNT = 8
type SafeMapShard struct {
	items map[int]interface{}
	sync.RWMutex
}
type SafeMap []*SafeMapShard

func NewSafeMap() *SafeMap {
	m := make(SafeMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &SafeMapShard{items: make(map[int]interface{})}
	}
	return &m
}

func makeHash(key int) int {
	hash := 2166136261
	const prime32 = 16777619
	for key > 0 {
		hash *= prime32
		hash ^= key%10
		key = key/10-1
	}
	idx := hash%SHARD_COUNT
	if idx < 0 {
		idx = -idx
	}
	return idx
}

func (p *SafeMap) GetShard(key int) *SafeMapShard {
	idx := makeHash(key)
	return (*p)[idx]
}

func (p *SafeMap) Set(key int, val interface{}) {
	shard := p.GetShard(key)
	shard.Lock()
	shard.items[key] = val
	shard.Unlock()
}

func (p *SafeMap) Get(key int) interface{} {
	shard := p.GetShard(key)
	shard.RLock()
	val := shard.items[key]
	shard.RUnlock()
	return val
}

func (p *SafeMap) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := (*p)[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (p *SafeMap) Has(key int) bool {
	shard := p.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

func (p *SafeMap) Remove(key int) {
	shard := p.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

func (p *SafeMap) RangeSafe(fn func(key int, v interface{}) bool) {
	var stop bool
	for idx := range *p {
		shard := (*p)[idx]
		shard.RLock()
		for key, value := range shard.items {
			if !fn(key, value) {
				stop = true
				break
			}
		}
		shard.RUnlock()
		if stop {
			return
		}
	}
}

///////////////////////////////
type SafeMapShardS struct {
	items map[string]interface{}
	sync.RWMutex
}
type SafeMapS []*SafeMapShardS

func NewSafeMapS() *SafeMapS {
	m := make(SafeMapS, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &SafeMapShardS{items: make(map[string]interface{})}
	}
	return &m
}

func fnv32(key string) int {
	hash := 2166136261
	const prime32 = 16777619
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= int(key[i])
	}
	if hash<0 {
		hash = -hash
	}
	return hash
}

func (p *SafeMapS) GetShard(key string) *SafeMapShardS {
	idx := fnv32(key)%SHARD_COUNT
	return (*p)[idx]
}

func (p *SafeMapS) Set(key string, val interface{}) {
	shard := p.GetShard(key)
	shard.Lock()
	shard.items[key] = val
	shard.Unlock()
}

func (p *SafeMapS) Get(key string) interface{} {
	shard := p.GetShard(key)
	shard.RLock()
	val := shard.items[key]
	shard.RUnlock()
	return val
}

func (p *SafeMapS) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := (*p)[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (p *SafeMapS) Has(key string) bool {
	shard := p.GetShard(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

func (p *SafeMapS) Remove(key string) {
	shard := p.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

func (p *SafeMapS) RangeSafe(fn func(key string, v interface{}) bool) {
	var stop bool
	for idx := range *p {
		shard := (*p)[idx]
		shard.RLock()
		for key, value := range shard.items {
			if !fn(key, value) {
				stop = true
				break
			}
		}
		shard.RUnlock()
		if stop {
			return
		}
	}
}
