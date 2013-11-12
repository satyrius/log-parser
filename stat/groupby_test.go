package stat

import (
	"github.com/satyrius/gonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGroupByValue(t *testing.T) {
	groupBy := GroupByValue("request")

	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	value, err := groupBy(entry)
	assert.NoError(t, err)
	assert.Equal(t, value, request)
}

func TestGroupByRegexp(t *testing.T) {
	groupBy := GroupByRegexp("request", `^\w+\s+(\S+)(?:\?|$)`, "")

	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	value, err := groupBy(entry)
	assert.NoError(t, err)
	assert.Equal(t, value, "/foo/bar")
}

func TestGroupByBadRegexp(t *testing.T) {
	groupBy := GroupByRegexp("request", `^(\d+)$`, "")

	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	value, err := groupBy(entry)
	assert.Error(t, err)
	// Return raw value on error
	assert.Equal(t, value, request)
}

func TestGroupByNoSubmatchRegexp(t *testing.T) {
	groupBy := GroupByRegexp("request", `^\w+`, "")
	request := "GET /foo/bar"
	entry := &gonx.Entry{"request": request}
	value, err := groupBy(entry)
	assert.NoError(t, err)
	assert.Equal(t, value, "GET")
}
