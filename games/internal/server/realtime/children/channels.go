package children

type Channels[T any] struct {
	connectChan    chan T
	disconnectChan chan T
}

func MakeChannels[T any]() Channels[T] {
	return Channels[T]{
		connectChan:    make(chan T),
		disconnectChan: make(chan T),
	}
}

func (ch *Channels[T]) ConnectChild(child T, done <-chan struct{}) {
	select {
	case ch.connectChan <- child:
		// connect child
	case <-done:
		// done
	}
}

func (ch *Channels[T]) DisconnectChild(child T, done <-chan struct{}) {
	select {
	case ch.disconnectChan <- child:
		// disconnect child
	case <-done:
		// done
	}
}

func (ch *Channels[T]) ToConnect() <-chan T {
	return ch.connectChan
}

func (ch *Channels[T]) ToDisconnect() <-chan T {
	return ch.disconnectChan
}
