package main

import (
	"github.com/TheInvader360/chip8/client"
	"github.com/hajimehoshi/ebiten"
)

func main() {
	ebiten.SetFullscreen(client.ScrF)
	ebiten.SetWindowSize(client.ScrW, client.ScrH)
	ebiten.SetWindowTitle("CHIP-8 by TheInvader360")
	if err := ebiten.RunGame(client.NewGame()); err != nil {
		panic(err)
	}
}
