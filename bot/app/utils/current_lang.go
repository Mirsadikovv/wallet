package utils

import "sync"

type LanguageCache struct {
	mu sync.RWMutex
	m  map[int64]string
}

func NewLanguageCache() *LanguageCache {
	return &LanguageCache{m: make(map[int64]string)}
}

func (lc *LanguageCache) Set(userID int64, language string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.m[userID] = language
}

func (lc *LanguageCache) Get(userID int64) string {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.m[userID]
}
