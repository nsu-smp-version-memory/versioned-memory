package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	s := NewSource()
	assert.Equal(t, s.currVersion.Load(), uint64(0))
}

func TestSource_Peek(t *testing.T) {
	s := NewSource()
	assert.Equal(t, s.currVersion.Load(), uint64(0))
	assert.Equal(t, s.Peek(), uint64(1))
	assert.Equal(t, s.currVersion.Load(), uint64(1))
}
