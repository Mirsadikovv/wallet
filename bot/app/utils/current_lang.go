package utils

type LanguageCache struct {
	m map[int64]string
}

func NewLanguageCache() *LanguageCache {
	return &LanguageCache{m: make(map[int64]string)}
}

func (lc *LanguageCache) Set(userID int64, language string) { lc.m[userID] = language }
func (lc *LanguageCache) Get(userID int64) string           { return lc.m[userID] }
