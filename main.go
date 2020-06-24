package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	fullscreen = false
	screenW    = 960
	screenH    = 480

	gfxW   = 64
	gfxH   = 32
	pixelW = float64(screenW / gfxW)
	pixelH = float64(screenH / gfxH)
)

var (
	opcode     uint16     //current opcode (each opcode is two bytes long)
	memory     [4096]byte //system memory (4kb total. 0x200-0xFFF: rom and ram)
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

	col0    = color.NRGBA{0x00, 0x00, 0x00, 0xff}
	col1    = color.NRGBA{0x00, 0xff, 0x00, 0xff}
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
)

type Game struct {
}

func NewGame() *Game {
	g := &Game{}
	return g
}

func init() {
	opcode = 0
	memory = [4096]byte{}
	for i := 0; i < len(fontset); i++ {
		memory[i] = fontset[i]
	}
	v = [16]byte{}
	i = 0
	pc = 0x200
	gfx = [gfxW * gfxH]byte{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	delayTimer = 0
	soundTimer = 0
	stack = [16]uint16{}
	sp = 0
	keys = [16]byte{}
	render = true
}

func loadRom() {
	path := flag.String("path", "./rom/test/test_opcode.ch8", "path to rom file")
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
		memory[i+0x200] = bytes[i]
	}

	fmt.Println(memory)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func (g *Game) Update(screen *ebiten.Image) error {
	if memory[0x200] == 0 {
		loadRom()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	emulateCycle()

	updateKeys()

	return nil
}

func emulateCycle() {
	/*
		Fetch opcode:
		Fetch and merge two bytes from memory locations pointed at by pc & pc+1
		e.g. memory[pc] = 0b10100010, memory[pc+1] = 0b11110000
		Convert first byte to uint16 and shift the bits left 8 times.
		e.g. 0b1010001000000000
		Use bitwise OR operation to merge the bytes.
		e.g. 0b1010001000000000 | 0b11110000 = 0b1010001011110000
	*/
	opcode = uint16(memory[pc])<<8 | uint16(memory[pc+1])

	/*
		Decode opcode:
		Read the first 4 bits of the current opcode using bitwise AND operation
		e.g. 0x2105 & 0xF000 = 0x2000
		We can't always rely on just the first nibble to decode opcodes
		e.g. 0x00E0 and 0x00EE both start with 0x0
		In these cases we go on to compare the last nibble or byte...
		e.g. 0x00EE & 0x00FF = 0x00EE
		Then use bitwise OR to merge and arrive at the final decoded opcode
		e.g. 0x0000 | 0x00EE = 0x00EE
	*/
	decoded := opcode & 0xF000
	if decoded == 0x8000 {
		lastNibble := opcode & 0x000F
		decoded = decoded | lastNibble
	}
	if decoded == 0x0000 || decoded == 0xE000 || decoded == 0xF000 {
		lastByte := opcode & 0x00FF
		decoded = decoded | lastByte
	}

	//Execute opcode
	switch opcode {
	case 0x0000:
		//TODO ?
	case 0x00E0:
		//disp_clear()
		gfx = [2048]byte{}
		pc += 2
		render = true
	case 0x00EE:
		//TODO return;
	case 0x1000:
		//TODO goto NNN;
	case 0x2000:
		//TODO *(0xNNN)()
	case 0x3000:
		//TODO if(Vx==NN)
	case 0x4000:
		//TODO if(Vx!=NN)
	case 0x5000:
		//TODO if(Vx==Vy)
	case 0x6000:
		//TODO Vx=NN
	case 0x7000:
		//TODO Vx+=NN
	case 0x8000:
		//TODO Vx=Vy
	case 0x8001:
		//TODO Vx=Vx|Vy
	case 0x8002:
		//TODO Vx=Vx&Vy
	case 0x8003:
		//TODO Vx=Vx^Vy
	case 0x8004:
		//TODO Vx+=Vy
	case 0x8005:
		//TODO Vx-=Vy
	case 0x8006:
		//TODO Vx>>=1
	case 0x8007:
		//TODO Vx=Vy-Vx
	case 0x800E:
		//TODO Vx<<=1
	case 0x9000:
		//TODO if(Vx!=Vy)
	case 0xA000:
		//TODO I=NNN
	case 0xB000:
		//TODO PC=V0+NNN
	case 0xC000:
		//TODO Vx=rand()&NN
	case 0xD000:
		//TODO draw(Vx,Vy,N)
	case 0xE09E:
		//TODO if(key()==Vx)
	case 0xE0A1:
		//TODO if(key()!=Vx)
	case 0xF007:
		//TODO Vx=get_delay()
	case 0xF00A:
		//TODO Vx=get_key()
	case 0xF015:
		//TODO delay_timer(Vx)
	case 0xF018:
		//TODO sound_timer(Vx)
	case 0xF01E:
		//TODO I+=Vx
	case 0xF029:
		//TODO I=sprite_addr[Vx]
	case 0xF033:
		//TODO set_BCD(Vx);*(I+0)=BCD(3);*(I+1)=BCD(2);*(I+2)=BCD(1);
	case 0xF055:
		//TODO reg_dump(Vx,&I)
	case 0xF065:
		//TODO reg_load(Vx,&I)
	}

	//Update timers
	if delayTimer > 0 {
		delayTimer--
	}
	if soundTimer > 0 {
		soundTimer--
	}
}

func updateKeys() {
	keys[0] = boolToByte(ebiten.IsKeyPressed(ebiten.Key1))
	keys[1] = boolToByte(ebiten.IsKeyPressed(ebiten.Key2))
	keys[2] = boolToByte(ebiten.IsKeyPressed(ebiten.Key3))
	keys[3] = boolToByte(ebiten.IsKeyPressed(ebiten.Key4))
	keys[4] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyQ))
	keys[5] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyW))
	keys[6] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyE))
	keys[7] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyR))
	keys[8] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyA))
	keys[9] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyS))
	keys[10] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyD))
	keys[11] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyF))
	keys[12] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyZ))
	keys[13] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyX))
	keys[14] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyC))
	keys[15] = boolToByte(ebiten.IsKeyPressed(ebiten.KeyV))
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func (g *Game) Draw(screen *ebiten.Image) {
	if render {
		screen.Fill(col0)
		for y := 0; y < gfxH; y++ {
			for x := 0; x < gfxW; x++ {
				if gfx[y*gfxW+x] == 1 {
					ebitenutil.DrawRect(screen, float64(x)*pixelW, float64(y)*pixelH, pixelW, pixelH, col1)
				}
			}
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	//fmt.Println(keys)
}

func main() {
	ebiten.SetFullscreen(fullscreen)
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("CHIP-8")
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
