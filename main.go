package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

//--- client

const (
	fullscreen     = false
	screenW        = 960
	screenH        = 480
	ticksPerSecond = 60

	gfxW   = 64
	gfxH   = 32
	pixelW = float64(screenW / gfxW)
	pixelH = float64(screenH / gfxH)
)

var (
	bg     = color.NRGBA{0x00, 0x00, 0x00, 0xff}
	fg     = color.NRGBA{0x00, 0xff, 0x00, 0xff}
	keyMap = map[ebiten.Key]uint16{
		ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
		ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
		ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
	}
	loaded = false
)

type Game struct {
	vm *chip8
}

func NewGame() *Game {
	return &Game{
		vm: newChip8(),
	}
}

func init() {
	ebiten.SetMaxTPS(ticksPerSecond)
	rand.Seed(time.Now().UnixNano())
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !loaded {
		loadRom(g.vm)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	g.vm.emulateCycle()
	updateKeys(g.vm)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.vm.render {
		screen.Fill(bg)
		for y := 0; y < gfxH; y++ {
			for x := 0; x < gfxW; x++ {
				if g.vm.gfx[y*gfxW+x] == 1 {
					ebitenutil.DrawRect(screen, float64(x)*pixelW, float64(y)*pixelH, pixelW, pixelH, fg)
				}
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	//fmt.Println(keys)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func loadRom(vm *chip8) {
	path := flag.String("path", "./rom/test/ti360.ch8", "path to rom file")
	flag.Parse()

	rom, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(rom)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(bytes); i++ {
		vm.mem[i+0x200] = bytes[i]
	}
	fmt.Println(vm.mem)

	loaded = true
}

func updateKeys(vm *chip8) {
	for phys, virt := range keyMap {
		vm.keys[virt] = boolToByte(ebiten.IsKeyPressed(phys))
	}
}

func main() {
	ebiten.SetFullscreen(fullscreen)
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("CHIP-8 by TheInvader360")
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}

//----- util

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

//----- vm

type chip8 struct {
	opcode     uint16     //current opcode (each opcode is two bytes long)
	mem        [4096]byte //system memory (4kb total. 0x200-0xFFF: rom and ram)
	v          [16]byte   //registers (v0-vE: general purpose. vF: carry flag)
	i          uint16     //index register
	pc         uint16     //program counter
	gfx        [2048]byte //vF is set upon pixel collision in draw instruction
	delayTimer byte       //counts down to zero at 60hz
	soundTimer byte       //counts down to zero at 60hz
	stack      [16]uint16 //store program counter in stack before jump/gosub
	sp         uint16     //stack pointer to remember the level of stack used
	keys       [16]byte   //stores the current state of the hex keypad (0-F)
	render     bool       //set by 0x00E0 (cls) and 0xDXYN (draw sprite)
}

type opcodeExecutor func() string

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

func newChip8() *chip8 {
	vm := chip8{
		pc:     0x0200,
		render: true,
	}
	for i := 0; i < len(fontset); i++ {
		vm.mem[i] = fontset[i]
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

func (vm *chip8) emulateCycle() {
	vm.opcode = vm.fetchOpcode()
	decoded := vm.decodeOpcode()
	fmt.Println(opcodeExecutors[decoded]())
	vm.updateTimers()
}

func (vm *chip8) fetchOpcode() uint16 {
	/*
		Fetch opcode:
		Fetch and merge two bytes from memory locations pointed at by pc & pc+1
		e.g. memory[pc] = 0b10100010, memory[pc+1] = 0b11110000
		Convert first byte to uint16 and shift the bits left 8 times.
		e.g. 0b1010001000000000
		Use bitwise OR operation to merge the bytes.
		e.g. 0b1010001000000000 | 0b11110000 = 0b1010001011110000
	*/
	return uint16(vm.mem[vm.pc])<<8 | uint16(vm.mem[vm.pc+1])
}

func (vm *chip8) decodeOpcode() uint16 {
	/*
		Decode opcode:
		Read the first 4 bits of the current opcode using bitwise AND operation
		e.g. 0x2105 & 0xF000 = 0x2000
		We can't always rely on just the first nibble to decode opcodes
		e.g. 0x00E0 and 0x00EE both start with 0x0
		In these cases we go on to compare the last nibble or byte...
		e.g. 0x00EE & 0x00FF = 0x00EE
	*/
	decoded := vm.opcode & 0xF000
	if decoded == 0x0000 {
		switch vm.opcode & 0x00FF {
		case 0x00E0:
			decoded = 0x00E0
		case 0x00EE:
			decoded = 0x00EE
		}
	}
	if decoded == 0x8000 {
		switch vm.opcode & 0x000F {
		case 0x0000:
			decoded = 0x8000
		case 0x0001:
			decoded = 0x8001
		case 0x0002:
			decoded = 0x8002
		case 0x0003:
			decoded = 0x8003
		case 0x0004:
			decoded = 0x8004
		case 0x0005:
			decoded = 0x8005
		case 0x0006:
			decoded = 0x8006
		case 0x0007:
			decoded = 0x8007
		case 0x000E:
			decoded = 0x800E
		}
	}
	if decoded == 0xE000 {
		switch vm.opcode & 0x00FF {
		case 0x009E:
			decoded = 0xE09E
		case 0x00A1:
			decoded = 0xE0A1
		}
	}
	if decoded == 0xF000 {
		switch vm.opcode & 0x00FF {
		case 0x0007:
			decoded = 0xF007
		case 0x000A:
			decoded = 0xF00A
		case 0x0015:
			decoded = 0xF015
		case 0x0018:
			decoded = 0xF018
		case 0x001E:
			decoded = 0xF01E
		case 0x0029:
			decoded = 0xF029
		case 0x0033:
			decoded = 0xF033
		case 0x0055:
			decoded = 0xF055
		case 0x0065:
			decoded = 0xF065
		}
	}
	return decoded
}

func (vm *chip8) updateTimers() {
	//Count down to zero
	if vm.delayTimer > 0 {
		vm.delayTimer--
	}
	if vm.soundTimer > 0 {
		vm.soundTimer--
	}
}

func (vm *chip8) exec0NNN() string {
	//not implemented
	return fmt.Sprintf("exec0NNN 0x%04X: pc=0x%04X (not implemented)", vm.opcode, vm.pc)
}

func (vm *chip8) exec00E0() string {
	//disp_clear()
	for i := range vm.gfx {
		vm.gfx[i] = 0
	}
	vm.pc += 2
	vm.render = true
	return fmt.Sprintf("exec00E0 0x%04X: pc=0x%04X gfx={cleared} render=%t", vm.opcode, vm.pc, vm.render)
}

func (vm *chip8) exec00EE() string {
	//TODO return
	return fmt.Sprintf("exec00EE 0x%04X", vm.opcode)
}

func (vm *chip8) exec1NNN() string {
	//goto nnn
	vm.pc = vm.opcode & 0x0FFF
	return fmt.Sprintf("exec1NNN 0x%04X: pc=0x%04X (goto)", vm.opcode, vm.pc)
}

func (vm *chip8) exec2NNN() string {
	//call subroutine (increment sp, put current pc on stack, set pc to nnn)
	nnn := vm.opcode & 0x0FFF
	vm.sp++
	vm.stack[vm.sp] = vm.pc
	vm.pc = nnn
	return fmt.Sprintf("exec2NNN 0x%04X: pc=0x%04X sp=0x%04X", vm.opcode, vm.pc, vm.sp)
}

func (vm *chip8) exec3XNN() string {
	//if(vx==nn) skip next instruction
	x := vm.opcode & 0x0F00 >> 8
	nn := byte(vm.opcode & 0x00FF)
	skip := false
	if vm.v[x] == nn {
		skip = true
		vm.pc += 2
	}
	vm.pc += 2
	return fmt.Sprintf("exec3XNN 0x%04X: pc=0x%04X {skip=%t}", vm.opcode, vm.pc, skip)
}

func (vm *chip8) exec4XNN() string {
	//if(vx!=nn) skip next instruction
	x := vm.opcode & 0x0F00 >> 8
	nn := byte(vm.opcode & 0x00FF)
	skip := false
	if vm.v[x] != nn {
		skip = true
		vm.pc += 2
	}
	vm.pc += 2
	return fmt.Sprintf("exec4XNN 0x%04X: pc=0x%04X {skip=%t}", vm.opcode, vm.pc, skip)
}

func (vm *chip8) exec5XY0() string {
	//if(vx==vy) skip next instruction
	x := vm.opcode & 0x0F00 >> 8
	y := vm.opcode & 0x00F0 >> 4
	skip := false
	if vm.v[x] == vm.v[y] {
		skip = true
		vm.pc += 2
	}
	vm.pc += 2
	return fmt.Sprintf("exec5XY0 0x%04X: pc=0x%04X {skip=%t}", vm.opcode, vm.pc, skip)
}

func (vm *chip8) exec6XNN() string {
	//vx=nn
	x := vm.opcode & 0x0F00 >> 8
	nn := byte(vm.opcode & 0x00FF)
	vm.v[x] = nn
	vm.pc += 2
	return fmt.Sprintf("exec6XNN 0x%04X: pc=0x%04X v[%01X]=%02X", vm.opcode, vm.pc, x, vm.v[x])
}

func (vm *chip8) exec7XNN() string {
	//vx+=nn
	x := vm.opcode & 0x0F00 >> 8
	nn := byte(vm.opcode & 0x00FF)
	vm.v[x] += nn
	vm.pc += 2
	return fmt.Sprintf("exec7XNN 0x%04X: pc=0x%04X v[%01X]=%02X", vm.opcode, vm.pc, x, vm.v[x])
}

func (vm *chip8) exec8XY0() string {
	//TODO vx=vy
	return fmt.Sprintf("exec8XY0 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XY1() string {
	//TODO vx=vx|vy
	return fmt.Sprintf("exec8XY1 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XY2() string {
	//TODO vx=vx&vy
	return fmt.Sprintf("exec8XY2 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XY3() string {
	//TODO vx=vx^vy
	return fmt.Sprintf("exec8XY3 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XY4() string {
	//TODO vx+=vy
	return fmt.Sprintf("exec8XY4 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XY5() string {
	//TODO vx-=vy
	return fmt.Sprintf("exec8XY5 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XY6() string {
	//TODO vx>>=1
	return fmt.Sprintf("exec8XY6 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XY7() string {
	//TODO vx=vy-vx
	return fmt.Sprintf("exec8XY7 0x%04X", vm.opcode)
}

func (vm *chip8) exec8XYE() string {
	//TODO vx<<=1
	return fmt.Sprintf("exec8XYE 0x%04X", vm.opcode)
}

func (vm *chip8) exec9XY0() string {
	//if(vx!=vy) skip next instruction
	x := vm.opcode & 0x0F00 >> 8
	y := vm.opcode & 0x00F0 >> 4
	skip := false
	if vm.v[x] != vm.v[y] {
		skip = true
		vm.pc += 2
	}
	vm.pc += 2
	return fmt.Sprintf("exec9XY0 0x%04X: pc=0x%04X {skip=%t}", vm.opcode, vm.pc, skip)
}

func (vm *chip8) execANNN() string {
	//i=nnn
	vm.i = vm.opcode & 0x0FFF
	vm.pc += 2
	return fmt.Sprintf("execANNN 0x%04X: pc=0x%04X i=0x%04X", vm.opcode, vm.pc, vm.i)
}

func (vm *chip8) execBNNN() string {
	//TODO pc=v0+nnn
	return fmt.Sprintf("execBNNN 0x%04X", vm.opcode)
}

func (vm *chip8) execCXNN() string {
	//vx=rand()&nn
	x := vm.opcode & 0x0F00 >> 8
	nn := byte(vm.opcode & 0x00FF)
	vm.v[x] = byte(rand.Intn(255)) & nn
	vm.pc += 2
	return fmt.Sprintf("execCXNN 0x%04X: pc=0x%04X v[%01X]=%02X", vm.opcode, vm.pc, x, vm.v[x])
}

func (vm *chip8) execDXYN() string {
	//draw(vx,vy,n)
	/*
		Read n bytes (data) from memory, starting at i.
		Display bytes (data) as sprites on screen at coordinates vx,vy.
		Sprites are XORed onto the existing screen.
		If any pixels are erased, v[F] is set to 1, otherwise it is set to 0.
		Sprites wrap to opposite side of screen if they overlap an edge.
	*/
	vx := uint16(vm.v[vm.opcode&0x0F00>>8])
	vy := uint16(vm.v[vm.opcode&0x00F0>>4])
	n := vm.opcode & 0x000F
	vm.v[0xF] = 0
	//iterate over all of the sprite's rows
	for row := uint16(0); row < n; row++ {
		//get the byte for the current row
		data := vm.mem[vm.i+row]
		//iterate over all of the current row's cols
		for col := uint16(0); col < 8; col++ {
			//calculate the gfx index for the current row and col
			idx := ((vy+row)*gfxW + vx + col)
			//apply bitwise AND mask to extract a single pixel from data
			if data&(0b10000000>>col) != 0 {
				//TODO: confirm that out of bounds draw really should wrap!
				if idx > uint16(len(vm.gfx)) {
					idx -= uint16(len(vm.gfx))
				}
				//set v[F] if pixel is to be erased
				if vm.gfx[idx] == 1 {
					vm.v[0xF] = 1
				}
				//bitwise XOR operation to toggle pixel value
				vm.gfx[idx] ^= 1
			}
		}
	}
	vm.pc += 2
	vm.render = true
	return fmt.Sprintf("execDXYN 0x%04X: pc=0x%04X gfx={updated} render=%t", vm.opcode, vm.pc, vm.render)
}

func (vm *chip8) execEX9E() string {
	//if the key stored in vx is pressed, skip next instruction
	x := vm.opcode & 0x0F00 >> 8
	skip := false
	if vm.keys[vm.v[x]] == 1 {
		skip = true
		vm.pc += 2
	}
	vm.pc += 2
	return fmt.Sprintf("execEX9E 0x%04X: pc=0x%04X {skip=%t}", vm.opcode, vm.pc, skip)
}

func (vm *chip8) execEXA1() string {
	//if the key stored in vx is not pressed, skip next instruction
	x := vm.opcode & 0x0F00 >> 8
	skip := false
	if vm.keys[vm.v[x]] == 0 {
		skip = true
		vm.pc += 2
	}
	vm.pc += 2
	return fmt.Sprintf("execEXA1 0x%04X: pc=0x%04X {skip=%t}", vm.opcode, vm.pc, skip)
}

func (vm *chip8) execFX07() string {
	//vx=delay_timer
	x := vm.opcode & 0x0F00 >> 8
	vm.v[x] = vm.delayTimer
	vm.pc += 2
	return fmt.Sprintf("execFX07 0x%04X: pc=0x%04X v[%01X]=%02X", vm.opcode, vm.pc, x, vm.v[x])
}

func (vm *chip8) execFX0A() string {
	//TODO vx=get_key()
	return fmt.Sprintf("execFX0A 0x%04X", vm.opcode)
}

func (vm *chip8) execFX15() string {
	//delay_timer=vx
	x := vm.opcode & 0x0F00 >> 8
	vm.delayTimer = vm.v[x]
	vm.pc += 2
	return fmt.Sprintf("execFX15 0x%04X: pc=0x%04X delayTimer=%02X", vm.opcode, vm.pc, vm.delayTimer)
}

func (vm *chip8) execFX18() string {
	//TODO sound_timer=vx
	return fmt.Sprintf("execFX18 0x%04X", vm.opcode)
}

func (vm *chip8) execFX1E() string {
	//TODO i+=vx
	return fmt.Sprintf("execFX1E 0x%04X", vm.opcode)
}

func (vm *chip8) execFX29() string {
	//TODO i=sprite_addr[vx]
	return fmt.Sprintf("execFX29 0x%04X", vm.opcode)
}

func (vm *chip8) execFX33() string {
	//TODO set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
	return fmt.Sprintf("execFX33 0x%04X", vm.opcode)
}

func (vm *chip8) execFX55() string {
	//TODO reg_dump(vx,&i)
	return fmt.Sprintf("execFX55 0x%04X", vm.opcode)
}

func (vm *chip8) execFX65() string {
	//TODO reg_load(vx,&i)
	return fmt.Sprintf("execFX65 0x%04X", vm.opcode)
}
