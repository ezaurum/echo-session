package memstore

import (
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/ezaurum/cthulthu/session"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestGetNew(t *testing.T) {
	sm := Default()

	s := sm.GetNew()

	assert.NotNil(t, s)
	assert.NotEmpty(t, s.ID(), "session ID cannot be empty")
}

func TestGetNewSerial(t *testing.T) {
	sm := Default()

	var wg sync.WaitGroup
	wg.Add(3)

	n := func(c chan string) {
		s := sm.GetNew()
		c <- s.ID()
	}

	c := make(chan string)

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			go n(c)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			go n(c)
		}
	}()

	go func() {
		defer wg.Done()
		for j := 0; j < 100000; j++ {
			s0 := <-c
			s1 := <-c
			assert.NotEqual(t, s1, s0)
		}
	}()

	wg.Wait()
}

func TestGet(t *testing.T) {
	sm := Default()

	s := sm.GetNew()
	sessionID := s.ID()

	s0, isExist := sm.Get(sessionID)

	assert.True(t, isExist, "Session is not exist.")
	assert.NotNil(t, s0, "Session is nil.")

	assert.Equal(t, s, s0)

}

func TestGetByGoroutine(t *testing.T) {
	sm := Default()

	var wg sync.WaitGroup
	wg.Add(3)

	n := func(c chan string) {
		s := sm.GetNew()
		c <- s.ID()
	}

	c := make(chan string)

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			go n(c)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100000; i++ {
			go n(c)
		}
	}()

	go func() {
		defer wg.Done()
		for j := 0; j < 100000; j++ {
			s0, e0 := sm.Get(<-c)
			s0b, e0b := sm.Get(s0.ID())
			assert.True(t, e0, "Session is not exist.")
			assert.True(t, e0b, "Session is not exist.")
			assert.Equal(t, s0, s0b)

			s1, e1 := sm.Get(<-c)
			s1b, e1b := sm.Get(s1.ID())
			assert.True(t, e1, "Session is not exist.")
			assert.True(t, e1b, "Session is not exist.")
			assert.Equal(t, s1, s1b)
		}
	}()

	wg.Wait()
}

func TestQueryAllList(t *testing.T) {
	sessionCount := 10000

	sm := Default()

	sMap := make(session.SessionMap)
	for i := 0; i < sessionCount; i++ {
		ss := sm.GetNew()
		sMap[ss.ID()] = ss
	}

	assert.True(t, sessionCount == sm.Count())

	for k, v := range sm.Sessions() {
		assert.Equal(t, sMap[k], v)
	}
}

//없으면 false 하도록
func TestGetFail(t *testing.T) {
	sm := Default()
	get, b := sm.Get("nothign")

	assert.True(t, !b)
	assert.Nil(t, get)
}

func TestExpire(t *testing.T) {
	sm := New(snowflake.New(0), time.Second, time.Hour)

	refreshed := sm.GetNew()
	notRefreshed := sm.GetNew()

	time.Sleep(time.Millisecond * 500)

	sm.Set(refreshed)

	time.Sleep(time.Millisecond * 500)

	get, b := sm.Get(notRefreshed.ID())

	assert.True(t, !b)
	assert.Nil(t, get)
	assert.True(t, notRefreshed.IsExpired())

	get0, b0 := sm.Get(refreshed.ID())

	assert.True(t, b0)
	assert.NotNil(t, get0)
	assert.True(t, !get0.IsExpired())

}
