package catalog

import (
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/core/Message"
)

func (a *ActivationData) GrainId() IDs.GrainId {
	return a.Address.GrainId
}

func (a *ActivationData) ActivationId() IDs.ActivationId {
	return a.Address.ActivationId
}

func (a *ActivationData) Activate() {
	a.ScheduleOperation(CommandActivate{})
}

func (a *ActivationData) ReceiveMessage(message Message.Message) {
	a.lock.Lock()
	a._waitingRequests.PushBack(message)
	a.lock.Unlock()

	a._workSignal.Signal()
}

// TODO
func (a *ActivationData) Post(f func()) {
	a.lock.Lock()
	a._waitingFuncs.PushBack(f)
	a.lock.Unlock()

	a._workSignal.Signal()
}

func (a *ActivationData) Deactivate(reason string) {
	a.StartDeactivating(reason)

	a.ScheduleOperation(CommandDeactivate{})
}
