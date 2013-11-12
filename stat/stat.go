package stat

import (
	"fmt"
	"github.com/satyrius/gonx"
	"regexp"
	"time"
)

// Entry data reducer. It accept links to the Item and Entry instances and
// should return reduced data. It should use Item.AggValue and Entry data.
type Aggregator func(item *Item, entry *gonx.Entry) (float64, error)

type Stat struct {
	StartedAt     time.Time
	ElapsedTime   time.Duration
	Logs          []string
	EntriesParsed int

	// Grouping and aggregation callbacks
	groupBy GroupBy
	agg     Aggregator

	// Store collected data as a slice of Items, but keep Item name/id index
	// tracking for sorting
	Data  []*Item
	index map[string]int
}

// Process group by field value and return exact value
type GroupBy func(entry *gonx.Entry) (string, error)

func GroupByValue(name string) GroupBy {
	return func(entry *gonx.Entry) (string, error) {
		return entry.Get(name)
	}
}

func GroupByRegexp(name string, pattern string) GroupBy {
	re := regexp.MustCompile(pattern)
	return func(entry *gonx.Entry) (value string, err error) {
		value, err = entry.Get(name)
		if err != nil {
			return
		}
		submatch := re.FindStringSubmatch(value)
		if submatch == nil {
			err = fmt.Errorf("Entry's '%v' value '%v' does not match Regexp '%v'",
				value, re)
			return
		}
		value = submatch[len(submatch)-1]
		return
	}
}

// Creates new stat aggregator for gonx.Entry processing. agg of type Aggregator is a callback
// function that reduce entries data. groupBy is GroupBy callback for data grouping.
func NewStat(agg Aggregator, groupBy GroupBy) *Stat {
	return &Stat{
		EntriesParsed: 0,
		StartedAt:     time.Now(),
		groupBy:       groupBy,
		agg:           agg,
		index:         make(map[string]int),
	}
}

func (s *Stat) Get(name string) *Item {
	if id, ok := s.index[name]; ok {
		return s.Data[id]
	}
	return nil
}

func (s *Stat) Add(entry *gonx.Entry) (err error) {
	value, err := s.groupBy(entry)
	if err != nil {
		return
	}

	// Update existing stat item or create new one
	if id, ok := s.index[value]; ok {
		err = s.Data[id].Update(entry)
	} else {
		item := NewItem(value, s.agg)
		err = item.Update(entry)
		if err == nil {
			s.Data = append(s.Data, item)
			s.index[value] = s.Len() - 1
		}
	}

	if err == nil {
		s.EntriesParsed++
	}

	return
}

func (s *Stat) AddLog(file string) {
	s.Logs = append(s.Logs, file)
}

func (s *Stat) Stop() time.Duration {
	s.ElapsedTime = time.Since(s.StartedAt)
	return s.ElapsedTime
}
