package main

import (
	"flag"

	"github.com/TheInvader360/chip8/client"
	"github.com/TheInvader360/chip8/lib"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	path := flag.String("path", "./rom/test/ti360.ch8", "path to rom file")
	clock := flag.Int("clock", 400, "cpu clock speed in hz (100-1000)")
	screenW := flag.Int("width", 640, "width of client screen in pixels")
	screenH := flag.Int("height", 320, "height of client screen in pixels")
	fullscreen := flag.Bool("fullscreen", false, "enable fullscreen mode")
	debug := flag.Bool("debug", false, "enable debug info in terminal")
	flag.Parse()

	opts := client.Opts{
		Path:  *path,
		Clock: lib.Clamp(*clock, 100, 1000),
		ScrW:  *screenW,
		ScrH:  *screenH,
		ScrF:  *fullscreen,
		Debug: *debug,
	}

	ebiten.SetFullscreen(opts.ScrF)
	ebiten.SetWindowSize(opts.ScrW, opts.ScrH)
	ebiten.SetWindowTitle("CHIP-8 by TheInvader360")
	if err := ebiten.RunGame(client.NewGame(&opts)); err != nil {
		panic(err)
	}
}
