package catalog

import (
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/core/Message"
	"github.com/0990/orleans-learn/pkg/syncx"
	"github.com/0990/orleans-learn/pkg/waiter"
	"sync"
	"time"
)

type ActivationData struct {
	Address IDs.GrainAddress
	State   ActivationState

	lock                 sync.Locker
	_workSignal          *waiter.Waiter
	_pendingOperations   []any
	_waitingRequests     *syncx.List[Message.Message]
	_waitingFuncs        *syncx.List[func()]
	_runningRequests     *syncx.SyncMap[Message.Message, time.Time]
	IsCurrentlyExecuting bool

	_shared       GrainTypeShardContext
	GrainInstance any

	DeactivationException error
	ForwaringAddress      string
	DeactivationReason    string
}

func NewActivationData(address IDs.GrainAddress, _shared GrainTypeShardContext) *ActivationData {
	ad := &ActivationData{
		Address:          address,
		_shared:          _shared,
		State:            Create,
		_waitingRequests: syncx.NewList[Message.Message](),
		_waitingFuncs:    syncx.NewList[func()](),
		_workSignal:      waiter.NewWaiter(),
		_runningRequests: syncx.NewSyncMap[Message.Message, time.Time](),
		lock:             syncx.NewReentrantLock(),
	}
	go ad.RunMessageLoop()
	return ad
}
