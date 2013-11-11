package stat

import (
	"fmt"
	"github.com/satyrius/gonx"
	"regexp"
	"time"
)

type Item struct {
	Name  string
	Count int
}

func NewItem(name string) *Item {
	return &Item{Name: name, Count: 0}
}

func (i *Item) Update(entry *gonx.Entry) {
	i.Count++
}

type Stat struct {
	StartedAt     time.Time
	ElapsedTime   time.Duration
	Logs          []string
	GroupBy       string
	GroupByRegexp *regexp.Regexp
	EntriesParsed int
	Data          []*Item
	index         map[string]int
}

func NewStat(groupBy string, regexpPattern string) *Stat {
	var re *regexp.Regexp
	if regexpPattern != "" {
		re = regexp.MustCompile(regexpPattern)
	}
	return &Stat{
		EntriesParsed: 0,
		StartedAt:     time.Now(),
		GroupBy:       groupBy,
		GroupByRegexp: re,
		index:         make(map[string]int),
	}
}

func (s *Stat) Get(name string) *Item {
	if id, ok := s.index[name]; ok {
		return s.Data[id]
	}
	return nil
}

func (s *Stat) Add(record *gonx.Entry) (err error) {
	value, ok := (*record)[s.GroupBy]
	if !ok {
		return fmt.Errorf("Field '%v' does not found in record %+v", s.GroupBy, *record)
	}

	if s.GroupByRegexp != nil {
		submatch := s.GroupByRegexp.FindStringSubmatch(value)
		if submatch == nil {
			return fmt.Errorf("Entry's '%v' value '%v' does not match Regexp '%v'",
				s.GroupBy, value, s.GroupByRegexp)
		}
		value = submatch[len(submatch)-1]
	}

	// Update existing stat item or create new one
	if id, ok := s.index[value]; ok {
		s.Data[id].Update(record)
	} else {
		item := NewItem(value)
		item.Update(record)
		s.Data = append(s.Data, item)
		s.index[value] = s.Len() - 1
	}

	s.EntriesParsed++
	return
}

func (s *Stat) AddLog(file string) {
	s.Logs = append(s.Logs, file)
}

func (s *Stat) Stop() time.Duration {
	s.ElapsedTime = time.Since(s.StartedAt)
	return s.ElapsedTime
}
