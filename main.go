package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/9lynd/Goch8/core"
)

const (
	Width  = 64
	Height = 32
)

func render(d *core.Display) {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if d.Pixels[y*Width+x] {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func anyPixel(d *core.Display) bool {
	for _, p := range d.Pixels {
		if p {
			return true
		}
	}
	return false
}

func main() {
	display := &core.Display{}
	keyboard := &core.Keyboard{}
	cpu := core.NewCPU(display, keyboard)

	romData, err := ioutil.ReadFile("roms/IBM_Logo.ch8")
	if err != nil {
		log.Fatalf("Failed to load ROM: %v", err)
	}

	err = cpu.LoadROM(romData)
	if err != nil {
		log.Fatalf("Error loading ROM: %v", err)
	}

	const maxCycles = 2000
	for i := 0; i < maxCycles; i++ {
		if err := cpu.Cycle(); err != nil {
			log.Printf("Stopped: %v", err)
			break
		}
		if cpu.Display.Dirty {
			if anyPixel(cpu.Display) {
				render(cpu.Display)
				break
			}
			cpu.Display.Dirty = false
		}
	}

	fmt.Println("ROM bytes at 0x200:")
	for i := 0; i < 16; i++ {
		fmt.Printf("0x%03X: 0x%02X\n", 0x200+i, cpu.Memory[0x200+i])
	}
}
