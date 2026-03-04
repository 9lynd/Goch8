package main

import (
	"fmt"
	"os"

	"github.com/9lynd/Goch8/core"
	ebitenrenderer "github.com/9lynd/Goch8/renderer/ebiten"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: chip8 <rom path>")
		os.Exit(1)
	}

	romPath := os.Args[1]
	romData, err := os.ReadFile(romPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read ROM: %v\n", err)
		os.Exit(1)
	}

	display := &core.Display{}
	keyboard := &core.Keyboard{}
	cpu := core.NewCPU(display, keyboard)

	if err := cpu.LoadROM(romData); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load ROM: %v\n", err)
		os.Exit(1)
	}

	renderer := &ebitenrenderer.Renderer{}
	if err := renderer.Run(cpu); err != nil {
		fmt.Fprintf(os.Stderr, "renderer error: %v\n", err)
		os.Exit(1)
	}
}
