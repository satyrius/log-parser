package stat

import (
	"github.com/satyrius/gonx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstructor(t *testing.T) {
	stat := NewStat("request", "")
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.Equal(t, stat.GroupBy, "request")
}

func TestTiming(t *testing.T) {
	start := time.Now()
	stat := NewStat("request", "")
	assert.WithinDuration(t, start, stat.StartedAt, time.Duration(time.Millisecond),
		"Constructor should setup StartedAt")
	assert.Equal(t, stat.ElapsedTime, 0)
	elapsed := stat.Stop()
	assert.NotEqual(t, stat.ElapsedTime, 0)
	assert.Equal(t, stat.ElapsedTime, elapsed)
}

func TestAddLog(t *testing.T) {
	stat := NewStat("request", "")
	assert.Empty(t, stat.Logs)
	file := "/var/log/nginx/access.log"
	stat.AddLog(file)
	assert.Equal(t, stat.Logs, []string{file})
}

func TestAdd(t *testing.T) {
	stat := NewStat("request", "")
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}

	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 1)
	item := stat.Get(request)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, request)
	assert.Equal(t, item.Count, 1)
	assert.Equal(t, stat.Len(), 1)

	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 2)
	item = stat.Get(request)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, request)
	assert.Equal(t, item.Count, 2)
	assert.Equal(t, stat.Len(), 1)
}

func TestAddInvalid(t *testing.T) {
	stat := NewStat("request", "")
	entry := &gonx.Entry{"foo": "bar"}
	assert.Error(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.Equal(t, stat.Len(), 0)
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

func TestGroupByRegexp(t *testing.T) {
	stat := NewStat("request", `^\w+\s+(\S+)`)
	uri := "/foo/bar"
	request := "GET " + uri
	entry := &gonx.Entry{"request": request}

	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 1)

	item := stat.Get(request)
	assert.Nil(t, item)

	// Uri should be used as data key because we have regexp to extract it
	item = stat.Get(uri)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, uri)
	assert.Equal(t, item.Count, 1)
}

func TestBadRegexp(t *testing.T) {
	// Invalid Regexp required request to be numeric field
	stat := NewStat("request", `^(\d+)$`)
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	assert.Error(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 0)
}

func TestNoSubmatchRegexp(t *testing.T) {
	// Invalid Regexp required request to be numeric field
	stat := NewStat("request", `^\w+`)
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	assert.NoError(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 1)

	item := stat.Get(request)
	assert.Nil(t, item)

	// Request method was used for grouping
	item = stat.Get("GET")
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, "GET")
	assert.Equal(t, item.Count, 1)
}
