package session

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultGet(t *testing.T) {
	s := New("test", nil)
	s.Set("t", "haha")

	assert.Equal(t, "test", s.ID())
	r, b := s.Get("t")
	assert.True(t, b)
	assert.Equal(t, "haha", r)

	r1, b1 := s.Get("q")
	assert.True(t, !b1)
	assert.Nil(t, r1)
}
