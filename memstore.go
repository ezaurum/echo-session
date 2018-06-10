package session

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/patrickmn/go-cache"
	"time"
)

var _ Store = memorySessionStore{}

const (
	DefaultDuration        = 30 * time.Minute
	DefaultCleanupInterval = 15 * time.Minute
)

func DefaultStore() Store {
	k := snowflake.New(0)
	return NewStore(k, DefaultDuration, DefaultCleanupInterval)
}

func NewStore(k generators.IDGenerator, duration time.Duration,
	cleanupInterval time.Duration) Store {
	return memorySessionStore{
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

func (sm memorySessionStore) GetNew() Session {
	id := sm.keyGenerator.Generate()
	s := New(id, sm)
	sm.Set(s)
	return s
}

func (sm memorySessionStore) Get(id string) (Session, bool) {
	s, e := sm.cache.Get(id)
	if !e {
		return nil, false
	}
	return s.(Session), true
}

func (sm memorySessionStore) Set(s Session) {
	sm.cache.Set(s.ID(), s, sm.duration)
	s.ExpiresAt(time.Now().Add(sm.duration).UnixNano())
}

func (sm memorySessionStore) Count() int {
	return sm.cache.ItemCount()
}

func (sm memorySessionStore) Sessions() SessionMap {
	m := make(SessionMap)
	for k, v := range sm.cache.Items() {
		m[k] = v.Object.(Session)
	}
	return m
}
