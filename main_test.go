package main

import (
	"testing"
)

func TestFetchOpcode(t *testing.T) {
	mem = [memS]byte{}
	mem[0x0000] = 0xA2
	mem[0x0001] = 0xF0
	mem[0x0002] = 0xC5
	mem[0x0003] = 0x02
	mem[0x0004] = 0x0F
	mem[0x0005] = 0xF0
	mem[0x0006] = 0xAB
	mem[0x0007] = 0x50
	mem[0x0008] = 0x20
	mem[0x0009] = 0xFF

	pc = 0x0000
	opcode := fetchOpcode()
	expected := uint16(0xA2F0)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x0002
	opcode = fetchOpcode()
	expected = uint16(0xC502)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x0004
	opcode = fetchOpcode()
	expected = uint16(0x0FF0)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x0006
	opcode = fetchOpcode()
	expected = uint16(0xAB50)
	if opcode != expected {
		t.Errorf("Expected %X, found %X.", expected, opcode)
	}

	pc = 0x0008
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
	//not implemented
	opcode = 0x00E0
	pc = 0x0000
	ertn := "exec0NNN 0x00E0: pc=0x0000 (not implemented)"
	frtn := exec0NNN()
	epc := uint16(0x0000)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec00E0(t *testing.T) {
	//disp_clear()
	opcode = 0x00E0
	pc = 0x0000
	gfx = [gfxS]byte{}
	gfx[8] = 1
	gfx[128] = 1
	gfx[278] = 1
	gfx[1080] = 1
	gfx[2000] = 1
	ertn := "exec00E0 0x00E0: pc=0x0002 gfx={cleared} render=true"
	frtn := exec00E0()
	egfx := [gfxS]byte{}
	fgfx := gfx
	epc := uint16(0x0002)
	fpc := pc
	erender := true
	frender := render
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fgfx != egfx {
		t.Errorf("Expected %v, found %v.", egfx, fgfx)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	if frender != erender {
		t.Errorf("Expected %t, found %t.", erender, frender)
	}
}

func TestExec00EE(t *testing.T) {
	//TODO return
	opcode = 0x00EE
}

func TestExec1NNN(t *testing.T) {
	//goto nnn
	opcode = 0x1A9F
	ertn := "exec1NNN 0x1A9F: pc=0x0A9F (goto)"
	frtn := exec1NNN()
	epc := uint16(0x0A9F)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec2NNN(t *testing.T) {
	//TODO *(0xnnn)()
	opcode = 0x2F08
}

func TestExec3XNN(t *testing.T) {
	//TODO if(vx==nn)
	opcode = 0x3A1D
}

func TestExec4XNN(t *testing.T) {
	//TODO if(vx!=nn)
	opcode = 0x4247
}

func TestExec5XY0(t *testing.T) {
	//TODO if(vx==vy)
	opcode = 0x5190
}

func TestExec6XNN(t *testing.T) {
	//vx=nn
	opcode = 0x6EFD
	pc = 0x0000
	ertn := "exec6XNN 0x6EFD: pc=0x0002 v[E]=FD"
	frtn := exec6XNN()
	evx := byte(0xFD)
	fvx := v[0xE]
	epc := uint16(0x0002)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fvx != evx {
		t.Errorf("Expected 0x%02X, found 0x%02X.", evx, fvx)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec7XNN(t *testing.T) {
	//vx+=nn
	opcode = 0x7315
	pc = 0x0000
	v[0x3] = 0xA1
	ertn := "exec7XNN 0x7315: pc=0x0002 v[3]=B6"
	frtn := exec7XNN()
	evx := byte(0xB6)
	fvx := v[0x3]
	epc := uint16(0x0002)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fvx != evx {
		t.Errorf("Expected 0x%02X, found 0x%02X.", evx, fvx)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec8XY0(t *testing.T) {
	//TODO vx=vy
	opcode = 0x89B0
}

func TestExec8XY1(t *testing.T) {
	//TODO vx=vx|vy
	opcode = 0x8DE1
}

func TestExec8XY2(t *testing.T) {
	//TODO vx=vx&vy
	opcode = 0x83F2
}

func TestExec8XY3(t *testing.T) {
	//TODO vx=vx^vy
	opcode = 0x8AC3
}

func TestExec8XY4(t *testing.T) {
	//TODO vx+=vy
	opcode = 0x8474
}

func TestExec8XY5(t *testing.T) {
	//TODO vx-=vy
	opcode = 0x89A5
}

func TestExec8XY6(t *testing.T) {
	//TODO vx>>=1
	opcode = 0x8206
}

func TestExec8XY7(t *testing.T) {
	//TODO vx=vy-vx
	opcode = 0x8197
}

func TestExec8XYE(t *testing.T) {
	//TODO vx<<=1
	opcode = 0x8EEE
}

func TestExec9XY0(t *testing.T) {
	//TODO if(vx!=vy)
	opcode = 0x9730
}

func TestExecANNN(t *testing.T) {
	//i=nnn
	opcode = 0xA259
	pc = 0x0000
	ertn := "execANNN 0xA259: pc=0x0002 i=0x0259"
	frtn := execANNN()
	ei := uint16(0x0259)
	fi := i
	epc := uint16(0x0002)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fi != ei {
		t.Errorf("Expected 0x%04X, found 0x%04X.", ei, fi)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecBNNN(t *testing.T) {
	//TODO pc=v0+nnn
	opcode = 0xBAF2
}

func TestExecCXNN(t *testing.T) {
	//TODO vx=rand()&nn
	opcode = 0xC4CE
}

func TestExecDXYN(t *testing.T) {
	//draw(x,y,n)

	//draw a 3x3 sprite at 5,2 (sprite: diagonal line topleft to bottomright)
	opcode = 0xD793
	pc = 0x0000
	gfx = [gfxS]byte{}
	mem = [memS]byte{}
	mem[0x0101] = 0x80
	mem[0x0102] = 0x40
	mem[0x0103] = 0x20
	i = 0x0101
	v[7] = 5
	v[9] = 2
	ertn := "execDXYN 0xD793: pc=0x0002 gfx={updated} render=true"
	frtn := execDXYN()
	egfx := [gfxS]byte{}
	egfx[133] = 1
	egfx[198] = 1
	egfx[263] = 1
	fgfx := gfx
	evf := byte(0)
	fvf := v[0xF]
	epc := uint16(0x0002)
	fpc := pc
	erender := true
	frender := render
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fgfx != egfx {
		t.Errorf("Expected %v, found %v.", egfx, fgfx)
	}
	if fvf != evf {
		t.Errorf("Expected %v, found %v.", evf, fvf)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	if frender != erender {
		t.Errorf("Expected %t, found %t.", erender, frender)
	}

	//draw the sprite again at 8,5 (no pixels erased)
	v[7] = 8
	v[9] = 5
	execDXYN()
	evf = byte(0)
	fvf = v[0xF]
	if fvf != evf {
		t.Errorf("Expected %v, found %v.", evf, fvf)
	}

	//draw the sprite again at 10,7 (pixel erased)
	v[7] = 10
	v[9] = 7
	execDXYN()
	evf = byte(1)
	fvf = v[0xF]
	if fvf != evf {
		t.Errorf("Expected %v, found %v.", evf, fvf)
	}

	//clear gfx then draw the sprite at 62,0 (exceed horizontal bounds)
	gfx = [gfxS]byte{}
	v[7] = 62
	v[9] = 0
	execDXYN()
	egfx = [gfxS]byte{}
	//egfx[62] = 1
	//egfx[127] = 1
	//egfx[128] = 1
	//fgfx = gfx
	//TODO: confirm expected behaviour then finish implementing test!

	//clear gfx then draw the sprite at 0,31 (exceed vertical bounds)
	gfx = [gfxS]byte{}
	v[7] = 0
	v[9] = 31
	execDXYN()
	egfx = [gfxS]byte{}
	//egfx[?] = 1
	//egfx[?] = 1
	//egfx[?] = 1
	//fgfx = gfx
	//TODO: confirm expected behaviour then finish implementing test!
}

func TestExecEX9E(t *testing.T) {
	//TODO if(key()==vx)
	opcode = 0xE39E
}

func TestExecEXA1(t *testing.T) {
	//TODO if(key()!=vx)
	opcode = 0xE2A1
}

func TestExecFX07(t *testing.T) {
	//TODO vx=get_delay()
	opcode = 0xF807
}

func TestExecFX0A(t *testing.T) {
	//TODO vx=get_key()
	opcode = 0xFC0A
}

func TestExecFX15(t *testing.T) {
	//TODO delay_timer(vx)
	opcode = 0xF815
}

func TestExecFX18(t *testing.T) {
	//TODO sound_timer(vx)
	opcode = 0xF118
}

func TestExecFX1E(t *testing.T) {
	//TODO i+=vx
	opcode = 0xFF1E
}

func TestExecFX29(t *testing.T) {
	//TODO i=sprite_addr[vx]
	opcode = 0xF329
}

func TestExecFX33(t *testing.T) {
	//TODO set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
	opcode = 0xF533
}

func TestExecFX55(t *testing.T) {
	//TODO reg_dump(vx,&i)
	opcode = 0xFC55
}

func TestExecFX65(t *testing.T) {
	//TODO reg_load(vx,&i)
	opcode = 0xFF65
}
