package stat

import (
	"github.com/satyrius/gonx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstructor(t *testing.T) {
	stat := NewStat(nil, GroupByValue("request"))
	assert.Equal(t, stat.EntriesParsed, 0)
}

func TestTiming(t *testing.T) {
	start := time.Now()
	stat := NewStat(nil, GroupByValue("request"))
	assert.WithinDuration(t, start, stat.StartedAt, time.Duration(time.Millisecond),
		"Constructor should setup StartedAt")
	assert.Equal(t, stat.ElapsedTime, 0)
	elapsed := stat.Stop()
	assert.NotEqual(t, stat.ElapsedTime, 0)
	assert.Equal(t, stat.ElapsedTime, elapsed)
}

func TestAddLog(t *testing.T) {
	stat := NewStat(nil, GroupByValue("request"))
	assert.Empty(t, stat.Logs)
	file := "/var/log/nginx/access.log"
	stat.AddLog(file)
	assert.Equal(t, stat.Logs, []string{file})
}

func TestAdd(t *testing.T) {
	stat := NewStat(nil, GroupByValue("request"))
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
	stat := NewStat(nil, GroupByValue("request"))
	entry := &gonx.Entry{"foo": "bar"}
	assert.Error(t, stat.Add(entry))
	assert.Equal(t, stat.EntriesParsed, 0)
	assert.Equal(t, stat.Len(), 0)
}
