package stat

import (
	"github.com/satyrius/gonx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstructor(t *testing.T) {
	stat := NewStat(nil, "request", "")
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.Equal(t, stat.GroupBy, "request")
}

func TestTiming(t *testing.T) {
	start := time.Now()
	stat := NewStat(nil, "request", "")
	assert.WithinDuration(t, start, stat.StartedAt, time.Duration(time.Millisecond),
		"Constructor should setup StartedAt")
	assert.Equal(t, stat.ElapsedTime, 0)
	elapsed := stat.Stop()
	assert.NotEqual(t, stat.ElapsedTime, 0)
	assert.Equal(t, stat.ElapsedTime, elapsed)
}

func TestAddLog(t *testing.T) {
	stat := NewStat(nil, "request", "")
	assert.Empty(t, stat.Logs)
	file := "/var/log/nginx/access.log"
	stat.AddLog(file)
	assert.Equal(t, stat.Logs, []string{file})
}

func TestAdd(t *testing.T) {
	stat := NewStat(nil, "request", "")
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}

	// Use whole GroupBy variable value, because there is
	// no GroupByRegexp specified for Stat
	value, err := stat.GetGroupByValue(entry)
	assert.NoError(t, err)
	assert.Equal(t, value, request)

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
	stat := NewStat(nil, "request", "")
	entry := &gonx.Entry{"foo": "bar"}
	assert.Error(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.Equal(t, stat.Len(), 0)
}

func TestEmptyRegexp(t *testing.T) {
	stat := NewStat(nil, "request", "")
	assert.Nil(t, stat.GroupByRegexp)
}

func TestRegexp(t *testing.T) {
	exp := `^\w+\s+(\S+)(?:\?|$)`
	stat := NewStat(nil, "request", exp)
	assert.Equal(t, stat.GroupByRegexp.String(), exp)

	uri := "/foo/bar"
	request := "GET " + uri
	entry := &gonx.Entry{"request": request}
	value, err := stat.GetGroupByValue(entry)
	assert.NoError(t, err)
	// Uri should be used as group by value because we have
	// regexp to extract it
	assert.Equal(t, value, uri)
}

func TestBadRegexp(t *testing.T) {
	// Invalid Regexp required request to be numeric field
	stat := NewStat(nil, "request", `^(\d+)$`)
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	_, err := stat.GetGroupByValue(entry)
	assert.Error(t, err)
	assert.Equal(t, stat.Add(entry), err)
	assert.Equal(t, stat.EntriesParsed, 0)
}

func TestNoSubmatchRegexp(t *testing.T) {
	stat := NewStat(nil, "request", `^\w+`)
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	value, err := stat.GetGroupByValue(entry)
	assert.NoError(t, err)
	// The group by regexp does not have submatch pattern, thats
	// why the whole regexp match (which is `GET`) will be used
	assert.Equal(t, value, "GET")
}
