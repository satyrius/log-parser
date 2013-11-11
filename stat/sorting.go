package stat

func (s *Stat) Len() int {
	return len(s.Data)
}

func (s *Stat) Less(i, j int) bool {
	if s.agg == nil {
		return s.Data[i].Count > s.Data[j].Count
	} else {
		return s.Data[i].AggValue > s.Data[j].AggValue
	}
}

func (s *Stat) Swap(i, j int) {
	a, b := s.Data[i], s.Data[j]
	// Swap item links
	s.Data[i], s.Data[j] = b, a
	// And fix name index
	s.index[a.Name], s.index[b.Name] = j, i
}
