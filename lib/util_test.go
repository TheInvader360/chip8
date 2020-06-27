package lib

import (
	"testing"
)

func TestBoolToByte(t *testing.T) {
	e := byte(0)
	f := BoolToByte(false)
	if f != e {
		t.Errorf("Expected %d, found %d.", e, f)
	}
	e = byte(1)
	f = BoolToByte(true)
	if f != e {
		t.Errorf("Expected %d, found %d.", e, f)
	}
}
