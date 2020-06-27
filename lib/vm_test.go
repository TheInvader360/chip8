package lib

import (
	"testing"
)

func TestFetchOpcode(t *testing.T) {
	vm := NewChip8()
	vm.Mem[0x0000] = 0xA2
	vm.Mem[0x0001] = 0xF0
	vm.Mem[0x0002] = 0xC5
	vm.Mem[0x0003] = 0x02
	vm.Mem[0x0004] = 0x0F
	vm.Mem[0x0005] = 0xF0
	vm.Mem[0x0006] = 0xAB
	vm.Mem[0x0007] = 0x50
	vm.Mem[0x0008] = 0x20
	vm.Mem[0x0009] = 0xFF
	for i := 0; i < 5; i++ {
		vm.pc = uint16(i * 2)
		oc := vm.fetchOpcode()
		e := uint16(0x0000)
		switch i {
		case 0:
			e = uint16(0xA2F0)
		case 1:
			e = uint16(0xC502)
		case 2:
			e = uint16(0x0FF0)
		case 3:
			e = uint16(0xAB50)
		case 4:
			e = uint16(0x20FF)
		}
		if oc != e {
			t.Errorf("Expected %X, found %X.", e, oc)
		}
	}
}

func TestDecodeOpcode(t *testing.T) {
	vm := NewChip8()
	m := map[uint16]uint16{
		uint16(0x0123): uint16(0x0000),
		uint16(0x00E0): uint16(0x00E0),
		uint16(0x00EE): uint16(0x00EE),
		uint16(0x1AB3): uint16(0x1000),
		uint16(0x205F): uint16(0x2000),
		uint16(0x3303): uint16(0x3000),
		uint16(0x45B0): uint16(0x4000),
		uint16(0x55FF): uint16(0x5000),
		uint16(0x6FA5): uint16(0x6000),
		uint16(0x7B09): uint16(0x7000),
		uint16(0x8120): uint16(0x8000),
		uint16(0x8341): uint16(0x8001),
		uint16(0x8562): uint16(0x8002),
		uint16(0x8783): uint16(0x8003),
		uint16(0x89A4): uint16(0x8004),
		uint16(0x8BC5): uint16(0x8005),
		uint16(0x8DE6): uint16(0x8006),
		uint16(0x8F07): uint16(0x8007),
		uint16(0x80FE): uint16(0x800E),
		uint16(0x9A90): uint16(0x9000),
		uint16(0xABCD): uint16(0xA000),
		uint16(0xB963): uint16(0xB000),
		uint16(0xC5B0): uint16(0xC000),
		uint16(0xD8AF): uint16(0xD000),
		uint16(0xEF9E): uint16(0xE09E),
		uint16(0xE5A1): uint16(0xE0A1),
		uint16(0xF507): uint16(0xF007),
		uint16(0xFA0A): uint16(0xF00A),
		uint16(0xF915): uint16(0xF015),
		uint16(0xF318): uint16(0xF018),
		uint16(0xF81E): uint16(0xF01E),
		uint16(0xFD29): uint16(0xF029),
		uint16(0xF233): uint16(0xF033),
		uint16(0xF155): uint16(0xF055),
		uint16(0xFF65): uint16(0xF065),
	}
	for in, out := range m {
		vm.oc = in
		d := vm.decodeOpcode()
		e := out
		if d != e {
			t.Errorf("Expected %X, found %X.", e, d)
		}
	}
}

func TestUpdateTimers(t *testing.T) {
	vm := NewChip8()
	vm.dt = 2
	vm.st = 1
	vm.updateTimers()
	if vm.dt != 1 {
		t.Errorf("Expected %X, found %X.", 1, vm.dt)
	}
	if vm.st != 0 {
		t.Errorf("Expected %X, found %X.", 0, vm.st)
	}
	vm.updateTimers()
	if vm.dt != 0 {
		t.Errorf("Expected %X, found %X.", 0, vm.dt)
	}
	if vm.st != 0 {
		t.Errorf("Expected %X, found %X.", 0, vm.st)
	}
	vm.updateTimers()
	if vm.dt != 0 {
		t.Errorf("Expected %X, found %X.", 0, vm.dt)
	}
}
