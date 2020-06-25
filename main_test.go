package main

import (
	"testing"
)

func reset() {
	opcode = 0
	mem = [memS]byte{}
	v = [vS]byte{}
	i = 0
	pc = 0x200
	gfx = [gfxS]byte{}
	delayTimer = 0
	soundTimer = 0
	stack = [sS]uint16{}
	sp = 0
	keys = [kS]byte{}
	render = true
	loaded = false
}

func TestFetchOpcode(t *testing.T) {
	reset()
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
	for i := 0; i < 5; i++ {
		pc = uint16(i * 2)
		opcode := fetchOpcode()
		expected := uint16(0x0000)
		switch i {
		case 0:
			expected = uint16(0xA2F0)
		case 1:
			expected = uint16(0xC502)
		case 2:
			expected = uint16(0x0FF0)
		case 3:
			expected = uint16(0xAB50)
		case 4:
			expected = uint16(0x20FF)
		}
		if opcode != expected {
			t.Errorf("Expected %X, found %X.", expected, opcode)
		}
	}
}

func TestDecodeOpcode(t *testing.T) {
	reset()
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
	for k, v := range m {
		opcode = k
		decoded := decodeOpcode()
		expected := v
		if decoded != expected {
			t.Errorf("Expected %X, found %X.", expected, decoded)
		}
	}
}

func TestUpdateTimers(t *testing.T) {
	reset()
	delayTimer = 2
	soundTimer = 1
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
}

func TestBoolToByte(t *testing.T) {
	reset()
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
	reset()
	opcode = 0x00E0
	ertn := "exec0NNN 0x00E0: pc=0x0200 (not implemented)"
	frtn := exec0NNN()
	epc := uint16(0x0200)
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
	reset()
	opcode = 0x00E0
	gfx[8] = 1
	gfx[278] = 1
	gfx[1080] = 1
	gfx[2000] = 1
	ertn := "exec00E0 0x00E0: pc=0x0202 gfx={cleared} render=true"
	frtn := exec00E0()
	egfx := [gfxS]byte{}
	fgfx := gfx
	epc := uint16(0x0202)
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
	reset()
	opcode = 0x00EE
}

func TestExec1NNN(t *testing.T) {
	//goto nnn
	reset()
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
	reset()
	opcode = 0x2F08
}

func TestExec3XNN(t *testing.T) {
	//if(vx==nn) skip next instruction
	reset()
	opcode = 0x3A1D
	v[0xA] = 0x1D
	ertn := "exec3XNN 0x3A1D: pc=0x0204 {skip=true}"
	frtn := exec3XNN()
	epc := uint16(0x0204)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	v[0xA] = 0xFD
	ertn = "exec3XNN 0x3A1D: pc=0x0206 {skip=false}"
	frtn = exec3XNN()
	epc = uint16(0x0206)
	fpc = pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec4XNN(t *testing.T) {
	//TODO if(vx!=nn)
	reset()
	opcode = 0x4247
}

func TestExec5XY0(t *testing.T) {
	//TODO if(vx==vy)
	reset()
	opcode = 0x5190
}

func TestExec6XNN(t *testing.T) {
	//vx=nn
	reset()
	opcode = 0x6EFD
	ertn := "exec6XNN 0x6EFD: pc=0x0202 v[E]=FD"
	frtn := exec6XNN()
	evx := byte(0xFD)
	fvx := v[0xE]
	epc := uint16(0x0202)
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
	reset()
	opcode = 0x7315
	v[0x3] = 0xA1
	ertn := "exec7XNN 0x7315: pc=0x0202 v[3]=B6"
	frtn := exec7XNN()
	evx := byte(0xB6)
	fvx := v[0x3]
	epc := uint16(0x0202)
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
	reset()
	opcode = 0x89B0
}

func TestExec8XY1(t *testing.T) {
	//TODO vx=vx|vy
	reset()
	opcode = 0x8DE1
}

func TestExec8XY2(t *testing.T) {
	//TODO vx=vx&vy
	reset()
	opcode = 0x83F2
}

func TestExec8XY3(t *testing.T) {
	//TODO vx=vx^vy
	reset()
	opcode = 0x8AC3
}

func TestExec8XY4(t *testing.T) {
	//TODO vx+=vy
	reset()
	opcode = 0x8474
}

func TestExec8XY5(t *testing.T) {
	//TODO vx-=vy
	reset()
	opcode = 0x89A5
}

func TestExec8XY6(t *testing.T) {
	//TODO vx>>=1
	reset()
	opcode = 0x8206
}

func TestExec8XY7(t *testing.T) {
	//TODO vx=vy-vx
	reset()
	opcode = 0x8197
}

func TestExec8XYE(t *testing.T) {
	//TODO vx<<=1
	reset()
	opcode = 0x8EEE
}

func TestExec9XY0(t *testing.T) {
	//TODO if(vx!=vy)
	reset()
	opcode = 0x9730
}

func TestExecANNN(t *testing.T) {
	//i=nnn
	reset()
	opcode = 0xA259
	ertn := "execANNN 0xA259: pc=0x0202 i=0x0259"
	frtn := execANNN()
	ei := uint16(0x0259)
	fi := i
	epc := uint16(0x0202)
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
	reset()
	opcode = 0xBAF2
}

func TestExecCXNN(t *testing.T) {
	//vx=rand()&nn
	reset()
	opcode = 0xC400
	ertn := "execCXNN 0xC400: pc=0x0202 v[4]=00"
	frtn := execCXNN()
	evx := byte(0x00)
	fvx := v[0x4]
	epc := uint16(0x0202)
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
	//TODO: reliable tests where nn is not zero...
}

func TestExecDXYN(t *testing.T) {
	//draw(x,y,n)
	reset()
	//draw a 3x3 sprite at 5,2 (sprite: diagonal line topleft to bottomright)
	opcode = 0xD793
	mem[0x0101] = 0x80
	mem[0x0102] = 0x40
	mem[0x0103] = 0x20
	i = 0x0101
	v[7] = 5
	v[9] = 2
	ertn := "execDXYN 0xD793: pc=0x0202 gfx={updated} render=true"
	frtn := execDXYN()
	egfx := [gfxS]byte{}
	egfx[133] = 1
	egfx[198] = 1
	egfx[263] = 1
	fgfx := gfx
	evf := byte(0)
	fvf := v[0xF]
	epc := uint16(0x0202)
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
	//TODO: confirm expected out of bounds behaviour then implement test!
}

func TestExecEX9E(t *testing.T) {
	//if the key stored in vx is pressed, skip next instruction
	reset()
	opcode = 0xE39E
	v[3] = 5
	ertn := "execEX9E 0xE39E: pc=0x0202 {skip=false}"
	frtn := execEX9E()
	epc := uint16(0x0202)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	keys[5] = 1
	ertn = "execEX9E 0xE39E: pc=0x0206 {skip=true}"
	frtn = execEX9E()
	epc = uint16(0x0206)
	fpc = pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecEXA1(t *testing.T) {
	//if the key stored in vx is not pressed, skip next instruction
	reset()
	opcode = 0xE2A1
	v[2] = 5
	ertn := "execEXA1 0xE2A1: pc=0x0204 {skip=true}"
	frtn := execEXA1()
	epc := uint16(0x0204)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	keys[5] = 1
	ertn = "execEXA1 0xE2A1: pc=0x0206 {skip=false}"
	frtn = execEXA1()
	epc = uint16(0x0206)
	fpc = pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecFX07(t *testing.T) {
	//vx=delay_timer
	reset()
	opcode = 0xF807
	v[8] = 0xFF
	delayTimer = 0xAB
	ertn := "execFX07 0xF807: pc=0x0202 v[8]=AB"
	frtn := execFX07()
	evx := byte(0xAB)
	fvx := v[8]
	epc := uint16(0x0202)
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

func TestExecFX0A(t *testing.T) {
	//TODO vx=get_key()
	reset()
	opcode = 0xFC0A
}

func TestExecFX15(t *testing.T) {
	//delay_timer=vx
	reset()
	opcode = 0xF715
	v[7] = 0xA1
	ertn := "execFX15 0xF715: pc=0x0202 delayTimer=A1"
	frtn := execFX15()
	edt := byte(0xA1)
	fdt := delayTimer
	epc := uint16(0x0202)
	fpc := pc
	if frtn != ertn {
		t.Errorf("Expected %s, found %s.", ertn, frtn)
	}
	if fdt != edt {
		t.Errorf("Expected 0x%02X, found 0x%02X.", edt, fdt)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecFX18(t *testing.T) {
	//TODO sound_timer=vx
	reset()
	opcode = 0xF118
}

func TestExecFX1E(t *testing.T) {
	//TODO i+=vx
	reset()
	opcode = 0xFF1E
}

func TestExecFX29(t *testing.T) {
	//TODO i=sprite_addr[vx]
	reset()
	opcode = 0xF329
}

func TestExecFX33(t *testing.T) {
	//TODO set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
	reset()
	opcode = 0xF533
}

func TestExecFX55(t *testing.T) {
	//TODO reg_dump(vx,&i)
	reset()
	opcode = 0xFC55
}

func TestExecFX65(t *testing.T) {
	//TODO reg_load(vx,&i)
	reset()
	opcode = 0xFF65
}
