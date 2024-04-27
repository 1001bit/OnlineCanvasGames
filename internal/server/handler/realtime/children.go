package realtime

type Children[T any] struct {
	childMap       map[*T]bool
	connectChan    chan *T
	disconnectChan chan *T
}

func MakeChildren[T any]() Children[T] {
	return Children[T]{
		childMap:       make(map[*T]bool),
		connectChan:    make(chan *T),
		disconnectChan: make(chan *T),
	}
}

func (children *Children[T]) DisconnectChild(child *T) {
	children.disconnectChan <- child
}

func (children *Children[T]) ConnectChild(child *T) {
	children.connectChan <- child
}

// Children but with int map
type ChildrenWithID[T any] struct {
	idMap          map[int]*T
	connectChan    chan *T
	disconnectChan chan *T
}

func MakeChildrenWithID[T any]() ChildrenWithID[T] {
	return ChildrenWithID[T]{
		idMap:          make(map[int]*T),
		connectChan:    make(chan *T),
		disconnectChan: make(chan *T),
	}
}

func (children *ChildrenWithID[T]) DisconnectChild(child *T) {
	children.disconnectChan <- child
}

func (children *ChildrenWithID[T]) ConnectChild(child *T) {
	children.connectChan <- child
}
