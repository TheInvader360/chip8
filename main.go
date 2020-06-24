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

type Game struct{}

type opcodeExecutor func() string

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
	loaded     bool       //set when rom is loaded into memory

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
	opcodeExecutors = map[uint16]opcodeExecutor{
		0x0000: exec0NNN, 0x00E0: exec00E0, 0x00EE: exec00EE, 0x1000: exec1NNN,
		0x2000: exec2NNN, 0x3000: exec3XNN, 0x4000: exec4XNN, 0x5000: exec5XY0,
		0x6000: exec6XNN, 0x7000: exec7XNN, 0x8000: exec8XY0, 0x8001: exec8XY1,
		0x8002: exec8XY2, 0x8003: exec8XY3, 0x8004: exec8XY4, 0x8005: exec8XY5,
		0x8006: exec8XY6, 0x8007: exec8XY7, 0x800E: exec8XYE, 0x9000: exec9XY0,
		0xA000: execANNN, 0xB000: execBNNN, 0xC000: execCXNN, 0xD000: execDXYN,
		0xE09E: execEX9E, 0xE0A1: execEXA1, 0xF007: execFX07, 0xF00A: execFX0A,
		0xF015: execFX15, 0xF018: execFX18, 0xF01E: execFX1E, 0xF029: execFX29,
		0xF033: execFX33, 0xF055: execFX55, 0xF065: execFX65,
	}
)

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
	loaded = false
}

func main() {
	ebiten.SetFullscreen(fullscreen)
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("CHIP-8 by TheInvader360")
	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}

func NewGame() *Game {
	g := &Game{}
	return g
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func (g *Game) Update(screen *ebiten.Image) error {
	if loaded == false {
		loadRom()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	emulateCycle()
	updateKeys()
	return nil
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
	//fmt.Println(memory)
	loaded = true
}

func emulateCycle() {
	opcode = fetchOpcode()
	decoded := decodeOpcode()
	fmt.Println(opcodeExecutors[decoded]())
	updateTimers()
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

func fetchOpcode() uint16 {
	/*
		Fetch opcode:
		Fetch and merge two bytes from memory locations pointed at by pc & pc+1
		e.g. memory[pc] = 0b10100010, memory[pc+1] = 0b11110000
		Convert first byte to uint16 and shift the bits left 8 times.
		e.g. 0b1010001000000000
		Use bitwise OR operation to merge the bytes.
		e.g. 0b1010001000000000 | 0b11110000 = 0b1010001011110000
	*/
	return uint16(memory[pc])<<8 | uint16(memory[pc+1])
}

func decodeOpcode() uint16 {
	/*
		Decode opcode:
		Read the first 4 bits of the current opcode using bitwise AND operation
		e.g. 0x2105 & 0xF000 = 0x2000
		We can't always rely on just the first nibble to decode opcodes
		e.g. 0x00E0 and 0x00EE both start with 0x0
		In these cases we go on to compare the last nibble or byte...
		e.g. 0x00EE & 0x00FF = 0x00EE
	*/
	decoded := opcode & 0xF000
	if decoded == 0x0000 {
		switch opcode & 0x00FF {
		case 0x00E0:
			decoded = 0x00E0
		case 0x00EE:
			decoded = 0x00EE
		}
	}
	if decoded == 0x8000 {
		switch opcode & 0x000F {
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
		switch opcode & 0x00FF {
		case 0x009E:
			decoded = 0xE09E
		case 0x00A1:
			decoded = 0xE0A1
		}
	}
	if decoded == 0xF000 {
		switch opcode & 0x00FF {
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

func updateTimers() {
	//Count down to zero
	if delayTimer > 0 {
		delayTimer--
	}
	if soundTimer > 0 {
		soundTimer--
	}
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func exec0NNN() string {
	//TODO ?
	return fmt.Sprintf("exec0NNN 0x%04X", opcode)
}

func exec00E0() string {
	//disp_clear()
	gfx = [2048]byte{}
	pc += 2
	render = true
	return fmt.Sprintf("exec00E0 0x%04X", opcode)
}

func exec00EE() string {
	//TODO return;
	return fmt.Sprintf("exec00EE 0x%04X", opcode)
}

func exec1NNN() string {
	//TODO goto NNN;
	return fmt.Sprintf("exec1NNN 0x%04X", opcode)
}

func exec2NNN() string {
	//TODO *(0xNNN)()
	return fmt.Sprintf("exec2NNN 0x%04X", opcode)
}

func exec3XNN() string {
	//TODO if(Vx==NN)
	return fmt.Sprintf("exec3XNN 0x%04X", opcode)
}

func exec4XNN() string {
	//TODO if(Vx!=NN)
	return fmt.Sprintf("exec4XNN 0x%04X", opcode)
}

func exec5XY0() string {
	//TODO if(Vx==Vy)
	return fmt.Sprintf("exec5XY0 0x%04X", opcode)
}

func exec6XNN() string {
	//TODO Vx=NN
	return fmt.Sprintf("exec6XNN 0x%04X", opcode)
}

func exec7XNN() string {
	//TODO Vx+=NN
	return fmt.Sprintf("exec7XNN 0x%04X", opcode)
}

func exec8XY0() string {
	//TODO Vx=Vy
	return fmt.Sprintf("exec8XY0 0x%04X", opcode)
}

func exec8XY1() string {
	//TODO Vx=Vx|Vy
	return fmt.Sprintf("exec8XY1 0x%04X", opcode)
}

func exec8XY2() string {
	//TODO Vx=Vx&Vy
	return fmt.Sprintf("exec8XY2 0x%04X", opcode)
}

func exec8XY3() string {
	//TODO Vx=Vx^Vy
	return fmt.Sprintf("exec8XY3 0x%04X", opcode)
}

func exec8XY4() string {
	//TODO Vx+=Vy
	return fmt.Sprintf("exec8XY4 0x%04X", opcode)
}

func exec8XY5() string {
	//TODO Vx-=Vy
	return fmt.Sprintf("exec8XY5 0x%04X", opcode)
}

func exec8XY6() string {
	//TODO Vx>>=1
	return fmt.Sprintf("exec8XY6 0x%04X", opcode)
}

func exec8XY7() string {
	//TODO Vx=Vy-Vx
	return fmt.Sprintf("exec8XY7 0x%04X", opcode)
}

func exec8XYE() string {
	//TODO Vx<<=1
	return fmt.Sprintf("exec8XYE 0x%04X", opcode)
}

func exec9XY0() string {
	//TODO if(Vx!=Vy)
	return fmt.Sprintf("exec9XY0 0x%04X", opcode)
}

func execANNN() string {
	//TODO I=NNN
	return fmt.Sprintf("execANNN 0x%04X", opcode)
}

func execBNNN() string {
	//TODO PC=V0+NNN
	return fmt.Sprintf("execBNNN 0x%04X", opcode)
}

func execCXNN() string {
	//TODO Vx=rand()&NN
	return fmt.Sprintf("execCXNN 0x%04X", opcode)
}

func execDXYN() string {
	//TODO draw(Vx,Vy,N)
	return fmt.Sprintf("execDXYN 0x%04X", opcode)
}

func execEX9E() string {
	//TODO if(key()==Vx)
	return fmt.Sprintf("execEX9E 0x%04X", opcode)
}

func execEXA1() string {
	//TODO if(key()!=Vx)
	return fmt.Sprintf("execEXA1 0x%04X", opcode)
}

func execFX07() string {
	//TODO Vx=get_delay()
	return fmt.Sprintf("execFX07 0x%04X", opcode)
}

func execFX0A() string {
	//TODO Vx=get_key()
	return fmt.Sprintf("execFX0A 0x%04X", opcode)
}

func execFX15() string {
	//TODO delay_timer(Vx)
	return fmt.Sprintf("execFX15 0x%04X", opcode)
}

func execFX18() string {
	//TODO sound_timer(Vx)
	return fmt.Sprintf("execFX18 0x%04X", opcode)
}

func execFX1E() string {
	//TODO I+=Vx
	return fmt.Sprintf("execFX1E 0x%04X", opcode)
}

func execFX29() string {
	//TODO I=sprite_addr[Vx]
	return fmt.Sprintf("execFX29 0x%04X", opcode)
}

func execFX33() string {
	//TODO set_BCD(Vx);*(I+0)=BCD(3);*(I+1)=BCD(2);*(I+2)=BCD(1);
	return fmt.Sprintf("execFX33 0x%04X", opcode)
}

func execFX55() string {
	//TODO reg_dump(Vx,&I)
	return fmt.Sprintf("execFX55 0x%04X", opcode)
}

func execFX65() string {
	//TODO reg_load(Vx,&I)
	return fmt.Sprintf("execFX65 0x%04X", opcode)
}
