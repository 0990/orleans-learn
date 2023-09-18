package main

import (
	"fmt"
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/catalog"
	"github.com/0990/orleans-learn/core"
	"github.com/0990/orleans-learn/core/Message"
	"github.com/0990/orleans-learn/pkg/util"
	"time"
)

func main() {
	catalog.RegisterNewGrain(func() any {
		return &Grain{}
	})

	cl := catalog.NewCatalog()
	mc := catalog.NewMessageCenter(cl)
	cl.Init(mc)

	mc.ReceiveMessage(Message.Message{
		Msg:          "1",
		MethodName:   "method1",
		ForwardCount: 0,
		TargetSilo:   "",
		TargetGrain: IDs.GrainId{
			GrainType: "hello",
			Key:       "hi",
		},
		SendingSilo: "",
	})

	mc.ReceiveMessage(Message.Message{
		Msg:          "two",
		MethodName:   "method1",
		ForwardCount: 0,
		TargetSilo:   "",
		TargetGrain: IDs.GrainId{
			GrainType: "hello",
			Key:       "hi",
		},
		SendingSilo: "",
		Deactivate:  true,
	})

	mc.ReceiveMessage(Message.Message{
		Msg:          "three",
		MethodName:   "method1",
		ForwardCount: 0,
		TargetSilo:   "",
		TargetGrain: IDs.GrainId{
			GrainType: "hello",
			Key:       "hi",
		},
		SendingSilo: "",
	})

	//a, ok := cl.GetOrCreateActivation(IDs.GrainId{
	//	GrainType: "hello",
	//	Key:       "hi",
	//})
	//
	//if !ok {
	//	fmt.Println("GetOrCreateActivation fail")
	//}
	//
	//a.ReceiveMessage(Message.Message{
	//	Msg: "1",
	//})
	//
	//a.ReceiveMessage(Message.Message{
	//	Msg: "2",
	//})
	//
	//time.Sleep(time.Second * 4)
	//a.Deactivate("test")
	//a.ReceiveMessage(Message.Message{
	//	Msg: "3",
	//})

	time.Sleep(time.Minute)
}

type Grain struct {
	core.IGrainBase
}

func (*Grain) OnActivateSync() {
	fmt.Println("OnActivateSync,", util.GetGoroutineID())
}

func (*Grain) OnDeactivateSync() {
	fmt.Println("OnDeactivateSync start")
	time.Sleep(time.Second * 3)
	fmt.Println("OnDeactivateSync end")
}

func (*Grain) OnMessage(msg Message.Message) {
	fmt.Println("OnMessage", msg.Msg)
}
