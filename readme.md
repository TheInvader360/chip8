# CHIP-8 [![Build Status](https://travis-ci.com/TheInvader360/chip8.svg?branch=master)](https://travis-ci.com/TheInvader360/chip8)

[CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) emulator written in [Go](https://golang.org/).


## Local Setup

[Download](https://golang.org/dl) and [install](https://golang.org/doc/install) Go, then install dependencies:

    sudo apt install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

Clone and run:

    git clone https://github.com/TheInvader360/chip8
    cd chip8
    go test ./... && go run main.go -path=./rom/test/ti360.ch8

Command line options:

    -path string : path to rom file (default "./rom/test/ti360.ch8")
    -clock int   : cpu clock speed in hz (100-1000) (default 400)
    -debug       : enable debug info in terminal
    -fullscreen  : enable fullscreen mode
    -height int  : height of client screen in pixels (default 320)
    -width int   : width of client screen in pixels (default 640)

Example launch commands:

    go run main.go
    go run main.go -path=./rom/BRIX.ch8
    go run main.go -path=./rom/BLINKY.ch8 -clock=1000
    go run main.go -path=./rom/CAVE.ch8 -clock=300 -debug
    go run main.go -path=./rom/PONG.ch8 -clock=500 -fullscreen
    go run main.go -path=./rom/UFO.ch8 -width=1280 -height=640
    go run main.go -path=./rom/WALL.ch8 -clock=300 -debug -width=640 -height=480

## Input

Mapping CHIP-8 keys to a QWERTY keyboard:

    CHIP-8:     QWERTY:
    1 2 3 C     1 2 3 4
    4 5 6 D     Q W E R
    7 8 9 E     A S D F
    A 0 B F     Z X C V

Additional controls:

    P   - Pause or unpause emulation loop
    O   - One iteration of the emulation loop (step through while paused)
    I   - Initialize (reset to a clean starting state)
    Esc - Exit

