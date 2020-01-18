package memstore

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/ezaurum/cthulthu/generators"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/session"
	"github.com/labstack/gommon/random"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
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
		cache:           cache.New(duration, cleanupInterval),
		random:          random.New(),
	}
}

type memorySessionStore struct {
	cache           *cache.Cache
	random          *random.Random
	duration        time.Duration
	cleanupInterval time.Duration
}

func (sm memorySessionStore) GetNew(args ...string) session.Session {
	var ipAddress string
	var agent string
	if len(args) > 1 {
		ipAddress = args[0]
		agent = args[1]
	} else {
		ipAddress = sm.random.String(10, random.Alphanumeric)
		agent = sm.random.String(10, random.Alphanumeric)
	}
	hash := sha256.New()
	hashTarget := strings.Join(args, " ")
	hashTarget += sm.random.String(10, random.Alphanumeric)
	hash.Write([]byte(hashTarget))
	sum := hash.Sum(nil)
	key := base64.RawURLEncoding.EncodeToString(sum)
	s := session.New(key, sm, ipAddress, agent)
	s.Extend(sm.duration)
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
