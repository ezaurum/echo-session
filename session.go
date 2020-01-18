package session

import "time"

type Store interface {
	GetNew(args ...string) Session
	Get(id string) (Session, bool)
	Set(session Session)
	Delete(key string)
	Sessions() StoreMap
	Count() int
}

type Session interface {
	Key() string
	Get(key string) (interface{}, bool)
	MustGet(key string) interface{}
	Set(key string, o interface{})
	Save()
	IsExpired() bool
	ExpiresAt(nano int64)
	Store() Store
	Remove(key string)
	Extend(duration time.Duration)
	Destroy()
	IPAddress() string
	Agent() string
}

type StoreMap map[string]Session
