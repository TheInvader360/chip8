package lib

import (
	"reflect"
	"testing"
)

func TestExec00E0(t *testing.T) {
	//clear gfx
	f := NewChip8()
	f.oc = 0x00E0
	f.Gfx[8] = 1
	f.Gfx[278] = 1
	f.Gfx[1080] = 1
	f.exec00E0()
	e := NewChip8()
	e.oc = 0x00E0
	e.pc = 0x0202
	e.Rg = true
	checkEqual(t, e, f)
}

func TestExec00EE(t *testing.T) {
	//TODO return
	f := NewChip8()
	f.oc = 0x00EE
}

func TestExec0NNN(t *testing.T) {
	//do nothing
	f := NewChip8()
	f.oc = 0x0123
	f.exec0NNN()
	e := NewChip8()
	e.oc = 0x0123
	checkEqual(t, e, f)
}

func TestExec1NNN(t *testing.T) {
	//goto nnn
	f := NewChip8()
	f.oc = 0x1A9F
	f.exec1NNN()
	e := NewChip8()
	e.oc = 0x1A9F
	e.pc = 0x0A9F
	checkEqual(t, e, f)
}

func TestExec2NNN(t *testing.T) {
	//call subroutine (increment sp, put current pc on stack, set pc to nnn)
	f := NewChip8()
	f.oc = 0x2F08
	f.exec2NNN()
	e := NewChip8()
	e.oc = 0x2F08
	e.sp = 0x0001
	e.stk[1] = 0x0200
	e.pc = 0x0F08
	checkEqual(t, e, f)
}

func TestExec3XNN(t *testing.T) {
	//if(vx==nn) skip next instruction
	f := NewChip8()
	f.oc = 0x3A1D
	f.vr[0xA] = 0x1D
	f.exec3XNN()
	e := NewChip8()
	e.oc = 0x3A1D
	e.vr[0xA] = 0x1D
	e.pc = 0x0204
	checkEqual(t, e, f)

	f.vr[0xA] = 0xFD
	f.exec3XNN()
	e.vr[0xA] = 0xFD
	e.pc = 0x0206
	checkEqual(t, e, f)
}

func TestExec4XNN(t *testing.T) {
	//if(vx!=nn) skip next instruction
	f := NewChip8()
	f.oc = 0x4247
	f.vr[0x2] = 0xAC
	f.exec4XNN()
	e := NewChip8()
	e.oc = 0x4247
	e.vr[0x2] = 0xAC
	e.pc = 0x0204
	checkEqual(t, e, f)

	f.vr[0x2] = 0x47
	f.exec4XNN()
	e.vr[0x2] = 0x47
	e.pc = 0x0206
	checkEqual(t, e, f)
}

func TestExec5XY0(t *testing.T) {
	//if(vx==vy) skip next instruction
	f := NewChip8()
	f.oc = 0x5190
	f.vr[0x1] = 0xAC
	f.vr[0x9] = 0xAC
	f.exec5XY0()
	e := NewChip8()
	e.oc = 0x5190
	e.vr[0x1] = 0xAC
	e.vr[0x9] = 0xAC
	e.pc = 0x0204
	checkEqual(t, e, f)

	f.vr[0x1] = 0x47
	f.exec5XY0()
	e.vr[0x1] = 0x47
	e.pc = 0x0206
	checkEqual(t, e, f)
}

func TestExec6XNN(t *testing.T) {
	//vx=nn
	f := NewChip8()
	f.oc = 0x6EFD
	f.exec6XNN()
	e := NewChip8()
	e.oc = 0x6EFD
	e.vr[0xE] = 0xFD
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExec7XNN(t *testing.T) {
	//vx+=nn
	f := NewChip8()
	f.oc = 0x7315
	f.vr[0x3] = 0xA1
	f.exec7XNN()
	e := NewChip8()
	e.oc = 0x7315
	e.vr[0x3] = 0xB6
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExec8XY0(t *testing.T) {
	//TODO vx=vy
	f := NewChip8()
	f.oc = 0x89B0
}

func TestExec8XY1(t *testing.T) {
	//TODO vx=vx|vy
	f := NewChip8()
	f.oc = 0x8DE1
}

func TestExec8XY2(t *testing.T) {
	//TODO vx=vx&vy
	f := NewChip8()
	f.oc = 0x83F2
}

func TestExec8XY3(t *testing.T) {
	//TODO vx=vx^vy
	f := NewChip8()
	f.oc = 0x8AC3
}

func TestExec8XY4(t *testing.T) {
	//TODO vx+=vy
	f := NewChip8()
	f.oc = 0x8474
}

func TestExec8XY5(t *testing.T) {
	//TODO vx-=vy
	f := NewChip8()
	f.oc = 0x89A5
}

func TestExec8XY6(t *testing.T) {
	//TODO vx>>=1
	f := NewChip8()
	f.oc = 0x8206
}

func TestExec8XY7(t *testing.T) {
	//TODO vx=vy-vx
	f := NewChip8()
	f.oc = 0x8197
}

func TestExec8XYE(t *testing.T) {
	//TODO vx<<=1
	f := NewChip8()
	f.oc = 0x8EEE
}

func TestExec9XY0(t *testing.T) {
	//if(vx!=vy) skip next instruction
	f := NewChip8()
	f.oc = 0x9730
	f.vr[0x7] = 0xAC
	f.vr[0x3] = 0x47
	f.exec9XY0()
	e := NewChip8()
	e.oc = 0x9730
	e.vr[0x7] = 0xAC
	e.vr[0x3] = 0x47
	e.pc = 0x0204
	checkEqual(t, e, f)

	f.vr[0x3] = 0xAC
	f.exec9XY0()
	e.vr[0x3] = 0xAC
	e.pc = 0x0206
	checkEqual(t, e, f)
}

func TestExecANNN(t *testing.T) {
	//i=nnn
	f := NewChip8()
	f.oc = 0xA259
	f.execANNN()
	e := NewChip8()
	e.oc = 0xA259
	e.ir = 0x0259
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecBNNN(t *testing.T) {
	//TODO pc=v0+nnn
	f := NewChip8()
	f.oc = 0xBAF2
}

func TestExecCXNN(t *testing.T) {
	//vx=rand()&nn
	f := NewChip8()
	f.oc = 0xC400
	f.execCXNN()
	e := NewChip8()
	e.oc = 0xC400
	e.vr[4] = 0x00
	e.pc = 0x0202
	checkEqual(t, e, f)
	//TODO: reliable tests where nn is not zero...
}

func TestExecDXYN(t *testing.T) {
	//draw(x,y,n)

	//draw a 2x2 sprite at 20,10 (sprite: filled square)
	f := NewChip8()
	f.oc = 0xD792
	f.Mem[0x0101] = 0xC0
	f.Mem[0x0102] = 0xC0
	f.ir = 0x0101
	f.vr[7] = 20
	f.vr[9] = 10
	f.execDXYN()
	e := NewChip8()
	e.Gfx[20+10*64] = 1
	e.Gfx[21+10*64] = 1
	e.Gfx[20+11*64] = 1
	e.Gfx[21+11*64] = 1
	e.vr[0xF] = 0
	e.pc = 0x0202
	e.Rg = true
	if f.Gfx != e.Gfx {
		t.Errorf("Expected %v, found %v.", e.Gfx, f.Gfx)
	}
	if f.vr[0xF] != e.vr[0xF] {
		t.Errorf("Expected %v, found %v.", e.vr[0xF], f.vr[0xF])
	}
	if f.pc != e.pc {
		t.Errorf("Expected 0x%04X, found 0x%04X.", e.pc, f.pc)
	}
	if f.Rg != e.Rg {
		t.Errorf("Expected %t, found %t.", e.Rg, f.Rg)
	}

	//draw the sprite again at 63,31 (overlaps all edges)
	f.vr[7] = 63
	f.vr[9] = 31
	f.execDXYN()
	e.Gfx[0+0*64] = 1
	e.Gfx[0+31*64] = 1
	e.Gfx[63+0*64] = 1
	e.Gfx[63+31*64] = 1
	e.Gfx[1+1*64] = 0
	e.Gfx[1+30*64] = 0
	e.Gfx[62+1*64] = 0
	e.Gfx[62+30*64] = 0
	if f.Gfx != e.Gfx {
		t.Errorf("Expected %v, found %v.", e.Gfx, f.Gfx)
	}

	//draw the sprite again at 22,12 (no pixels erased)
	f.vr[7] = 22
	f.vr[9] = 12
	f.execDXYN()
	e.Gfx[22+12*64] = 1
	e.Gfx[23+12*64] = 1
	e.Gfx[22+13*64] = 1
	e.Gfx[23+13*64] = 1
	e.vr[0xF] = 0
	if f.Gfx != e.Gfx {
		t.Errorf("Expected %v, found %v.", e.Gfx, f.Gfx)
	}
	if f.vr[0xF] != e.vr[0xF] {
		t.Errorf("Expected %v, found %v.", e.vr[0xF], f.vr[0xF])
	}

	//draw the sprite again at 23,13 (pixel erased)
	f.vr[7] = 23
	f.vr[9] = 13
	f.execDXYN()
	e.Gfx[23+13*64] = 0
	e.Gfx[24+13*64] = 1
	e.Gfx[23+14*64] = 1
	e.Gfx[24+14*64] = 1
	e.vr[0xF] = 1
	if f.Gfx != e.Gfx {
		t.Errorf("Expected %v, found %v.", e.Gfx, f.Gfx)
	}
	if f.vr[0xF] != e.vr[0xF] {
		t.Errorf("Expected %v, found %v.", e.vr[0xF], f.vr[0xF])
	}
}

func TestExecEX9E(t *testing.T) {
	//if the key stored in vx is pressed, skip next instruction
	f := NewChip8()
	f.oc = 0xE39E
	f.vr[3] = 5
	f.execEX9E()
	e := NewChip8()
	e.oc = 0xE39E
	e.vr[3] = 5
	e.pc = 0x0202
	checkEqual(t, e, f)

	f.Key[5] = 1
	f.execEX9E()
	e.Key[5] = 1
	e.pc = 0x0206
	checkEqual(t, e, f)
}

func TestExecEXA1(t *testing.T) {
	//if the key stored in vx is not pressed, skip next instruction
	f := NewChip8()
	f.oc = 0xE2A1
	f.vr[2] = 5
	f.execEXA1()
	e := NewChip8()
	e.oc = 0xE2A1
	e.vr[2] = 5
	e.pc = 0x0204
	checkEqual(t, e, f)

	f.Key[5] = 1
	f.execEXA1()
	e.Key[5] = 1
	e.pc = 0x0206
	checkEqual(t, e, f)
}

func TestExecFX07(t *testing.T) {
	//vx=delay_timer
	f := NewChip8()
	f.oc = 0xF807
	f.dt = 0xAB
	f.execFX07()
	e := NewChip8()
	e.oc = 0xF807
	e.dt = 0xAB
	e.vr[8] = 0xAB
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX0A(t *testing.T) {
	//TODO vx=get_key()
	f := NewChip8()
	f.oc = 0xFC0A
}

func TestExecFX15(t *testing.T) {
	//delay_timer=vx
	f := NewChip8()
	f.oc = 0xF715
	f.vr[7] = 0xA1
	f.execFX15()
	e := NewChip8()
	e.oc = 0xF715
	e.vr[7] = 0xA1
	e.dt = 0xA1
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX18(t *testing.T) {
	//TODO sound_timer=vx
	f := NewChip8()
	f.oc = 0xF118
}

func TestExecFX1E(t *testing.T) {
	//TODO i+=vx
	f := NewChip8()
	f.oc = 0xFF1E
}

func TestExecFX29(t *testing.T) {
	//TODO i=sprite_addr[vx]
	f := NewChip8()
	f.oc = 0xF329
}

func TestExecFX33(t *testing.T) {
	//TODO set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
	f := NewChip8()
	f.oc = 0xF533
}

func TestExecFX55(t *testing.T) {
	//TODO reg_dump(vx,&i)
	f := NewChip8()
	f.oc = 0xFC55
}

func TestExecFX65(t *testing.T) {
	//TODO reg_load(vx,&i)
	f := NewChip8()
	f.oc = 0xFF65
}

func checkEqual(t *testing.T, e *Chip8, f *Chip8) {
	t.Helper()
	if !reflect.DeepEqual(e, f) {
		t.Errorf("Expected %v, found %v.", e, f)
	}
}
