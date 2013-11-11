package stat

import (
	"fmt"
	"github.com/satyrius/gonx"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestItemUpdateCount(t *testing.T) {
	item := NewItem("foo", nil)
	assert.Equal(t, item.Count, 0)
	item.Update(&gonx.Entry{})
	assert.Equal(t, item.Count, 1)
	item.Update(&gonx.Entry{})
	assert.Equal(t, item.Count, 2)
}

func TestItemAggeagator(t *testing.T) {
	sum := func(item *Item, entry *gonx.Entry) (val float64, err error) {
		if strVal, ok := (*entry)["foo"]; ok {
			v, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return 0, err
			}
			val = item.AggValue + v
		} else {
			err = fmt.Errorf("Invalid entry data")
		}
		return
	}
	item := NewItem("foo", sum)
	assert.Equal(t, item.AggValue, 0)

	item.Update(&gonx.Entry{"foo": "1"})
	assert.Equal(t, item.AggValue, 1)

	item.Update(&gonx.Entry{"foo": "2"})
	assert.Equal(t, item.AggValue, 3)
}
