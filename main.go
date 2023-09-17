package main

import (
	"fmt"
	"github.com/0990/orleans-learn/IDs"
	"github.com/0990/orleans-learn/catalog"
	"github.com/0990/orleans-learn/core/Message"
	"time"
)

func main() {
	cl := catalog.NewCatalog()
	a, ok := cl.GetOrCreateActivation(IDs.GrainId{
		GrainType: "hello",
		Key:       "hi",
	})
	if !ok {
		fmt.Println("GetOrCreateActivation fail")
	}

	a.ReceiveMessage(Message.Message{
		Msg: "1",
	})

	a.ReceiveMessage(Message.Message{
		Msg: "2",
	})

	time.Sleep(time.Minute)
}
