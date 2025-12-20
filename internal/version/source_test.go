package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	s := NewSource()
	assert.Equal(t, s.currVersion.Load(), uint64(0))
}

func TestSource_Next(t *testing.T) {
	s := NewSource()
	assert.Equal(t, s.currVersion.Load(), uint64(0))
	assert.Equal(t, s.Next(), uint64(1))
	assert.Equal(t, s.currVersion.Load(), uint64(1))
}
