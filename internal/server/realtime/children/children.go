package children

type Children[T any] struct {
	Channels[T]
	ChildMap map[*T]bool
}

func MakeChildren[T any]() Children[T] {
	return Children[T]{
		ChildMap: make(map[*T]bool),
		Channels: MakeChannels[T](),
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
