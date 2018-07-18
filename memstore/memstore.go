package memstore

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/patrickmn/go-cache"
	"time"
	"github.com/ezaurum/echo-session"
)

var _ session.Store = &memorySessionStore{}

const (
	DefaultDuration        = 30 * time.Minute
	DefaultCleanupInterval = 15 * time.Minute
)

func DefaultStore() session.Store {
	k := snowflake.New(0)
	return NewStore(k, DefaultDuration, DefaultCleanupInterval)
}

func NewStore(k generators.IDGenerator, duration time.Duration,
	cleanupInterval time.Duration) session.Store {
	return &memorySessionStore{
		duration:        duration,
		cleanupInterval: cleanupInterval,
		keyGenerator:    k,
		cache:           cache.New(duration, cleanupInterval),
	}
}

type memorySessionStore struct {
	cache           *cache.Cache
	keyGenerator    generators.IDGenerator
	duration        time.Duration
	cleanupInterval time.Duration
}

func (sm memorySessionStore) GetNew(args...string) session.Session {
	id := sm.keyGenerator.Generate()
	s := session.New(id, sm)
	sm.Set(s)
	return s
}

func (sm memorySessionStore) Get(id string) (session.Session, bool) {
	s, e := sm.cache.Get(id)
	if !e {
		return nil, false
	}
	return s.(session.Session), true
}

func (sm memorySessionStore) Set(s session.Session) {
	s.Extend(sm.duration)
	sm.cache.Set(s.Key(), s, sm.duration)
}

func (sm memorySessionStore) Delete(key string) {
	sm.cache.Delete(key)
}

func (sm memorySessionStore) Count() int {
	return sm.cache.ItemCount()
}

func (sm memorySessionStore) Sessions() session.StoreMap {
	m := make(session.StoreMap)
	for k, v := range sm.cache.Items() {
		m[k] = v.Object.(session.Session)
	}
	return m
}
