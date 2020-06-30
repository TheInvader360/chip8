package client

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/TheInvader360/chip8/lib"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	tps = 60 //max ticks per second
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
	counter float64       //cpu counter
	rom     = false       //is rom loaded into memory
	paused  = false       //is emulation loop paused
	view    *ebiten.Image //screen buffer
)

type Game struct {
	opts *Opts
	vm   *lib.Chip8
}

func init() {
	ebiten.SetMaxTPS(tps)
	rand.Seed(time.Now().UnixNano())
	view, _ = ebiten.NewImage(lib.GfxW, lib.GfxH, ebiten.FilterDefault)
}

func NewGame(opts *Opts) *Game {
	return &Game{
		opts: opts,
		vm:   lib.NewChip8(),
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !rom {
		loadRom(g)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		paused = !paused
	}
	if paused && inpututil.IsKeyJustPressed(ebiten.KeyO) {
		step(g)
	}
	if !paused {
		for counter > 0 {
			step(g)
			counter -= tps
		}
		counter += float64(g.opts.Clock)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.vm = lib.NewChip8()
		rom = false
	}
	g.vm.UpdateTimers()
	//play sound (dependent on sound timer)
	Noise(g.vm.St > 0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//only draw vm.gfx to view (buffer) if vm.rg flag is set
	if g.vm.Rg {
		view.Fill(bg)
		for y := 0; y < lib.GfxH; y++ {
			for x := 0; x < lib.GfxW; x++ {
				if g.vm.Gfx[y*lib.GfxW+x] == 1 {
					ebitenutil.DrawRect(view, float64(x), float64(y), 1, 1, fg)
				}
			}
		}
		g.vm.Rg = false
	}
	//draw to screen (ebiten clears the screen each frame)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(view, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return lib.GfxW, lib.GfxH
}

func loadRom(g *Game) {
	file, err := os.Open(g.opts.Path)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(bytes); i++ {
		g.vm.Mem[i+0x200] = bytes[i]
	}
	if g.opts.Debug {
		fmt.Println(g.vm.Mem)
	}

	rom = true
}

func step(g *Game) {
	g.vm.EmulateCycle()
	if g.opts.Debug {
		fmt.Println(g.vm.DebugInfo())
	}
	updateKeys(g.vm)
}

func updateKeys(vm *lib.Chip8) {
	for phys, virt := range keyMap {
		vm.Key[virt] = lib.BoolToByte(ebiten.IsKeyPressed(phys))
	}
}
