package stat

import (
	"github.com/satyrius/gonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItemUpdateCount(t *testing.T) {
	item := NewItem("foo")
	assert.Equal(t, item.Count, 0)
	item.Update(&gonx.Entry{})
	assert.Equal(t, item.Count, 1)
	item.Update(&gonx.Entry{})
	assert.Equal(t, item.Count, 2)
}
