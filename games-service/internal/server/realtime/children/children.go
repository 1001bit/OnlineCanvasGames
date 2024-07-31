package children

import "github.com/neinBit/ocg-games-service/pkg/concurrent"

type SetChildren[T comparable] struct {
	Channels[T]
	ChildrenSet concurrent.ConcurrentSet[T]
}

func MakeSetChildren[T comparable]() SetChildren[T] {
	return SetChildren[T]{
		ChildrenSet: concurrent.MakeSet[T](),
		Channels:    MakeChannels[T](),
	}
}

type MapChildren[K comparable, V any] struct {
	Channels[V]
	ChildrenMap concurrent.ConcurrentMap[K, V]
}

func MakeMapChildren[K comparable, V any]() MapChildren[K, V] {
	return MapChildren[K, V]{
		ChildrenMap: concurrent.MakeMap[K, V](),
		Channels:    MakeChannels[V](),
	}
}
