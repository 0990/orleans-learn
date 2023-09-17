package activation

import (
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/core"
)

type GrainContextActivator struct {
}

func (a *GrainContextActivator) CreateInstance(address IDs.GrainAddress) core.IGrainContext {
	return nil
	//return catalog.NewActivationData(address)
}
