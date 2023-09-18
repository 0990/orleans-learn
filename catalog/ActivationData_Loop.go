package catalog

import (
	"github.com/0990/orleans-learn/core"
	"github.com/0990/orleans-learn/core/Message"
	"log/slog"
	"time"
)

func (a *ActivationData) RunMessageLoop() {
	for {
		if !a.IsCurrentlyExecuting {
			var operations []any
			a.lock.Lock()
			operations = a._pendingOperations
			a._pendingOperations = nil
			a.lock.Unlock()
			if len(operations) > 0 {
				a.ProcessOperationsAsync(operations)
			}
		}
		a.ProcessPendingRequests()
		a._workSignal.Wait()
	}
}

func (a *ActivationData) ProcessOperationsAsync(operations []any) {
	for _, op := range operations {
		switch cmd := op.(type) {
		case CommandActivate:
			a.ActivateSync()
		case CommandDeactivate:
			a.FinishDeactivating()
		case CommandDelay:
			time.Sleep(cmd.Duration)
		case CommandUnregisterFromCatalog:
			a.UnregisterMessageTarget()
		default:
			panic("unknown command")
		}
	}
}

func (a *ActivationData) ProcessPendingRequests() {
	for {
		a.lock.Lock()
		if a._waitingRequests.Len() <= 0 {
			break
		}
		if a.State != Valid {
			a.ProcessRequestsToInvalidActivation()
			break
		}
		message := a._waitingRequests.PopFront()
		a.RecordRunning(message, true)
		a.lock.Unlock()
		a.InvokeIncomingRequest(message)
	}

	a.lock.Unlock()
}

func (a *ActivationData) ProcessRequestsToInvalidActivation() {
	if a.State == Create || a.State == Activating {
		// Do nothing until the activation becomes either valid or invalid
		return
	}

	if a.State == Deactivating {
		IsStuckDeactivating := false
		if !IsStuckDeactivating {

			return
		}
	}

	if a.DeactivationException == nil || a.ForwaringAddress == "" {
		a.RerouteAllQueuedMessages()
	} else {
		a.RejectAllQueuedMessages()
	}
}

func (a *ActivationData) RerouteAllQueuedMessages() {
	a.lock.Lock()
	defer a.lock.Unlock()

	msgs := a.DequeueAllWaitingRequests()
	msgs.Range(func(message Message.Message) {
		slog.Info("RerouteAllQueuedMessages", "msg", message)
	})

	a._shared.InternalGrainRuntime.MessageCenter.ProcessRequestsToInvalidActivation(msgs, a.Address, a.ForwaringAddress, false)
}

func (a *ActivationData) RejectAllQueuedMessages() {
	a.lock.Lock()
	defer a.lock.Unlock()

	msgs := a.DequeueAllWaitingRequests()
	a._shared.InternalGrainRuntime.MessageCenter.ProcessRequestsToInvalidActivation(msgs, a.Address, a.ForwaringAddress, true)
}

func (a *ActivationData) InvokeIncomingRequest(message Message.Message) {
	//test only
	if message.Deactivate {
		a.Deactivate("receive deactivate message")
		return
	}

	if grainBase, ok := a.GrainInstance.(core.IGrainBase); ok {
		grainBase.OnMessage(message)
	}

	a.OnCompletedRequest(message)
}

func (a *ActivationData) OnCompletedRequest(message Message.Message) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a._runningRequests.Delete(message)
	a._workSignal.Signal()
}
