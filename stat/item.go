package stat

import (
	"github.com/satyrius/gonx"
)

type Item struct {
	Name     string
	Count    int
	AggValue float64
	agg      Aggregator
}

func NewItem(name string, agg Aggregator) *Item {
	item := &Item{Name: name, Count: 0, AggValue: 0, agg: agg}
	return item
}

func (i *Item) Update(entry *gonx.Entry) (err error) {
	i.Count++
	if i.agg != nil {
		val, err := i.agg(i, entry)
		if err != nil {
			return err
		}
		i.AggValue = val
	}
	return
}
