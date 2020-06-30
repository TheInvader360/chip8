# CHIP-8 [![Build Status](https://travis-ci.com/TheInvader360/chip8.svg?branch=master)](https://travis-ci.com/TheInvader360/chip8)

[CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) emulator written in [Go](https://golang.org/).


## Local Setup

[Download](https://golang.org/dl) and [install](https://golang.org/doc/install) Go, then install dependencies:

    sudo apt install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

Clone and run:

    git clone https://github.com/TheInvader360/chip8
    cd chip8
    go test ./... && go run main.go -path=./rom/test/ti360.ch8

Various public domain ROMs are included. Examples:

    go run main.go -path=./rom/BLINKY.ch8 -clock=700
    go run main.go -path=./rom/BRIX.ch8
    go run main.go -path=./rom/CAVE.ch8 -clock=300
    go run main.go -path=./rom/PONG.ch8 -clock=500

The CPU clock defaults to 400Hz. Set the optional -clock argument to override.

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

