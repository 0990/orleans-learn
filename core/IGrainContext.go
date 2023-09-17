package core

import (
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/core/Message"
)

type IGrainContext interface {
	GrainId() IDs.GrainId
	ActivationId() IDs.ActivationId

	ReceiveMessage(message Message.Message)
	Activate()
}
