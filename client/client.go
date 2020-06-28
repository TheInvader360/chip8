package client

import (
	"flag"
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
	ScrF = false //client fullscreen
	ScrW = 960   //client screen width in pixels
	ScrH = 480   //client screen height in pixels

	tps  = 60                       //client max ticks per second
	pixW = float64(ScrW / lib.GfxW) //pixel width scale factor
	pixH = float64(ScrH / lib.GfxH) //pixel height scale factor
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
	path   *string
	rom    = false
	paused = false
	view   *ebiten.Image
)

type Game struct {
	vm *lib.Chip8
}

func init() {
	ebiten.SetMaxTPS(tps)
	rand.Seed(time.Now().UnixNano())
	view, _ = ebiten.NewImage(ScrW, ScrH, ebiten.FilterDefault)
}

func NewGame() *Game {
	return &Game{
		vm: lib.NewChip8(),
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	if !rom {
		loadRom(g.vm)
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		paused = !paused
	}
	if !paused || inpututil.IsKeyJustPressed(ebiten.KeyO) {
		g.vm.EmulateCycle()
		fmt.Println(g.vm.DebugInfo())
		updateKeys(g.vm)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.vm = lib.NewChip8()
		rom = false
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//only draw vm.gfx to view (buffer) if vm.rg flag is set
	if g.vm.Rg {
		view.Fill(bg)
		for y := 0; y < lib.GfxH; y++ {
			for x := 0; x < lib.GfxW; x++ {
				if g.vm.Gfx[y*lib.GfxW+x] == 1 {
					ebitenutil.DrawRect(view,
						float64(x)*pixW, float64(y)*pixH, pixW, pixH, fg)
				}
			}
		}
		g.vm.Rg = false
	}
	//draw to screen (ebiten clears the screen each frame)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(view, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScrW, ScrH
}

func loadRom(vm *lib.Chip8) {
	if path == nil {
		path = flag.String("path", "./rom/test/ti360.ch8", "path to rom file")
		flag.Parse()
	}

	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(bytes); i++ {
		vm.Mem[i+0x200] = bytes[i]
	}
	fmt.Println(vm.Mem)

	rom = true
}

func updateKeys(vm *lib.Chip8) {
	for phys, virt := range keyMap {
		vm.Key[virt] = lib.BoolToByte(ebiten.IsKeyPressed(phys))
	}
}
