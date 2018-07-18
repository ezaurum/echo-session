package session

import "time"

var _ Session = &DefaultSession{}

func New(id string, store Store, ipAddress string, agent string) Session {
	return &DefaultSession{
		key:       id,
		M:         make(map[string]interface{}),
		store:     store,
		ipAddress: ipAddress,
		agent:     agent,
	}
}

type DefaultSession struct {
	key       string
	M         map[string]interface{}
	store     Store
	expires   int64
	ipAddress string
	agent     string
}

func (s DefaultSession) IPAddress() string {
	return s.ipAddress
}

func (s DefaultSession) Agent() string {
	return s.agent
}

func (s DefaultSession) Key() string {
	return s.key
}

func (s DefaultSession) Get(k string) (interface{}, bool) {
	o, b := s.M[k]
	if b {
		return o, true
	}
	return nil, false
}

func (s *DefaultSession) Set(k string, o interface{}) {
	s.M[k] = o
}

func (s *DefaultSession) Remove(k string) {
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

func (s *DefaultSession) Extend(duration time.Duration) {
	s.ExpiresAt(time.Now().Add(duration).UnixNano())
}

func (s DefaultSession) Destroy() {
	s.store.Delete(s.Key())
}
