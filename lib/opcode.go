package lib

import (
	"math/rand"
)

func (vm *Chip8) exec0NNN() {
	//do nothing
}

func (vm *Chip8) exec00CN() {
	//TODO - scroll down n lines
	vm.pc += 2
}

func (vm *Chip8) exec00E0() {
	//clear gfx
	for i := range vm.Gfx {
		vm.Gfx[i] = 0
	}
	vm.pc += 2
	vm.Rg = true
}

func (vm *Chip8) exec00EE() {
	//return (set pc from stack, decrement sp)
	vm.pc = vm.stk[vm.sp] + 2
	vm.sp--
}

func (vm *Chip8) exec00FB() {
	//TODO - scroll right 4 pixels
	vm.pc += 2
}

func (vm *Chip8) exec00FC() {
	//TODO - scroll left 4 pixels
	vm.pc += 2
}

func (vm *Chip8) exec00FD() {
	//TODO - exit (reset/quit?)
	vm.pc += 2
}

func (vm *Chip8) exec00FE() {
	//lo-res mode
	vm.mode = sclr
	vm.pc += 2
}

func (vm *Chip8) exec00FF() {
	//hi-res mode
	vm.mode = schr
	vm.pc += 2
}

func (vm *Chip8) exec1NNN() {
	//goto nnn
	vm.pc = vm.oc & 0x0FFF
}

func (vm *Chip8) exec2NNN() {
	//call subroutine (increment sp, put current pc on stack, set pc to nnn)
	nnn := vm.oc & 0x0FFF
	vm.sp++
	vm.stk[vm.sp] = vm.pc
	vm.pc = nnn
}

func (vm *Chip8) exec3XNN() {
	//if(vx==nn) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	if vm.vr[x] == nn {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *Chip8) exec4XNN() {
	//if(vx!=nn) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	if vm.vr[x] != nn {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *Chip8) exec5XY0() {
	//if(vx==vy) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	if vm.vr[x] == vm.vr[y] {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *Chip8) exec6XNN() {
	//vx=nn
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	vm.vr[x] = nn
	vm.pc += 2
}

func (vm *Chip8) exec7XNN() {
	//vx+=nn
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	vm.vr[x] += nn
	vm.pc += 2
}

func (vm *Chip8) exec8XY0() {
	//vx=vy
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	vm.vr[x] = vm.vr[y]
	vm.pc += 2
}

func (vm *Chip8) exec8XY1() {
	//vx=vx|vy (bitwise OR)
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	vm.vr[x] = vm.vr[x] | vm.vr[y]
	vm.pc += 2
}

func (vm *Chip8) exec8XY2() {
	//vx=vx&vy (bitwise AND)
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	vm.vr[x] = vm.vr[x] & vm.vr[y]
	vm.pc += 2
}

func (vm *Chip8) exec8XY3() {
	//vx=vx^vy (bitwise XOR)
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	vm.vr[x] = vm.vr[x] ^ vm.vr[y]
	vm.pc += 2
}

func (vm *Chip8) exec8XY4() {
	//vx+=vy (only stores lowest 8 bits of result, if result > 0xFF then vF=1)
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	vm.vr[0xF] = 0
	//use vy>0xFF-vx as vx+vy>0xFF fails e.g. a:=byte(0xF0) b:=byte(0x15) a+b=5
	if vm.vr[y] > 0xFF-vm.vr[x] {
		vm.vr[0xF] = 1
	}
	vm.vr[x] += vm.vr[y]
	vm.pc += 2
}

func (vm *Chip8) exec8XY5() {
	//vx-=vy (if vx>vy then vF=1)
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	vm.vr[0xF] = 0
	if vm.vr[x] > vm.vr[y] {
		vm.vr[0xF] = 1
	}
	vm.vr[x] -= vm.vr[y]
	vm.pc += 2
}

func (vm *Chip8) exec8XY6() {
	//vx>>=1 (vF=the lsb of vx, then vx is divided by 2)
	x := vm.oc & 0x0F00 >> 8
	vm.vr[0xF] = vm.vr[x] & 0b00000001
	vm.vr[x] >>= 1
	vm.pc += 2
}

func (vm *Chip8) exec8XY7() {
	//vx=vy-vx (if vy>vx then vF=1)
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	vm.vr[0xF] = 0
	if vm.vr[y] > vm.vr[x] {
		vm.vr[0xF] = 1
	}
	vm.vr[x] = vm.vr[y] - vm.vr[x]
	vm.pc += 2
}

func (vm *Chip8) exec8XYE() {
	//vx<<=1 (vF=the msb of vx, then vx is multiplied by 2)
	x := vm.oc & 0x0F00 >> 8
	vm.vr[0xF] = vm.vr[x] & 0b10000000
	vm.vr[0xF] >>= 7
	vm.vr[x] <<= 1
	vm.pc += 2
}

func (vm *Chip8) exec9XY0() {
	//if(vx!=vy) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	if vm.vr[x] != vm.vr[y] {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *Chip8) execANNN() {
	//i=nnn
	vm.ir = vm.oc & 0x0FFF
	vm.pc += 2
}

func (vm *Chip8) execBNNN() {
	//pc=v0+nnn
	nnn := vm.oc & 0x0FFF
	vm.pc = uint16(vm.vr[0x0]) + nnn
}

func (vm *Chip8) execCXNN() {
	//vx=rand()&nn
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	vm.vr[x] = byte(rand.Intn(255)) & nn
	vm.pc += 2
}

func (vm *Chip8) execDXYN() {
	//TODO - call execDXYNHR/execDXYNLR dependent on mode (hi-res/lo-res)
	vm.execDXYNLR()
}

func (vm *Chip8) execDXYNLR() {
	//draw(x,y,n) - draw n byte sprite from mem[i] at vx,xy (vf=collision)
	/*
		Read n bytes (data) from memory, starting at i.
		Display bytes (data) as sprites on screen at coordinates vx,vy.
		Sprites are XORed onto the existing screen.
		If any pixels are erased, v[F] is set to 1, otherwise it is set to 0.
		Sprites wrap to opposite side of screen if they overlap an edge.
	*/
	vx := uint16(vm.vr[vm.oc&0x0F00>>8])
	vy := uint16(vm.vr[vm.oc&0x00F0>>4])
	n := vm.oc & 0x000F
	vm.vr[0xF] = 0
	//iterate over all of the sprite's rows
	for row := uint16(0); row < n; row++ {
		//get the byte for the current sprite row
		data := vm.Mem[vm.ir+row]
		//iterate over all cols in the current sprite row
		for col := uint16(0); col < 8; col++ {
			//calculate the gfx index for the current pixel (wrap if necessary)
			x := ((vx + col) % 64) * 2 //scale logical pix to gfx pix (64:128)
			y := ((vy + row) % 32) * 2 //scale logical pix to gfx pix (32:64)
			idx := (x + y*GfxW)
			//apply bitwise AND mask to extract the pixel's state from data
			if data&(0b10000000>>col) != 0 {
				//set v[F]=1 if pixel is to be erased
				if vm.Gfx[idx] == 1 {
					vm.vr[0xF] = 1
				}
				//bitwise XOR operation to toggle 2x2 gfx pixels state
				vm.Gfx[idx] ^= 1        //draw top-left gfx pixel
				vm.Gfx[idx+1] ^= 1      //draw top-right gfx pixel
				vm.Gfx[idx+GfxW] ^= 1   //draw bottom-left gfx pixel
				vm.Gfx[idx+GfxW+1] ^= 1 //draw bottom-right gfx pixel
			}
		}
	}
	vm.pc += 2
	vm.Rg = true
}

func (vm *Chip8) execDXYNHR() {
	//TODO - draw(x,y) - draw 16x16 sprite from mem[i] at vx,xy (vf=collision)
	vm.pc += 2
}

func (vm *Chip8) execEX9E() {
	//if the key stored in vx is pressed, skip next instruction
	x := vm.oc & 0x0F00 >> 8
	if vm.Key[vm.vr[x]] == 1 {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *Chip8) execEXA1() {
	//if the key stored in vx is not pressed, skip next instruction
	x := vm.oc & 0x0F00 >> 8
	if vm.Key[vm.vr[x]] == 0 {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *Chip8) execFX07() {
	//vx=delay_timer
	x := vm.oc & 0x0F00 >> 8
	vm.vr[x] = vm.dt
	vm.pc += 2
}

func (vm *Chip8) execFX0A() {
	//vx=get_key() (wait for a key press then store the key value in vx)
	x := vm.oc & 0x0F00 >> 8
	for i, k := range vm.Key {
		if k == 1 {
			vm.vr[x] = byte(i) //e.g. if R/0xD key is pressed, vr[x]=0xD
			vm.pc += 2         //execution stopped until a key is pressed
			break
		}
	}
}

func (vm *Chip8) execFX15() {
	//delay_timer=vx
	x := vm.oc & 0x0F00 >> 8
	vm.dt = vm.vr[x]
	vm.pc += 2
}

func (vm *Chip8) execFX18() {
	//sound_timer=vx
	x := vm.oc & 0x0F00 >> 8
	vm.St = vm.vr[x]
	vm.pc += 2
}

func (vm *Chip8) execFX1E() {
	// i+=vx
	x := vm.oc & 0x0F00 >> 8
	vm.ir += uint16(vm.vr[x])
	vm.pc += 2
}

func (vm *Chip8) execFX29() {
	//i=sprite_addr[vx] (point i at 5 byte font sprite for hex char at vx)
	x := vm.oc & 0x0F00 >> 8
	vm.ir = uint16(vm.vr[x] * 5)
	vm.pc += 2
}

func (vm *Chip8) execFX30() {
	//i=sprite_addr[vx] (point i at 10 byte font sprite for hex char at vx)
	x := vm.oc & 0x0F00 >> 8
	vm.ir = uint16(80 + vm.vr[x]*10)
	vm.pc += 2
}

func (vm *Chip8) execFX33() {
	//set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
	//store a decimal of vx in memory (e.g. if i=0 and vx=128, m0=1 m1=2 m2=8)
	x := vm.oc & 0x0F00 >> 8
	vm.Mem[vm.ir] = vm.vr[x] / 100         //e.g. 128/100=1
	vm.Mem[vm.ir+1] = (vm.vr[x] / 10) % 10 //e.g. 128/10=12, then 12%10=2
	vm.Mem[vm.ir+2] = vm.vr[x] % 10        //e.g. 128%10=8
	vm.pc += 2
}

func (vm *Chip8) execFX55() {
	//reg_dump(vx,&i) (store v0:vx inclusive from addr i, i is not modified)
	x := vm.oc & 0x0F00 >> 8
	for i := uint16(0); i <= x; i++ {
		vm.Mem[vm.ir+i] = vm.vr[i]
	}
	vm.pc += 2
}

func (vm *Chip8) execFX65() {
	//reg_load(vx,&i) (fill v0:vx inclusive from addr i, i is not modified)
	x := vm.oc & 0x0F00 >> 8
	for i := uint16(0); i <= x; i++ {
		vm.vr[i] = vm.Mem[vm.ir+i]
	}
	vm.pc += 2
}

func (vm *Chip8) execFX75() {
	//TODO - store v0:vx in rpl user flags (x <= 7)... investigate!
	vm.pc += 2
}

func (vm *Chip8) execFX85() {
	//TODO - read v0:vx from rpl user flags (x <= 7)... investigate!
	vm.pc += 2
}
