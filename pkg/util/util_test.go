package util

import (
	"fmt"
	"testing"
)

func TestGetGoroutineID(t *testing.T) {
	id := GetGoroutineID()
	if id <= 0 {
		t.Fail()
	}
	fmt.Println(id)
}
