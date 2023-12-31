package catalog

import (
	"github.com/0990/orleans-learn/GrainDirectory"
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/core"
	"sync"
)

type Catalog struct {
	activations *ActivationDirectory
	//grainActivator *activation.GrainContextActivator
	lock    sync.Mutex
	_shared GrainTypeShardContext
}

func NewCatalog() *Catalog {
	c := &Catalog{
		activations: NewActivationDirectory(),
		//grainActivator: activation.NewGrainContextActivator(),
	}

	shared := GrainTypeShardContext{
		InternalGrainRuntime: &InternalGrainRuntime{
			GrainLocator: GrainDirectory.NewGrainLocator(),
			Catalog:      c,
		},
	}
	c._shared = shared
	return c
}

func (c *Catalog) Init(mc *MessageCenter) {
	c._shared.InternalGrainRuntime.MessageCenter = mc
}

func (c *Catalog) GetOrCreateActivation(grainId IDs.GrainId) (core.IGrainContext, bool) {
	result, ok := c.TryGetGrainContext(grainId)
	if ok {
		return result, true
	}

	c.lock.Lock()
	result, ok = c.TryGetGrainContext(grainId)
	if ok {
		return result, true
	}
	//create grain
	var address = IDs.NewGrainAddress(grainId, IDs.NewActivationId())
	//result = c.grainActivator.CreateInstance(address)

	ad := NewActivationData(address, c._shared)
	c.RegisterMessageTarget(ad)
	c.lock.Unlock()

	if ad == nil {
		return nil, false
	}

	ad.SetGrainInstance(CreateNewGrainInstance())
	ad.Activate()
	return ad, true
}

func (c *Catalog) TryGetGrainContext(grainId IDs.GrainId) (core.IGrainContext, bool) {
	data, ok := c.activations.FindTarget(grainId)
	return data, ok
}

func (c *Catalog) RegisterMessageTarget(activation core.IGrainContext) {
	c.activations.RecordNewTarget(activation)
}

func (c *Catalog) UnregisterMessageTarget(activation core.IGrainContext) {
	c.activations.RemoveTarget(activation)
}
