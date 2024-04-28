package children

type ChildrenWithID[T any] struct {
	IDMap          map[int]*T
	connectChan    chan *T
	disconnectChan chan *T
}

func MakeChildrenWithID[T any]() ChildrenWithID[T] {
	return ChildrenWithID[T]{
		IDMap:          make(map[int]*T),
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

func (children *ChildrenWithID[T]) ToConnect() <-chan *T {
	return children.connectChan
}

func (children *ChildrenWithID[T]) ToDisconnect() <-chan *T {
	return children.disconnectChan
}
