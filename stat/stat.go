package stat

import (
	"fmt"
	"github.com/satyrius/gonx"
	"regexp"
	"time"
)

type Aggregator func(item *Item, entry *gonx.Entry) (float64, error)

type Stat struct {
	StartedAt     time.Time
	ElapsedTime   time.Duration
	Logs          []string
	GroupBy       string
	GroupByRegexp *regexp.Regexp
	EntriesParsed int
	Data          []*Item
	agg           Aggregator
	index         map[string]int
}

func NewStat(agg Aggregator, groupBy string, regexpPattern string) *Stat {
	var re *regexp.Regexp
	if regexpPattern != "" {
		re = regexp.MustCompile(regexpPattern)
	}
	return &Stat{
		EntriesParsed: 0,
		StartedAt:     time.Now(),
		GroupBy:       groupBy,
		GroupByRegexp: re,
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

func (s *Stat) GetGroupByValue(entry *gonx.Entry) (value string, err error) {
	value, ok := (*entry)[s.GroupBy]
	if !ok {
		err = fmt.Errorf("Field '%v' does not found in record %+v", s.GroupBy, *entry)
		return
	}
	if s.GroupByRegexp != nil {
		submatch := s.GroupByRegexp.FindStringSubmatch(value)
		if submatch == nil {
			err = fmt.Errorf("Entry's '%v' value '%v' does not match Regexp '%v'",
				s.GroupBy, value, s.GroupByRegexp)
			return
		}
		value = submatch[len(submatch)-1]
	}
	return
}

func (s *Stat) Add(entry *gonx.Entry) (err error) {
	value, err := s.GetGroupByValue(entry)
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
