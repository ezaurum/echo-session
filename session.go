package session

import "time"

type Store interface {
	GetNew() Session
	Get(id string) (Session, bool)
	Set(session Session)
	Sessions() SessionMap
	Count() int
}

type Session interface {
	ID() string
	Get(key string) (interface{}, bool)
	MustGet(key string) interface{}
	Set(key string, o interface{})
	Save()
	IsExpired() bool
	ExpiresAt(nano int64)
	Store() Store
}

type SessionMap map[string]Session

func New(id string, store Store) Session {
	return &DefaultSession{
		id:    id,
		M:     make(map[string]interface{}),
		store: store,
	}
}

var _ Session = &DefaultSession{}

type DefaultSession struct {
	id      string
	M       map[string]interface{}
	store   Store
	expires int64
}

func (s DefaultSession) ID() string {
	return s.id
}

func (s DefaultSession) Get(k string) (interface{}, bool) {
	o, b := s.M[k]
	if b {
		return o, true
	}
	return nil, false
}

func (s DefaultSession) Set(k string, o interface{}) {
	s.M[k] = o
}

func (s DefaultSession) Remove(k string) {
	delete(s.M, k)
}

func (s *DefaultSession) Save() {
	s.store.Set(s)
}

// Returns true if the item has expired.
func (s DefaultSession) IsExpired() bool {
	if s.expires == 0 {
		return false
	}
	return time.Now().UnixNano() > s.expires
}

func (s *DefaultSession) ExpiresAt(nano int64) {
	s.expires = nano
}

func (s *DefaultSession) Store() Store {
	return s.store
}

func (s DefaultSession) MustGet(k string) interface{} {
	o, b := s.M[k]
	if b {
		return o
	}
	return nil
}
