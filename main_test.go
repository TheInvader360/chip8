package main

import "testing"

func TestFetchOpcode(t *testing.T) {
	memory = [4096]byte{}
	memory[0x000] = 0xA2
	memory[0x001] = 0xF0
	memory[0x002] = 0xC5
	memory[0x003] = 0x02
	memory[0x004] = 0x0F
	memory[0x005] = 0xF0
	memory[0x006] = 0xAB
	memory[0x007] = 0x50
	memory[0x008] = 0x20
	memory[0x009] = 0xFF

	pc = 0x000
	opcode := fetchOpcode()
	expected := uint16(0xA2F0)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x002
	opcode = fetchOpcode()
	expected = uint16(0xC502)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x004
	opcode = fetchOpcode()
	expected = uint16(0x0FF0)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x006
	opcode = fetchOpcode()
	expected = uint16(0xAB50)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x008
	opcode = fetchOpcode()
	expected = uint16(0x20FF)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}
}

func TestDecodeOpcode(t *testing.T) {
	opcode = 0x0123
	decoded := decodeOpcode()
	expected := uint16(0x0000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x00E0
	decoded = decodeOpcode()
	expected = uint16(0x00E0)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x00EE
	decoded = decodeOpcode()
	expected = uint16(0x00EE)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x1AB5
	decoded = decodeOpcode()
	expected = uint16(0x1000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x205F
	decoded = decodeOpcode()
	expected = uint16(0x2000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x3303
	decoded = decodeOpcode()
	expected = uint16(0x3000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x45B0
	decoded = decodeOpcode()
	expected = uint16(0x4000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x55FF
	decoded = decodeOpcode()
	expected = uint16(0x5000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x6FA5
	decoded = decodeOpcode()
	expected = uint16(0x6000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x7B09
	decoded = decodeOpcode()
	expected = uint16(0x7000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x8120
	decoded = decodeOpcode()
	expected = uint16(0x8000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x8341
	decoded = decodeOpcode()
	expected = uint16(0x8001)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x8562
	decoded = decodeOpcode()
	expected = uint16(0x8002)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x8783
	decoded = decodeOpcode()
	expected = uint16(0x8003)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x89A4
	decoded = decodeOpcode()
	expected = uint16(0x8004)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x8BC5
	decoded = decodeOpcode()
	expected = uint16(0x8005)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x8DE6
	decoded = decodeOpcode()
	expected = uint16(0x8006)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x8F07
	decoded = decodeOpcode()
	expected = uint16(0x8007)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0x80FE
	decoded = decodeOpcode()
	expected = uint16(0x800E)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0x9A90
	decoded = decodeOpcode()
	expected = uint16(0x9000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0xABCD
	decoded = decodeOpcode()
	expected = uint16(0xA000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0xB963
	decoded = decodeOpcode()
	expected = uint16(0xB000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0xC5B0
	decoded = decodeOpcode()
	expected = uint16(0xC000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0xD8AF
	decoded = decodeOpcode()
	expected = uint16(0xD000)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0xEF9E
	decoded = decodeOpcode()
	expected = uint16(0xE09E)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xE5A1
	decoded = decodeOpcode()
	expected = uint16(0xE0A1)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}

	opcode = 0xF507
	decoded = decodeOpcode()
	expected = uint16(0xF007)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xFA0A
	decoded = decodeOpcode()
	expected = uint16(0xF00A)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xF915
	decoded = decodeOpcode()
	expected = uint16(0xF015)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xF318
	decoded = decodeOpcode()
	expected = uint16(0xF018)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xF81E
	decoded = decodeOpcode()
	expected = uint16(0xF01E)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xFD29
	decoded = decodeOpcode()
	expected = uint16(0xF029)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xF233
	decoded = decodeOpcode()
	expected = uint16(0xF033)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xF155
	decoded = decodeOpcode()
	expected = uint16(0xF055)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
	opcode = 0xFF65
	decoded = decodeOpcode()
	expected = uint16(0xF065)
	if decoded != expected {
		t.Errorf("Expected %X, found %X.", expected, decoded)
	}
}

func TestUpdateTimers(t *testing.T) {
	delayTimer = 4
	soundTimer = 2

	updateTimers()
	if delayTimer != 3 {
		t.Errorf("Expected %X, found %X.", 3, delayTimer)
	}
	if soundTimer != 1 {
		t.Errorf("Expected %X, found %X.", 1, soundTimer)
	}

	updateTimers()
	if delayTimer != 2 {
		t.Errorf("Expected %X, found %X.", 2, delayTimer)
	}
	if soundTimer != 0 {
		t.Errorf("Expected %X, found %X.", 0, soundTimer)
	}

	updateTimers()
	if delayTimer != 1 {
		t.Errorf("Expected %X, found %X.", 1, delayTimer)
	}
	if soundTimer != 0 {
		t.Errorf("Expected %X, found %X.", 0, soundTimer)
	}

	updateTimers()
	if delayTimer != 0 {
		t.Errorf("Expected %X, found %X.", 0, delayTimer)
	}
	if soundTimer != 0 {
		t.Errorf("Expected %X, found %X.", 0, soundTimer)
	}
}

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

func TestExec0NNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x0B99
	expected := "exec0NNN 0x0B99"
	found := exec0NNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec00E0(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x00E0
	expected := "exec00E0 0x00E0"
	found := exec00E0()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec00EE(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x00EE
	expected := "exec00EE 0x00EE"
	found := exec00EE()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec1NNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x1A9F
	expected := "exec1NNN 0x1A9F"
	found := exec1NNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec2NNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x2F08
	expected := "exec2NNN 0x2F08"
	found := exec2NNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec3XNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x3A1D
	expected := "exec3XNN 0x3A1D"
	found := exec3XNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec4XNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x4247
	expected := "exec4XNN 0x4247"
	found := exec4XNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec5XY0(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x5190
	expected := "exec5XY0 0x5190"
	found := exec5XY0()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec6XNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x6EFD
	expected := "exec6XNN 0x6EFD"
	found := exec6XNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec7XNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x7015
	expected := "exec7XNN 0x7015"
	found := exec7XNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY0(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x89B0
	expected := "exec8XY0 0x89B0"
	found := exec8XY0()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY1(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x8DE1
	expected := "exec8XY1 0x8DE1"
	found := exec8XY1()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY2(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x83F2
	expected := "exec8XY2 0x83F2"
	found := exec8XY2()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY3(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x8AC3
	expected := "exec8XY3 0x8AC3"
	found := exec8XY3()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY4(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x8474
	expected := "exec8XY4 0x8474"
	found := exec8XY4()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY5(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x89A5
	expected := "exec8XY5 0x89A5"
	found := exec8XY5()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY6(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x8206
	expected := "exec8XY6 0x8206"
	found := exec8XY6()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XY7(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x8197
	expected := "exec8XY7 0x8197"
	found := exec8XY7()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec8XYE(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x8EEE
	expected := "exec8XYE 0x8EEE"
	found := exec8XYE()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExec9XY0(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0x9730
	expected := "exec9XY0 0x9730"
	found := exec9XY0()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecANNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xA259
	expected := "execANNN 0xA259"
	found := execANNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecBNNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xBAF2
	expected := "execBNNN 0xBAF2"
	found := execBNNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecCXNN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xC4CE
	expected := "execCXNN 0xC4CE"
	found := execCXNN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecDXYN(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xD81C
	expected := "execDXYN 0xD81C"
	found := execDXYN()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecEX9E(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xE39E
	expected := "execEX9E 0xE39E"
	found := execEX9E()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecEXA1(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xE2A1
	expected := "execEXA1 0xE2A1"
	found := execEXA1()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX07(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xF807
	expected := "execFX07 0xF807"
	found := execFX07()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX0A(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xFC0A
	expected := "execFX0A 0xFC0A"
	found := execFX0A()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX15(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xF815
	expected := "execFX15 0xF815"
	found := execFX15()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX18(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xF118
	expected := "execFX18 0xF118"
	found := execFX18()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX1E(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xFF1E
	expected := "execFX1E 0xFF1E"
	found := execFX1E()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX29(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xF329
	expected := "execFX29 0xF329"
	found := execFX29()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX33(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xF533
	expected := "execFX33 0xF533"
	found := execFX33()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX55(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xFC55
	expected := "execFX55 0xFC55"
	found := execFX55()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}

func TestExecFX65(t *testing.T) {
	//TODO rewrite test for real implementation!
	opcode = 0xFF65
	expected := "execFX65 0xFF65"
	found := execFX65()
	if found != expected {
		t.Errorf("Expected %s, found %s.", expected, found)
	}
}
