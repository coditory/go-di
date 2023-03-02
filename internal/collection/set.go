package collection

import "sort"

type Set[T comparable] struct {
	indexByValue map[T]int
	valueByIndex map[int]T
	index        int
}

func NewSet[T comparable]() *Set[T] {
	return NewSetWithSize[T](1)
}

func NewSetWithSize[T comparable](size int) *Set[T] {
	return &Set[T]{
		indexByValue: make(map[T]int, size),
		valueByIndex: make(map[int]T, size),
	}
}

func (s *Set[T]) Add(value T) {
	s.index += 1
	index := s.index
	s.indexByValue[value] = index
	s.valueByIndex[index] = value
}

func (s Set[T]) Remove(value T) {
	index, ok := s.indexByValue[value]
	if !ok {
		return
	}
	delete(s.indexByValue, value)
	delete(s.valueByIndex, index)
}

func (s Set[T]) Contains(key T) bool {
	_, ok := s.indexByValue[key]
	return ok
}

func (s Set[T]) ToSlice() []T {
	indexes := make([]int, len(s.indexByValue))
	i := 0
	for _, v := range s.indexByValue {
		indexes[i] = v
		i += 1
	}
	sort.Ints(indexes)
	result := make([]T, len(s.indexByValue))
	i = 0
	for _, index := range indexes {
		result[i] = s.valueByIndex[index]
		i += 1
	}
	return result
}
