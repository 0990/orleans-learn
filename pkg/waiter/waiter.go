package waiter

type Waiter struct {
	signal chan struct{}
}

func NewWaiter() *Waiter {
	return &Waiter{
		signal: make(chan struct{}, 1),
	}
}

func (w *Waiter) Wait() {
	<-w.signal
}

func (w *Waiter) Signal() {
	select {
	case w.signal <- struct{}{}:
	default:

	}
}
