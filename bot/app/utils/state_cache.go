package utils

import "sync"

type StateCache struct {
	mu sync.RWMutex
	m  map[int64]string
}

func NewStateCache() *StateCache {
	return &StateCache{m: make(map[int64]string)}
}

func (sc *StateCache) Set(userID int64, state string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.m[userID] = state
}

func (sc *StateCache) Get(userID int64) string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.m[userID]
}

func (sc *StateCache) Clear(userID int64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.m, userID)
}
