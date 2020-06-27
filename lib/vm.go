package lib

import (
	"fmt"
	"strings"
)

const (
	GfxW = 64 //gfx width
	GfxH = 32 //gfx height
)

var (
	fontset = []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	opcodeExecutors map[uint16]opcodeExecutor
)

type Chip8 struct {
	Mem [4096]byte //system memory (4kb total. 0x200-0xFFF: rom and ram)
	Gfx [2048]byte //vF is set upon pixel collision in draw instruction
	pc  uint16     //program counter
	oc  uint16     //current opcode (each opcode is two bytes long)
	vr  [16]byte   //v registers (v0-vE: general purpose. vF: carry flag)
	ir  uint16     //index register
	stk [16]uint16 //store program counter in stack before jump/gosub
	sp  uint16     //stack pointer to remember the level of stack used
	Key [16]byte   //stores the current state of the hex keypad (0-F)
	dt  byte       //delay timer counts down to zero at 60hz
	st  byte       //sound timer counts down to zero at 60hz
	Rg  bool       //redraw gfx - set by 0x00E0 (cls) and 0xDXYN (sprite)
}

type opcodeExecutor func()

func NewChip8() *Chip8 {
	vm := Chip8{
		pc: 0x0200,
		Rg: true,
	}
	for i := 0; i < len(fontset); i++ {
		vm.Mem[i] = fontset[i]
	}
	opcodeExecutors = map[uint16]opcodeExecutor{
		0x0000: vm.exec0NNN, 0x00E0: vm.exec00E0, 0x00EE: vm.exec00EE,
		0x1000: vm.exec1NNN, 0x2000: vm.exec2NNN, 0x3000: vm.exec3XNN,
		0x4000: vm.exec4XNN, 0x5000: vm.exec5XY0, 0x6000: vm.exec6XNN,
		0x7000: vm.exec7XNN, 0x8000: vm.exec8XY0, 0x8001: vm.exec8XY1,
		0x8002: vm.exec8XY2, 0x8003: vm.exec8XY3, 0x8004: vm.exec8XY4,
		0x8005: vm.exec8XY5, 0x8006: vm.exec8XY6, 0x8007: vm.exec8XY7,
		0x800E: vm.exec8XYE, 0x9000: vm.exec9XY0, 0xA000: vm.execANNN,
		0xB000: vm.execBNNN, 0xC000: vm.execCXNN, 0xD000: vm.execDXYN,
		0xE09E: vm.execEX9E, 0xE0A1: vm.execEXA1, 0xF007: vm.execFX07,
		0xF00A: vm.execFX0A, 0xF015: vm.execFX15, 0xF018: vm.execFX18,
		0xF01E: vm.execFX1E, 0xF029: vm.execFX29, 0xF033: vm.execFX33,
		0xF055: vm.execFX55, 0xF065: vm.execFX65,
	}
	return &vm
}

func (vm *Chip8) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("pc:%04X oc:%04X vr:", vm.pc, vm.oc))
	for i := 0; i < len(vm.vr); i++ {
		b.WriteString(fmt.Sprintf("%02X ", vm.vr[i]))
	}
	b.WriteString(fmt.Sprintf("ir:%04X key:", vm.ir))
	for i := 0; i < len(vm.Key); i++ {
		b.WriteString(fmt.Sprintf("%d", vm.Key[i]))
	}
	b.WriteString(fmt.Sprintf(" dt:%02X st:%02X rg:%t", vm.dt, vm.st, vm.Rg))
	return b.String()
}

func (vm *Chip8) EmulateCycle() {
	vm.oc = vm.fetchOpcode()
	d := vm.decodeOpcode()
	opcodeExecutors[d]()
	vm.updateTimers()
}

func (vm *Chip8) fetchOpcode() uint16 {
	/*
		Fetch opcode:
		Fetch and merge two bytes from memory locations pointed at by pc & pc+1
		e.g. memory[pc] = 0b10100010, memory[pc+1] = 0b11110000
		Convert first byte to uint16 and shift the bits left 8 times.
		e.g. 0b1010001000000000
		Use bitwise OR operation to merge the bytes.
		e.g. 0b1010001000000000 | 0b11110000 = 0b1010001011110000
	*/
	return uint16(vm.Mem[vm.pc])<<8 | uint16(vm.Mem[vm.pc+1])
}

func (vm *Chip8) decodeOpcode() uint16 {
	/*
		Decode opcode:
		Read the first 4 bits of the current opcode using bitwise AND operation
		e.g. 0x2105 & 0xF000 = 0x2000
		We can't always rely on just the first nibble to decode opcodes
		e.g. 0x00E0 and 0x00EE both start with 0x0
		In these cases we go on to compare the last nibble or byte...
		e.g. 0x00EE & 0x00FF = 0x00EE
	*/
	d := vm.oc & 0xF000
	if d == 0x0000 {
		switch vm.oc & 0x00FF {
		case 0x00E0:
			d = 0x00E0
		case 0x00EE:
			d = 0x00EE
		}
	}
	if d == 0x8000 {
		switch vm.oc & 0x000F {
		case 0x0000:
			d = 0x8000
		case 0x0001:
			d = 0x8001
		case 0x0002:
			d = 0x8002
		case 0x0003:
			d = 0x8003
		case 0x0004:
			d = 0x8004
		case 0x0005:
			d = 0x8005
		case 0x0006:
			d = 0x8006
		case 0x0007:
			d = 0x8007
		case 0x000E:
			d = 0x800E
		}
	}
	if d == 0xE000 {
		switch vm.oc & 0x00FF {
		case 0x009E:
			d = 0xE09E
		case 0x00A1:
			d = 0xE0A1
		}
	}
	if d == 0xF000 {
		switch vm.oc & 0x00FF {
		case 0x0007:
			d = 0xF007
		case 0x000A:
			d = 0xF00A
		case 0x0015:
			d = 0xF015
		case 0x0018:
			d = 0xF018
		case 0x001E:
			d = 0xF01E
		case 0x0029:
			d = 0xF029
		case 0x0033:
			d = 0xF033
		case 0x0055:
			d = 0xF055
		case 0x0065:
			d = 0xF065
		}
	}
	return d
}

func (vm *Chip8) updateTimers() {
	//Count down to zero
	if vm.dt > 0 {
		vm.dt--
	}
	if vm.st > 0 {
		vm.st--
	}
}
