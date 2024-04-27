package children

type Children[T any] struct {
	ChildMap       map[*T]bool
	connectChan    chan *T
	disconnectChan chan *T
}

func MakeChildren[T any]() Children[T] {
	return Children[T]{
		ChildMap:       make(map[*T]bool),
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

func (children *Children[T]) ToConnect() <-chan *T {
	return children.connectChan
}

func (children *Children[T]) ToDisconnect() <-chan *T {
	return children.disconnectChan
}
