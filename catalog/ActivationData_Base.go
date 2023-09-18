package catalog

import (
	"github.com/0990/orleans-learn/core/Message"
	"github.com/0990/orleans-learn/pkg/syncx"
	"time"
)

func (a *ActivationData) ScheduleOperation(operation any) {
	a.lock.Lock()
	a._pendingOperations = append(a._pendingOperations, operation)
	a.lock.Unlock()

	a._workSignal.Signal()
}

func (a *ActivationData) SetState(state ActivationState) {
	a.State = state
}

func (a *ActivationData) UnregisterMessageTarget() {
	a._shared.InternalGrainRuntime.Catalog.UnregisterMessageTarget(a)
	if a.GrainInstance != nil {
		a.SetGrainInstance(nil)
	}
}

func (a *ActivationData) SetGrainInstance(grainInstance any) {
	a.GrainInstance = grainInstance
}

func (a *ActivationData) DequeueAllWaitingRequests() *syncx.List[Message.Message] {
	a.lock.Lock()
	defer a.lock.Unlock()

	result := syncx.NewList[Message.Message]()
	a._waitingRequests.Range(func(msg Message.Message) {
		result.PushBack(msg)
	})
	a._waitingRequests.Clear()
	return result
}

func (a *ActivationData) RecordRunning(message Message.Message, isInterleaable bool) {
	a._runningRequests.Store(message, time.Now())
	if isInterleaable {
		return
	}
}
