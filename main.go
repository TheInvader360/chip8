package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

//--- client

const (
	scrF = false //client fullscreen
	scrW = 960   //client screen width in pixels
	scrH = 480   //client screen height in pixels
	tps  = 60    //client max ticks per second

	gfxW = 64                   //chip8 gfx width
	gfxH = 32                   //chip8 gfx height
	pixW = float64(scrW / gfxW) //pixel width scale factor
	pixH = float64(scrH / gfxH) //pixel height scale factor
)

var (
	bg = color.NRGBA{0x00, 0x00, 0x00, 0xff}
	fg = color.NRGBA{0x00, 0xff, 0x00, 0xff}
	km = map[ebiten.Key]uint16{
		ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
		ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
		ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
	}
	rom = false
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
	ebiten.SetMaxTPS(tps)
	rand.Seed(time.Now().UnixNano())
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !rom {
		loadRom(g.vm)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	g.vm.emulateCycle()
	fmt.Println(g.vm)
	updateKeys(g.vm)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.vm.rg {
		screen.Fill(bg)
		for y := 0; y < gfxH; y++ {
			for x := 0; x < gfxW; x++ {
				if g.vm.gfx[y*gfxW+x] == 1 {
					ebitenutil.DrawRect(screen,
						float64(x)*pixW, float64(y)*pixH, pixW, pixH, fg)
				}
			}
		}
		g.vm.rg = false
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return scrW, scrH
}

func loadRom(vm *chip8) {
	path := flag.String("path", "./rom/test/ti360.ch8", "path to rom file")
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(bytes); i++ {
		vm.mem[i+0x200] = bytes[i]
	}
	fmt.Println(vm.mem)

	rom = true
}

func updateKeys(vm *chip8) {
	for phys, virt := range km {
		vm.key[virt] = boolToByte(ebiten.IsKeyPressed(phys))
	}
}

func main() {
	ebiten.SetFullscreen(scrF)
	ebiten.SetWindowSize(scrW, scrH)
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
	mem [4096]byte //system memory (4kb total. 0x200-0xFFF: rom and ram)
	gfx [2048]byte //vF is set upon pixel collision in draw instruction
	pc  uint16     //program counter
	oc  uint16     //current opcode (each opcode is two bytes long)
	vr  [16]byte   //v registers (v0-vE: general purpose. vF: carry flag)
	ir  uint16     //index register
	stk [16]uint16 //store program counter in stack before jump/gosub
	sp  uint16     //stack pointer to remember the level of stack used
	key [16]byte   //stores the current state of the hex keypad (0-F)
	dt  byte       //delay timer counts down to zero at 60hz
	st  byte       //sound timer counts down to zero at 60hz
	rg  bool       //redraw gfx - set by 0x00E0 (cls) and 0xDXYN (sprite)
}

type opcodeExecutor func()

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
		pc: 0x0200,
		rg: true,
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

func (vm *chip8) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("pc:%04X oc:%04X vr:", vm.pc, vm.oc))
	for i := 0; i < len(vm.vr); i++ {
		b.WriteString(fmt.Sprintf("%02X ", vm.vr[i]))
	}
	b.WriteString(fmt.Sprintf("ir:%04X key:", vm.ir))
	for i := 0; i < len(vm.key); i++ {
		b.WriteString(fmt.Sprintf("%d", vm.key[i]))
	}
	b.WriteString(fmt.Sprintf(" dt:%02X st:%02X rg:%t", vm.dt, vm.st, vm.rg))
	return b.String()
}

func (vm *chip8) emulateCycle() {
	vm.oc = vm.fetchOpcode()
	d := vm.decodeOpcode()
	opcodeExecutors[d]()
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

func (vm *chip8) updateTimers() {
	//Count down to zero
	if vm.dt > 0 {
		vm.dt--
	}
	if vm.st > 0 {
		vm.st--
	}
}

func (vm *chip8) exec0NNN() {
	//do nothing
}

func (vm *chip8) exec00E0() {
	//clear gfx
	for i := range vm.gfx {
		vm.gfx[i] = 0
	}
	vm.pc += 2
	vm.rg = true
}

func (vm *chip8) exec00EE() {
	//TODO return
}

func (vm *chip8) exec1NNN() {
	//goto nnn
	vm.pc = vm.oc & 0x0FFF
}

func (vm *chip8) exec2NNN() {
	//call subroutine (increment sp, put current pc on stack, set pc to nnn)
	nnn := vm.oc & 0x0FFF
	vm.sp++
	vm.stk[vm.sp] = vm.pc
	vm.pc = nnn
}

func (vm *chip8) exec3XNN() {
	//if(vx==nn) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	if vm.vr[x] == nn {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *chip8) exec4XNN() {
	//if(vx!=nn) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	if vm.vr[x] != nn {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *chip8) exec5XY0() {
	//if(vx==vy) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	if vm.vr[x] == vm.vr[y] {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *chip8) exec6XNN() {
	//vx=nn
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	vm.vr[x] = nn
	vm.pc += 2
}

func (vm *chip8) exec7XNN() {
	//vx+=nn
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	vm.vr[x] += nn
	vm.pc += 2
}

func (vm *chip8) exec8XY0() {
	//TODO vx=vy
}

func (vm *chip8) exec8XY1() {
	//TODO vx=vx|vy
}

func (vm *chip8) exec8XY2() {
	//TODO vx=vx&vy
}

func (vm *chip8) exec8XY3() {
	//TODO vx=vx^vy
}

func (vm *chip8) exec8XY4() {
	//TODO vx+=vy
}

func (vm *chip8) exec8XY5() {
	//TODO vx-=vy
}

func (vm *chip8) exec8XY6() {
	//TODO vx>>=1
}

func (vm *chip8) exec8XY7() {
	//TODO vx=vy-vx
}

func (vm *chip8) exec8XYE() {
	//TODO vx<<=1
}

func (vm *chip8) exec9XY0() {
	//if(vx!=vy) skip next instruction
	x := vm.oc & 0x0F00 >> 8
	y := vm.oc & 0x00F0 >> 4
	if vm.vr[x] != vm.vr[y] {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *chip8) execANNN() {
	//i=nnn
	vm.ir = vm.oc & 0x0FFF
	vm.pc += 2
}

func (vm *chip8) execBNNN() {
	//TODO pc=v0+nnn
}

func (vm *chip8) execCXNN() {
	//vx=rand()&nn
	x := vm.oc & 0x0F00 >> 8
	nn := byte(vm.oc & 0x00FF)
	vm.vr[x] = byte(rand.Intn(255)) & nn
	vm.pc += 2
}

func (vm *chip8) execDXYN() {
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
		data := vm.mem[vm.ir+row]
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
					vm.vr[0xF] = 1
				}
				//bitwise XOR operation to toggle pixel value
				vm.gfx[idx] ^= 1
			}
		}
	}
	vm.pc += 2
	vm.rg = true
}

func (vm *chip8) execEX9E() {
	//if the key stored in vx is pressed, skip next instruction
	x := vm.oc & 0x0F00 >> 8
	if vm.key[vm.vr[x]] == 1 {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *chip8) execEXA1() {
	//if the key stored in vx is not pressed, skip next instruction
	x := vm.oc & 0x0F00 >> 8
	if vm.key[vm.vr[x]] == 0 {
		vm.pc += 2
	}
	vm.pc += 2
}

func (vm *chip8) execFX07() {
	//vx=delay_timer
	x := vm.oc & 0x0F00 >> 8
	vm.vr[x] = vm.dt
	vm.pc += 2
}

func (vm *chip8) execFX0A() {
	//TODO vx=get_key()
}

func (vm *chip8) execFX15() {
	//delay_timer=vx
	x := vm.oc & 0x0F00 >> 8
	vm.dt = vm.vr[x]
	vm.pc += 2
}

func (vm *chip8) execFX18() {
	//TODO sound_timer=vx
}

func (vm *chip8) execFX1E() {
	//TODO i+=vx
}

func (vm *chip8) execFX29() {
	//TODO i=sprite_addr[vx]
}

func (vm *chip8) execFX33() {
	//TODO set_bcd(vx);*(i+0)=bcd(3);*(i+1)=bcd(2);*(i+2)=bcd(1);
}

func (vm *chip8) execFX55() {
	//TODO reg_dump(vx,&i)
}

func (vm *chip8) execFX65() {
	//TODO reg_load(vx,&i)
}
