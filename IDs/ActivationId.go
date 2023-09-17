package IDs

import "github.com/google/uuid"

type ActivationId struct {
	Key string
}

func NewActivationId() ActivationId {
	uuid := uuid.New().String()
	return ActivationId{
		Key: uuid,
	}
}
