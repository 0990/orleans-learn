package util

import "github.com/petermattis/goid"

func GetGoroutineID() int64 {
	return goid.Get()
}
