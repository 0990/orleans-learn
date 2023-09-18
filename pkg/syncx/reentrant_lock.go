package syncx

import (
	"github.com/0990/orleans-learn/pkg/util"
	"sync"
)

func NewReentrantLock() sync.Locker {
	res := &ReentrantLock{
		lock:      new(sync.Mutex),
		cond:      nil,
		recursion: 0,
		host:      0,
	}
	res.cond = sync.NewCond(res.lock)
	return res
}

type ReentrantLock struct {
	lock      *sync.Mutex
	cond      *sync.Cond
	recursion int32
	host      int64
}

func (rt *ReentrantLock) Lock() {
	id := util.GetGoroutineID()

	rt.lock.Lock()
	defer rt.lock.Unlock()

	if rt.host == id {
		rt.recursion++
		return
	}

	for rt.recursion != 0 {
		rt.cond.Wait()
	}

	rt.host = id
	rt.recursion = 1
}

func (rt *ReentrantLock) Unlock() {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	if rt.recursion == 0 || rt.host != util.GetGoroutineID() {
		panic("unlock error")
	}

	rt.recursion--
	if rt.recursion == 0 {
		rt.cond.Signal()
	}
}
