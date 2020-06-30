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

func TestMax(t *testing.T) {
	f := Max(0, 1)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
	f = Max(1, 0)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
	f = Max(1, 1)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
}

func TestMin(t *testing.T) {
	f := Min(2, 1)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
	f = Min(1, 2)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
	f = Min(1, 1)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
}

func TestClamp(t *testing.T) {
	f := Clamp(0, 1, 3)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
	f = Clamp(1, 1, 3)
	if f != 1 {
		t.Errorf("Expected %d, found %d.", 1, f)
	}
	f = Clamp(2, 1, 3)
	if f != 2 {
		t.Errorf("Expected %d, found %d.", 2, f)
	}
	f = Clamp(3, 1, 3)
	if f != 3 {
		t.Errorf("Expected %d, found %d.", 3, f)
	}
	f = Clamp(4, 1, 3)
	if f != 3 {
		t.Errorf("Expected %d, found %d.", 3, f)
	}
}
