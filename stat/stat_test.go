package stat

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstructor(t *testing.T) {
	start := time.Now()
	stat := NewStat("request", "")
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.WithinDuration(t, start, stat.StartedAt, time.Duration(time.Millisecond),
		"Constructor should setup StartedAt")
	assert.Equal(t, stat.GroupBy, "request")
}

func TestEmptyRegexp(t *testing.T) {
	stat := NewStat("request", "")
	assert.Nil(t, stat.GroupByRegexp)
}

func TestRegexp(t *testing.T) {
	exp := `^\w+\s+(\S+)(?:\?|$)`
	stat := NewStat("request", exp)
	assert.Equal(t, stat.GroupByRegexp.String(), exp)
}
