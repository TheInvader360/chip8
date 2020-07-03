package lib

import (
	"reflect"
	"testing"
)

func TestExec0NNN(t *testing.T) {
	//do nothing
	f := NewChip8()
	f.oc = 0x0123
	f.exec0NNN()
	e := NewChip8()
	e.oc = 0x0123
	checkEqual(t, e, f)
}

func TestExec00CN(t *testing.T) {
	//TODO - scroll down n lines
}

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
	//return (set pc from stack, decrement sp)
	f := NewChip8()
	f.oc = 0x00EE
	f.stk[0x1] = 0x1234
	f.sp = 1
	f.exec00EE()
	e := NewChip8()
	e.oc = 0x00EE
	e.stk[0x1] = 0x1234
	e.sp = 0
	e.pc = 0x1236
	checkEqual(t, e, f)
}

func TestExec00FB(t *testing.T) {
	//TODO - scroll right 4 pixels
}

func TestExec00FC(t *testing.T) {
	//TODO - scroll left 4 pixels
}

func TestExec00FD(t *testing.T) {
	//TODO - exit (reset/quit?)
}

func TestExec00FE(t *testing.T) {
	//lo-res mode
	f := NewChip8()
	f.oc = 0x00FE
	f.exec00FE()
	e := NewChip8()
	e.oc = 0x00FE
	e.mode = sclr
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExec00FF(t *testing.T) {
	//hi-res mode
	f := NewChip8()
	f.oc = 0x00FF
	f.exec00FF()
	e := NewChip8()
	e.oc = 0x00FF
	e.mode = schr
	e.pc = 0x0202
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
	e.stk[0x1] = 0x0200
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
	//vx=vy
	f := NewChip8()
	f.oc = 0x89B0
	f.vr[0x9] = 0x0F
	f.vr[0xB] = 0x0A
	f.exec8XY0()
	e := NewChip8()
	e.oc = 0x89B0
	e.vr[0x9] = 0x0A
	e.vr[0xB] = 0x0A
	e.pc = 0x202
	checkEqual(t, e, f)
}

func TestExec8XY1(t *testing.T) {
	//vx=vx|vy (bitwise OR)
	f := NewChip8()
	f.oc = 0x8DE1
	f.vr[0xD] = 0b01010101
	f.vr[0xE] = 0b00110011
	f.exec8XY1()
	e := NewChip8()
	e.oc = 0x8DE1
	e.vr[0xD] = 0b01110111
	e.vr[0xE] = 0b00110011
	e.pc = 0x202
	checkEqual(t, e, f)
}

func TestExec8XY2(t *testing.T) {
	//vx=vx&vy (bitwise AND)
	f := NewChip8()
	f.oc = 0x83F2
	f.vr[0x3] = 0b01010101
	f.vr[0xF] = 0b00110011
	f.exec8XY2()
	e := NewChip8()
	e.oc = 0x83F2
	e.vr[0x3] = 0b00010001
	e.vr[0xF] = 0b00110011
	e.pc = 0x202
	checkEqual(t, e, f)
}

func TestExec8XY3(t *testing.T) {
	//vx=vx^vy (bitwise XOR)
	f := NewChip8()
	f.oc = 0x8AC3
	f.vr[0xA] = 0b01010101
	f.vr[0xC] = 0b00110011
	f.exec8XY3()
	e := NewChip8()
	e.oc = 0x8AC3
	e.vr[0xA] = 0b01100110
	e.vr[0xC] = 0b00110011
	e.pc = 0x202
	checkEqual(t, e, f)
}

func TestExec8XY4(t *testing.T) {
	//vx+=vy (only stores lowest 8 bits of result, if result > 0xFF then vF=1)
	f := NewChip8()
	f.oc = 0x8474
	f.vr[0x4] = 0b11111111
	f.vr[0x7] = 0b00000001
	f.exec8XY4()
	e := NewChip8()
	e.oc = 0x8474
	e.vr[0x4] = 0b00000000
	e.vr[0x7] = 0b00000001
	e.vr[0xF] = 1
	e.pc = 0x202
	checkEqual(t, e, f)

	f.vr[0x4] = 0b11000011
	f.vr[0x7] = 0b00001111
	f.exec8XY4()
	e.vr[0x4] = 0b11010010
	e.vr[0x7] = 0b00001111
	e.vr[0xF] = 0
	e.pc = 0x204
	checkEqual(t, e, f)
}

func TestExec8XY5(t *testing.T) {
	//vx-=vy (if vx>vy then vF=1)
	f := NewChip8()
	f.oc = 0x89A5
	f.vr[0x9] = 0b00001111
	f.vr[0xA] = 0b00000110
	f.exec8XY5()
	e := NewChip8()
	e.oc = 0x89A5
	e.vr[0x9] = 0b00001001
	e.vr[0xA] = 0b00000110
	e.vr[0xF] = 1
	e.pc = 0x202
	checkEqual(t, e, f)

	f.vr[0x9] = 0b00000110
	f.vr[0xA] = 0b00001111
	f.exec8XY5()
	e.vr[0x9] = 0b11110111 //i.e. 0b100000110-0b1111
	e.vr[0xA] = 0b00001111
	e.vr[0xF] = 0
	e.pc = 0x204
	checkEqual(t, e, f)
}

func TestExec8XY6(t *testing.T) {
	//vx>>=1 (vF=the lsb of vx, then vx is divided by 2)
	f := NewChip8()
	f.oc = 0x8206
	f.vr[0x2] = 0b01010101
	f.exec8XY6()
	e := NewChip8()
	e.oc = 0x8206
	e.vr[0x2] = 0b00101010
	e.vr[0xF] = 1
	e.pc = 0x202
	checkEqual(t, e, f)

	f.vr[0x2] = 0b00111100
	f.exec8XY6()
	e.vr[0x2] = 0b00011110
	e.vr[0xF] = 0
	e.pc = 0x204
	checkEqual(t, e, f)
}

func TestExec8XY7(t *testing.T) {
	//vx=vy-vx (if vy>vx then vF=1)
	f := NewChip8()
	f.oc = 0x8197
	f.vr[0x1] = 0b01010101
	f.vr[0x9] = 0b00011100
	f.exec8XY7()
	e := NewChip8()
	e.oc = 0x8197
	e.vr[0x1] = 0b11000111 //i.e. 0b100011100-0b1010101
	e.vr[0x9] = 0b00011100
	e.vr[0xF] = 0
	e.pc = 0x202
	checkEqual(t, e, f)

	f.vr[0x1] = 0b00011100
	f.vr[0x9] = 0b01010101
	f.exec8XY7()
	e.vr[0x1] = 0b00111001 //i.e. 0b01010101-0b00011100
	e.vr[0x9] = 0b01010101
	e.vr[0xF] = 1
	e.pc = 0x204
	checkEqual(t, e, f)
}

func TestExec8XYE(t *testing.T) {
	//vx<<=1 (vF=the msb of vx, then vx is multiplied by 2)
	f := NewChip8()
	f.oc = 0x8EEE
	f.vr[0xE] = 0b01010101
	f.exec8XYE()
	e := NewChip8()
	e.oc = 0x8EEE
	e.vr[0xE] = 0b10101010
	e.vr[0xF] = 0
	e.pc = 0x202
	checkEqual(t, e, f)

	f.vr[0xE] = 0b11000000
	f.exec8XYE()
	e.vr[0xE] = 0b10000000
	e.vr[0xF] = 1
	e.pc = 0x204
	checkEqual(t, e, f)
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
	//pc=v0+nnn
	f := NewChip8()
	f.oc = 0xBAF2
	f.vr[0x0] = 0x23
	f.execBNNN()
	e := NewChip8()
	e.oc = 0xBAF2
	e.vr[0x0] = 0x23
	e.pc = 0x0B15 //i.e. 0x23 + 0x0AF2
	checkEqual(t, e, f)
}

func TestExecCXNN(t *testing.T) {
	//vx=rand()&nn
	f := NewChip8()
	f.oc = 0xC400
	f.execCXNN()
	e := NewChip8()
	e.oc = 0xC400
	e.vr[0x4] = 0x00
	e.pc = 0x0202
	checkEqual(t, e, f)
	//TODO: reliable tests where nn is not zero...
}

func TestExecDXYN(t *testing.T) {
	//TODO - call execDXYNHR/execDXYNLR dependent on mode (hi-res/lo-res)
}

func TestExecDXYNLR(t *testing.T) {
	//draw(x,y,n) - draw n byte sprite from mem[i] at vx,xy (vf=collision)

	//draw a 2x2 sprite at 20,10 (sprite: filled square)
	f := NewChip8()
	f.oc = 0xD792
	f.Mem[0x0101] = 0xC0
	f.Mem[0x0102] = 0xC0
	f.ir = 0x0101
	f.vr[0x7] = 20
	f.vr[0x9] = 10
	f.execDXYNLR()
	e := NewChip8()
	//translate logical 2x2 at 20,10 to gfx 4x4 at 40,20
	s := (20*2 + 10*2*128)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			e.Gfx[s+x+y*128] = 1
		}
	}
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
	f.vr[0x7] = 63
	f.vr[0x9] = 31
	f.execDXYNLR()
	//top-left corner - logical 1x1 / gfx 2x2
	e.Gfx[0+0*128] = 1
	e.Gfx[1+0*128] = 1
	e.Gfx[0+1*128] = 1
	e.Gfx[1+1*128] = 1
	//top-right corner - logical 1x1 / gfx 2x2
	e.Gfx[126+0*128] = 1
	e.Gfx[127+0*128] = 1
	e.Gfx[126+1*128] = 1
	e.Gfx[127+1*128] = 1
	//bottom-left corner - logical 1x1 / gfx 2x2
	e.Gfx[0+62*128] = 1
	e.Gfx[1+62*128] = 1
	e.Gfx[0+63*128] = 1
	e.Gfx[1+63*128] = 1
	//bottom-right corner - logical 1x1 / gfx 2x2
	e.Gfx[126+62*128] = 1
	e.Gfx[127+62*128] = 1
	e.Gfx[126+63*128] = 1
	e.Gfx[127+63*128] = 1
	if f.Gfx != e.Gfx {
		t.Errorf("Expected %v, found %v.", e.Gfx, f.Gfx)
	}

	//draw the sprite again at 22,12 (no pixels erased)
	f.vr[0x7] = 22
	f.vr[0x9] = 12
	f.execDXYNLR()
	//translate logical 2x2 at 22,12 to gfx 4x4 at 44,24
	s = (22*2 + 12*2*128)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			e.Gfx[s+x+y*128] = 1
		}
	}
	e.vr[0xF] = 0
	if f.Gfx != e.Gfx {
		t.Errorf("Expected %v, found %v.", e.Gfx, f.Gfx)
	}
	if f.vr[0xF] != e.vr[0xF] {
		t.Errorf("Expected %v, found %v.", e.vr[0xF], f.vr[0xF])
	}

	//draw the sprite again at 23,13 (pixel erased)
	f.vr[0x7] = 23
	f.vr[0x9] = 13
	f.execDXYNLR()
	//translate logical 2x2 at 23,13 to gfx 4x4 at 46,26 (on)
	s = (23*2 + 13*2*128)
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			e.Gfx[s+x+y*128] = 1
		}
	}
	//translate logical 1x1 at 23,13 to gfx 2x2 at 46,26 (off)
	s = (23*2 + 13*2*128)
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			e.Gfx[s+x+y*128] = 0
		}
	}
	e.vr[0xF] = 1
	if f.Gfx != e.Gfx {
		t.Errorf("Expected %v, found %v.", e.Gfx, f.Gfx)
	}
	if f.vr[0xF] != e.vr[0xF] {
		t.Errorf("Expected %v, found %v.", e.vr[0xF], f.vr[0xF])
	}
}

func TestExecDXYNHR(t *testing.T) {
	//TODO - draw(x,y) - draw 16x16 sprite from mem[i] at vx,xy (vf=collision)
}

func TestExecEX9E(t *testing.T) {
	//if the key stored in vx is pressed, skip next instruction
	f := NewChip8()
	f.oc = 0xE39E
	f.vr[0x3] = 5
	f.execEX9E()
	e := NewChip8()
	e.oc = 0xE39E
	e.vr[0x3] = 5
	e.pc = 0x0202
	checkEqual(t, e, f)

	f.Key[0x5] = 1
	f.execEX9E()
	e.Key[0x5] = 1
	e.pc = 0x0206
	checkEqual(t, e, f)
}

func TestExecEXA1(t *testing.T) {
	//if the key stored in vx is not pressed, skip next instruction
	f := NewChip8()
	f.oc = 0xE2A1
	f.vr[0x2] = 5
	f.execEXA1()
	e := NewChip8()
	e.oc = 0xE2A1
	e.vr[0x2] = 5
	e.pc = 0x0204
	checkEqual(t, e, f)

	f.Key[0x5] = 1
	f.execEXA1()
	e.Key[0x5] = 1
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
	e.vr[0x8] = 0xAB
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX0A(t *testing.T) {
	//vx=get_key() (wait for a key press then store the key value in vx)
	f := NewChip8()
	f.oc = 0xFC0A
	f.execFX0A()
	e := NewChip8()
	e.oc = 0xFC0A
	checkEqual(t, e, f)

	f.execFX0A()
	checkEqual(t, e, f)

	f.execFX0A()
	checkEqual(t, e, f)

	f.Key[0xD] = 1
	f.execFX0A()
	e.Key[0xD] = 1
	e.vr[0xC] = 0xD
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX15(t *testing.T) {
	//delay_timer=vx
	f := NewChip8()
	f.oc = 0xF715
	f.vr[0x7] = 0xA1
	f.execFX15()
	e := NewChip8()
	e.oc = 0xF715
	e.vr[0x7] = 0xA1
	e.dt = 0xA1
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX18(t *testing.T) {
	//sound_timer=vx
	f := NewChip8()
	f.oc = 0xF118
	f.vr[0x1] = 0xD4
	f.execFX18()
	e := NewChip8()
	e.oc = 0xF118
	e.vr[0x1] = 0xD4
	e.St = 0xD4
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX1E(t *testing.T) {
	// i+=vx
	f := NewChip8()
	f.oc = 0xFE1E
	f.vr[0xE] = 0xA6
	f.ir = 0x1111
	f.execFX1E()
	e := NewChip8()
	e.oc = 0xFE1E
	e.vr[0xE] = 0xA6
	e.ir = 0x11B7
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX29(t *testing.T) {
	//i=sprite_addr[vx] (point i at 5 byte font sprite for hex char at vx)
	f := NewChip8()
	f.oc = 0xF329
	f.vr[0x3] = 0xF
	f.execFX29()
	e := NewChip8()
	e.oc = 0xF329
	e.vr[0x3] = 0xF
	e.ir = 0x004B //"F" sprite should be stored in Mem[0x4B]:Mem[0x4F]
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX30(t *testing.T) {
	//i=sprite_addr[vx] (point i at 10 byte font sprite for hex char at vx)
	f := NewChip8()
	f.oc = 0xF230
	f.vr[0x2] = 0xA
	f.execFX30()
	e := NewChip8()
	e.oc = 0xF230
	e.vr[0x2] = 0xA
	e.ir = 0x00B4 //"A" sprite should be stored in Mem[0xB4]:Mem[0xBD]
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX33(t *testing.T) {
	//set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
	//store a decimal of vx in memory (e.g. if i=0 and vx=128, m0=1 m1=2 m2=8)
	f := NewChip8()
	f.oc = 0xF533
	f.ir = 0x00A0
	f.vr[0x5] = 128
	f.execFX33()
	e := NewChip8()
	e.oc = 0xF533
	e.ir = 0x00A0
	e.vr[0x5] = 128
	e.Mem[0x00A0] = 1
	e.Mem[0x00A1] = 2
	e.Mem[0x00A2] = 8
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX55(t *testing.T) {
	//reg_dump(vx,&i) (store v0:vx inclusive from addr i, i is not modified)
	f := NewChip8()
	f.oc = 0xF255
	f.ir = 0x0100
	f.vr[0x0] = 0x42
	f.vr[0x1] = 0xD4
	f.vr[0x2] = 0xAB
	f.execFX55()
	e := NewChip8()
	e.oc = 0xF255
	e.ir = 0x0100
	e.vr[0x0] = 0x42
	e.vr[0x1] = 0xD4
	e.vr[0x2] = 0xAB
	e.Mem[0x0100] = 0x42
	e.Mem[0x0101] = 0xD4
	e.Mem[0x0102] = 0xAB
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX65(t *testing.T) {
	//reg_load(vx,&i) (fill v0:vx inclusive from addr i, i is not modified)
	f := NewChip8()
	f.oc = 0xF265
	f.ir = 0x0100
	f.Mem[0x0100] = 0x42
	f.Mem[0x0101] = 0xD4
	f.Mem[0x0102] = 0xAB
	f.execFX65()
	e := NewChip8()
	e.oc = 0xF265
	e.ir = 0x0100
	e.Mem[0x0100] = 0x42
	e.Mem[0x0101] = 0xD4
	e.Mem[0x0102] = 0xAB
	e.vr[0x0] = 0x42
	e.vr[0x1] = 0xD4
	e.vr[0x2] = 0xAB
	e.pc = 0x0202
	checkEqual(t, e, f)
}

func TestExecFX75(t *testing.T) {
	//TODO - store v0:vx in rpl user flags (x <= 7)... investigate!
}

func TestExecFX85(t *testing.T) {
	//TODO - read v0:vx from rpl user flags (x <= 7)... investigate!
}

func checkEqual(t *testing.T, e *Chip8, f *Chip8) {
	t.Helper()
	if !reflect.DeepEqual(e, f) {
		t.Errorf("Expected %v, found %v.", e, f)
	}
}
