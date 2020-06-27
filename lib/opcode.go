package lib

import (
	"math/rand"
)

func (vm *Chip8) exec0NNN() {
	//do nothing
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
	//TODO return
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
	//TODO vx=vy
}

func (vm *Chip8) exec8XY1() {
	//TODO vx=vx|vy
}

func (vm *Chip8) exec8XY2() {
	//TODO vx=vx&vy
}

func (vm *Chip8) exec8XY3() {
	//TODO vx=vx^vy
}

func (vm *Chip8) exec8XY4() {
	//TODO vx+=vy
}

func (vm *Chip8) exec8XY5() {
	//TODO vx-=vy
}

func (vm *Chip8) exec8XY6() {
	//TODO vx>>=1
}

func (vm *Chip8) exec8XY7() {
	//TODO vx=vy-vx
}

func (vm *Chip8) exec8XYE() {
	//TODO vx<<=1
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
	//TODO pc=v0+nnn
}

func (vm *Chip8) execCXNN() {
	//vx=rand()&nn
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	vm.vr[x] = byte(rand.Intn(255)) & nn
	vm.pc += 2
}

func (vm *Chip8) execDXYN() {
	//draw(vx,vy,n)
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
		//get the byte for the current row
		data := vm.Mem[vm.ir+row]
		//iterate over all of the current row's cols
		for col := uint16(0); col < 8; col++ {
			//calculate the gfx index for the current row and col
			idx := ((vy+row)*GfxW + vx + col)
			//apply bitwise AND mask to extract a single pixel from data
			if data&(0b10000000>>col) != 0 {
				//TODO: confirm that out of bounds draw really should wrap!
				if idx > uint16(len(vm.Gfx)) {
					idx -= uint16(len(vm.Gfx))
				}
				//set v[F] if pixel is to be erased
				if vm.Gfx[idx] == 1 {
					vm.vr[0xF] = 1
				}
				//bitwise XOR operation to toggle pixel value
				vm.Gfx[idx] ^= 1
			}
		}
	}
	vm.pc += 2
	vm.Rg = true
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
	//TODO vx=get_key()
}

func (vm *Chip8) execFX15() {
	//delay_timer=vx
	x := vm.oc & 0x0F00 >> 8
	vm.dt = vm.vr[x]
	vm.pc += 2
}

func (vm *Chip8) execFX18() {
	//TODO sound_timer=vx
}

func (vm *Chip8) execFX1E() {
	//TODO i+=vx
}

func (vm *Chip8) execFX29() {
	//TODO i=sprite_addr[vx]
}

func (vm *Chip8) execFX33() {
	//TODO set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
}

func (vm *Chip8) execFX55() {
	//TODO reg_dump(vx,&i)
}

func (vm *Chip8) execFX65() {
	//TODO reg_load(vx,&i)
}
