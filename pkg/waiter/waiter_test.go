package waiter

import (
	"testing"
	"time"
)

func TestWaiter(t *testing.T) {
	w := NewWaiter()

	go func() {
		time.Sleep(1 * time.Second)
		w.Signal()
		w.Signal()
		w.Signal()
	}()

	start := time.Now()
	w.Wait()
	elapsed := time.Since(start)

	if elapsed < 1*time.Second {
		t.Errorf("Waiter.Wait() returned before 1 second")
	}
}
