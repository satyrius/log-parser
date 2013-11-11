package stat

import (
	"github.com/satyrius/gonx"
	"regexp"
	"time"
)

type Stat struct {
	StartedAt     time.Time
	ElapsedTime   time.Duration
	Logs          []string
	GroupBy       string
	GroupByRegexp *regexp.Regexp
	EntriesParsed int
	Data          map[string]int
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
		Data:          make(map[string]int),
	}
}

func (s *Stat) Add(record *gonx.Entry) (err error) {
	s.EntriesParsed++
	value := (*record)[s.GroupBy]
	if _, ok := s.Data[value]; ok {
		s.Data[value]++
	} else {
		s.Data[value] = 1
	}
	return
}
