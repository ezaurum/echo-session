package memstore

import (
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/session"
	"github.com/patrickmn/go-cache"
	"time"
)

var _ session.Store = memorySessionStore{}

const (
	DefaultDuration        = 30 * time.Minute
	DefaultCleanupInterval = 15 * time.Minute
)

func Default() session.Store {
	k := snowflake.New(0)
	return New(k, DefaultDuration, DefaultCleanupInterval)
}

func New(k generators.IDGenerator, duration time.Duration,
	cleanupInterval time.Duration) session.Store {
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

func (sm memorySessionStore) GetNew() session.Session {
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
	sm.cache.Set(s.ID(), s, sm.duration)
	s.ExpiresAt(time.Now().Add(sm.duration).UnixNano())
}

func (sm memorySessionStore) Count() int {
	return sm.cache.ItemCount()
}

func (sm memorySessionStore) Sessions() session.SessionMap {
	m := make(session.SessionMap)
	for k, v := range sm.cache.Items() {
		m[k] = v.Object.(session.Session)
	}
	return m
}
