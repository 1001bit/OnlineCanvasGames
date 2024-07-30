package children

import "github.com/neinBit/ocg-games-service/pkg/concurrent"

type Children[T any] struct {
	Channels[T]
	ChildrenSet concurrent.ConcurrentSet[*T]
}

func MakeChildren[T any]() Children[T] {
	return Children[T]{
		ChildrenSet: concurrent.MakeSet[*T](),
		Channels:    MakeChannels[T](),
	}
}

type ChildrenWithID[T any] struct {
	Channels[T]
	IDMap concurrent.ConcurrentMap[int, *T]
}

func MakeChildrenWithID[T any]() ChildrenWithID[T] {
	return ChildrenWithID[T]{
		IDMap:    concurrent.MakeMap[int, *T](),
		Channels: MakeChannels[T](),
	}
}
