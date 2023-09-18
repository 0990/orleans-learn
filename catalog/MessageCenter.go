package catalog

import (
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/core/Message"
	"github.com/0990/orleans-learn/pkg/syncx"
	"log/slog"
)

// 本来在core/Messaging文件夹下的文件，但go不允许循环依赖，所以把这个文件移这了
type MessageCenter struct {
	_siloAddresses string
	_catalog       *Catalog
}

func NewMessageCenter(cl *Catalog) *MessageCenter {
	return &MessageCenter{
		_catalog: cl,
	}
}

func (mc *MessageCenter) ProcessRequestsToInvalidActivation(messages *syncx.List[Message.Message], oldAddress IDs.GrainAddress, forwardingAddress string, rejectMessage bool) {
	if rejectMessage {
		messages.Range(func(message Message.Message) {
			slog.Info("reject message", "msg", message)
		})
		return
	}

	messages.Range(func(message Message.Message) {
		mc.TryForwardRequest(message, oldAddress, forwardingAddress)
	})
}

func (mc *MessageCenter) TryForwardRequest(message Message.Message, oldAddress IDs.GrainAddress, forwardingAddress string) {
	mc.TryForwardMessage(message, forwardingAddress)
}

func (mc *MessageCenter) TryForwardMessage(message Message.Message, forwardingAddress string) {
	message.ForwardCount++
	mc.ResendMessageImpl(message, forwardingAddress)
}

func (mc *MessageCenter) ResendMessageImpl(message Message.Message, forwardingAddress string) {
	if forwardingAddress != "" {
		message.TargetSilo = forwardingAddress
		mc.SendMessage(message)
		return
	}
	message.TargetSilo = ""
	mc.AddressAndSendMessage(message)
}

func (mc *MessageCenter) AddressAndSendMessage(message Message.Message) {
	//TODO try get address from placementService
	//var messageAddressingTask = placementService.AddressMessage(message);
	mc.SendMessage(message)
}

func (mc *MessageCenter) SendMessage(message Message.Message) {
	message.SendingSilo = mc._siloAddresses

	mc.ReceiveMessage(message)
}

func (mc *MessageCenter) ReceiveMessage(message Message.Message) {
	targetActivation, ok := mc._catalog.GetOrCreateActivation(message.TargetGrain)
	if !ok {
		mc.ProcessMessageToNonExistentActivation(message)
		return
	}
	targetActivation.ReceiveMessage(message)
}

func (mc *MessageCenter) ProcessMessageToNonExistentActivation(msg Message.Message) {
	nonExistentActivation := IDs.GrainAddress{
		GrainId:      msg.TargetGrain,
		ActivationId: IDs.ActivationId{},
	}

	l := syncx.NewList[Message.Message]()
	l.PushBack(msg)
	mc.ProcessRequestsToInvalidActivation(l, nonExistentActivation, "", false)
}
