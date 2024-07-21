package children

import "github.com/1001bit/OnlineCanvasGames/pkg/set"

type Children[T any] struct {
	Channels[T]
	ChildrenSet set.Set[*T]
}

func MakeChildren[T any]() Children[T] {
	return Children[T]{
		ChildrenSet: set.MakeEmptySet[*T](),
		Channels:    MakeChannels[T](),
	}
}

type ChildrenWithID[T any] struct {
	Channels[T]
	IDMap map[int]*T
}

func MakeChildrenWithID[T any]() ChildrenWithID[T] {
	return ChildrenWithID[T]{
		IDMap:    make(map[int]*T),
		Channels: MakeChannels[T](),
	}
}
