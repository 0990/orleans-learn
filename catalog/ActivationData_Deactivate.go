package catalog

import "github.com/0990/orleans-learn/core"

func (a *ActivationData) StartDeactivating(reason string) bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.State == Deactivating || a.State == Invalid || a.State == FailedToActivate {
		return false
	}

	if a.State == Activating || a.State == Create {
		panic("Calling DeactivateOnIdle from within OnActivateAsync is not supported")
	}

	if a.DeactivationReason == "" {
		a.DeactivationReason = reason
	}

	a.SetState(Deactivating)
	if !a.IsCurrentlyExecuting {
		//StopAllTimers()
	}

	return true
}

func (a *ActivationData) FinishDeactivating() {
	a.CallGrainDeactivate()

	a.lock.Lock()
	a.SetState(Invalid)
	a.lock.Unlock()

	a.UnregisterMessageTarget()
	//DisposeAsync()
	//GetDeactivationCompletionSource().TrySetResult(true)
	a._workSignal.Signal()

	a._shared.InternalGrainRuntime.GrainLocator.Unregister(a.Address)
}

func (a *ActivationData) CallGrainDeactivate() {
	if a.State == Deactivating {
		if grainBase, ok := a.GrainInstance.(core.IGrainBase); ok {
			grainBase.OnDeactivateSync()
		}
	}
}
