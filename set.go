package main

type Set struct {
	content map[string]struct{}
}

func NewSet(elements ...string) *Set {
	s := &Set{content: map[string]struct{}{}}
	s.Add(elements...)
	return s
}

func (s *Set) Add(elements ...string) {
	for _, element := range elements {
		s.content[element] = struct{}{}
	}
}

func (s *Set) Values() (values []string) {
	if s == nil {
		return nil
	}
	for k := range s.content {
		values = append(values, k)
	}
	return values
}
