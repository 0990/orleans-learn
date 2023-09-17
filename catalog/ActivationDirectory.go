package catalog

import (
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/core"
	"github.com/0990/orleans-learn/pkg/syncx"
	"sync/atomic"
)

type ActivationDirectory struct {
	activations      *syncx.SyncMap[IDs.GrainId, core.IGrainContext]
	_activationCount int32
}

func NewActivationDirectory() *ActivationDirectory {
	return &ActivationDirectory{
		activations: syncx.NewSyncMap[IDs.GrainId, core.IGrainContext](),
	}
}

func (ad *ActivationDirectory) FindTarget(key IDs.GrainId) (core.IGrainContext, bool) {
	data, ok := ad.activations.Load(key)
	return data, ok
}

func (ad *ActivationDirectory) RecordNewTarget(target core.IGrainContext) {
	ok := ad.activations.Store(target.GrainId(), target)
	if ok {
		atomic.AddInt32(&ad._activationCount, 1)
	}
}

func (ad *ActivationDirectory) RemoveTarget(target core.IGrainContext) {
	ok := ad.activations.Delete(target.GrainId())
	if ok {
		atomic.AddInt32(&ad._activationCount, -1)
	}
}
