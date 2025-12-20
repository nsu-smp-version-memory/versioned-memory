package versioned

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	s := NewSet()
	assert.Empty(t, s.values)
}

func TestSet_Add(t *testing.T) {
	s := NewSet()

	s = s.Add(1)
	assert.Contains(t, s.values, 1)
}

func TestSet_Remove(t *testing.T) {
	s := NewSet()

	s = s.Remove(1)
	assert.Empty(t, s.values)

	s = s.Add(1)
	assert.Contains(t, s.values, 1)

	s = s.Remove(1)
	assert.Empty(t, s.values)
}
