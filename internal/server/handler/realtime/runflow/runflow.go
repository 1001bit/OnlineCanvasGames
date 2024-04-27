package runflow

type RunFlow struct {
	stopChan chan struct{}
	doneChan chan struct{}
}

func MakeRunFlow() RunFlow {
	return RunFlow{
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
	}
}

// Stop
func (rf *RunFlow) Stopped() <-chan struct{} {
	return rf.stopChan
}

func (rf *RunFlow) Stop() {
	rf.stopChan <- struct{}{}
}

// Done
func (rf *RunFlow) Done() <-chan struct{} {
	return rf.doneChan
}

func (rf *RunFlow) CloseDone() {
	close(rf.doneChan)
}
