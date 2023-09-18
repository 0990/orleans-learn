package IDs

type GrainAddress struct {
	GrainId      GrainId
	ActivationId ActivationId
}

func NewGrainAddress(grainId GrainId, activationId ActivationId) GrainAddress {
	return GrainAddress{
		GrainId:      grainId,
		ActivationId: activationId,
	}
}

func (ga *GrainAddress) Matches(other GrainAddress) bool {
	return ga.GrainId == other.GrainId && ga.ActivationId == other.ActivationId
}
