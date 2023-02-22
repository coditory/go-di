package di

import "sort"

type Set[T comparable] struct {
	set      map[T]int
	inverted map[int]T
	index    int
}

func NewSet[T comparable]() *Set[T] {
	return NewSetWithSize[T](1)
}

func NewSetWithSize[T comparable](size int) *Set[T] {
	return &Set[T]{
		set:      make(map[T]int, size),
		inverted: make(map[int]T, size),
	}
}

func (s *Set[T]) Add(value T) {
	s.index += 1
	index := s.index
	s.set[value] = index
	s.inverted[index] = value
}

func (s Set[T]) Remove(value T) {
	index, ok := s.set[value]
	if !ok {
		return
	}
	delete(s.set, value)
	delete(s.inverted, index)
}

func (s Set[T]) Contains(key T) bool {
	_, ok := s.set[key]
	return ok
}

func (s Set[T]) ToSlice() []T {
	indexes := make([]int, len(s.set))
	i := 0
	for _, v := range s.set {
		indexes[i] = v
		i += 1
	}
	sort.Ints(indexes)
	result := make([]T, len(s.set))
	i = 0
	for _, index := range indexes {
		result[i] = s.inverted[index]
		i += 1
	}
	return result
}
