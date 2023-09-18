package syncx

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewReentrantLock(t *testing.T) {
	l := NewReentrantLock()
	F(l)

	go func() {
		F(l)
	}()

	time.Sleep(time.Minute)
}

func F(locker sync.Locker) {
	locker.Lock()
	G(locker)
	locker.Unlock()
}

func G(locker sync.Locker) {
	locker.Lock()
	time.Sleep(time.Second * 5)
	fmt.Println("G")
	locker.Unlock()
}
