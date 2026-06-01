package utils

type StateCache struct {
	m map[int64]string
}

func NewStateCache() *StateCache {
	return &StateCache{m: make(map[int64]string)}
}

func (sc *StateCache) Set(userID int64, state string) { sc.m[userID] = state }
func (sc *StateCache) Get(userID int64) string        { return sc.m[userID] }
func (sc *StateCache) Clear(userID int64)              { delete(sc.m, userID) }
