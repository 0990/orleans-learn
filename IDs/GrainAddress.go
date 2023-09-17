package IDs

type GrainAddress struct {
	_grainId      GrainId
	_activationId ActivationId
}

func NewGrainAddress(grainId GrainId, activationId ActivationId) GrainAddress {
	return GrainAddress{
		_grainId:      grainId,
		_activationId: activationId,
	}
}

func (ga *GrainAddress) GrainId() GrainId {
	return ga._grainId
}

func (ga *GrainAddress) ActivationId() ActivationId {
	return ga._activationId
}

func (ga *GrainAddress) Matches(other GrainAddress) bool {
	return ga._grainId == other._grainId && ga._activationId == other._activationId
}
