package Message

import "github.com/0990/orleans-learn/IDs"

type Message struct {
	Msg          string `json:"msg"`
	MethodName   string `json:"methodName"`
	ForwardCount int32
	TargetSilo   string
	TargetGrain  IDs.GrainId
	SendingSilo  string

	Deactivate bool //debug only
}
