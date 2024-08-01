package concurrent

import (
	"sync"

	"github.com/1001bit/ocg-games-service/pkg/set"
)

type ConcurrentSet[T comparable] struct {
	mutex sync.RWMutex
	items set.Set[T]
}

func MakeSet[T comparable]() ConcurrentSet[T] {
	return ConcurrentSet[T]{
		mutex: sync.RWMutex{},
		items: make(set.Set[T]),
	}
}

func (s *ConcurrentSet[T]) Insert(elem T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.items.Insert(elem)
}

func (s *ConcurrentSet[T]) Delete(elem T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.items, elem)
}

func (s *ConcurrentSet[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clear(s.items)
}

func (s *ConcurrentSet[T]) Has(elem T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.items.Has(elem)
}

func (s *ConcurrentSet[T]) GetSetForRead() (set.Set[T], func()) {
	s.mutex.RLock()
	return s.items, func() { s.mutex.RUnlock() }
}

func (s *ConcurrentSet[T]) Length() int {
	return len(s.items)
}
