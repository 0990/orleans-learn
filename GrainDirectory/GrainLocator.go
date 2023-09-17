package GrainDirectory

import "github.com/0990/orleans-learn/IDs"

type GrainLocator struct {
	addresses map[IDs.GrainId]IDs.GrainAddress
}

func NewGrainLocator() *GrainLocator {
	return &GrainLocator{
		addresses: make(map[IDs.GrainId]IDs.GrainAddress),
	}
}

func (gl *GrainLocator) Register(address IDs.GrainAddress) (IDs.GrainAddress, bool) {
	result, ok := gl.addresses[address.GrainId()]
	if ok {
		return result, false
	}

	gl.addresses[address.GrainId()] = address
	return address, true
}

func (gl *GrainLocator) Unregister(address IDs.GrainAddress) {
	delete(gl.addresses, address.GrainId())
}
