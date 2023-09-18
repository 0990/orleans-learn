package core

import "github.com/0990/orleans-learn/core/Message"

type IGrainBase interface {
	OnActivateSync()
	OnDeactivateSync()

	OnMessage(msg Message.Message) //debug only
}
