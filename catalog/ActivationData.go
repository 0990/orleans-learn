package catalog

import (
	"fmt"
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

	lock                 sync.Mutex
	_workSignal          *waiter.Waiter
	_pendingOperations   []any
	_waitingRequests     *syncx.List[Message.Message]
	IsCurrentlyExecuting bool

	_shared GrainTypeShardContext
	//GrainInstance any

	ForwaringAddress   string
	DeactivationReason string
}

func NewActivationData(address IDs.GrainAddress, _shared GrainTypeShardContext) *ActivationData {
	ad := &ActivationData{
		Address:          address,
		_shared:          _shared,
		State:            Create,
		_waitingRequests: syncx.NewList[Message.Message](),
		_workSignal:      waiter.NewWaiter(),
	}
	go ad.RunMessageLoop()
	return ad
}

func (a *ActivationData) GrainId() IDs.GrainId {
	return a.Address.GrainId()
}

func (a *ActivationData) ActivationId() IDs.ActivationId {
	return a.Address.ActivationId()
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

func (a *ActivationData) ScheduleOperation(operation any) {
	a.lock.Lock()
	a._pendingOperations = append(a._pendingOperations, operation)
	a.lock.Unlock()

	a._workSignal.Signal()
}

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
		//a.FinishDeactivating()
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
			//TODO
			//ProcessRequestsToInvalidActivation()
			break
		}
		message := a._waitingRequests.PopFront()
		a.lock.Unlock()
		a.InvokeIncomingRequest(message)
	}
}

func (a *ActivationData) InvokeIncomingRequest(message Message.Message) {
	fmt.Println("InvokeIncomingRequest", message)
}

func (a *ActivationData) ActivateSync() {
	success := a.RegisterActivationInGrainDirectoryAndValidate()
	if !success {
		return
	}

	a.lock.Lock()
	a.SetState(Activating)
	a.lock.Unlock()

	success = a.CallActivateSync()
	if !success {
		return
	}

	a._workSignal.Signal()
}

func (a *ActivationData) SetState(state ActivationState) {
	a.State = state
}

func (a *ActivationData) RegisterActivationInGrainDirectoryAndValidate() bool {
	var success bool
	result, ok := a._shared.InternalGrainRuntime.GrainLocator.Register(a.Address)
	if ok {
		if a.Address.Matches(result) {
			success = true
		} else {
			a.ForwaringAddress = "RegisterActivationInGrainDirectoryAndValidate"
			a.DeactivationReason = "This grain hasbeen activated elsewhere"
			success = false
		}
	}

	if !success {
		a.lock.Lock()
		a.SetState(Invalid)
		a.lock.Unlock()
		a.UnregisterMessageTarget()
		a.DeactivationReason = "Failed to register activation in grain directory."
	}

	return success
}

func (a *ActivationData) UnregisterMessageTarget() {
	a._shared.InternalGrainRuntime.Catalog.UnregisterMessageTarget(a)
	//if a.GrainInstance != nil {
	//	a.SetGrainInstance(nil)
	//}
}

func (a *ActivationData) CallActivateSync() bool {
	//TODO lifecycle init
	var initSuccess bool = true
	if initSuccess {
		//if grainBase, ok := a.GrainInstance.(core.IGrainContext); ok {
		//	grainBase.OnActivateAsync()
		//}

		a.lock.Lock()
		if a.State == Activating {
			a.SetState(Valid)
		}
		a.lock.Unlock()
		return true
	}

	a.lock.Lock()
	a.SetState(FailedToActivate)
	a.DeactivationReason = "Failed to activate grain."
	a.lock.Unlock()

	//TODO
	//GetDeactivationCompletionSource().TrySetResult(true);

	if a.ForwaringAddress == "" {
		a._shared.InternalGrainRuntime.GrainLocator.Unregister(a.Address)
	}

	a.ScheduleOperation(CommandUnregisterFromCatalog{})
	a.lock.Lock()
	a.SetState(Invalid)
	a.lock.Unlock()

	return false
}

//func (a *ActivationData) SetGrainInstance(grainInstance any) {
//	a.GrainInstance = grainInstance
//}
