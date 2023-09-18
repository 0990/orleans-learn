package catalog

import (
	"github.com/0990/orleans-learn/core"
	"log/slog"
)

func (a *ActivationData) ActivateSync() {
	slog.Info("start Activate")
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

func (a *ActivationData) CallActivateSync() bool {
	//TODO lifecycle init
	var initSuccess bool = true
	if initSuccess {
		if grainBase, ok := a.GrainInstance.(core.IGrainBase); ok {
			grainBase.OnActivateSync()
		}

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
