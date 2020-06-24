package main

import "testing"

func TestBoolToByte(t *testing.T) {
	expected := byte(0)
	found := boolToByte(false)
	if found != expected {
		t.Errorf("Expected %d, found %d.", expected, found)
	}
	expected = byte(1)
	found = boolToByte(true)
	if found != expected {
		t.Errorf("Expected %d, found %d.", expected, found)
	}
}
