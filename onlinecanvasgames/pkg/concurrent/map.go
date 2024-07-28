package concurrent

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	mutex sync.RWMutex
	items map[K]V
}

func MakeMap[K comparable, V any]() ConcurrentMap[K, V] {
	return ConcurrentMap[K, V]{
		mutex: sync.RWMutex{},
		items: make(map[K]V),
	}
}

func (m *ConcurrentMap[K, V]) Set(key K, val V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.items[key] = val
}

func (m *ConcurrentMap[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.items, key)
}

func (m *ConcurrentMap[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	clear(m.items)
}

func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	val, ok := m.items[key]
	return val, ok
}

func (m *ConcurrentMap[K, V]) GetMapForRead() (map[K]V, func()) {
	m.mutex.RLock()
	return m.items, func() { m.mutex.RUnlock() }
}

func (m *ConcurrentMap[K, V]) Length() int {
	return len(m.items)
}
