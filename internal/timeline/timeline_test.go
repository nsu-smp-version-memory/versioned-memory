package timeline

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	timeline := New[int]()
	assert.Nil(t, timeline.last)
}

func TestTimeline_NextChange(t *testing.T) {
	timeline := New[int]()

	v := rand.Int()
	timeline = timeline.NextChange(v)
	assert.NotNil(t, timeline.last)
	assert.Equal(t, timeline.last.version, uint64(1))
	assert.Nil(t, timeline.last.prev)
	assert.Equal(t, timeline.last.diff, v)

	v1 := rand.Int()
	timeline1 := timeline.NextChange(v1)
	assert.NotNil(t, timeline1.last)
	assert.Equal(t, timeline1.last.version, uint64(2))
	assert.Equal(t, timeline1.last.prev, timeline.last)
	assert.Equal(t, timeline1.last.diff, v1)

	v2 := rand.Int()
	timeline2 := timeline.NextChange(v2)
	assert.NotNil(t, timeline2.last)
	assert.Equal(t, timeline2.last.version, uint64(3))
	assert.Equal(t, timeline2.last.prev, timeline.last)
	assert.Equal(t, timeline2.last.diff, v2)
}
