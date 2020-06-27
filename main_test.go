package main

import (
	"testing"
)

func TestFetchOpcode(t *testing.T) {
	vm := newChip8()
	vm.mem[0x0000] = 0xA2
	vm.mem[0x0001] = 0xF0
	vm.mem[0x0002] = 0xC5
	vm.mem[0x0003] = 0x02
	vm.mem[0x0004] = 0x0F
	vm.mem[0x0005] = 0xF0
	vm.mem[0x0006] = 0xAB
	vm.mem[0x0007] = 0x50
	vm.mem[0x0008] = 0x20
	vm.mem[0x0009] = 0xFF
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
	vm := newChip8()
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
	vm := newChip8()
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
	//do nothing
	vm := newChip8()
	vm.oc = 0x0123
	vm.exec0NNN()
	epc := uint16(0x0200)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec00E0(t *testing.T) {
	//clear gfx
	vm := newChip8()
	vm.oc = 0x00E0
	vm.gfx[8] = 1
	vm.gfx[278] = 1
	vm.gfx[1080] = 1
	vm.gfx[2000] = 1
	vm.exec00E0()
	egfx := [2048]byte{}
	fgfx := vm.gfx
	epc := uint16(0x0202)
	fpc := vm.pc
	erg := true
	frg := vm.rg
	if fgfx != egfx {
		t.Errorf("Expected %v, found %v.", egfx, fgfx)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	if frg != erg {
		t.Errorf("Expected %t, found %t.", erg, frg)
	}
}

func TestExec00EE(t *testing.T) {
	//TODO return
	vm := newChip8()
	vm.oc = 0x00EE
}

func TestExec1NNN(t *testing.T) {
	//goto nnn
	vm := newChip8()
	vm.oc = 0x1A9F
	vm.exec1NNN()
	epc := uint16(0x0A9F)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec2NNN(t *testing.T) {
	//call subroutine (increment sp, put current pc on stack, set pc to nnn)
	vm := newChip8()
	vm.oc = 0x2F08
	vm.exec2NNN()
	esp := uint16(0x0001)
	fsp := vm.sp
	essp := uint16(0x0200)
	fssp := vm.stk[vm.sp]
	epc := uint16(0x0F08)
	fpc := vm.pc
	if fsp != esp {
		t.Errorf("Expected 0x%04X, found 0x%04X.", esp, fsp)
	}
	if fssp != essp {
		t.Errorf("Expected 0x%04X, found 0x%04X.", essp, fssp)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec3XNN(t *testing.T) {
	//if(vx==nn) skip next instruction
	vm := newChip8()
	vm.oc = 0x3A1D
	vm.vr[0xA] = 0x1D
	vm.exec3XNN()
	epc := uint16(0x0204)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	vm.vr[0xA] = 0xFD
	vm.exec3XNN()
	epc = uint16(0x0206)
	fpc = vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec4XNN(t *testing.T) {
	//if(vx!=nn) skip next instruction
	vm := newChip8()
	vm.oc = 0x4247
	vm.vr[0x2] = 0xAC
	vm.exec4XNN()
	epc := uint16(0x0204)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	vm.vr[0x2] = 0x47
	vm.exec4XNN()
	epc = uint16(0x0206)
	fpc = vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec5XY0(t *testing.T) {
	//if(vx==vy) skip next instruction
	vm := newChip8()
	vm.oc = 0x5190
	vm.vr[0x1] = 0xAC
	vm.vr[0x9] = 0xAC
	vm.exec5XY0()
	epc := uint16(0x0204)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	vm.vr[0x1] = 0x47
	vm.exec5XY0()
	epc = uint16(0x0206)
	fpc = vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec6XNN(t *testing.T) {
	//vx=nn
	vm := newChip8()
	vm.oc = 0x6EFD
	vm.exec6XNN()
	evx := byte(0xFD)
	fvx := vm.vr[0xE]
	epc := uint16(0x0202)
	fpc := vm.pc
	if fvx != evx {
		t.Errorf("Expected 0x%02X, found 0x%02X.", evx, fvx)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec7XNN(t *testing.T) {
	//vx+=nn
	vm := newChip8()
	vm.oc = 0x7315
	vm.vr[0x3] = 0xA1
	vm.exec7XNN()
	evx := byte(0xB6)
	fvx := vm.vr[0x3]
	epc := uint16(0x0202)
	fpc := vm.pc
	if fvx != evx {
		t.Errorf("Expected 0x%02X, found 0x%02X.", evx, fvx)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExec8XY0(t *testing.T) {
	//TODO vx=vy
	vm := newChip8()
	vm.oc = 0x89B0
}

func TestExec8XY1(t *testing.T) {
	//TODO vx=vx|vy
	vm := newChip8()
	vm.oc = 0x8DE1
}

func TestExec8XY2(t *testing.T) {
	//TODO vx=vx&vy
	vm := newChip8()
	vm.oc = 0x83F2
}

func TestExec8XY3(t *testing.T) {
	//TODO vx=vx^vy
	vm := newChip8()
	vm.oc = 0x8AC3
}

func TestExec8XY4(t *testing.T) {
	//TODO vx+=vy
	vm := newChip8()
	vm.oc = 0x8474
}

func TestExec8XY5(t *testing.T) {
	//TODO vx-=vy
	vm := newChip8()
	vm.oc = 0x89A5
}

func TestExec8XY6(t *testing.T) {
	//TODO vx>>=1
	vm := newChip8()
	vm.oc = 0x8206
}

func TestExec8XY7(t *testing.T) {
	//TODO vx=vy-vx
	vm := newChip8()
	vm.oc = 0x8197
}

func TestExec8XYE(t *testing.T) {
	//TODO vx<<=1
	vm := newChip8()
	vm.oc = 0x8EEE
}

func TestExec9XY0(t *testing.T) {
	//if(vx!=vy) skip next instruction
	vm := newChip8()
	vm.oc = 0x9730
	vm.vr[0x7] = 0xAC
	vm.vr[0x3] = 0x47
	vm.exec9XY0()
	epc := uint16(0x0204)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	vm.vr[0x3] = 0xAC
	vm.exec9XY0()
	epc = uint16(0x0206)
	fpc = vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecANNN(t *testing.T) {
	//i=nnn
	vm := newChip8()
	vm.oc = 0xA259
	vm.execANNN()
	ei := uint16(0x0259)
	fi := vm.ir
	epc := uint16(0x0202)
	fpc := vm.pc
	if fi != ei {
		t.Errorf("Expected 0x%04X, found 0x%04X.", ei, fi)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecBNNN(t *testing.T) {
	//TODO pc=v0+nnn
	vm := newChip8()
	vm.oc = 0xBAF2
}

func TestExecCXNN(t *testing.T) {
	//vx=rand()&nn
	vm := newChip8()
	vm.oc = 0xC400
	vm.execCXNN()
	evx := byte(0x00)
	fvx := vm.vr[0x4]
	epc := uint16(0x0202)
	fpc := vm.pc
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
	vm := newChip8()
	//draw a 3x3 sprite at 5,2 (sprite: diagonal line topleft to bottomright)
	vm.oc = 0xD793
	vm.mem[0x0101] = 0x80
	vm.mem[0x0102] = 0x40
	vm.mem[0x0103] = 0x20
	vm.ir = 0x0101
	vm.vr[7] = 5
	vm.vr[9] = 2
	vm.execDXYN()
	egfx := [2048]byte{}
	egfx[133] = 1
	egfx[198] = 1
	egfx[263] = 1
	fgfx := vm.gfx
	evf := byte(0)
	fvf := vm.vr[0xF]
	epc := uint16(0x0202)
	fpc := vm.pc
	erg := true
	frg := vm.rg
	if fgfx != egfx {
		t.Errorf("Expected %v, found %v.", egfx, fgfx)
	}
	if fvf != evf {
		t.Errorf("Expected %v, found %v.", evf, fvf)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	if frg != erg {
		t.Errorf("Expected %t, found %t.", erg, frg)
	}
	//draw the sprite again at 8,5 (no pixels erased)
	vm.vr[7] = 8
	vm.vr[9] = 5
	vm.execDXYN()
	evf = byte(0)
	fvf = vm.vr[0xF]
	if fvf != evf {
		t.Errorf("Expected %v, found %v.", evf, fvf)
	}
	//draw the sprite again at 10,7 (pixel erased)
	vm.vr[7] = 10
	vm.vr[9] = 7
	vm.execDXYN()
	evf = byte(1)
	fvf = vm.vr[0xF]
	if fvf != evf {
		t.Errorf("Expected %v, found %v.", evf, fvf)
	}
	//TODO: confirm expected out of bounds behaviour then implement test!
}

func TestExecEX9E(t *testing.T) {
	//if the key stored in vx is pressed, skip next instruction
	vm := newChip8()
	vm.oc = 0xE39E
	vm.vr[3] = 5
	vm.execEX9E()
	epc := uint16(0x0202)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	vm.key[5] = 1
	vm.execEX9E()
	epc = uint16(0x0206)
	fpc = vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecEXA1(t *testing.T) {
	//if the key stored in vx is not pressed, skip next instruction
	vm := newChip8()
	vm.oc = 0xE2A1
	vm.vr[2] = 5
	vm.execEXA1()
	epc := uint16(0x0204)
	fpc := vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
	vm.key[5] = 1
	vm.execEXA1()
	epc = uint16(0x0206)
	fpc = vm.pc
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecFX07(t *testing.T) {
	//vx=delay_timer
	vm := newChip8()
	vm.oc = 0xF807
	vm.vr[8] = 0xFF
	vm.dt = 0xAB
	vm.execFX07()
	evx := byte(0xAB)
	fvx := vm.vr[8]
	epc := uint16(0x0202)
	fpc := vm.pc
	if fvx != evx {
		t.Errorf("Expected 0x%02X, found 0x%02X.", evx, fvx)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecFX0A(t *testing.T) {
	//TODO vx=get_key()
	vm := newChip8()
	vm.oc = 0xFC0A
}

func TestExecFX15(t *testing.T) {
	//delay_timer=vx
	vm := newChip8()
	vm.oc = 0xF715
	vm.vr[7] = 0xA1
	vm.execFX15()
	edt := byte(0xA1)
	fdt := vm.dt
	epc := uint16(0x0202)
	fpc := vm.pc
	if fdt != edt {
		t.Errorf("Expected 0x%02X, found 0x%02X.", edt, fdt)
	}
	if fpc != epc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", epc, fpc)
	}
}

func TestExecFX18(t *testing.T) {
	//TODO sound_timer=vx
	vm := newChip8()
	vm.oc = 0xF118
}

func TestExecFX1E(t *testing.T) {
	//TODO i+=vx
	vm := newChip8()
	vm.oc = 0xFF1E
}

func TestExecFX29(t *testing.T) {
	//TODO i=sprite_addr[vx]
	vm := newChip8()
	vm.oc = 0xF329
}

func TestExecFX33(t *testing.T) {
	//TODO set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
	vm := newChip8()
	vm.oc = 0xF533
}

func TestExecFX55(t *testing.T) {
	//TODO reg_dump(vx,&i)
	vm := newChip8()
	vm.oc = 0xFC55
}

func TestExecFX65(t *testing.T) {
	//TODO reg_load(vx,&i)
	vm := newChip8()
	vm.oc = 0xFF65
}
