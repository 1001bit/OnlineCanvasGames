package set

type Set[T comparable] map[T]struct{}

func MakeSet[T comparable](slice []T) Set[T] {
	set := make(Set[T])

	for i := range slice {
		set[slice[i]] = struct{}{}
	}

	return set
}

func MakeEmptySet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Has(elem T) bool {
	_, ok := s[elem]
	return ok
}

func (s Set[T]) Insert(elem T) {
	s[elem] = struct{}{}
}

func (s Set[T]) Delete(elem T) {
	delete(s, elem)
}

func (s Set[T]) Clear() {
	clear(s)
}